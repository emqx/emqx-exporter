package prober

import (
	"emqx-exporter/config"

	"fmt"

	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Handler(w http.ResponseWriter, r *http.Request, probes []config.Probe, logger log.Logger, params url.Values) {
	var probe config.Probe
	if params == nil {
		params = r.URL.Query()
	}
	target := params.Get("target")
	for i := 0; i < len(probes); i++ {
		if probes[i].Target == target {
			probe = probes[i]
			break
		}
	}
	if probe.Target == "" {
		http.Error(w, fmt.Sprintf("Unknown probe target %q", target), http.StatusBadRequest)
		level.Debug(logger).Log("msg", "Unknown probe target", "target", target)
		return
	}

	probeSuccessGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "emqx",
		Subsystem: "mqtt",
		Name:      "probe_success",
		Help:      "Displays whether or not the probe was a success",
		ConstLabels: prometheus.Labels{
			"target": probe.Target,
		},
	})
	probeDurationGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "emqx",
		Subsystem: "mqtt",
		Name:      "probe_duration_seconds",
		Help:      "Returns how long the probe took to complete in seconds",
		ConstLabels: prometheus.Labels{
			"target": probe.Target,
		},
	})

	registry := prometheus.NewRegistry()
	registry.MustRegister(probeSuccessGauge)
	registry.MustRegister(probeDurationGauge)

	start := time.Now()
	mp := newMQTTProbe(probe, logger)
	if mp != nil && mp.Probe(probe, logger) {
		probeSuccessGauge.Set(1)
	} else {
		probeSuccessGauge.Set(0)
	}

	probeDurationGauge.Set(time.Since(start).Seconds())
	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}
