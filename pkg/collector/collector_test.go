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
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

var golden = flag.Bool("golden", false, "")

// verify
type testMetrics struct {
	Desc   string
	Metric dto.Metric
}

func TestCollect(t *testing.T) {
	var (
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
						Name:           "2m.16g",
						UUID:           "test-uuid1",
						InstanceID:     1,
						PlacementStart: 0,
						PlacementSize:  2,
					},
					{
						Name:           "2m.16g",
						UUID:           "test-uuid2",
						InstanceID:     2,
						PlacementStart: 2,
						PlacementSize:  2,
					},
				},
				link:       2,
				linkActive: map[int]bool{0: true, 1: true},
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
						Name:       "12.000m.78.784gb",
						UUID:       "test-uuid-1",
						InstanceID: 1,
						IpuUtil:    10,
						MemUsed:    3400000,
						MemTotal:   100000000,
					},
					{
						Name:       "6.000m.39.392gb",
						UUID:       "test-uuid-2",
						InstanceID: 2,
						IpuUtil:    12,
						MemUsed:    34000000,
						MemTotal:   1000000000,
					},
				},
				link:       2,
				linkActive: map[int]bool{0: true, 1: true},
			},
			uuid3: {
				slot:       2,
				model:      "MLU590",
				uuid:       uuid3,
				sn:         "sn3",
				mcu:        "v1.1.1",
				driver:     "v2.2.2",
				link:       2,
				linkActive: map[int]bool{0: true, 1: true},
			},
			uuid4: {
				slot:        3,
				model:       "MLU590",
				uuid:        uuid4,
				sn:          "sn4",
				mcu:         "v1.1.1",
				driver:      "v2.2.2",
				link:        2,
				linkActive:  map[int]bool{0: true, 1: true},
				crcDisabled: true,
			},
		}

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

		// health
		health = []int{1, 0, 1, 0}

		// parity
		parityError = []int{3, 0, 0, 0}

		// power
		power = []int{30, 31, 32, 33}

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

		pcieRead   = []int64{100, 200, 300, 400}
		pcieReplay = []uint32{100, 200, 300, 400}
		pcieWrite  = []int64{100, 200, 300, 400}

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
		processJpuUtil       = processIpuUtil
		processMemUtil       = processIpuUtil
		processVpuDecodeUtil = processIpuUtil
		processVpuEncodeUtil = processIpuUtil

		// chip
		chipCPUUtil     = []uint16{25, 25, 26, 26}
		chipCPUCoreUtil = [][]uint8{
			{0, 0, 100, 0},
			{0, 100, 0, 0},
			{0, 50, 50, 0},
			{25, 25, 25, 25},
		}
		armOsMemUsed  = []int64{1024, 4096, 512, 2048}
		armOsMemTotal = int64(10240)
		tinyCoreUtil  = [][]int{
			{10, 21},
			{12, 21},
			{15, 11},
			{21, 12},
		}
		numaNodeID   = []int{0, 0, 1, 1}
		ddrBandWidth = []float64{3.1, 3.2, 3.3, 3.4}
		ddrDataWidth = []int{1, 2, 3, 4}

		// frequency
		ddrFrequency = []int{100, 50, 150, 200}
		frequency    = []int{1000, 500, 1500, 2000}

		// codec
		videoCodecUtil = [][]int{
			{0, 2, 0, 1, 0, 1},
			{0, 1, 0, 3, 0, 1},
			{0, 1, 0, 1, 0, 2},
			{0, 4, 0, 1, 0, 1},
		}
		imageCodecUtil = videoCodecUtil

		// retired page
		cause      = []int{0, 1, 0, 1}
		pageCounts = []uint32{100, 10, 0, 0}

		// remapped rows
		correctRows   = []uint32{2, 1, 1, 2}
		failedRows    = []uint32{1, 1, 1, 1}
		pendingRows   = []uint32{2, 1, 2, 1}
		uncorrectRows = []uint32{1, 1, 1, 1}

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
		dramEccDbeCount          = []uint32{100, 100, 100, 100}
		dramEccSbeCount          = dramEccDbeCount
		sramEccDbeCount          = dramEccDbeCount
		sramEccParityCount       = dramEccDbeCount
		sramEccSbeCount          = dramEccDbeCount

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
		mluLinkPortMode = [][]int{
			{1, 1},
			{1, 1},
			{1, 1},
			{1, 1},
		}
		mluLinkSpeedValue = [][]float32{
			{12, 13},
			{22, 23},
			{32, 31},
			{42, 43},
		}
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
		mluLinkSpeedFormat                = mluLinkPortMode
		mluLinkStatusIsActive             = mluLinkPortMode
		mluLinkStatusSerdesState          = mluLinkPortMode
		mluLinkMajor                      = mluLinkCapabilityP2PTransfer
		mluLinkMinor                      = mluLinkCapabilityP2PTransfer
		mluLinkBuild                      = mluLinkCapabilityP2PTransfer
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mcndev mock response
	mcndev := mock.NewCndev(ctrl)
	mcndev.EXPECT().Init(false).Return(nil).AnyTimes()

	slots := []int{}
	for _, stat := range mst {
		mcndev.EXPECT().GetDeviceTemperature(stat.slot).Return(temperature[stat.slot], memTemperature[stat.slot], chipTemperature[stat.slot], clusterTemperatures[stat.slot], memDieTemperatures[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceHealth(stat.slot).Return(health[stat.slot], nil).AnyTimes()
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
		mcndev.EXPECT().GetDevicePCIeInfo(stat.slot).Return(pcieSlotID[stat.slot], pcieSubsystemID, pcieDeviceID, pcieVendor, pcieSubsystemVendor, pcieDomain, pcieBus[stat.slot], pcieDevice, pcieFunction, nil).AnyTimes()
		mcndev.EXPECT().GetDevicePCIeThroughput(stat.slot).Return(pcieRead[stat.slot], pcieWrite[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceCurrentPCIeInfo(stat.slot).Return(pcieSpeed, pcieWidth, nil).AnyTimes()
		mcndev.EXPECT().GetDevicePCIeReplayCount(stat.slot).Return(pcieReplay[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceCPUUtil(stat.slot).Return(chipCPUUtil[stat.slot], chipCPUCoreUtil[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceArmOsMemory(stat.slot).Return(armOsMemUsed[stat.slot], armOsMemTotal, nil).AnyTimes()
		mcndev.EXPECT().GetDeviceTinyCoreUtil(stat.slot).Return(tinyCoreUtil[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceNUMANodeID(stat.slot).Return(numaNodeID[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceDDRInfo(stat.slot).Return(ddrDataWidth[stat.slot], ddrBandWidth[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceFrequency(stat.slot).Return(frequency[stat.slot], ddrFrequency[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceRetiredPageInfo(stat.slot).Return(cause[stat.slot], pageCounts[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceRemappedRows(stat.slot).Return(correctRows[stat.slot], failedRows[stat.slot], pendingRows[stat.slot], uncorrectRows[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceVideoCodecUtil(stat.slot).Return(videoCodecUtil[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceImageCodecUtil(stat.slot).Return(imageCodecUtil[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceECCInfo(stat.slot).Return(eccAddressForbiddenError[stat.slot], eccCorrectedError[stat.slot], eccMultipleError[stat.slot], eccMultipleMultipleError[stat.slot], eccMultipleOneError[stat.slot], eccOneBitError[stat.slot], eccTotalError[stat.slot], eccUncorrectedError[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceCRCInfo(stat.slot).Return(d2dCRCError[stat.slot], d2dCRCErrorOverflow[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceMemEccCounter(stat.slot).Return(sramEccSbeCount[stat.slot], sramEccDbeCount[stat.slot], sramEccParityCount[stat.slot], dramEccSbeCount[stat.slot], dramEccDbeCount[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceHeartbeatCount(stat.slot).Return(heartbeatCount[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceCount().Return(deviceCount[stat.slot], nil).AnyTimes()
		for link := 0; link < mluLinkPortNumber; link++ {
			mcndev.EXPECT().GetDeviceMLULinkCapability(stat.slot, uint(link)).Return(mluLinkCapabilityP2PTransfer[stat.slot][link], mluLinkCapabilityInterlakenSerdes[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkCounter(stat.slot, uint(link)).Return(mluLinkCounterCntrReadByte[stat.slot][link], mluLinkCounterCntrReadPackage[stat.slot][link], mluLinkCounterCntrWriteByte[stat.slot][link], mluLinkCounterCntrWritePackage[stat.slot][link],
				mluLinkCounterErrCorrected[stat.slot][link], mluLinkCounterErrCRC24[stat.slot][link], mluLinkCounterErrCRC32[stat.slot][link], mluLinkCounterErrEccDouble[stat.slot][link], mluLinkCounterErrFatal[stat.slot][link], mluLinkCounterErrReplay[stat.slot][link],
				mluLinkCounterErrUncorrected[stat.slot][link], mluLinkCounterCntrCnpPackage[stat.slot][link], mluLinkCounterCntrPfcPackage[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkPortMode(stat.slot, uint(link)).Return(mluLinkPortMode[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkSpeedInfo(stat.slot, uint(link)).Return(mluLinkSpeedValue[stat.slot][link], mluLinkSpeedFormat[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkStatus(stat.slot, uint(link)).Return(mluLinkStatusIsActive[stat.slot][link], mluLinkStatusSerdesState[stat.slot][link], nil).AnyTimes()
			mcndev.EXPECT().GetDeviceMLULinkVersion(stat.slot, uint(link)).Return(mluLinkMajor[stat.slot][link], mluLinkMinor[stat.slot][link], mluLinkBuild[stat.slot][link], nil).AnyTimes()
		}
		mcndev.EXPECT().GetDeviceProcessUtil(stat.slot).Return(pid[stat.slot], processIpuUtil[stat.slot], processJpuUtil[stat.slot], processMemUtil[stat.slot], processVpuDecodeUtil[stat.slot], processVpuEncodeUtil[stat.slot], nil).AnyTimes()
		mcndev.EXPECT().GetDeviceParityError(stat.slot).Return(parityError[stat.slot], nil).AnyTimes()

		slots = append(slots, int(stat.slot))
	}
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
	expectedMetrics := collectMetrics(node, mst, host, mcndev, pres, metricConfig)
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
	expectedMetrics = collectMetrics(node, mst, host, mcndev, pres, metricConfig)
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

func collectMetrics(node string, mst map[string]MLUStat, host *mock.Host, cndv *mock.Cndev, pres *mock.PodResources, m map[string]metrics.CollectorMetrics) []testMetrics {
	bi := BaseInfo{
		host: node,
	}
	cndevCollector := NewCndevCollector(m[Cndev], bi).(*cndevCollector)
	cndevCollector.client = cndv
	cndevCollector.lastXID = map[cndev.DeviceInfo]int64{
		{
			Device:            0,
			ComputeInstanceID: 0,
			MLUInstanceID:     0,
		}: 3155969,
		{
			Device:            1,
			ComputeInstanceID: 0,
			MLUInstanceID:     0,
		}: 3155969,
	}
	cndevCollector.mluXIDCounter = map[cndev.XIDInfo]int{
		{
			DeviceInfo: cndev.DeviceInfo{
				Device:            0,
				ComputeInstanceID: 0,
				MLUInstanceID:     0,
			},
			XID: 3155969,
		}: 5,
		{
			DeviceInfo: cndev.DeviceInfo{
				Device:            1,
				ComputeInstanceID: 0,
				MLUInstanceID:     0,
			},
			XID: 3155969,
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
