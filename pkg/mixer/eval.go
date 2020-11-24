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

	"github.com/google/go-jsonnet"
)

func evaluatePrometheusAlerts(vm *jsonnet.VM, filename string) (string, error) {
	snippet := fmt.Sprintf(`
local mixin = (import %q);

if std.objectHasAll(mixin, "prometheusAlerts")
then mixin.prometheusAlerts
else {}
`, filename)

	return vm.EvaluateSnippet("", snippet)
}

func evaluatePrometheusRules(vm *jsonnet.VM, filename string) (string, error) {
	snippet := fmt.Sprintf(`
local mixin = (import %q);

if std.objectHasAll(mixin, "prometheusRules")
then mixin.prometheusRules
else {}
`, filename)

	return vm.EvaluateSnippet("", snippet)
}

func evaluatePrometheusRulesAlerts(vm *jsonnet.VM, filename string) (string, error) {
	snippet := fmt.Sprintf(`
local mixin = (import %q);

if std.objectHasAll(mixin, "prometheusRules") && std.objectHasAll(mixin, "prometheusAlerts")
then mixin.prometheusRules + mixin.prometheusAlerts
else if std.objectHasAll(mixin, "prometheusRules")
then mixin.prometheusRules 
else if std.objectHasAll(mixin, "prometheusAlerts")
then mixin.prometheusAlerts
else {}
`, filename)

	return vm.EvaluateSnippet("", snippet)
}

func evaluateGrafanaDashboards(vm *jsonnet.VM, filename string) (string, error) {
	snippet := fmt.Sprintf(`
local mixin = (import %q);

if std.objectHasAll(mixin, "grafanaDashboards")
then mixin.grafanaDashboards
else {}
`, filename)

	return vm.EvaluateSnippet("", snippet)
}
