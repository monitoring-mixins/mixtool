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
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/urfave/cli"
)

func serverCommand() cli.Command {
	return cli.Command{
		Name:        "server",
		Usage:       "Start a server to provision Prometheus rule file(s) with.",
		Description: "Start a server to provision Prometheus rule file(s) with.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "bind-address",
				Usage: "Address to bind HTTP server to.",
			},
			cli.StringFlag{
				Name:  "prometheus-reload-url",
				Value: "http://127.0.0.1:9090/-/reload",
				Usage: "Prometheus address to reload after provisioning the rule file(s).",
			},
			cli.StringFlag{
				Name:  "config-file",
				Usage: "Prometheus configuration file",
			},
		},
		Action: serverAction,
	}
}

func serverAction(c *cli.Context) error {
	bindAddress := c.String("bind-address")
	http.Handle("/api/v1/rules", &ruleProvisioningHandler{
		ruleProvisioner: &ruleProvisioner{
			configFile: c.String("config-file"),
		},
		prometheusReloader: &prometheusReloader{
			prometheusReloadURL: c.String("prometheus-reload-url"),
		},
	})
	return http.ListenAndServe(bindAddress, nil)
}

type ruleProvisioningHandler struct {
	ruleProvisioner    *ruleProvisioner
	prometheusReloader *prometheusReloader
}

func (h *ruleProvisioningHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != "PUT" {
		http.Error(w, "Bad request: only PUT requests supported", http.StatusBadRequest)
		return
	}

	reloadNecessary, err := h.ruleProvisioner.provision(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	if reloadNecessary {
		if err := h.prometheusReloader.triggerReload(ctx); err != nil {
			http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

type ruleProvisioner struct {
	configFile string
}

// PUT request
// /api/v1/rules/<name>
// name specify
// filename determined by server

// provision attempts to provision the rule files read from r
// expects {rule-filename: "filename", data: "groups: ...."}
// makes new file and dumps data into rule-filename
// edits prometheus configuration to include that entry
func (p *ruleProvisioner) provision(r io.Reader) (bool, error) {

	configReader, err := os.OpenFile(p.configFile, os.O_RDWR, 0644)
	if err != nil {
		return false, fmt.Errorf("unable to open prometheus config file: %w", err)
	}

	configBuf, err := ioutil.ReadAll(configReader)
	if err != nil {
		return false, fmt.Errorf("unable to read prometheus config file: %w", err)
	}
	m := make(map[string]interface{})
	err = yaml.Unmarshal(configBuf, &m)
	if err != nil {
		return false, fmt.Errorf("unable to unmarshal prometheus config file: %w", err)
	}

	return true, nil
}

type prometheusReloader struct {
	prometheusReloadURL string
}

func (r *prometheusReloader) triggerReload(ctx context.Context) error {
	req, err := http.NewRequest("POST", r.prometheusReloadURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("reload request: %w", err)
	}

	if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
		return fmt.Errorf("exhausting request body: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("received non-200 response: %s; have you set `--web.enable-lifecycle` Prometheus flag?", resp.Status)
	}
	return nil
}
