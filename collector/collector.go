// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package collector includes all individual collectors to gather and export emqx metrics.
package collector

import (
	"errors"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// Namespace defines the common namespace to be used by all metrics.
const namespace = "emqx"

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
		"emqx-exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_success"),
		"emqx-exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)
)

var (
	factories = make(map[string]func(client *client) (Collector, error))
)

func registerCollector(collector string, factory func(client *client) (Collector, error)) {
	factories[collector] = factory
}

// EMQXCollector implements the prometheus.Collector interface.
type EMQXCollector struct {
	Collectors map[string]Collector
	logger     log.Logger
}

// NewEMQXCollector creates a new EMQXCollector.
func NewEMQXCollector(client *client, logger log.Logger) (*EMQXCollector, error) {
	collectors := make(map[string]Collector)
	for key, factory := range factories {
		collector, err := factory(client)
		if err != nil {
			return nil, err
		}
		collectors[key] = collector
	}
	return &EMQXCollector{Collectors: collectors, logger: logger}, nil
}

// Describe implements the prometheus.Collector interface.
func (n EMQXCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

// Collect implements the prometheus.Collector interface.
func (n EMQXCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	for name, c := range n.Collectors {
		wg.Add(1)
		go func(name string, c Collector) {
			defer wg.Done()
			execute(name, c, ch, n.logger)
		}(name, c)
	}
}

func execute(name string, c Collector, ch chan<- prometheus.Metric, logger log.Logger) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		if IsNoDataError(err) {
			level.Debug(logger).Log("msg", "collector returned no data", "name", name, "duration_seconds", duration.Seconds(), "err", err)
		} else {
			level.Error(logger).Log("msg", "collector failed", "name", name, "duration_seconds", duration.Seconds(), "err", err)
		}
		success = 0
	} else {
		level.Debug(logger).Log("msg", "collector succeeded", "name", name, "duration_seconds", duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(ch chan<- prometheus.Metric) error
}

// ErrNoData indicates the collector found no data to collect, but had no other error.
var ErrNoData = errors.New("collector returned no data")

func IsNoDataError(err error) bool {
	return err == ErrNoData
}
