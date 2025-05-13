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
	lastXID       map[cndev.DeviceInfo]cndev.XIDInfo
}

func NewCndevCollector(m metrics.CollectorMetrics, bi BaseInfo) Collector {
	c := &cndevCollector{
		baseInfo:      bi,
		client:        cndev.NewCndevClient(),
		metrics:       m,
		mluXIDCounter: map[cndev.XIDInfo]int{},
		lastXID:       map[cndev.DeviceInfo]cndev.XIDInfo{},
	}

	c.fnMap = map[string]interface{}{
		AddressSwapsCorrectCount:          c.collectAddressSwapsCorrectCount,
		AddressSwapsMax:                   c.collectAddressSwapsMax,
		AddressSwapsNone:                  c.collectAddressSwapsNone,
		AddressSwapsPartial:               c.collectAddressSwapsPartial,
		AddressSwapsUncorrectCount:        c.collectAddressSwapsUncorrectCount,
		AggregateDramEccDbeCount:          c.collectAggregateDramEccDbeCount,
		AggregateDramEccSbeCount:          c.collectAggregateDramEccSbeCount,
		AggregateSramEccDbeCount:          c.collectAggregateSramEccDbeCount,
		AggregateSramEccParityCount:       c.collectAggregateSramEccParityCount,
		AggregateSramEccSbeCount:          c.collectAggregateSramEccSbeCount,
		AggregateSramEccThresholdExceed:   c.collectAggregateSramEccThresholdExceed,
		BAR4MemoryTotal:                   c.collectBAR4MemoryTotal,
		BAR4MemoryUsed:                    c.collectBAR4MemoryUsed,
		BAR4MemoryFree:                    c.collectBAR4MemoryFree,
		ChipCPUUtil:                       c.collectChipCPUUtil,
		ChipCPUHi:                         c.collectChipCPUHi,
		ChipCPUNice:                       c.collectChipCPUNice,
		ChipCPUSi:                         c.collectChipCPUSi,
		ChipCPUSys:                        c.collectChipCPUSys,
		ChipCPUUser:                       c.collectChipCPUUser,
		ChipTemperature:                   c.collectChipTemperature,
		ClusterTemperature:                c.collectClusterTemperature,
		ComputeMode:                       c.collectComputeMode,
		CoreCurrent:                       c.collectCoreCurrent,
		CoreUtil:                          c.collectCoreUtil,
		CoreVoltage:                       c.collectCoreVoltage,
		D2DCRCError:                       c.collectD2DCRCError,
		D2DCRCErrorOverflow:               c.collectD2DCRCErrorOverflow,
		DDRDataWidth:                      c.collectDDRDataWidth,
		DDRBandWidth:                      c.collectDDRBandWidth,
		DDRFrequency:                      c.collectDDRFrequency,
		DeviceOsMemTotal:                  c.collectDeviceOsMemTotal,
		DeviceOsMemUsed:                   c.collectDeviceOsMemUsed,
		DoubleBitRetiredPageCount:         c.collectDoubleBitRetiredPageCount,
		DramEccDbeCount:                   c.collectDramEccDbeCount,
		DramEccSbeCount:                   c.collectDramEccSbeCount,
		DriverState:                       c.collectDriverState,
		ECCAddressForbiddenError:          c.collectECCAddressForbiddenError,
		ECCCorrectedError:                 c.collectECCCorrectedError,
		ECCCurrentMode:                    c.collectECCCurrentMode,
		ECCMultipleError:                  c.collectECCMultipleError,
		ECCMultipleMultipleError:          c.collectECCMultipleMultipleError,
		ECCMultipleOneError:               c.collectECCMultipleOneError,
		ECCOneBitError:                    c.collectECCOneBitError,
		ECCPendingMode:                    c.collectECCPendingMode,
		ECCTotalError:                     c.collectECCTotalError,
		ECCUncorrectedError:               c.collectECCUncorrectedError,
		FanSpeed:                          c.collectFanSpeed,
		Frequency:                         c.collectFrequency,
		FrequencyStatus:                   c.collectFrequencyStatus,
		Health:                            c.collectHealth,
		HeartbeatCount:                    c.collectHeartbeatCount,
		ImageCodecUtil:                    c.collectImageCodecUtil,
		MachineMLUNums:                    c.collectMachineMLUNums,
		MemCurrent:                        c.collectMemCurrent,
		MemFree:                           c.collectMemFree,
		MemTemperature:                    c.collectMemTemperature,
		MemTotal:                          c.collectMemTotal,
		MemUsed:                           c.collectMemUsed,
		MemUtil:                           c.collectMemUtil,
		MemVoltage:                        c.collectMemVoltage,
		MimInstanceInfo:                   c.collectMimInstanceInfo,
		MimProfileInfo:                    c.collectMimProfileInfo,
		MimProfileTotalCapacity:           c.collectMimProfileTotalCapacity,
		MLURDMADeviceInfo:                 c.collectMLURDMADeviceInfo,
		MLULinkCapabilityP2PTransfer:      c.collectMLULinkCapabilityP2PTransfer,
		MLULinkCapabilityInterlakenSerdes: c.collectMLULinkCapabilityInterlakenSerdes,
		MLULinkConnectType:                c.collectMLULinkConnectType,
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
		MLULinkCounterIllegalAccess:       c.collectMLULinkCounterIllegalAccess,
		MLULinkCounterLinkDown:            c.collectMLULinkCounterLinkDown,
		MLULinkPortMode:                   c.collectMLULinkPortMode,
		MLULinkPortNumber:                 c.collectMLULinkPortNumber,
		MLULinkRemoteInfo:                 c.collectMLULinkRemoteInfo,
		MLULinkSpeedFormat:                c.collectMLULinkSpeedFormat,
		MLULinkSpeedValue:                 c.collectMLULinkSpeedValue,
		MLULinkStatusCableState:           c.collectMLULinkStatusCableState,
		MLULinkStatusIsActive:             c.collectMLULinkStatusIsActive,
		MLULinkStatusSerdesState:          c.collectMLULinkStatusSerdesState,
		MLULinkVersion:                    c.collectMLULinkVersion,
		NUMANodeID:                        c.collectNUMANodeID,
		OpticalIsPresent:                  c.collectOpticalIsPresent,
		OpticalRxpwr:                      c.collectOpticalRxpwr,
		OpticalTemp:                       c.collectOpticalTemp,
		OpticalTxpwr:                      c.collectOpticalTxpwr,
		OpticalVolt:                       c.collectOpticalVolt,
		OverTempPowerOffCounter:           c.collectOverTempPowerOffCounter,
		OverTempPowerOffThreshold:         c.collectOverTempPowerOffThreshold,
		OverTempUnderClockCounter:         c.collectOverTempUnderClockCounter,
		OverTempUnderClockThreshold:       c.collectOverTempUnderClockThreshold,
		ParityError:                       c.collectParityError,
		PCIeCurrentSpeed:                  c.collectPCIeCurrentSpeed,
		PCIeCurrentWidth:                  c.collectPCIeCurrentWidth,
		PCIeInfo:                          c.collectPCIeInfo,
		PCIeMaxSpeed:                      c.collectePCIeMaxSpeed,
		PCIeMaxWidth:                      c.collectePCIeMaxWidth,
		PCIeReadThroughput:                c.collectPCIeReadThroughput,
		PCIeReplayCount:                   c.collectPCIeReplayCount,
		PCIeWriteThroughput:               c.collectPCIeWriteThroughput,
		PowerCurrentLimit:                 c.collectPowerCurrentLimit,
		PowerDefaultLimit:                 c.collectPowerDefaultLimit,
		PowerMaxLimit:                     c.collectPowerMaxLimit,
		PowerMinLimit:                     c.collectPowerMinLimit,
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
		RepairRetireStatus:                c.collectRepairRetireStatus,
		RetiredPagesOperationStatus:       c.collectRetiredPagesOperationStatus,
		SingleBitRetiredPageCount:         c.collectSingleBitRetiredPageCount,
		SmluInstanceInfo:                  c.collectSmluInstanceInfo,
		SmluProfileInfo:                   c.collectSmluProfileInfo,
		SmluProfileTotalCapacity:          c.collectSmluProfileTotalCapacity,
		SocCurrent:                        c.collectSocCurrent,
		SocVoltage:                        c.collectSocVoltage,
		SramEccDbeCount:                   c.collectSramEccDbeCount,
		SramEccParityCount:                c.collectSramEccParityCount,
		SramEccSbeCount:                   c.collectSramEccSbeCount,
		Temperature:                       c.collectTemperature,
		TensorUtil:                        c.collectTensorUtil,
		ThermalSlowdown:                   c.collectThermalSlowdown,
		TinyCoreUtil:                      c.collectTinyCoreUtil,
		Util:                              c.collectUtil,
		Version:                           c.collectVersion,
		VideoCodecUtil:                    c.collectVideoCodecUtil,
		VideoDecoderUtil:                  c.collectVideoDecoderUtil,
		VirtualFunctionMemUtil:            c.collectVirtualFunctionMemUtil,
		VirtualFunctionPowerUsage:         c.collectVirtualFunctionPowerUsage,
		VirtualFunctionUtil:               c.collectVirtualFunctionUtil,
		VirtualMemFree:                    c.collectVirtualMemFree,
		VirtualMemTotal:                   c.collectVirtualMemTotal,
		VirtualMemUsed:                    c.collectVirtualMemUsed,
		VirtualMode:                       c.collectMLUMode,
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

func (c *cndevCollector) collectAddressSwapsCorrectCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["addressSwapsDisabled"] {
			continue
		}
		correctCount, _, _, _, _, err := c.client.GetDeviceAddressSwaps(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRetiredPageInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(correctCount), labelValues...)
	}
}

