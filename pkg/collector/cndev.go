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

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

var capacity = map[string]float64{
	"MLU100":      64,
	"MLU270":      128,
	"MLU270-X5K":  140,
	"MLU220":      16,
	"MLU220-M2":   8,
	"MLU290":      512,
	"MLU370-EVBS": 131,
	"MLU370-S4":   192,
	"MLU370-M8":   332.8,
	"MLU370":      256,
	"MLU365":      64,
}

const maxFanSpeed = 5400

func getModelCapacity(model string) float64 {
	var modelCapacity float64
	if c, ok := capacity[model]; ok {
		modelCapacity = c
	} else {
		m := strings.Split(model, "-")[0]
		modelCapacity = capacity[m]
	}
	return modelCapacity
}

func init() {
	registerCollector(Cndev, NewCndevCollector)
}

type cndevCollector struct {
	metrics    collectorMetrics
	host       string
	cndev      cndev.Cndev
	sharedInfo map[string]mluStat
	fnMap      map[string]interface{}
}

func NewCndevCollector(m collectorMetrics, host string) Collector {
	c := &cndevCollector{
		metrics: m,
		host:    host,
		cndev:   cndev.NewCndevClient(),
	}
	c.fnMap = map[string]interface{}{
		Temperature:         c.collectTemperature,
		Health:              c.collectHealth,
		MemTotal:            c.collectMemTotal,
		MemUsed:             c.collectMemUsed,
		MemUtil:             c.collectMemUtil,
		Util:                c.collectUtil,
		CoreUtil:            c.collectCoreUtil,
		FanSpeed:            c.collectFanSpeed,
		PowerUsage:          c.collectPowerUsage,
		Capacity:            c.collectCapacity,
		Usage:               c.collectUsage,
		Version:             c.collectVersion,
		VirtualMemUtil:      c.collectVirtualMemUtil,
		VirtualFunctionUtil: c.collectVirtualFunctionUtil,
	}
	return c
}

func (c *cndevCollector) init(info map[string]mluStat) error {
	c.sharedInfo = info
	return c.cndev.Init()
}

func (c *cndevCollector) updateMetrics(m collectorMetrics) {
	c.metrics = m
}

func (c *cndevCollector) collect(ch chan<- prometheus.Metric) {
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

func (c *cndevCollector) collectVirtualMemUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		state, err := c.cndev.GetDeviceVfState(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceVfState", stat.slot))
			continue
		}
		vf := 0
		for state != 0 {
			vf++
			if (state & 0x1) != 0 {
				_, _, used, total, err := c.cndev.GetDeviceMemory(uint(vf<<8 | int(stat.slot)))
				if err != nil {
					log.Println(errors.Wrapf(err, "Slot %d VF %d GetDeviceMemory", stat.slot, vf))
					continue
				}
				util := float64(used) / float64(total) * 100
				labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, vf: fmt.Sprintf("%d", vf)})
				ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, util, labelValues...)
			}
			state >>= 1
		}
	}
}

func (c *cndevCollector) collectVirtualFunctionUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		state, err := c.cndev.GetDeviceVfState(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceVfState", stat.slot))
			continue
		}
		vfNum := 0
		for state != 0 {
			vfNum++
			state >>= 1
		}
		_, coreUtil, err := c.cndev.GetDeviceUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceUtil", stat.slot))
			continue
		}
		for vfID := 1; vfID < vfNum+1; vfID++ {
			util, err := calcVFUtil(coreUtil, vfNum, vfID)
			if err != nil {
				log.Println(errors.Wrapf(err, "vf %d calcVFUtil, util %v, vfNum: %d", vfID, coreUtil, vfNum))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, vf: fmt.Sprintf("%d", vfID)})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, util, labelValues...)
		}
	}
}

func (c *cndevCollector) collectTemperature(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		board, clusters, err := c.cndev.GetDeviceTemperature(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceTemperature", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, cluster: ""})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(board), labelValues...)
		for index, cluster := range clusters {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, cluster: fmt.Sprintf("%d", index)})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(cluster), labelValues...)
		}
	}
}

func (c *cndevCollector) collectHealth(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		health, err := c.cndev.GetDeviceHealth(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceHealth", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(health), labelValues...)
	}
}

func (c *cndevCollector) collectMemTotal(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, total, _, _, err := c.cndev.GetDeviceMemory(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(total*1024*1024), labelValues...)
	}
}

func (c *cndevCollector) collectMemUsed(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		used, _, _, _, err := c.cndev.GetDeviceMemory(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(used*1024*1024), labelValues...)
	}
}

func (c *cndevCollector) collectMemUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		used, total, _, _, err := c.cndev.GetDeviceMemory(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceMemory", stat.slot))
			continue
		}
		util := float64(used) / float64(total) * 100
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, util, labelValues...)
	}
}

func (c *cndevCollector) collectUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		board, _, err := c.cndev.GetDeviceUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceUtil", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(board), labelValues...)
	}
}

func (c *cndevCollector) collectCoreUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, coreUtil, err := c.cndev.GetDeviceUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceUtil", stat.slot))
			continue
		}
		for core, util := range coreUtil {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, core: core})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectFanSpeed(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		fan, err := c.cndev.GetDeviceFanSpeed(stat.slot)
		var speed float64
		if fan == 0 {
			speed = -1
		} else {
			speed = float64(maxFanSpeed) * float64(fan) / 100
		}
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceFanSpeed", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, speed, labelValues...)
	}
}

func (c *cndevCollector) collectPowerUsage(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		power, err := c.cndev.GetDevicePower(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDevicePower", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(power), labelValues...)
	}
}

func (c *cndevCollector) collectCapacity(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, getModelCapacity(stat.model), labelValues...)
	}
}

func (c *cndevCollector) collectUsage(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		board, _, err := c.cndev.GetDeviceUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceUtil", stat.slot))
			continue
		}
		usage := getModelCapacity(stat.model) * float64(board) / 100
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, usage, labelValues...)
	}
}

func (c *cndevCollector) collectVersion(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, 1, labelValues...)
	}
}

func calcVFUtil(utils []uint, vfNum int, vfID int) (float64, error) {
	sum := uint(0)
	offset := 0
	switch vfNum {
	case 1, 2, 4, 8:
		offset = len(utils) / vfNum
	default:
		return 0, fmt.Errorf("invalid vfNum %d", vfNum)
	}
	for i := (vfID - 1) * offset; i < vfID*offset; i++ {
		sum += utils[i]
	}
	return float64(sum) / float64(offset), nil
}
