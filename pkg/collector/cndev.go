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
	"strconv"
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
		ArmOsMemTotal:                     c.collectArmOsMemTotal,
		ArmOsMemUsed:                      c.collectArmOsMemUsed,
		Capacity:                          c.collectCapacity,
		ChipCPUUtil:                       c.collectChipCPUUtil,
		CoreUtil:                          c.collectCoreUtil,
		D2DCRCError:                       c.collectD2DCRCError,
		D2DCRCErrorOverflow:               c.collectD2DCRCErrorOverflow,
		DDRDataWidth:                      c.collectDDRDataWidth,
		DDRBandWidth:                      c.collectDDRBandWidth,
		ECCAddressForbiddenError:          c.collectECCAddressForbiddenError,
		ECCCorrectedError:                 c.collectECCCorrectedError,
		ECCMultipleError:                  c.collectECCMultipleError,
		ECCMultipleMultipleError:          c.collectECCMultipleMultipleError,
		ECCMultipleOneError:               c.collectECCMultipleOneError,
		ECCOneBitError:                    c.collectECCOneBitError,
		ECCTotalError:                     c.collectECCTotalError,
		ECCUncorrectedError:               c.collectECCUncorrectedError,
		FanSpeed:                          c.collectFanSpeed,
		Health:                            c.collectHealth,
		ImageCodecUtil:                    c.collectImageCodecUtil,
		MemTotal:                          c.collectMemTotal,
		MemUsed:                           c.collectMemUsed,
		MemUtil:                           c.collectMemUtil,
		MLULinkCapabilityP2PTransfer:      c.collectMLULinkCapabilityP2PTransfer,
		MLULinkCapabilityInterlakenSerdes: c.collectMLULinkCapabilityInterlakenSerdes,
		MLULinkCounterCntrReadByte:        c.collectMLULinkCounterCntrReadByte,
		MLULinkCounterCntrReadPackage:     c.collectMLULinkCounterCntrReadPackage,
		MLULinkCounterCntrWriteByte:       c.collectMLULinkCounterCntrWriteByte,
		MLULinkCounterCntrWritePackage:    c.collectMLULinkCounterCntrWritePackage,
		MLULinkCounterErrCorrected:        c.collectMLULinkCounterErrCorrected,
		MLULinkCounterErrCRC24:            c.collectMLULinkCounterErrCRC24,
		MLULinkCounterErrCRC32:            c.collectMLULinkCounterErrCRC32,
		MLULinkCounterErrEccDouble:        c.collectMLULinkCounterErrEccDouble,
		MLULinkCounterErrFatal:            c.collectMLULinkCounterErrFatal,
		MLULinkCounterErrReplay:           c.collectMLULinkCounterErrReplay,
		MLULinkCounterErrUncorrected:      c.collectMLULinkCounterErrUncorrected,
		MLULinkPortMode:                   c.collectMLULinkPortMode,
		MLULinkSpeedFormat:                c.collectMLULinkSpeedFormat,
		MLULinkSpeedValue:                 c.collectMLULinkSpeedValue,
		MLULinkStatusIsActive:             c.collectMLULinkStatusIsActive,
		MLULinkStatusSerdesState:          c.collectMLULinkStatusSerdesState,
		MLULinkVersion:                    c.collectMLULinkVersion,
		NUMANodeID:                        c.collectNUMANodeID,
		PCIeInfo:                          c.collectPCIeInfo,
		PowerUsage:                        c.collectPowerUsage,
		ProcessIpuUtil:                    c.collectProcessIpuUtil,
		ProcessJpuUtil:                    c.collectProcessJpuUtil,
		ProcessMemoryUtil:                 c.collectProcessMemoryUtil,
		ProcessVpuDecodeUtil:              c.collectProcessVpuDecodeUtil,
		ProcessVpuEncodeUtil:              c.collectProcessVpuEncodeUtil,
		Temperature:                       c.collectTemperature,
		TinyCoreUtil:                      c.collectTinyCoreUtil,
		Usage:                             c.collectUsage,
		Util:                              c.collectUtil,
		Version:                           c.collectVersion,
		VideoCodecUtil:                    c.collectVideoCodecUtil,
		VirtualFunctionMemUtil:            c.collectVirtualFunctionMemUtil,
		VirtualFunctionUtil:               c.collectVirtualFunctionUtil,
		VirtualMemTotal:                   c.collectVirtualMemTotal,
		VirtualMemUsed:                    c.collectVirtualMemUsed,
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

