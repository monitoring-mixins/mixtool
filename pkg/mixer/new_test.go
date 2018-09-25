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

func TestNewGrafanaDashboard(t *testing.T) {
	d, err := NewGrafanaDashboard()
	if err != nil {
		t.Errorf("failed to create new Grafana dashboard: %v", err)
	}
	if string(d) != dashboard {
		t.Error("returned dashboard is not correct")
	}
}

func TestNewPrometheusAlerts(t *testing.T) {
	a, err := NewPrometheusAlerts()
	if err != nil {
		t.Errorf("failed to create new Prometheus alerts: %v", err)
	}
	if string(a) != alerts {
		t.Error("returned Prometheus alerts is not correct")
	}
}

func TestNewPrometheusRules(t *testing.T) {
	r, err := NewPrometheusRules()
	if err != nil {
		t.Errorf("failed to create new Prometheus rules: %v", err)
	}
	if string(r) != rules {
		t.Error("returned Prometheus rules is not correct")
	}
}
