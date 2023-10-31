package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
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
	Target    string `yaml:"target"`
	APIKey    string `yaml:"api_key"`
	APISecret string `yaml:"api_secret"`
}

type Probe struct {
	Target    string     `yaml:"target"`
	Scheme    string     `yaml:"scheme,omitempty"`
	ClientID  string     `yaml:"client_id,omitempty"`
	Username  string     `yaml:"username,omitempty"`
	Password  string     `yaml:"password,omitempty"`
	Topic     string     `yaml:"topic,omitempty"`
	QoS       byte       `yaml:"qos,omitempty"`
	SSLConfig *SSLConfig `yaml:"ssl_config,omitempty"`
}

type SSLConfig struct {
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
		if c.Metrics.Target == "" {
			return fmt.Errorf("metrics.target is required")
		}
		if c.Metrics.APIKey == "" {
			return fmt.Errorf("metrics.api_key is required")
		}

		if c.Metrics.APISecret == "" {
			return fmt.Errorf("metrics.api_secret is required")
		}
	}

	for index, probe := range c.Probes {
		if probe.Target == "" {
			return fmt.Errorf("probes[%d].target is required", index)
		}
		if probe.SSLConfig != nil {
			if probe.Scheme == "" {
				probe.Scheme = "ssl"
			}
			if probe.SSLConfig.CAData, err = dataFromSliceOrFile(probe.SSLConfig.CAData, probe.SSLConfig.CAFile); err != nil {
				return fmt.Errorf("probes[%d].ssl_config.ca_data: %s", index, err)
			}
			if probe.SSLConfig.CertData, err = dataFromSliceOrFile(probe.SSLConfig.CertData, probe.SSLConfig.CertFile); err != nil {
				return fmt.Errorf("probes[%d].ssl_config.cert_data: %s", index, err)
			}
			if probe.SSLConfig.KeyData, err = dataFromSliceOrFile(probe.SSLConfig.KeyData, probe.SSLConfig.KeyFile); err != nil {
				return fmt.Errorf("probes[%d].ssl_config.key_data: %s", index, err)
			}
		}
		if probe.Scheme == "" {
			probe.Scheme = "tcp"
		}
		if probe.ClientID == "" {
			probe.ClientID = "emqx_exporter_probe_" + fmt.Sprintf("%d", index)
		}
		if probe.Topic == "" {
			probe.Topic = "emqx-exporter-probe-" + fmt.Sprintf("%d", index)
		}
		c.Probes[index] = probe
	}

	sc.Lock()
	sc.C = c
	sc.Unlock()

	return nil
}

func (sslConfig *SSLConfig) ToTLSConfig() *tls.Config {
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(sslConfig.CAData)
	clientKeyPair, _ := tls.X509KeyPair(sslConfig.CertData, sslConfig.KeyData)
	return &tls.Config{
		InsecureSkipVerify: sslConfig.InsecureSkipVerify,
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
