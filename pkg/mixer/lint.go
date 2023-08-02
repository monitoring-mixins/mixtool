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
	"path"
	"regexp"

	"github.com/fatih/color"
	"github.com/google/go-jsonnet"
	"github.com/grafana/dashboard-linter/lint"
	"github.com/prometheus/prometheus/model/rulefmt"
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

	// Lint using the config file from grafana/dashboard-linter
	config := lint.NewConfigurationFile()
	configFilename := path.Join(path.Dir(filename), ".lint")
	if err := config.Load(configFilename); err != nil {
		errsOut <- err
	}

	j, err := evaluatePrometheusAlerts(vm, filename)
	if err != nil {
		errsOut <- err
		return
	}
	groups, errs := rulefmt.Parse([]byte(j))
	for _, err := range errs {
		errsOut <- err
	}

	for _, g := range groups.Groups {
		for _, r := range g.Rules {
			errs = lintPrometheusAlertsGuidelines(&r, config)
			for _, err := range errs {
				errsOut <- err
			}
		}
	}

	_, err = evaluatePrometheusRules(vm, filename)
	if err != nil {
		errsOut <- err
		return
	}

}

var camelCaseRegexp = regexp.MustCompile(`^([A-Z]+[a-z0-9]+)+$`)
var goTemplateRegexp = regexp.MustCompile(`\{\{.+}\}`)
var sentenceRegexp = regexp.MustCompile(`^[A-Z].+\.$`)

// Enforces alerting guidelines.
// https://monitoring.mixins.dev/#guidelines-for-alert-names-labels-and-annotations
func lintPrometheusAlertsGuidelines(rule *rulefmt.RuleNode, cf *lint.ConfigurationFile) (errs []error) {
	if !isLintExcluded("alert-name-camelcase", rule.Alert.Value, cf) {
		if !camelCaseRegexp.MatchString(rule.Alert.Value) {
			errs = append(errs, fmt.Errorf("[alert-name-camelcase] Alert '%s' name is not in camel case", rule.Alert.Value))
		}
	}

	if !isLintExcluded("alert-severity-rule", rule.Alert.Value, cf) {
		if rule.Labels["severity"] != "warning" && rule.Labels["severity"] != "critical" && rule.Labels["severity"] != "info" {
			errs = append(errs, fmt.Errorf("[alert-severity-rule] Alert '%s' severity must be 'warning', 'critical' or 'info', is currently '%s'", rule.Alert.Value, rule.Labels["severity"]))
		}
	}

	if _, ok := rule.Annotations["description"]; !ok {
		if !isLintExcluded("alert-description-missing-rule", rule.Alert.Value, cf) {
			errs = append(errs, fmt.Errorf("[alert-description-missing-rule] Alert '%s' must have annotation 'description'", rule.Alert.Value))
		}
	} else {
		if !isLintExcluded("alert-description-templating", rule.Alert.Value, cf) {
			if !goTemplateRegexp.MatchString(rule.Annotations["description"]) {
				errs = append(errs, fmt.Errorf("[alert-description-templating] Alert %s annotation 'description' must use templates, is currently '%s'", rule.Alert.Value, rule.Annotations["description"]))
			}
		}
	}

	if _, ok := rule.Annotations["summary"]; !ok {
		if !isLintExcluded("alert-summary-missing-rule", rule.Alert.Value, cf) {
			errs = append(errs, fmt.Errorf("[alert-summary-missing-rule] Alert '%s' must have annotation 'summary'", rule.Alert.Value))
		}
	} else {
		if goTemplateRegexp.MatchString(rule.Annotations["summary"]) {
			if !isLintExcluded("alert-summary-templating", rule.Alert.Value, cf) {
				errs = append(errs, fmt.Errorf("[alert-summary-templating] Alert %s annotation 'summary' must not use templates", rule.Alert.Value))
			}
		}
		if !sentenceRegexp.MatchString(rule.Annotations["summary"]) {
			if !isLintExcluded("alert-summary-style", rule.Alert.Value, cf) {
				errs = append(errs, fmt.Errorf("[alert-summary-style] Alert %s annotation 'summary' must start with capital letter and end with period, is currently '%s'", rule.Alert.Value, rule.Annotations["summary"]))
			}
		}
	}
	return errs
}

func lintGrafanaDashboards(filename string, vm *jsonnet.VM, errsOut chan<- error) {
	defer close(errsOut)

	j, err := evaluateGrafanaDashboards(vm, filename)
	if err != nil {
		errsOut <- err
		return
	}

	var dashboards map[string]json.RawMessage
	if err := json.Unmarshal([]byte(j), &dashboards); err != nil {
		errsOut <- err
		return
	}

	rules := lint.NewRuleSet()

	for dashboardFilename, raw := range dashboards {
		var db map[string]interface{}
		if err := json.Unmarshal(raw, &db); err != nil {
			errsOut <- err
			continue
		}

		var title, uid string
		if t, ok := db["title"]; ok {
			title, _ = t.(string)
		}
		if u, ok := db["uid"]; ok {
			uid, _ = u.(string)
		}

		if title == "" {
			errsOut <- fmt.Errorf("dashboard has no title: %s", dashboardFilename)
		}
		if uid == "" {
			errsOut <- fmt.Errorf("dashboard has no UID, please set one for links to work: %s", dashboardFilename)
		}

		// Lint using the new grafana/dashboard-linter project.
		config := lint.NewConfigurationFile()
		configFilename := path.Join(path.Dir(filename), ".lint")
		if err := config.Load(configFilename); err != nil {
			errsOut <- err
			continue
		}

		dash, err := lint.NewDashboard(raw)
		if err != nil {
			errsOut <- err
			continue
		}

		rs, err := rules.Lint([]lint.Dashboard{dash})
		if err != nil {
			errsOut <- err
			continue
		}

		for rule, results := range rs.ByRule() {
			for _, result := range results {
				result = config.Apply(result)
				for _, r := range result.Result.Results {
					switch r.Severity {
					case lint.Exclude, lint.Success, lint.Quiet:
					default:
						errsOut <- fmt.Errorf("[%s] '%s': %s", rule, result.Dashboard.Title, r.Message)
					}
				}
			}
		}
	}
}

func isLintExcluded(ruleName string, alertName string, cf *lint.ConfigurationFile) bool {
	exclusions, ok := cf.Exclusions[ruleName]
	if exclusions != nil {
		for _, ce := range exclusions.Entries {
			if alertName == ce.Alert {
				return true
			}
		}
		if len(exclusions.Entries) == 0 {
			return true
		}
	} else if ok {
		return true
	}
	return false
}
