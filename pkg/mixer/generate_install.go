package mixer

import (
	"encoding/json"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

func GenerateInstallAlerts(importStr string, opts GenerateOptions) ([]byte, error) {
	vm := NewVM(opts.JPaths)

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

func GenerateInstallRules(importStr string, opts GenerateOptions) ([]byte, error) {
	vm := NewVM(opts.JPaths)

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

func GenerateInstallDashboards(importStr string, opts GenerateOptions) (map[string]json.RawMessage, error) {
	vm := NewVM(opts.JPaths)

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
