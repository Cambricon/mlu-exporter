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
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector(Cnpapi, NewCnpapiCollector)
}

type cnpapiCollector struct {
	metrics    collectorMetrics
	host       string
	cnpapi     cnpapi.Cnpapi
	sharedInfo map[string]mluStat
	fnMap      map[string]interface{}
}

func NewCnpapiCollector(m collectorMetrics, host string) Collector {
	c := &cnpapiCollector{
		metrics: m,
		host:    host,
		cnpapi:  cnpapi.NewCnpapiClient(),
	}
	c.fnMap = map[string]interface{}{
		PCIeRead:     c.collectPCIeRead,
		PCIeWrite:    c.collectPCIeWrite,
		DramRead:     c.collectDramRead,
		DramWrite:    c.collectDramWrite,
		MLULinkRead:  c.collectMLULinkRead,
		MLULinkWrite: c.collectMLULinkWrite,
	}
	return c
}

func (c *cnpapiCollector) init(info map[string]mluStat) error {
	c.sharedInfo = info
	return c.cnpapi.Init()
}

func (c *cnpapiCollector) updateMetrics(m collectorMetrics) {
	c.metrics = m
}

func (c *cnpapiCollector) collect(ch chan<- prometheus.Metric) {
	for i := 0; i < len(c.sharedInfo); i++ {
		if err := c.cnpapi.PmuFlushData(uint(i)); err != nil {
			log.Printf("cnpapi flush data, err: %v", err)
			return
		}
	}
	for name, m := range c.metrics {
		fn := c.fnMap[name]
		f, ok := fn.(func(chan<- prometheus.Metric, metric))
		if !ok {
			log.Printf("type assertion for fn %s failed, skip", name)
		} else {
			f(ch, m)
		}
	}
}

func (c *cnpapiCollector) collectPCIeRead(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		data, err := c.cnpapi.GetPCIeReadBytes(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetPCIeReadBytes", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(data)/1024/1024/1024, labelValues...)
	}
}

func (c *cnpapiCollector) collectPCIeWrite(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		data, err := c.cnpapi.GetPCIeWriteBytes(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetPCIeWriteBytes", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(data)/1024/1024/1024, labelValues...)
	}
}

func (c *cnpapiCollector) collectDramRead(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		data, err := c.cnpapi.GetDramReadBytes(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDramReadBytes", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(data)/1024/1024/1024, labelValues...)
	}
}

func (c *cnpapiCollector) collectDramWrite(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		data, err := c.cnpapi.GetDramWriteBytes(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDramWriteBytes", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(data)/1024/1024/1024, labelValues...)
	}
}

func (c *cnpapiCollector) collectMLULinkRead(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		data, err := c.cnpapi.GetMLULinkReadBytes(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetMLULinkReadBytes", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(data)/1024/1024/1024, labelValues...)
	}
}

func (c *cnpapiCollector) collectMLULinkWrite(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		data, err := c.cnpapi.GetMLULinkWriteBytes(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetMLULinkWriteBytes", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(data)/1024/1024/1024, labelValues...)
	}
}
