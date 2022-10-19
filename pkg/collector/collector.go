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
	"sort"
	"strings"
	"sync"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/Cambricon/mlu-exporter/pkg/podresources"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector interface {
	collect(ch chan<- prometheus.Metric)
	init(info map[string]mluStat) error
	updateMetrics(m collectorMetrics)
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
	slot   uint
	model  string
	uuid   string
	sn     string
	mcu    string
	driver string
}

type pcieInfo struct {
	pcieSlot            string
	pcieSubsystem       string
	pcieDeviceID        string
	pcieVendor          string
	pcieSubsystemVendor string
	pcieID              string
	pcieSpeed           string
	pcieWidth           string
}

type Collectors struct {
	collectors map[string]Collector
	metrics    map[string]collectorMetrics
	mutex      sync.Mutex
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
	}
	c.init()
	go c.syncMetrics(config, prefix)
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
			ch <- metric.desc
		}
	}
}

func collectSharedInfo() map[string]mluStat {
	cli := cndev.NewCndevClient()
	if err := cli.Init(); err != nil {
		log.Panic(errors.Wrap(err, "Init"))
	}
	num, err := cli.GetDeviceCount()
	if err != nil {
		log.Panic(errors.Wrap(err, "GetDeviceCount"))
	}
	info := make(map[string]mluStat)
	model := cli.GetDeviceModel(0)
	for i := uint(0); i < num; i++ {
		uuid, err := cli.GetDeviceUUID(i)
		if err != nil {
			log.Panic(errors.Wrap(err, "GetDeviceUUID"))
		}
		sn, err := cli.GetDeviceSN(i)
		if err != nil {
			log.Panic(errors.Wrap(err, "GetDeviceSN"))
		}
		mcuMajor, mcuMinor, mcuBuild, driverMajor, driverMinor, driverBuild, err := cli.GetDeviceVersion(i)
		if err != nil {
			log.Panic(errors.Wrap(err, "GetDeviceVersion"))
		}
		info[uuid] = mluStat{
			sn:     sn,
			uuid:   uuid,
			model:  model,
			slot:   i,
			mcu:    calcVersion(mcuMajor, mcuMinor, mcuBuild),
			driver: calcVersion(driverMajor, driverMinor, driverBuild),
		}
	}
	return info
}

func (c *Collectors) syncMetrics(config string, prefix string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("failed to create watcher, err %v", err)
		return
	}
	defer watcher.Close()
	err = watcher.Add(config)
	if err != nil {
		log.Printf("failed to add config to watcher, err %v", err)
		return
	}
	log.Printf("start watching metrics config %s", config)
	defer log.Printf("stop watching metrics config %s", config)
	for {
		select {
		case event := <-watcher.Events:
			if event.Name == config && event.Op&fsnotify.Remove == fsnotify.Remove {
				log.Printf("inotify: %v, reload config", event.String())
				c.mutex.Lock()
				watcher.Remove(event.Name)
				watcher.Add(config)
				cfg := metrics.GetOrDie(config)
				m := getMetrics(cfg, prefix)
				c.metrics = m
				for name, collector := range c.collectors {
					collector.updateMetrics(m[name])
				}
				c.mutex.Unlock()
			}
		case err := <-watcher.Errors:
			log.Printf("watcher error: %v", err)
		}
	}
}

func (c *Collectors) init() {
	info := collectSharedInfo()
	for name, collector := range c.collectors {
		if err := collector.init(info); err != nil {
			log.Panic(errors.Wrapf(err, "init collector %s", name))
		}
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
			sort.Strings(labels)
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

type labelInfo struct {
	stat        mluStat
	host        string
	core        int
	cpuCore     string
	link        int
	linkVersion string
	cluster     string
	vf          string
	pcieInfo    pcieInfo
	pid         uint32
	podInfo     podresources.PodInfo
}

func getLabelValues(labels []string, info labelInfo) []string {
	values := []string{}
	for _, l := range labels {
		switch l {
		case MLU:
			values = append(values, fmt.Sprintf("%d", info.stat.slot))
		case Model:
			values = append(values, info.stat.model)
		case Type:
			// suffix needs to be stripped except for "mlu270-x5k"
			if strings.EqualFold(info.stat.model, "mlu270-x5k") {
				values = append(values, strings.ToLower(info.stat.model))
			} else {
				values = append(values, strings.ToLower(strings.Split(info.stat.model, "-")[0]))
			}
		case SN:
			values = append(values, info.stat.sn)
		case UUID:
			values = append(values, info.stat.uuid)
		case Node:
			values = append(values, info.host)
		case Cluster:
			values = append(values, info.cluster)
		case Core:
			values = append(values, fmt.Sprintf("%d", info.core))
		case CPUCore:
			values = append(values, info.cpuCore)
		case MCU:
			values = append(values, info.stat.mcu)
		case Driver:
			values = append(values, info.stat.driver)
		case Namespace:
			values = append(values, info.podInfo.Namespace)
		case Pod:
			values = append(values, info.podInfo.Pod)
		case Container:
			values = append(values, info.podInfo.Container)
		case VF:
			values = append(values, info.vf)
		case Link:
			values = append(values, fmt.Sprintf("%d", info.link))
		case LinkVersion:
			values = append(values, info.linkVersion)
		case PCIeSlot:
			values = append(values, info.pcieInfo.pcieSlot)
		case PCIeSubsystem:
			values = append(values, info.pcieInfo.pcieSubsystem)
		case PCIeDeviceID:
			values = append(values, info.pcieInfo.pcieDeviceID)
		case PCIeVendor:
			values = append(values, info.pcieInfo.pcieVendor)
		case PCIeSubsystemVendor:
			values = append(values, info.pcieInfo.pcieSubsystemVendor)
		case PCIeID:
			values = append(values, info.pcieInfo.pcieID)
		case PCIeSpeed:
			values = append(values, info.pcieInfo.pcieSpeed)
		case PCIeWidth:
			values = append(values, info.pcieInfo.pcieWidth)
		case Pid:
			values = append(values, fmt.Sprintf("%d", info.pid))
		default:
			values = append(values, "") // configured label not applicable, add this to prevent panic
		}
	}
	return values
}
