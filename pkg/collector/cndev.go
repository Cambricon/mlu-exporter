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
	"sort"
	"strconv"
	"strings"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const maxFanSpeed = 5400

func init() {
	registerCollector(Cndev, NewCndevCollector)
}

type cndevCollector struct {
	baseInfo      BaseInfo
	client        cndev.Cndev
	fnMap         map[string]interface{}
	metrics       metrics.CollectorMetrics
	sharedInfo    map[string]MLUStat
	mluXIDCounter map[cndev.XIDInfo]int
	lastXID       map[cndev.DeviceInfo]int64
}

func NewCndevCollector(m metrics.CollectorMetrics, bi BaseInfo) Collector {
	c := &cndevCollector{
		baseInfo:      bi,
		client:        cndev.NewCndevClient(),
		metrics:       m,
		mluXIDCounter: map[cndev.XIDInfo]int{},
		lastXID:       map[cndev.DeviceInfo]int64{},
	}

	c.fnMap = map[string]interface{}{
		ArmOsMemTotal:                     c.collectArmOsMemTotal,
		ArmOsMemUsed:                      c.collectArmOsMemUsed,
		ChipCPUUtil:                       c.collectChipCPUUtil,
		ChipTemperature:                   c.collectChipTemperature,
		ClusterTemperature:                c.collectClusterTemperature,
		CoreUtil:                          c.collectCoreUtil,
		D2DCRCError:                       c.collectD2DCRCError,
		D2DCRCErrorOverflow:               c.collectD2DCRCErrorOverflow,
		DDRDataWidth:                      c.collectDDRDataWidth,
		DDRBandWidth:                      c.collectDDRBandWidth,
		DDRFrequency:                      c.collectDDRFrequency,
		DoubleBitRetiredPageCount:         c.collectDoubleBitRetiredPageCount,
		DramEccDbeCount:                   c.collectDramEccDbeCount,
		DramEccSbeCount:                   c.collectDramEccSbeCount,
		ECCAddressForbiddenError:          c.collectECCAddressForbiddenError,
		ECCCorrectedError:                 c.collectECCCorrectedError,
		ECCMultipleError:                  c.collectECCMultipleError,
		ECCMultipleMultipleError:          c.collectECCMultipleMultipleError,
		ECCMultipleOneError:               c.collectECCMultipleOneError,
		ECCOneBitError:                    c.collectECCOneBitError,
		ECCTotalError:                     c.collectECCTotalError,
		ECCUncorrectedError:               c.collectECCUncorrectedError,
		FanSpeed:                          c.collectFanSpeed,
		Frequency:                         c.collectFrequency,
		Health:                            c.collectHealth,
		HeartbeatCount:                    c.collectHeartbeatCount,
		ImageCodecUtil:                    c.collectImageCodecUtil,
		MachineMLUNums:                    c.collectMachineMLUNums,
		MemTemperature:                    c.collectMemTemperature,
		MemTotal:                          c.collectMemTotal,
		MemUsed:                           c.collectMemUsed,
		MemUtil:                           c.collectMemUtil,
		MLULinkCapabilityP2PTransfer:      c.collectMLULinkCapabilityP2PTransfer,
		MLULinkCapabilityInterlakenSerdes: c.collectMLULinkCapabilityInterlakenSerdes,
		MLULinkCounterCntrCnpPackage:      c.collectMLULinkCounterCntrCnpPackage,
		MLULinkCounterCntrPfcPackage:      c.collectMLULinkCounterCntrPfcPackage,
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
		ParityError:                       c.collectParityError,
		PCIeInfo:                          c.collectPCIeInfo,
		PCIeReadThroughput:                c.collectPCIeReadThroughput,
		PCIeReplayCount:                   c.collectPCIeReplayCount,
		PCIeWriteThroughput:               c.collectPCIeWriteThroughput,
		PowerUsage:                        c.collectPowerUsage,
		ProcessIpuUtil:                    c.collectProcessIpuUtil,
		ProcessJpuUtil:                    c.collectProcessJpuUtil,
		ProcessMemoryUtil:                 c.collectProcessMemoryUtil,
		ProcessVpuDecodeUtil:              c.collectProcessVpuDecodeUtil,
		ProcessVpuEncodeUtil:              c.collectProcessVpuEncodeUtil,
		RemappedCorrectRows:               c.collectRemappedCorrectRows,
		RemappedFailedRows:                c.collectRemappedFailedRows,
		RemappedPendingRows:               c.collectRemappedPendingRows,
		RemappedUncorrectRows:             c.collectRemappedUncorrectRows,
		SingleBitRetiredPageCount:         c.collectSingleBitRetiredPageCount,
		SramEccDbeCount:                   c.collectSramEccDbeCount,
		SramEccParityCount:                c.collectSramEccParityCount,
		SramEccSbeCount:                   c.collectSramEccSbeCount,
		Temperature:                       c.collectTemperature,
		TinyCoreUtil:                      c.collectTinyCoreUtil,
		Util:                              c.collectUtil,
		Version:                           c.collectVersion,
		VideoCodecUtil:                    c.collectVideoCodecUtil,
		VirtualFunctionMemUtil:            c.collectVirtualFunctionMemUtil,
		VirtualFunctionPowerUsage:         c.collectVirtualFunctionPowerUsage,
		VirtualFunctionUtil:               c.collectVirtualFunctionUtil,
		VirtualMemTotal:                   c.collectVirtualMemTotal,
		VirtualMemUsed:                    c.collectVirtualMemUsed,
		LastXIDError:                      c.collectLastXIDError,
		XIDErrorCount:                     c.collectXIDErrorCount,
	}

	log.Debugf("cndevCollector: %+v", c)
	return c
}

