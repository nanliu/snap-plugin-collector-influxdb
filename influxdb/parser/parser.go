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

package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/intelsdi-x/snap-plugin-collector-influxdb/influxdb/dtype"
)

// Parser holds data from unmarhalled json
type Parser struct {
	output dtype.Results
}

// To unmarshal JSON into a struct, structs have to contain exported fields
type Response struct {
	Results []ResultType `json:"results"`
}

type ResultType struct {
	Series []SeriesType `json:"series"`
}

type SeriesType struct {
	Name    string            `json:"name"`
	Tags    map[string]string `json:"tags"`
	Columns []string          `json:"columns"`
	Values  [][]interface{}   `json:"values"`
}

// FromJSON unmarshals `data` and parses results to structure defined in dtype package (dtype.Results)
func FromJSON(data []byte) (dtype.Results, error) {
	p := &Parser{output: dtype.Results{}}
	var response Response

	err := json.Unmarshal(data, &response)

	if err != nil {
		return nil, err
	}

	for _, result := range response.Results {
		for _, s := range result.Series {

			if err := p.addSeries(s); err != nil {
				return nil, err
			}

		}
	}

	return p.output, nil
}

// addSeries adds series defined by SeriesType structure to Parser `output` item
func (p *Parser) addSeries(st SeriesType) error {
	key := st.Name
	data := map[string]interface{}{}

	if id, ok := st.Tags["id"]; ok {
		// if tag 'id' exist, add its value to series name (key)
		key = key + "/" + id

	} else if _, ok := st.Tags["path"]; ok {
		// otherwise, treat the last path's element as an id
		path := strings.Split(st.Tags["path"], "/")
		key = key + "/" + path[len(path)-1]
	}

	if _, exist := p.output[key]; exist {
		return fmt.Errorf("Series %+s is not unique", st.Name)
	}

	// validation content of values
	if len(st.Values) < 1 {
		return fmt.Errorf("Series %+s has no 'values' section", st.Name)
	}

	values := st.Values[0]

	if len(values) != len(st.Columns) {
		return fmt.Errorf("For series %+s can not assign %d values to %d columns", st.Name, len(values), len(st.Columns))
	}

	for i, column := range st.Columns {
		data[column] = values[i]
	}

	// adding series to parser output
	p.output[key] = &dtype.Series{
		Data: data,
		Tags: st.Tags,
	}

	return nil
}
