package mixer

import (
	"fmt"

	"github.com/google/go-jsonnet"
)

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
