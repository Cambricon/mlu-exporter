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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Cambricon/mlu-exporter/pkg/metrics"
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
	registerCollector(metrics.PodResources, NewPodResourcesCollector)
}

type podResourcesCollector struct {
	metrics       collectorMetrics
	host          string
	devicePodInfo map[string]podresources.PodInfo
	client        podresources.PodResources
}

func NewPodResourcesCollector(m collectorMetrics, host string) Collector {
	return &podResourcesCollector{
		metrics: m,
		host:    host,
		client:  podresources.NewPodResourcesClient(timeout, socket, resources, maxSize),
	}
}

func (c *podResourcesCollector) init() error {
	_, err := os.Stat(socket)
	return err
}

func (c *podResourcesCollector) updateMetrics(m collectorMetrics) {
	c.metrics = m
}

func (c *podResourcesCollector) collect(ch chan<- prometheus.Metric, mluInfo map[string]mluStat) {
	info, err := c.client.GetDeviceToPodInfo()
	check(err)
	c.devicePodInfo = info

	for name, metric := range c.metrics {
		switch name {
		case metrics.ContainerMLUVFUtil:
			c.collectVFUtil(ch, metric, mluInfo)
		case metrics.BoardAllocated:
			c.collectBoardAllocated(ch, metric, mluInfo)
		case metrics.ContainerMLUUtil:
			c.collectMLUUtil(ch, metric, mluInfo)
		case metrics.ContainerMLUMemUtil:
			c.collectMLUMemUtil(ch, metric, mluInfo)
		case metrics.ContainerMLUBoardPower:
			c.collectMLUBoardPower(ch, metric, mluInfo)
		case metrics.MLUContainer:
			c.collectMLUContainer(ch, metric, mluInfo)
		}
	}
}

func (c *podResourcesCollector) collectMLUUtil(ch chan<- prometheus.Metric, m metric, mluInfo map[string]mluStat) {
	for device, podInfo := range c.devicePodInfo {
		if strings.Contains(device, sriovSubString) {
			continue
		}
		uuid := strings.TrimLeft(device, uuidPrefix)
		if strings.Contains(uuid, envShareSubString) {
			uuid = strings.Split(uuid, envShareSubString)[0]
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: mluInfo[uuid], host: c.host, podInfo: podInfo})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(mluInfo[uuid].boardUtil), labelValues...)
	}
}

func (c *podResourcesCollector) collectMLUMemUtil(ch chan<- prometheus.Metric, m metric, mluInfo map[string]mluStat) {
	for device, podInfo := range c.devicePodInfo {
		if strings.Contains(device, sriovSubString) {
			continue
		}
		uuid := strings.TrimLeft(device, uuidPrefix)
		if strings.Contains(uuid, envShareSubString) {
			uuid = strings.Split(uuid, envShareSubString)[0]
		}
		util := float64(mluInfo[uuid].memUsed) / float64(mluInfo[uuid].memTotal) * 100
		labelValues := getLabelValues(m.labels, labelInfo{stat: mluInfo[uuid], host: c.host, podInfo: podInfo})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, util, labelValues...)
	}
}

func (c *podResourcesCollector) collectMLUBoardPower(ch chan<- prometheus.Metric, m metric, mluInfo map[string]mluStat) {
	for device, podInfo := range c.devicePodInfo {
		if strings.Contains(device, sriovSubString) {
			continue
		}
		uuid := strings.TrimLeft(device, uuidPrefix)
		if strings.Contains(uuid, envShareSubString) {
			uuid = strings.Split(uuid, envShareSubString)[0]
		}
		power := float64(mluInfo[uuid].power)
		labelValues := getLabelValues(m.labels, labelInfo{stat: mluInfo[uuid], host: c.host, podInfo: podInfo})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, power, labelValues...)
	}
}

func (c *podResourcesCollector) collectBoardAllocated(ch chan<- prometheus.Metric, m metric, mluInfo map[string]mluStat) {
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
	for _, info := range mluInfo {
		model = info.model
		driver = info.driver
		break
	}
	labelValues := getLabelValues(m.labels, labelInfo{stat: mluStat{model: model, driver: driver}, host: c.host})
	ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(len(mlus)), labelValues...)
}

func (c *podResourcesCollector) collectMLUContainer(ch chan<- prometheus.Metric, m metric, mluInfo map[string]mluStat) {
	for device, podInfo := range c.devicePodInfo {
		if strings.Contains(device, sriovSubString) {
			continue
		}
		uuid := strings.TrimLeft(device, uuidPrefix)
		if strings.Contains(uuid, envShareSubString) {
			uuid = strings.Split(uuid, envShareSubString)[0]
		}
		labelValues := getLabelValues(m.labels, labelInfo{stat: mluInfo[uuid], host: c.host, podInfo: podInfo})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, 1, labelValues...)
	}
}

func (c *podResourcesCollector) collectVFUtil(ch chan<- prometheus.Metric, m metric, mluInfo map[string]mluStat) {
	for device, podInfo := range c.devicePodInfo {
		if !strings.Contains(device, sriovSubString) {
			continue
		}
		uuid, vfNum, vfID, err := getSriovDeviceInfo(device, mluInfo)
		check(err)
		util, err := calcVFUtil(mluInfo[uuid].coreUtil, vfNum, vfID)
		check(err)
		labelValues := getLabelValues(m.labels, labelInfo{stat: mluInfo[uuid], host: c.host, vfID: vfID, podInfo: podInfo})
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, util, labelValues...)
	}
}

func getSriovDeviceInfo(device string, mluInfo map[string]mluStat) (string, int, int, error) {
	uuid := strings.TrimLeft(device, uuidPrefix)
	uuid = strings.Split(uuid, sriovSubString)[0]
	vf := strings.Split(device, sriovSubString)[1]
	vfID, err := strconv.ParseInt(vf, 10, 64)
	if err != nil {
		return "", 0, 0, err
	}
	pcie := mluInfo[uuid].pcie
	path := "/sys/bus/pci/devices/" + pcie + "/sriov_numvfs"
	vfNum, err := getNumFromFile(path)
	if err != nil {
		return "", 0, 0, err
	}
	return uuid, vfNum, int(vfID), nil
}

func getNumFromFile(path string) (int, error) {
	output, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	output = bytes.Trim(output, "\n")
	num, err := strconv.ParseInt(string(output), 10, 64)
	return int(num), err
}

func calcVFUtil(utils []uint, vfNum int, vfID int) (float64, error) {
	sum := uint(0)
	offset := 0
	switch vfNum {
	case 1, 2, 4:
		offset = len(utils) / vfNum
	default:
		return 0, fmt.Errorf("invalid vfNum %d", vfNum)
	}
	for i := (vfID - 1) * offset; i < vfID*offset; i++ {
		sum += utils[i]
	}
	return float64(sum) / float64(offset), nil
}
