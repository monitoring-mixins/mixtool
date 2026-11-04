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
	"context"
	"fmt"
	"os"

	"github.com/monitoring-mixins/mixtool/pkg/mixer"
	"github.com/urfave/cli/v3"
)

type lintConfig struct {
	Prometheus bool
	Grafana    bool
	Vendor     []string
}

func lintCommand() *cli.Command {
	config := lintConfig{
		Prometheus: true,
		Grafana:    true,
	}

	return &cli.Command{
		Name:        "lint",
		Usage:       "Lint jsonnet files",
		Description: "Lint jsonnet files for correct structure of JSON objects",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Value:       true,
				Name:        "grafana",
				Usage:       "Lint Grafana dashboards against Grafana's schema",
				Destination: &config.Grafana,
			},
			&cli.BoolFlag{
				Value:       true,
				Name:        "prometheus",
				Usage:       "Lint Prometheus alerts and rules and their given expressions",
				Destination: &config.Prometheus,
			},
			&cli.StringSliceFlag{
				Name:    "jpath",
				Aliases: []string{"J"},
				Usage:   "Add folders to be used as vendor folders",
			},
		},
		Action: lintAction,
	}
}

func lintAction(_ context.Context, c *cli.Command) error {
	filename := c.Args().First()
	if filename == "" {
		return fmt.Errorf("expected one argument, the mixin filename")
	}

	jPath := c.StringSlice("jpath")
	jPath, err := availableVendor(filename, jPath)
	if err != nil {
		return err
	}

	options := mixer.LintOptions{
		JPaths:     jPath,
		Grafana:    c.Bool("grafana"),
		Prometheus: c.Bool("prometheus"),
	}

	if err := mixer.Lint(os.Stdout, filename, options); err != nil {
		return fmt.Errorf("failed to lint the file %s: %v", filename, err)
	}
	return nil
}