func (c *cndevCollector) collectAddressSwapsMax(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["addressSwapsDisabled"] {
			continue
		}
		_, _, _, _, maxAddressSwaps, err := c.client.GetDeviceAddressSwaps(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRetiredPageInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(maxAddressSwaps), labelValues...)
	}
}

func (c *cndevCollector) collectAddressSwapsNone(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["addressSwapsDisabled"] {
			continue
		}
		_, _, none, _, _, err := c.client.GetDeviceAddressSwaps(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRetiredPageInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(none), labelValues...)
	}
}

func (c *cndevCollector) collectAddressSwapsPartial(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["addressSwapsDisabled"] {
			continue
		}
		_, _, _, partial, _, err := c.client.GetDeviceAddressSwaps(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRetiredPageInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(partial), labelValues...)
	}
}

func (c *cndevCollector) collectAddressSwapsUncorrectCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["addressSwapsDisabled"] {
			continue
		}
		_, uncorrectCount, _, _, _, err := c.client.GetDeviceAddressSwaps(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRetiredPageInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(uncorrectCount), labelValues...)
	}
}

func (c *cndevCollector) collectAggregateDramEccDbeCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["memEccDisabled"] {
			continue
		}
		_, aggregate, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(aggregate[1]), labelValues...)
	}
}