func (c *cndevCollector) init(info map[string]MLUStat) error {
	c.sharedInfo = info
	if err := c.client.Init(false); err != nil {
		return err
	}
	go c.waitXIDEvents()
	return nil
}

func (c *cndevCollector) updateMetrics(m metrics.CollectorMetrics) {
	c.metrics = m
}

func (c *cndevCollector) collect(ch chan<- prometheus.Metric) {
	for name, m := range c.metrics {
		fn := c.fnMap[name]
		f, ok := fn.(func(chan<- prometheus.Metric, metrics.Metric))
		if !ok {
			log.Warnf("type assertion for fn %s failed, skip", name)
		} else {
			f(ch, m)
		}
	}
}

func (c *cndevCollector) collectArmOsMemTotal(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, total, err := c.client.GetDeviceArmOsMemory(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceArmOsMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(total*1024), labelValues...)
	}
}

func (c *cndevCollector) collectArmOsMemUsed(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		used, _, err := c.client.GetDeviceArmOsMemory(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceArmOsMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(used*1024), labelValues...)
	}
}

func (c *cndevCollector) collectChipCPUUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		board, coreUtil, err := c.client.GetDeviceCPUUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCPUUtil", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, cpuCore: ""})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(board), labelValues...)
		for core, util := range coreUtil {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, cpuCore: fmt.Sprintf("%d", core)})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectChipTemperature(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, _, chip, _, _, err := c.client.GetDeviceTemperature(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceTemperature", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(chip), labelValues...)
	}
}

func (c *cndevCollector) collectClusterTemperature(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, _, _, clusters, _, err := c.client.GetDeviceTemperature(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceTemperature", stat.slot))
			continue
		}
		for index, cluster := range clusters {
			if cluster < 0 {
				log.Debugf("Slot %d, cluster index %d, temperature %d is invalid", stat.slot, index, cluster)
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, cluster: fmt.Sprintf("%d", index)})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(cluster), labelValues...)
		}
	}
}

func (c *cndevCollector) collectCoreUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, coreUtil, err := c.client.GetDeviceUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceUtil", stat.slot))
			continue
		}
		for core, util := range coreUtil {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, core: core})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectD2DCRCError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.crcDisabled {
			continue
		}
		d2dCRCError, _, err := c.client.GetDeviceCRCInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCRCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(d2dCRCError), labelValues...)
	}
}

