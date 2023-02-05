// Copyright 2020 mixtool authors
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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testAlertJsonnet = `
{
	prometheusAlerts+:: {
		groups+: [
		  {
			name: 'test-alerts',
			rules: [
			  {
				alert: 'TestAlert',
				expr: |||
				   test_alert == 1
				|||,
				labels: {
				  severity: 'warning',
				},
				annotations: {
				  message: 'test alert',
				},
				'for': '5m',
			  },
			],
		  },
		],
	  },
}
`

const expectedYaml = `groups:
- name: test-alerts
  rules:
  - alert: TestAlert
    annotations:
      message: test alert
    expr: |
      test_alert == 1
    for: 5m
    labels:
      severity: warning
`

func TestEvalAlerts(t *testing.T) {
	inFile, err := ioutil.TempFile(os.TempDir(), "mixtool-")
	assert.NoError(t, err)
	defer os.Remove(inFile.Name())

	_, err = inFile.Write([]byte(testAlertJsonnet))
	assert.NoError(t, err)

	out, err := GenerateAlerts(inFile.Name(), GenerateOptions{
		YAML: true,
	})
	assert.NoError(t, err)
	assert.YAMLEq(t, expectedYaml, string(out))
}
