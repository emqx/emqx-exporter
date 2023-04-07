package client

import (
	"emqx-exporter/collector"
)

const (
	unknown = iota
	unhealthy
	healthy
)

const (
	version4x = 440
	version5x = 500
)

type edition = int

const (
	openSource edition = iota
	enterprise
)

type client interface {
	getEdition() edition
	GetLicense() (collector.LicenseInfo, error)
	GetClusterStatus() (collector.ClusterStatus, error)
	GetBrokerMetrics() (collector.Broker, error)
	GetDataBridge() ([]collector.DataBridge, error)
	GetRuleEngineMetrics() ([]collector.RuleEngine, error)
	GetAuthenticationMetrics() ([]collector.DataSource, []collector.Authentication, error)
	GetAuthorizationMetrics() ([]collector.DataSource, []collector.Authorization, error)
}
