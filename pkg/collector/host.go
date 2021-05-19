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
	"github.com/Cambricon/mlu-exporter/pkg/host"
	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector(metrics.Host, NewHostCollector)
}

type hostCollector struct {
	metrics collectorMetrics
	host    string
	client  host.Host
}

func NewHostCollector(m collectorMetrics, h string) Collector {
	return &hostCollector{
		metrics: m,
		host:    h,
		client:  host.NewHostClient(),
	}
}

func (c *hostCollector) init() error {
	return nil
}

func (c *hostCollector) updateMetrics(m collectorMetrics) {
	c.metrics = m
}

func (c *hostCollector) collect(ch chan<- prometheus.Metric, mluInfo map[string]mluStat) {
	for name, metric := range c.metrics {
		switch name {
		case metrics.CPUIdle:
			c.collectCPUIdle(ch, metric)
		case metrics.CPUTotal:
			c.collectCPUTotal(ch, metric)
		case metrics.MemoryTotal:
			c.collectMemoryTotal(ch, metric)
		case metrics.MemoryFree:
			c.collectMemoryFree(ch, metric)
		}
	}
}

func (c *hostCollector) collectCPUIdle(ch chan<- prometheus.Metric, m metric) {
	_, idle, err := c.client.GetCPUStats()
	check(err)
	ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, idle, c.host)
}

func (c *hostCollector) collectCPUTotal(ch chan<- prometheus.Metric, m metric) {
	total, _, err := c.client.GetCPUStats()
	check(err)
	ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, total, c.host)
}

func (c *hostCollector) collectMemoryTotal(ch chan<- prometheus.Metric, m metric) {
	total, _, err := c.client.GetMemoryStats()
	check(err)
	ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, total/1024/1024, c.host)
}

func (c *hostCollector) collectMemoryFree(ch chan<- prometheus.Metric, m metric) {
	_, free, err := c.client.GetMemoryStats()
	check(err)
	ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, free/1024/1024, c.host)
}
