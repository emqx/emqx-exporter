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
	GetLicense() (*collector.LicenseInfo, error)
	GetClusterStatus() (collector.ClusterStatus, error)
	GetBrokerMetrics() (*collector.Broker, error)
	GetDataBridge() ([]collector.DataBridge, error)
	GetRuleEngineMetrics() ([]collector.RuleEngine, error)
	GetAuthenticationMetrics() ([]collector.DataSource, []collector.Authentication, error)
	GetAuthorizationMetrics() ([]collector.DataSource, []collector.Authorization, error)
}
