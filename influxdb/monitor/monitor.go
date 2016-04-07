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

package monitor

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/intelsdi-x/snap-plugin-collector-influxdb/influxdb/dtype"
	"github.com/intelsdi-x/snap-plugin-collector-influxdb/influxdb/parser"
)

// Monitoring is an interface represents data monitoring service
// (needed for mocking purposes)
type Monitoring interface {
	GetStatistics() (dtype.Results, error)
	GetDiagnostics() (dtype.Results, error)
	InitURLs(sets map[string]interface{}) error
}

// Monitor holds urls
type Monitor struct {
	urlStatistic  *url.URL
	urlDiagnostic *url.URL
}

// GetStatistics returns statistics information (url contains query "SHOW STATS")
func (m *Monitor) GetStatistics() (dtype.Results, error) {
	return getURLResults(m.urlStatistic.String())
}

// GetDiagnostics returns diagnostics information (url contains query "SHOW DIAGNOSTICS")
func (m *Monitor) GetDiagnostics() (dtype.Results, error) {
	return getURLResults(m.urlDiagnostic.String())
}

// InitURLs initializes URLs based on settings
func (m *Monitor) InitURLs(settings map[string]interface{}) error {
	var err1, err2 error

	m.urlStatistic, err1 = createURL(settings, "show stats")
	if err1 != nil {
		fmt.Fprintln(os.Stderr, "Cannot parse raw url into a URL structure with `show stats` query, err=", err1)
	}

	m.urlDiagnostic, err2 = createURL(settings, "show diagnostics")
	if err2 != nil {
		fmt.Fprintln(os.Stderr, "Cannot parse raw url into a URL structure with `show diagnostics` query, err=", err2)
	}

	if err1 != nil && err2 != nil {
		return errors.New("Invalid URL-encoding")
	}

	return nil
}

// createURL returns URL structure created base on `sets` (keeps info about hostname, port, etc.) and query state
func createURL(sets map[string]interface{}, query string) (*url.URL, error) {
	p, ok := sets["port"].(int)
	if !ok {
		return nil, fmt.Errorf("Invalid port value")
	}
	u, err := url.Parse(fmt.Sprintf("http://%s:%d/query?u=%s&p=%s&pretty=true",
		sets["host"].(string),
		p,
		sets["user"].(string),
		sets["password"].(string),
	))

	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("q", query)
	u.RawQuery = q.Encode()

	return u, nil
}

// getURLResults returns result of GET to the specified url which has been parsed into a dtype.Results structure
func getURLResults(url string) (dtype.Results, error) {
	response, err := getHTTPResponse(url)
	if err != nil {
		return nil, err
	}
	results, err := parser.FromJSON(response)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// getHTTPResponse returns HTTP response of GET to the specified url
func getHTTPResponse(url string) ([]byte, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
