package client

import (
	"emqx-exporter/collector"
	"emqx-exporter/config"

	stdlog "log"
	"net/http"
	"sort"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

func NewHandler(disableExporterMetrics bool, maxRequests int, metrics *config.Metrics, logger log.Logger) http.Handler {
	registry := prometheus.NewRegistry()
	registry.MustRegister(version.NewCollector("emqx_exporter"))

	if metrics == nil {
		level.Info(logger).Log("msg", "No metrics configured, skipping cluster metrics")
	} else {
		emqxCluster := NewCluster(metrics, logger)
		nc, err := collector.NewEMQXCollector(emqxCluster, logger)
		if err != nil {
			level.Debug(logger).Log("msg", "Couldn't create collector", "err", err)
			panic("Couldn't create collector")
		}

		level.Info(logger).Log("msg", "Enabled collectors")
		collectors := make([]string, 0, len(nc.Collectors))
		for n := range nc.Collectors {
			collectors = append(collectors, n)
		}
		sort.Strings(collectors)
		for _, c := range collectors {
			level.Info(logger).Log("collector", c)
		}
		registry.MustRegister(nc)
	}

	var h http.Handler
	if disableExporterMetrics {
		level.Info(logger).Log("msg", "Excluding metrics about the exporter itself")
		h = promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{
				ErrorLog:            stdlog.New(log.NewStdlibAdapter(level.Error(logger)), "", 0),
				ErrorHandling:       promhttp.ContinueOnError,
				MaxRequestsInFlight: maxRequests,
			})
	} else {
		level.Info(logger).Log("msg", "Including metrics about the exporter itself")
		exporterMetricsRegistry := prometheus.NewRegistry()
		exporterMetricsRegistry.MustRegister(
			promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
			promcollectors.NewGoCollector(),
		)

		h = promhttp.HandlerFor(
			prometheus.Gatherers{exporterMetricsRegistry, registry},
			promhttp.HandlerOpts{
				ErrorLog:            stdlog.New(log.NewStdlibAdapter(level.Error(logger)), "", 0),
				ErrorHandling:       promhttp.ContinueOnError,
				MaxRequestsInFlight: maxRequests,
				Registry:            exporterMetricsRegistry,
			})
		h = promhttp.InstrumentMetricHandler(
			exporterMetricsRegistry, h,
		)
	}

	return h
}
