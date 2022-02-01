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
	"os"
	"path"
	"strings"

	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/urfave/cli"

	"github.com/monitoring-mixins/mixtool/pkg/mixer"
)

func installCommand() cli.Command {
	return cli.Command{
		Name:        "install",
		Usage:       "Install a mixin",
		Description: "Install a mixin from a repository",
		Action:      installAction,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name: "jpath, J",
			},
		},
	}
}

func installAction(c *cli.Context) error {
	filename := c.Args().First()
	if filename == "" {
		return fmt.Errorf("no jsonnet file given")
	}

	grafanaURL := os.Getenv("GRAFANA_URL")
	grafanaKey := os.Getenv("GRAFANA_TOKEN")
	client, err := gapi.New(grafanaURL, gapi.Config{APIKey: grafanaKey})
	if err != nil {
		return err
	}

	if _, err := client.Folders(); err != nil {
		return fmt.Errorf("failed to ping grafana: %v", err)
	}

	jPathFlag := c.StringSlice("jpath")
	jPathFlag, err = availableVendor(filename, jPathFlag)
	if err != nil {
		return err
	}

	generateCfg := mixer.GenerateOptions{
		AlertsFilename: "alerts.yaml",
		RulesFilename:  "rules.yaml",
		Directory:      "dashboards_out",
		JPaths:         jPathFlag,
		YAML:           true,
	}

	if err := generateAll(filename, generateCfg); err != nil {
		return err
	}

	ds, err := os.ReadDir("dashboards_out")
	if err != nil {
		return err
	}

	for _, d := range ds {
		if d.IsDir() {
			continue
		}

		buf, err := os.ReadFile(path.Join("dashboards_out", d.Name()))
		if err != nil {
			return err
		}

		var dashboardJson map[string]interface{}
		if err := json.Unmarshal(buf, &dashboardJson); err != nil {
			return err
		}

		uploadDashboard(client, dashboardJson)
	}

	return nil
}

func uploadDashboard(client *gapi.Client, dashboardJson map[string]interface{}) error {
	var uid string
	tmp, ok := dashboardJson["uid"]
	if !ok {
		return fmt.Errorf("missing uid from dashboard")
	}

	if uid, ok = tmp.(string); !ok {
		return fmt.Errorf("bad uid in dashboard")
	}

	dashboard, err := client.DashboardByUID(uid)
	if err != nil && !strings.HasPrefix(err.Error(), "status: 404") {
		return err
	}

	fmt.Printf("Updating dashboard %s (exists: %t)\n", uid, err == nil)

	if err != nil {
		dashboard = &gapi.Dashboard{}
	}

	dashboard.Model = dashboardJson
	dashboard.Overwrite = true
	if _, err := client.NewDashboard(*dashboard); err != nil {
		return err
	}

	return nil
}
