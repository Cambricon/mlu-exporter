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
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/Cambricon/mlu-exporter/pkg/podresources"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	timeout   = 10 * time.Second
	socket    = "/var/lib/kubelet/pod-resources/kubelet.sock"
	resources = []string{
		"cambricon.com/mlu",
		"cambricon.com/whole-mlu-vf-count",
		"cambricon.com/half-mlu-vf-count",
		"cambricon.com/quarter-mlu-vf-count",
	}
	envShareSubString = "-_-"
	mimSubString      = "-mim-"
	smluSubString     = "-smlu-"
	uuidPrefix        = "MLU-"
	maxSize           = 1024 * 1024 * 16 // 16 Mb
	void              struct{}
)

func init() {
	registerCollector(PodResources, NewPodResourcesCollector)
}

type podResourcesCollector struct {
	baseInfo      BaseInfo
	client        podresources.PodResources
	devicePodInfo map[string]podresources.PodInfo
	fnMap         map[string]interface{}
	metrics       metrics.CollectorMetrics
	sharedInfo    map[string]MLUStat
}

func NewPodResourcesCollector(m metrics.CollectorMetrics, bi BaseInfo) Collector {
	c := &podResourcesCollector{
		baseInfo: bi,
		client:   podresources.NewPodResourcesClient(timeout, socket, resources, maxSize, bi.mode, bi.host, bi.client),
		metrics:  m,
	}
	c.fnMap = map[string]interface{}{
		Allocated: c.collectBoardAllocated,
		Container: c.collectMLUContainer,
	}
	log.Debugf("podResourcesCollector: %+v", c)
	return c
}

func (c *podResourcesCollector) init(info map[string]MLUStat) error {
	c.sharedInfo = info
	_, err := os.Stat(socket)
	return err
}

func (c *podResourcesCollector) updateMetrics(m metrics.CollectorMetrics) {
	c.metrics = m
}

func (c *podResourcesCollector) collect(ch chan<- prometheus.Metric) {
	info, err := c.client.GetDeviceToPodInfo()
	log.Debugf("GetDeviceToPodInfo: %+v", info)
	if err != nil {
		log.Errorf("GetDeviceToPodInfo err: %v", err)
		return
	}
	c.devicePodInfo = info
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

func (c *podResourcesCollector) collectBoardAllocated(ch chan<- prometheus.Metric, m metrics.Metric) {
	mlus := make(map[string]bool)
	for device := range c.devicePodInfo {
		uuid := device
		if strings.Contains(device, envShareSubString) {
			uuid = strings.Split(device, envShareSubString)[0]
		}
		mlus[uuid] = true
	}
	var model, driver string
	for _, info := range c.sharedInfo {
		model = info.model
		driver = info.driver
		break
	}
	labelValues := getLabelValues(m.Labels, labelInfo{stat: MLUStat{model: model, driver: driver}, host: c.baseInfo.host})
	ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, float64(len(mlus)), labelValues...)
}

func (c *podResourcesCollector) collectMLUContainer(ch chan<- prometheus.Metric, m metrics.Metric) {
	devs := map[string]struct{}{}
	for dev := range c.sharedInfo {
		devs[dev] = void
	}
	filterDup := map[string]struct{}{}
	for device, podInfo := range c.devicePodInfo {
		var vf, uuid string
		if c.baseInfo.mode == "dynamic-smlu" {
			for _, inf := range c.sharedInfo {
				if inf.slot == podInfo.Index {
					// Here uuid is the device id and vf is the instanceID.
					// We can ensure that the combination of uuid and instanceID is unique at the same moment.
					uuid = inf.uuid
					vf = podInfo.VF
					break
				}
			}
		} else {
			uuid = strings.TrimLeft(device, uuidPrefix)
			if strings.Contains(uuid, envShareSubString) {
				s := strings.Split(uuid, envShareSubString)
				uuid = s[0]
				vf = s[1]
			}
			if strings.Contains(uuid, mimSubString) {
				s := strings.Split(uuid, mimSubString)
				uuid = s[0]
				uid := strings.TrimLeft(s[1], uuidPrefix)
				st := c.sharedInfo[uuid]
				for _, inf := range st.mimInfos {
					if inf.UUID == uid {
						vf = strconv.Itoa(inf.InstanceID)
						break
					}
				}
			}
			if strings.Contains(uuid, smluSubString) {
				s := strings.Split(uuid, smluSubString)
				uuid = s[0]
				uid := strings.TrimLeft(s[1], uuidPrefix)
				st := c.sharedInfo[uuid]
				for _, inf := range st.smluInfos {
					if inf.UUID == uid {
						vf = strconv.Itoa(inf.InstanceID)
						break
					}
				}
			}
		}
		if _, ok := c.sharedInfo[uuid]; !ok {
			log.Warnf("failed to find mlu with uuid %s, skip", uuid)
			continue
		}
		if _, ok := filterDup[uuid+vf]; ok {
			continue
		}
		filterDup[uuid+vf] = void
		delete(devs, uuid)
		if c.baseInfo.mode == "dynamic-smlu" {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: c.sharedInfo[uuid], host: c.baseInfo.host, podInfo: podInfo, typ: "vcore", vf: vf})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, 1, labelValues...)
			labelValues = getLabelValues(m.Labels, labelInfo{stat: c.sharedInfo[uuid], host: c.baseInfo.host, podInfo: podInfo, typ: "vmemory", vf: vf})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, 1, labelValues...)
		} else {
			labelValues := getLabelValues(m.Labels, labelInfo{stat: c.sharedInfo[uuid], host: c.baseInfo.host, podInfo: podInfo, vf: vf})
			ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, 1, labelValues...)
		}
	}
	for dev := range devs {
		labelValues := getLabelValues(m.Labels, labelInfo{stat: c.sharedInfo[dev], host: c.baseInfo.host})
		ch <- prometheus.MustNewConstMetric(m.Desc, prometheus.GaugeValue, 0, labelValues...)
	}
}
