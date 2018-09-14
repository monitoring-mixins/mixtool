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
	"io"
	"os"
	"path/filepath"

	"github.com/metalmatze/mixtool/pkg/mixer"
	"github.com/urfave/cli"
)

func buildCommand() cli.Command {
	return cli.Command{
		Name:  "build",
		Usage: "Build manifests from jsonnet input",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "jpath, J",
			},
			cli.StringFlag{
				Name: "multi, m",
			},
			cli.StringFlag{
				Name: "output-file, o",
			},
			cli.BoolTFlag{
				Name: "yaml, y",
			},
		},
		Action: buildAction,
	}
}

func buildAction(c *cli.Context) error {
	outputFileFlag := c.String("output-file")
	jPathFlag := c.String("jpath")
	multiFlag := c.String("multi")
	yamlFlag := c.BoolT("yaml")

	filename := c.Args().First()
	if filename == "" {
		return fmt.Errorf("no jsonnet file given")
	}

	var out io.Writer
	out = os.Stdout

	if outputFileFlag != "" {
		f, err := os.Create(outputFileFlag)
		if err != nil {
			return err
		}
		defer f.Close()

		out = f
	}

	buildCfg := mixer.BuildConfig{
		JPath: jPathFlag,
		YAML:  yamlFlag,
	}

	if multiFlag != "" {
		if err := os.MkdirAll(multiFlag, 0755); err != nil {
			return err
		}

		files, err := mixer.BuildMulti(filename, buildCfg)
		if err != nil {
			return err
		}

		for filename, content := range files {
			if yamlFlag {
				filename = filename + ".yaml"
			} else {
				filename = filename + ".json"
			}

			if err := writeMultiFileToDisk(multiFlag, filename, content); err != nil {
				return err
			}
		}

		return nil
	}

	b, err := mixer.Build(filename, buildCfg)
	if err != nil {
		return err
	}

	fmt.Fprint(out, b)

	return nil
}

func writeMultiFileToDisk(dir string, filename string, content string) error {
	f, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	return nil
}
