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
	filename, delete := writeTempFile(t, "alerts.jsonnet", testAlerts)
	defer delete()

	vm := jsonnet.MakeVM()
	errs := make(chan error)
	go lintPrometheus(filename, vm, errs)
	for err := range errs {
		t.Errorf("linting wrote unexpected output: %v", err)
	}
}

func TestLintPrometheusRules(t *testing.T) {
	filename, delete := writeTempFile(t, "rules.jsonnet", rules)
	defer delete()

	vm := jsonnet.MakeVM()
	errs := make(chan error)
	go lintPrometheus(filename, vm, errs)
	for err := range errs {
		t.Errorf("linting wrote unexpected output: %v", err)
	}
}

func TestLintGrafana(t *testing.T) {
	vm := jsonnet.MakeVM()
	errs := make(chan error)
	go lintGrafanaDashboards("lint_test_dashboard.json", vm, errs)
	for err := range errs {
		t.Errorf("linting wrote unexpected output: %v", err)
	}
}

func writeTempFile(t *testing.T, pattern string, contents string) (filename string, delete func()) {
	f, err := ioutil.TempFile("", pattern)
	if err != nil {
		t.Errorf("failed to create temp file: %v", err)
	}

	if _, err := f.WriteString(rules); err != nil {
		t.Errorf("failed to write rules.jsonnet to disk: %v", err)
	}

	if err := f.Close(); err != nil {
		t.Errorf("failed to close temp file: %v", err)
	}

	return f.Name(), func() { os.Remove(f.Name()) }
}
