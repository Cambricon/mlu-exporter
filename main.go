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

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cambricon/mlu-exporter/pkg/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newHandler(host string, cfg string, collectors []string, prefix string) http.Handler {
	c := collector.NewCollectors(host, cfg, collectors, prefix)
	r := prometheus.NewRegistry()
	r.MustRegister(c)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	return handler
}

func main() {
	options := ParseFlags()
	http.Handle(options.MetricsPath, newHandler(options.Hostname, options.MetricsConfig, options.Collector, options.MetricsPrefix))
	log.Printf("Start serving at http://0.0.0.0:%d%s", options.Port, options.MetricsPath)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", options.Port), nil))
}
