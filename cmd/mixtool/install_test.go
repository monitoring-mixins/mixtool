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
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/monitoring-mixins/mixtool/pkg/mixer"
	"github.com/stretchr/testify/assert"
)

// Try to install every mixin from the mixin repository
// verify that each package generated has the yaml files
func TestInstallMixin(t *testing.T) {
	t.Skip("Test is unreliable as it depends on external mixins.")

	body, err := queryWebsite(defaultWebsite)
	if err != nil {
		t.Errorf("failed to query website %v", err)
	}
	mixins, err := parseMixinJSON(body)
	if err != nil {
		t.Errorf("failed to parse mixin body %v", err)
	}

	// download each mixin in turn
	for _, m := range mixins {
		t.Run(m.Name, func(t *testing.T) {
			testInstallMixin(t, m)
		})
	}
}

func testInstallMixin(t *testing.T, m mixin) {
	tmpdir, err := os.CreateTemp("", "mixtool-install")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	generateCfg := mixer.GenerateOptions{
		AlertsFilename: "alerts.yaml",
		RulesFilename:  "rules.yaml",
		Directory:      "dashboards_out",
		JPaths:         []string{"vendor"},
		YAML:           true,
	}

	mixinURL := path.Join(m.URL, m.Subdir)

	fmt.Printf("installing %v\n", mixinURL)
	dldir := path.Join(tmpdir, m.Name+"mixin-test")

	err = os.Mkdir(dldir, 0755)
	assert.NoError(t, err)

	jsonnetHome := "vendor"

	err = downloadMixin(mixinURL, jsonnetHome, dldir)
	assert.NoError(t, err)

	_, err = generateMixin(dldir, jsonnetHome, mixinURL, generateCfg)
	assert.NoError(t, err)

	// verify that alerts, rules, dashboards exist
	err = os.Chdir(dldir)
	assert.NoError(t, err)

	if _, err := os.Stat("alerts.yaml"); os.IsNotExist(err) {
		t.Errorf("expected alerts.yaml in %s", dldir)
	}

	if _, err := os.Stat("rules.yaml"); os.IsNotExist(err) {
		t.Errorf("expected rules.yaml in %s", dldir)
	}

	if _, err := os.Stat("dashboards_out"); os.IsNotExist(err) {
		t.Errorf("expected dashboards_out in %s", dldir)
	}

	// verify that the output of alerts and rules matches using jsonnet
}
