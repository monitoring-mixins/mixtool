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
	"encoding/json"
	"fmt"

	"github.com/google/go-jsonnet"
	"github.com/monitoring-mixins/mixtool/pkg/mixer"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

func evaluateInstallAlerts(importStr string, opts mixer.GenerateOptions) ([]byte, error) {
	vm := mixer.NewVM(opts.JPaths)

	j, err := evaluateAlerts(vm, importStr)
	if err != nil {
		return nil, err
	}

	output := []byte(j)

	if opts.YAML {
		output, err = yaml.JSONToYAML(output)
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}

func evaluateInstallRules(importStr string, opts mixer.GenerateOptions) ([]byte, error) {
	vm := mixer.NewVM(opts.JPaths)

	j, err := evaluateRules(vm, importStr)
	if err != nil {
		return nil, err
	}

	output := []byte(j)

	if opts.YAML {
		output, err = yaml.JSONToYAML(output)
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}

func evaluateInstallDashboards(importStr string, opts mixer.GenerateOptions) (map[string]json.RawMessage, error) {
	vm := mixer.NewVM(opts.JPaths)

	j, err := evaluateDashboards(vm, importStr)
	if err != nil {
		return nil, err
	}

	var dashboards map[string]json.RawMessage
	if err := json.Unmarshal([]byte(j), &dashboards); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal dashboards")
	}

	return dashboards, nil

}

func evaluateRules(vm *jsonnet.VM, importStr string) (string, error) {
	snippet := fmt.Sprintf(`
	local mixin = (%s);
	if std.objectHasAll(mixin, "prometheusRules")
	then mixin.prometheusRules
	else {}
	`, importStr)
	return vm.EvaluateSnippet("", snippet)
}

func evaluateAlerts(vm *jsonnet.VM, importStr string) (string, error) {
	snippet := fmt.Sprintf(`
	local mixin = (%s);
	if std.objectHasAll(mixin, "prometheusAlerts")
	then mixin.prometheusAlerts
	else {}
	`, importStr)
	return vm.EvaluateSnippet("", snippet)
}

func evaluateDashboards(vm *jsonnet.VM, importStr string) (string, error) {
	snippet := fmt.Sprintf(`
	local mixin = (%s);
	if std.objectHasAll(mixin, "grafanaDashboards")
	then mixin.grafanaDashboards
	else {}
	`, importStr)

	return vm.EvaluateSnippet("", snippet)
}
