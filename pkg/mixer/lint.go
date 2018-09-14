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
	vm := jsonnet.MakeVM()

	vm.Importer(&jsonnet.FileImporter{
		JPaths: options.JPaths,
	})

	if options.Prometheus {
		if err := lintPrometheusAlerts(w, filename, vm); err != nil {
			return err
		}
		if err := lintPrometheusRules(w, filename, vm); err != nil {
			return err
		}
	}

	if options.Grafana {
		if err := lintGrafanaDashboards(w, filename, vm); err != nil {
			return err
		}
	}

	return nil
}

func lintPrometheusAlerts(w io.Writer, filename string, vm *jsonnet.VM) error {
	snippet := fmt.Sprintf("(import \"%s\").prometheusAlerts", filename)

	j, err := vm.EvaluateSnippet("", snippet)
	if err != nil {
		return err
	}

	_, errs := rulefmt.Parse([]byte(j))
	if errs != nil {
		for _, err := range errs {
			fmt.Fprintln(w, color.RedString(err.Error()))
		}
	}

	// TODO: Make some more verbose printing?
	//for _, g := range rg.Groups {
	//	fmt.Fprintf(w, "Group %s has %d alerts\n", g.Name, len(g.Rules))
	//}

	return err
}

func lintPrometheusRules(w io.Writer, filename string, vm *jsonnet.VM) error {
	snippet := fmt.Sprintf("(import '%s').prometheusRules", filename)

	j, err := vm.EvaluateSnippet("", snippet)
	if err != nil {
		return err
	}

	_, errs := rulefmt.Parse([]byte(j))
	if errs != nil {
		for _, err := range errs {
			fmt.Fprintln(w, color.RedString(err.Error()))
		}
	}

	// TODO: Make some more verbose printing?
	//for _, g := range rg.Groups {
	//	fmt.Fprintf(w, "Group %s has %d rules\n", g.Name, len(g.Rules))
	//}

	return err
}

func lintGrafanaDashboards(w io.Writer, filename string, vm *jsonnet.VM) error {
	snippet := fmt.Sprintf("(import '%s').grafanaDashboards", filename)

	j, err := vm.EvaluateSnippet("", snippet)
	if err != nil {
		return err
	}

	var dashboards map[string]interface{}
	if err := json.Unmarshal([]byte(j), &dashboards); err != nil {
		return err
	}

	for filename, dashboard := range dashboards {
		d := models.NewDashboardFromJson(simplejson.NewFromAny(dashboard))
		if d.Title == "" {
			fmt.Fprintln(w, color.RedString("Dashboard has no title: %s", filename))
		}
		if d.Uid == "" {
			fmt.Fprintln(w, color.YellowString("Dashboard has no UID, please set one for links to work"))
		}
	}

	return nil
}
