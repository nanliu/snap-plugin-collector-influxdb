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
	"testing"

	"github.com/intelsdi-x/snap-plugin-collector-influxdb/influxdb/dtype"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

type mcMock struct {
	mock.Mock
}

func (mc *mcMock) GetStatistics() (dtype.Results, error) {
	args := mc.Called()
	return args.Get(0).(dtype.Results), args.Error(1)
}

func (mc *mcMock) GetDiagnostics() (dtype.Results, error) {
	args := mc.Called()
	return args.Get(0).(dtype.Results), args.Error(1)
}

func (mc *mcMock) InitURLs(settings map[string]interface{}) error {
	args := mc.Called()
	return args.Error(0)
}

var mockStats = dtype.Results{

	"shard/1": &dtype.Series{
		Data: map[string]interface{}{
			"columnA": 1,
			"columnB": 10.1,
			"columnC": "value",
		},
		Tags: map[string]string{
			"tag1": "v1",
			"tag2": "v2",
		},
	},
	"httpd": &dtype.Series{
		Data: map[string]interface{}{
			"columnA": 1,
			"columnB": 10.1,
			"columnC": "value",
		},
		Tags: map[string]string{
			"tag1": "v1",
			"tag2": "v2",
		},
	},
}

var mockDiagn = dtype.Results{
	"build": &dtype.Series{
		Data: map[string]interface{}{
			"columnA": 1,
			"columnB": 10.1,
			"columnC": "value",
		},
	},
}

var mockMtsStat = []plugin.MetricType{
	plugin.MetricType{Namespace_: core.NewNamespace("intel", "influxdb", "stat", "shard", "1", "columnA")},
	plugin.MetricType{Namespace_: core.NewNamespace("intel", "influxdb", "stat", "shard", "1", "columnB")},
	plugin.MetricType{Namespace_: core.NewNamespace("intel", "influxdb", "stat", "shard", "1", "columnC")},
	plugin.MetricType{Namespace_: core.NewNamespace("intel", "influxdb", "stat", "httpd", "columnA")},
	plugin.MetricType{Namespace_: core.NewNamespace("intel", "influxdb", "stat", "httpd", "columnB")},
	plugin.MetricType{Namespace_: core.NewNamespace("intel", "influxdb", "stat", "httpd", "columnC")},
}

var mockMtsDiagn = []plugin.MetricType{
	plugin.MetricType{Namespace_: core.NewNamespace("intel", "influxdb", "diagn", "build", "columnA")},
	plugin.MetricType{Namespace_: core.NewNamespace("intel", "influxdb", "diagn", "build", "columnB")},
	plugin.MetricType{Namespace_: core.NewNamespace("intel", "influxdb", "diagn", "build", "columnC")},
}

var mockMtsWildCards = []plugin.MetricType{
	plugin.MetricType{Namespace_: core.NewNamespace("intel", "influxdb", "diagn", "*")},
	plugin.MetricType{Namespace_: core.NewNamespace("intel", "influxdb", "stat", "*")},
}

// Mts is a mocked metrics, both statistical and diagnostic
var mockMts = append(mockMtsStat, mockMtsDiagn...)

func TestGetConfigPolicy(t *testing.T) {
	influxdbPlugin := New()

	Convey("getting config policy", t, func() {
		So(func() { influxdbPlugin.GetConfigPolicy() }, ShouldNotPanic)
		_, err := influxdbPlugin.GetConfigPolicy()
		So(err, ShouldBeNil)
	})
}

