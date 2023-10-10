package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Probes []Probe `yaml:"probes"`
}

type Probe struct {
	Target   string `yaml:"target"`
	Scheme   string `yaml:"scheme,omitempty"`
	ClientID string `yaml:"client_id,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Topic    string `yaml:"topic,omitempty"`
	QoS      byte   `yaml:"qos,omitempty"`
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

	for index, probe := range c.Probes {
		if probe.Scheme == "" {
			probe.Scheme = "tcp"
		}
		if probe.ClientID == "" {
			probe.ClientID = "emqx_exporter_probe"
		}
		if probe.Topic == "" {
			probe.Topic = "emqx_exporter_probe"
		}
		c.Probes[index] = probe
	}

	sc.Lock()
	sc.C = c
	sc.Unlock()

	return nil
}
