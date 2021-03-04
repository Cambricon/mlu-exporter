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
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	MetricsConfig string   `long:"metrics-config" description:"configuration file of MLU exporter metrics" default:"/var/lib/mlu-exporter/metrics.yaml"`
	MetricsPath   string   `long:"metrics-path" description:"metrics path of the exporter service" default:"/metrics"`
	Hostname      string   `long:"hostname" description:"machine hostname" env:"ENV_NODE_NAME"`
	Port          uint     `long:"port" description:"exporter service port" default:"30108" env:"ENV_SERVE_PORT"`
	Collector     []string `long:"collector" description:"enabled collectors" choice:"cndev" choice:"podresources" choice:"host" choice:"cnpapi" default:"cndev"`
	MetricsPrefix string   `long:"metrics-prefix" description:"prefix of all metric names" env:"ENV_METRICS_PREFIX"`
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
	log.Printf("Options: %v\n", options)
	return options
}
