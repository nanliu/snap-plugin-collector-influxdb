// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt
Copyright 2016 Intel Corporation
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package influxdb

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"

	"github.com/intelsdi-x/snap-plugin-collector-influxdb/influxdb/dtype"
	"github.com/intelsdi-x/snap-plugin-collector-influxdb/influxdb/monitor"
	"github.com/intelsdi-x/snap-plugin-utilities/config"
)

const (
	// Name of plugin
	Name = "influxdb"
	// Version of plugin
	Version = 2
	// Type of plugin
	Type = plugin.CollectorPluginType

	nsVendor = "intel"
	nsClass  = "influxdb"

	nsTypeStats = "stat"
	nsTypeDiagn = "diagn"
)

const (
	typeUnknown = iota
	typeStats
	typeDiagn
)

// prefix in metric namespace
var prefix = []string{nsVendor, nsClass}

// InfluxdbCollector holds data retrived from influxDB system monitoring
type InfluxdbCollector struct {
	data        map[string]datum
	service     monitor.Monitoring
	initialized bool
}

type datum struct {
	value interface{}
	tags  map[string]string
}

// New returns new instance of snap-plugin-collector-influxdb
func New() *InfluxdbCollector {
	return &InfluxdbCollector{initialized: false, service: &monitor.Monitor{}, data: map[string]datum{}}
}

// GetConfigPolicy returns a ConfigPolicy
func (ic *InfluxdbCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}

// GetMetricTypes returns list of metrics based on influxDB system monitoring
func (ic *InfluxdbCollector) GetMetricTypes(cfg plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	mts := []plugin.PluginMetricType{}

	ic.init(cfg)
	ic.getStatistics()  // get statistical information about influxDB
	ic.getDiagnostics() // get diagnostic information about influxDB

	for ns, dat := range ic.data {
		mts = append(mts, plugin.PluginMetricType{Namespace_: splitNamespace(ns), Tags_: dat.tags})
	}

	return mts, nil
}

// CollectMetrics collects given metrics
func (ic *InfluxdbCollector) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	metrics := []plugin.PluginMetricType{}
	hostname, _ := os.Hostname()

	if !ic.initialized {
		ic.init(mts[0])     // if CollectMetrics() is called, mts has one item at least
		ic.getDiagnostics() // get diagnostic information (once only)
	}

	ic.getStatistics() // get statistical information

	for _, m := range mts {
		if dat, ok := ic.data[joinNamespace(m.Namespace())]; ok {
			metric := plugin.PluginMetricType{
				Namespace_: m.Namespace(),
				Data_:      dat.value,
				Source_:    hostname,
				Timestamp_: time.Now(),
				Tags_:      dat.tags,
			}

			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

// init initializes InfluxdbCollector instance based on config `cfg`
func (ic *InfluxdbCollector) init(cfg interface{}) {
	items := []string{"host", "port", "user", "password"}
	settings, err := config.GetConfigItems(cfg, items)
	handleErr(err)

	err = ic.service.InitURLs(settings)
	handleErr(err)

	ic.initialized = true
}

// getDiagnostics executes the command "SHOW DIAGNOSTICS" (indirectly)
func (ic *InfluxdbCollector) getDiagnostics() {
	err := ic.getData(typeDiagn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot get influxdb diagnostic information, err=", err)
	}
}

// getStatistics executes the command "SHOW STATS" (indirectly)
func (ic *InfluxdbCollector) getStatistics() {

	err := ic.getData(typeStats)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot get influxdb internal statistics, err=", err)
	}
}

// getData executes a command specified by given `type` of desired data
// and assignes its results to InfluxdbCollector structure item `data`
func (ic *InfluxdbCollector) getData(kind int) error {
	var results dtype.Results
	var err error
	var nsType string

	switch kind {
	case typeStats:
		nsType = nsTypeStats
		results, err = ic.service.GetStatistics()

	case typeDiagn:
		nsType = nsTypeDiagn
		results, err = ic.service.GetDiagnostics()

	default:
		err = errors.New("Inalid type of monitoring service")
	}

	if err != nil {
		return err
	}

	for seriesName, series := range results {
		for columnName := range series.Data {
			key := createNamespace(nsType, seriesName, columnName)
			ic.data[key] = datum{
				value: series.Data[columnName],
				tags:  series.Tags,
			}

		}
	}

	return nil
}

// handleErr handles critical error indicated with an abnormal state of plugin
func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

// createNamespace returns namespace composed from prefix, type of metric and metric name's components (the elements are joined to a single string)
func createNamespace(nsType, seriesName, columnName string) string {
	ns := append(prefix, nsType)
	ns = append(ns, seriesName)
	ns = append(ns, columnName)
	return strings.Join(ns, "/")

}

// joinNamespace concatenates the elements of `ns` to create a single string combined by slash
func joinNamespace(ns []string) string {
	return strings.Join(ns, "/")
}

//splitNamespace splits namespace (repesented by single string `s`) and returns  a slice of the substrings beetween slash separator
func splitNamespace(ns string) []string {
	return strings.Split(ns, "/")
}
