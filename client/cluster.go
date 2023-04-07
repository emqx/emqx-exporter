package client

import (
	"context"
	"emqx-exporter/collector"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	emqxNodes    = kingpin.Flag("emqx.nodes", "The list of EMQX cluster node addr").Default("").String()
	emqxUsername = kingpin.Flag("emqx.auth-username", "The username used for emqx api basic auth").Default("").String()
	emqxPassword = kingpin.Flag("emqx.auth-password", "The password used for emqx api basic auth").Default("").String()
)

type cluster struct {
	client   client
	nodeLock sync.RWMutex
	logger   log.Logger
}

func NewCluster(logger log.Logger) collector.Cluster {
	addrs := strings.Split(*emqxNodes, ",")
	if len(addrs) == 0 {
		panic(fmt.Sprintf("Invalid emqx node addrs: %s", *emqxNodes))
	}
	for _, addr := range addrs {
		if !strings.ContainsRune(addr, ':') {
			panic(fmt.Sprintf("Invalid emqx node addr: %s", addr))
		}
	}

	if *emqxUsername == "" {
		panic("Missing username used for emqx api basic auth")
	}
	if *emqxPassword == "" {
		panic("Missing password used for emqx api basic auth")
	}

	c := &cluster{
		logger: logger,
	}
	go c.checkNodes()
	return c
}

func (c *cluster) checkNodes() {
	httpClient := getHTTPClient(*emqxNodes)
	for {
		var client client
		client = &cluster4x{client: httpClient}
		_, err := client.GetClusterStatus()
		if err != nil {
			client = &cluster5x{client: httpClient}
			_, err = client.GetClusterStatus()
		}
		if err != nil {
			level.Warn(c.logger).Log("check nodes", "couldn't get node info", "addr", *emqxNodes, "err", err.Error())
			client = nil
		}
		c.nodeLock.Lock()
		c.client = client
		c.nodeLock.Unlock()

		select {
		case <-context.Background().Done():
			return
		case <-time.After(5 * time.Second):
		}
	}
}

func (c *cluster) GetLicense() (lic *collector.LicenseInfo, err error) {
	client := c.getNode()
	if client == nil {
		return
	}
	if client.getEdition() == openSource {
		return
	}

	l, err := client.GetLicense()
	if err != nil {
		return
	}
	l.RemainingDays = time.UnixMilli(l.Expiration).Sub(time.Now()).Hours() / 24
	l.RemainingDays, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", l.RemainingDays), 64)
	return &l, nil
}

func (c *cluster) GetClusterStatus() (cluster collector.ClusterStatus, err error) {
	client := c.getNode()
	if client == nil {
		cluster.Status = unknown
		return
	}
	return client.GetClusterStatus()
}

func (c *cluster) GetBrokerMetrics() (brokers collector.Broker, err error) {
	client := c.getNode()
	if client == nil {
		return
	}
	return client.GetBrokerMetrics()
}

func (c *cluster) GetRuleEngineMetrics() (bridges []collector.DataBridge, res []collector.RuleEngine, err error) {
	client := c.getNode()
	if client == nil {
		return
	}
	bridges, err = client.GetDataBridge()
	if err != nil {
		return
	}
	res, err = client.GetRuleEngineMetrics()
	if err != nil {
		return
	}
	return
}

func (c *cluster) GetAuthenticationMetrics() (dataSources []collector.DataSource, auths []collector.Authentication, err error) {
	client := c.getNode()
	if client == nil {
		return
	}
	return client.GetAuthenticationMetrics()
}

func (c *cluster) GetAuthorizationMetrics() (dataSources []collector.DataSource, auths []collector.Authorization, err error) {
	client := c.getNode()
	if client == nil {
		return
	}
	return client.GetAuthorizationMetrics()
}

func (c *cluster) getNode() client {
	c.nodeLock.RLock()
	client := c.client
	c.nodeLock.RUnlock()
	return client
}
