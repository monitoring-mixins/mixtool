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

	"github.com/metalmatze/mixtool/pkg/mixer"
	"github.com/urfave/cli"
)

func generateCommand() cli.Command {
	return cli.Command{
		Name:        "generate",
		Usage:       "Generate jsonnet mixin files",
		Description: "Generate files for Prometheus alerts & rules and Grafana dashboards as jsonnet mixin",
		Subcommands: cli.Commands{
			cli.Command{
				Name:   "grafana-dashboard",
				Usage:  "Generate a new file for a Grafana dashboard",
				Action: generateGrafanaDashboard,
			},
			cli.Command{
				Name:   "prometheus-alerts",
				Usage:  "Generate a new file for Prometheus alerts",
				Action: generatePrometheusAlerts,
			},
			cli.Command{
				Name:   "prometheus-rules",
				Usage:  "Generate a new file for Prometheus rules",
				Action: generatePrometheusRules,
			},
		},
	}
}

func generateGrafanaDashboard(c *cli.Context) error {
	return writeFileToDisk(c, mixer.GenerateGrafanaDashboard)
}

func generatePrometheusAlerts(c *cli.Context) error {
	return writeFileToDisk(c, mixer.GeneratePrometheusAlerts)
}

func generatePrometheusRules(c *cli.Context) error {
	return writeFileToDisk(c, mixer.GeneratePrometheusRules)
}

func writeFileToDisk(c *cli.Context, creator func() ([]byte, error)) error {
	filename := c.Args().First()

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer f.Close()

	out, err := creator()
	if err != nil {
		return fmt.Errorf("failed to generate rules: %v", err)
	}

	f.Write(out)

	if err := f.Sync(); err != nil {
		return fmt.Errorf("failed to sync file to disk %s: %v", filename, err)
	}

	return nil
}
