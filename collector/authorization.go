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
	AuthorizationSubsystem = "authorization"
)

const (
	authorizationResStatus      = "resource_status"
	authorizationTotal          = "total"
	authorizationAllowCount     = "allow_count"
	authorizationDenyCount      = "deny_count"
	authorizationExecRate       = "exec_rate"
	authorizationExecLast5mRate = "exec_last5m_rate"
	authorizationExecMaxRate    = "exec_max_rate"
	authorizationExecTimeCost   = "exec_time_cost"
)

func init() {
	registerCollector(AuthorizationSubsystem, NewAuthorizationCollector)
}

type AuthorizationCollector struct {
	desc    map[string]*prometheus.Desc
	logger  log.Logger
	cluster Cluster
}

// NewAuthorizationCollector returns a new authorization collector
func NewAuthorizationCollector(cluster Cluster, logger log.Logger) (Collector, error) {
	collector := &AuthorizationCollector{
		desc:    make(map[string]*prometheus.Desc),
		logger:  logger,
		cluster: cluster,
	}

	metrics := []struct {
		name   string
		help   string
		labels []string
	}{
		{
			name:   authorizationResStatus,
			help:   "The status of authorization resource",
			labels: []string{"resource"},
		},
		{
			name:   authorizationTotal,
			help:   "The total of authorization",
			labels: []string{"node", "resource"},
		},
		{
			name:   authorizationAllowCount,
			help:   "The count of allowable authorization",
			labels: []string{"node", "resource"},
		},
		{
			name:   authorizationDenyCount,
			help:   "The count of denied authorization",
			labels: []string{"node", "resource"},
		},
		{
			name:   authorizationExecRate,
			help:   "The rate of authorization exec",
			labels: []string{"node", "resource"},
		},
		{
			name:   authorizationExecLast5mRate,
			help:   "The last 5m average rate of authorization exec",
			labels: []string{"node", "resource"},
		},
		{
			name:   authorizationExecMaxRate,
			help:   "The max rate of authorization exec",
			labels: []string{"node", "resource"},
		},
		{
			name:   authorizationExecTimeCost,
			help:   "The time cost of authorization exec",
			labels: []string{"node", "resource"},
		},
	}

	for _, m := range metrics {
		collector.desc[m.name] = prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				AuthorizationSubsystem,
				m.name,
			),
			m.help,
			m.labels,
			nil,
		)
	}
	return collector, nil
}

// Update implements the Collector interface and will collect authorization metrics.
func (c *AuthorizationCollector) Update(ch chan<- prometheus.Metric) error {
	dataSources, metrics, err := c.cluster.GetAuthorizationMetrics()
	if err != nil {
		return err
	}

	for i := range dataSources {
		ds := &dataSources[i]
		ch <- prometheus.MustNewConstMetric(
			c.desc[authorizationResStatus],
			prometheus.GaugeValue, float64(ds.Status), ds.ResType,
		)
	}

	for i := range metrics {
		metric := &metrics[i]
		bucket, err := getBucket(metric.ExecTimeCost)
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstHistogram(c.desc[authorizationExecTimeCost],
			metric.ExecTimeCost["count"],
			float64(metric.ExecTimeCost["sum"]),
			bucket, metric.NodeName, metric.ResType)

		ch <- prometheus.MustNewConstMetric(
			c.desc[authorizationTotal],
			prometheus.CounterValue, float64(metric.Total), metric.NodeName, metric.ResType,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[authorizationAllowCount],
			prometheus.CounterValue, float64(metric.AllowCount), metric.NodeName, metric.ResType,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[authorizationDenyCount],
			prometheus.CounterValue, float64(metric.DenyCount), metric.NodeName, metric.ResType,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[authorizationExecRate],
			prometheus.GaugeValue, metric.ExecRate, metric.NodeName, metric.ResType,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[authorizationExecLast5mRate],
			prometheus.GaugeValue, metric.ExecLast5mRate, metric.NodeName, metric.ResType,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[authorizationExecMaxRate],
			prometheus.GaugeValue, metric.ExecMaxRate, metric.NodeName, metric.ResType,
		)

	}
	return nil
}
