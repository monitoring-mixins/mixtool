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

	"github.com/urfave/cli"
)

// type mixin struct {
// 	URL         string `json:"source"`
// 	Description string `json:"description,omitempty"`
// 	Name        string `json:"name"`
// 	Subdir      string `json:"subdir"`
// }

// type mixins struct {
// 	d map[string][]mixin
// }

func installCommand() cli.Command {
	return cli.Command{
		Name:        "install",
		Usage:       "Install a mixin",
		Description: "Install a mixin from a repository",
		Action:      installAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "bind-address",
				Usage: "Address to bind HTTP server to.",
			},
			cli.StringFlag{
				Name:  "prometheus",
				Value: "http://127.0.0.1:9090/",
				Usage: "location of the prometheus server",
			},
		},
	}
}

func installAction(c *cli.Context) error {
	filename := c.Args().First()
	if filename == "" {
		return fmt.Errorf("expected one argument, the name of the mixin. Show available mixins using mixtool list")
	}

	// process:

	// check if the name of the mixin exists in mixtool list

	// if not, check if the mixin url is valid - if so, interpret as a url to a repo

	// run jb install the mixin

	// run mixtool generate all

	// read files and run mixtool server

	return nil
}