func (c *cndevCollector) collectAggregateDramEccSbeCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["memEccDisabled"] {
			continue
		}
		_, aggregate, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(aggregate[0]), labelValues...)
	}
}

func (c *cndevCollector) collectAggregateSramEccDbeCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["memEccDisabled"] {
			continue
		}
		_, aggregate, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(aggregate[3]), labelValues...)
	}
}

func (c *cndevCollector) collectAggregateSramEccParityCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["memEccDisabled"] {
			continue
		}
		_, aggregate, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(aggregate[4]), labelValues...)
	}
}

func (c *cndevCollector) collectAggregateSramEccSbeCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["memEccDisabled"] {
			continue
		}
		_, aggregate, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(aggregate[2]), labelValues...)
	}
}

func (c *cndevCollector) collectAggregateSramEccThresholdExceed(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["memEccDisabled"] {
			continue
		}
		_, aggregate, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(aggregate[5]), labelValues...)
	}
}

func (c *cndevCollector) collectBAR4MemoryTotal(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["bar4MemoryInfoDisabled"] {
			continue
		}
		total, _, _, err := c.client.GetDeviceBAR4MemoryInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceBAR4MemoryInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(total), labelValues...)
	}
}

func (c *cndevCollector) collectBAR4MemoryUsed(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["bar4MemoryInfoDisabled"] {
			continue
		}
		_, used, _, err := c.client.GetDeviceBAR4MemoryInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceBAR4MemoryInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(used), labelValues...)
	}
}

func (c *cndevCollector) collectBAR4MemoryFree(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["bar4MemoryInfoDisabled"] {
			continue
		}
		_, _, free, err := c.client.GetDeviceBAR4MemoryInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceBAR4MemoryInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(free), labelValues...)
	}
}

func (c *cndevCollector) collectChipCPUUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["cpuUtilDisabled"] {
			continue
		}
		board, coreUtil, _, _, _, _, _, err := c.client.GetDeviceCPUUtil(stat.slot)
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

func (c *cndevCollector) collectChipCPUHi(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["cpuUtilDisabled"] {
			continue
		}
		_, _, _, _, _, _, hi, err := c.client.GetDeviceCPUUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCPUUtil", stat.slot))
			continue
		}
		for core, util := range hi {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, cpuCore: fmt.Sprintf("%d", core)})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectChipCPUNice(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["cpuUtilDisabled"] {
			continue
		}
		_, _, _, nice, _, _, _, err := c.client.GetDeviceCPUUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCPUUtil", stat.slot))
			continue
		}
		for core, util := range nice {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, cpuCore: fmt.Sprintf("%d", core)})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectChipCPUSi(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["cpuUtilDisabled"] {
			continue
		}
		_, _, _, _, _, si, _, err := c.client.GetDeviceCPUUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCPUUtil", stat.slot))
			continue
		}
		for core, util := range si {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, cpuCore: fmt.Sprintf("%d", core)})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectChipCPUSys(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["cpuUtilDisabled"] {
			continue
		}
		_, _, _, _, sys, _, _, err := c.client.GetDeviceCPUUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCPUUtil", stat.slot))
			continue
		}
		for core, util := range sys {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, cpuCore: fmt.Sprintf("%d", core)})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
		}
	}
}

func (c *cndevCollector) collectChipCPUUser(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["cpuUtilDisabled"] {
			continue
		}
		_, _, user, _, _, _, _, err := c.client.GetDeviceCPUUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCPUUtil", stat.slot))
			continue
		}
		for core, util := range user {
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
		if chip < 0 {
			log.Debugf("Slot %d chip temperature %d is invalid", stat.slot, chip)
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

func (c *cndevCollector) collectComputeMode(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["computeModeDisabled"] {
			continue
		}
		mode, err := c.client.GetDeviceComputeMode(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceComputeMode", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(mode), labelValues...)
	}
}

func (c *cndevCollector) collectCoreCurrent(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["currentInfoDisabled"] {
			continue
		}
		core, _, _, err := c.client.GetDeviceCurrentInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCurrentInfo", stat.slot))
			continue
		}
		if core < 0 {
			log.Debugf("Slot %d core current %d is invalid", stat.slot, core)
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(core), labelValues...)
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

func (c *cndevCollector) collectCoreVoltage(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["voltageInfoDisabled"] {
			continue
		}
		core, _, _, err := c.client.GetDeviceVoltageInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceVoltageInfo", stat.slot))
			continue
		}
		if core < 0 {
			log.Debugf("Slot %d core voltage %d is invalid", stat.slot, core)
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(core), labelValues...)
	}
}