func (c *cndevCollector) collectD2DCRCErrorOverflow(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.crcDisabled {
			continue
		}
		_, d2dCRCErrorOverflow, err := c.client.GetDeviceCRCInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCRCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(d2dCRCErrorOverflow), labelValues...)
	}
}

func (c *cndevCollector) collectDDRBandWidth(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, bandWidth, err := c.client.GetDeviceDDRInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceDDRInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, bandWidth, labelValues...)
	}
}

func (c *cndevCollector) collectDDRDataWidth(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		dataWidth, _, err := c.client.GetDeviceDDRInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceDDRInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(dataWidth), labelValues...)
	}
}

func (c *cndevCollector) collectDDRFrequency(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, ddrFrequency, err := c.client.GetDeviceFrequency(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceFrequency", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(ddrFrequency), labelValues...)
	}
}

func (c *cndevCollector) collectDoubleBitRetiredPageCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.retiredPageDisabled {
			continue
		}
		cause, pageCount, err := c.client.GetDeviceRetiredPageInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRetiredPageInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		if cause != 1 {
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, 0, labelValues...)
			continue
		}
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(pageCount), labelValues...)
	}
}

func (c *cndevCollector) collectDramEccDbeCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.memEccDisabled {
			continue
		}
		_, _, _, _, dramEccDbeCounter, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(dramEccDbeCounter), labelValues...)
	}
}

func (c *cndevCollector) collectDramEccSbeCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.memEccDisabled {
			continue
		}
		_, _, _, dramEccSbeCounter, _, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(dramEccSbeCounter), labelValues...)
	}
}

func (c *cndevCollector) collectECCAddressForbiddenError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.eccDisabled {
			continue
		}
		addressForbiddenError, _, _, _, _, _, _, _, err := c.client.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(addressForbiddenError), labelValues...)
	}
}

func (c *cndevCollector) collectECCCorrectedError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.eccDisabled {
			continue
		}
		_, correctedError, _, _, _, _, _, _, err := c.client.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(correctedError), labelValues...)
	}
}

func (c *cndevCollector) collectECCMultipleError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.eccDisabled {
			continue
		}
		_, _, multipleError, _, _, _, _, _, err := c.client.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(multipleError), labelValues...)
	}
}

func (c *cndevCollector) collectECCMultipleMultipleError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.eccDisabled {
			continue
		}
		_, _, _, multipleMultipleError, _, _, _, _, err := c.client.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(multipleMultipleError), labelValues...)
	}
}

func (c *cndevCollector) collectECCMultipleOneError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.eccDisabled {
			continue
		}
		_, _, _, _, multipleOneError, _, _, _, err := c.client.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(multipleOneError), labelValues...)
	}
}

func (c *cndevCollector) collectECCOneBitError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.eccDisabled {
			continue
		}
		_, _, _, _, _, oneBitError, _, _, err := c.client.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(oneBitError), labelValues...)
	}
}

func (c *cndevCollector) collectECCTotalError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.eccDisabled {
			continue
		}
		_, _, _, _, _, _, totalError, _, err := c.client.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(totalError), labelValues...)
	}
}

func (c *cndevCollector) collectECCUncorrectedError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.eccDisabled {
			continue
		}
		_, _, _, _, _, _, _, uncorrectedError, err := c.client.GetDeviceECCInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceECCInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(uncorrectedError), labelValues...)
	}
}

func (c *cndevCollector) collectFanSpeed(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		fan, err := c.client.GetDeviceFanSpeed(stat.slot)
		var speed float64
		if fan == 0 {
			speed = -1
		} else {
			speed = float64(maxFanSpeed) * float64(fan) / 100
		}
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceFanSpeed", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, speed, labelValues...)
	}
}

func (c *cndevCollector) collectFrequency(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		frequency, _, err := c.client.GetDeviceFrequency(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceFrequency", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(frequency), labelValues...)
	}
}

func (c *cndevCollector) collectHealth(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		health, err := c.client.GetDeviceHealth(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceHealth", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(health), labelValues...)
	}
}

