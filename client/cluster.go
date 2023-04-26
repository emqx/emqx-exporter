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
	var currentVersion string
	for {
		var client client
		var err4, err5 error
		client = &cluster4x{client: httpClient}
		_, err4 = client.getClusterStatus()
		if err4 != nil {
			client = &cluster5x{client: httpClient}
			_, err5 = client.getClusterStatus()
		}
		if err4 != nil && err5 != nil {
			_ = level.Warn(c.logger).Log("check nodes", "couldn't get node info", "addr", *emqxNodes,
				"err4", err4.Error(), "err5", err5.Error())
			client = nil
		} else if currentVersion != client.getVersion() {
			currentVersion = client.getVersion()
			_ = level.Info(c.logger).Log("ClusterVersion", currentVersion)
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

	lic, err = client.getLicense()
	if err != nil || lic == nil {
		return
	}

	lic.RemainingDays = time.Until(time.UnixMilli(lic.Expiration)).Hours() / 24
	lic.RemainingDays, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", lic.RemainingDays), 64)
	return
}

func (c *cluster) GetClusterStatus() (cluster collector.ClusterStatus, err error) {
	client := c.getNode()
	if client == nil {
		cluster.Status = unknown
		return
	}
	cluster, err = client.getClusterStatus()
	if err != nil {
		err = fmt.Errorf("collect cluster status failed. %w", err)
		return
	}
	return
}

func (c *cluster) GetBrokerMetrics() (brokers *collector.Broker, err error) {
	client := c.getNode()
	if client == nil {
		return
	}
	brokers, err = client.getBrokerMetrics()
	if err != nil {
		err = fmt.Errorf("collect broker metrics failed. %w", err)
		return
	}
	return
}

func (c *cluster) GetRuleEngineMetrics() (bridges []collector.DataBridge, res []collector.RuleEngine, err error) {
	client := c.getNode()
	if client == nil {
		return
	}
	bridges, err = client.getDataBridge()
	if err != nil {
		err = fmt.Errorf("collect rule engine data bridge failed. %w", err)
		return
	}
	res, err = client.getRuleEngineMetrics()
	if err != nil {
		err = fmt.Errorf("collect rule engine metrics failed. %w", err)
		return
	}
	return
}

func (c *cluster) GetAuthenticationMetrics() (dataSources []collector.DataSource, auths []collector.Authentication, err error) {
	client := c.getNode()
	if client == nil {
		return
	}
	dataSources, auths, err = client.getAuthenticationMetrics()
	if err != nil {
		err = fmt.Errorf("collect authentication metrics failed. %w", err)
		return
	}
	return
}

func (c *cluster) GetAuthorizationMetrics() (dataSources []collector.DataSource, auths []collector.Authorization, err error) {
	client := c.getNode()
	if client == nil {
		return
	}
	dataSources, auths, err = client.getAuthorizationMetrics()
	if err != nil {
		err = fmt.Errorf("collect authorization metrics failed. %w", err)
		return
	}
	return
}

func (c *cluster) getNode() client {
	c.nodeLock.RLock()
	client := c.client
	c.nodeLock.RUnlock()
	return client
}
