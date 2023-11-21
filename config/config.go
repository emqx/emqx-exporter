package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Metrics *Metrics `yaml:"metrics,omitempty"`
	Probes  []Probe  `yaml:"probes,omitempty"`
}

type Metrics struct {
	APIKey          string           `yaml:"api_key"`
	APISecret       string           `yaml:"api_secret"`
	Target          string           `yaml:"target"`
	Scheme          string           `yaml:"scheme,omitempty"`
	TLSClientConfig *TLSClientConfig `yaml:"tls_config,omitempty"`
}

type Probe struct {
	// Target is the address of the EMQX node to probe. Required.
	Target string `yaml:"target"`
	// Scheme is the protocol scheme of the EMQX node to probe.
	// Enum: [mqtt | tcp | mqtts | ssl | tls | ws | wss]
	// Default: tcp
	Scheme string `yaml:"scheme,omitempty"`
	// ClientID is the MQTT client ID to use when probing.
	// Default: emqx_exporter_probe_<index>
	ClientID string `yaml:"client_id,omitempty"`
	// Username is the MQTT username to use when probing.
	Username string `yaml:"username,omitempty"`
	// Password is the MQTT password to use when probing.
	Password string `yaml:"password,omitempty"`
	// Topic is the MQTT topic to use when probing.
	// Default: emqx-exporter-probe-<index>
	Topic string `yaml:"topic,omitempty"`
	// QoS is the MQTT QoS to use when probing.
	// Default: 0
	QoS byte `yaml:"qos,omitempty"`
	// KeepAlive is the keep alive period in seconds. Defaults to 30 seconds.
	KeepAlive int64 `yaml:"keep_alive,omitempty"`
	// TLSClientConfig is the TLS configuration to use when probing.
	TLSClientConfig *TLSClientConfig `yaml:"tls_config,omitempty"`
}

type TLSClientConfig struct {
	// Server should be accessed without verifying the TLS certificate. For testing only.
	InsecureSkipVerify bool `yaml:"insecure_skip_verify,omitempty"`

	// Server requires TLS client certificate authentication
	CertFile string `yaml:"cert_file,omitempty"`
	// Server requires TLS client certificate authentication
	KeyFile string `yaml:"key_file,omitempty"`
	// Trusted root certificates for server
	CAFile string `yaml:"ca_file,omitempty"`

	// CertData holds PEM-encoded bytes (typically read from a client certificate file).
	// CertData takes precedence over CertFile
	CertData []byte `yaml:"cert_data,omitempty"`
	// KeyData holds PEM-encoded bytes (typically read from a client certificate key file).
	// KeyData takes precedence over KeyFile
	KeyData []byte `yaml:"key_data,omitempty"`
	// CAData holds PEM-encoded bytes (typically read from a root certificates bundle).
	// CAData takes precedence over CAFile
	CAData []byte `yaml:"ca_data,omitempty"`
}

type SafeConfig struct {
	sync.RWMutex
	C                   *Config
	configReloadSuccess prometheus.Gauge
	configReloadSeconds prometheus.Gauge
}

func NewSafeConfig(reg prometheus.Registerer) *SafeConfig {
	configReloadSuccess := promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: "emqx_exporter",
		Name:      "config_last_reload_successful",
		Help:      "EMQX exporter config loaded successfully.",
	})

	configReloadSeconds := promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: "emqx_exporter",
		Name:      "config_last_reload_success_timestamp_seconds",
		Help:      "Timestamp of the last successful configuration reload.",
	})
	return &SafeConfig{C: &Config{}, configReloadSuccess: configReloadSuccess, configReloadSeconds: configReloadSeconds}
}