func (c *cndevCollector) collectHeartbeatCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.heartBeatDisabled {
			continue
		}
		heartbeatCount, err := c.client.GetDeviceHeartbeatCount(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceHeartbeatCount", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(heartbeatCount), labelValues...)
	}
}

func (c *cndevCollector) collectImageCodecUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		vpuUtil, err := c.client.GetDeviceImageCodecUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceImageCodecUtil", stat.slot))
			continue
		}
		for core, util := range vpuUtil {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, core: core})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMachineMLUNums(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		num, err := c.client.GetDeviceCount()
		if err != nil {
			log.Error(errors.Wrap(err, "GetDeviceCount"))
			return
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: MLUStat{model: stat.model, driver: stat.driver}, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(num), labelValues...)
		return
	}
}

func (c *cndevCollector) collectMemTemperature(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, mem, _, _, memDies, err := c.client.GetDeviceTemperature(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceTemperature", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, memoryDie: ""})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(mem), labelValues...)
		for index, memDie := range memDies {
			if memDie < 0 {
				log.Debugf("Slot %d, cluster index %d, temperature %d is invalid", stat.slot, index, memDie)
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, memoryDie: fmt.Sprintf("%d", index)})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(memDie), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMemTotal(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, total, _, _, err := c.client.GetDeviceMemory(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(total*1024), labelValues...)
	}
}

func (c *cndevCollector) collectMemUsed(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		used, _, _, _, err := c.client.GetDeviceMemory(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(used*1024), labelValues...)
	}
}

func (c *cndevCollector) collectMemUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		used, total, _, _, err := c.client.GetDeviceMemory(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemory", stat.slot))
			continue
		}
		util := float64(used) / float64(total) * 100
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: ""})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, util, labelValues...)

		if stat.smluEnabled {
			smluInfos, err := c.client.GetAllSMluInfo(stat.slot)
			if err != nil {
				log.Warn(errors.Wrap(err, "GetAllSMluInfo"))
				continue
			}
			stat.smluInfos = smluInfos
			for _, info := range smluInfos {
				var labelValues []string
				if c.baseInfo.mode == "dynamic-smlu" {
					labelValues = getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, typ: "vmemory", vf: fmt.Sprintf("%d", info.InstanceID)})
				} else {
					labelValues = getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", info.InstanceID)})
				}
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(info.MemUsed)/float64(info.MemTotal)*100, labelValues...)
			}
			continue
		}

		if c.baseInfo.num > 0 {
			for i := uint(0); i < c.baseInfo.num; i++ {
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", i+1)})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, util, labelValues...)
			}
			continue
		}

		// for mim
		state, err := c.client.GetDeviceVfState(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceVfState", stat.slot))
			continue
		}
		vf := 0
		for state != 0 {
			vf++
			if (state & 0x1) != 0 {
				used, total, _, _, err := c.client.GetDeviceMemory(uint(vf<<8 | int(stat.slot)))
				if err != nil {
					log.Errorln(errors.Wrapf(err, "Slot %d VF %d GetDeviceMemory", stat.slot, vf))
					continue
				}
				util := float64(used) / float64(total) * 100
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", vf)})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, util, labelValues...)
			}
			state >>= 1
		}
	}
}

