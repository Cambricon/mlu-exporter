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
	"log"
	"strings"
	"sync"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/Cambricon/mlu-exporter/pkg/podresources"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector interface {
	collect(ch chan<- prometheus.Metric, mluInfo map[string]mluStat)
	init() error
}

var factories = make(map[string]func(m collectorMetrics, host string) Collector)

func registerCollector(name string, factory func(m collectorMetrics, host string) Collector) {
	factories[name] = factory
}

type collectorMetrics map[string]metric

type metric struct {
	desc   *prometheus.Desc
	labels []string
}

type mluStat struct {
	slot      uint
	model     string
	sn        string
	pcie      string
	boardUtil uint
	coreUtil  []uint
	memTotal  uint
	memUsed   uint
	power     uint
	mcu       string
	driver    string
}

type Collectors struct {
	collectors map[string]Collector
	metrics    map[string]collectorMetrics
	mutex      sync.Mutex
	sharedInfo map[string]mluStat
	cndev      cndev.Cndev
}

func NewCollectors(host string, config string, enabled []string, prefix string) *Collectors {
	cfg := metrics.GetOrDie(config)
	m := getMetrics(cfg, prefix)
	cs := make(map[string]Collector)
	for _, name := range enabled {
		c := factories[name](m[name], host)
		cs[name] = c
	}
	c := &Collectors{
		collectors: cs,
		metrics:    m,
		cndev:      cndev.NewCndevClient(),
	}
	c.init()
	return c
}

func (c *Collectors) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.collectSharedInfo()

	wg := sync.WaitGroup{}
	wg.Add(len(c.collectors))
	for _, col := range c.collectors {
		go func(col Collector) {
			col.collect(ch, c.sharedInfo)
			wg.Done()
		}(col)
	}
	wg.Wait()
}

func (c *Collectors) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		for _, metric := range m {
			ch <- metric.desc
		}
	}
}

func (c *Collectors) collectSharedInfo() {
	var errs *multierror.Error
	num, err := c.cndev.GetDeviceCount()
	errs = multierror.Append(errs, err)
	info := make(map[string]mluStat)
	for i := uint(0); i < num; i++ {
		sn, err := c.cndev.GetDeviceSN(i)
		errs = multierror.Append(errs, err)
		model := c.cndev.GetDeviceModel(i)
		pcie, err := c.cndev.GetDevicePCIeID(i)
		errs = multierror.Append(errs, err)
		board, core, err := c.cndev.GetDeviceUtil(i)
		errs = multierror.Append(errs, err)
		used, total, err := c.cndev.GetDeviceMemory(i)
		errs = multierror.Append(errs, err)
		power, err := c.cndev.GetDevicePower(i)
		errs = multierror.Append(errs, err)
		mcu, driver, err := c.cndev.GetDeviceVersion(i)
		errs = multierror.Append(errs, err)
		info[sn] = mluStat{
			sn:        sn,
			model:     model,
			slot:      i,
			pcie:      pcie,
			boardUtil: board,
			coreUtil:  core,
			memTotal:  total,
			memUsed:   used,
			power:     power,
			mcu:       mcu,
			driver:    driver,
		}
	}
	check(errs.ErrorOrNil())
	c.sharedInfo = info
}

func (c *Collectors) init() {
	for _, collector := range c.collectors {
		err := collector.init()
		check(err)
	}
}

func getMetrics(cfg metrics.Conf, prefix string) map[string]collectorMetrics {
	m := make(map[string]collectorMetrics)
	for collector, metrics := range cfg {
		cms := make(collectorMetrics)
		for key, value := range metrics {
			fqName := value.Name
			if prefix != "" {
				fqName = fmt.Sprintf("%s_%s", prefix, fqName)
			}
			labels := value.Labels.GetKeys()
			labelNames := []string{}
			for _, l := range labels {
				labelNames = append(labelNames, value.Labels[l])
			}
			cms[key] = metric{
				desc:   prometheus.NewDesc(fqName, value.Help, labelNames, nil),
				labels: labels,
			}
		}
		m[collector] = cms
	}
	return m
}

func check(err error) {
	if err != nil {
		log.Panicln("Fatal:", err)
	}
}

type labelInfo struct {
	stat    mluStat
	host    string
	core    int
	cluster string
	vfID    int
	podInfo podresources.PodInfo
}

func getLabelValues(labels []string, info labelInfo) []string {
	values := []string{}
	for _, l := range labels {
		switch l {
		case metrics.Slot:
			values = append(values, fmt.Sprintf("%d", info.stat.slot))
		case metrics.Model:
			values = append(values, info.stat.model)
		case metrics.MluType:
			// suffix needs to be stripped except for "mlu270-x5k"
			if strings.EqualFold(info.stat.model, "mlu270-x5k") {
				values = append(values, info.stat.model)
			} else {
				values = append(values, strings.Split(info.stat.model, "-")[0])
			}
		case metrics.SN:
			values = append(values, info.stat.sn)
		case metrics.Hostname:
			values = append(values, info.host)
		case metrics.Cluster:
			values = append(values, info.cluster)
		case metrics.Core:
			values = append(values, fmt.Sprintf("%d", info.core))
		case metrics.MCU:
			values = append(values, info.stat.mcu)
		case metrics.Driver:
			values = append(values, info.stat.driver)
		case metrics.Namespace:
			values = append(values, info.podInfo.Namespace)
		case metrics.Pod:
			values = append(values, info.podInfo.Pod)
		case metrics.Container:
			values = append(values, info.podInfo.Container)
		case metrics.VFID:
			values = append(values, fmt.Sprintf("%d", info.vfID))
		}
	}
	return values
}
