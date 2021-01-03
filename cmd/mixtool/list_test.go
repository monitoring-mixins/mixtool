package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleMixins = `
{
    "mixins": [
        {
            "name": "ceph",
            "source": "https://github.com/ceph/ceph-mixins",
            "subdir": "",
            "description": "A set of Prometheus alerts for Ceph.\n\nThe scope of this project is to provide Ceph specific Prometheus rule files using Prometheus Mixins.\n"
        },
        {
            "name": "cortex",
            "source": "https://github.com/grafana/cortex-jsonnet",
            "subdir": "cortex-mixin"
        },
        {
            "name": "cool-mixin",
            "source": "https://github.com",
			"subdir": "cool-mixin",
			"description": "A fantastic mixin"
        }
    ]
}
`

func TestList(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "exampleMixinsTest.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	err = ioutil.WriteFile(tempFile.Name(), []byte(exampleMixins), 0644)
	assert.NoError(t, err)

	body, err := ioutil.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	mixins, err := parseMixinJSON([]byte(body))
	assert.NoError(t, err)

	exampleMixins := map[string]bool{"ceph": true, "cool-mixin": true, "cortex": true}
	for _, m := range mixins {
		if _, ok := exampleMixins[m.Name]; !ok {
			t.Errorf("failed to find %v in exampleMixinsTest", m.Name)
		}
	}
}
