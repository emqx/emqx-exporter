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
	AuthenticationSubsystem = "authentication"
)

const (
	authenticationResStatus      = "resource_status"
	authenticationTotal          = "total"
	authenticationAllowCount     = "allow_count"
	authenticationDenyCount      = "deny_count"
	authenticationExecRate       = "exec_rate"
	authenticationExecLast5mRate = "exec_last5m_rate"
	authenticationExecMaxRate    = "exec_max_rate"
	authenticationExecTimeCost   = "exec_time_cost"
)

func init() {
	registerCollector(AuthenticationSubsystem, NewAuthenticationCollector)
}

type authenticationCollector struct {
	desc    map[string]*prometheus.Desc
	logger  log.Logger
	cluster Cluster
}

// NewAuthenticationCollector returns a new authentication collector
func NewAuthenticationCollector(client Cluster, logger log.Logger) (Collector, error) {
	collector := &authenticationCollector{
		desc:    make(map[string]*prometheus.Desc),
		logger:  logger,
		cluster: client,
	}

	metrics := []struct {
		name   string
		help   string
		labels []string
	}{
		{
			name:   authenticationResStatus,
			help:   "The status of authentication resource",
			labels: []string{"resource"},
		},
		{
			name:   authenticationTotal,
			help:   "The total of authentication",
			labels: []string{"node", "resource"},
		},
		{
			name:   authenticationAllowCount,
			help:   "The count of allowable authentication",
			labels: []string{"node", "resource"},
		},
		{
			name:   authenticationDenyCount,
			help:   "The count of denied authentication",
			labels: []string{"node", "resource"},
		},
		{
			name:   authenticationExecRate,
			help:   "The rate of authentication exec",
			labels: []string{"node", "resource"},
		},
		{
			name:   authenticationExecLast5mRate,
			help:   "The last 5m average rate of authentication exec",
			labels: []string{"node", "resource"},
		},
		{
			name:   authenticationExecMaxRate,
			help:   "The max rate of authentication exec",
			labels: []string{"node", "resource"},
		},
		{
			name:   authenticationExecTimeCost,
			help:   "The time cost of authentication exec",
			labels: []string{"node", "resource"},
		},
	}

	for _, m := range metrics {
		collector.desc[m.name] = prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				AuthenticationSubsystem,
				m.name,
			),
			m.help,
			m.labels,
			nil,
		)
	}
	return collector, nil
}

// Update implements the Collector interface and will collect authentication metrics.
func (c *authenticationCollector) Update(ch chan<- prometheus.Metric) error {
	dataSources, metrics, err := c.cluster.GetAuthenticationMetrics()
	if err != nil {
		return err
	}

	for i := range dataSources {
		ds := &dataSources[i]
		ch <- prometheus.MustNewConstMetric(
			c.desc[authenticationResStatus],
			prometheus.GaugeValue, float64(ds.Status), ds.ResType,
		)
	}

	for i := range metrics {
		metric := &metrics[i]
		bucket, err := getBucket(metric.ExecTimeCost)
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstHistogram(c.desc[authenticationExecTimeCost],
			metric.ExecTimeCost["count"],
			float64(metric.ExecTimeCost["sum"]),
			bucket, metric.NodeName, metric.ResType)

		ch <- prometheus.MustNewConstMetric(
			c.desc[authenticationTotal],
			prometheus.CounterValue, float64(metric.Total), metric.NodeName, metric.ResType,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[authenticationAllowCount],
			prometheus.CounterValue, float64(metric.AllowCount), metric.NodeName, metric.ResType,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[authenticationDenyCount],
			prometheus.CounterValue, float64(metric.DenyCount), metric.NodeName, metric.ResType,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[authenticationExecRate],
			prometheus.GaugeValue, metric.ExecRate, metric.NodeName, metric.ResType,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[authenticationExecLast5mRate],
			prometheus.GaugeValue, metric.ExecLast5mRate, metric.NodeName, metric.ResType,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[authenticationExecMaxRate],
			prometheus.GaugeValue, metric.ExecMaxRate, metric.NodeName, metric.ResType,
		)

	}
	return nil
}
