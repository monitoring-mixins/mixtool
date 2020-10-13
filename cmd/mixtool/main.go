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
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/urfave/cli"
)

// Version of the mixtool.
// This is overridden at compile time.
var version = "0.0.0"

func main() {
	app := cli.NewApp()
	app.Name = "mixtool"
	app.Usage = "Improves your jsonnet mixins workflow"
	app.Description = "mixtool helps with generating, building and linting jsonnet mixins"
	app.Version = version

	app.Commands = cli.Commands{
		generateCommand(),
		lintCommand(),
		newCommand(),
		serverCommand(),
		// runbookCommand(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// If no jPath is given, we check if ./vendor exists in the current directory and use it.
// filename should be the path to the mixin root file (ie 'mixin.libsonnet').
func availableVendor(filename string, jPathsFlag []string) ([]string, error) {
	if len(jPathsFlag) == 0 {
		workdir, err := filepath.Abs(filename)
		if err != nil {
			return nil, err
		}

		workdir = filepath.Dir(workdir)
		root, err := jpath.FindParentFile("jsonnetfile.json", workdir, "/")
		if err != nil {
			if _, ok := err.(jpath.ErrorFileNotFound); ok {
				return jPathsFlag, nil
			}
			return nil, err
		}

		vendor := path.Join(root, "vendor")
		_, err = os.Stat(vendor)
		if err == nil {
			return []string{vendor}, nil
		}
	}

	return jPathsFlag, nil
}