func (c *cndevCollector) collectMLULinkCapabilityP2PTransfer(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			p2p, _, err := c.client.GetDeviceMLULinkCapability(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCapability", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(p2p), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCapabilityInterlakenSerdes(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			_, interlakenSerdes, err := c.client.GetDeviceMLULinkCapability(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCapability", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(interlakenSerdes), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterCntrCnpPackage(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, _, _, _, _, _, _, _, _, _, _, cntrCnpPackage, _, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(cntrCnpPackage), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterCntrPfcPackage(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, _, _, _, _, _, _, _, _, _, _, _, cntrPfcPackage, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(cntrPfcPackage), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterCntrReadByte(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			cntrReadByte, _, _, _, _, _, _, _, _, _, _, _, _, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(cntrReadByte), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterCntrReadPackage(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, cntrReadPackage, _, _, _, _, _, _, _, _, _, _, _, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(cntrReadPackage), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterCntrWriteByte(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, _, cntrWriteByte, _, _, _, _, _, _, _, _, _, _, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(cntrWriteByte), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterCntrWritePackage(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, _, _, cntrWritePackage, _, _, _, _, _, _, _, _, _, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(cntrWritePackage), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrCorrected(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, _, _, _, errCorrected, _, _, _, _, _, _, _, _, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(errCorrected), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrCRC24(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, _, _, _, _, errCRC24, _, _, _, _, _, _, _, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(errCRC24), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrCRC32(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, _, _, _, _, _, errCRC32, _, _, _, _, _, _, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(errCRC32), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrEccDouble(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, _, _, _, _, _, _, errEccDouble, _, _, _, _, _, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(errEccDouble), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrFatal(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, _, _, _, _, _, _, _, errFatal, _, _, _, _, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(errFatal), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrReplay(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, _, _, _, _, _, _, _, _, errReplay, _, _, _, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(errReplay), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterErrUncorrected(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, _, _, _, _, _, _, _, _, _, errUncorrected, _, _, err := c.client.GetDeviceMLULinkCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(errUncorrected), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkPortMode(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			mode, err := c.client.GetDeviceMLULinkPortMode(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkPortMode", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(mode), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkSpeedFormat(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, speedFormat, err := c.client.GetDeviceMLULinkSpeedInfo(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkSpeedInfo", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(speedFormat), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkSpeedValue(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			speedValue, _, err := c.client.GetDeviceMLULinkSpeedInfo(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkSpeedInfo", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(speedValue), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkStatusIsActive(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			var active int
			if stat.linkActive[i] {
				active = 1
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(active), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkStatusSerdesState(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			_, serdes, err := c.client.GetDeviceMLULinkStatus(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkStatus", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(serdes), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkVersion(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			major, minor, build, err := c.client.GetDeviceMLULinkVersion(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkVersion", stat.slot, i))
				continue
			}
			v := calcVersion(major, minor, build)
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i, linkVersion: v})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(1), labelValues...)
		}
	}
}

func (c *cndevCollector) collectNUMANodeID(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		numaNodeID, err := c.client.GetDeviceNUMANodeID(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceNUMANodeId", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(numaNodeID), labelValues...)
	}
}

func (c *cndevCollector) collectParityError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		switch strings.Split(strings.ToLower(stat.model), "-")[0] {
		case "mlu220", "mlu270", "mlu290", "mlu370":
			return
		}
		parityError, err := c.client.GetDeviceParityError(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceParityError", stat.slot))
			break
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(parityError), labelValues...)
	}
}

func (c *cndevCollector) collectPCIeInfo(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		slotID, subsystemID, deviceID, vendor, subsystemVendor, domain, bus, device, function, err := c.client.GetDevicePCIeInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDevicePCIeInfo", stat.slot))
			continue
		}
		speed, width, err := c.client.GetDeviceCurrentPCIeInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCurrentPCIeInfo", stat.slot))
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
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, pcieInfo: pcieInfo})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, 1, labelValues...)
	}
}

func (c *cndevCollector) collectPCIeReadThroughput(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		read, _, err := c.client.GetDevicePCIeThroughput(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDevicePCIeThroughput", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(read), labelValues...)
	}
}

func (c *cndevCollector) collectPCIeReplayCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		count, err := c.client.GetDevicePCIeReplayCount(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDevicePCIeReplayCount", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(count), labelValues...)
	}
}

func (c *cndevCollector) collectPCIeWriteThroughput(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, write, err := c.client.GetDevicePCIeThroughput(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDevicePCIeThroughput", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(write), labelValues...)
	}
}

func (c *cndevCollector) collectPowerUsage(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		power, err := c.client.GetDevicePower(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDevicePower", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: ""})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(power), labelValues...)

		if c.baseInfo.num > 0 {
			for i := uint(0); i < c.baseInfo.num; i++ {
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", i+1)})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(power), labelValues...)
			}
			continue
		}

		if stat.mimEnabled {
			for _, info := range stat.mimInfos {
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", info.InstanceID)})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(power), labelValues...)
			}
			continue
		}

		if stat.smluEnabled {
			var smluInfos []cndev.SmluInfo
			smluInfos, err = c.client.GetAllSMluInfo(stat.slot)
			if err != nil {
				log.Warn(errors.Wrap(err, "GetAllSMluInfo"))
				continue
			}
			stat.smluInfos = smluInfos
			for _, info := range smluInfos {
				var labelValues []string
				if c.baseInfo.mode == "dynamic-smlu" {
					labelValues = getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, typ: "vcore", vf: fmt.Sprintf("%d", info.InstanceID)})
				} else {
					labelValues = getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", info.InstanceID)})
				}
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(power), labelValues...)
			}
			continue
		}

		state, err := c.client.GetDeviceVfState(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceVfState", stat.slot))
			continue
		}
		vf := 0
		for state != 0 {
			vf++
			if (state & 0x1) != 0 {
				power, err := c.client.GetDevicePower(uint(vf<<8 | int(stat.slot)))
				if err != nil {
					log.Errorln(errors.Wrapf(err, "Slot %d VF %d GetDevicePower", stat.slot, vf))
					continue
				}
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", vf)})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(power), labelValues...)
			}
			state >>= 1
		}
	}
}

func (c *cndevCollector) collectProcessIpuUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.processUtilDisabled {
			continue
		}
		pid, ipu, _, _, _, _, err := c.client.GetDeviceProcessUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", stat.slot))
			continue
		}
		for i, util := range ipu {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, pid: pid[i]})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectProcessJpuUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.processUtilDisabled {
			continue
		}
		pid, _, jpu, _, _, _, err := c.client.GetDeviceProcessUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", stat.slot))
			continue
		}
		for i, util := range jpu {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, pid: pid[i]})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectProcessMemoryUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.processUtilDisabled {
			continue
		}
		pid, _, _, mem, _, _, err := c.client.GetDeviceProcessUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", stat.slot))
			continue
		}
		for i, util := range mem {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, pid: pid[i]})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectProcessVpuDecodeUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.processUtilDisabled {
			continue
		}
		pid, _, _, _, vpuDecode, _, err := c.client.GetDeviceProcessUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", stat.slot))
			continue
		}
		for i, util := range vpuDecode {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, pid: pid[i]})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectProcessVpuEncodeUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.processUtilDisabled {
			continue
		}
		pid, _, _, _, _, vpuEncode, err := c.client.GetDeviceProcessUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceProcessUtil", stat.slot))
			continue
		}
		for i, util := range vpuEncode {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, pid: pid[i]})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectRemappedCorrectRows(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.remappedRowsDisabled {
			continue
		}
		correctRows, _, _, _, err := c.client.GetDeviceRemappedRows(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRemappedRows", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(correctRows), labelValues...)
	}
}

