// Copyright 2021 Cambricon, Inc.
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
	"log"

	"github.com/Cambricon/mlu-exporter/pkg/cnpapi"
	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector(metrics.Cnpapi, NewCnpapiCollector)
}

type cnpapiCollector struct {
	metrics collectorMetrics
	host    string
	cnpapi  cnpapi.Cnpapi
}

func NewCnpapiCollector(m collectorMetrics, host string) Collector {
	return &cnpapiCollector{
		metrics: m,
		host:    host,
		cnpapi:  cnpapi.NewCnpapiClient(),
	}
}

func (c *cnpapiCollector) init() error {
	return c.cnpapi.Init()
}

func (c *cnpapiCollector) collect(ch chan<- prometheus.Metric, mluInfo map[string]mluStat) {
	for i := 0; i < len(mluInfo); i++ {
		if err := c.cnpapi.PmuFlushData(uint(i)); err != nil {
			log.Printf("cnpapi flush data, err: %v", err)
			return
		}
	}
	for name, metric := range c.metrics {
		switch name {
		case metrics.PCIeRead:
			c.collectPCIeRead(ch, metric, mluInfo)
		case metrics.PCIeWrite:
			c.collectPCIeWrite(ch, metric, mluInfo)
		case metrics.DramRead:
			c.collectDramRead(ch, metric, mluInfo)
		case metrics.DramWrite:
			c.collectDramWrite(ch, metric, mluInfo)
		case metrics.MLULinkRead:
			c.collectMLULinkRead(ch, metric, mluInfo)
		case metrics.MLULinkWrite:
			c.collectMLULinkWrite(ch, metric, mluInfo)
		}
	}
}

func (c *cnpapiCollector) collectPCIeRead(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		data, err := c.cnpapi.GetPCIeReadBytes(stat.slot)
		check(err)
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(data)/1024/1024/1024, labelValues...)
	}
}

func (c *cnpapiCollector) collectPCIeWrite(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		data, err := c.cnpapi.GetPCIeWriteBytes(stat.slot)
		check(err)
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(data)/1024/1024/1024, labelValues...)
	}
}

func (c *cnpapiCollector) collectDramRead(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		data, err := c.cnpapi.GetDramReadBytes(stat.slot)
		check(err)
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(data)/1024/1024/1024, labelValues...)
	}
}

func (c *cnpapiCollector) collectDramWrite(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		data, err := c.cnpapi.GetDramWriteBytes(stat.slot)
		check(err)
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(data)/1024/1024/1024, labelValues...)
	}
}

func (c *cnpapiCollector) collectMLULinkRead(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		data, err := c.cnpapi.GetMLULinkReadBytes(stat.slot)
		check(err)
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(data)/1024/1024/1024, labelValues...)
	}
}

func (c *cnpapiCollector) collectMLULinkWrite(ch chan<- prometheus.Metric, m metric, info map[string]mluStat) {
	for _, stat := range info {
		data, err := c.cnpapi.GetMLULinkWriteBytes(stat.slot)
		check(err)
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(data)/1024/1024/1024, labelValues...)
	}
}
