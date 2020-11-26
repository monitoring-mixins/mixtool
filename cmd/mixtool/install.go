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
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/monitoring-mixins/mixtool/pkg/jsonnetbundler"
	"github.com/monitoring-mixins/mixtool/pkg/mixer"

	"github.com/urfave/cli"
)

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
				Value: "http://127.0.0.1:8080",
			},
			cli.StringFlag{
				Name:  "directory, d",
				Usage: "Path where the downloaded mixin is saved. If it doesn't exist already it will be created",
			},
			cli.BoolFlag{
				Name:  "put, p",
				Usage: "Specify this flag when you want to send PUT request to mixtool server once the mixins are generated",
			},
		},
	}
}

// Downloads a mixin from a given repository given by url and places into directory
// by running jb init and jb install
func downloadMixin(url string, jsonnetHome string, directory string) error {
	// intialize the jsonnet bundler library
	err := jsonnetbundler.InitCommand(directory)
	if err != nil {
		return fmt.Errorf("jsonnet bundler init failed %v", err)
	}

	// by default, set the single flag to false
	err = jsonnetbundler.InstallCommand(directory, jsonnetHome, []string{url}, false)
	if err != nil {
		return fmt.Errorf("jsonnet bundler install failed %v", err)
	}

	return nil
}

// Gets mixins from default website - mostly copied from list.go
func getMixins() ([]mixin, error) {
	body, err := queryWebsite(defaultWebsite)
	if err != nil {
		return nil, err
	}
	mixins, err := parseMixinJSON(body)
	if err != nil {
		return nil, err
	}
	return mixins, nil
}

func generateMixin(directory string, jsonnetHome string, mixinURL string, options mixer.GenerateOptions) ([]byte, error) {

	mixinBaseDirectory := filepath.Join(directory)

	err := os.Chdir(mixinBaseDirectory)
	if err != nil {
		return nil, fmt.Errorf("Cannot cd into directory %s", err)
	}

	files, err := filepath.Glob("*")
	fmt.Println("in generatemixin, directory is ", directory, files)

	// generate alerts, rules, grafana dashboards
	// empty files if not present

	u, err := url.Parse(mixinURL)
	if err != nil {
		return nil, fmt.Errorf("url parse %v", err)
	}

	// absolute directory is the same as the download url stripped of the scheme
	absDirectory := path.Join(u.Host, u.Path)
	absDirectory = strings.TrimLeft(absDirectory, "/:")
	absDirectory = strings.TrimRight(absDirectory, ".git")

	fmt.Println("absDirectory is", absDirectory)

	importFile := filepath.Join(absDirectory, "mixin.libsonnet")

	// generate rules, dashboards, alerts
	err = generateAll(importFile, options)
	if err != nil {
		return nil, fmt.Errorf("generateAll: %w", err)
	}

	out, err := generateRulesAlerts(importFile, options)
	if err != nil {
		return nil, fmt.Errorf("generateRulesAlerts %w", err)
	}

	return out, nil

}

func putMixin(content []byte, bindAddress string) error {
	u, err := url.Parse(bindAddress)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "/api/v1/rules")

	r := bytes.NewReader(content)

	req, err := http.NewRequest("PUT", u.String(), r)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("response from server %v", err)
	}
	if resp.StatusCode == 200 {
		fmt.Println("PUT alerts OK")
	} else {
		return fmt.Errorf("response code: %d resp is %v", resp.StatusCode, resp.Body)
	}
	return nil
}

func installAction(c *cli.Context) error {
	directory := c.String("directory")
	if directory == "" {
		return fmt.Errorf("Must specify a directory to download mixin")
	}

	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			return fmt.Errorf("could not create directory %v", err)
		}
	}

	mixinPath := c.Args().First()
	if mixinPath == "" {
		return fmt.Errorf("Expected the url of mixin repository or name of the mixin. Show available mixins using mixtool list")
	}

	mixinsList, err := getMixins()
	if err != nil {
		return fmt.Errorf("getMixins failed %v", err)
	}

	var mixinURL string
	if _, err := url.ParseRequestURI(mixinPath); err != nil {
		// check if the name exists in mixinsList
		found := false
		for _, m := range mixinsList {
			if m.Name == mixinPath {
				// join paths together
				u, err := url.Parse(m.URL)
				if err != nil {
					return fmt.Errorf("url parse failed %v", err)
				}
				u.Path = path.Join(u.Path, m.Subdir)
				mixinURL = u.String()
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("Could not find mixin with name %s", mixinPath)
		}
	} else {
		mixinURL = mixinPath
	}

	if mixinURL == "" {
		return fmt.Errorf("Empty mixinURL")
	}

	// by default jsonnet packages are downloaded under vendor
	jsonnetHome := "vendor"

	err = downloadMixin(mixinURL, jsonnetHome, directory)
	if err != nil {
		return err
	}

	generateCfg := mixer.GenerateOptions{
		AlertsFilename: "alerts.yaml",
		RulesFilename:  "rules.yaml",
		Directory:      "dashboards_out",
		JPaths:         []string{"./vendor"},
		YAML:           true,
	}

	rulesAlerts, err := generateMixin(directory, jsonnetHome, mixinURL, generateCfg)
	if err != nil {
		return err
	}

	// check if put address flag was set

	if c.Bool("put") {
		bindAddress := c.String("bind-address")
		// run put requests onto the server
		err = putMixin(rulesAlerts, bindAddress)
		if err != nil {
			return err
		}
	}

	return nil
}