func (c *cndevCollector) collectD2DCRCError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["crcDisabled"] {
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
		if stat.cndevInterfaceDisabled["crcDisabled"] {
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

func (c *cndevCollector) collectDeviceOsMemTotal(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, total, err := c.client.GetDeviceOsMemory(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceArmOsMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(total*1024), labelValues...)
	}
}

func (c *cndevCollector) collectDeviceOsMemUsed(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		used, _, err := c.client.GetDeviceOsMemory(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceArmOsMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(used*1024), labelValues...)
	}
}

func (c *cndevCollector) collectDoubleBitRetiredPageCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["retiredPageDisabled"] {
			continue
		}
		_, pageCount, err := c.client.GetDeviceRetiredPageInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRetiredPageInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(pageCount), labelValues...)
	}
}

func (c *cndevCollector) collectDramEccDbeCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["memEccDisabled"] {
			continue
		}
		volatile, _, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(volatile[4]), labelValues...)
	}
}

func (c *cndevCollector) collectDramEccSbeCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["memEccDisabled"] {
			continue
		}
		volatile, _, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(volatile[3]), labelValues...)
	}
}

func (c *cndevCollector) collectDriverState(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		var state int
		_, _, isRunning, err := c.client.GetDeviceHealth(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceHealth", stat.slot))
			continue
		}
		if isRunning {
			state = 1
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(state), labelValues...)
	}
}

func (c *cndevCollector) collectECCAddressForbiddenError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["eccDisabled"] {
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
		if stat.cndevInterfaceDisabled["eccDisabled"] {
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

func (c *cndevCollector) collectECCCurrentMode(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["eccModeDisabled"] {
			continue
		}
		currentMode, _, err := c.client.GetDeviceEccMode(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceEccMode", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(currentMode), labelValues...)
	}
}

func (c *cndevCollector) collectECCMultipleError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["eccDisabled"] {
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
		if stat.cndevInterfaceDisabled["eccDisabled"] {
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
		if stat.cndevInterfaceDisabled["eccDisabled"] {
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
		if stat.cndevInterfaceDisabled["eccDisabled"] {
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

func (c *cndevCollector) collectECCPendingMode(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["eccModeDisabled"] {
			continue
		}
		_, pendingMode, err := c.client.GetDeviceEccMode(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceEccMode", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(pendingMode), labelValues...)
	}
}

func (c *cndevCollector) collectECCTotalError(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["eccDisabled"] {
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
		if stat.cndevInterfaceDisabled["eccDisabled"] {
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

func (c *cndevCollector) collectFrequencyStatus(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["frequencyStatusDisabled"] {
			continue
		}
		stats, err := c.client.GetDeviceFrequencyStatus(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceFrequencyStatus", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(stats), labelValues...)
	}
}

func (c *cndevCollector) collectHealth(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		health, _, _, err := c.client.GetDeviceHealth(stat.slot)
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
		if stat.cndevInterfaceDisabled["heartBeatDisabled"] {
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

func (c *cndevCollector) collectMemCurrent(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["currentInfoDisabled"] {
			continue
		}
		_, _, mem, err := c.client.GetDeviceCurrentInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCurrentInfo", stat.slot))
			continue
		}
		if mem < 0 {
			log.Debugf("Slot %d memory current %d is invalid", stat.slot, mem)
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(mem), labelValues...)
	}
}

func (c *cndevCollector) collectMemFree(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		used, total, _, _, err := c.client.GetDeviceMemory(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64((total-used)*1024*1024), labelValues...)
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
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(total*1024*1024), labelValues...)
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
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(used*1024*1024), labelValues...)
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
					labelValues = getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, typ: "vmemory", vf: fmt.Sprintf("%d", info.InstanceInfo.InstanceID)})
				} else {
					labelValues = getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", info.InstanceInfo.InstanceID)})
				}
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(info.MemUsed)/float64(info.ProfileInfo.MemTotal)*100, labelValues...)
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

func (c *cndevCollector) collectMemVoltage(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["voltageInfoDisabled"] {
			continue
		}
		_, _, mem, err := c.client.GetDeviceVoltageInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceVoltageInfo", stat.slot))
			continue
		}
		if mem < 0 {
			log.Debugf("Slot %d memory voltage %d is invalid", stat.slot, mem)
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(mem), labelValues...)
	}
}

func (c *cndevCollector) collectMimInstanceInfo(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.mimEnabled {
			mimInfos, err := c.client.GetAllMLUInstanceInfo(stat.slot)
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d GetAllMLUInstanceInfo", stat.slot))
				continue
			}
			for _, info := range mimInfos {
				mimInstanceInfo := cndev.MimInstanceInfo{
					InstanceName: info.InstanceInfo.InstanceName,
					InstanceID:   info.InstanceInfo.InstanceID,
					UUID:         info.InstanceInfo.UUID,
				}
				mimProfileInfo := cndev.MimProfileInfo{
					GDMACount:    info.ProfileInfo.GDMACount,
					JPUCount:     info.ProfileInfo.JPUCount,
					MemorySize:   info.ProfileInfo.MemorySize,
					MLUCoreCount: info.ProfileInfo.MLUCoreCount,
					ProfileID:    info.ProfileInfo.ProfileID,
					ProfileName:  info.ProfileInfo.ProfileName,
					VPUCount:     info.ProfileInfo.VPUCount,
				}
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, mimInstanceInfo: mimInstanceInfo, mimProfileInfo: mimProfileInfo})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(1), labelValues...)
			}
		}
	}
}

