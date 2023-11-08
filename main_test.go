// Copyright 2017 The Prometheus Authors
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

package main

import (
	"context"
	"emqx-exporter/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v3"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

var ctx = context.Background()

type testServerContainer struct {
	id          string
	name        string
	image       string
	hostPortMap map[nat.Port]string
}

var emqxContainer = testServerContainer{
	name:  "emqx-for-emqx-exporter-test",
	image: "docker.io/emqx/emqx-enterprise:5.3",
	hostPortMap: map[nat.Port]string{
		"18083/tcp": "28083",
		"18084/tcp": "28084",
		"1883/tcp":  "31883",
		"8883/tcp":  "38883",
		"8083/tcp":  "38083",
		"8084/tcp":  "38084",
	},
}

type testClientContainer struct {
	image string
	pubID string
	subID string
}

var mqttxContainer = testClientContainer{
	image: "docker.io/emqx/mqttx-cli:latest",
}

type testExporter struct {
	httpServer *http.Server
	binDir     string
	bin        string
	port       int
}

var emqxExporter = &testExporter{
	port: 65534,
}

var emqxExporterRunningPort int

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	format.MaxLength = 0

	// fetch the current config
	suiteConfig, reporterConfig := GinkgoConfiguration()
	// adjust it
	suiteConfig.SkipStrings = []string{"NEVER-RUN"}
	reporterConfig.FullTrace = true
	// pass it in to RunSpecs
	RunSpecs(t, "Controller Suite", suiteConfig, reporterConfig)
}

var _ = BeforeSuite(func() {
	var err error
	binName := "emqx-exporter"

	emqxExporter.binDir, err = os.MkdirTemp("/tmp", binName+"-test-bindir-")
	Expect(err).NotTo(HaveOccurred())
	emqxExporter.bin = emqxExporter.binDir + "/" + binName

	copyCert := exec.Command("cp", "-r", "config/example/certs", emqxExporter.binDir+"/certs")
	Expect(copyCert.Run()).NotTo(HaveOccurred())
	createEMQXContainer()
})

var _ = AfterSuite(func() {
	Expect(os.RemoveAll(emqxExporter.binDir)).NotTo(HaveOccurred())
	deleteEMQXContainer()
})