func (c *cndevCollector) collectRemappedFailedRows(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.remappedRowsDisabled {
			continue
		}
		_, _, _, failedRows, err := c.client.GetDeviceRemappedRows(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRemappedRows", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(failedRows), labelValues...)
	}
}

func (c *cndevCollector) collectRemappedPendingRows(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.remappedRowsDisabled {
			continue
		}
		_, _, pendingRows, _, err := c.client.GetDeviceRemappedRows(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRemappedRows", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(pendingRows), labelValues...)
	}
}

func (c *cndevCollector) collectRemappedUncorrectRows(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.remappedRowsDisabled {
			continue
		}
		_, uncorrectRows, _, _, err := c.client.GetDeviceRemappedRows(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRemappedRows", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(uncorrectRows), labelValues...)
	}
}

func (c *cndevCollector) collectSingleBitRetiredPageCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.retiredPageDisabled {
			continue
		}
		cause, pageCount, err := c.client.GetDeviceRetiredPageInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRetiredPageInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		if cause != 0 {
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, 0, labelValues...)
			continue
		}
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(pageCount), labelValues...)
	}
}

func (c *cndevCollector) collectSramEccDbeCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.memEccDisabled {
			continue
		}
		_, sramEccDbeCounter, _, _, _, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(sramEccDbeCounter), labelValues...)
	}
}

