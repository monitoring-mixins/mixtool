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
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/monitoring-mixins/mixtool/pkg/jsonnetbundler"
	"github.com/monitoring-mixins/mixtool/pkg/mixer"
	"github.com/pkg/errors"

	"github.com/urfave/cli"
)

// type mixin struct {
// 	URL         string `json:"source"`
// 	Description string `json:"description,omitempty"`
// 	Name        string `json:"name"`
// 	Subdir      string `json:"subdir"`
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
				Value: "http://127.0.0.1:8080",
			},
			cli.StringFlag{
				Name:  "rule-file",
				Usage: "File to provision rules into.",
			},
			cli.StringFlag{
				Name:  "directory, d",
				Usage: "Path where the downloaded mixin is saved. If it doesn't exist already it will be created",
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
		return err
	}

	// by default, set the single flag to false
	err = jsonnetbundler.InstallCommand(directory, jsonnetHome, []string{url}, false)
	if err != nil {
		fmt.Println(err)
		return err
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

func generateMixin(directory string, jsonnetHome string, mixinURL string, options mixer.GenerateOptions) error {
	fmt.Println("running generate all")

	mixinBaseDirectory := filepath.Join(directory)

	err := os.Chdir(mixinBaseDirectory)
	if err != nil {
		return fmt.Errorf("Cannot cd into directory %s", err)
	}

	files, err := filepath.Glob("*")
	fmt.Println("in generatemixin, directory is ", directory, files)

	// generate alerts, rules, grafana dashboards
	// empty files if not present

	u, err := url.Parse(mixinURL)
	if err != nil {
		return err
	}

	// absolute directory is the same as the download url stripped of the scheme
	absDirectory := path.Join(u.Host, u.Path)

	fmt.Println("absDirectory is", absDirectory)

	// create a temporary jsonnet file that
	// imports mixin.libsonnet as the "main" file in the mixin configfuration
	// then run generate all passing in the vendor folder to find dependencies
	// the name of the mixin folder inside vendor seems to be the last fragment of mixin's url + subdir
	// TODO: need to get absolute path of the mixin

	// note - need to somehow explicitly pick up +:: hidden fields in thanos
	tempContent := fmt.Sprintf(
		`import "%s"`, filepath.Join(absDirectory, "mixin.libsonnet"))

	// generate rules, dashboards, alerts
	err = evaluateMixin(tempContent, options)
	if err != nil {
		return err
	}

	// 	tempContent := fmt.Sprintf(
	// 		`local mixin = (import "%s");
	// mixin.grafanaDashboards`, filepath.Join(absDirectory, "mixin.libsonnet"))

	// evaluate prometheus rules and alerts
	// since generateall expects a filename but here we do not need to provide a filename
	// err = ioutil.WriteFile("temp.jsonnet", []byte(tempContent), 0644)
	// if err != nil {
	// 	return err
	// }

	// err = generateAll("temp.jsonnet", options)
	// if err != nil {
	// 	return err
	// }

	// err = os.Remove("temp.jsonnet")
	// if err != nil {
	// 	return err
	// }
	// return nil
	return nil

}

// generateMixin generates the mixin given an jsonnet importString and writes
// to files specified in options
func evaluateMixin(importStr string, options mixer.GenerateOptions) error {
	fmt.Println("inside evaluateMixin")

	// rules

	out, err := evaluateInstallRules(importStr, options)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(options.RulesFilename, out, 0644)
	if err != nil {
		return err
	}

	// alerts

	out, err = evaluateInstallAlerts(importStr, options)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(options.AlertsFilename, out, 0644)
	if err != nil {
		return err
	}

	// dashboards

	dashboards, err := evaluateInstallDashboards(importStr, options)
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

func putMixin(directory string, mixinURL string, bindAddress string, options mixer.GenerateOptions) error {
	fmt.Println("running put Mixin")

	wd, err := os.Getwd()
	fmt.Println("path is", wd)
	if err != nil {
		return err
	}

	// err := os.Chdir(filepath.Join()
	// if err != nil {
	// 	return fmt.Errorf("Cannot cd into directory %s", err)
	// }

	// alerts.yaml
	alertsFilename := options.AlertsFilename
	alertsReader, err := os.Open(alertsFilename)
	if err != nil {
		return err
	}

	// rules.yaml
	rulesFilename := options.RulesFilename
	rulesReader, err := os.Open(rulesFilename)
	if err != nil {
		return err
	}

	u, err := url.Parse(bindAddress)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "/api/v1/rules")

	// request for rules
	req, err := http.NewRequest("PUT", u.String(), rulesReader)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("response from server %v", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("OK")
	} else {
		fmt.Printf("resp is %v\n", resp.Body)
	}

	// same request but for alerts
	req, err = http.NewRequest("PUT", u.String(), alertsReader)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("response from server %v", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("OK")
	} else {
		fmt.Printf("resp is %v\n", resp.Body)
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
			return err
		}
	}

	mixinPath := c.Args().First()
	if mixinPath == "" {
		return fmt.Errorf("Expected the url of mixin repository or name of the mixin. Show available mixins using mixtool list")
	}

	mixinsList, err := getMixins()
	if err != nil {
		return err
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
					return err
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
		fmt.Println(err)
		return err
	}

	generateCfg := mixer.GenerateOptions{
		AlertsFilename: "alerts.yaml",
		RulesFilename:  "rules.yaml",
		Directory:      "dashboards_out",
		JPaths:         []string{"./vendor"},
		YAML:           true,
	}

	err = generateMixin(directory, jsonnetHome, mixinURL, generateCfg)
	if err != nil {
		return err
	}

	// bindAddress := c.String("bind-address")
	// // run put requests onto the server
	// err = putMixin(directory, mixinURL, bindAddress, generateCfg)
	// if err != nil {
	// 	return err
	// }

	return nil
}
