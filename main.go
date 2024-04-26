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
	"net/http"
	"os"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/Cambricon/mlu-exporter/pkg/collector"
	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// Here version is a variable that stores the version of the application.
// The value of this variable will be set at build time with the -ldflags option.
// Using `-ldflags="-X 'main.Version=v1.0.0'"`.
var version string

type Options struct {
	Collector     []string `long:"collector" description:"enabled collectors" choice:"cndev" choice:"podresources" choice:"host" default:"cndev"`
	EnvShareNum   uint     `long:"env-share-num" description:"numbers of vfs under env share mode, 0 means env share disabled" default:"0"`
	Hostname      string   `long:"hostname" description:"machine hostname" env:"ENV_NODE_NAME"`
	LogLevel      string   `long:"log-level" description:"set log level: trace/debug/info/warn/error/fatal/panic" default:"info"`
	MetricsConfig string   `long:"metrics-config" description:"configuration file of MLU exporter metrics" default:"/etc/mlu-exporter/metrics.yaml"`
	MetricsPath   string   `long:"metrics-path" description:"metrics path of the exporter service" default:"/metrics"`
	MetricsPrefix string   `long:"metrics-prefix" description:"prefix of all metric names" env:"ENV_METRICS_PREFIX"`
	Port          uint     `long:"port" description:"exporter service port" default:"30108" env:"ENV_SERVE_PORT"`
	Version       bool     `long:"version" description:"print out version"`
	VirtualMode   string   `long:"virtual-mode" description:"virtual mode for devices" default:"" choice:"dynamic-smlu" choice:"env-share"`
}

func ParseFlags() Options {
	options := Options{}
	parser := flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}
	log.Infof("Options: %v\n", options)
	return options
}

func main() {
	options := ParseFlags()
	if options.Version {
		fmt.Println("Version:", version)
		return
	}

	switch options.LogLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	}

	c := collector.NewCollectors(
		options.Collector,
		options.EnvShareNum,
		options.Hostname,
		options.MetricsConfig,
		options.MetricsPrefix,
		options.VirtualMode,
	)
	r := prometheus.NewRegistry()
	r.MustRegister(c)

	http.Handle(options.MetricsPath, promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		cli := cndev.NewCndevClient()
		if err := cli.Init(true); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error: %v", errors.Wrap(err, "Init"))))
			log.Errorln(errors.Wrap(err, "Init"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}
	})

	server := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", options.Port),
	}

	log.Printf("start serving at %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
