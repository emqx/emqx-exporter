package collector

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var _ emqxClientInterface = &client4x{}

type client4x struct {
	edition   edition
	requester *requester
}

func (n *client4x) getLicense() (lic *LicenseInfo, err error) {
	if n.edition == openSource {
		return
	}

	resp := struct {
		Data struct {
			MaxConnections int64  `json:"max_connections"`
			ExpiryAt       string `json:"expiry_at"`
		}
		Code int
	}{}
	err = n.requester.callHTTPGetWithResp("/api/v4/license", &resp)
	if err != nil {
		return
	}

	expiryAt, err := time.Parse("2006-01-02 15:04:05", resp.Data.ExpiryAt)
	if err != nil {
		err = fmt.Errorf("parse expiry time failed: %s", resp.Data.ExpiryAt)
		return
	}

	lic = &LicenseInfo{
		MaxClientLimit: resp.Data.MaxConnections,
		Expiration:     expiryAt.UnixMilli(),
	}
	return
}

func (n *client4x) getClusterStatus() (cluster ClusterStatus, err error) {
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
	err = n.requester.callHTTPGetWithResp("/api/v4/nodes", &resp)
	if err != nil {
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("get err from nodes api: %d", resp.Code)
		return
	}

	cluster.Status = unhealthy
	cluster.NodeUptime = make(map[string]int64)
	cluster.NodeMaxFDs = make(map[string]int)
	cluster.CPULoads = make(map[string]CPULoad)

	for _, data := range resp.Data {
		if data.NodeStatus == "Running" {
			cluster.Status = healthy
		}
		nodeName := cutNodeName(data.Node)
		cluster.NodeUptime[nodeName] = parseUptimeFor4x(data.Uptime)
		cluster.NodeMaxFDs[nodeName] = data.MaxFds

		load := CPULoad{}
		load.Load1, _ = strconv.ParseFloat(data.Load1, 64)
		load.Load5, _ = strconv.ParseFloat(data.Load5, 64)
		load.Load15, _ = strconv.ParseFloat(data.Load15, 64)
		cluster.CPULoads[nodeName] = load
	}
	return
}

func (n *client4x) getBrokerMetrics() (metrics *Broker, err error) {
	resp := struct {
		Data struct {
			Sent     int64 `json:"sent"`
			Received int64 `json:"received"`
		}
		Code int
	}{}
	err = n.requester.callHTTPGetWithResp("/api/v4/monitor/current_metrics", &resp)
	if err != nil {
		return
	}

	metrics = &Broker{
		MsgInputPeriodSec:  resp.Data.Received,
		MsgOutputPeriodSec: resp.Data.Sent,
	}
	return
}

func (n *client4x) getRuleEngineMetrics() (metrics []RuleEngine, err error) {
	resp := struct {
		Data []struct {
			Metrics []struct {
				Node        string  `json:"node"`
				SpeedMax    float64 `json:"speed_max"`
				SpeedLast5m float64 `json:"speed_last5m"`
				Speed       float64 `json:"speed"`
				Matched     int64   `json:"matched"`
				Passed      int64   `json:"passed"`
				NoResult    int64   `json:"no_result"`
				Exception   int64   `json:"exception"`
				Failed      int64   `json:"failed"`
			}
			Actions []struct {
				Metrics []struct {
					Node    string `json:"node"`
					Taken   int64  `json:"taken"`
					Success int64  `json:"success"`
					Failed  int64  `json:"failed"`
				}
			}
			ID      string `json:"id"`
			Enabled bool
		}
		Code int
	}{}
	err = n.requester.callHTTPGetWithResp("/api/v4/rules?_limit=10000", &resp)
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

		fillActionMetrics := func(node string, m *RuleEngine) {
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
			re := RuleEngine{
				NodeName: cutNodeName(m.Node),
				RuleID:   rule.ID,
				//ResStatus:           unknown,
				TopicHitCount:      m.Matched,
				ExecPassCount:      m.Passed,
				ExecFailureCount:   m.Failed,
				NoResultCount:      m.NoResult,
				ExecExceptionCount: m.Exception,
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

func (n *client4x) getDataBridge() (bridges []DataBridge, err error) {
	resp := struct {
		Data []struct {
			ID     string `json:"id"`
			Type   string
			Status bool
		}
		Code int
	}{}
	err = n.requester.callHTTPGetWithResp("/api/v4/resources", &resp)
	if err != nil {
		return
	}

	bridges = make([]DataBridge, len(resp.Data))
	for i, data := range resp.Data {
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

func (n *client4x) getAuthenticationMetrics() ([]DataSource, []Authentication, error) {
	return nil, nil, nil
}

func (n *client4x) getAuthorizationMetrics() ([]DataSource, []Authorization, error) {
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
