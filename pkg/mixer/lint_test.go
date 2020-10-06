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
	"os"
	"testing"

	"github.com/google/go-jsonnet"
)

func TestLintPrometheusAlerts(t *testing.T) {
	const testAlerts = alerts + `+
{
  _config+:: {
     kubeStateMetricsSelector: 'job="ksm"',
  }
}`

	vm := jsonnet.MakeVM()

	f, err := ioutil.TempFile("", "alerts.jsonnet")
	if err != nil {
		t.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(testAlerts); err != nil {
		t.Errorf("failed to write alerts.jsonnet to disk: %v", err)
	}

	if err := f.Close(); err != nil {
		t.Errorf("failed to close temp file: %v", err)
	}

	errs := make(chan error)
	go lintPrometheus(f.Name(), vm, errs)
	for err := range errs {
		t.Errorf("linting wrote unexpected output: %v", err)
	}
}

func TestLintPrometheusRules(t *testing.T) {
	vm := jsonnet.MakeVM()

	f, err := ioutil.TempFile("", "rules.jsonnet")
	if err != nil {
		t.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(rules); err != nil {
		t.Errorf("failed to write rules.jsonnet to disk: %v", err)
	}

	if err := f.Close(); err != nil {
		t.Errorf("failed to close temp file: %v", err)
	}

	errs := make(chan error)
	go lintPrometheus(f.Name(), vm, errs)
	for err := range errs {
		t.Errorf("linting wrote unexpected output: %v", err)
	}
}
