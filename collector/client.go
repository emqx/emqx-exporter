package collector

import (
	"context"
	"emqx-exporter/config"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
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

type emqxClientInterface interface {
	getLicense() (*LicenseInfo, error)
	getClusterStatus() (ClusterStatus, error)
	getBrokerMetrics() (*Broker, error)
	getDataBridge() ([]DataBridge, error)
	getRuleEngineMetrics() ([]RuleEngine, error)
	getAuthenticationMetrics() ([]DataSource, []Authentication, error)
	getAuthorizationMetrics() ([]DataSource, []Authorization, error)
}

type client struct {
	emqxClient emqxClientInterface
	nodeLock   sync.RWMutex
}

func newClient(metrics *config.Metrics, logger log.Logger) *client {
	c := &client{emqxClient: nil}

	go func() {
		requester := newRequester(metrics)
		for {
			client4 := &client4x{
				requester: requester,
			}
			if _, err := client4.getClusterStatus(); err == nil {
				c.emqxClient = client4
				level.Info(logger).Log("msg", "client4x client created")
				return
			} else {
				level.Debug(logger).Log("msg", "client4x client failed", "err", err)
			}

			client5 := &client5x{
				requester: requester,
			}
			if _, err := client5.getClusterStatus(); err == nil {
				c.emqxClient = client5
				level.Info(logger).Log("msg", "client5x client created")
				return
			} else {
				level.Debug(logger).Log("msg", "client5x client failed", "err", err)
			}

			level.Error(logger).Log("msg", "Couldn't create scraper client, will retry it after 5 seconds", "err", "no scraper node found")
			select {
			case <-context.Background().Done():
			case <-time.After(5 * time.Second):
			}
		}
	}()
	return c
}

func (c *client) getClient() emqxClientInterface {
	c.nodeLock.RLock()
	client := c.emqxClient
	c.nodeLock.RUnlock()
	return client
}