func (sc *SafeConfig) ReloadConfig(confFile string) (err error) {
	var c = &Config{}
	defer func() {
		if err != nil {
			sc.configReloadSuccess.Set(0)
		} else {
			sc.configReloadSuccess.Set(1)
			sc.configReloadSeconds.SetToCurrentTime()
		}
	}()

	yamlReader, err := os.Open(confFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %s", err)
	}
	defer yamlReader.Close()
	decoder := yaml.NewDecoder(yamlReader)
	decoder.KnownFields(true)

	if err = decoder.Decode(c); err != nil {
		return fmt.Errorf("error parsing config file: %s", err)
	}

	if c.Metrics != nil {
		if c.Metrics.APIKey == "" {
			return fmt.Errorf("metrics.api_key is required")
		}
		if c.Metrics.APISecret == "" {
			return fmt.Errorf("metrics.api_secret is required")
		}
		if c.Metrics.Target == "" {
			return fmt.Errorf("metrics.target is required")
		}
		if c.Metrics.TLSClientConfig != nil {
			if c.Metrics.Scheme == "" {
				c.Metrics.Scheme = "https"
			}
			if c.Metrics.TLSClientConfig.CAData, err = dataFromSliceOrFile(c.Metrics.TLSClientConfig.CAData, c.Metrics.TLSClientConfig.CAFile); err != nil {
				return fmt.Errorf("metrics.ssl_config.ca_data: %s", err)
			}
			if c.Metrics.TLSClientConfig.CertData, err = dataFromSliceOrFile(c.Metrics.TLSClientConfig.CertData, c.Metrics.TLSClientConfig.CertFile); err != nil {
				return fmt.Errorf("metrics.ssl_config.cert_data: %s", err)
			}
			if c.Metrics.TLSClientConfig.KeyData, err = dataFromSliceOrFile(c.Metrics.TLSClientConfig.KeyData, c.Metrics.TLSClientConfig.KeyFile); err != nil {
				return fmt.Errorf("metrics.ssl_config.key_data: %s", err)
			}
		}
		if c.Metrics.Scheme == "" {
			c.Metrics.Scheme = "http"
		}
	}

	for index, probe := range c.Probes {
		if probe.Target == "" {
			return fmt.Errorf("probes[%d].target is required", index)
		}
		if probe.TLSClientConfig != nil {
			if probe.Scheme == "" {
				probe.Scheme = "ssl"
			}
			if probe.TLSClientConfig.CAData, err = dataFromSliceOrFile(probe.TLSClientConfig.CAData, probe.TLSClientConfig.CAFile); err != nil {
				return fmt.Errorf("probes[%d].ssl_config.ca_data: %s", index, err)
			}
			if probe.TLSClientConfig.CertData, err = dataFromSliceOrFile(probe.TLSClientConfig.CertData, probe.TLSClientConfig.CertFile); err != nil {
				return fmt.Errorf("probes[%d].ssl_config.cert_data: %s", index, err)
			}
			if probe.TLSClientConfig.KeyData, err = dataFromSliceOrFile(probe.TLSClientConfig.KeyData, probe.TLSClientConfig.KeyFile); err != nil {
				return fmt.Errorf("probes[%d].ssl_config.key_data: %s", index, err)
			}
		}
		if probe.Scheme == "" {
			probe.Scheme = "tcp"
		}
		if probe.ClientID == "" {
			hostname, _ := os.Hostname()
			hostname = strings.Replace(hostname, ".", "-", -1)
			probe.ClientID = fmt.Sprintf("emqx-exporter-probe-%s-%d", hostname, index)
		}
		if probe.Topic == "" {
			hostname, _ := os.Hostname()
			hostname = strings.Replace(hostname, ".", "-", -1)
			probe.Topic = fmt.Sprintf("emqx-exporter-probe/%s/%d", hostname, index)
		}
		if probe.KeepAlive == 0 {
			probe.KeepAlive = 30
		}
		c.Probes[index] = probe
	}

	sc.Lock()
	sc.C = c
	sc.Unlock()

	return nil
}

func (conf *TLSClientConfig) ToTLSConfig() *tls.Config {
	if conf == nil {
		return nil
	}
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(conf.CAData)
	clientKeyPair, _ := tls.X509KeyPair(conf.CertData, conf.KeyData)
	return &tls.Config{
		InsecureSkipVerify: conf.InsecureSkipVerify,
		RootCAs:            certpool,
		Certificates:       []tls.Certificate{clientKeyPair},
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
	}
}

// dataFromSliceOrFile returns data from the slice (if non-empty), or from the file,
// or an error if an error occurred reading the file
func dataFromSliceOrFile(data []byte, file string) ([]byte, error) {
	if len(data) > 0 {
		return data, nil
	}
	if len(file) > 0 {
		fileData, err := os.ReadFile(file)
		if err != nil {
			return []byte{}, err
		}
		return fileData, nil
	}
	return nil, nil
}