func (c *cndevCollector) collectArmOsMemTotal(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, total, err := c.cndev.GetDeviceArmOsMemory(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceArmOsMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(total*1024), labelValues...)
	}
}

func (c *cndevCollector) collectArmOsMemUsed(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		used, _, err := c.cndev.GetDeviceArmOsMemory(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceArmOsMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(used*1024), labelValues...)
	}
}

func (c *cndevCollector) collectCapacity(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, getModelCapacity(stat.model), labelValues...)
	}
}

func (c *cndevCollector) collectChipCPUUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		board, coreUtil, err := c.cndev.GetDeviceCPUUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceCPUUtil", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, cpuCore: ""})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(board), labelValues...)
		for core, util := range coreUtil {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, cpuCore: fmt.Sprintf("%d", core)})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
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

func (c *cndevCollector) collectD2DCRCError(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		d2dCRCError, _, err := c.cndev.GetDeviceCRCInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceCRCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(d2dCRCError), labelValues...)
	}
}

func (c *cndevCollector) collectD2DCRCErrorOverflow(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, d2dCRCErrorOverflow, err := c.cndev.GetDeviceCRCInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceCRCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(d2dCRCErrorOverflow), labelValues...)
	}
}

func (c *cndevCollector) collectDDRBandWidth(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, bandWidth, err := c.cndev.GetDeviceDDRInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceDDRInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, bandWidth, labelValues...)
	}
}

func (c *cndevCollector) collectDDRDataWidth(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		dataWidth, _, err := c.cndev.GetDeviceDDRInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceDDRInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(dataWidth), labelValues...)
	}
}

func (c *cndevCollector) collectECCAddressForbiddenError(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		addressForbiddenError, _, _, _, _, _, _, _, err := c.cndev.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(addressForbiddenError), labelValues...)
	}
}

func (c *cndevCollector) collectECCCorrectedError(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, correctedError, _, _, _, _, _, _, err := c.cndev.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(correctedError), labelValues...)
	}
}

func (c *cndevCollector) collectECCMultipleError(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, _, multipleError, _, _, _, _, _, err := c.cndev.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(multipleError), labelValues...)
	}
}

func (c *cndevCollector) collectECCMultipleMultipleError(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, _, _, multipleMultipleError, _, _, _, _, err := c.cndev.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(multipleMultipleError), labelValues...)
	}
}

func (c *cndevCollector) collectECCMultipleOneError(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, _, _, _, multipleOneError, _, _, _, err := c.cndev.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(multipleOneError), labelValues...)
	}
}

func (c *cndevCollector) collectECCOneBitError(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, _, _, _, _, oneBitError, _, _, err := c.cndev.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(oneBitError), labelValues...)
	}
}

func (c *cndevCollector) collectECCTotalError(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, _, _, _, _, _, totalError, _, err := c.cndev.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(totalError), labelValues...)
	}
}

func (c *cndevCollector) collectECCUncorrectedError(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, _, _, _, _, _, _, uncorrectedError, err := c.cndev.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(uncorrectedError), labelValues...)
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

func (c *cndevCollector) collectImageCodecUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		vpuUtil, err := c.cndev.GetDeviceImageCodecUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceImageCodecUtil", stat.slot))
			continue
		}
		for core, util := range vpuUtil {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, core: core})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
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

