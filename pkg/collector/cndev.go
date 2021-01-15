// Copyright 2020 Cambricon, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"fmt"
	"strings"

	"github.com/cambricon/mlu-exporter/pkg/cndev"
	"github.com/cambricon/mlu-exporter/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var capacity = map[string]int{
	"MLU100":     64,
	"MLU270":     128,
	"MLU270-X5K": 140,
	"MLU220":     16,
	"MLU220-M2":  8,
	"MLU290":     512,
}

const maxFanSpeed = 5400

func getModelCapacity(model string) int {
	var modelCapacity int
	if c, ok := capacity[model]; ok {
		modelCapacity = c
	} else {
		m := strings.Split(model, "-")[0]
		modelCapacity = capacity[m]
	}
	return modelCapacity
}

func init() {
	registerCollector(metrics.Cndev, NewCndevCollector)
}

type cndevCollector struct {
	metrics collectorMetrics
	host    string
	cndev   cndev.Cndev
}

func NewCndevCollector(m collectorMetrics, host string) Collector {
	return &cndevCollector{
		metrics: m,
		host:    host,
		cndev:   cndev.NewCndevClient(),
	}
}

func (c *cndevCollector) init() error {
	return c.cndev.Init()
}

func (c *cndevCollector) collect(ch chan<- prometheus.Metric, mluInfo map[string]mluStat) {
	for name, metric := range c.metrics {
		switch name {
		case metrics.Temperature:
			c.collectTemperature(ch, metric, mluInfo)
		case metrics.BoardHealth:
			c.collectBoardHealth(ch, metric, mluInfo)
		case metrics.MemTotal:
			c.collectMemTotal(ch, metric, mluInfo)
		case metrics.MemUsed:
			c.collectMemUsed(ch, metric, mluInfo)
		case metrics.MemUtil:
			c.collectMemUtil(ch, metric, mluInfo)
		case metrics.BoardUtil:
			c.collectBoardUtil(ch, metric, mluInfo)
		case metrics.CoreUtil:
			c.collectCoreUtil(ch, metric, mluInfo)
		case metrics.FanSpeed:
			c.collectFanSpeed(ch, metric, mluInfo)
		case metrics.BoardPower:
			c.collectBoardPower(ch, metric, mluInfo)
		case metrics.BoardCapacity:
			c.collectBoardCapacity(ch, metric, mluInfo)
		case metrics.BoardUsage:
			c.collectBoardUsage(ch, metric, mluInfo)
		}
	}
}

func (c *cndevCollector) collectTemperature(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		board, clusters, err := c.cndev.GetDeviceTemperature(stat.slot)
		check(err)
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, cluster: ""})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(board), labelValues...)
		for index, cluster := range clusters {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, cluster: fmt.Sprintf("%d", index)})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(cluster), labelValues...)
		}
	}
}

func (c *cndevCollector) collectBoardHealth(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		health, err := c.cndev.GetDeviceHealth(stat.slot)
		check(err)
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(health), labelValues...)
	}
}

func (c *cndevCollector) collectMemTotal(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(stat.memTotal*1024*1024), labelValues...)
	}
}

func (c *cndevCollector) collectMemUsed(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(stat.memUsed*1024*1024), labelValues...)
	}
}

func (c *cndevCollector) collectMemUtil(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		util := float64(stat.memUsed) / float64(stat.memTotal) * 100
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, util, labelValues...)
	}
}

func (c *cndevCollector) collectBoardUtil(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(stat.boardUtil), labelValues...)
	}
}

func (c *cndevCollector) collectCoreUtil(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		for core, util := range stat.coreUtil {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, core: core})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectFanSpeed(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		fan, err := c.cndev.GetDeviceFanSpeed(stat.slot)
		var speed float64
		if fan == 0 {
			speed = -1
		} else {
			speed = float64(maxFanSpeed) * float64(fan) / 100
		}
		check(err)
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, speed, labelValues...)
	}
}

func (c *cndevCollector) collectBoardPower(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		power, err := c.cndev.GetDevicePower(stat.slot)
		check(err)
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(power), labelValues...)
	}
}

func (c *cndevCollector) collectBoardCapacity(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		modelCapacity := getModelCapacity(stat.model)
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(modelCapacity), labelValues...)
	}
}

func (c *cndevCollector) collectBoardUsage(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		modelCapacity := getModelCapacity(stat.model)
		usage := float64(modelCapacity) * float64(stat.boardUtil) / 100
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, usage, labelValues...)
	}
}
