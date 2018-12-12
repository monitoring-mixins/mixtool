// Copyright 2018 mixtool authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mixer

import (
	"encoding/json"

	"github.com/ghodss/yaml"
	"github.com/google/go-jsonnet"
	"github.com/pkg/errors"
)

type GenerateOptions struct {
	Directory string
	JPaths    []string
	YAML      bool
}

func GenerateAlerts(filename string, opts GenerateOptions) ([]byte, error) {
	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{
		JPaths: opts.JPaths,
	})

	j, err := evaluatePrometheusAlerts(vm, filename)
	if err != nil {
		return nil, err
	}

	output := []byte(j)

	if opts.YAML {
		output, err = yaml.JSONToYAML(output)
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}

func GenerateDashboards(filename string, opts GenerateOptions) (map[string]json.RawMessage, error) {
	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{
		JPaths: opts.JPaths,
	})

	j, err := evaluateGrafanaDashboards(vm, filename)
	if err != nil {
		return nil, err
	}

	var dashboards map[string]json.RawMessage
	if err := json.Unmarshal([]byte(j), &dashboards); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal dashboards")
	}

	return dashboards, nil
}