var _ = Describe("Check EMQX Exporter Metrics", Label("metics"), func() {
	Context("check https", func() {
		BeforeEach(func() {
			exporterConfig := config.Config{
				Metrics: &config.Metrics{
					APIKey:    "some_api_key",
					APISecret: "some_api_secret",
					Target:    "127.0.0.1:28084",
					Scheme:    "https",
					TLSClientConfig: &config.TLSClientConfig{
						InsecureSkipVerify: true,
						CAFile:             emqxExporter.binDir + "/certs/cacert.pem",
						CertFile:           emqxExporter.binDir + "/certs/client-cert.pem",
						KeyFile:            emqxExporter.binDir + "/certs/client-key.pem",
					},
				},
			}

			configFile, _ := yaml.Marshal(exporterConfig)
			configFilePath := emqxExporter.binDir + "/config.yml"
			Expect(os.WriteFile(configFilePath, configFile, 0644)).ToNot(HaveOccurred())

			emqxExporterRunningPort = emqxExporter.port
			emqxExporter.port--
			emqxExporter.httpServer = new(http.Server)

			go func() {
				app := kingpin.New("emqx-exporter", "EMQX Exporter")
				_ = run(app, []string{
					"--config.file", configFilePath,
					"--web.listen-address", fmt.Sprintf(":%d", emqxExporterRunningPort),
				}, emqxExporter.httpServer)
			}()

			Eventually(func() (err error) {
				_, err = callExporterAPI("http://127.0.0.1:" + strconv.Itoa(emqxExporterRunningPort) + "/metrics")
				return err
			}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			Expect(emqxExporter.httpServer.Shutdown(ctx)).NotTo(HaveOccurred())
			Expect(os.Remove(emqxExporter.binDir + "/config.yml")).NotTo(HaveOccurred())
		})

		It("should return metrics", func() {
			uri := &fasthttp.URI{}
			uri.SetScheme("http")
			uri.SetHost("127.0.0.1:" + strconv.Itoa(emqxExporterRunningPort))
			uri.SetPath("/metrics")

			By("check emqx_scrape_collector")
			Eventually(func() map[string]*dto.MetricFamily {
				mf, _ := callExporterAPI(uri.String())
				return mf
			}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(And(
				HaveKey(MatchRegexp(`^go_.*`)),
				HaveKey(MatchRegexp(`^promhttp_.*`)),
				HaveKey("emqx_scrape_collector_success"),
				HaveKey("emqx_scrape_collector_duration_seconds"),
			))

			// By("check emqx_authentication")
			// Eventually(func() map[string]*dto.MetricFamily {
			// 	mf, _ := callExporterAPI(uri.String())
			// 	return mf
			// }).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(And(
			// 	HaveKey("emqx_authentication_resource_status"),
			// 	HaveKey("emqx_authentication_total"),
			// 	HaveKey("emqx_authentication_allow_count"),
			// 	HaveKey("emqx_authentication_deny_count"),
			// 	HaveKey("emqx_authentication_exec_rate"),
			// 	HaveKey("emqx_authentication_exec_last5m_rate"),
			// 	HaveKey("emqx_authentication_exec_max_rate"),
			// 	HaveKey("emqx_authentication_exec_time_cost"),
			// ))

			By("check emqx_authorization")
			Eventually(func() map[string]*dto.MetricFamily {
				mf, _ := callExporterAPI(uri.String())
				return mf
			}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(And(
				HaveKey("emqx_authorization_resource_status"),
				HaveKey("emqx_authorization_total"),
				HaveKey("emqx_authorization_allow_count"),
				HaveKey("emqx_authorization_deny_count"),
				HaveKey("emqx_authorization_exec_rate"),
				HaveKey("emqx_authorization_exec_last5m_rate"),
				HaveKey("emqx_authorization_exec_max_rate"),
				HaveKey("emqx_authorization_exec_time_cost"),
			))

			By("check emqx_broker")
			Eventually(func() map[string]*dto.MetricFamily {
				mf, _ := callExporterAPI(uri.String())
				return mf
			}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(And(
				HaveKey("emqx_messages_consume_time_cost"),
				HaveKey("emqx_messages_input_period_second"),
				HaveKey("emqx_messages_output_period_second"),
			))

			By("check emqx_cluster_status")
			Eventually(func() map[string]*dto.MetricFamily {
				mf, _ := callExporterAPI(uri.String())
				return mf
			}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(And(
				HaveKey("emqx_cluster_status"),
				HaveKeyWithValue("emqx_cluster_status", WithTransform(func(m *dto.MetricFamily) int {
					return int(*m.Metric[0].Gauge.Value)
				}, Equal(2)),
				),
				HaveKey("emqx_cluster_node_uptime"),
				HaveKey("emqx_cluster_node_max_fds"),
				HaveKey("emqx_cluster_cpu_load"),
			))

			By("check emqx_license")
			Eventually(func() map[string]*dto.MetricFamily {
				mf, _ := callExporterAPI(uri.String())
				return mf
			}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(And(
				HaveKey("emqx_license_max_client_limit"),
				HaveKey("emqx_license_expiration_time"),
				HaveKey("emqx_license_remaining_days"),
			))

			By("check emqx_rule_engine")
			Eventually(func() map[string]*dto.MetricFamily {
				mf, _ := callExporterAPI(uri.String())
				return mf
			}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).Should(And(
				HaveKey("emqx_rule_bridge_status"),
				HaveKey("emqx_rule_bridge_queuing"),
				HaveKey("emqx_rule_bridge_last5m_rate"),
				HaveKey("emqx_rule_bridge_max_rate"),
				HaveKey("emqx_rule_bridge_failed"),
				HaveKey("emqx_rule_bridge_dropped"),

				HaveKey("emqx_rule_topic_hit_count"),
				HaveKey("emqx_rule_exec_pass_count"),
				HaveKey("emqx_rule_exec_failure_count"),
				HaveKey("emqx_rule_exec_no_result_count"),
				HaveKey("emqx_rule_exec_exception_count"),
				HaveKey("emqx_rule_exec_rate"),
				HaveKey("emqx_rule_exec_last5m_rate"),
				HaveKey("emqx_rule_exec_max_rate"),
				HaveKey("emqx_rule_exec_time_cost"),
				HaveKey("emqx_rule_action_total"),
				HaveKey("emqx_rule_action_success"),
				HaveKey("emqx_rule_action_failed"),
			))
		})
	})
})

var _ = Describe("Check EMQX Exporter Probe", Label("probes"), func() {
	var exporterConfig config.Config
	BeforeEach(func() {
		exporterConfig = config.Config{
			Probes: []config.Probe{
				{Target: "127.0.0.1:31883", Scheme: "tcp"},
				{Target: "127.0.0.1:38883", Scheme: "ssl",
					TLSClientConfig: &config.TLSClientConfig{
						InsecureSkipVerify: true,
						CAFile:             emqxExporter.binDir + "/certs/cacert.pem",
						CertFile:           emqxExporter.binDir + "/certs/client-cert.pem",
						KeyFile:            emqxExporter.binDir + "/certs/client-key.pem",
					},
				},
				{Target: "127.0.0.1:38083/mqtt", Scheme: "ws"},
				{Target: "127.0.0.1:38084/mqtt", Scheme: "wss",
					TLSClientConfig: &config.TLSClientConfig{
						InsecureSkipVerify: true,
						CAFile:             emqxExporter.binDir + "/certs/cacert.pem",
						CertFile:           emqxExporter.binDir + "/certs/client-cert.pem",
						KeyFile:            emqxExporter.binDir + "/certs/client-key.pem",
					},
				},
			},
		}

		configFile, _ := yaml.Marshal(exporterConfig)
		configFilePath := emqxExporter.binDir + "/config.yml"
		Expect(os.WriteFile(configFilePath, configFile, 0644)).ToNot(HaveOccurred())

		emqxExporterRunningPort = emqxExporter.port
		emqxExporter.port--
		emqxExporter.httpServer = new(http.Server)

		go func() {
			app := kingpin.New("emqx-exporter-test", "EMQX Exporter")
			_ = run(app, []string{
				"--config.file", configFilePath,
				"--web.listen-address", fmt.Sprintf(":%d", emqxExporterRunningPort),
			}, emqxExporter.httpServer)
		}()

		Eventually(func() (err error) {
			_, err = callExporterAPI("http://127.0.0.1:" + strconv.Itoa(emqxExporterRunningPort) + "/metrics")
			return err
		}).WithTimeout(30 * time.Second).WithPolling(500 * time.Millisecond).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		Expect(emqxExporter.httpServer.Shutdown(ctx)).NotTo(HaveOccurred())
		Expect(os.Remove(emqxExporter.binDir + "/config.yml")).NotTo(HaveOccurred())
	})

	Context("when the exporter is running", func() {
		DescribeTable("should return probe success",
			func(target string) {
				uri := &fasthttp.URI{}
				uri.SetScheme("http")
				uri.SetHost("127.0.0.1:" + strconv.Itoa(emqxExporterRunningPort))
				uri.SetPath("/probe")
				uri.SetQueryString("target=" + target)

				var mf map[string]*dto.MetricFamily
				Eventually(func() (err error) {
					mf, err = callExporterAPI(uri.String())
					return err
				}).WithTimeout(10 * time.Second).WithPolling(500 * time.Millisecond).ShouldNot(HaveOccurred())

				Expect(mf).Should(And(
					HaveKeyWithValue("emqx_mqtt_probe_duration_seconds", And(
						WithTransform(func(m *dto.MetricFamily) string {
							return m.GetName()
						}, Equal("emqx_mqtt_probe_duration_seconds")),
						WithTransform(func(m *dto.MetricFamily) string {
							return m.GetType().String()
						}, Equal("GAUGE")),
						WithTransform(func(m *dto.MetricFamily) float64 {
							return *m.Metric[0].Gauge.Value
						}, Not(BeZero())),
					)),
					HaveKeyWithValue("emqx_mqtt_probe_success", And(
						WithTransform(func(m *dto.MetricFamily) string {
							return m.GetName()
						}, Equal("emqx_mqtt_probe_success")),
						WithTransform(func(m *dto.MetricFamily) dto.MetricType {
							return m.GetType()
						}, Equal(dto.MetricType_GAUGE)),
						WithTransform(func(m *dto.MetricFamily) float64 {
							return *m.GetMetric()[0].Gauge.Value
						}, BeNumerically("==", 1)),
						WithTransform(func(m *dto.MetricFamily) map[string]string {
							return map[string]string{
								m.GetMetric()[0].Label[0].GetName(): m.GetMetric()[0].Label[0].GetValue(),
							}
						}, HaveKeyWithValue("target", target)),
					)),
				))

			},
			Entry("mqtt", "127.0.0.1:31883"),
			Entry("ssl", "127.0.0.1:38883"),
			Entry("ws", "127.0.0.1:38083/mqtt"),
			Entry("wss", "127.0.0.1:38084/mqtt"),
		)
	})
})

func callExporterAPI(url string) (map[string]*dto.MetricFamily, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", url, http.StatusText(resp.StatusCode))
	}

	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		return nil, err
	}
	return mf, nil
}

func deleteEMQXContainer() {
	cli, err := dockerClient.NewClientWithOpts(
		dockerClient.FromEnv,
		dockerClient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStop(ctx, mqttxContainer.pubID, container.StopOptions{}); err != nil {
		panic(err)
	}
	if err := cli.ContainerRemove(ctx, mqttxContainer.pubID, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}

	if err := cli.ContainerStop(ctx, mqttxContainer.subID, container.StopOptions{}); err != nil {
		panic(err)
	}
	if err := cli.ContainerRemove(ctx, mqttxContainer.subID, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}

	if err := cli.ContainerStop(ctx, emqxContainer.id, container.StopOptions{}); err != nil {
		panic(err)
	}
	if err := cli.ContainerRemove(ctx, emqxContainer.id, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}
}

func createEMQXContainer() {
	bootstrapAPI := "some_api_key:some_api_secret"
	bootstrapAPIFilePath := emqxExporter.binDir + "/bootstrap_api"
	if err := os.WriteFile(bootstrapAPIFilePath, []byte(bootstrapAPI), 0644); err != nil {
		panic(err)
	}

	emqxConfig := `
node {
	name = "emqx@127.0.0.1"
	cookie = "emqxsecretcookie"
	data_dir = "data"
}
cluster {
	name = emqxcl
	discovery_strategy = manual
}
dashboard {
	listeners.http {
		bind = 18083
	}
	listeners.https {
		bind = 18084
		ssl_options {
			cacertfile  = "/opt/emqx/etc/certs/cacert.pem"
			certfile = "/opt/emqx/etc/certs/cert.pem"
			keyfile = "/opt/emqx/etc/certs/key.pem"
			verify = verify_peer
		}
	}
}
listeners {
	tcp.fake{
		bind = 11883
		authentication = [
			{
				mechanism = password_based
				backend	= built_in_database
			}
		]
	}
	tcp.default {
		bind = 1883
	}
	ssl.default {
		bind = 8883
		ssl_options {
			verify = verify_peer
		}
	}
	ws.default {
		bind = 8083
	}
	wss.default {
		bind = 8084
		ssl_options {
			verify = verify_peer
		}
	}
}
bridges {
	mqtt {
		public_broker {
			enable = true
			server = broker.emqx.io
			egress {
				remote {
					payload = "${payload}"
					qos = 1
					retain = false
					topic = test
				}
			}
		}
	}
}
rule_engine {
	rules {
		test {
			actions = ["mqtt:public_broker"]
			sql = "SELECT * FROM \"#\""
		}
	}
}
	`
	emqxConfigPath := emqxExporter.binDir + "/emqx.conf"
	if err := os.WriteFile(emqxConfigPath, []byte(emqxConfig), 0644); err != nil {
		panic(err)
	}

	cli, err := dockerClient.NewClientWithOpts(
		dockerClient.FromEnv,
		dockerClient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, emqxContainer.image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	_, _ = io.Copy(io.Discard, reader)

	containerConf := &container.Config{
		Image: emqxContainer.image,
		Env: []string{
			"EMQX_API_KEY__BOOTSTRAP_FILE=/opt/emqx/data/bootstrap-api",
		},
		ExposedPorts: nat.PortSet{},
		Healthcheck: &container.HealthConfig{
			Test: []string{"CMD", "emqx", "ping"},
		},
	}
	containerHostConf := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: bootstrapAPIFilePath,
				Target: "/opt/emqx/data/bootstrap-api",
			},
			{
				Type:   mount.TypeBind,
				Source: emqxConfigPath,
				Target: "/opt/emqx/etc/emqx.conf",
			},
		},
		PortBindings: nat.PortMap{},
	}

	for containerPort, hostPort := range emqxContainer.hostPortMap {
		containerConf.ExposedPorts[containerPort] = struct{}{}
		containerHostConf.PortBindings[containerPort] = []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: hostPort}}
	}
	emqxResp, err := cli.ContainerCreate(ctx, containerConf, containerHostConf, nil, nil, emqxContainer.name)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, emqxResp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	emqxContainer.id = emqxResp.ID

	var emqxInfo types.ContainerJSON
	Eventually(func() string {
		emqxInfo, _ = cli.ContainerInspect(ctx, emqxResp.ID)
		return emqxInfo.State.Health.Status
	}).WithTimeout(60 * time.Second).WithPolling(1 * time.Second).Should(Equal("healthy"))

	// mqttx client
	reader, err = cli.ImagePull(ctx, "emqx/mqttx-cli", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	_, _ = io.Copy(io.Discard, reader)

	mqttxSubResp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: mqttxContainer.image,
		Cmd:   []string{"mqttx", "bench", "sub", "-t", "test", "-h", emqxInfo.NetworkSettings.IPAddress},
	}, nil, nil, nil, "mqttx-sub")
	if err != nil {
		panic(err)
	}
	if err := cli.ContainerStart(ctx, mqttxSubResp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	mqttxContainer.subID = mqttxSubResp.ID

	mqttxPubResp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: mqttxContainer.image,
		Cmd:   []string{"mqttx", "bench", "pub", "-c", "1", "-t", "test", "-h", emqxInfo.NetworkSettings.IPAddress},
	}, nil, nil, nil, "mqttx-pub")
	if err != nil {
		panic(err)
	}
	if err := cli.ContainerStart(ctx, mqttxPubResp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	mqttxContainer.pubID = mqttxPubResp.ID
}
