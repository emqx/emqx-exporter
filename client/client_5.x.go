package client

import (
	"emqx-exporter/collector"
	"fmt"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

var _ client = &cluster5x{}

type cluster5x struct {
	version string
	edition edition
	client  *fasthttp.Client
}

func (n *cluster5x) getVersion() string {
	return n.version + "-" + n.edition.String()
}

func (n *cluster5x) getLicense() (lic *collector.LicenseInfo, err error) {
	if n.edition == openSource {
		return
	}

	resp := struct {
		MaxConnections int64  `json:"max_connections"`
		ExpiryAt       string `json:"expiry_at"`
	}{}
	err = callHTTPGetWithResp(n.client, "/api/v5/license", &resp)
	if err != nil {
		return
	}

	expiryAt, err := time.Parse("2006-01-02", resp.ExpiryAt)
	if err != nil {
		err = fmt.Errorf("parse expiry time failed: %s", resp.ExpiryAt)
		return
	}

	lic = &collector.LicenseInfo{
		MaxClientLimit: resp.MaxConnections,
		Expiration:     expiryAt.UnixMilli(),
	}
	return
}

func (n *cluster5x) getClusterStatus() (cluster collector.ClusterStatus, err error) {
	resp := []struct {
		Version     string
		Uptime      int64
		NodeStatus  string `json:"node_status"`
		Node        string
		MaxFds      int `json:"max_fds"`
		Connections int64
		Edition     string
		Load1       any `json:"load1"`
		Load5       any `json:"load5"`
		Load15      any `json:"load15"`
	}{{}}
	err = callHTTPGetWithResp(n.client, "/api/v5/nodes", &resp)
	if err != nil {
		return
	}

	cluster.Status = healthy
	cluster.NodeUptime = make(map[string]int64)
	cluster.NodeMaxFDs = make(map[string]int)
	cluster.CPULoads = make(map[string]collector.CPULoad)

	for _, data := range resp {
		if data.NodeStatus != "running" {
			cluster.Status = unhealthy
		}
		nodeName := cutNodeName(data.Node)
		cluster.NodeUptime[nodeName] = data.Uptime / 1000
		cluster.NodeMaxFDs[nodeName] = data.MaxFds

		cpuLoad := collector.CPULoad{}
		// the cpu load value may be string or float64
		if value, ok := data.Load1.(float64); ok {
			cpuLoad.Load1 = value
			cpuLoad.Load5, _ = data.Load5.(float64)
			cpuLoad.Load15, _ = data.Load15.(float64)
		} else if strValue, ok := data.Load1.(string); ok {
			cpuLoad.Load1, _ = strconv.ParseFloat(strValue, 64)
			cpuLoad.Load5, _ = strconv.ParseFloat(data.Load5.(string), 64)
			cpuLoad.Load15, _ = strconv.ParseFloat(data.Load15.(string), 64)
		}
		cluster.CPULoads[nodeName] = cpuLoad

		if data.Edition == "Opensource" {
			n.edition = openSource
		} else {
			n.edition = enterprise
		}
		n.version = data.Version
	}
	return
}

func (n *cluster5x) getBrokerMetrics() (metrics *collector.Broker, err error) {
	resp := struct {
		SentMsgRate     int64 `json:"sent_msg_rate"`
		ReceivedMsgRate int64 `json:"received_msg_rate"`
	}{}
	err = callHTTPGetWithResp(n.client, "/api/v5/monitor_current", &resp)
	if err != nil {
		return
	}

	metrics = &collector.Broker{
		MsgInputPeriodSec:  resp.ReceivedMsgRate,
		MsgOutputPeriodSec: resp.SentMsgRate,
	}
	return
}

func (n *cluster5x) getRuleEngineMetrics() (metrics []collector.RuleEngine, err error) {
	resp := struct {
		Data []struct {
			ID     string `json:"id"`
			Name   string
			Enable bool
		}
	}{}
	err = callHTTPGetWithResp(n.client, "/api/v5/rules?limit=10000", &resp)
	if err != nil {
		return
	}

	for _, rule := range resp.Data {
		if !rule.Enable {
			continue
		}

		metricsResp := struct {
			NodeMetrics []struct {
				Node    string
				Metrics struct {
					Rate       float64 `json:"matched.rate"`
					RateLast5m float64 `json:"matched.rate.last5m"`
					RateMax    float64 `json:"matched.rate.max"`
					Matched    int64
					Passed     int64
					Failed     int64
					//Exception     int64 `json:"failed.exception"`
					NoResult      int64 `json:"failed.no_result"`
					ActionTotal   int64 `json:"actions.total"`
					ActionSuccess int64 `json:"actions.success"`
					ActionFailed  int64 `json:"actions.failed"`
				}
			} `json:"node_metrics"`
		}{}
		err = callHTTPGetWithResp(n.client, fmt.Sprintf("/api/v5/rules/%s/metrics", rule.ID), &metricsResp)
		if err != nil {
			return
		}

		for _, node := range metricsResp.NodeMetrics {
			metrics = append(metrics, collector.RuleEngine{
				NodeName:           cutNodeName(node.Node),
				RuleID:             rule.ID,
				TopicHitCount:      node.Metrics.Matched,
				ExecPassCount:      node.Metrics.Passed,
				ExecFailureCount:   node.Metrics.Failed,
				NoResultCount:      node.Metrics.NoResult,
				ExecRate:           node.Metrics.Rate,
				ExecLast5mRate:     node.Metrics.RateLast5m,
				ExecMaxRate:        node.Metrics.RateMax,
				ActionTotal:        node.Metrics.ActionTotal,
				ActionSuccess:      node.Metrics.ActionSuccess,
				ActionFailed:       node.Metrics.ActionFailed,
				ActionExecTimeCost: nil,
			})
		}
	}
	return
}

func (n *cluster5x) getDataBridge() (bridges []collector.DataBridge, err error) {
	bridgesResp := []struct {
		Name    string
		Type    string
		Status  string
		Metrics struct {
			Queuing    int64
			RateLast5m float64 `json:"rate_last5m"`
			RateMax    float64 `json:"rate_max"`
			Failed     int64
			Dropped    int64
		}
	}{{}}
	err = callHTTPGetWithResp(n.client, "/api/v5/bridges", &bridgesResp)
	if err != nil {
		return
	}

	bridges = make([]collector.DataBridge, len(bridgesResp))
	for i, data := range bridgesResp {
		enabled := unhealthy
		if data.Status == "connected" {
			enabled = healthy
		}
		bridges[i].Type = data.Type
		bridges[i].Name = data.Name
		bridges[i].Status = enabled
		bridges[i].Queuing = data.Metrics.Queuing
		bridges[i].RateLast5m = data.Metrics.RateLast5m
		bridges[i].RateMax = data.Metrics.RateMax
		bridges[i].Failed = data.Metrics.Failed
		bridges[i].Dropped = data.Metrics.Dropped
	}
	return
}

func (n *cluster5x) getAuthenticationMetrics() (dataSources []collector.DataSource, metrics []collector.Authentication, err error) {
	resp := []struct {
		ID      string `json:"id"`
		Backend string
		Enable  bool
	}{{}}
	err = callHTTPGetWithResp(n.client, "/api/v5/authentication", &resp)
	if err != nil {
		return
	}

	for _, plugin := range resp {
		if !plugin.Enable {
			continue
		}

		status := struct {
			NodeMetrics []struct {
				Metrics struct {
					Total      int64
					Success    int64
					Failed     int64
					Rate       float64
					RateLast5m float64 `json:"rate_last5m"`
					RateMax    float64 `json:"rate_max"`
				}
				Node string
			} `json:"node_metrics"`
			Status string
		}{}
		err = callHTTPGetWithResp(n.client, fmt.Sprintf("/api/v5/authentication/%s/status", plugin.ID), &status)
		if err != nil {
			return
		}

		ds := collector.DataSource{
			ResType: plugin.Backend,
			Status:  unhealthy,
		}
		if status.Status == "connected" {
			ds.Status = healthy
		}
		dataSources = append(dataSources, ds)

		for _, node := range status.NodeMetrics {
			m := collector.Authentication{
				NodeName:       cutNodeName(node.Node),
				ResType:        plugin.Backend,
				Total:          node.Metrics.Total,
				AllowCount:     node.Metrics.Success,
				DenyCount:      node.Metrics.Failed,
				ExecRate:       node.Metrics.Rate,
				ExecLast5mRate: node.Metrics.RateLast5m,
				ExecMaxRate:    node.Metrics.RateMax,
				ExecTimeCost:   nil,
			}
			metrics = append(metrics, m)
		}
	}
	return
}

func (n *cluster5x) getAuthorizationMetrics() (dataSources []collector.DataSource, metrics []collector.Authorization, err error) {
	resp := struct {
		Sources []struct {
			Type   string
			Enable bool
		}
	}{}
	err = callHTTPGetWithResp(n.client, "/api/v5/authorization/sources", &resp)
	if err != nil {
		return
	}

	for _, plugin := range resp.Sources {
		if !plugin.Enable {
			continue
		}

		status := struct {
			NodeMetrics []struct {
				Metrics struct {
					Total      int64
					Allow      int64
					Deny       int64
					Rate       float64
					RateLast5m float64 `json:"rate_last5m"`
					RateMax    float64 `json:"rate_max"`
				}
				Node string
			} `json:"node_metrics"`
			Status string
		}{}
		err = callHTTPGetWithResp(n.client, fmt.Sprintf("/api/v5/authorization/sources/%s/status", plugin.Type), &status)
		if err != nil {
			return
		}

		ds := collector.DataSource{
			ResType: plugin.Type,
			Status:  unhealthy,
		}
		if status.Status == "connected" {
			ds.Status = healthy
		}
		dataSources = append(dataSources, ds)

		for _, node := range status.NodeMetrics {
			m := collector.Authorization{
				NodeName:       cutNodeName(node.Node),
				ResType:        plugin.Type,
				Total:          node.Metrics.Total,
				AllowCount:     node.Metrics.Allow,
				DenyCount:      node.Metrics.Deny,
				ExecRate:       node.Metrics.Rate,
				ExecLast5mRate: node.Metrics.RateLast5m,
				ExecMaxRate:    node.Metrics.RateMax,
				ExecTimeCost:   nil,
			}
			metrics = append(metrics, m)
		}
	}
	return
}