func (c *cndevCollector) collectMLULinkCapabilityP2PTransfer(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			p2p, _, err := c.cndev.GetDeviceMLULinkCapability(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCapability", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(p2p), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCapabilityInterlakenSerdes(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, interlakenSerdes, err := c.cndev.GetDeviceMLULinkCapability(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCapability", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(interlakenSerdes), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterCntrReadByte(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			cntrReadByte, _, _, _, _, _, _, _, _, _, _, err := c.cndev.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(cntrReadByte), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterCntrReadPackage(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, cntrReadPackage, _, _, _, _, _, _, _, _, _, err := c.cndev.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(cntrReadPackage), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterCntrWriteByte(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, _, cntrWriteByte, _, _, _, _, _, _, _, _, err := c.cndev.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(cntrWriteByte), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterCntrWritePackage(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, _, _, cntrWritePackage, _, _, _, _, _, _, _, err := c.cndev.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(cntrWritePackage), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrCorrected(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, _, _, _, errCorrected, _, _, _, _, _, _, err := c.cndev.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(errCorrected), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrCRC24(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, _, _, _, _, errCRC24, _, _, _, _, _, err := c.cndev.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(errCRC24), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrCRC32(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, _, _, _, _, _, errCRC32, _, _, _, _, err := c.cndev.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(errCRC32), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrEccDouble(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, _, _, _, _, _, _, errEccDouble, _, _, _, err := c.cndev.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(errEccDouble), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrFatal(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, _, _, _, _, _, _, _, errFatal, _, _, err := c.cndev.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(errFatal), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrReplay(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, _, _, _, _, _, _, _, _, errReplay, _, err := c.cndev.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(errReplay), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrUncorrected(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, _, _, _, _, _, _, _, _, _, errUncorrected, err := c.cndev.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.CounterValue, float64(errUncorrected), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkPortMode(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			mode, err := c.cndev.GetDeviceMLULinkPortMode(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkPortMode", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(mode), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkSpeedFormat(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, speedFormat, err := c.cndev.GetDeviceMLULinkSpeedInfo(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkSpeedInfo", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(speedFormat), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkSpeedValue(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			speedValue, _, err := c.cndev.GetDeviceMLULinkSpeedInfo(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkSpeedInfo", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(speedValue), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkStatusIsActive(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			active, _, err := c.cndev.GetDeviceMLULinkStatus(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkStatus", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(active), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkStatusSerdesState(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			_, serdes, err := c.cndev.GetDeviceMLULinkStatus(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkStatus", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(serdes), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkVersion(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		link := c.cndev.GetDeviceMLULinkPortNumber(stat.slot)
		for i := 0; i < link; i++ {
			major, minor, build, err := c.cndev.GetDeviceMLULinkVersion(stat.slot, uint(i))
			if err != nil {
				log.Println(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkVersion", stat.slot, i))
				continue
			}
			v := calcVersion(major, minor, build)
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, link: i, linkVersion: v})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(1), labelValues...)
		}
	}
}

func (c *cndevCollector) collectNUMANodeID(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		numaNodeID, err := c.cndev.GetDeviceNUMANodeID(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceNUMANodeId", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(numaNodeID), labelValues...)
	}
}

func (c *cndevCollector) collectPCIeInfo(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		slotID, subsystemID, deviceID, vendor, subsystemVendor, domain, bus, device, function, err := c.cndev.GetDevicePCIeInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDevicePCIeInfo", stat.slot))
			continue
		}
		speed, width, err := c.cndev.GetDeviceCurrentPCIeInfo(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceCurrentPCIeInfo", stat.slot))
			continue
		}
		pcieInfo := pcieInfo{
			pcieSlot:            fmt.Sprintf("%d", slotID),
			pcieSubsystem:       fmt.Sprintf("%x", subsystemID),
			pcieDeviceID:        fmt.Sprintf("%x", deviceID),
			pcieVendor:          fmt.Sprintf("%x", vendor),
			pcieSubsystemVendor: fmt.Sprintf("%x", subsystemVendor),
			pcieID:              calcPCIeID(domain, bus, device, function),
			pcieSpeed:           calcPCIeSpeed(speed),
			pcieWidth:           fmt.Sprintf("%d", width),
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, pcieInfo: pcieInfo})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, 1, labelValues...)
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

func (c *cndevCollector) collectProcessIpuUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		pid, ipu, _, _, _, _, err := c.cndev.GetDeviceProcessUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", stat.slot))
			continue
		}
		for i, util := range ipu {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, pid: pid[i]})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectProcessJpuUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		pid, _, jpu, _, _, _, err := c.cndev.GetDeviceProcessUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", stat.slot))
			continue
		}
		for i, util := range jpu {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, pid: pid[i]})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectProcessMemoryUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		pid, _, _, mem, _, _, err := c.cndev.GetDeviceProcessUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", stat.slot))
			continue
		}
		for i, util := range mem {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, pid: pid[i]})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectProcessVpuDecodeUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		pid, _, _, _, vpuDecode, _, err := c.cndev.GetDeviceProcessUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", stat.slot))
			continue
		}
		for i, util := range vpuDecode {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, pid: pid[i]})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectProcessVpuEncodeUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		pid, _, _, _, _, vpuEncode, err := c.cndev.GetDeviceProcessUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", stat.slot))
			continue
		}
		for i, util := range vpuEncode {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, pid: pid[i]})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(util), labelValues...)
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
		if clusters[0] == -100 {
			continue
		}
		for index, cluster := range clusters {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, cluster: fmt.Sprintf("%d", index)})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(cluster), labelValues...)
		}
	}
}

