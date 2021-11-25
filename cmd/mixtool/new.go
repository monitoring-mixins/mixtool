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

	"github.com/monitoring-mixins/mixtool/pkg/mixer"
	"github.com/urfave/cli"
)

func newCommand() cli.Command {
	return cli.Command{
		Name:        "new",
		Usage:       "Create new jsonnet mixin files",
		Description: "Create new files for Prometheus alerts & rules and Grafana dashboards as jsonnet mixin",
		Subcommands: cli.Commands{
			cli.Command{
				Name:   "grafana-dashboard",
				Usage:  "Create a new file with a Grafana dashboard mixin inside",
				Action: newGrafanaDashboard,
			},
			cli.Command{
				Name:   "prometheus-alerts",
				Usage:  "Create a new file with Prometheus alert mixins inside",
				Action: newPrometheusAlerts,
			},
			cli.Command{
				Name:   "prometheus-rules",
				Usage:  "Create a new file with Prometheus rule mixins inside",
				Action: newPrometheusRules,
			},
		},
	}
}

func newGrafanaDashboard(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return fmt.Errorf("expected filename as only argument")
	}

	filename := c.Args().First()
	if fileExists(filename) {
		return fmt.Errorf("file already exists. not overwriting")
	}

	return writeFileToDisk(filename, mixer.NewGrafanaDashboard)
}

func newPrometheusAlerts(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return fmt.Errorf("expected filename as only argument")
	}

	filename := c.Args().First()
	if fileExists(filename) {
		return fmt.Errorf("file already exists. not overwriting")
	}

	return writeFileToDisk(filename, mixer.NewPrometheusAlerts)
}

func newPrometheusRules(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return fmt.Errorf("expected filename as only argument")
	}

	filename := c.Args().First()
	if fileExists(filename) {
		return fmt.Errorf("file already exists. not overwriting")
	}

	return writeFileToDisk(filename, mixer.NewPrometheusRules)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}

	return true
}

func writeFileToDisk(filename string, creator func() ([]byte, error)) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer f.Close()

	out, err := creator()
	if err != nil {
		return fmt.Errorf("failed to create new rules: %v", err)
	}

	if _, err := f.Write(out); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("failed to sync file to disk %s: %v", filename, err)
	}

	return nil
}