func TestGetMetricTypes(t *testing.T) {

	Convey("initialization fails", t, func() {

		Convey("when no config items available", func() {
			influxdbPlugin := New()
			cfg := plugin.NewPluginConfigType()

			So(func() { influxdbPlugin.GetMetricTypes(cfg) }, ShouldPanic)
		})

		Convey("when one of config item is not available", func() {
			influxdbPlugin := New()
			cfg := plugin.NewPluginConfigType()
			cfg.DeleteItem("user")

			So(func() { influxdbPlugin.GetMetricTypes(cfg) }, ShouldPanic)
		})

		Convey("when config item has different type than expected", func() {
			influxdbPlugin := New()
			cfg := plugin.NewPluginConfigType()
			cfg.DeleteItem("port")
			cfg.AddItem("port", ctypes.ConfigValueStr{Value: "1234"})

			So(func() { influxdbPlugin.GetMetricTypes(cfg) }, ShouldPanic)
		})

		Convey("when initialization of URLs returns error", func() {
			mc := &mcMock{}
			influxdbPlugin := New()
			influxdbPlugin.service = mc
			cfg := plugin.NewPluginConfigType()

			mc.On("InitURLs").Return(errors.New("x"))

			So(func() { influxdbPlugin.GetMetricTypes(cfg) }, ShouldPanic)
		})

	})

	Convey("Metrics are not available", t, func() {

		Convey("when cannot obtain any data", func() {
			mc := &mcMock{}
			influxdbPlugin := New()
			influxdbPlugin.service = mc
			cfg := getMockPluginConfig()

			mc.On("InitURLs").Return(nil)
			mc.On("GetStatistics").Return(dtype.Results{}, errors.New("x"))
			mc.On("GetDiagnostics").Return(dtype.Results{}, errors.New("x"))

			So(func() { influxdbPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := influxdbPlugin.GetMetricTypes(cfg)

			So(results, ShouldBeEmpty)
			So(err, ShouldBeNil)
		})

		Convey("when cannot obtain statictics data", func() {
			mc := &mcMock{}
			influxdbPlugin := New()
			influxdbPlugin.service = mc
			cfg := getMockPluginConfig()

			mc.On("InitURLs").Return(nil)
			mc.On("GetStatistics").Return(dtype.Results{}, errors.New("x"))
			mc.On("GetDiagnostics").Return(mockDiagn, nil)

			So(func() { influxdbPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := influxdbPlugin.GetMetricTypes(cfg)

			So(results, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
		})

		Convey("when cannot obtain diagnostics data", func() {
			mc := &mcMock{}
			influxdbPlugin := New()
			influxdbPlugin.service = mc
			cfg := getMockPluginConfig()

			mc.On("InitURLs").Return(nil)
			mc.On("GetStatistics").Return(mockStats, nil)
			mc.On("GetDiagnostics").Return(dtype.Results{}, errors.New("x"))

			So(func() { influxdbPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)

			So(func() { influxdbPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
			results, err := influxdbPlugin.GetMetricTypes(cfg)

			So(results, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
		})
	})

	Convey("successfull getting metrics types", t, func() {
		mc := &mcMock{}
		influxdbPlugin := New()
		influxdbPlugin.service = mc
		cfg := getMockPluginConfig()

		mc.On("InitURLs").Return(nil)
		mc.On("GetStatistics").Return(mockStats, nil)
		mc.On("GetDiagnostics").Return(mockDiagn, nil)

		So(func() { influxdbPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)

		So(func() { influxdbPlugin.GetMetricTypes(cfg) }, ShouldNotPanic)
		results, err := influxdbPlugin.GetMetricTypes(cfg)

		So(results, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})

}

func TestCollectMetrics(t *testing.T) {

	mts := getMockMetricsConfigured()

	Convey("initialization fails", t, func() {
		mc := &mcMock{}
		influxdbPlugin := New()
		influxdbPlugin.service = mc

		mc.On("InitURLs").Return(errors.New("x"))

		So(func() { influxdbPlugin.CollectMetrics(mts) }, ShouldPanic)
	})

	Convey("metrics are not available", t, func() {
		mc := &mcMock{}
		influxdbPlugin := New()
		influxdbPlugin.service = mc

		mc.On("InitURLs").Return(nil)

		Convey("when statistics data are not available", func() {
			mc.On("GetStatistics").Return(dtype.Results{}, errors.New("x"))
			mc.On("GetDiagnostics").Return(mockDiagn, nil)

			So(func() { influxdbPlugin.CollectMetrics(mts) }, ShouldNotPanic)

			results, err := influxdbPlugin.CollectMetrics(mts)

			So(len(results), ShouldEqual, len(mockMtsDiagn))
			So(err, ShouldBeNil)
		})

		Convey("when diagnostics data are not available", func() {
			mc.On("GetStatistics").Return(mockStats, nil)
			mc.On("GetDiagnostics").Return(dtype.Results{}, nil)

			So(func() { influxdbPlugin.CollectMetrics(mts) }, ShouldNotPanic)

			results, err := influxdbPlugin.CollectMetrics(mts)

			So(len(results), ShouldEqual, len(mockMtsStat))
			So(err, ShouldBeNil)
		})

	})

	Convey("invalid metrics namespaces", t, func() {
		mc := &mcMock{}
		influxdbPlugin := New()
		influxdbPlugin.service = mc
		cfg := getMockMetricConfig()

		mc.On("InitURLs").Return(nil)
		mc.On("GetStatistics").Return(mockStats, nil)
		mc.On("GetDiagnostics").Return(mockDiagn, nil)

		Convey("when series type is unknown", func() {
			mtsInvalid := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace("intel", "influxdb", "unknown", "shard", "1", "columnA"),
					Config_:    cfg,
				},
			}

			So(func() { influxdbPlugin.CollectMetrics(mtsInvalid) }, ShouldNotPanic)

			results, err := influxdbPlugin.CollectMetrics(mtsInvalid)

			So(results, ShouldBeEmpty)
			So(err, ShouldBeNil)
		})

		Convey("when series name is empty", func() {
			mtsInvalid := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace("intel", "influxdb", "stats", ""),
					Config_:    cfg,
				},
			}
			So(func() { influxdbPlugin.CollectMetrics(mtsInvalid) }, ShouldNotPanic)

			results, err := influxdbPlugin.CollectMetrics(mtsInvalid)

			So(results, ShouldBeEmpty)
			So(err, ShouldBeNil)
		})

	})

	Convey("successfull collect metrics", t, func() {

		mc := &mcMock{}
		influxdbPlugin := New()
		influxdbPlugin.service = mc

		mc.On("InitURLs").Return(nil)
		mc.On("GetStatistics").Return(mockStats, nil)
		mc.On("GetDiagnostics").Return(mockDiagn, nil)

		Convey("when metrics namespaces are deliver directly", func() {

			So(func() { influxdbPlugin.CollectMetrics(mts) }, ShouldNotPanic)

			results, err := influxdbPlugin.CollectMetrics(mts)

			So(results, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
		})
	})

}

func getMockPluginConfig() plugin.ConfigType {
	// mocking global config
	cfg := plugin.NewPluginConfigType()
	cfg.AddItem("host", ctypes.ConfigValueStr{Value: "hostname"})
	cfg.AddItem("port", ctypes.ConfigValueInt{Value: 1234})
	cfg.AddItem("user", ctypes.ConfigValueStr{Value: "test"})
	cfg.AddItem("password", ctypes.ConfigValueStr{Value: "passwd"})

	return cfg
}

func getMockMetricConfig() *cdata.ConfigDataNode {
	// mocking metric config
	cfg := cdata.NewNode()
	cfg.AddItem("host", ctypes.ConfigValueStr{Value: "hostname"})
	cfg.AddItem("port", ctypes.ConfigValueInt{Value: 1234})
	cfg.AddItem("user", ctypes.ConfigValueStr{Value: "test"})
	cfg.AddItem("password", ctypes.ConfigValueStr{Value: "passwd"})

	return cfg
}

func getMockMetricsConfigured() []plugin.MetricType {
	mts := mockMts
	cfg := getMockMetricConfig()

	// add mocked config to each metric
	for i := range mts {
		mts[i].Config_ = cfg
	}

	return mts
}
