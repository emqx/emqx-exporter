package client

import (
	"emqx-exporter/collector"
)

const (
	unknown = iota
	unhealthy
	healthy
)

type edition int

const (
	openSource edition = iota
	enterprise
)

func (e edition) String() string {
	switch e {
	case enterprise:
		return "Enterprise"
	default:
		return "OpenSource"
	}
}

type client interface {
	getVersion() string
	getLicense() (*collector.LicenseInfo, error)
	getClusterStatus() (collector.ClusterStatus, error)
	getBrokerMetrics() (*collector.Broker, error)
	getDataBridge() ([]collector.DataBridge, error)
	getRuleEngineMetrics() ([]collector.RuleEngine, error)
	getAuthenticationMetrics() ([]collector.DataSource, []collector.Authentication, error)
	getAuthorizationMetrics() ([]collector.DataSource, []collector.Authorization, error)
}
