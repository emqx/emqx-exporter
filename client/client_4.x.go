package client

import (
	"emqx-exporter/collector"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var _ client = &cluster4x{}

type cluster4x struct {
	version string
	client  *fasthttp.Client
}

func (n *cluster4x) getVersion() string {
	return n.version
}

func (n *cluster4x) getLicense() (lic *collector.LicenseInfo, err error) {
	resp := struct {
		Data struct {
			MaxConnections int64  `json:"max_connections"`
			ExpiryAt       string `json:"expiry_at"`
		}
		Code int
	}{}

	data, statusCode, err := callHTTPGet(n.client, "/api/v4/license")
	if statusCode == http.StatusNotFound {
		// open source version doesn't support license api
		err = nil
		return
	}
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(data, &resp)
	if err != nil {
		err = fmt.Errorf("unmarshal license failed: /api/v4/license")
		return
	}

	if resp.Code != 0 {
		err = fmt.Errorf("get err from license api: %d", resp.Code)
		return
	}

	expiryAt, err := time.Parse("2006-01-02 15:04:05", resp.Data.ExpiryAt)
	if err != nil {
		err = fmt.Errorf("parse expiry time failed: %s", resp.Data.ExpiryAt)
		return
	}

	lic = &collector.LicenseInfo{
		MaxClientLimit: resp.Data.MaxConnections,
		Expiration:     expiryAt.UnixMilli(),
	}
	return
}

func (n *cluster4x) getClusterStatus() (cluster collector.ClusterStatus, err error) {
	resp := struct {
		Data []struct {
			Version     string
			Uptime      string
			NodeStatus  string `json:"node_status"`
			Node        string
			MaxFds      int `json:"max_fds"`
			Connections int64
			Load1       string `json:"load1"`
			Load5       string `json:"load5"`
			Load15      string `json:"load15"`
		}
		Code int
	}{}
	err = callHTTPGetWithResp(n.client, "/api/v4/nodes", &resp)
	if err != nil {
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("get err from nodes api: %d", resp.Code)
		return
	}

	cluster.Status = healthy
	cluster.NodeUptime = make(map[string]int64)
	cluster.NodeMaxFDs = make(map[string]int)
	cluster.CPULoads = make(map[string]collector.CPULoad)

	for _, data := range resp.Data {
		if data.NodeStatus != "Running" {
			cluster.Status = unhealthy
		}
		nodeName := cutNodeName(data.Node)
		cluster.NodeUptime[nodeName] = parseUptimeFor4x(data.Uptime)
		cluster.NodeMaxFDs[nodeName] = data.MaxFds

		load := collector.CPULoad{}
		load.Load1, _ = strconv.ParseFloat(data.Load1, 64)
		load.Load5, _ = strconv.ParseFloat(data.Load5, 64)
		load.Load15, _ = strconv.ParseFloat(data.Load15, 64)
		cluster.CPULoads[nodeName] = load

		n.version = data.Version
	}
	return
}

func (n *cluster4x) getBrokerMetrics() (metrics *collector.Broker, err error) {
	resp := struct {
		Data struct {
			Sent     int64
			Received int64
		}
		Code int
	}{}
	data, statusCode, err := callHTTPGet(n.client, "/api/v4/monitor/current_metrics")
	if statusCode == http.StatusNotFound {
		// open source version doesn't support this api
		err = nil
		return
	}
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(data, &resp)
	if err != nil {
		err = fmt.Errorf("unmarshal license failed: /api/v4/monitor/current_metrics")
		return
	}

	if resp.Code != 0 {
		err = fmt.Errorf("get err from monitor api: %d", resp.Code)
		return
	}

	metrics = &collector.Broker{
		MsgInputPeriodSec:  resp.Data.Received,
		MsgOutputPeriodSec: resp.Data.Sent,
	}
	return
}

func (n *cluster4x) getRuleEngineMetrics() (metrics []collector.RuleEngine, err error) {
	resp := struct {
		Data []struct {
			Metrics []struct {
				Node        string
				SpeedMax    float64 `json:"speed_max"`
				SpeedLast5m float64 `json:"speed_last5m"`
				Speed       float64 `json:"speed"`
				Matched     int64
				Passed      int64
				//NoResult    int64
				//Exception   int64
				Failed int64
			}
			Actions []struct {
				Metrics []struct {
					Node    string
					Taken   int64
					Success int64
					Failed  int64
				}
			}
			ID      string `json:"id"`
			Enabled bool
		}
		Code int
	}{}
	err = callHTTPGetWithResp(n.client, "/api/v4/rules?_limit=10000", &resp)
	if err != nil {
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("get err from rules api: %d", resp.Code)
		return
	}

	for _, rule := range resp.Data {
		if !rule.Enabled {
			continue
		}

		fillActionMetrics := func(node string, m *collector.RuleEngine) {
			for i := range rule.Actions {
				for j := range rule.Actions[i].Metrics {
					am := rule.Actions[i].Metrics[j]
					if am.Node == node {
						m.ActionSuccess = am.Success
						m.ActionTotal = am.Taken
						m.ActionFailed = am.Failed
						break
					}
				}
			}
		}
		for _, m := range rule.Metrics {
			re := collector.RuleEngine{
				NodeName: cutNodeName(m.Node),
				RuleID:   rule.ID,
				//ResStatus:           unknown,
				TopicHitCount:    m.Matched,
				ExecPassCount:    m.Passed,
				ExecFailureCount: m.Failed,
				//NoResultCount:      m.NoResult,
				ExecRate:           m.Speed,
				ExecLast5mRate:     m.SpeedLast5m,
				ExecMaxRate:        m.SpeedMax,
				ActionExecTimeCost: nil,
			}
			fillActionMetrics(m.Node, &re)
			metrics = append(metrics, re)
		}
	}
	return
}

func (n *cluster4x) getDataBridge() (bridges []collector.DataBridge, err error) {
	bridgesResp := struct {
		Data []struct {
			ID     string `json:"id"`
			Type   string
			Status bool
		}
		Code int
	}{}
	err = callHTTPGetWithResp(n.client, "/api/v4/resources", &bridgesResp)
	if err != nil {
		return
	}

	bridges = make([]collector.DataBridge, len(bridgesResp.Data))
	for i, data := range bridgesResp.Data {
		enabled := unhealthy
		if data.Status {
			enabled = healthy
		}
		bridges[i].Type = data.Type
		bridges[i].Name = data.ID
		bridges[i].Status = enabled
	}
	return
}

func (n *cluster4x) getAuthenticationMetrics() ([]collector.DataSource, []collector.Authentication, error) {
	return nil, nil, nil
}

func (n *cluster4x) getAuthorizationMetrics() ([]collector.DataSource, []collector.Authorization, error) {
	return nil, nil, nil
}

// parse uptime to second, exp: "2 days, 19 hours, 41 minutes, 47 seconds"
func parseUptimeFor4x(uptime string) int64 {
	times := strings.Split(uptime, ", ")
	var upSecond int64
	for _, t := range times {
		timeUnit := strings.Split(t, " ")
		digit, _ := strconv.Atoi(timeUnit[0])
		switch timeUnit[1] {
		case "days":
			upSecond += int64(digit * 60 * 60 * 24)
		case "hours":
			upSecond += int64(digit * 60 * 60)
		case "minutes":
			upSecond += int64(digit * 60)
		case "seconds":
			upSecond += int64(digit)
		}
	}
	return upSecond
}