func (c *cndevCollector) collectMimProfileInfo(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.mimEnabled {
			infos, err := c.client.GetDeviceMimProfileInfo(stat.slot)
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMimProfileInfo", stat.slot))
				continue
			}
			for _, info := range infos {
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, mimProfileInfo: info})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(1), labelValues...)
			}
		}
	}
}

func (c *cndevCollector) collectMimProfileTotalCapacity(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.smluEnabled {
			infos, err := c.client.GetDeviceMimProfileInfo(stat.slot)
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMimProfileInfo", stat.slot))
				continue
			}
			for _, info := range infos {
				count, err := c.client.GetDeviceMimProfileMaxInstanceCount(stat.slot, uint(info.ProfileID))
				if err != nil {
					log.Errorln(errors.Wrapf(err, "Slot %d ProfileID %d GetDeviceMimProfileMaxInstanceCount", stat.slot, info.ProfileID))
					continue
				}
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, mimProfileInfo: info})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(count), labelValues...)
			}
		}
	}
}

func (c *cndevCollector) collectMLURDMADeviceInfo(ch chan<- prometheus.Metric, m metrics.Metric) {
	if len(c.baseInfo.rdmaDevice) == 0 {
		return
	}
	for _, stat := range c.sharedInfo {
		_, _, _, _, _, domain, bus, device, function, err := c.client.GetDevicePCIeInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDevicePCIeInfo", stat.slot))
			continue
		}
		for _, rdmaDevice := range c.baseInfo.rdmaDevice {
			relation, err := c.client.GetTopologyRelationship(domain, bus, device, function, rdmaDevice.domain, rdmaDevice.bus, rdmaDevice.device, rdmaDevice.function)
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d RDMADevice %s GetTopologyRelationship", stat.slot, rdmaDevice.name))
				continue
			}
			// 2 indicates all devices that only need to traverse a single PCIe switch.
			if relation == 2 {
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, rdmaDevice: rdmaDevice})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, 1, labelValues...)
			}
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

func (c *cndevCollector) collectMLULinkConnectType(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["mluLinkRemoteInfoDisabled"] {
			continue
		}
		for i := 0; i < stat.link; i++ {
			_, _, _, _, connectType, _, _, _, _, _, err := c.client.GetDeviceMLULinkRemoteInfo(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMLULinkRemoteInfo", stat.slot))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(connectType), labelValues...)
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

func (c *cndevCollector) collectMLULinkCounterIllegalAccess(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["mluLinkErrorCounterDisabled"] {
			continue
		}
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			counter, err := c.client.GetDeviceMLULinkErrorCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkErrorCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(counter), labelValues...)
		}
	}
}

func (c *cndevCollector) collectMLULinkCounterLinkDown(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["mluLinkEventCounterDisabled"] {
			continue
		}
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			counter, err := c.client.GetDeviceMLULinkEventCounter(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkEventCounter", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(counter), labelValues...)
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

func (c *cndevCollector) collectMLULinkPortNumber(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		number := c.client.GetDeviceMLULinkPortNumber(stat.slot)
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(number), labelValues...)
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

func (c *cndevCollector) collectMLULinkRemoteInfo(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["mluLinkRemoteInfoDisabled"] {
			continue
		}
		for i := 0; i < stat.link; i++ {
			mcSn, baSn, slotID, portID, connectType, ncsUUID64, ip, mac, uuid, portName, err := c.client.GetDeviceMLULinkRemoteInfo(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMLULinkRemoteInfo", stat.slot))
				continue
			}
			mluLinkRemoteInfo := mluLinkRemoteInfo{
				remoteMac:      mac,
				remotePortName: portName,
			}
			if connectType == 0 {
				mluLinkRemoteInfo.remoteMcSn = fmt.Sprintf("%d", mcSn)
				mluLinkRemoteInfo.remoteBaSn = fmt.Sprintf("%d", baSn)
				mluLinkRemoteInfo.remoteSlotID = fmt.Sprintf("%d", slotID)
				mluLinkRemoteInfo.remotePortID = fmt.Sprintf("%d", portID)
				mluLinkRemoteInfo.remoteNcsUUID64 = fmt.Sprintf("%d", ncsUUID64)
				mluLinkRemoteInfo.remoteIP = ip
				mluLinkRemoteInfo.remoteUUID = uuid
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i, mluLinkRemoteInfo: mluLinkRemoteInfo})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, 1, labelValues...)
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

func (c *cndevCollector) collectMLULinkStatusCableState(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			if !stat.linkActive[i] {
				continue
			}
			_, _, cable, err := c.client.GetDeviceMLULinkStatus(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceMLULinkStatus", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(cable), labelValues...)
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
			if !stat.linkActive[i] {
				continue
			}
			_, serdes, _, err := c.client.GetDeviceMLULinkStatus(stat.slot, uint(i))
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

func (c *cndevCollector) collectMLUMode(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		var mluMode int
		if stat.mimEnabled {
			mluMode = 1
		} else if stat.smluEnabled {
			mluMode = 2
		} else if c.baseInfo.mode == "env-share" {
			mluMode = 3
		} else {
			mluMode = 0
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(mluMode), labelValues...)
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

func (c *cndevCollector) collectOpticalIsPresent(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i := 0; i < stat.link; i++ {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(stat.opticalPresent[i]), labelValues...)
		}
	}
}

func (c *cndevCollector) collectOpticalRxpwr(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i, present := range stat.opticalPresent {
			if present != uint8(1) {
				continue
			}
			_, _, _, _, rxpwrs, err := c.client.GetDeviceOpticalInfo(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceOpticalInfo", stat.slot, i))
				continue
			}

			for index, rxpwr := range rxpwrs {
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i, lane: index})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(rxpwr), labelValues...)
			}
		}
	}
}

