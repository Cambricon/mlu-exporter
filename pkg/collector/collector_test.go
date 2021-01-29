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
	"testing"

	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/Cambricon/mlu-exporter/pkg/mock"
	"github.com/Cambricon/mlu-exporter/pkg/podresources"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

var (
	cardNum             = uint(4)
	sn                  = []string{"sn0", "sn1", "sn2", "sn3"}
	mcu                 = "v1.1.1"
	driver              = "v2.2.2"
	model               = "MLU270-X5K"
	mluType             = "MLU270-X5K"
	pcie                = "0000:1c:00.0"
	boardUtil           = []uint{11, 12, 13, 14}
	coreUtil            = []uint{11, 12, 13, 14}
	memUsed             = []uint{21, 22, 23, 24}
	memTotal            = uint(1000)
	temperature         = []uint{50, 51, 52, 53}
	clusterTemperatures = []uint{22, 23, 24, 26}
	health              = uint(1)
	power               = []uint{44, 45, 46, 46}
	nodeName            = "machine1"
	devicePodInfo       = map[string]podresources.PodInfo{
		"MLU-sn1": {
			Pod:       "pod1",
			Namespace: "namespace1",
			Container: "container1",
		},
		"MLU-sn2-_-2": {
			Pod:       "pod2",
			Namespace: "namespace2",
			Container: "container2",
		},
	}
)

func TestInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cndev := mock.NewCndev(ctrl)
	cndev.EXPECT().Init().Return(nil).Times(1)
	c := &Collectors{
		collectors: map[string]Collector{
			"cndev": &cndevCollector{
				cndev: cndev,
			},
		},
	}
	c.init()
}

