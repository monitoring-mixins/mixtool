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
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/google/go-jsonnet"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/models"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

type LintOptions struct {
	JPaths     []string
	Grafana    bool
	Prometheus bool
}

func Lint(w io.Writer, filename string, options LintOptions) error {
	errCount := 0

	if options.Prometheus {
		vm := NewVM(options.JPaths)
		errs := make(chan error)
		go lintPrometheus(filename, vm, errs)
		errCount += printErrs(w, errs)
	}

	if options.Grafana {
		vm := NewVM(options.JPaths)
		errs := make(chan error)
		go lintGrafanaDashboards(filename, vm, errs)
		errCount += printErrs(w, errs)
	}

	if errCount > 0 {
		return fmt.Errorf("%d lint errors found", errCount)
	}
	return nil
}

func printErrs(w io.Writer, errs <-chan error) int {
	errCount := 0
	for err := range errs {
		fmt.Fprintln(w, color.RedString(err.Error()))
		errCount++
	}
	return errCount
}

func lintPrometheus(filename string, vm *jsonnet.VM, errsOut chan<- error) {
	defer close(errsOut)

	j, err := evaluatePrometheusAlerts(vm, filename)
	if err != nil {
		errsOut <- err
		return
	}

	_, errs := rulefmt.Parse([]byte(j))
	for _, err := range errs {
		errsOut <- err
	}

	j, err = evaluatePrometheusRules(vm, filename)
	if err != nil {
		errsOut <- err
		return
	}

	_, errs = rulefmt.Parse([]byte(j))
	for _, err := range errs {
		errsOut <- err
	}

	return
}

func lintGrafanaDashboards(filename string, vm *jsonnet.VM, errsOut chan<- error) {
	defer close(errsOut)

	j, err := evaluateGrafanaDashboards(vm, filename)
	if err != nil {
		errsOut <- err
		return
	}

	var dashboards map[string]interface{}
	if err := json.Unmarshal([]byte(j), &dashboards); err != nil {
		errsOut <- err
		return
	}

	for filename, dashboard := range dashboards {
		d := models.NewDashboardFromJson(simplejson.NewFromAny(dashboard))
		if d.Title == "" {
			errsOut <- fmt.Errorf("dashboard has no title: %s", filename)
		}
		if d.Uid == "" {
			errsOut <- fmt.Errorf("dashboard has no UID, please set one for links to work: %s", filename)
		}
	}

	return
}