func (c *cndevCollector) collectOpticalTemp(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i, present := range stat.opticalPresent {
			if present != uint8(1) {
				continue
			}
			_, temp, _, _, _, err := c.client.GetDeviceOpticalInfo(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceOpticalInfo", stat.slot, i))
				continue
			}
			if temp < 0 {
				log.Debugf("Slot %d, link %d optical temperature %f is invalid", stat.slot, i, temp)
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(temp), labelValues...)
		}
	}
}

func (c *cndevCollector) collectOpticalTxpwr(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i, present := range stat.opticalPresent {
			if present != uint8(1) {
				continue
			}
			_, _, _, txpwrs, _, err := c.client.GetDeviceOpticalInfo(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceOpticalInfo", stat.slot, i))
				continue
			}

			for index, txpwr := range txpwrs {
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i, lane: index})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(txpwr), labelValues...)
			}
		}
	}
}

func (c *cndevCollector) collectOpticalVolt(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		for i, present := range stat.opticalPresent {
			if present != uint8(1) {
				continue
			}
			_, _, volt, _, _, err := c.client.GetDeviceOpticalInfo(stat.slot, uint(i))
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d link %d GetDeviceOpticalInfo", stat.slot, i))
				continue
			}
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, link: i})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(volt), labelValues...)
		}
	}
}

func (c *cndevCollector) collectOverTempPowerOffCounter(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["overTemperatureInfoDisabled"] {
			continue
		}
		counter, _, err := c.client.GetDeviceOverTemperatureInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceOverTemperatureInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(counter), labelValues...)
	}
}

func (c *cndevCollector) collectOverTempPowerOffThreshold(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["overTemperatureThresholdDisabled"] {
			continue
		}
		threshold, err := c.client.GetDeviceOverTemperatureShutdownThreshold(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceOverTemperatureShutdownThreshold", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(threshold), labelValues...)
	}
}

func (c *cndevCollector) collectOverTempUnderClockCounter(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["overTemperatureInfoDisabled"] {
			continue
		}
		_, counter, err := c.client.GetDeviceOverTemperatureInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceOverTemperatureInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(counter), labelValues...)
	}
}

func (c *cndevCollector) collectOverTempUnderClockThreshold(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["overTemperatureThresholdDisabled"] {
			continue
		}
		threshold, err := c.client.GetDeviceOverTemperatureSlowdownThreshold(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceOverTemperatureSlowdownThreshold", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(threshold), labelValues...)
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

func (c *cndevCollector) collectPCIeCurrentSpeed(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		speed, _, err := c.client.GetDeviceCurrentPCIeInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCurrentPCIeInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(speed), labelValues...)
	}
}

func (c *cndevCollector) collectPCIeCurrentWidth(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, width, err := c.client.GetDeviceCurrentPCIeInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCurrentPCIeInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(width), labelValues...)
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

func (c *cndevCollector) collectePCIeMaxSpeed(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["maxPCIeInfoDisabled"] {
			continue
		}
		speed, _, err := c.client.GetDeviceMaxPCIeInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMaxPCIeInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(speed), labelValues...)
	}
}

func (c *cndevCollector) collectePCIeMaxWidth(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["maxPCIeInfoDisabled"] {
			continue
		}
		_, width, err := c.client.GetDeviceMaxPCIeInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMaxPCIeInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(width), labelValues...)
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

func (c *cndevCollector) collectPowerCurrentLimit(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["powerCurrentLimitDisabled"] {
			continue
		}
		currentLimit, err := c.client.GetDevicePowerManagementLimitation(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDevicePowerManagementDefaultLimitation", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(currentLimit), labelValues...)
	}
}

