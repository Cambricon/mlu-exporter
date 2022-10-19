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
	"strings"
	"testing"

	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/Cambricon/mlu-exporter/pkg/mock"
	"github.com/Cambricon/mlu-exporter/pkg/podresources"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestCollect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mluStat := map[string]mluStat{
		"uuid1": {
			slot:   0,
			model:  "MLU270-X5K",
			uuid:   "uuid1",
			sn:     "sn1",
			mcu:    "v1.1.1",
			driver: "v2.2.2",
		},
		"uuid2": {
			slot:   1,
			model:  "MLU270-X5K",
			uuid:   "uuid2",
			sn:     "sn2",
			mcu:    "v1.1.1",
			driver: "v2.2.2",
		},
		"uuid3": {
			slot:   2,
			model:  "MLU270-X5K",
			uuid:   "uuid3",
			sn:     "sn3",
			mcu:    "v1.1.1",
			driver: "v2.2.2",
		},
		"uuid4": {
			slot:   3,
			model:  "MLU270-X5K",
			uuid:   "uuid4",
			sn:     "sn4",
			mcu:    "v1.1.1",
			driver: "v2.2.2",
		},
	}

	node := "machine1"
	devicePodInfo := map[string]podresources.PodInfo{
		"MLU-uuid1": {
			Pod:       "pod1",
			Namespace: "namespace1",
			Container: "container1",
		},
		"MLU-uuid2-_-2": {
			Pod:       "pod2",
			Namespace: "namespace2",
			Container: "container2",
		},
		"MLU-uuid3--fake--4": {
			Pod:       "pod3",
			Namespace: "namespace3",
			Container: "container3",
		},
	}
	hostCPUTotal := float64(6185912)
	hostCPUIdle := float64(34459)
	hostMemTotal := float64(24421820)
	hostMemFree := float64(18470732)
	mluMetrics := map[string]struct {
		boardUtil                         int
		coreUtil                          []int
		memUsed                           int64
		memTotal                          int64
		virtualMemUsed                    int64
		virtualMemTotal                   int64
		virtualFunctionMemUsed            []int64
		virtualFunctionMemTotal           int64
		temperature                       int
		clusterTemperatures               []int
		health                            int
		power                             int
		pcieRead                          uint64
		pcieWrite                         uint64
		dramRead                          uint64
		dramWrite                         uint64
		mlulinkRead                       uint64
		mlulinkWrite                      uint64
		vfState                           int
		pcieInfo                          []string
		pcieSlotID                        int
		pcieSubsystemID                   uint
		pcieDeviceID                      uint
		pcieVendor                        uint16
		pcieSubsystemVendor               uint16
		pcieDomain                        uint
		pcieBus                           uint
		pcieDevice                        uint
		pcieFunction                      uint
		pcieSpeed                         int
		pcieWidth                         int
		pid                               []uint32
		processIpuUtil                    []uint32
		processJpuUtil                    []uint32
		processMemUtil                    []uint32
		processVpuDecodeUtil              []uint32
		processVpuEncodeUtil              []uint32
		chipCPUUtil                       uint16
		chipCPUCoreUtil                   []uint8
		armOsMemUsed                      int64
		armOsMemTotal                     int64
		videoCodecUtil                    []int
		imageCodecUtil                    []int
		tinyCoreUtil                      []int
		numaNodeID                        int
		ddrDataWidth                      int
		ddrBandWidth                      float64
		eccAddressForbiddenError          uint64
		eccCorrectedError                 uint64
		eccMultipleError                  uint64
		eccMultipleMultipleError          uint64
		eccMultipleOneError               uint64
		eccOneBitError                    uint64
		eccTotalError                     uint64
		eccUncorrectedError               uint64
		mluLinkCapabilityP2PTransfer      []uint
		mluLinkCapabilityInterlakenSerdes []uint
		mluLinkCounterCntrReadByte        []uint64
		mluLinkCounterCntrReadPackage     []uint64
		mluLinkCounterCntrWriteByte       []uint64
		mluLinkCounterCntrWritePackage    []uint64
		mluLinkCounterErrCorrected        []uint64
		mluLinkCounterErrCRC24            []uint64
		mluLinkCounterErrCRC32            []uint64
		mluLinkCounterErrEccDouble        []uint64
		mluLinkCounterErrFatal            []uint64
		mluLinkCounterErrReplay           []uint64
		mluLinkCounterErrUncorrected      []uint64
		mluLinkPortMode                   []int
		mluLinkPortNumber                 int
		mluLinkSpeedFormat                []int
		mluLinkSpeedValue                 []float32
		mluLinkStatusIsActive             []int
		mluLinkStatusSerdesState          []int
		mluLinkMajor                      []uint
		mluLinkMinor                      []uint
		mluLinkBuild                      []uint
		d2dCRCError                       uint64
		d2dCRCErrorOverflow               uint64
	}{
		"uuid1": {
			11,
			[]int{11, 11, 11, 11, 13, 13, 13, 13, 15, 15, 15, 15, 17, 17, 17, 17},
			21,
			1000,
			100,
			1000,
			[]int64{},
			1000,
			31,
			[]int{22, 23, 24, 26},
			1,
			41,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			0,
			[]string{"0", "12", "270", "cabc", "cabc", "0000:1a:00.0", "16GT/s", "4"},
			0,
			0x12,
			0x270,
			0xcabc,
			0xcabc,
			0,
			26,
			0,
			0,
			4,
			4,
			[]uint32{1100, 1101},
			[]uint32{1, 2},
			[]uint32{2, 4},
			[]uint32{3, 6},
			[]uint32{4, 8},
			[]uint32{5, 10},
			25,
			[]uint8{0, 0, 100, 0},
			1024,
			10240,
			[]int{0, 1, 0, 1, 0, 1},
			[]int{0, 1, 0, 1, 0, 1},
			[]int{0, 1},
			0,
			1,
			3.1415926535,
			100,
			0,
			1000,
			10,
			100000,
			19,
			20,
			5,
			[]uint{1, 1},
			[]uint{1, 1},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]int{1, 1},
			2,
			[]int{1, 1},
			[]float32{12, 13},
			[]int{1, 1},
			[]int{1, 1},
			[]uint{1, 1},
			[]uint{1, 1},
			[]uint{1, 1},
			128,
			256,
		},
		"uuid2": {
			11,
			[]int{11, 11, 11, 11, 13, 13, 13, 13, 15, 15, 15, 15, 17, 17, 17, 17},
			21,
			1000,
			100,
			1000,
			[]int64{},
			1000,
			31,
			[]int{22, 23, 24, 26},
			1,
			41,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			0,
			[]string{"1", "12", "270", "cabc", "cabc", "0000:1b:00.0", "16GT/s", "4"},
			1,
			0x12,
			0x270,
			0xcabc,
			0xcabc,
			0,
			27,
			0,
			0,
			4,
			4,
			[]uint32{1103, 1104},
			[]uint32{1, 2},
			[]uint32{2, 4},
			[]uint32{3, 6},
			[]uint32{4, 8},
			[]uint32{5, 10},
			25,
			[]uint8{0, 0, 100, 0},
			1024,
			10240,
			[]int{0, 1, 0, 1, 0, 1},
			[]int{0, 1, 0, 1, 0, 1},
			[]int{0, 1},
			0,
			1,
			3.1415926535,
			100,
			0,
			1000,
			10,
			100000,
			19,
			20,
			21,
			[]uint{1, 1},
			[]uint{1, 1},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]int{1, 1},
			2,
			[]int{1, 1},
			[]float32{12, 13},
			[]int{1, 1},
			[]int{1, 1},
			[]uint{1, 1},
			[]uint{1, 1},
			[]uint{1, 1},
			128,
			256,
		},
		"uuid3": {
			11,
			[]int{11, 11, 11, 11, 13, 13, 13, 13, 15, 15, 15, 15, 17, 17, 17, 17},
			21,
			1000,
			100,
			1000,
			[]int64{11, 13, 15, 17},
			250,
			31,
			[]int{22, 23, 24, 26},
			1,
			41,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			15,
			[]string{"2", "12", "270", "cabc", "cabc", "0000:1c:00.0", "16GT/s", "4"},
			2,
			0x12,
			0x270,
			0xcabc,
			0xcabc,
			0,
			28,
			0,
			0,
			4,
			4,
			[]uint32{1104, 1105},
			[]uint32{1, 2},
			[]uint32{2, 4},
			[]uint32{3, 6},
			[]uint32{4, 8},
			[]uint32{5, 10},
			25,
			[]uint8{0, 0, 100, 0},
			1024,
			10240,
			[]int{0, 1, 0, 1, 0, 1},
			[]int{0, 1, 0, 1, 0, 1},
			[]int{0, 1},
			1,
			1,
			3.1415926535,
			100,
			0,
			1000,
			10,
			100000,
			19,
			20,
			21,
			[]uint{1, 1},
			[]uint{1, 1},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]int{1, 1},
			2,
			[]int{1, 1},
			[]float32{12, 13},
			[]int{1, 1},
			[]int{1, 1},
			[]uint{1, 1},
			[]uint{1, 1},
			[]uint{1, 1},
			128,
			256,
		},
		"uuid4": {
			11,
			[]int{11, 11, 11, 11, 13, 13, 13, 13, 15, 15, 15, 15, 17, 17, 17, 17},
			21,
			1000,
			100,
			1000,
			[]int64{12, 16},
			500,
			31,
			[]int{22, 23, 24, 26},
			1,
			41,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			3,
			[]string{"3", "12", "270", "cabc", "cabc", "0000:1d:00.0", "16GT/s", "4"},
			3,
			0x12,
			0x270,
			0xcabc,
			0xcabc,
			0,
			29,
			0,
			0,
			4,
			4,
			[]uint32{1106, 1107},
			[]uint32{1, 2},
			[]uint32{2, 4},
			[]uint32{3, 6},
			[]uint32{4, 8},
			[]uint32{5, 10},
			25,
			[]uint8{0, 0, 100, 0},
			1024,
			10240,
			[]int{0, 1, 0, 1, 0, 1},
			[]int{0, 1, 0, 1, 0, 1},
			[]int{0, 1},
			1,
			1,
			3.1415926535,
			100,
			0,
			1000,
			10,
			100000,
			19,
			20,
			2,
			[]uint{1, 1},
			[]uint{1, 1},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]uint64{10000, 10000},
			[]int{1, 1},
			2,
			[]int{1, 1},
			[]float32{12, 13},
			[]int{1, 1},
			[]int{1, 1},
			[]uint{1, 1},
			[]uint{1, 1},
			[]uint{1, 1},
			128,
			256,
		},
	}

	cndev := mock.NewCndev(ctrl)
	cndev.EXPECT().Init().Return(nil).Times(1)
	cnpapi := mock.NewCnpapi(ctrl)
	cnpapi.EXPECT().Init().Return(nil).Times(1)
	for mlu, stat := range mluStat {
		cndev.EXPECT().GetDeviceTemperature(stat.slot).Return(mluMetrics[mlu].temperature, mluMetrics[mlu].clusterTemperatures, nil).Times(1)
		cndev.EXPECT().GetDeviceHealth(stat.slot).Return(mluMetrics[mlu].health, nil).Times(1)
		cndev.EXPECT().GetDeviceVfState(stat.slot).Return(mluMetrics[mlu].vfState, nil).Times(2)
		for vf := 0; vf < len(mluMetrics[stat.uuid].virtualFunctionMemUsed); vf++ {
			cndev.EXPECT().GetDeviceMemory(uint((vf+1)<<8|int(stat.slot))).Return(mluMetrics[mlu].virtualFunctionMemUsed[vf], mluMetrics[mlu].virtualFunctionMemTotal, mluMetrics[mlu].virtualMemUsed, mluMetrics[mlu].virtualMemTotal, nil).Times(1)
		}
		cndev.EXPECT().GetDeviceMemory(stat.slot).Return(mluMetrics[mlu].memUsed, mluMetrics[mlu].memTotal, mluMetrics[mlu].virtualMemUsed, mluMetrics[mlu].virtualMemTotal, nil).Times(5)
		cndev.EXPECT().GetDeviceUtil(stat.slot).Return(mluMetrics[mlu].boardUtil, mluMetrics[mlu].coreUtil, nil).Times(4)
		cndev.EXPECT().GetDeviceFanSpeed(stat.slot).Return(0, nil).Times(1)
		cndev.EXPECT().GetDevicePower(stat.slot).Return(mluMetrics[mlu].power, nil).Times(1)
		cndev.EXPECT().GetDevicePCIeInfo(stat.slot).Return(mluMetrics[mlu].pcieSlotID, mluMetrics[mlu].pcieSubsystemID, mluMetrics[mlu].pcieDeviceID, mluMetrics[mlu].pcieVendor, mluMetrics[mlu].pcieSubsystemVendor, mluMetrics[mlu].pcieDomain, mluMetrics[mlu].pcieBus, mluMetrics[mlu].pcieDevice, mluMetrics[mlu].pcieFunction, nil).Times(1)
		cndev.EXPECT().GetDeviceCurrentPCIeInfo(stat.slot).Return(mluMetrics[mlu].pcieSpeed, mluMetrics[mlu].pcieWidth, nil).Times(1)
		cndev.EXPECT().GetDeviceCPUUtil(stat.slot).Return(mluMetrics[mlu].chipCPUUtil, mluMetrics[mlu].chipCPUCoreUtil, nil).Times(1)
		cndev.EXPECT().GetDeviceArmOsMemory(stat.slot).Return(mluMetrics[mlu].armOsMemUsed, mluMetrics[mlu].armOsMemTotal, nil).Times(2)
		cndev.EXPECT().GetDeviceVideoCodecUtil(stat.slot).Return(mluMetrics[mlu].videoCodecUtil, nil).Times(1)
		cndev.EXPECT().GetDeviceImageCodecUtil(stat.slot).Return(mluMetrics[mlu].imageCodecUtil, nil).Times(1)
		cndev.EXPECT().GetDeviceTinyCoreUtil(stat.slot).Return(mluMetrics[mlu].tinyCoreUtil, nil).Times(1)
		cndev.EXPECT().GetDeviceNUMANodeID(stat.slot).Return(mluMetrics[mlu].numaNodeID, nil).Times(1)
		cndev.EXPECT().GetDeviceDDRInfo(stat.slot).Return(mluMetrics[mlu].ddrDataWidth, mluMetrics[mlu].ddrBandWidth, nil).Times(2)
		cndev.EXPECT().GetDeviceECCInfo(stat.slot).Return(mluMetrics[mlu].eccAddressForbiddenError, mluMetrics[mlu].eccCorrectedError, mluMetrics[mlu].eccMultipleError, mluMetrics[mlu].eccMultipleMultipleError, mluMetrics[mlu].eccMultipleOneError, mluMetrics[mlu].eccOneBitError, mluMetrics[mlu].eccTotalError, mluMetrics[mlu].eccUncorrectedError, nil).Times(8)
		cndev.EXPECT().GetDeviceMLULinkPortNumber(stat.slot).Return(mluMetrics[stat.uuid].mluLinkPortNumber).Times(19)
		for link := 0; link < mluMetrics[stat.uuid].mluLinkPortNumber; link++ {
			cndev.EXPECT().GetDeviceMLULinkCapability(stat.slot, uint(link)).Return(mluMetrics[mlu].mluLinkCapabilityP2PTransfer[link], mluMetrics[mlu].mluLinkCapabilityInterlakenSerdes[link], nil).Times(2)
			cndev.EXPECT().GetDeviceMLULinkCounter(stat.slot, uint(link)).Return(mluMetrics[mlu].mluLinkCounterCntrReadByte[link], mluMetrics[mlu].mluLinkCounterCntrReadPackage[link], mluMetrics[mlu].mluLinkCounterCntrWriteByte[link], mluMetrics[mlu].mluLinkCounterCntrWritePackage[link],
				mluMetrics[mlu].mluLinkCounterErrCorrected[link], mluMetrics[mlu].mluLinkCounterErrCRC24[link], mluMetrics[mlu].mluLinkCounterErrCRC32[link], mluMetrics[mlu].mluLinkCounterErrEccDouble[link], mluMetrics[mlu].mluLinkCounterErrFatal[link], mluMetrics[mlu].mluLinkCounterErrReplay[link], mluMetrics[mlu].mluLinkCounterErrUncorrected[link], nil).Times(11)
			cndev.EXPECT().GetDeviceMLULinkPortMode(stat.slot, uint(link)).Return(mluMetrics[mlu].mluLinkPortMode[link], nil).Times(1)
			cndev.EXPECT().GetDeviceMLULinkSpeedInfo(stat.slot, uint(link)).Return(mluMetrics[mlu].mluLinkSpeedValue[link], mluMetrics[mlu].mluLinkSpeedFormat[link], nil).Times(2)
			cndev.EXPECT().GetDeviceMLULinkStatus(stat.slot, uint(link)).Return(mluMetrics[mlu].mluLinkStatusIsActive[link], mluMetrics[mlu].mluLinkStatusSerdesState[link], nil).Times(2)
			cndev.EXPECT().GetDeviceMLULinkVersion(stat.slot, uint(link)).Return(mluMetrics[mlu].mluLinkMajor[link], mluMetrics[mlu].mluLinkMinor[link], mluMetrics[mlu].mluLinkBuild[link], nil).Times(1)
		}
		cndev.EXPECT().GetDeviceCRCInfo(stat.slot).Return(mluMetrics[mlu].d2dCRCError, mluMetrics[mlu].d2dCRCErrorOverflow, nil).Times(2)
		cndev.EXPECT().GetDeviceProcessUtil(stat.slot).Return(mluMetrics[mlu].pid, mluMetrics[mlu].processIpuUtil, mluMetrics[mlu].processJpuUtil, mluMetrics[mlu].processMemUtil, mluMetrics[mlu].processVpuDecodeUtil, mluMetrics[mlu].processVpuEncodeUtil, nil).Times(5)

		cnpapi.EXPECT().PmuFlushData(stat.slot).Return(nil).Times(1)
		cnpapi.EXPECT().GetPCIeReadBytes(stat.slot).Return(mluMetrics[mlu].pcieRead, nil).Times(1)
		cnpapi.EXPECT().GetPCIeWriteBytes(stat.slot).Return(mluMetrics[mlu].pcieWrite, nil).Times(1)
		cnpapi.EXPECT().GetDramReadBytes(stat.slot).Return(mluMetrics[mlu].dramRead, nil).Times(1)
		cnpapi.EXPECT().GetDramWriteBytes(stat.slot).Return(mluMetrics[mlu].dramWrite, nil).Times(1)
		cnpapi.EXPECT().GetMLULinkReadBytes(stat.slot).Return(mluMetrics[mlu].mlulinkRead, nil).Times(1)
		cnpapi.EXPECT().GetMLULinkWriteBytes(stat.slot).Return(mluMetrics[mlu].mlulinkWrite, nil).Times(1)
	}
	podresources := mock.NewPodResources(ctrl)
	podresources.EXPECT().GetDeviceToPodInfo().Return(devicePodInfo, nil).Times(1)
	host := mock.NewHost(ctrl)
	host.EXPECT().GetCPUStats().Return(hostCPUTotal, hostCPUIdle, nil).Times(2)
	host.EXPECT().GetMemoryStats().Return(hostMemTotal, hostMemFree, nil).Times(2)

	m := getMetrics(metrics.GetOrDie("../../examples/metrics.yaml"), "")

	cndevCollector := &cndevCollector{
		host:    node,
		cndev:   cndev,
		metrics: m["cndev"],
	}
	cndevCollector.fnMap = map[string]interface{}{
		ArmOsMemTotal:                     cndevCollector.collectArmOsMemTotal,
		ArmOsMemUsed:                      cndevCollector.collectArmOsMemUsed,
		Capacity:                          cndevCollector.collectCapacity,
		ChipCPUUtil:                       cndevCollector.collectChipCPUUtil,
		CoreUtil:                          cndevCollector.collectCoreUtil,
		D2DCRCError:                       cndevCollector.collectD2DCRCError,
		D2DCRCErrorOverflow:               cndevCollector.collectD2DCRCErrorOverflow,
		DDRDataWidth:                      cndevCollector.collectDDRDataWidth,
		DDRBandWidth:                      cndevCollector.collectDDRBandWidth,
		ECCAddressForbiddenError:          cndevCollector.collectECCAddressForbiddenError,
		ECCCorrectedError:                 cndevCollector.collectECCCorrectedError,
		ECCMultipleError:                  cndevCollector.collectECCMultipleError,
		ECCMultipleMultipleError:          cndevCollector.collectECCMultipleMultipleError,
		ECCMultipleOneError:               cndevCollector.collectECCMultipleOneError,
		ECCOneBitError:                    cndevCollector.collectECCOneBitError,
		ECCTotalError:                     cndevCollector.collectECCTotalError,
		ECCUncorrectedError:               cndevCollector.collectECCUncorrectedError,
		FanSpeed:                          cndevCollector.collectFanSpeed,
		Health:                            cndevCollector.collectHealth,
		ImageCodecUtil:                    cndevCollector.collectImageCodecUtil,
		MemTotal:                          cndevCollector.collectMemTotal,
		MemUsed:                           cndevCollector.collectMemUsed,
		MemUtil:                           cndevCollector.collectMemUtil,
		MLULinkCapabilityP2PTransfer:      cndevCollector.collectMLULinkCapabilityP2PTransfer,
		MLULinkCapabilityInterlakenSerdes: cndevCollector.collectMLULinkCapabilityInterlakenSerdes,
		MLULinkCounterCntrReadByte:        cndevCollector.collectMLULinkCounterCntrReadByte,
		MLULinkCounterCntrReadPackage:     cndevCollector.collectMLULinkCounterCntrReadPackage,
		MLULinkCounterCntrWriteByte:       cndevCollector.collectMLULinkCounterCntrWriteByte,
		MLULinkCounterCntrWritePackage:    cndevCollector.collectMLULinkCounterCntrWritePackage,
		MLULinkCounterErrCorrected:        cndevCollector.collectMLULinkCounterErrCorrected,
		MLULinkCounterErrCRC24:            cndevCollector.collectMLULinkCounterErrCRC24,
		MLULinkCounterErrCRC32:            cndevCollector.collectMLULinkCounterErrCRC32,
		MLULinkCounterErrEccDouble:        cndevCollector.collectMLULinkCounterErrEccDouble,
		MLULinkCounterErrFatal:            cndevCollector.collectMLULinkCounterErrFatal,
		MLULinkCounterErrReplay:           cndevCollector.collectMLULinkCounterErrReplay,
		MLULinkCounterErrUncorrected:      cndevCollector.collectMLULinkCounterErrUncorrected,
		MLULinkPortMode:                   cndevCollector.collectMLULinkPortMode,
		MLULinkSpeedFormat:                cndevCollector.collectMLULinkSpeedFormat,
		MLULinkSpeedValue:                 cndevCollector.collectMLULinkSpeedValue,
		MLULinkStatusIsActive:             cndevCollector.collectMLULinkStatusIsActive,
		MLULinkStatusSerdesState:          cndevCollector.collectMLULinkStatusSerdesState,
		MLULinkVersion:                    cndevCollector.collectMLULinkVersion,
		NUMANodeID:                        cndevCollector.collectNUMANodeID,
		PCIeInfo:                          cndevCollector.collectPCIeInfo,
		PowerUsage:                        cndevCollector.collectPowerUsage,
		ProcessIpuUtil:                    cndevCollector.collectProcessIpuUtil,
		ProcessJpuUtil:                    cndevCollector.collectProcessJpuUtil,
		ProcessMemoryUtil:                 cndevCollector.collectProcessMemoryUtil,
		ProcessVpuDecodeUtil:              cndevCollector.collectProcessVpuDecodeUtil,
		ProcessVpuEncodeUtil:              cndevCollector.collectProcessVpuEncodeUtil,
		Temperature:                       cndevCollector.collectTemperature,
		TinyCoreUtil:                      cndevCollector.collectTinyCoreUtil,
		Usage:                             cndevCollector.collectUsage,
		Util:                              cndevCollector.collectUtil,
		Version:                           cndevCollector.collectVersion,
		VideoCodecUtil:                    cndevCollector.collectVideoCodecUtil,
		VirtualFunctionMemUtil:            cndevCollector.collectVirtualFunctionMemUtil,
		VirtualFunctionUtil:               cndevCollector.collectVirtualFunctionUtil,
		VirtualMemTotal:                   cndevCollector.collectVirtualMemTotal,
		VirtualMemUsed:                    cndevCollector.collectVirtualMemUsed,
	}

	podResourcesCollector := &podResourcesCollector{
		metrics: m["podresources"],
		host:    node,
		client:  podresources,
	}

	podResourcesCollector.fnMap = map[string]interface{}{
		Allocated: podResourcesCollector.collectBoardAllocated,
		Container: podResourcesCollector.collectMLUContainer,
	}

	cnpapiCollector := &cnpapiCollector{
		metrics: m["cnpapi"],
		host:    node,
		cnpapi:  cnpapi,
	}

	cnpapiCollector.fnMap = map[string]interface{}{
		PCIeRead:     cnpapiCollector.collectPCIeRead,
		PCIeWrite:    cnpapiCollector.collectPCIeWrite,
		DramRead:     cnpapiCollector.collectDramRead,
		DramWrite:    cnpapiCollector.collectDramWrite,
		MLULinkRead:  cnpapiCollector.collectMLULinkRead,
		MLULinkWrite: cnpapiCollector.collectMLULinkWrite,
	}

	hostCollector := &hostCollector{
		metrics: m["host"],
		host:    node,
		client:  host,
	}

	hostCollector.fnMap = map[string]interface{}{
		HostCPUIdle:  hostCollector.collectCPUIdle,
		HostCPUTotal: hostCollector.collectCPUTotal,
		HostMemTotal: hostCollector.collectMemoryTotal,
		HostMemFree:  hostCollector.collectMemoryFree,
	}

	c := &Collectors{
		collectors: map[string]Collector{
			Cndev:        cndevCollector,
			PodResources: podResourcesCollector,
			Cnpapi:       cnpapiCollector,
			Host:         hostCollector,
		},
		metrics: m,
	}

	for _, collector := range c.collectors {
		collector.init(mluStat)
	}

	expected := map[*prometheus.Desc][]metricValue{}

	for mlu, stat := range mluStat {
		l := []string{stat.sn, stat.model, strings.ToLower(stat.model), node, fmt.Sprintf("%d", stat.slot), stat.mcu, stat.driver, stat.uuid}
		expected[m[Cndev][Temperature].desc] = append(expected[m[Cndev][Temperature].desc], append(l, "", fmt.Sprintf("%v", mluMetrics[mlu].temperature)))
		for cluster, temp := range mluMetrics[mlu].clusterTemperatures {
			expected[m[Cndev][Temperature].desc] = append(expected[m[Cndev][Temperature].desc], append(l, fmt.Sprintf("%d", cluster), fmt.Sprintf("%v", temp)))
		}
		expected[m[Cndev][Health].desc] = append(expected[m[Cndev][Health].desc], append(l, fmt.Sprintf("%v", mluMetrics[mlu].health)))
		expected[m[Cndev][MemTotal].desc] = append(expected[m[Cndev][MemTotal].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].memTotal*1024*1024))))
		expected[m[Cndev][MemUsed].desc] = append(expected[m[Cndev][MemUsed].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].memUsed*1024*1024))))
		expected[m[Cndev][MemUtil].desc] = append(expected[m[Cndev][MemUtil].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].memUsed)/float64(mluMetrics[mlu].memTotal)*100)))
		expected[m[Cndev][Util].desc] = append(expected[m[Cndev][Util].desc], append(l, fmt.Sprintf("%v", mluMetrics[mlu].boardUtil)))
		expected[m[Cndev][Capacity].desc] = append(expected[m[Cndev][Capacity].desc], append(l, fmt.Sprintf("%v", getModelCapacity(stat.model))))
		expected[m[Cndev][Usage].desc] = append(expected[m[Cndev][Usage].desc], append(l, fmt.Sprintf("%v", getModelCapacity(stat.model)*float64(mluMetrics[mlu].boardUtil)/100)))
		for core, u := range mluMetrics[mlu].coreUtil {
			expected[m[Cndev][CoreUtil].desc] = append(expected[m[Cndev][CoreUtil].desc], append(l, fmt.Sprintf("%d", core), fmt.Sprintf("%v", u)))
		}
		expected[m[Cndev][FanSpeed].desc] = append(expected[m[Cndev][FanSpeed].desc], append(l, fmt.Sprintf("%v", -1)))
		expected[m[Cndev][PowerUsage].desc] = append(expected[m[Cndev][PowerUsage].desc], append(l, fmt.Sprintf("%v", mluMetrics[mlu].power)))
		expected[m[Cndev][Version].desc] = append(expected[m[Cndev][Version].desc], append(l, "1"))
		for vf := 0; vf < len(mluMetrics[stat.uuid].virtualFunctionMemUsed); vf++ {
			memUtil := float64(mluMetrics[mlu].virtualFunctionMemUsed[vf]) / float64(mluMetrics[mlu].virtualFunctionMemTotal) * 100
			util, _ := calcVFUtil(mluMetrics[mlu].coreUtil, len(mluMetrics[stat.uuid].virtualFunctionMemUsed), vf+1)
			expected[m[Cndev][VirtualFunctionMemUtil].desc] = append(expected[m[Cndev][VirtualFunctionMemUtil].desc], append(l, fmt.Sprintf("%d", vf+1), fmt.Sprintf("%v", memUtil)))
			expected[m[Cndev][VirtualFunctionUtil].desc] = append(expected[m[Cndev][VirtualFunctionUtil].desc], append(l, fmt.Sprintf("%d", vf+1), fmt.Sprintf("%v", util)))
		}
		expected[m[Cndev][PCIeInfo].desc] = append(expected[m[Cndev][PCIeInfo].desc], append(mluMetrics[mlu].pcieInfo, append(l, "1")...))
		expected[m[Cndev][ChipCPUUtil].desc] = append(expected[m[Cndev][ChipCPUUtil].desc], append(l, "", fmt.Sprintf("%v", mluMetrics[mlu].chipCPUUtil)))
		for core, u := range mluMetrics[mlu].chipCPUCoreUtil {
			expected[m[Cndev][ChipCPUUtil].desc] = append(expected[m[Cndev][ChipCPUUtil].desc], append(l, fmt.Sprintf("%d", core), fmt.Sprintf("%v", u)))
		}
		expected[m[Cndev][VirtualMemTotal].desc] = append(expected[m[Cndev][VirtualMemTotal].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].virtualMemTotal*1024*1024))))
		expected[m[Cndev][VirtualMemUsed].desc] = append(expected[m[Cndev][VirtualMemUsed].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].virtualMemUsed*1024*1024))))
		expected[m[Cndev][ArmOsMemTotal].desc] = append(expected[m[Cndev][ArmOsMemTotal].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].armOsMemTotal*1024))))
		expected[m[Cndev][ArmOsMemUsed].desc] = append(expected[m[Cndev][ArmOsMemUsed].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].armOsMemUsed*1024))))
		for core, u := range mluMetrics[mlu].videoCodecUtil {
			expected[m[Cndev][VideoCodecUtil].desc] = append(expected[m[Cndev][VideoCodecUtil].desc], append(l, fmt.Sprintf("%d", core), fmt.Sprintf("%v", u)))
		}
		for core, u := range mluMetrics[mlu].imageCodecUtil {
			expected[m[Cndev][ImageCodecUtil].desc] = append(expected[m[Cndev][ImageCodecUtil].desc], append(l, fmt.Sprintf("%d", core), fmt.Sprintf("%v", u)))
		}
		for core, u := range mluMetrics[mlu].tinyCoreUtil {
			expected[m[Cndev][TinyCoreUtil].desc] = append(expected[m[Cndev][TinyCoreUtil].desc], append(l, fmt.Sprintf("%d", core), fmt.Sprintf("%v", u)))
		}
		expected[m[Cndev][NUMANodeID].desc] = append(expected[m[Cndev][NUMANodeID].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].numaNodeID))))
		expected[m[Cndev][DDRDataWidth].desc] = append(expected[m[Cndev][DDRDataWidth].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].ddrDataWidth))))
		expected[m[Cndev][DDRBandWidth].desc] = append(expected[m[Cndev][DDRBandWidth].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].ddrBandWidth))))
		expected[m[Cndev][ECCAddressForbiddenError].desc] = append(expected[m[Cndev][ECCAddressForbiddenError].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].eccAddressForbiddenError))))
		expected[m[Cndev][ECCCorrectedError].desc] = append(expected[m[Cndev][ECCCorrectedError].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].eccCorrectedError))))
		expected[m[Cndev][ECCMultipleError].desc] = append(expected[m[Cndev][ECCMultipleError].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].eccMultipleError))))
		expected[m[Cndev][ECCMultipleMultipleError].desc] = append(expected[m[Cndev][ECCMultipleMultipleError].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].eccMultipleMultipleError))))
		expected[m[Cndev][ECCMultipleOneError].desc] = append(expected[m[Cndev][ECCMultipleOneError].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].eccMultipleOneError))))
		expected[m[Cndev][ECCOneBitError].desc] = append(expected[m[Cndev][ECCOneBitError].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].eccOneBitError))))
		expected[m[Cndev][ECCTotalError].desc] = append(expected[m[Cndev][ECCTotalError].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].eccTotalError))))
		expected[m[Cndev][ECCUncorrectedError].desc] = append(expected[m[Cndev][ECCUncorrectedError].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].eccUncorrectedError))))
		for link := 0; link < mluMetrics[stat.uuid].mluLinkPortNumber; link++ {
			expected[m[Cndev][MLULinkCapabilityP2PTransfer].desc] = append(expected[m[Cndev][MLULinkCapabilityP2PTransfer].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCapabilityP2PTransfer[link]))))
			expected[m[Cndev][MLULinkCapabilityInterlakenSerdes].desc] = append(expected[m[Cndev][MLULinkCapabilityInterlakenSerdes].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCapabilityInterlakenSerdes[link]))))
			expected[m[Cndev][MLULinkCounterCntrReadByte].desc] = append(expected[m[Cndev][MLULinkCounterCntrReadByte].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCounterCntrReadByte[link]))))
			expected[m[Cndev][MLULinkCounterCntrReadPackage].desc] = append(expected[m[Cndev][MLULinkCounterCntrReadPackage].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCounterCntrReadPackage[link]))))
			expected[m[Cndev][MLULinkCounterCntrWriteByte].desc] = append(expected[m[Cndev][MLULinkCounterCntrWriteByte].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCounterCntrWriteByte[link]))))
			expected[m[Cndev][MLULinkCounterCntrWritePackage].desc] = append(expected[m[Cndev][MLULinkCounterCntrWritePackage].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCounterCntrWritePackage[link]))))
			expected[m[Cndev][MLULinkCounterErrCorrected].desc] = append(expected[m[Cndev][MLULinkCounterErrCorrected].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCounterErrCorrected[link]))))
			expected[m[Cndev][MLULinkCounterErrCRC24].desc] = append(expected[m[Cndev][MLULinkCounterErrCRC24].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCounterErrCRC24[link]))))
			expected[m[Cndev][MLULinkCounterErrCRC32].desc] = append(expected[m[Cndev][MLULinkCounterErrCRC32].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCounterErrCRC32[link]))))
			expected[m[Cndev][MLULinkCounterErrEccDouble].desc] = append(expected[m[Cndev][MLULinkCounterErrEccDouble].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCounterErrEccDouble[link]))))
			expected[m[Cndev][MLULinkCounterErrFatal].desc] = append(expected[m[Cndev][MLULinkCounterErrFatal].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCounterErrFatal[link]))))
			expected[m[Cndev][MLULinkCounterErrReplay].desc] = append(expected[m[Cndev][MLULinkCounterErrReplay].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCounterErrReplay[link]))))
			expected[m[Cndev][MLULinkCounterErrUncorrected].desc] = append(expected[m[Cndev][MLULinkCounterErrUncorrected].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkCounterErrUncorrected[link]))))
			expected[m[Cndev][MLULinkPortMode].desc] = append(expected[m[Cndev][MLULinkPortMode].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkPortMode[link]))))
			expected[m[Cndev][MLULinkSpeedFormat].desc] = append(expected[m[Cndev][MLULinkSpeedFormat].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkSpeedFormat[link]))))
			expected[m[Cndev][MLULinkSpeedValue].desc] = append(expected[m[Cndev][MLULinkSpeedValue].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkSpeedValue[link]))))
			expected[m[Cndev][MLULinkStatusIsActive].desc] = append(expected[m[Cndev][MLULinkStatusIsActive].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkStatusIsActive[link]))))
			expected[m[Cndev][MLULinkStatusSerdesState].desc] = append(expected[m[Cndev][MLULinkStatusSerdesState].desc], append(l, fmt.Sprintf("%d", link), fmt.Sprintf("%v", float64(mluMetrics[mlu].mluLinkStatusSerdesState[link]))))
			expected[m[Cndev][MLULinkVersion].desc] = append(expected[m[Cndev][MLULinkVersion].desc], append(l, fmt.Sprintf("%d", link), "1", calcVersion(mluMetrics[mlu].mluLinkMajor[link], mluMetrics[mlu].mluLinkMinor[link], mluMetrics[mlu].mluLinkBuild[link])))
		}
		expected[m[Cndev][D2DCRCError].desc] = append(expected[m[Cndev][D2DCRCError].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].d2dCRCError))))
		expected[m[Cndev][D2DCRCErrorOverflow].desc] = append(expected[m[Cndev][D2DCRCErrorOverflow].desc], append(l, fmt.Sprintf("%v", float64(mluMetrics[mlu].d2dCRCErrorOverflow))))
		for i, pid := range mluMetrics[mlu].pid {
			expected[m[Cndev][ProcessIpuUtil].desc] = append(expected[m[Cndev][ProcessIpuUtil].desc], append(l, fmt.Sprintf("%d", pid), fmt.Sprintf("%v", float64(mluMetrics[mlu].processIpuUtil[i]))))
			expected[m[Cndev][ProcessJpuUtil].desc] = append(expected[m[Cndev][ProcessJpuUtil].desc], append(l, fmt.Sprintf("%d", pid), fmt.Sprintf("%v", float64(mluMetrics[mlu].processJpuUtil[i]))))
			expected[m[Cndev][ProcessMemoryUtil].desc] = append(expected[m[Cndev][ProcessMemoryUtil].desc], append(l, fmt.Sprintf("%d", pid), fmt.Sprintf("%v", float64(mluMetrics[mlu].processMemUtil[i]))))
			expected[m[Cndev][ProcessVpuDecodeUtil].desc] = append(expected[m[Cndev][ProcessVpuDecodeUtil].desc], append(l, fmt.Sprintf("%d", pid), fmt.Sprintf("%v", float64(mluMetrics[mlu].processVpuDecodeUtil[i]))))
			expected[m[Cndev][ProcessVpuEncodeUtil].desc] = append(expected[m[Cndev][ProcessVpuEncodeUtil].desc], append(l, fmt.Sprintf("%d", pid), fmt.Sprintf("%v", float64(mluMetrics[mlu].processVpuEncodeUtil[i]))))
		}

		expected[m[Cnpapi][PCIeRead].desc] = append(expected[m[Cnpapi][PCIeRead].desc], append(l, fmt.Sprintf("%v", mluMetrics[mlu].pcieRead/1024/1024/1024)))
		expected[m[Cnpapi][PCIeWrite].desc] = append(expected[m[Cnpapi][PCIeWrite].desc], append(l, fmt.Sprintf("%v", mluMetrics[mlu].pcieWrite/1024/1024/1024)))
		expected[m[Cnpapi][DramRead].desc] = append(expected[m[Cnpapi][DramRead].desc], append(l, fmt.Sprintf("%v", mluMetrics[mlu].dramRead/1024/1024/1024)))
		expected[m[Cnpapi][DramWrite].desc] = append(expected[m[Cnpapi][DramWrite].desc], append(l, fmt.Sprintf("%v", mluMetrics[mlu].dramWrite/1024/1024/1024)))
		expected[m[Cnpapi][MLULinkRead].desc] = append(expected[m[Cnpapi][MLULinkRead].desc], append(l, fmt.Sprintf("%v", mluMetrics[mlu].mlulinkRead/1024/1024/1024)))
		expected[m[Cnpapi][MLULinkWrite].desc] = append(expected[m[Cnpapi][MLULinkWrite].desc], append(l, fmt.Sprintf("%v", mluMetrics[mlu].mlulinkWrite/1024/1024/1024)))
	}

	var allocatedDriver, allocatedModel, allocatedType string
	for id, info := range devicePodInfo {
		uuid := strings.Split(id, uuidPrefix)[1]
		vf := ""
		if strings.Contains(uuid, envShareSubString) {
			uuid = strings.Split(uuid, envShareSubString)[0]
		}
		if strings.Contains(uuid, sriovSubString) {
			s := strings.Split(uuid, sriovSubString)
			uuid = s[0]
			vf = s[1]
		}
		stat := mluStat[uuid]
		allocatedDriver = stat.driver
		allocatedModel = stat.model
		allocatedType = strings.ToLower(stat.model)
		l := []string{stat.sn, stat.model, strings.ToLower(stat.model), node, fmt.Sprintf("%d", stat.slot), stat.mcu, stat.driver, stat.uuid}
		expected[m[PodResources][Container].desc] = append(expected[m[PodResources][Container].desc], append(l, "1", info.Container, info.Namespace, info.Pod, vf))
	}
	expected[m[PodResources][Allocated].desc] = append(expected[m[PodResources][Allocated].desc], []string{node, allocatedDriver, allocatedModel, allocatedType, fmt.Sprintf("%d", len(devicePodInfo))})
	stat := mluStat["uuid4"]
	expected[m[PodResources][Container].desc] = append(expected[m[PodResources][Container].desc], []string{stat.sn, stat.model, strings.ToLower(stat.model), node, fmt.Sprintf("%d", stat.slot), stat.mcu, stat.driver, stat.uuid, "0", "", "", "", ""})

	expected[m[Host][HostMemTotal].desc] = append(expected[m[Host][HostMemTotal].desc], []string{fmt.Sprintf("%v", hostMemTotal/1024/1024), node})
	expected[m[Host][HostMemFree].desc] = append(expected[m[Host][HostMemFree].desc], []string{fmt.Sprintf("%v", hostMemFree/1024/1024), node})
	expected[m[Host][HostCPUIdle].desc] = append(expected[m[Host][HostCPUIdle].desc], []string{fmt.Sprintf("%v", hostCPUIdle), node})
	expected[m[Host][HostCPUTotal].desc] = append(expected[m[Host][HostCPUTotal].desc], []string{fmt.Sprintf("%v", hostCPUTotal), node})
	for desc, values := range expected {
		sortMetricValues(values)
		expected[desc] = values
	}

	ch := make(chan prometheus.Metric)
	go func() {
		c.Collect(ch)
		close(ch)
	}()
	got := map[*prometheus.Desc][]metricValue{}
	for m := range ch {
		metric := dto.Metric{}
		m.Write(&metric)
		var values []string
		if metric.GetGauge() != nil {
			values = []string{fmt.Sprintf("%v", metric.GetGauge().GetValue())}
		} else {
			values = []string{fmt.Sprintf("%v", metric.GetCounter().GetValue())}
		}
		for _, l := range metric.GetLabel() {
			values = append(values, l.GetValue())
		}
		got[m.Desc()] = append(got[m.Desc()], values)
	}
	assert.Equal(t, len(expected), len(got))
	for desc, values := range got {
		sortMetricValues(values)
		assert.Equal(t, expected[desc], values)
	}
}

type metricValue []string

func sortMetricValues(values []metricValue) {
	for i, v := range values {
		sort.Strings(v)
		values[i] = v
	}
	sort.Slice(values, func(i, j int) bool {
		for index := range values[i] {
			if values[i][index] == values[j][index] {
				continue
			}
			return values[i][index] < values[j][index]
		}
		return false
	})
}
