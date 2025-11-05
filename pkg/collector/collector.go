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
	"sync"

	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Collector interface {
	collect(ch chan<- prometheus.Metric)
	init(info *MLUStatMap) error
	updateMetrics(m metrics.CollectorMetrics)
}

var factories = make(map[string]func(m metrics.CollectorMetrics, bi BaseInfo) Collector)

func registerCollector(name string, factory func(m metrics.CollectorMetrics, bi BaseInfo) Collector) {
	factories[name] = factory
}

type Collectors struct {
	collectors map[string]Collector
	metrics    map[string]metrics.CollectorMetrics
	filterPush bool
	mutex      sync.Mutex
}

type rdmaDevice struct {
	name        string
	pcieAddress string
	nicName     string
	ipAddress   string

	domain   uint
	bus      uint
	device   uint
	function uint
}

type BaseInfo struct {
	client     kubernetes.Interface
	host       string
	hostIP     string
	mode       string
	num        uint
	rdmaDevice []rdmaDevice
}

func NewCollectors(enabled []string, metricConfig map[string]metrics.CollectorMetrics, num uint, host string, hostIP string, mode string, shareInfo *MLUStatMap, filterPush bool) *Collectors {
	m := filter(metricConfig, filterPush)
	cs := make(map[string]Collector)
	bi := BaseInfo{
		host:       host,
		hostIP:     hostIP,
		mode:       mode,
		rdmaDevice: getRDMAPCIeInfo(),
	}
	if mode == "env-share" {
		bi.num = num
	}
	if mode == "dynamic-smlu" {
		for _, stat := range shareInfo.Range {
			if !stat.smluEnabled {
				log.Panicf("MLU %d is not in smlu mode", stat.slot)
			}
		}
		bi.client = initClientSet()
	}
	for _, name := range enabled {
		c := factories[name](metricConfig[name], bi)
		cs[name] = c
	}
	c := &Collectors{
		collectors: cs,
		metrics:    m,
		filterPush: filterPush,
	}
	c.init(shareInfo)
	log.Debugf("Collectors: %+v", c)
	return c
}

func (c *Collectors) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(len(c.collectors))
	for _, col := range c.collectors {
		go func(col Collector) {
			col.collect(ch)
			wg.Done()
		}(col)
	}
	wg.Wait()
}

func (c *Collectors) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		for _, metric := range m {
			ch <- metric.Desc
		}
	}
}

func (c *Collectors) UpdateMetrics(metricConfig map[string]metrics.CollectorMetrics) {
	m := filter(metricConfig, c.filterPush)

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.metrics = m

	for name, collector := range c.collectors {
		collector.updateMetrics(m[name])
	}
}

func (c *Collectors) init(info *MLUStatMap) {
	for name, collector := range c.collectors {
		if err := collector.init(info); err != nil {
			log.Error(errors.Wrapf(err, "init collector %s", name))
		}
	}
}

func filter(input map[string]metrics.CollectorMetrics, filterPush bool) map[string]metrics.CollectorMetrics {
	o := make(map[string]metrics.CollectorMetrics)
	for collector, cm := range input {
		cms := make(metrics.CollectorMetrics)
		for key, value := range cm {
			if filterPush && !value.Push {
				continue
			}
			if !filterPush && value.Push {
				continue
			}

			cms[key] = value
		}
		o[collector] = cms
	}
	return o
}

func initClientSet() kubernetes.Interface {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Failed to get in cluser config, err: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to init clientset, err: %v", err)
	}
	return clientset
}
