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

	"github.com/metalmatze/mixtool/pkg/mixer"
	"github.com/urfave/cli"
)

func runbookCommand() cli.Command {
	return cli.Command{
		Name:        "runbook",
		Usage:       "Generate a runbook markdown file",
		Description: "Generate a runbook markdown file from the jsonnet mixins",
		Action:      runbookAction,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name: "jpath, J",
			},
			cli.StringFlag{
				Name: "output-file, o",
			},
		},
	}
}

func runbookAction(c *cli.Context) error {
	outputFileFlag := c.String("output-file")
	jPathFlag := c.StringSlice("jpath")

	filename := c.Args().First()
	if filename == "" {
		return fmt.Errorf("no jsonnet file given")
	}

	jPathFlag = availableVendor(jPathFlag)

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

	err := mixer.Runbook(out, filename, mixer.RunbookOptions{
		JPaths: jPathFlag,
	})
	if err != nil {
		return err
	}

	return nil
}
