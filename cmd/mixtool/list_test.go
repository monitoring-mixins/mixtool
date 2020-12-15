package main

import (
	"io/ioutil"
	"os"
	"testing"
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
	err := ioutil.WriteFile("exampleMixinsTest.json", []byte(exampleMixins), 0644)
	if err != nil {
		t.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove("exampleMixinsTest.json")

	body, err := ioutil.ReadFile("exampleMixinsTest.json")
	if err != nil {
		t.Errorf("failed to read exampleMixinsTest.json %v", err)
	}
	mixins, err := parseMixinJSON(body)
	if err != nil {
		t.Errorf("failed to read exampleMixinsTest.json %v", err)
	}
	exampleMixins := map[string]bool{"ceph": true, "cool-mixin": true, "cortex": true}
	for _, m := range mixins {
		if _, ok := exampleMixins[m.Name]; !ok {
			t.Errorf("failed to find %v in exampleMixinsTest", m.Name)
		}
	}
}