func (c *cndevCollector) collectTinyCoreUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		vpuUtil, err := c.cndev.GetDeviceTinyCoreUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceTinyCoreUtil", stat.slot))
			continue
		}
		for core, util := range vpuUtil {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, core: core})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
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

func (c *cndevCollector) collectVersion(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, 1, labelValues...)
	}
}

func (c *cndevCollector) collectVideoCodecUtil(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		vpuUtil, err := c.cndev.GetDeviceVideoCodecUtil(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceVideoCodecUtil", stat.slot))
			continue
		}
		for core, util := range vpuUtil {
			labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host, core: core})
			ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectVirtualFunctionMemUtil(ch chan<- prometheus.Metric, m metric) {
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
				used, total, _, _, err := c.cndev.GetDeviceMemory(uint(vf<<8 | int(stat.slot)))
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

func (c *cndevCollector) collectVirtualMemTotal(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, _, _, total, err := c.cndev.GetDeviceMemory(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(total*1024*1024), labelValues...)
	}
}

func (c *cndevCollector) collectVirtualMemUsed(ch chan<- prometheus.Metric, m metric) {
	for _, stat := range c.sharedInfo {
		_, _, used, _, err := c.cndev.GetDeviceMemory(stat.slot)
		if err != nil {
			log.Println(errors.Wrapf(err, "Slot %d GetDeviceMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: stat, host: c.host})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(used*1024*1024), labelValues...)
	}
}

func calcVFUtil(utils []int, vfNum int, vfID int) (float64, error) {
	sum := 0
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

func calcPCIeID(domain uint, bus uint, device uint, function uint) string {
	domainStr := strconv.FormatInt(int64(domain), 16)
	domainStr = strings.Repeat("0", 4-len([]byte(domainStr))) + domainStr
	busStr := strconv.FormatInt(int64(bus), 16)
	if bus < 16 {
		busStr = "0" + busStr
	}
	deviceStr := strconv.FormatInt(int64(device), 16)
	if device < 16 {
		deviceStr = "0" + deviceStr
	}
	functionStr := strconv.FormatInt(int64(function), 16)
	return domainStr + ":" + busStr + ":" + deviceStr + "." + functionStr
}

func calcPCIeSpeed(speed int) string {
	var pcieSpeed string
	switch speed {
	case 1:
		pcieSpeed = "2.5GT/s"
	case 2:
		pcieSpeed = "5GT/s"
	case 3:
		pcieSpeed = "8GT/s"
	case 4:
		pcieSpeed = "16GT/s"
	case 5:
		pcieSpeed = "32GT/s"
	case 6:
		pcieSpeed = "64GT/s"
	default:
	}
	return pcieSpeed
}

func calcVersion(major uint, minor uint, build uint) string {
	return fmt.Sprintf("v%d.%d.%d", major, minor, build)
}
