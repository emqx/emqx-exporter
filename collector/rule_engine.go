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
	RuleEngineSubsystem = "rule"
)

const (
	bridgeResStatus  = "bridge_status"
	bridgeQueuing    = "bridge_queuing"
	bridgeLast5mRate = "bridge_last5m_rate"
	bridgeRateMax    = "bridge_max_rate"
	bridgeFailed     = "bridge_failed"
	bridgeDropped    = "bridge_dropped"

	ruleTopicHitCount    = "topic_hit_count"
	ruleExecPassCount    = "exec_pass_count"
	ruleExecFailureCount = "exec_failure_count"
	ruleNoResultCount    = "exec_no_result_count"
	ruleExecRate         = "exec_rate"
	ruleExecLast5mRate   = "exec_last5m_rate"
	ruleExecMaxRate      = "exec_max_rate"
	ruleActionTotal      = "action_total"
	ruleActionSuccess    = "action_success"
	ruleActionFailed     = "action_failed"
	ruleExecTimeCost     = "exec_time_cost"
)

func init() {
	registerCollector(RuleEngineSubsystem, NewRuleEngineCollector)
}

type ruleEngineCollector struct {
	desc   map[string]*prometheus.Desc
	logger log.Logger
	client Cluster
}

// NewRuleEngineCollector returns a new rule engine collector
func NewRuleEngineCollector(client Cluster, logger log.Logger) (Collector, error) {
	collector := &ruleEngineCollector{
		desc:   make(map[string]*prometheus.Desc),
		logger: logger,
		client: client,
	}

	metrics := []struct {
		name   string
		help   string
		labels []string
	}{
		{
			name:   bridgeResStatus,
			help:   "The status of rule engine resource",
			labels: []string{"type", "name"},
		},
		{
			name:   bridgeQueuing,
			help:   "The count of messages that are currently queuing",
			labels: []string{"type", "name"},
		},
		{
			name:   bridgeLast5mRate,
			help:   "The last 5m average rate of rule engine resource",
			labels: []string{"type", "name"},
		},
		{
			name:   bridgeRateMax,
			help:   "The max rate of rule engine resource",
			labels: []string{"type", "name"},
		},
		{
			name:   bridgeFailed,
			help:   "The failure messages count of rule engine resource",
			labels: []string{"type", "name"},
		},
		{
			name:   bridgeDropped,
			help:   "The dropped messages count of rule engine resource",
			labels: []string{"type", "name"},
		},
		{
			name:   ruleTopicHitCount,
			help:   "The count of topic hit",
			labels: []string{"node", "rule"},
		},
		{
			name:   ruleExecPassCount,
			help:   "The pass count of rule exec",
			labels: []string{"node", "rule"},
		},
		{
			name:   ruleExecFailureCount,
			help:   "The failure count of rule exec",
			labels: []string{"node", "rule"},
		},
		{
			name:   ruleNoResultCount,
			help:   "The no result count of rule exec",
			labels: []string{"node", "rule"},
		},
		{
			name:   ruleExecRate,
			help:   "The current rate of rule exec",
			labels: []string{"node", "rule"},
		},
		{
			name:   ruleExecLast5mRate,
			help:   "The last 5m average rate of rule exec",
			labels: []string{"node", "rule"},
		},
		{
			name:   ruleExecMaxRate,
			help:   "The max rate of rule exec",
			labels: []string{"node", "rule"},
		},
		{
			name:   ruleActionTotal,
			help:   "The total of rule action exec",
			labels: []string{"node", "rule"},
		},
		{
			name:   ruleActionSuccess,
			help:   "The success count of rule action exec",
			labels: []string{"node", "rule"},
		},
		{
			name:   ruleActionFailed,
			help:   "The failure count of rule action exec",
			labels: []string{"node", "rule"},
		},
		{
			name:   ruleExecTimeCost,
			help:   "The time cost of rule exec",
			labels: []string{"node", "rule"},
		},
	}

	for _, m := range metrics {
		collector.desc[m.name] = prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				RuleEngineSubsystem,
				m.name,
			),
			m.help,
			m.labels,
			nil,
		)
	}
	return collector, nil
}

// Update implements the Collector interface and will collect rule engine metrics.
func (c *ruleEngineCollector) Update(ch chan<- prometheus.Metric) error {
	bridges, metrics, err := c.client.GetRuleEngineMetrics()
	if err != nil {
		return err
	}

	for i := range metrics {
		metric := &metrics[i]
		bucket, err := getBucket(metric.ActionExecTimeCost)
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstHistogram(c.desc[ruleExecTimeCost],
			metric.ActionExecTimeCost["count"],
			float64(metric.ActionExecTimeCost["sum"]),
			bucket, metric.NodeName, metric.RuleID)

		ch <- prometheus.MustNewConstMetric(
			c.desc[ruleTopicHitCount],
			prometheus.CounterValue, float64(metric.TopicHitCount), metric.NodeName, metric.RuleID,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[ruleExecPassCount],
			prometheus.CounterValue, float64(metric.ExecPassCount), metric.NodeName, metric.RuleID,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[ruleExecFailureCount],
			prometheus.CounterValue, float64(metric.ExecFailureCount), metric.NodeName, metric.RuleID,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[ruleNoResultCount],
			prometheus.CounterValue, float64(metric.NoResultCount), metric.NodeName, metric.RuleID,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[ruleExecRate],
			prometheus.GaugeValue, metric.ExecRate, metric.NodeName, metric.RuleID,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[ruleExecLast5mRate],
			prometheus.GaugeValue, metric.ExecLast5mRate, metric.NodeName, metric.RuleID,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[ruleExecMaxRate],
			prometheus.GaugeValue, metric.ExecMaxRate, metric.NodeName, metric.RuleID,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[ruleActionTotal],
			prometheus.CounterValue, float64(metric.ActionTotal), metric.NodeName, metric.RuleID,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[ruleActionSuccess],
			prometheus.CounterValue, float64(metric.ActionSuccess), metric.NodeName, metric.RuleID,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[ruleActionFailed],
			prometheus.CounterValue, float64(metric.ActionFailed), metric.NodeName, metric.RuleID,
		)
	}

	for i := range bridges {
		ch <- prometheus.MustNewConstMetric(
			c.desc[bridgeResStatus],
			prometheus.GaugeValue, float64(bridges[i].Status), bridges[i].Type, bridges[i].Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[bridgeQueuing],
			prometheus.GaugeValue, float64(bridges[i].Queuing), bridges[i].Type, bridges[i].Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[bridgeLast5mRate],
			prometheus.GaugeValue, bridges[i].RateLast5m, bridges[i].Type, bridges[i].Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[bridgeRateMax],
			prometheus.GaugeValue, bridges[i].RateMax, bridges[i].Type, bridges[i].Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[bridgeFailed],
			prometheus.CounterValue, float64(bridges[i].Failed), bridges[i].Type, bridges[i].Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[bridgeDropped],
			prometheus.CounterValue, float64(bridges[i].Dropped), bridges[i].Type, bridges[i].Name,
		)
	}
	return nil
}
