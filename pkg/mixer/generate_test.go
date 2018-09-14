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
	"testing"
)

func TestGenerateGrafanaDashboard(t *testing.T) {
	d, err := GenerateGrafanaDashboard()
	if err != nil {
		t.Errorf("failed to generate Grafana dashboard: %v", err)
	}
	if string(d) != dashboard {
		t.Error("returned dashboard is not correct")
	}
}

func TestGeneratePrometheusAlerts(t *testing.T) {
	a, err := GeneratePrometheusAlerts()
	if err != nil {
		t.Errorf("failed to generate Prometheus alerts: %v", err)
	}
	if string(a) != alerts {
		t.Error("returned Prometheus alerts is not correct")
	}
}

func TestGeneratePrometheusRules(t *testing.T) {
	r, err := GeneratePrometheusRules()
	if err != nil {
		t.Errorf("failed to generate Prometheus rules: %v", err)
	}
	if string(r) != rules {
		t.Error("returned Prometheus rules is not correct")
	}
}