func (c *cndevCollector) collectSramEccParityCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.memEccDisabled {
			continue
		}
		_, _, sramEccParityCounter, _, _, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(sramEccParityCounter), labelValues...)
	}
}

func (c *cndevCollector) collectSramEccSbeCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.memEccDisabled {
			continue
		}
		sramEccSbeCounter, _, _, _, _, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(sramEccSbeCounter), labelValues...)
	}
}

func (c *cndevCollector) collectTemperature(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		board, _, _, _, _, err := c.client.GetDeviceTemperature(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceTemperature", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(board), labelValues...)
	}
}

func (c *cndevCollector) collectTinyCoreUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		vpuUtil, err := c.client.GetDeviceTinyCoreUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceTinyCoreUtil", stat.slot))
			continue
		}
		for core, util := range vpuUtil {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, core: core})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		board, coreUtil, err := c.client.GetDeviceUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceUtil", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: ""})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(board), labelValues...)

		if c.baseInfo.num > 0 {
			for i := uint(0); i < c.baseInfo.num; i++ {
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", i+1)})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(board), labelValues...)
			}
			continue
		}

		if stat.mimEnabled {
			for _, info := range stat.mimInfos {
				util := calcMimUtil(info, coreUtil)
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", info.InstanceID)})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
			}
			continue
		}

		if stat.smluEnabled {
			var smluInfos []cndev.SmluInfo
			smluInfos, err = c.client.GetAllSMluInfo(stat.slot)
			if err != nil {
				log.Warn(errors.Wrap(err, "GetAllSMluInfo"))
				continue
			}
			stat.smluInfos = smluInfos
			for _, info := range smluInfos {
				var labelValues []string
				if c.baseInfo.mode == "dynamic-smlu" {
					labelValues = getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, typ: "vcore", vf: fmt.Sprintf("%d", info.InstanceID)})
				} else {
					labelValues = getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", info.InstanceID)})
				}
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(info.IpuUtil), labelValues...)
			}
			continue
		}

		state, err := c.client.GetDeviceVfState(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceVfState", stat.slot))
			continue
		}
		vfNum := 0
		for state != 0 {
			vfNum++
			state >>= 1
		}
		for vfID := 1; vfID < vfNum+1; vfID++ {
			util, err := calcVFUtil(coreUtil, vfNum, vfID)
			if err != nil {
				log.Errorln(errors.Wrapf(err, "vf %d calcVFUtil, util %v, vfNum: %d", vfID, coreUtil, vfNum))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", vfID)})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, util, labelValues...)
		}
	}
}

func (c *cndevCollector) collectVersion(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, 1, labelValues...)
	}
}

func (c *cndevCollector) collectVideoCodecUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		vpuUtil, err := c.client.GetDeviceVideoCodecUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceVideoCodecUtil", stat.slot))
			continue
		}
		for core, util := range vpuUtil {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, core: core})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectVirtualFunctionMemUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		state, err := c.client.GetDeviceVfState(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceVfState", stat.slot))
			continue
		}
		vf := 0
		for state != 0 {
			vf++
			if (state & 0x1) != 0 {
				used, total, _, _, err := c.client.GetDeviceMemory(uint(vf<<8 | int(stat.slot)))
				if err != nil {
					log.Errorln(errors.Wrapf(err, "Slot %d VF %d GetDeviceMemory", stat.slot, vf))
					continue
				}
				util := float64(used) / float64(total) * 100
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", vf)})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, util, labelValues...)
			}
			state >>= 1
		}
	}
}

func (c *cndevCollector) collectVirtualFunctionPowerUsage(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		state, err := c.client.GetDeviceVfState(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceVfState", stat.slot))
			continue
		}
		vf := 0
		for state != 0 {
			vf++
			if (state & 0x1) != 0 {
				power, err := c.client.GetDevicePower(uint(vf<<8 | int(stat.slot)))
				if err != nil {
					log.Errorln(errors.Wrapf(err, "Slot %d VF %d GetDevicePower", stat.slot, vf))
					continue
				}
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", vf)})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(power), labelValues...)
			}
			state >>= 1
		}
	}
}

