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
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func init() {
	registerCollector(Host, NewHostCollector)
}

type hostCollector struct {
	metrics    collectorMetrics
	host       string
	client     host.Host
	sharedInfo map[string]mluStat
	fnMap      map[string]interface{}
}

func NewHostCollector(m collectorMetrics, bi BaseInfo) Collector {
	c := &hostCollector{
		metrics: m,
		host:    bi.host,
		client:  host.NewHostClient(),
	}
	c.fnMap = map[string]interface{}{
		HostCPUIdle:  c.collectCPUIdle,
		HostCPUTotal: c.collectCPUTotal,
		HostMemTotal: c.collectMemoryTotal,
		HostMemFree:  c.collectMemoryFree,
	}
	log.Debugf("hostCollector: %+v", c)
	return c
}

func (c *hostCollector) init(info map[string]mluStat) error {
	c.sharedInfo = info
	return nil
}

func (c *hostCollector) updateMetrics(m collectorMetrics) {
	c.metrics = m
}

func (c *hostCollector) collect(ch chan<- prometheus.Metric) {
	for name, m := range c.metrics {
		fn := c.fnMap[name]
		f, ok := fn.(func(chan<- prometheus.Metric, metric))
		if !ok {
			log.Warnf("type assertion for fn %s failed, skip", name)
		} else {
			f(ch, m)
		}
	}
}

func (c *hostCollector) collectCPUIdle(ch chan<- prometheus.Metric, m metric) {
	_, idle, err := c.client.GetCPUStats()
	if err != nil {
		log.Errorln(errors.Wrap(err, "GetCPUStats"))
		return
	}
	ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, idle, c.host)
}

func (c *hostCollector) collectCPUTotal(ch chan<- prometheus.Metric, m metric) {
	total, _, err := c.client.GetCPUStats()
	if err != nil {
		log.Errorln(errors.Wrap(err, "GetCPUStats"))
		return
	}
	ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, total, c.host)
}

func (c *hostCollector) collectMemoryTotal(ch chan<- prometheus.Metric, m metric) {
	total, _, err := c.client.GetMemoryStats()
	if err != nil {
		log.Errorln(errors.Wrap(err, "GetMemoryStats"))
		return
	}
	ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, total/1024/1024, c.host)
}

func (c *hostCollector) collectMemoryFree(ch chan<- prometheus.Metric, m metric) {
	_, free, err := c.client.GetMemoryStats()
	if err != nil {
		log.Errorln(errors.Wrap(err, "GetMemoryStats"))
		return
	}
	ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, free/1024/1024, c.host)
}
