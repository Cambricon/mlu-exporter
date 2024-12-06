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
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Callback struct {
	cndevcli   cndev.Cndev
	client     *HTTPClient
	metric     *metrics.Metric
	metricName string
	sharedInfo map[string]MLUStat
	retryTimes int
	logFile    string

	metricMutex sync.Mutex
}

func NewCallback(
	importURL string,
	metricConfig map[string]metrics.CollectorMetrics,
	info map[string]MLUStat,
	metricName, metricPrefix, host, logFile string,
	retryTimes int,
	jobName string,
) (*Callback, error) {

	// xid callback was disabled
	disabled := true
	for _, i := range info {
		if !i.cndevInterfaceDisabled["xidCallbackDisabled"] {
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
	fqName := metricName
	if metricPrefix != "" {
		fqName = fmt.Sprintf("%s_%s", metricPrefix, fqName)
	}
	client, err := NewClient(importURL, info, m.Labels, host, fqName, jobName)
	if err != nil {
		return nil, err
	}
	cndevcli := cndev.NewCndevClient()
	if err := cndevcli.Init(false); err != nil {
		return nil, err
	}
	c := &Callback{
		cndevcli:   cndevcli,
		client:     client,
		sharedInfo: info,
		metric:     m,
		metricName: metricName,
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

func (c *Callback) Start() {
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
		c.work(event)
	}
}

func (c *Callback) work(event cndev.XIDInfoWithTimestamp) {
	err := c.client.PushWithRetries(event, c.retryTimes)
	if err == nil {
		return
	}
	log.Errorln(errors.Wrap(err, "push failed"))

	if c.logFile == "" {
		return
	}

	file, err := os.OpenFile(c.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorln(errors.Wrapf(err, "open log file: %v failed", c.logFile))
	}
	defer file.Close()

	entry := struct {
		Event cndev.XIDInfoWithTimestamp `json:"event"`
		Retry int                        `json:"retry"`
	}{
		Event: event,
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