func (c *cndevCollector) collectPowerDefaultLimit(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["powerDefaultLimitDisabled"] {
			continue
		}
		defaultLimit, err := c.client.GetDevicePowerManagementDefaultLimitation(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDevicePowerManagementDefaultLimitation", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(defaultLimit), labelValues...)
	}
}

func (c *cndevCollector) collectPowerMaxLimit(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["powerLimitRangeDisabled"] {
			continue
		}
		_, maxLimit, err := c.client.GetDevicePowerManagementLimitRange(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDevicePowerManagementLimitRange", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(maxLimit), labelValues...)
	}
}

func (c *cndevCollector) collectPowerMinLimit(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["powerLimitRangeDisabled"] {
			continue
		}
		minLimit, _, err := c.client.GetDevicePowerManagementLimitRange(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDevicePowerManagementLimitRange", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(minLimit), labelValues...)
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
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", info.InstanceInfo.InstanceID)})
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
					labelValues = getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, typ: "vcore", vf: fmt.Sprintf("%d", info.InstanceInfo.InstanceID)})
				} else {
					labelValues = getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", info.InstanceInfo.InstanceID)})
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
		if stat.cndevInterfaceDisabled["processUtilDisabled"] {
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
		if stat.cndevInterfaceDisabled["processUtilDisabled"] {
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
		if stat.cndevInterfaceDisabled["processUtilDisabled"] {
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
		if stat.cndevInterfaceDisabled["processUtilDisabled"] {
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
		if stat.cndevInterfaceDisabled["processUtilDisabled"] {
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
		if stat.cndevInterfaceDisabled["remappedRowsDisabled"] {
			continue
		}
		correctRows, _, _, _, _, err := c.client.GetDeviceRemappedRows(stat.slot)
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
		if stat.cndevInterfaceDisabled["remappedRowsDisabled"] {
			continue
		}
		_, _, _, failedRows, _, err := c.client.GetDeviceRemappedRows(stat.slot)
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
		if stat.cndevInterfaceDisabled["remappedRowsDisabled"] {
			continue
		}
		_, _, pendingRows, _, _, err := c.client.GetDeviceRemappedRows(stat.slot)
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
		if stat.cndevInterfaceDisabled["remappedRowsDisabled"] {
			continue
		}
		_, uncorrectRows, _, _, _, err := c.client.GetDeviceRemappedRows(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRemappedRows", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(uncorrectRows), labelValues...)
	}
}

func (c *cndevCollector) collectRepairRetireStatus(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["remappedRowsDisabled"] {
			continue
		}
		_, _, _, _, retirePending, err := c.client.GetDeviceRemappedRows(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRemappedRows", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(retirePending), labelValues...)
	}
}

func (c *cndevCollector) collectRetiredPagesOperationStatus(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["retiredPagesOperationDisabled"] {
			continue
		}
		retirement, err := c.client.GetDeviceRetiredPagesOperation(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRetiredPagesOperation", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(retirement), labelValues...)
	}
}

func (c *cndevCollector) collectSingleBitRetiredPageCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["retiredPageDisabled"] {
			continue
		}
		pageCount, _, err := c.client.GetDeviceRetiredPageInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceRetiredPageInfo", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(pageCount), labelValues...)
	}
}

func (c *cndevCollector) collectSmluInstanceInfo(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.smluEnabled {
			smluInfos, err := c.client.GetAllSMluInfo(stat.slot)
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d GetAllSMluInfo", stat.slot))
				continue
			}
			for _, info := range smluInfos {
				smluInstanceInfo := cndev.SmluInstanceInfo{
					InstanceName: info.InstanceInfo.InstanceName,
					InstanceID:   info.InstanceInfo.InstanceID,
					UUID:         info.InstanceInfo.UUID,
				}
				smluProfileInfo := cndev.SmluProfileInfo{
					IpuTotal:    info.ProfileInfo.IpuTotal,
					MemTotal:    info.ProfileInfo.MemTotal,
					ProfileID:   info.ProfileInfo.ProfileID,
					ProfileName: info.ProfileInfo.ProfileName,
				}
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, smluInstanceInfo: smluInstanceInfo, smluProfileInfo: smluProfileInfo})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(1), labelValues...)
			}
		}
	}
}

func (c *cndevCollector) collectSmluProfileInfo(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.smluEnabled {
			info, err := c.client.GetDeviceSMluProfileIDInfo(stat.slot)
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceSMluProfileIdInfo", stat.slot))
				continue
			}
			for _, profileID := range info {
				smluProfileInfo, _, _, err := c.client.GetDeviceSMluProfileInfo(stat.slot, uint(profileID))
				if err != nil {
					log.Errorln(errors.Wrapf(err, "Slot %d profileId %d GetDeviceSMluProfileInfo", stat.slot, profileID))
					continue
				}
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, smluProfileInfo: smluProfileInfo})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(1), labelValues...)
			}
		}
	}
}

