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
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const runbookFile = `
{
  prometheusAlerts+:: {
    groups+: [
      {
        name: 'kubernetes-resources',
        rules: [
          {
            alert: 'KubeNodeNotReady',
            expr: |||
              kube_node_status_condition{%(kubeStateMetricsSelector)s,condition="Ready",status="true"} == 0
            ||| % $._config,
            labels: {
              severity: 'warning',
            },
            annotations: {
              message: 'Overcommited CPU resource requests on Pods, cannot tolerate node failure.',
            },
            'for': '1h',
            runbook: |||
              This is an awesome runbook text. :tada:
            |||,
          },
        ],
      },
    ],
  }, 
} + {
  _config+:: {
    kubeStateMetricsSelector: 'job="ksm"',
  },
}
`

const runbookExpected = `# Alert Runbook

### kubernetes-resources

##### KubeNodeNotReady
+ *Severity*: warning
+ *Message*: ` + "`" + `Overcommited CPU resource requests on Pods, cannot tolerate node failure.` + "`" + `

This is an awesome runbook text. :tada:


`

func TestRunbook(t *testing.T) {
	b := &bytes.Buffer{}
	w := bufio.NewWriter(b)

	f, err := ioutil.TempFile("", "alerts.jsonnet")
	if err != nil {
		t.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(runbookFile); err != nil {
		t.Errorf("failed to write alerts.jsonnet to disk: %v", err)
	}

	if err := f.Close(); err != nil {
		t.Errorf("failed to close temp file: %v", err)
	}

	err = Runbook(w, f.Name(), RunbookOptions{})
	if err != nil {
		t.Errorf("failed to generate runbook: %v", err)
	}

	w.Flush()

	if strings.TrimSpace(b.String()) != strings.TrimSpace(runbookExpected) {
		t.Errorf("failed to generate correct runbook")
		t.Logf("expected:\n%s\n", runbookExpected)
		t.Logf("actual:\n%s\n", b.String())
	}
}
