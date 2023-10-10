package prober

import (
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Handler(w http.ResponseWriter, r *http.Request, logger log.Logger) {
	probeSuccessGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "emqx",
		Subsystem: "mqtt",
		Name:      "probe_success",
		Help:      "Displays whether or not the probe was a success",
	})
	probeDurationGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "emqx",
		Subsystem: "mqtt",
		Name:      "probe_duration_seconds",
		Help:      "Returns how long the probe took to complete in seconds",
	})

	registry := prometheus.NewRegistry()
	registry.MustRegister(probeSuccessGauge)
	registry.MustRegister(probeDurationGauge)

	start := time.Now()
	if ProbeMQTT(logger) {
		probeSuccessGauge.Set(1)
	} else {
		probeSuccessGauge.Set(0)
	}
	probeDurationGauge.Set(time.Since(start).Seconds())

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
