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

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleMixins = `
{
    "mixins": [
        {
            "name": "ceph",
            "source": "https://github.com/ceph/ceph-mixins",
            "subdir": "",
            "description": "A set of Prometheus alerts for Ceph.\n\nThe scope of this project is to provide Ceph specific Prometheus rule files using Prometheus Mixins.\n"
        },
        {
            "name": "cortex",
            "source": "https://github.com/grafana/cortex-jsonnet",
            "subdir": "cortex-mixin"
        },
        {
            "name": "cool-mixin",
            "source": "https://github.com",
			"subdir": "cool-mixin",
			"description": "A fantastic mixin"
        }
    ]
}
`

func TestList(t *testing.T) {
	tempFile, err := os.CreateTemp("", "exampleMixinsTest.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	err = os.WriteFile(tempFile.Name(), []byte(exampleMixins), 0644)
	assert.NoError(t, err)

	body, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	mixins, err := parseMixinJSON([]byte(body))
	assert.NoError(t, err)

	exampleMixins := map[string]bool{"ceph": true, "cool-mixin": true, "cortex": true}
	for _, m := range mixins {
		if _, ok := exampleMixins[m.Name]; !ok {
			t.Errorf("failed to find %v in exampleMixinsTest", m.Name)
		}
	}
}
