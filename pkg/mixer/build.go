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
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/google/go-jsonnet"
)

type BuildConfig struct {
	JPath string
	YAML  bool
}

func Build(filename string, config BuildConfig) ([]byte, error) {
	vm := jsonnet.MakeVM()

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	vm.Importer(&jsonnet.FileImporter{
		JPaths: []string{config.JPath},
	})

	j, err := vm.EvaluateSnippet(filename, string(contents))
	if err != nil {
		return nil, err
	}
	output := []byte(j)

	if config.YAML {
		output, err = yaml.JSONToYAML(output)
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}

func BuildMulti(filename string, config BuildConfig) (map[string]string, error) {
	vm := jsonnet.MakeVM()

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	vm.Importer(&jsonnet.FileImporter{
		JPaths: []string{config.JPath},
	})

	files, err := vm.EvaluateSnippetMulti(filename, string(contents))
	if err != nil {
		return nil, err
	}

	if config.YAML {
		for filename, content := range files {
			y, err := yaml.JSONToYAML([]byte(content))
			if err != nil {
				return nil, err
			}
			files[filename] = string(y)
		}
	}

	return files, nil
}
