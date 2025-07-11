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
	"encoding/json"
	"flag"
	"os"
	"sort"
	"testing"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/Cambricon/mlu-exporter/pkg/mock"
	"github.com/Cambricon/mlu-exporter/pkg/podresources"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var golden = flag.Bool("golden", false, "")

// verify
type testMetrics struct {
	Desc   string
	Metric dto.Metric
}

func TestCollect(t *testing.T) {
	var (
		// fake chassis info
		chassisSn          = uint64(8039172446100346)
		chassisProductDate = "2024-11-23"
		chassisProductName = "HBBU-M9DE"
		chassisVendorName  = "Cambricon"
		chassisPartNumber  = "DS591026"
		chassisBmcIP       = "192.168.100.50"
		chassisNvme        = []cndev.ChassisDevInfo{
			{
				SN:    "S63UNE0W300023",
				Model: "PM9A3",
				Fw:    "00-1.5-1.2",
				Mfc:   "Samsung",
			},
			{
				SN:    "S63UNE0W300024",
				Model: "PM9A3",
				Fw:    "00-1.5-1.2",
				Mfc:   "Samsung",
			},
		}
		chassisIb = []cndev.ChassisDevInfo{
			{
				SN:    "MT22394008L0",
				Model: "MCX653105A-HDAT",
				Fw:    "AG",
				Mfc:   "Mellanox",
			},
			{
				SN:    "MT22394008L1",
				Model: "MCX653105A-HDAT",
				Fw:    "AG",
				Mfc:   "Mellanox",
			},
		}
		chassisPsu = []cndev.ChassisDevInfo{
			{
				SN:    "S63UNE0W300022",
				Model: "ECD16010099",
				Fw:    "00-1.5-1.2-1.3",
				Mfc:   "DELTA",
			},
			{
				SN:    "S63UNE0W300023",
				Model: "ECD16010099",
				Fw:    "00-1.5-1.2-1.3",
				Mfc:   "DELTA",
			},
		}

		// fake node info
		node         = "machine1"
		hostCPUTotal = float64(6185912)
		hostCPUIdle  = float64(34459)
		hostMemTotal = float64(24421820)
		hostMemFree  = float64(18470732)

		// fake card info
		uuid1 = "uuid1"
		uuid2 = "uuid2"
		uuid3 = "uuid3"
		uuid4 = "uuid4"
		mst   = map[string]MLUStat{
			uuid1: {
				slot:       0,
				model:      "MLU590",
				uuid:       uuid1,
				sn:         "sn1",
				mcu:        "v1.1.1",
				driver:     "v2.2.2",
				mimEnabled: true,
				mimInfos: []cndev.MimInfo{
					{
						InstanceInfo: cndev.MimInstanceInfo{
							InstanceName: "instance-1",
							InstanceID:   1,
							UUID:         "test-uuid1",
						},
						ProfileInfo: cndev.MimProfileInfo{
							GDMACount:    1,
							JPUCount:     1,
							MemorySize:   20,
							MLUCoreCount: 3,
							ProfileID:    1,
							ProfileName:  "2m.16g",
							VPUCount:     2,
						},
						PlacementStart: 0,
						PlacementSize:  2,
					},
					{
						InstanceInfo: cndev.MimInstanceInfo{
							InstanceName: "instance-2",
							InstanceID:   2,
							UUID:         "test-uuid1",
						},
						ProfileInfo: cndev.MimProfileInfo{
							GDMACount:    1,
							JPUCount:     1,
							MemorySize:   20,
							MLUCoreCount: 3,
							ProfileID:    1,
							ProfileName:  "2m.16g",
							VPUCount:     2,
						},
						PlacementStart: 0,
						PlacementSize:  2,
					},
				},
				link: 2,
			},
			uuid2: {
				slot:        1,
				model:       "MLU590",
				uuid:        uuid2,
				sn:          "sn2",
				mcu:         "v1.1.1",
				driver:      "v2.2.2",
				smluEnabled: true,
				smluInfos: []cndev.SmluInfo{
					{
						InstanceInfo: cndev.SmluInstanceInfo{
							InstanceID:   1,
							InstanceName: "instance-1",
							UUID:         "test-uuid-1",
						},
						ProfileInfo: cndev.SmluProfileInfo{
							IpuTotal:    20,
							MemTotal:    100000000,
							ProfileID:   1,
							ProfileName: "12.000m.78.784gb",
						},
						IpuUtil: 10,
						MemUsed: 3400000,
					},
					{
						InstanceInfo: cndev.SmluInstanceInfo{
							InstanceID:   2,
							InstanceName: "instance-2",
							UUID:         "test-uuid-2",
						},
						ProfileInfo: cndev.SmluProfileInfo{
							IpuTotal:    20,
							MemTotal:    100000000,
							ProfileID:   1,
							ProfileName: "12.000m.78.784gb",
						},
						IpuUtil: 10,
						MemUsed: 3400000,
					},
				},
				link: 2,
			},
			uuid3: {
				slot:   2,
				model:  "MLU590",
				uuid:   uuid3,
				sn:     "sn3",
				mcu:    "v1.1.1",
				driver: "v2.2.2",
				link:   2,
			},
			uuid4: {
				slot:                   3,
				model:                  "MLU590",
				uuid:                   uuid4,
				sn:                     "sn4",
				mcu:                    "v1.1.1",
				driver:                 "v2.2.2",
				link:                   2,
				cndevInterfaceDisabled: map[string]bool{"crcDisabled": true},
			},
		}

		mimProfileInfos = []cndev.MimProfileInfo{
			{
				GDMACount:    1,
				JPUCount:     1,
				MemorySize:   20,
				MLUCoreCount: 3,
				ProfileID:    1,
				ProfileName:  "2m.16g",
				VPUCount:     2,
			},
		}
		mimProfileMaxInstanceCount = []int{1, 0, 0, 0}
		sMluProfileIDInfo          = [][]int{
			{1},
			{1},
			{1},
			{1},
		}
		sMluProfileInfo = cndev.SmluProfileInfo{
			IpuTotal:    20,
			MemTotal:    100000000,
			ProfileID:   1,
			ProfileName: "12.000m.78.784gb",
		}
		sMluProfileTotal  = []uint32{0, 2, 0, 0}
		sMluProfileRemain = []uint32{0, 1, 0, 0}
		// fake container info
		devicePodInfo = map[string]podresources.PodInfo{
			"MLU-uuid1": {
				Pod:       "pod1",
				Namespace: "namespace1",
				Container: "container1",
			},
			"MLU-uuid2-mim-MLU-test-uuid-1": {
				Pod:       "pod2",
				Namespace: "namespace2",
				Container: "container2",
			},
		}

		// fake metrics
		cndevVersion      = []uint{3, 12, 1}
		computeCapability = []uint{3, 12}
		computeMode       = []uint{0, 2, 2, 0}

		// util
		boardUtil = []int{11, 12, 13, 21}
		coreUtil  = [][]int{
			{11, 11, 11, 11, 13, 13, 13, 13, 15, 15, 15, 15, 17, 17, 17, 17},
			{12, 12, 12, 12, 14, 14, 14, 14, 16, 16, 16, 16, 18, 18, 18, 18},
			{21, 21, 21, 21, 13, 13, 13, 13, 17, 17, 17, 17, 19, 19, 19, 19},
			{31, 31, 31, 31, 33, 33, 33, 13, 35, 35, 35, 35, 17, 17, 17, 17},
		}

		// memory
		memUsed         = []int64{20, 21, 22, 23}
		memTotal        = int64(1000)
		virtualMemUsed  = []int64{100, 101, 102, 103}
		virtualMemTotal = int64(1000)

		// vf
		virtualFunctionMemUsed = [][]int64{
			{11, 12},
			{},
			{11, 13, 15, 17},
			{12, 16},
		}
		virtualFunctionMemTotal   = int64(1000)
		virtualFunctionPowerUsage = []int{30, 30, 30, 30}
		vfState                   = []int{3, 0, 15, 3}

		// temperature
		temperature         = []int{23, 23, 22, -100}
		memTemperature      = []int{22, 24, 10, 11}
		chipTemperature     = []int{22, 24, 10, 11}
		clusterTemperatures = [][]int{
			{22, 23, 24, 26},
			{21, 23, 24, 36},
			{22, 43, 24, 26},
			{-100, -100, -100, -100},
		}
		memDieTemperatures = [][]int{
			{21, 23},
			{22, 26},
			{10, 10},
			{11, 11},
		}
		overTemperatureInfo = [][]uint32{
			{2, 23},
			{2, 26},
			{0, 10},
			{1, 11},
		}
		overTemperatureThreshold = [][]int{
			{120, 110},
			{120, 110},
			{110, 100},
			{120, 100},
		}

		throttleReason = []bool{true, false, false, true}

		// health
		health   = []int{1, 0, 1, 0}
		unhealth = [][]string{
			{},
			{"CNDEV_FR_TEMP_VIOLATION"},
			{},
			{"CNDEV_FR_REPAIR_FAILURE"},
		}

		// parity
		parityError = []int{3, 0, 0, 0}

		// power
		power             = []int{30, 31, 32, 33}
		powerDefaultLimit = []uint16{300, 310, 320, 330}
		powerCurrentLimit = []uint16{300, 310, 320, 330}
		powerLimitRange   = [][]uint16{
			{350, 470},
			{330, 450},
			{350, 470},
			{360, 480},
		}
		current = [][]int{
			{31000, 42000, 13000},
			{32000, 43000, 14000},
			{35000, 45000, 16000},
			{37000, 47000, 15000},
		}
		voltage = [][]int{
			{310, 420, 130},
			{320, 430, 140},
			{350, 450, 160},
			{370, 470, 150},
		}

		// pcie
		pcieSlotID          = []int{0, 1, 2, 3}
		pcieSubsystemID     = uint(0x12)
		pcieDeviceID        = uint(0x590)
		pcieVendor          = uint16(0xcabc)
		pcieSubsystemVendor = uint16(0xcabc)
		pcieDomain          = uint(0)
		pcieBus             = []uint{26, 27, 28, 29}
		pcieDevice          = uint(0)
		pcieFunction        = uint(0)
		pcieSpeed           = 4
		pcieWidth           = 4
		moduleID            = []uint16{0, 1, 2, 3}

		pcieRead       = []int64{100, 200, 300, 400}
		pcieReplay     = []uint32{100, 200, 300, 400}
		pcieWrite      = []int64{100, 200, 300, 400}
		bar4MemoryInfo = [][]uint64{
			{310, 420, 130},
			{320, 430, 140},
			{350, 450, 160},
			{370, 470, 150},
		}
		maxPCIeInfo = [][]int{
			{3, 16},
			{4, 16},
			{5, 8},
			{5, 16},
		}

		// process
		pid = [][]uint32{
			{1100, 1101},
			{1102, 1103},
			{1104, 1105},
			{1106, 1107},
		}
		processIpuUtil = [][]uint32{
			{11, 12},
			{13, 32},
			{15, 21},
			{41, 22},
		}
		processJpuUtil         = processIpuUtil
		processMemUtil         = processIpuUtil
		processVpuDecodeUtil   = processIpuUtil
		processVpuEncodeUtil   = processIpuUtil
		processPhysicalMemUsed = [][]uint64{
			{1100, 1200},
			{1300, 3200},
			{1500, 2100},
			{4100, 2200},
		}
		processVirtualMemUsed = processPhysicalMemUsed

		// chip
		chipCPUUtil     = []uint16{25, 25, 26, 26}
		chipCPUCoreUtil = [][]uint8{
			{0, 0, 100, 0},
			{0, 100, 0, 0},
			{0, 50, 50, 0},
			{25, 25, 25, 25},
		}
		coreCPUUser = chipCPUCoreUtil
		coreCPUNice = chipCPUCoreUtil
		coreCPUSys  = chipCPUCoreUtil
		coreCPUSi   = chipCPUCoreUtil
		coreCPUHi   = chipCPUCoreUtil

		deviceOsMemUsed  = []int64{1024, 4096, 512, 2048}
		deviceOsMemTotal = int64(10240)
		tinyCoreUtil     = [][]int{
			{10, 21},
			{12, 21},
			{15, 11},
			{21, 12},
		}
		numaNodeID   = []int{0, 0, 1, 1}
		ddrBandWidth = []float64{3.1, 3.2, 3.3, 3.4}
		ddrDataWidth = []int{1, 2, 3, 4}

		// frequency
		ddrFrequency     = []int{100, 50, 150, 200}
		defaultFrequency = []int{2000, 2000, 2000, 2000}
		frequency        = []int{1000, 500, 1500, 2000}
		frequencyStatus  = []int{0, 0, 1, 1}

		// codec
		videoCodecUtil = [][]int{
			{0, 2, 0, 1, 0, 1},
			{0, 1, 0, 3, 0, 1},
			{0, 1, 0, 1, 0, 2},
			{0, 4, 0, 1, 0, 1},
		}
		imageCodecUtil   = videoCodecUtil
		videoDecoderUtil = videoCodecUtil

		// retired page
		pageCounts = []uint32{100, 10, 0, 0}
		retirement = []int{1, 0, 1, 1}

		// remapped rows
		correctRows   = []uint32{2, 1, 1, 2}
		failedRows    = []uint32{1, 1, 1, 1}
		pendingRows   = []uint32{2, 1, 2, 1}
		uncorrectRows = []uint32{1, 1, 1, 1}
		retirePending = []uint32{1, 1, 0, 1}

		// repair status
		repairPending       = []bool{false, false, true, false}
		repairFailed        = repairPending
		repairRetirePending = repairPending

		// ecc & crc
		eccAddressForbiddenError = []uint64{100, 100, 100, 100}
		eccCorrectedError        = eccAddressForbiddenError
		eccMultipleError         = eccAddressForbiddenError
		eccMultipleMultipleError = eccAddressForbiddenError
		eccMultipleOneError      = eccAddressForbiddenError
		eccOneBitError           = eccAddressForbiddenError
		eccTotalError            = eccAddressForbiddenError
		eccUncorrectedError      = eccAddressForbiddenError
		d2dCRCError              = []uint64{128, 128, 128, 128}
		d2dCRCErrorOverflow      = d2dCRCError
		eccVolatile              = [][]uint32{
			{0, 2, 0, 1, 0},
			{0, 1, 0, 3, 0},
			{0, 1, 0, 1, 0},
			{0, 4, 0, 1, 0},
		}
		eccAggregate = [][]uint64{
			{0, 2, 0, 1, 0, 1},
			{0, 1, 0, 3, 0, 1},
			{0, 1, 0, 1, 0, 0},
			{0, 4, 0, 1, 0, 0},
		}
		addressSwap = [][]uint32{
			{0, 2},
			{0, 1},
			{0, 1},
			{0, 4},
		}
		eccMode = [][]int{
			{0, 1},
			{0, 1},
			{0, 1},
			{1, 0},
		}
		// heartbeat count
		heartbeatCount = []uint32{10, 10, 10, 10}

		// device count
		deviceCount = []uint{4, 4, 4, 4}

		// mlulink
		mluLinkCapabilityP2PTransfer = [][]uint{
			{1, 0},
			{1, 1},
			{0, 1},
			{0, 0},
		}
		mluLinkCounterCntrReadByte = [][]uint64{
			{1000, 10000},
			{20000, 40000},
			{30000, 5000},
			{50000, 30000},
		}
		mluLinkEventCounter = [][]uint64{
			{0, 1},
			{1, 0},
			{1, 3},
			{2, 1},
		}
		mluLinkPortMode = [][]int{
			{1, 1},
			{1, 1},
			{1, 1},
			{1, 1},
		}
		mluLinkPortPPI = [][]string{
			{"1:0", "1:1"},
			{"1:0", "1:1"},
			{"N/A", "1:1"},
			{"1:0", "N/A"},
		}
		mluLinkSpeedValue = [][]float32{
			{12, 13},
			{22, 23},
			{32, 31},
			{42, 43},
		}
		mluLinkRemoteMcSn = [][]uint64{
			{000001, 000002},
			{000011, 000012},
			{000021, 000022},
			{000031, 000032},
		}
		mluLinkRemoteSlotID = [][]uint32{
			{0, 1},
			{1, 0},
			{1, 3},
			{2, 1},
		}
		mluLinkRemoteConnectType = [][]int32{
			{0, 0},
			{0, 1},
			{1, 0},
			{0, 0},
		}
		mluLinkRemoteIP = [][]string{
			{"1.1.1.1", "2.2.2.2"},
			{"1.1.1.2", "2.2.2.3"},
			{"1.1.1.3", "2.2.2.4"},
			{"1.1.1.4", "2.2.2.5"},
		}
		mluLinkState = [][]int{
			{1, 0},
			{1, 1},
			{0, 1},
			{1, 1},
		}
		mluLinkOpticalPresent = [][]uint8{
			{1, 0},
			{1, 1},
			{0, 2},
			{0, 0},
		}
		mluLinkOpticalTxpwr = [][][]float32{
			{{0.5, 0.7}, {1.5, 1.7}},
			{{0.5, 1.7}, {1.5, 0.7}},
			{{0.5, 2.7}, {1.5, 0.7}},
			{{0.5, 3.7}, {1.5, 0.7}},
		}
		mluLinkOpticalTemp                = mluLinkSpeedValue
		mluLinkOpticalVolt                = mluLinkSpeedValue
		mluLinkOpticalRxpwr               = mluLinkOpticalTxpwr
		mluLinkRemoteBaSn                 = mluLinkRemoteMcSn
		mluLinkRemotePortID               = mluLinkRemoteSlotID
		mluLinkRemoteNcsUUID64            = mluLinkRemoteMcSn
		mluLinkRemoteMac                  = mluLinkRemoteIP
		mluLinkRemoteUUID                 = mluLinkRemoteIP
		mluLinkRemotePortName             = mluLinkRemoteIP
		mluLinkPortNumber                 = 2
		mluLinkCapabilityInterlakenSerdes = mluLinkCapabilityP2PTransfer
		mluLinkCounterCntrCnpPackage      = mluLinkCounterCntrReadByte
		mluLinkCounterCntrPfcPackage      = mluLinkCounterCntrReadByte
		mluLinkCounterCntrReadPackage     = mluLinkCounterCntrReadByte
		mluLinkCounterCntrWriteByte       = mluLinkCounterCntrReadByte
		mluLinkCounterCntrWritePackage    = mluLinkCounterCntrReadByte
		mluLinkCounterErrCorrected        = mluLinkCounterCntrReadByte
		mluLinkCounterErrCRC24            = mluLinkCounterCntrReadByte
		mluLinkCounterErrCRC32            = mluLinkCounterCntrReadByte
		mluLinkCounterErrEccDouble        = mluLinkCounterCntrReadByte
		mluLinkCounterErrFatal            = mluLinkCounterCntrReadByte
		mluLinkCounterErrReplay           = mluLinkCounterCntrReadByte
		mluLinkCounterErrUncorrected      = mluLinkCounterCntrReadByte
		mluLinkErrorCounter               = mluLinkEventCounter
		mluLinkCorrectFecCounter          = mluLinkEventCounter
		mluLinkUncorrectFecCounter        = mluLinkEventCounter
		mluLinkSpeedFormat                = mluLinkPortMode
		mluLinkStatusIsActive             = mluLinkPortMode
		mluLinkStatusSerdesState          = mluLinkPortMode
		mluLinkStatusCableState           = mluLinkPortMode
		mluLinkMajor                      = mluLinkCapabilityP2PTransfer
		mluLinkMinor                      = mluLinkCapabilityP2PTransfer
		mluLinkBuild                      = mluLinkCapabilityP2PTransfer
		rdmaDevice                        = []rdmaDevice{
			{
				name:        "mlx5_1",
				pcieAddress: "0000:69:00.0",
				domain:      0,
				bus:         105,
				device:      0,
				function:    0,
				nicName:     "roce1",
				ipAddress:   "192.168.1.1",
			},
			{
				name:        "mlx5_2",
				pcieAddress: "0000:70:00.0",
				domain:      0,
				bus:         106,
				device:      0,
				function:    0,
				nicName:     "roce2",
				ipAddress:   "192.168.2.1",
			},
		}

		// tensor util
		tensorUtil = []int{80, 81, 82, 102}
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mcndev mock response
	mcndev := mock.NewCndev(ctrl)
	mcndev.EXPECT().Init(false).Return(nil).AnyTimes()
	mcndev.EXPECT().Init(true).Return(nil).AnyTimes()
	stub := gomonkey.ApplyFunc(fetchMLUCounts, func() (uint, error) {
		return 4, nil
	})
	defer stub.Reset()

	slots := []int{}
	for _, stat := range mst {
		mcndev.EXPECT().GetDeviceTemperature(stat.slot).Return(temperature[stat.slot], memTemperature[stat.slot], chipTemperature[stat.slot], clusterTemperatures[stat.slot], memDieTemperatures[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceHealth(stat.slot).Return(health[stat.slot], true, true, unhealth[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceVfState(stat.slot).Return(vfState[stat.slot], nil).AnyTimes()
		for vf := 0; vf < len(virtualFunctionMemUsed[stat.slot]); vf++ {
			mcndev.EXPECT().GetDeviceMemory(uint((vf+1)<<8|int(stat.slot))).Return(virtualFunctionMemUsed[stat.slot][vf], virtualFunctionMemTotal, virtualMemUsed[stat.slot], virtualMemTotal, nil).AnyTimes()
			mcndev.EXPECT().GetDevicePower(uint((vf+1)<<8|int(stat.slot))).Return(virtualFunctionPowerUsage[stat.slot], nil).AnyTimes()
		}
		mcndev.EXPECT().GetDeviceMemory(stat.slot).Return(memUsed[stat.slot], memTotal, virtualMemUsed[stat.slot], virtualMemTotal, nil).AnyTimes()
		mcndev.EXPECT().GetDeviceUtil(stat.slot).Return(boardUtil[stat.slot], coreUtil[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceFanSpeed(stat.slot).Return(0, nil).AnyTimes()
		mcndev.EXPECT().GetAllSMluInfo(stat.slot).Return(stat.smluInfos, nil).AnyTimes()
		mcndev.EXPECT().GetDevicePower(stat.slot).Return(power[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDevicePCIeInfo(stat.slot).Return(pcieSlotID[stat.slot], pcieSubsystemID, pcieDeviceID, pcieVendor, pcieSubsystemVendor, pcieDomain, pcieBus[stat.slot], pcieDevice, pcieFunction, moduleID[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDevicePCIeThroughput(stat.slot).Return(pcieRead[stat.slot], pcieWrite[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceCurrentPCIeInfo(stat.slot).Return(pcieSpeed, pcieWidth, nil).AnyTimes()
		mcndev.EXPECT().GetDevicePCIeReplayCount(stat.slot).Return(pcieReplay[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceCPUUtil(stat.slot).Return(chipCPUUtil[stat.slot], chipCPUCoreUtil[stat.slot], coreCPUUser[stat.slot], coreCPUNice[stat.slot], coreCPUSys[stat.slot], coreCPUSi[stat.slot], coreCPUHi[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceOsMemory(stat.slot).Return(deviceOsMemUsed[stat.slot], deviceOsMemTotal, nil).AnyTimes()
		mcndev.EXPECT().GetDeviceTinyCoreUtil(stat.slot).Return(tinyCoreUtil[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceNUMANodeID(stat.slot).Return(numaNodeID[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceDDRInfo(stat.slot).Return(ddrDataWidth[stat.slot], ddrBandWidth[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceFrequency(stat.slot).Return(frequency[stat.slot], ddrFrequency[stat.slot], defaultFrequency[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceFrequencyStatus(stat.slot).Return(frequencyStatus[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceRetiredPageInfo(stat.slot).Return(pageCounts[stat.slot], pageCounts[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceRetiredPagesOperation(stat.slot).Return(retirement[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceRemappedRows(stat.slot).Return(correctRows[stat.slot], uncorrectRows[stat.slot], pendingRows[stat.slot], failedRows[stat.slot], retirePending[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceRepairStatus(stat.slot).Return(repairPending[stat.slot], repairFailed[stat.slot], repairRetirePending[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceVideoCodecUtil(stat.slot).Return(videoCodecUtil[stat.slot], videoDecoderUtil[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceImageCodecUtil(stat.slot).Return(imageCodecUtil[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceECCInfo(stat.slot).Return(eccAddressForbiddenError[stat.slot], eccCorrectedError[stat.slot], eccMultipleError[stat.slot], eccMultipleMultipleError[stat.slot], eccMultipleOneError[stat.slot], eccOneBitError[stat.slot], eccTotalError[stat.slot], eccUncorrectedError[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceEccMode(stat.slot).Return(eccMode[stat.slot][0], eccMode[stat.slot][1], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceCRCInfo(stat.slot).Return(d2dCRCError[stat.slot], d2dCRCErrorOverflow[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceMemEccCounter(stat.slot).Return(eccVolatile[stat.slot], eccAggregate[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceHeartbeatCount(stat.slot).Return(heartbeatCount[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceCount().Return(deviceCount[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceTensorUtil(stat.slot).Return(tensorUtil[stat.slot], nil).AnyTimes()
		for link := 0; link < mluLinkPortNumber; link++ {
			mcndev.EXPECT().GetDeviceMLULinkCapability(stat.slot, uint(link)).Return(mluLinkCapabilityP2PTransfer[stat.slot][link], mluLinkCapabilityInterlakenSerdes[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkCounter(stat.slot, uint(link)).Return(mluLinkCounterCntrReadByte[stat.slot][link], mluLinkCounterCntrReadPackage[stat.slot][link], mluLinkCounterCntrWriteByte[stat.slot][link], mluLinkCounterCntrWritePackage[stat.slot][link],
				mluLinkCounterErrCorrected[stat.slot][link], mluLinkCounterErrCRC24[stat.slot][link], mluLinkCounterErrCRC32[stat.slot][link], mluLinkCounterErrEccDouble[stat.slot][link], mluLinkCounterErrFatal[stat.slot][link], mluLinkCounterErrReplay[stat.slot][link],
				mluLinkCounterErrUncorrected[stat.slot][link], mluLinkCounterCntrCnpPackage[stat.slot][link], mluLinkCounterCntrPfcPackage[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkErrorCounter(stat.slot, uint(link)).Return(mluLinkErrorCounter[stat.slot][link], mluLinkCorrectFecCounter[stat.slot][link], mluLinkUncorrectFecCounter[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkEventCounter(stat.slot, uint(link)).Return(mluLinkEventCounter[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceOpticalInfo(stat.slot, uint(link)).Return(mluLinkOpticalPresent[stat.slot][link], mluLinkOpticalTemp[stat.slot][link], mluLinkOpticalVolt[stat.slot][link], mluLinkOpticalTxpwr[stat.slot][link], mluLinkOpticalRxpwr[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkPortMode(stat.slot, uint(link)).Return(mluLinkPortMode[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkPPI(stat.slot, uint(link)).Return(mluLinkPortPPI[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkRemoteInfo(stat.slot, uint(link)).Return(mluLinkRemoteMcSn[stat.slot][link], mluLinkRemoteBaSn[stat.slot][link], mluLinkRemoteSlotID[stat.slot][link], mluLinkRemotePortID[stat.slot][link], mluLinkRemoteConnectType[stat.slot][link],
				mluLinkRemoteNcsUUID64[stat.slot][link], mluLinkRemoteIP[stat.slot][link], mluLinkRemoteMac[stat.slot][link], mluLinkRemoteUUID[stat.slot][link], mluLinkRemotePortName[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkSpeedInfo(stat.slot, uint(link)).Return(mluLinkSpeedValue[stat.slot][link], mluLinkSpeedFormat[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkStatus(stat.slot, uint(link)).Return(mluLinkStatusIsActive[stat.slot][link], mluLinkStatusSerdesState[stat.slot][link], mluLinkStatusCableState[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkState(stat.slot, uint(link)).Return(mluLinkState[stat.slot][link], mluLinkState[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkVersion(stat.slot, uint(link)).Return(mluLinkMajor[stat.slot][link], mluLinkMinor[stat.slot][link], mluLinkBuild[stat.slot][link], nil).AnyTimes()
		}
		mcndev.EXPECT().GetDeviceMLULinkPortNumber(stat.slot).Return(mluLinkPortNumber).AnyTimes()
		mcndev.EXPECT().GetDeviceProcessUtil(stat.slot).Return(pid[stat.slot], processIpuUtil[stat.slot], processJpuUtil[stat.slot], processMemUtil[stat.slot], processVpuDecodeUtil[stat.slot], processVpuEncodeUtil[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceProcessInfo(stat.slot).Return(pid[stat.slot], processPhysicalMemUsed[stat.slot], processVirtualMemUsed[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceParityError(stat.slot).Return(parityError[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceAddressSwaps(stat.slot).Return(addressSwap[stat.slot][0], addressSwap[stat.slot][1], addressSwap[stat.slot][1], addressSwap[stat.slot][1], addressSwap[stat.slot][1], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceCurrentInfo(stat.slot).Return(current[stat.slot][0], current[stat.slot][1], current[stat.slot][2], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceVoltageInfo(stat.slot).Return(voltage[stat.slot][0], current[stat.slot][1], current[stat.slot][2], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceOverTemperatureInfo(stat.slot).Return(overTemperatureInfo[stat.slot][0], overTemperatureInfo[stat.slot][1], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceOverTemperatureShutdownThreshold(stat.slot).Return(overTemperatureThreshold[stat.slot][0], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceOverTemperatureSlowdownThreshold(stat.slot).Return(overTemperatureThreshold[stat.slot][1], nil).AnyTimes()
		mcndev.EXPECT().GetDevicePerformanceThrottleReason(stat.slot).Return(throttleReason[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDevicePowerManagementDefaultLimitation(stat.slot).Return(powerDefaultLimit[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDevicePowerManagementLimitation(stat.slot).Return(powerCurrentLimit[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDevicePowerManagementLimitRange(stat.slot).Return(powerLimitRange[stat.slot][0], powerLimitRange[stat.slot][1], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceBAR4MemoryInfo(stat.slot).Return(bar4MemoryInfo[stat.slot][0], bar4MemoryInfo[stat.slot][1], bar4MemoryInfo[stat.slot][2], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceMaxPCIeInfo(stat.slot).Return(maxPCIeInfo[stat.slot][0], maxPCIeInfo[stat.slot][1], nil).AnyTimes()
		mcndev.EXPECT().GetAllMLUInstanceInfo(stat.slot).Return(stat.mimInfos, nil).AnyTimes()
		mcndev.EXPECT().GetAllSMluInfo(stat.slot).Return(stat.smluInfos, nil).AnyTimes()
		mcndev.EXPECT().GetDeviceMimProfileInfo(stat.slot).Return(mimProfileInfos, nil).AnyTimes()
		mcndev.EXPECT().GetDeviceMimProfileMaxInstanceCount(stat.slot, uint(mimProfileInfos[0].ProfileID)).Return(mimProfileMaxInstanceCount[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceSMluProfileIDInfo(stat.slot).Return(sMluProfileIDInfo[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceSMluProfileInfo(stat.slot, uint(sMluProfileIDInfo[stat.slot][0])).Return(sMluProfileInfo, sMluProfileTotal[stat.slot], sMluProfileRemain[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceComputeCapability(stat.slot).Return(computeCapability[0], computeCapability[1], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceComputeMode(stat.slot).Return(computeMode[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceChassisInfo(stat.slot).Return(chassisSn, chassisProductDate, chassisProductName, chassisVendorName, chassisPartNumber, chassisBmcIP, chassisNvme, chassisIb, chassisPsu, nil).AnyTimes()

		slots = append(slots, int(stat.slot))
	}

	mcndev.EXPECT().GetTopologyRelationship(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(2, nil).AnyTimes()
	mcndev.EXPECT().GetDeviceCndevVersion().Return(cndevVersion[0], cndevVersion[1], cndevVersion[2], nil).AnyTimes()
	sort.Ints(slots)
	mcndev.EXPECT().RegisterEventsHandleAndWait(slots, gomock.Any()).Return(nil).Times(2)

	// pres mock response
	pres := mock.NewPodResources(ctrl)
	pres.EXPECT().GetDeviceToPodInfo().Return(devicePodInfo, nil).AnyTimes()

	// host mock response
	host := mock.NewHost(ctrl)
	host.EXPECT().GetCPUStats().Return(hostCPUTotal, hostCPUIdle, nil).AnyTimes()
	host.EXPECT().GetMemoryStats().Return(hostMemTotal, hostMemFree, nil).AnyTimes()

	metricConfig := metrics.GetMetrics("../../examples/metrics.yaml", "")
	// collect metrics
	expectedMetrics := collectMetrics(node, rdmaDevice, mst, host, mcndev, pres, metricConfig)
	goldFile := "testdata/collect_metrics.json"
	out, err := json.MarshalIndent(expectedMetrics, "", "  ")
	assert.NoError(t, err)
	if *golden {
		assert.NoError(t, os.WriteFile(goldFile, out, 0o644))
	} else {
		expect, err := os.ReadFile(goldFile)
		assert.NoError(t, err)
		assert.JSONEq(t, string(expect), string(out))
	}

	// collect push metrics
	metricConfig = metrics.GetMetrics("../../examples/metrics-push.yaml", "")
	expectedMetrics = collectMetrics(node, rdmaDevice, mst, host, mcndev, pres, metricConfig)
	goldFile = "testdata/collect_push_metrics.json"
	out, err = json.MarshalIndent(expectedMetrics, "", "  ")
	assert.NoError(t, err)
	if *golden {
		assert.NoError(t, os.WriteFile(goldFile, out, 0o644))
	} else {
		expect, err := os.ReadFile(goldFile)
		assert.NoError(t, err)
		assert.JSONEq(t, string(expect), string(out))
	}
}

func collectMetrics(node string, rdmaDevice []rdmaDevice, mst map[string]MLUStat, host *mock.Host, cndv *mock.Cndev, pres *mock.PodResources, m map[string]metrics.CollectorMetrics) []testMetrics {
	bi := BaseInfo{
		host:       node,
		rdmaDevice: rdmaDevice,
	}
	cndevCollector := NewCndevCollector(m[Cndev], bi).(*cndevCollector)
	cndevCollector.client = cndv
	cndevCollector.lastXID = map[cndev.DeviceInfo]cndev.XIDInfo{
		{
			Device:            0,
			ComputeInstanceID: 0,
			MLUInstanceID:     0,
		}: {
			XIDBase10: 3155969,
			XID:       "0x000000000030d641",
		},
		{
			Device:            1,
			ComputeInstanceID: 0,
			MLUInstanceID:     0,
		}: {
			XIDBase10: 3155969,
			XID:       "0x000000000030d641",
		},
	}
	cndevCollector.mluXIDCounter = map[cndev.XIDInfo]int{
		{
			DeviceInfo: cndev.DeviceInfo{
				Device:            0,
				ComputeInstanceID: 0,
				MLUInstanceID:     0,
			},
			XID:       "0x000000000030d641",
			XIDBase10: 3155969,
		}: 5,
		{
			DeviceInfo: cndev.DeviceInfo{
				Device:            1,
				ComputeInstanceID: 0,
				MLUInstanceID:     0,
			},
			XID:       "0x000000000030d641",
			XIDBase10: 3155969,
		}: 10,
	}
	podResourcesCollector := NewPodResourcesCollector(m[PodResources], bi).(*podResourcesCollector)
	podResourcesCollector.client = pres
	hostCollector := NewHostCollector(m[Host], bi).(*hostCollector)
	hostCollector.client = host
	c := &Collectors{
		collectors: map[string]Collector{
			Cndev:        cndevCollector,
			PodResources: podResourcesCollector,
			Host:         hostCollector,
		},
		metrics: m,
	}
	for _, collector := range c.collectors {
		collector.init(mst)
	}
	ch := make(chan prometheus.Metric)
	go func() {
		c.Collect(ch)
		close(ch)
	}()

	expectedMetrics := []testMetrics{}
	for res := range ch {
		ms := dto.Metric{}
		res.Write(&ms)
		expectedMetrics = append(expectedMetrics, testMetrics{
			Desc:   res.Desc().String(),
			Metric: ms,
		})
	}
	sort.Slice(expectedMetrics, func(i, j int) bool {
		if expectedMetrics[i].Desc == expectedMetrics[j].Desc {
			for c, v := range expectedMetrics[i].Metric.Label {
				if *v.Name == *expectedMetrics[j].Metric.Label[c].Name && *v.Value != *expectedMetrics[j].Metric.Label[c].Value {
					return *expectedMetrics[i].Metric.Label[c].Value < *expectedMetrics[j].Metric.Label[c].Value
				}
			}
		}
		return expectedMetrics[i].Desc < expectedMetrics[j].Desc
	})
	return expectedMetrics
}