func (c *cndevCollector) collectVirtualFunctionUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, coreUtil, err := c.client.GetDeviceUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceUtil", stat.slot))
			continue
		}
		if stat.mimEnabled {
			for _, info := range stat.mimInfos {
				util := calcMimUtil(info, coreUtil)
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", info.InstanceID)})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
			}
			continue
		}

		state, err := c.client.GetDeviceVfState(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceVfState", stat.slot))
			continue
		}
		vfNum := 0
		for state != 0 {
			vfNum++
			state >>= 1
		}
		for vfID := 1; vfID < vfNum+1; vfID++ {
			util, err := calcVFUtil(coreUtil, vfNum, vfID)
			if err != nil {
				log.Errorln(errors.Wrapf(err, "vf %d calcVFUtil, util %v, vfNum: %d", vfID, coreUtil, vfNum))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", vfID)})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, util, labelValues...)
		}
	}
}

func (c *cndevCollector) collectVirtualMemTotal(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, _, _, total, err := c.client.GetDeviceMemory(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(total*1024), labelValues...)
	}
}

func (c *cndevCollector) collectVirtualMemUsed(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, _, used, _, err := c.client.GetDeviceMemory(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(used*1024), labelValues...)
	}
}

func (c *cndevCollector) collectLastXIDError(ch chan<- prometheus.Metric, m metrics.Metric) {
	slotInfo := map[uint]MLUStat{}
	for _, stat := range c.sharedInfo {
		slotInfo[stat.slot] = stat
	}

	for info, xid := range c.lastXID {
		labelValues := getLabelValues(m.Labels, labelInfo{
			stat: slotInfo[info.Device],
			host: c.baseInfo.host,
			xidInfo: cndev.XIDInfoWithTimestamp{
				XIDInfo: cndev.XIDInfo{
					DeviceInfo: info,
				},
			},
		})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(xid), labelValues...)
	}
}

func (c *cndevCollector) collectXIDErrorCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	slotInfo := map[uint]MLUStat{}
	for _, stat := range c.sharedInfo {
		slotInfo[stat.slot] = stat
	}

	for info, count := range c.mluXIDCounter {
		labelValues := getLabelValues(m.Labels, labelInfo{
			stat: slotInfo[info.Device],
			host: c.baseInfo.host,
			xidInfo: cndev.XIDInfoWithTimestamp{
				XIDInfo: info,
			},
		})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(count), labelValues...)
	}
}

func (c *cndevCollector) waitXIDEvents() error {
	slots := []int{}
	for _, stat := range c.sharedInfo {
		if stat.xidCallbackDisabled {
			continue
		}
		slots = append(slots, int(stat.slot))
	}
	if len(slots) == 0 {
		return nil
	}
	sort.Ints(slots)

	ch := make(chan cndev.XIDInfoWithTimestamp, 10)
	if err := c.client.RegisterEventsHandleAndWait(slots, ch); err != nil {
		log.Errorln(errors.Wrap(err, "register event handle"))
		return err
	}

	for event := range ch {
		c.mluXIDCounter[event.XIDInfo]++
		c.lastXID[event.DeviceInfo] = event.XID
	}

	return nil
}

func calcVFUtil(utils []int, vfNum int, vfID int) (float64, error) {
	sum := 0
	offset := 0
	switch vfNum {
	case 1, 2, 3, 4, 8:
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
	log.Debugf("PCIeID: %v:%v:%v.%v", domainStr, busStr, deviceStr, functionStr)
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
	log.Debugf("PCIeSpeed: %v", pcieSpeed)
	return pcieSpeed
}

func calcMimUtil(info cndev.MimInfo, utils []int) float64 {
	sum := 0
	offset := info.PlacementSize
	for i := 0; i < offset; i++ {
		sum += utils[info.PlacementStart+i]
	}
	return float64(sum) / float64(offset)
}
