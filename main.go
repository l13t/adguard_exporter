// Copyright (C) 2016 Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"adguard_exporter/adguard"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const (
	banner    = "adguard_exporter - %s\n"
	namespace = "adguard"
)

var (
	debug             bool
	version           bool
	listenAddress     string
	metricsPath       string
	endpoint          string
	logLevel          string
	logFormat         string
	avgProcessingTime = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "avg_processing_time"),
		"Average processing time.",
		nil, nil,
	)
	dnsQueries = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "dns_queries"),
		"Total number of DNS queries.",
		nil, nil,
	)
	blockedFiltering = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "blocked_filtering"),
		"Domains blocked.",
		nil, nil,
	)
	replacedParental = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "replaced_parental"),
		"Parental Control blocked.",
		nil, nil,
	)
	replacedSafebrowsing = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "replaced_safebrowsing"),
		"Safebrowsing Control blocked.",
		nil, nil,
	)
	replacedSafesearch = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "replaced_safesearch"),
		"Safesearch Control blocked.",
		nil, nil,
	)
)

// Exporter collects Adguard stats from the given server and exports them using
// the prometheus metrics package.
type Exporter struct {
	Adguard *adguard.Client
}

// NewExporter returns an initialized Exporter.
func NewExporter(endpoint string) (*Exporter, error) {
	log.Infoln("Setup Adguard exporter using URL: ", endpoint)
	adguard, err := adguard.NewClient(endpoint)
	if err != nil {
		return nil, err
	}
	return &Exporter{
		Adguard: adguard,
	}, nil
}

// Describe describes all the metrics ever exported by the Adguard exporter.
// It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- avgProcessingTime
	ch <- dnsQueries
	ch <- blockedFiltering
	ch <- replacedParental
	ch <- replacedSafebrowsing
	ch <- replacedSafesearch
}

// Collect the stats from channel and delivers them as Prometheus metrics.
// It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	resp, err := e.Adguard.GetMetrics()
	if err != nil {
		log.Errorf("Adguard error: %s", err.Error())
		return
	}
	log.Debugf("Adguard metrics: %#v", resp)
	ch <- prometheus.MustNewConstMetric(
		avgProcessingTime, prometheus.CounterValue, float64(resp.AvgProcessingTime))
	ch <- prometheus.MustNewConstMetric(
		dnsQueries, prometheus.CounterValue, float64(resp.DnsQueries))
	ch <- prometheus.MustNewConstMetric(
		blockedFiltering, prometheus.CounterValue, float64(resp.BlockedFiltering))
	ch <- prometheus.MustNewConstMetric(
		replacedParental, prometheus.CounterValue, float64(resp.ReplacedParental))
	ch <- prometheus.MustNewConstMetric(
		replacedSafebrowsing, prometheus.CounterValue, float64(resp.ReplacedSafebrowsing))
	ch <- prometheus.MustNewConstMetric(
		replacedSafesearch, prometheus.CounterValue, float64(resp.ReplacedSafesearch))
}

func init() {
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.StringVar(&listenAddress, "web.listen-address", ":9311", "Address to listen on for web interface and telemetry.")
	flag.StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	flag.StringVar(&endpoint, "adguard", "", "Endpoint of Adguard")
	flag.StringVar(&logLevel, "log.level", "info", "Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]")
	flag.StringVar(&logFormat, "log.format", "logger:stderr", `Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"`)
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(banner, version))
		flag.PrintDefaults()
	}

	flag.Parse()

	if version {
		fmt.Printf("%s", adguard.Version)
		os.Exit(0)
	}
	if logLevel != "" {
		if err := log.Base().SetLevel(logLevel); err != nil {
			log.Errorf("Failed to set log level: %s", err.Error())
			os.Exit(1)
		}

	}
	if logFormat != "" {
		if err := log.Base().SetFormat(logFormat); err != nil {
			log.Errorf("Failed to set log format: %s", err.Error())
			os.Exit(1)
		}
	}

	if len(endpoint) == 0 {
		usageAndExit("Adguard endpoint cannot be empty.", 1)
	}
}

func main() {
	exporter, err := NewExporter(endpoint)
	if err != nil {
		log.Errorf("Can't create exporter : %s", err)
		os.Exit(1)
	}
	log.Infoln("Register exporter")
	prometheus.MustRegister(exporter)

	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Adguard Exporter</title></head>
             <body>
             <h1>Adguard Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	log.Infoln("Listening on", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Usage()
	os.Exit(exitCode)
}

func storeMetric(ch chan<- prometheus.Metric, value string, desc *prometheus.Desc, labels ...string) {
	if val, err := strconv.ParseFloat(value, 64); err == nil {
		ch <- prometheus.MustNewConstMetric(
			desc, prometheus.GaugeValue, val, labels...)
	} else {
		log.Errorf("Can't store metric %s into %s: %s", value, desc, err.Error())
	}
}
