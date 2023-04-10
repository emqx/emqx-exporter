// Copyright 2019 The Prometheus Authors
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

package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	BrokerSubsystem = "messages"
)

const (
	consumeTimeCost = "consume_time_cost"
	inputPeriodSec  = "input_period_second"
	outputPeriodSec = "output_period_second"
)

func init() {
	registerCollector(BrokerSubsystem, NewBrokerCollector)
}

type brokerCollector struct {
	desc    map[string]*prometheus.Desc
	logger  log.Logger
	cluster Cluster
}

// NewBrokerCollector returns a new broker msg collector
func NewBrokerCollector(cluster Cluster, logger log.Logger) (Collector, error) {
	collector := &brokerCollector{
		desc:    make(map[string]*prometheus.Desc),
		logger:  logger,
		cluster: cluster,
	}

	metrics := []struct {
		name string
		help string
	}{
		{
			name: consumeTimeCost,
			help: "The time cost of msg consumed",
		},
		{
			name: inputPeriodSec,
			help: "The input msg period second",
		},
		{
			name: outputPeriodSec,
			help: "The output msg period second",
		},
	}

	for _, m := range metrics {
		collector.desc[m.name] = prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				BrokerSubsystem,
				m.name,
			),
			m.help,
			nil,
			nil,
		)
	}
	return collector, nil
}

// Update implements the Collector interface and will collect license info.
func (c *brokerCollector) Update(ch chan<- prometheus.Metric) error {
	metrics, err := c.cluster.GetBrokerMetrics()
	if err != nil {
		return err
	}
	if metrics == nil {
		return nil
	}

	bucket, err := getBucket(metrics.MsgConsumeTimeCosts)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstHistogram(c.desc[consumeTimeCost],
		metrics.MsgConsumeTimeCosts["count"],
		float64(metrics.MsgConsumeTimeCosts["sum"]),
		bucket)

	ch <- prometheus.MustNewConstMetric(
		c.desc[inputPeriodSec],
		prometheus.GaugeValue, float64(metrics.MsgInputPeriodSec),
	)
	ch <- prometheus.MustNewConstMetric(
		c.desc[outputPeriodSec],
		prometheus.GaugeValue, float64(metrics.MsgOutputPeriodSec),
	)
	return nil
}
