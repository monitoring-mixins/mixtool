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
	"fmt"
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

var alertTests = []struct {
	alert           string // input
	expectedLintErr string // expected lint error
}{
	// valid alert
	{
		`{
			alert: 'TestAlert',
			expr: 'up == 0',
			labels: {
				severity: 'warning',
			},
			annotations: {
				description: '{{ $labels.instance }} has been unready for more than 15 minutes.',
				summary: 'Instance is not ready.',
			},
			'for': '15m',
		}`,
		``,
	},
	{
		`{
			alert: 'SNMPDown',
			expr: 'up == 0',
			labels: {
				severity: 'warning',
			},
			annotations: {
				description: '{{ $labels.instance }} has been unready for more than 15 minutes.',
				summary: 'SNMP is down.',
			},
			'for': '15m',
		}`,
		``,
	},
	// alertnames
	{
		`{
			alert: 'testAlert',
			expr: 'up == 0',
			labels: {
				severity: 'warning',
			},
			annotations: {
				description: '{{ $labels.instance }} has been unready for more than 15 minutes.',
				summary: 'Instance is not ready.',
			},
			'for': '15m',
		}`,
		`[alert-name-camelcase] Alert 'testAlert' name is not in camel case`,
	},
	{
		`{
			alert: 'test_Alert',
			expr: 'up == 0',
			labels: {
				severity: 'warning',
			},
			annotations: {
				description: '{{ $labels.instance }} has been unready for more than 15 minutes.',
				summary: 'Instance is not ready.',
			},
			'for': '15m',
		}`,
		`[alert-name-camelcase] Alert 'test_Alert' name is not in camel case`,
	},
	{
		`{
			alert: 'test Alert',
			expr: 'up == 0',
			labels: {
				severity: 'warning',
			},
			annotations: {
				description: '{{ $labels.instance }} has been unready for more than 15 minutes.',
				summary: 'Instance is not ready.',
			},
			'for': '15m',
		}`,
		`[alert-name-camelcase] Alert 'test Alert' name is not in camel case`,
	},

	// severity
	{
		`{
			alert: 'TestAlert',
			expr: 'up == 0',
			labels: {
				severity: 'disaster',
			},
			annotations: {
				description: '{{ $labels.instance }} has been unready for more than 15 minutes.',
				summary: 'Instance is not ready.',
			},
			'for': '15m',
		}`,
		`[alert-severity-rule] Alert 'TestAlert' severity must be 'warning', 'critical' or 'info', is currently 'disaster'`,
	},
	{
		`{
			alert: 'TestAlert',
			expr: 'up == 0',
			labels: {
			},
			annotations: {
				description: '{{ $labels.instance }} has been unready for more than 15 minutes.',
				summary: 'Instance is not ready.',
			},
			'for': '15m',
		}`,
		`[alert-severity-rule] Alert 'TestAlert' severity must be 'warning', 'critical' or 'info', is currently ''`,
	},
	// summary
	{
		`{
			alert: 'TestAlert',
			expr: 'up == 0',
			labels: {
				severity: 'warning',
			},
			annotations: {
				summary: 'Instance is not ready.',
			},
			'for': '15m',
		}`,
		`[alert-description-missing-rule] Alert 'TestAlert' must have annotation 'description'`,
	},
	{
		`{
			alert: 'TestAlert',
			expr: 'up == 0',
			labels: {
				severity: 'warning',
			},
			annotations: {
				summary: 'Instance {{ $labels.instance}} is not ready.',
				description: '{{ $labels.instance }} has been unready for more than 15 minutes.',
			},
			'for': '15m',
		}`,
		`[alert-summary-templating] Alert TestAlert annotation 'summary' must not use templates`,
	},
	{
		`{
			alert: 'TestAlert',
			expr: 'up == 0',
			labels: {
				severity: 'warning',
			},
			annotations: {
				summary: 'Instance is not ready',
				description: '{{ $labels.instance }} has been unready for more than 15 minutes.',
			},
			'for': '15m',
		}`,
		`[alert-summary-style] Alert TestAlert annotation 'summary' must start with capital letter and end with period, is currently 'Instance is not ready'`,
	},
	// description
	{
		`{
			alert: 'TestAlert',
			expr: 'up == 0',
			labels: {
				severity: 'warning',
			},
			annotations: {
				summary: 'Instance is not ready.',
				description: '{{ $labels.instance }} has been unready for more than 15 minutes.',
			},
			'for': '15m',
		}`,
		`[alert-summary-missing-rule] Alert 'TestAlert' must have annotation 'summary'`,
	},
	{
		`{
			alert: 'TestAlert',
			expr: 'up == 0',
			labels: {
				severity: 'warning',
			},
			annotations: {
				summary: 'Instance is not ready.',
				description: 'Instance has been unready for more than 15 minutes.',
			},
			'for': '15m',
		}`,
		`[alert-description-templating] Alert TestAlert annotation 'description' must use templates, is currently 'Instance has been unready for more than 15 minutes.'`,
	},
}

func TestLintPrometheusAlertsGuidelines(t *testing.T) {

	for _, alertTest := range alertTests {

		alerts := fmt.Sprintf(`
		{
			_config+:: {},
			prometheusAlerts+: {
				groups+: [
				  {
					name: 'test',
					rules: [
					  %s,
					],
				  },
				],
			},
		}
		`, alertTest.alert)

		filename, delete := writeTempFile(t, "alerts.jsonnet", alerts)
		defer delete()

		vm := jsonnet.MakeVM()
		errs := make(chan error)
		go lintPrometheus(filename, vm, errs)
		for err := range errs {
			if err.Error() != alertTest.expectedLintErr {
				t.Errorf("linting wrote unexpected output, expected '%s', got: %v", alertTest.expectedLintErr, err)
			}
		}
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

	if _, err := f.WriteString(contents); err != nil {
		t.Errorf("failed to write temp file to disk: %v", err)
	}

	if err := f.Close(); err != nil {
		t.Errorf("failed to close temp file: %v", err)
	}

	return f.Name(), func() { os.Remove(f.Name()) }
}
