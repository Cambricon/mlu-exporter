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
	"os"
	"sort"
	"sync"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type Callback struct {
	cndevcli   cndev.Cndev
	metric     *metrics.Metric
	metricName string
	sharedInfo map[string]MLUStat
	host       string
	retryTimes int
	logFile    string

	currentXID  cndev.XIDInfoWithTimestamp
	metricMutex sync.Mutex
	xidMutex    sync.Mutex
}

func NewCallback(
	metricConfig map[string]metrics.CollectorMetrics,
	info map[string]MLUStat,
	metricName, host, logFile string,
	retryTimes int,
) (*Callback, error) {

	// xid callback was disabled
	disabled := true
	for _, i := range info {
		if !i.xidCallbackDisabled {
			disabled = false
			break
		}
	}
	if disabled {
		return nil, nil
	}

	m := filterMetric(metricConfig, metricName)
	if m == nil {
		return nil, nil
	}
	cndevcli := cndev.NewCndevClient()
	if err := cndevcli.Init(false); err != nil {
		return nil, err
	}
	c := &Callback{
		cndevcli:   cndevcli,
		sharedInfo: info,
		metric:     m,
		metricName: metricName,
		host:       host,
		retryTimes: retryTimes,
		logFile:    logFile,
	}

	return c, nil
}

func filterMetric(cfg map[string]metrics.CollectorMetrics, metricName string) *metrics.Metric {
	for _, metricsMap := range cfg {
		for _, value := range metricsMap {
			fqName := value.Name
			if fqName != metricName {
				continue
			}
			return &value
		}
	}
	return nil
}

func (c *Callback) Start(f func() error) {
	slots := []int{}
	for _, stat := range c.sharedInfo {
		slots = append(slots, int(stat.slot))
	}
	sort.Ints(slots)

	ch := make(chan cndev.XIDInfoWithTimestamp, 10)
	if err := c.cndevcli.RegisterEventsHandleAndWait(slots, ch); err != nil {
		log.Errorln(errors.Wrap(err, "register event handle"))
		return
	}

	for event := range ch {
		c.xidMutex.Lock()
		c.currentXID = event
		c.work(f)
		c.currentXID = cndev.XIDInfoWithTimestamp{}
		c.xidMutex.Unlock()
	}
}

func (c *Callback) work(f func() error) {
	for i := 0; i < c.retryTimes; i++ {
		err := f()
		if err == nil {
			return
		}
		log.Errorln(errors.Wrap(err, "push failed"))
	}

	if c.logFile == "" {
		return
	}

	var file *os.File
	file, err := os.OpenFile(c.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorln(errors.Wrapf(err, "open log file: %v failed", c.logFile))
	}
	defer file.Close()

	entry := struct {
		Event cndev.XIDInfoWithTimestamp `json:"event"`
		Retry int                        `json:"retry"`
	}{
		Event: c.currentXID,
		Retry: c.retryTimes,
	}
	entryS, err := json.Marshal(entry)
	if err != nil {
		log.Errorln(errors.Wrap(err, "marshal failed"))
		return
	}
	_, err = file.WriteString(string(entryS) + "\n")
	if err != nil {
		log.Errorln(errors.Wrapf(err, "write log file: %v failed", c.logFile))
	}
}

func (c *Callback) UpdateMetrics(metricConfig map[string]metrics.CollectorMetrics) {
	m := filterMetric(metricConfig, c.metricName)
	c.metricMutex.Lock()
	defer c.metricMutex.Unlock()
	c.metric = m
}

func (c *Callback) Collect(ch chan<- prometheus.Metric) {
	if c.metric == nil {
		return
	}
	slotInfo := map[uint]MLUStat{}
	for _, stat := range c.sharedInfo {
		slotInfo[stat.slot] = stat
	}
	event := c.currentXID
	labelValues := getLabelValues(c.metric.Labels, labelInfo{
		stat:    slotInfo[uint(event.Device)],
		host:    c.host,
		xidInfo: event,
	})
	ch <- prometheus.MustNewConstMetric(c.metric.Desc, prometheus.GaugeValue, 1, labelValues...)
}

func (c *Callback) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.metric.Desc
}
