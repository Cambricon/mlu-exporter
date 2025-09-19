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
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func init() {
	registerCollector(Host, NewHostCollector)
}

type hostCollector struct {
	metrics    metrics.CollectorMetrics
	host       string
	client     host.Host
	sharedInfo *MLUStatMap
	fnMap      map[string]interface{}
}

func NewHostCollector(m metrics.CollectorMetrics, bi BaseInfo) Collector {
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

func (c *hostCollector) init(info *MLUStatMap) error {
	c.sharedInfo = info
	return nil
}

func (c *hostCollector) updateMetrics(m metrics.CollectorMetrics) {
	c.metrics = m
}

func (c *hostCollector) collect(ch chan<- prometheus.Metric) {
	for name, m := range c.metrics {
		fn, ok := c.fnMap[name]
		if !ok {
			continue
		}
		f, ok := fn.(func(chan<- prometheus.Metric, metrics.Metric))
		if !ok {
			log.Warnf("type assertion for fn %s failed, skip", name)
		} else {
			f(ch, m)
		}
	}
}

func (c *hostCollector) collectCPUIdle(ch chan<- prometheus.Metric, m metrics.Metric) {
	_, idle, err := c.client.GetCPUStats()
	if err != nil {
		log.Errorln(errors.Wrap(err, "GetCPUStats"))
		return
	}
	ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, idle, c.host)
}

func (c *hostCollector) collectCPUTotal(ch chan<- prometheus.Metric, m metrics.Metric) {
	total, _, err := c.client.GetCPUStats()
	if err != nil {
		log.Errorln(errors.Wrap(err, "GetCPUStats"))
		return
	}
	ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, total, c.host)
}

func (c *hostCollector) collectMemoryTotal(ch chan<- prometheus.Metric, m metrics.Metric) {
	total, _, err := c.client.GetMemoryStats()
	if err != nil {
		log.Errorln(errors.Wrap(err, "GetMemoryStats"))
		return
	}
	ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, total/1024/1024, c.host)
}

func (c *hostCollector) collectMemoryFree(ch chan<- prometheus.Metric, m metrics.Metric) {
	_, free, err := c.client.GetMemoryStats()
	if err != nil {
		log.Errorln(errors.Wrap(err, "GetMemoryStats"))
		return
	}
	ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, free/1024/1024, c.host)
}
