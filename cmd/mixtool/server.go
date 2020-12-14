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
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/urfave/cli"
)

const apiPath = "/api/v1/rules/"

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
	http.Handle(apiPath, &ruleProvisioningHandler{
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

	// TODO: might not be the best place to put this
	mixin := r.URL.Path[len(apiPath):]

	reloadNecessary, err := h.ruleProvisioner.provision(r.Body, mixin)
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
// specify mixin name
// filename determined by server
func (p *ruleProvisioner) provision(r io.Reader, mixinName string) (bool, error) {
	newRules, err := ioutil.ReadAll(r)
	if err != nil {
		return false, fmt.Errorf("unable to read new rules: %w", err)
	}

	mixinName = mixinName + ".yaml"
	dir := filepath.Dir(p.configFile)
	mixinFilename := filepath.Join(dir, mixinName)

	// if the filename under filepath.Join(dir, mixinName) already exists, do nothing
	if _, err = os.Stat(mixinFilename); err == nil {
		return true, nil
	}

	f, err := os.OpenFile(mixinFilename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return false, fmt.Errorf("unable to create new mixin file: %w", err)
	}

	// write all the contents into file
	n, err := f.Write(newRules)
	if err != nil {
		return false, fmt.Errorf("error when writing new rules: %w", err)
	}
	if n != len(newRules) {
		return false, fmt.Errorf("writing error, wrote %d bytes, expected %d", n, len(newRules))
	}

	f.Sync()
	f.Close()

	// add file's name to config file
	configBuf, err := ioutil.ReadFile(p.configFile)
	if err != nil {
		return false, fmt.Errorf("unable to open prometheus config file: %w", err)
	}

	m := make(map[string]interface{})
	err = yaml.Unmarshal(configBuf, &m)
	if err != nil {
		return false, fmt.Errorf("unable to unmarshal prometheus config file: %w", err)
	}

	for k, v := range m {
		if k == "rule_files" {
			// TODO: not entirely sure if this type assertion is safe
			rulemap := v.([]interface{})
			rulemap = append(rulemap, mixinName)
			m[k] = rulemap
			break
		}
	}

	// create a temporary config file
	tempfile, err := ioutil.TempFile(filepath.Dir(p.configFile), "temp-config")

	// marshal back into yaml
	newConfig, err := yaml.Marshal(m)
	if err != nil {
		return false, fmt.Errorf("failed to marhsal yaml: %w", err)
	}

	// write contents to temp config file
	n, err = tempfile.Write(newConfig)
	if err != nil {
		return false, fmt.Errorf("error when writing new rules: %w", err)
	}
	if n != len(newConfig) {
		return false, fmt.Errorf("writing error, wrote %d bytes, expected %d", n, len(newConfig))
	}

	tempfile.Sync()

	configReader, err := os.OpenFile(p.configFile, os.O_RDONLY, 0644)
	if err != nil {
		return false, fmt.Errorf("unable to read existing config: %w", err)
	}

	newConfigReader, err := os.OpenFile(tempfile.Name(), os.O_RDONLY, 0644)
	if err != nil {
		return false, fmt.Errorf("unable to open new config file: %w", err)
	}

	equal, err := readersEqual(configReader, newConfigReader)
	if err != nil {
		return false, fmt.Errorf("error from readersEqual: %w", err)
	}

	if equal {
		return false, nil
	}

	if err = os.Rename(tempfile.Name(), p.configFile); err != nil {
		return false, fmt.Errorf("cannot rename config file: %w", err)
	}
	return true, nil
}

type prometheusReloader struct {
	prometheusReloadURL string
}

func readersEqual(r1, r2 io.Reader) (bool, error) {
	buf1 := bufio.NewReader(r1)
	buf2 := bufio.NewReader(r2)
	for {
		b1, err1 := buf1.ReadByte()
		b2, err2 := buf2.ReadByte()
		if err1 != nil && !errors.Is(err1, io.EOF) {
			return false, err1
		}
		if err2 != nil && !errors.Is(err2, io.EOF) {
			return false, err2
		}
		if errors.Is(err1, io.EOF) || errors.Is(err2, io.EOF) {
			return err1 == err2, nil
		}
		if b1 != b2 {
			return false, nil
		}
	}
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
