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

package metrics

import (
	"fmt"
	"sort"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type CollectorMetrics map[string]Metric

type Metric struct {
	Name   string
	Desc   *prometheus.Desc
	Labels []string
	Push   bool
}

var listeners []func(map[string]CollectorMetrics)

func RegisterWatcher(f func(map[string]CollectorMetrics)) {
	listeners = append(listeners, f)
}

func GetMetrics(configFile string, prefix string) map[string]CollectorMetrics {
	return getMetrics(configFile, prefix)
}

func getMetrics(configFile string, prefix string) map[string]CollectorMetrics {
	cfg := getConfOrDie(configFile)
	m := make(map[string]CollectorMetrics)
	for collector, metricsMap := range cfg {
		cms := make(CollectorMetrics)
		for key, value := range metricsMap {
			fqName := value.Name
			if prefix != "" {
				fqName = fmt.Sprintf("%s_%s", prefix, fqName)
			}
			labels := value.Labels.keys()
			sort.Strings(labels)
			labelNames := []string{}
			for _, l := range labels {
				labelNames = append(labelNames, value.Labels[l])
			}
			cms[key] = Metric{
				Name:   value.Name,
				Desc:   prometheus.NewDesc(fqName, value.Help, labelNames, nil),
				Labels: labels,
				Push:   value.Push,
			}
		}
		m[collector] = cms
	}
	log.Debugf("collectorMetrics: %+v", m)
	return m
}

func WatchMetrics(configFile string, prefix string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("failed to create watcher, err %v", err)
		return
	}
	defer watcher.Close()
	err = watcher.Add(configFile)
	if err != nil {
		log.Infof("failed to add config to watcher, err %v", err)
		return
	}
	log.Infof("start watching metrics config %s", configFile)
	defer log.Infof("stop watching metrics config %s", configFile)
	for {
		select {
		case event := <-watcher.Events:
			if event.Name == configFile && event.Op&fsnotify.Remove == fsnotify.Remove {
				log.Infof("inotify: %v, reload config", event.String())
				watcher.Remove(event.Name)
				watcher.Add(configFile)
				m := getMetrics(configFile, prefix)
				for _, l := range listeners {
					l(m)
				}
			}
		case err := <-watcher.Errors:
			log.Errorf("watcher error: %v", err)
		}
	}
}
