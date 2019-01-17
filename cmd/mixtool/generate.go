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
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/monitoring-mixins/mixtool/pkg/mixer"
	"github.com/urfave/cli"
)

func generateCommand() cli.Command {
	flags := []cli.Flag{
		cli.StringSliceFlag{
			Name: "jpath, J",
		},
		cli.BoolTFlag{
			Name: "yaml, y",
		},
	}

	return cli.Command{
		Name:  "generate",
		Usage: "Generate manifests from jsonnet input",
		Subcommands: cli.Commands{
			cli.Command{
				Name:   "alerts",
				Usage:  "Generate Prometheus alerts based on the mixins",
				Flags:  flags,
				Action: generateAction(generateAlerts),
			},
			cli.Command{
				Name:   "rules",
				Usage:  "Generate Prometheus rules based on the mixins",
				Flags:  flags,
				Action: generateAction(generateRules),
			},
			cli.Command{
				Name:  "dashboards",
				Usage: "Generate Grafana dashboards based on the mixins",
				Flags: append(flags,
					cli.StringFlag{
						Name:  "directory, d",
						Usage: "The directory where Grafana dashboards are written to",
					},
				),
				Action: generateAction(generateDashboards),
			},
		},
	}
}

type generatorFunc func(io.Writer, string, mixer.GenerateOptions) error

func generateAction(generator generatorFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		jPathFlag := c.StringSlice("jpath")
		filename := c.Args().First()
		if filename == "" {
			return fmt.Errorf("no jsonnet file given")
		}

		jPathFlag = availableVendor(jPathFlag)

		generateCfg := mixer.GenerateOptions{
			Directory: c.String("directory"),
			JPaths:    jPathFlag,
			YAML:      c.BoolT("yaml"),
		}

		return generator(os.Stdout, filename, generateCfg)
	}
}

func generateAlerts(w io.Writer, filename string, options mixer.GenerateOptions) error {
	out, err := mixer.GenerateAlerts(filename, options)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}

func generateRules(w io.Writer, filename string, options mixer.GenerateOptions) error {
	out, err := mixer.GenerateRules(filename, options)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}

func generateDashboards(w io.Writer, filename string, opts mixer.GenerateOptions) error {
	if opts.Directory == "" {
		return errors.New("missing directory flag to tell where to write to")
	}

	dashboards, err := mixer.GenerateDashboards(filename, opts)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(opts.Directory, 0755); err != nil {
		return err
	}

	// Creating this func so that we can make proper use of defer
	writeDashboard := func(name string, dashboard json.RawMessage) error {
		file, err := os.Create(filepath.Join(opts.Directory, name))
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
