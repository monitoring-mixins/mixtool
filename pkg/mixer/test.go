package mixer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/google/go-jsonnet"
	"github.com/metalmatze/mixtool/pkg/promtool"
)

type TestOptions struct {
	JPaths []string
}

func Test(w io.Writer, filename string, options TestOptions) error {
	vm := jsonnet.MakeVM()

	vm.Importer(&jsonnet.FileImporter{
		JPaths: options.JPaths,
	})

	alerts, err := tempAlerts(vm, filename)
	if err != nil {
		return err
	}
	defer os.Remove(alerts)

	tests, err := tempTests(vm, filename, alerts)
	if err != nil {
		return err
	}
	defer os.Remove(tests)

	exit := promtool.RulesUnitTest(tests)
	if exit != 0 {
		return fmt.Errorf("exit code %d", exit)
	}

	return nil
}

func tempAlerts(vm *jsonnet.VM, filename string) (string, error) {
	snippet := fmt.Sprintf(`(import "%s").prometheusAlerts`, filename)
	j, err := vm.EvaluateSnippet("", snippet)
	if err != nil {
		return "", err
	}

	y, err := yaml.JSONToYAML([]byte(j))
	if err != nil {
		return "", err
	}

	tmpAlerts, err := ioutil.TempFile("", "alerts")
	if err != nil {
		return "", err
	}

	if _, err := tmpAlerts.Write(y); err != nil {
		return "", err
	}

	if err := tmpAlerts.Close(); err != nil {
		return "", err
	}

	return tmpAlerts.Name(), nil
}

var testsOutput = `
(import '%s') {
  prometheusAlerts+::
    local mapRuleGroups(f) = {
      groups: [
        group {
          rules: [
            f(rule)
            for rule in super.rules
          ],
        }
        for group in super.groups
      ],
    };
    local outputTests(rule) = rule {
      [if 'alert' in rule then 'testsOutput']+:
        if 'tests' in rule then super.tests,
    };

    mapRuleGroups(outputTests),
}.prometheusAlerts
`

type alertsAndTests struct {
	Groups []struct {
		Rules []struct {
			TestOutput []map[string]interface{} `json:"testsOutput"`
		} `json:"rules"`
	} `json:"groups"`
}

type testFile struct {
	RuleFiles []string                 `json:"rule_files"`
	Tests     []map[string]interface{} `json:"tests"`
}

func tempTests(vm *jsonnet.VM, filename string, alerts string) (string, error) {
	snippet := fmt.Sprintf(testsOutput, filename)

	j, err := vm.EvaluateSnippet("", snippet)
	if err != nil {
		return "", err
	}

	var at alertsAndTests
	if err := json.Unmarshal([]byte(j), &at); err != nil {
		return "", err
	}

	tf := testFile{
		RuleFiles: []string{alerts},
	}

	for _, gr := range at.Groups {
		for _, rules := range gr.Rules {
			if len(rules.TestOutput) > 0 {
				tf.Tests = append(tf.Tests, rules.TestOutput...)
			}
		}
	}

	y, err := yaml.Marshal(tf)
	if err != nil {
		return "", err
	}

	tests, err := ioutil.TempFile("", "tests")
	if err != nil {
		return "", err
	}

	if _, err := tests.Write(y); err != nil {
		return "", err
	}

	if err := tests.Close(); err != nil {
		return "", err
	}

	return tests.Name(), nil
}
