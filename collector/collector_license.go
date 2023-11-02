// Copyright 2019 The Prometheus Authors
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

package collector

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	LicenseSubsystem = "license"
)

const (
	maxClientLimit    = "max_client_limit"
	licenseExpiration = "expiration_time"
	remainingDays     = "remaining_days"
)

func init() {
	registerCollector(LicenseSubsystem, NewLicenseCollector)
}

type licenseCollector struct {
	desc   map[string]*prometheus.Desc
	client *client
}

// NewLicenseCollector returns a new license based collector
func NewLicenseCollector(client *client) (Collector, error) {
	collector := &licenseCollector{
		desc:   make(map[string]*prometheus.Desc),
		client: client,
	}

	metrics := []struct {
		name string
		help string
	}{
		{
			name: maxClientLimit,
			help: "The client limit of license",
		},
		{
			name: licenseExpiration,
			help: "The expiration time of license",
		},
		{
			name: remainingDays,
			help: "The remaining days of license before expiring",
		},
	}

	for _, m := range metrics {
		collector.desc[m.name] = prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				LicenseSubsystem,
				m.name,
			),
			m.help,
			nil,
			nil,
		)
	}
	return collector, nil
}

// Update implements the Collector interface and will collect license info.
func (c *licenseCollector) Update(ch chan<- prometheus.Metric) error {
	lic, err := doGetLicense(c.client)
	if err != nil {
		return err
	}

	if lic == nil {
		return nil
	}

	ch <- prometheus.MustNewConstMetric(
		c.desc[maxClientLimit],
		prometheus.GaugeValue, float64(lic.MaxClientLimit),
	)
	ch <- prometheus.MustNewConstMetric(
		c.desc[licenseExpiration],
		prometheus.GaugeValue, float64(lic.Expiration),
	)
	ch <- prometheus.MustNewConstMetric(
		c.desc[remainingDays],
		prometheus.GaugeValue, lic.RemainingDays,
	)

	return nil
}

type LicenseInfo struct {
	MaxClientLimit int64
	Expiration     int64
	RemainingDays  float64
}

func doGetLicense(c *client) (lic *LicenseInfo, err error) {
	client := c.getClient()
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
