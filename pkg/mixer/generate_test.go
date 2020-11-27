package mixer

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testAlertJsonnet = `
{
	prometheusAlerts+:: {
		groups+: [
		  {
			name: 'test-alerts',
			rules: [
			  {
				alert: 'TestAlert',
				expr: |||
				   test_alert == 1
				|||,
				labels: {
				  severity: 'warning',
				},
				annotations: {
				  message: 'test alert',
				},
				'for': '5m',
			  },
			],
		  },
		],
	  },
}
`

const expectedYaml = `groups:
- name: test-alerts
  rules:
  - alert: TestAlert
    annotations:
      message: test alert
    expr: |
      test_alert == 1
    for: 5m
    labels:
      severity: warning
`

func TestEvalAlerts(t *testing.T) {
	inFile, err := ioutil.TempFile(os.TempDir(), "mixtool-")
	assert.NoError(t, err)
	defer os.Remove(inFile.Name())

	_, err = inFile.Write([]byte(testAlertJsonnet))
	assert.NoError(t, err)

	out, err := GenerateAlerts(inFile.Name(), GenerateOptions{
		YAML: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, expectedYaml, string(out))
}