func TestCollect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cndevOuter := mock.NewCndev(ctrl)
	cndevOuter.EXPECT().GetDeviceCount().Return(uint(cardNum), nil).Times(1)
	for i := uint(0); i < cardNum; i++ {
		cndevOuter.EXPECT().GetDeviceSN(i).Return(sn[i], nil).Times(1)
		cndevOuter.EXPECT().GetDeviceModel(i).Return(model).Times(1)
		cndevOuter.EXPECT().GetDevicePCIeID(i).Return(pcie, nil).Times(1)
		cndevOuter.EXPECT().GetDeviceUtil(i).Return(boardUtil[i], coreUtil, nil).Times(1)
		cndevOuter.EXPECT().GetDeviceMemory(i).Return(memUsed[i], memTotal, nil).Times(1)
		cndevOuter.EXPECT().GetDevicePower(i).Return(power[i], nil).Times(1)
		cndevOuter.EXPECT().GetDeviceVersion(i).Return(mcu, driver, nil).Times(1)
	}
	cndevInner := mock.NewCndev(ctrl)
	for i := uint(0); i < cardNum; i++ {
		cndevInner.EXPECT().GetDeviceTemperature(i).Return(temperature[i], clusterTemperatures, nil).Times(1)
		cndevInner.EXPECT().GetDeviceHealth(i).Return(health, nil).Times(1)
		cndevInner.EXPECT().GetDeviceFanSpeed(i).Return(uint(0), nil).Times(1)
	}
	podresources := mock.NewPodResources(ctrl)
	podresources.EXPECT().GetDeviceToPodInfo().Return(devicePodInfo, nil).Times(1)
	m := getMetrics(metrics.GetOrDie("../../examples/metrics.yaml"), "")
	expected := map[*prometheus.Desc][]metricValue{}
	for i := uint(0); i < cardNum; i++ {
		l := []string{sn[i], model, mluType, nodeName, fmt.Sprintf("%d", i), mcu, driver}
		expected[m[metrics.Cndev][metrics.Temperature].desc] = append(expected[m[metrics.Cndev][metrics.Temperature].desc], append(l, "", fmt.Sprintf("%v", temperature[i])))
		for cluster, temp := range clusterTemperatures {
			expected[m[metrics.Cndev][metrics.Temperature].desc] = append(expected[m[metrics.Cndev][metrics.Temperature].desc], append(l, fmt.Sprintf("%d", cluster), fmt.Sprintf("%v", temp)))
		}
		expected[m[metrics.Cndev][metrics.BoardHealth].desc] = append(expected[m[metrics.Cndev][metrics.BoardHealth].desc], append(l, fmt.Sprintf("%v", health)))
		expected[m[metrics.Cndev][metrics.MemTotal].desc] = append(expected[m[metrics.Cndev][metrics.MemTotal].desc], append(l, fmt.Sprintf("%v", float64(memTotal*1024*1024))))
		expected[m[metrics.Cndev][metrics.MemUsed].desc] = append(expected[m[metrics.Cndev][metrics.MemUsed].desc], append(l, fmt.Sprintf("%v", float64(memUsed[i]*1024*1024))))
		expected[m[metrics.Cndev][metrics.MemUtil].desc] = append(expected[m[metrics.Cndev][metrics.MemUtil].desc], append(l, fmt.Sprintf("%v", float64(memUsed[i])/float64(memTotal)*100)))
		expected[m[metrics.Cndev][metrics.BoardUtil].desc] = append(expected[m[metrics.Cndev][metrics.BoardUtil].desc], append(l, fmt.Sprintf("%v", boardUtil[i])))
		expected[m[metrics.Cndev][metrics.BoardCapacity].desc] = append(expected[m[metrics.Cndev][metrics.BoardCapacity].desc], append(l, fmt.Sprintf("%v", capacity[model])))
		expected[m[metrics.Cndev][metrics.BoardUsage].desc] = append(expected[m[metrics.Cndev][metrics.BoardUsage].desc], append(l, fmt.Sprintf("%v", float64(capacity[model])*float64(boardUtil[i])/100)))
		for core, u := range coreUtil {
			expected[m[metrics.Cndev][metrics.CoreUtil].desc] = append(expected[m[metrics.Cndev][metrics.CoreUtil].desc], append(l, fmt.Sprintf("%d", core), fmt.Sprintf("%v", u)))
		}
		expected[m[metrics.Cndev][metrics.FanSpeed].desc] = append(expected[m[metrics.Cndev][metrics.FanSpeed].desc], append(l, fmt.Sprintf("%v", -1)))
		expected[m[metrics.Cndev][metrics.BoardPower].desc] = append(expected[m[metrics.Cndev][metrics.BoardPower].desc], append(l, fmt.Sprintf("%v", power[i])))
		expected[m[metrics.Cndev][metrics.BoardVersion].desc] = append(expected[m[metrics.Cndev][metrics.BoardVersion].desc], append(l, "1"))
	}
	expected[m[metrics.PodResources][metrics.BoardAllocated].desc] = append(expected[m[metrics.PodResources][metrics.BoardAllocated].desc], []string{model, mluType, driver, nodeName, "2"})
	expected[m[metrics.PodResources][metrics.ContainerMLUUtil].desc] = append(expected[m[metrics.PodResources][metrics.ContainerMLUUtil].desc], []string{"1", model, mluType, sn[1], nodeName, "namespace1", "pod1", "container1", fmt.Sprintf("%v", boardUtil[1])})
	expected[m[metrics.PodResources][metrics.ContainerMLUUtil].desc] = append(expected[m[metrics.PodResources][metrics.ContainerMLUUtil].desc], []string{"2", model, mluType, sn[2], nodeName, "namespace2", "pod2", "container2", fmt.Sprintf("%v", boardUtil[2])})
	expected[m[metrics.PodResources][metrics.ContainerMLUMemUtil].desc] = append(expected[m[metrics.PodResources][metrics.ContainerMLUMemUtil].desc], []string{"1", model, mluType, sn[1], nodeName, "namespace1", "pod1", "container1", fmt.Sprintf("%v", float64(memUsed[1])/float64(memTotal)*100)})
	expected[m[metrics.PodResources][metrics.ContainerMLUMemUtil].desc] = append(expected[m[metrics.PodResources][metrics.ContainerMLUMemUtil].desc], []string{"2", model, mluType, sn[2], nodeName, "namespace2", "pod2", "container2", fmt.Sprintf("%v", float64(memUsed[2])/float64(memTotal)*100)})
	expected[m[metrics.PodResources][metrics.ContainerMLUBoardPower].desc] = append(expected[m[metrics.PodResources][metrics.ContainerMLUBoardPower].desc], []string{"1", model, mluType, sn[1], nodeName, "namespace1", "pod1", "container1", fmt.Sprintf("%v", power[1])})
	expected[m[metrics.PodResources][metrics.ContainerMLUBoardPower].desc] = append(expected[m[metrics.PodResources][metrics.ContainerMLUBoardPower].desc], []string{"2", model, mluType, sn[2], nodeName, "namespace2", "pod2", "container2", fmt.Sprintf("%v", power[2])})
	expected[m[metrics.PodResources][metrics.MLUContainer].desc] = append(expected[m[metrics.PodResources][metrics.MLUContainer].desc], []string{"1", sn[1], model, mluType, mcu, driver, nodeName, "namespace1", "pod1", "container1", "1"})
	expected[m[metrics.PodResources][metrics.MLUContainer].desc] = append(expected[m[metrics.PodResources][metrics.MLUContainer].desc], []string{"2", sn[2], model, mluType, mcu, driver, nodeName, "namespace2", "pod2", "container2", "1"})
	for desc, values := range expected {
		sortMetricValues(values)
		expected[desc] = values
	}
	c := &Collectors{
		collectors: map[string]Collector{
			"cndev": &cndevCollector{
				host:    nodeName,
				cndev:   cndevInner,
				metrics: m["cndev"],
			},
			"podresources": &podResourcesCollector{
				metrics: m["podresources"],
				host:    nodeName,
				client:  podresources,
			},
		},
		metrics: m,
		cndev:   cndevOuter,
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
	for desc, values := range got {
		sortMetricValues(values)
		got[desc] = values
	}
	assert.Equal(t, expected, got)
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
