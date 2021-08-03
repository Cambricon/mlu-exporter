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
	"log"
	"os"
	"strings"
	"time"

	"github.com/Cambricon/mlu-exporter/pkg/podresources"
	"github.com/prometheus/client_golang/prometheus"
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
	sriovSubString    = "--fake--"
	envShareSubString = "-_-"
	uuidPrefix        = "MLU-"
	maxSize           = 1024 * 1024 * 16 // 16 Mb
)

func init() {
	registerCollector(PodResources, NewPodResourcesCollector)
}

type podResourcesCollector struct {
	metrics       collectorMetrics
	host          string
	devicePodInfo map[string]podresources.PodInfo
	client        podresources.PodResources
	sharedInfo    map[string]mluStat
	fnMap         map[string]interface{}
}

func NewPodResourcesCollector(m collectorMetrics, host string) Collector {
	c := &podResourcesCollector{
		metrics: m,
		host:    host,
		client:  podresources.NewPodResourcesClient(timeout, socket, resources, maxSize),
	}
	c.fnMap = map[string]interface{}{
		Allocated: c.collectBoardAllocated,
		Container: c.collectMLUContainer,
	}
	return c
}

func (c *podResourcesCollector) init(info map[string]mluStat) error {
	c.sharedInfo = info
	_, err := os.Stat(socket)
	return err
}

func (c *podResourcesCollector) updateMetrics(m collectorMetrics) {
	c.metrics = m
}

func (c *podResourcesCollector) collect(ch chan<- prometheus.Metric) {
	info, err := c.client.GetDeviceToPodInfo()
	if err != nil {
		log.Printf("GetDeviceToPodInfo err: %v", err)
		return
	}
	c.devicePodInfo = info
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

func (c *podResourcesCollector) collectBoardAllocated(ch chan<- prometheus.Metric, m metric) {
	mlus := make(map[string]bool)
	for device := range c.devicePodInfo {
		uuid := device
		for _, sub := range []string{envShareSubString, sriovSubString} {
			if strings.Contains(device, sub) {
				uuid = strings.Split(device, sub)[0]
				break
			}
		}
		mlus[uuid] = true
	}
	var model, driver string
	for _, info := range c.sharedInfo {
		model = info.model
		driver = info.driver
		break
	}
	labelValues := getLabelValues(m.labels, labelInfo{stat: mluStat{model: model, driver: driver}, host: c.host})
	ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(len(mlus)), labelValues...)
}

func (c *podResourcesCollector) collectMLUContainer(ch chan<- prometheus.Metric, m metric) {
	for device, podInfo := range c.devicePodInfo {
		vf := ""
		uuid := strings.TrimLeft(device, uuidPrefix)
		if strings.Contains(uuid, envShareSubString) {
			uuid = strings.Split(uuid, envShareSubString)[0]
		}
		if strings.Contains(device, sriovSubString) {
			s := strings.Split(uuid, sriovSubString)
			uuid = s[0]
			vf = s[1]
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: c.sharedInfo[uuid], host: c.host, podInfo: podInfo, vf: vf})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, 1, labelValues...)
	}
}
