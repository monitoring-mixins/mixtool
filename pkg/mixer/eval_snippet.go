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

// functions almost identical to eval.go except instead of reading contents from a file
// it formats in a importString and evaluates it instead

func evaluatePrometheusAlertsSnippet(vm *jsonnet.VM, importStr string) (string, error) {
	snippet := fmt.Sprintf(`
local mixin = (%s);

if std.objectHas(mixin, "prometheusAlerts")
then mixin.prometheusAlerts
else {}
`, importStr)

	return vm.EvaluateSnippet("", snippet)
}

func evaluatePrometheusRulesSnippet(vm *jsonnet.VM, importStr string) (string, error) {
	snippet := fmt.Sprintf(`
local mixin = (%s);

if std.objectHas(mixin, "prometheusRules")
then mixin.prometheusRules
else {}
`, importStr)

	return vm.EvaluateSnippet("", snippet)
}

func evaluateGrafanaDashboardsSnippet(vm *jsonnet.VM, importStr string) (string, error) {
	snippet := fmt.Sprintf(`
local mixin = (%s);

if std.objectHas(mixin, "grafanaDashboards")
then mixin.grafanaDashboards
else {}
`, importStr)

	return vm.EvaluateSnippet("", snippet)
}
