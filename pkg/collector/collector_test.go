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
			pcie:   "0000:1a:00.0",
			mcu:    "v1.1.1",
			driver: "v2.2.2",
		},
		"uuid2": {
			slot:   1,
			model:  "MLU270-X5K",
			uuid:   "uuid2",
			sn:     "sn2",
			pcie:   "0000:1b:00.0",
			mcu:    "v1.1.1",
			driver: "v2.2.2",
		},
		"uuid3": {
			slot:   2,
			model:  "MLU270-X5K",
			uuid:   "uuid3",
			sn:     "sn3",
			pcie:   "0000:1c:00.0",
			mcu:    "v1.1.1",
			driver: "v2.2.2",
		},
		"uuid4": {
			slot:   3,
			model:  "MLU270-X5K",
			uuid:   "uuid4",
			sn:     "sn4",
			pcie:   "0000:1d:00.0",
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
		boardUtil           uint
		coreUtil            []uint
		memUsed             uint
		memTotal            uint
		virtualMemUsed      []uint
		virtualMemTotal     uint
		temperature         uint
		clusterTemperatures []uint
		health              uint
		power               uint
		pcieRead            uint64
		pcieWrite           uint64
		dramRead            uint64
		dramWrite           uint64
		mlulinkRead         uint64
		mlulinkWrite        uint64
		vfState             int
	}{
		"uuid1": {
			11,
			[]uint{11, 11, 11, 11, 13, 13, 13, 13, 15, 15, 15, 15, 17, 17, 17, 17},
			21,
			1000,
			[]uint{},
			1000,
			31,
			[]uint{22, 23, 24, 26},
			1,
			41,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			0,
		},
		"uuid2": {
			11,
			[]uint{11, 11, 11, 11, 13, 13, 13, 13, 15, 15, 15, 15, 17, 17, 17, 17},
			21,
			1000,
			[]uint{},
			1000,
			31,
			[]uint{22, 23, 24, 26},
			1,
			41,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			0,
		},
		"uuid3": {
			11,
			[]uint{11, 11, 11, 11, 13, 13, 13, 13, 15, 15, 15, 15, 17, 17, 17, 17},
			21,
			1000,
			[]uint{11, 13, 15, 17},
			250,
			31,
			[]uint{22, 23, 24, 26},
			1,
			41,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			15,
		},
		"uuid4": {
			11,
			[]uint{11, 11, 11, 11, 13, 13, 13, 13, 15, 15, 15, 15, 17, 17, 17, 17},
			21,
			1000,
			[]uint{12, 16},
			500,
			31,
			[]uint{22, 23, 24, 26},
			1,
			41,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			1073741824,
			3,
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
		for vf := 0; vf < len(mluMetrics[stat.uuid].virtualMemUsed); vf++ {
			cndev.EXPECT().GetDeviceMemory(uint((vf+1)<<8|int(stat.slot))).Return(mluMetrics[mlu].memUsed, mluMetrics[mlu].memTotal, mluMetrics[mlu].virtualMemUsed[vf], mluMetrics[mlu].virtualMemTotal, nil).Times(1)
		}
		cndev.EXPECT().GetDeviceMemory(stat.slot).Return(mluMetrics[mlu].memUsed, mluMetrics[mlu].memTotal, uint(0), uint(0), nil).Times(3)
		cndev.EXPECT().GetDeviceUtil(stat.slot).Return(mluMetrics[mlu].boardUtil, mluMetrics[mlu].coreUtil, nil).Times(4)
		cndev.EXPECT().GetDeviceFanSpeed(stat.slot).Return(uint(0), nil).Times(1)
		cndev.EXPECT().GetDevicePower(stat.slot).Return(mluMetrics[mlu].power, nil).Times(1)
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
		Temperature:         cndevCollector.collectTemperature,
		Health:              cndevCollector.collectHealth,
		MemTotal:            cndevCollector.collectMemTotal,
		MemUsed:             cndevCollector.collectMemUsed,
		MemUtil:             cndevCollector.collectMemUtil,
		Util:                cndevCollector.collectUtil,
		CoreUtil:            cndevCollector.collectCoreUtil,
		FanSpeed:            cndevCollector.collectFanSpeed,
		PowerUsage:          cndevCollector.collectPowerUsage,
		Capacity:            cndevCollector.collectCapacity,
		Usage:               cndevCollector.collectUsage,
		Version:             cndevCollector.collectVersion,
		VirtualMemUtil:      cndevCollector.collectVirtualMemUtil,
		VirtualFunctionUtil: cndevCollector.collectVirtualFunctionUtil,
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
		for vf := 0; vf < len(mluMetrics[stat.uuid].virtualMemUsed); vf++ {
			memUtil := float64(mluMetrics[mlu].virtualMemUsed[vf]) / float64(mluMetrics[mlu].virtualMemTotal) * 100
			util, _ := calcVFUtil(mluMetrics[mlu].coreUtil, len(mluMetrics[stat.uuid].virtualMemUsed), vf+1)
			expected[m[Cndev][VirtualMemUtil].desc] = append(expected[m[Cndev][VirtualMemUtil].desc], append(l, fmt.Sprintf("%d", vf+1), fmt.Sprintf("%v", memUtil)))
			expected[m[Cndev][VirtualFunctionUtil].desc] = append(expected[m[Cndev][VirtualFunctionUtil].desc], append(l, fmt.Sprintf("%d", vf+1), fmt.Sprintf("%v", util)))
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
		values := []string{fmt.Sprintf("%v", metric.GetGauge().GetValue())}
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
