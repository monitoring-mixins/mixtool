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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/monitoring-mixins/mixtool/pkg/mixer"
	"github.com/pkg/errors"
)

// generateMixin generates the mixin given an jsonnet importString and writes
// to files specified in options
func generateAllMixin(importStr string, options mixer.GenerateOptions) error {
	fmt.Println("inside evaluateMixin")

	// rules

	out, err := mixer.GenerateInstallRules(importStr, options)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(options.RulesFilename, out, 0644)
	if err != nil {
		return err
	}

	// alerts

	out, err = mixer.GenerateInstallAlerts(importStr, options)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(options.AlertsFilename, out, 0644)
	if err != nil {
		return err
	}

	// dashboards

	dashboards, err := mixer.GenerateInstallDashboards(importStr, options)
	if err != nil {
		return err
	}

	if options.Directory == "" {
		return errors.New("missing directory flag to tell where to write to")
	}

	if err := os.MkdirAll(options.Directory, 0755); err != nil {
		return err
	}

	// Creating this func so that we can make proper use of defer
	writeDashboard := func(name string, dashboard json.RawMessage) error {
		file, err := os.Create(filepath.Join(options.Directory, name))
		if err != nil {
			return errors.Wrap(err, "failed to create dashboard file")
		}
		defer file.Close()

		file.Write(dashboard)

		return nil
	}

	for name, dashboard := range dashboards {
		if err := writeDashboard(name, dashboard); err != nil {
			return err
		}
	}

	return nil

}
