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
	clusterStatusSubsystem = "cluster"
)

const (
	clusterStatus = "status"
	nodeUptime    = "node_uptime"
	nodeMaxFDs    = "node_max_fds"
	cpuLoad       = "cpu_load"
)

func init() {
	registerCollector(clusterStatusSubsystem, NewClusterStatusCollector)
}

type clusterStatusCollector struct {
	desc    map[string]*prometheus.Desc
	logger  log.Logger
	cluster Cluster
}

// NewClusterStatusCollector returns a new cluster status collector
func NewClusterStatusCollector(cluster Cluster, logger log.Logger) (Collector, error) {
	collector := &clusterStatusCollector{
		desc:    map[string]*prometheus.Desc{},
		logger:  logger,
		cluster: cluster,
	}

	metrics := []struct {
		name   string
		help   string
		labels []string
	}{
		{
			name: clusterStatus,
			help: "The status of cluster",
		},
		{
			name:   nodeUptime,
			help:   "the node uptime",
			labels: []string{"node"},
		},
		{
			name:   nodeMaxFDs,
			help:   "The max fds of node",
			labels: []string{"node"},
		},
		{
			name:   cpuLoad,
			help:   "The load of node cpu",
			labels: []string{"node", "load"},
		},
	}

	for _, m := range metrics {
		collector.desc[m.name] = prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				clusterStatusSubsystem,
				m.name,
			),
			m.help,
			m.labels,
			nil,
		)
	}
	return collector, nil
}

// Update implements the Collector interface and will collect cluster status.
func (c *clusterStatusCollector) Update(ch chan<- prometheus.Metric) error {
	status, err := c.cluster.GetClusterStatus()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.desc[clusterStatus],
		prometheus.GaugeValue, float64(status.Status),
	)
	for node, uptime := range status.NodeUptime {
		ch <- prometheus.MustNewConstMetric(
			c.desc[nodeUptime],
			prometheus.GaugeValue, float64(uptime), node,
		)
	}
	for node, fd := range status.NodeMaxFDs {
		ch <- prometheus.MustNewConstMetric(
			c.desc[nodeMaxFDs],
			prometheus.GaugeValue, float64(fd), node,
		)
	}
	for node, load := range status.CPULoads {
		ch <- prometheus.MustNewConstMetric(
			c.desc[cpuLoad],
			prometheus.GaugeValue, load.Load1, node, "load1",
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[cpuLoad],
			prometheus.GaugeValue, load.Load5, node, "load5",
		)
		ch <- prometheus.MustNewConstMetric(
			c.desc[cpuLoad],
			prometheus.GaugeValue, load.Load15, node, "load15",
		)
	}
	return nil
}