func (c *cndevCollector) collectSmluProfileTotalCapacity(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.smluEnabled {
			info, err := c.client.GetDeviceSMluProfileIDInfo(stat.slot)
			if err != nil {
				log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceSMluProfileIdInfo", stat.slot))
				continue
			}
			for _, profileID := range info {
				smluProfileInfo, total, _, err := c.client.GetDeviceSMluProfileInfo(stat.slot, uint(profileID))
				if err != nil {
					log.Errorln(errors.Wrapf(err, "Slot %d profileId %d GetDeviceSMluProfileInfo", stat.slot, profileID))
					continue
				}
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, smluProfileInfo: smluProfileInfo})
				ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(total), labelValues...)
			}
		}
	}
}

func (c *cndevCollector) collectSocCurrent(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["currentInfoDisabled"] {
			continue
		}
		_, soc, _, err := c.client.GetDeviceCurrentInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCurrentInfo", stat.slot))
			continue
		}
		if soc < 0 {
			log.Debugf("Slot %d soc current %d is invalid", stat.slot, soc)
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(soc), labelValues...)
	}
}

func (c *cndevCollector) collectSocVoltage(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["voltageInfoDisabled"] {
			continue
		}
		_, soc, _, err := c.client.GetDeviceVoltageInfo(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceVoltageInfo", stat.slot))
			continue
		}
		if soc < 0 {
			log.Debugf("Slot %d soc voltage %d is invalid", stat.slot, soc)
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(soc), labelValues...)
	}
}

func (c *cndevCollector) collectSramEccDbeCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["memEccDisabled"] {
			continue
		}
		volatile, _, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(volatile[1]), labelValues...)
	}
}

func (c *cndevCollector) collectSramEccParityCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["memEccDisabled"] {
			continue
		}
		volatile, _, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(volatile[2]), labelValues...)
	}
}

func (c *cndevCollector) collectSramEccSbeCount(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["memEccDisabled"] {
			continue
		}
		volatile, _, err := c.client.GetDeviceMemEccCounter(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemEccCounter", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.CounterValue, float64(volatile[0]), labelValues...)
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

func (c *cndevCollector) collectTensorUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["tensorUtilDisabled"] {
			continue
		}
		util, err := c.client.GetDeviceTensorUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceTensorUtil", stat.slot))
			continue
		}
		if util < 0 || util > 100 {
			log.Debugf("Slot %d tensor util %d is abnormal", stat.slot, util)
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(util), labelValues...)
	}
}

func (c *cndevCollector) collectThermalSlowdown(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		if stat.cndevInterfaceDisabled["performanceThrottleDisabled"] {
			continue
		}
		thermal, err := c.client.GetDevicePerformanceThrottleReason(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDevicePerformanceThrottleReason", stat.slot))
			continue
		}
		if thermal {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(1), labelValues...)
		} else {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(0), labelValues...)
		}
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
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", info.InstanceInfo.InstanceID)})
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
					labelValues = getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, typ: "vcore", vf: fmt.Sprintf("%d", info.InstanceInfo.InstanceID)})
				} else {
					labelValues = getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", info.InstanceInfo.InstanceID)})
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
		major, minor, build, err := c.client.GetDeviceCndevVersion()
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceCndevVersion", stat.slot))
			break
		}
		cndevVersion := calcVersion(major, minor, build)
		if stat.cndevInterfaceDisabled["computeCapabilityDisabled"] {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, cndevVersion: cndevVersion, computeCapability: ""})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, 1, labelValues...)
			continue
		}
		capMajor, capMinor, err := c.client.GetDeviceComputeCapability(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceComputeCapability", stat.slot))
			continue
		}
		computeCapability := fmt.Sprintf("v%d.%d", capMajor, capMinor)
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, cndevVersion: cndevVersion, computeCapability: computeCapability})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, 1, labelValues...)
	}
}

func (c *cndevCollector) collectVideoCodecUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		vpuUtil, _, err := c.client.GetDeviceVideoCodecUtil(stat.slot)
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

func (c *cndevCollector) collectVideoDecoderUtil(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, decoderUtil, err := c.client.GetDeviceVideoCodecUtil(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceVideoCodecUtil", stat.slot))
			continue
		}
		for core, util := range decoderUtil {
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
				labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host, vf: fmt.Sprintf("%d", info.InstanceInfo.InstanceID)})
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

func (c *cndevCollector) collectVirtualMemFree(ch chan<- prometheus.Metric, m metrics.Metric) {
	for _, stat := range c.sharedInfo {
		_, _, used, total, err := c.client.GetDeviceMemory(stat.slot)
		if err != nil {
			log.Errorln(errors.Wrapf(err, "Slot %d GetDeviceMemory", stat.slot))
			continue
		}
		labelValues := getLabelValues(m.Labels, labelInfo{stat: stat, host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64((total-used)*1024*1024), labelValues...)
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
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(total*1024*1024), labelValues...)
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
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(used*1024*1024), labelValues...)
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
					XID:        xid.XID,
				},
			},
		})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(xid.XIDBase10), labelValues...)
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
		if stat.cndevInterfaceDisabled["xidCallbackDisabled"] {
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
		c.lastXID[event.DeviceInfo] = event.XIDInfo
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
