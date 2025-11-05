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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Cambricon/mlu-exporter/pkg/cndev"
	"github.com/Cambricon/mlu-exporter/pkg/collector"
	"github.com/Cambricon/mlu-exporter/pkg/metrics"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/prometheus/common/expfmt"
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
	HostIP        string   `long:"host-ip" description:"machine host ip" env:"ENV_NODE_IP"`
	LogLevel      string   `long:"log-level" description:"set log level: trace/debug/info/warn/error/fatal/panic" default:"info"`
	MetricsConfig string   `long:"metrics-config" description:"configuration file of MLU exporter metrics" default:"/etc/mlu-exporter/metrics.yaml"`
	MetricsPath   string   `long:"metrics-path" description:"metrics path of the exporter service" default:"/metrics"`
	MetricsPrefix string   `long:"metrics-prefix" description:"prefix of all metric names" env:"ENV_METRICS_PREFIX"`
	Port          uint     `long:"port" description:"exporter service port" default:"30108" env:"ENV_SERVE_PORT"`
	Version       bool     `long:"version" description:"print out version"`
	VirtualMode   string   `long:"virtual-mode" description:"virtual mode for devices" default:"" choice:"dynamic-smlu" choice:"env-share"`

	PushGatewayURL string `long:"push-gateway-url" description:"If set, metrics with push enabled will push to this server via prometheus push gateway protocol" env:"PUSH_GATEWAY_URL"`
	PushIntervalMS uint   `long:"push-interval-ms" description:"numbers of metrics push interval in milliseconds, minimum 100" default:"500" env:"PUSH_INTERVAL_MS"`
	PushJobName    string `long:"push-job-name" description:"metrics push job name" default:"mlu-push-monitoring" env:"PUSH_JOB_NAME"`
	ClusterName    string `long:"cluster-name" description:"cluster name, add cluster label for metrics push" env:"CLUSTER_NAME"`
	PushCAFile     string `long:"push-ca-file" description:"Optional,CA certificate for pushing data" env:"PUSH_CA_FILE"`
	PushTLSFile    string `long:"push-tls-file" description:"Optional,TLS certificate for pushing data" env:"PUSH_TLS_FILE"`
	PushKeyFile    string `long:"push-key-file" description:"Optional,Key certificate for pushing data" env:"PUSH_KEY_FILE"`

	XIDErrorMetricName        string `long:"xid-error-metric-name" description:"xid error metric name in config, if not set, not push this metric data" env:"XID_ERROR_METRIC_NAME"`
	XIDErrorRetryTimes        int    `long:"xid-error-retry-times" description:"retry times when push xid error metric failed" default:"10" env:"XID_ERROR_RETRY_TIMES"`
	LogFileForXIDMetricFailed string `long:"log-file-for-xid-metric-failed" description:"when report xid metric failed, write log to this file" env:"LOG_FILE_FOR_XID_METRIC_FAILED"`
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
	log.Printf("Final options: %v", options)

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

	log.Info("Start loading cndev")
	cndevcli := cndev.NewCndevClient()
	if err := collector.EnsureCndevLib(); err != nil {
		log.Panicf("Failed to ensure CNDEV lib %v", err)
	}

	mluInfo := &collector.MLUStatMap{}
	collector.EnsureMLUAllOK(cndevcli, mluInfo, true)
	if mluInfo.InProblem.Load() {
		log.Warn("MLU is in problem state")
		go collector.EnsureMLUAllOK(cndevcli, mluInfo, false)
	}
	defer func() { log.Println("Shutdown of CNDEV returned:", cndevcli.Release()) }()

	metricConfig := metrics.GetMetrics(options.MetricsConfig, options.MetricsPrefix)
	log.Debug("Start WatchMetrics")
	go metrics.WatchMetrics(options.MetricsConfig, options.MetricsPrefix)

	if options.PushGatewayURL != "" {
		client := &http.Client{}
		var err error
		if options.PushCAFile != "" {
			client, err = buildClient(options)
			if err != nil {
				log.Panicf("Build push client %v", err)
			}
		}
		startPushMode(client, options, metricConfig, mluInfo)
		startCallbackMode(client, options, metricConfig, mluInfo)
	}

	c := collector.NewCollectors(
		options.Collector,
		metricConfig,
		options.EnvShareNum,
		options.Hostname,
		options.HostIP,
		options.VirtualMode,
		mluInfo,
		false,
	)
	log.Debug("Start RegisterWatcher")
	metrics.RegisterWatcher(c.UpdateMetrics)
	r := prometheus.NewRegistry()
	r.MustRegister(c)

	http.Handle(options.MetricsPath, promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
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

	http.HandleFunc("/logLevel", func(w http.ResponseWriter, r *http.Request) {
		level := r.URL.Query().Get("level")
		logLevel, err := log.ParseLevel(level)
		if err != nil {
			fmt.Fprintf(w, "Invalid log level: %s", level)
			return
		}
		log.SetLevel(logLevel)
		log.Printf("Log level set to %s", level)
	})

	server := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", options.Port),
	}

	log.Printf("Start serving at %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func startPushMode(client *http.Client, options Options, metricConfig map[string]metrics.CollectorMetrics, mluInfo *collector.MLUStatMap) {
	if options.PushIntervalMS < 100 {
		log.Fatal("Minimum of push-interval-ms is 100")
	}

	pushc := collector.NewCollectors(
		options.Collector,
		metricConfig,
		options.EnvShareNum,
		options.Hostname,
		options.HostIP,
		options.VirtualMode,
		mluInfo,
		true,
	)

	metrics.RegisterWatcher(pushc.UpdateMetrics)
	pushr := prometheus.NewRegistry()
	pushr.MustRegister(pushc)
	pusher := push.New(options.PushGatewayURL, options.PushJobName).
		Client(client).
		Format(expfmt.FmtText).
		Collector(pushr)

	if options.ClusterName != "" {
		pusher = pusher.Grouping("cluster", options.ClusterName)
	}

	ticker := time.NewTicker(time.Duration(options.PushIntervalMS) * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if err := pusher.Push(); err != nil {
				log.Errorln("Could not push to Pushgateway:", err)
			} else {
				log.Debugln("Metrics pushed to Pushgateway")
			}
		}
	}()
}

func startCallbackMode(client *http.Client, options Options, metricConfig map[string]metrics.CollectorMetrics, mluInfo *collector.MLUStatMap) {
	cb, err := collector.NewCallback(
		client,
		options.PushGatewayURL,
		metricConfig,
		mluInfo,
		options.XIDErrorMetricName,
		options.MetricsPrefix,
		options.Hostname,
		options.HostIP,
		options.LogFileForXIDMetricFailed,
		options.XIDErrorRetryTimes,
		options.PushJobName,
	)
	if err != nil {
		log.Debugln("XID Callback not enabled")
		return
	}
	if cb != nil {
		go cb.Start()
	}
}

func buildClient(options Options) (*http.Client, error) {
	caCert, err := os.ReadFile(options.PushCAFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		RootCAs:    caCertPool,
		MinVersion: tls.VersionTLS12,
	}

	if options.PushTLSFile != "" && options.PushKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(options.PushTLSFile, options.PushKeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}, nil
}
