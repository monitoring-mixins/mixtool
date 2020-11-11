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
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

type mixin struct {
	URL         string `json:"source"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Subdir      string `json:"subdir"`
}

const defaultWebsite = "https://monitoring.mixins.dev/mixins.json"

func listCommand() cli.Command {
	return cli.Command{
		Name:        "list",
		Usage:       "List all available mixins",
		Description: "List all available mixins as presented on monitoring.mixins.dev",
		Action:      listAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "path, p",
				Usage: "provide the path of a url with a json endpoint or a local json file",
			},
		},
	}
}

func queryWebsite(mixinsWebsite string) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, mixinsWebsite, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "mixtool-list")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func printMixins(mixinsList []mixin) error {
	writer := tabwriter.NewWriter(os.Stdout, 4, 8, 0, '\t', tabwriter.TabIndent)
	fmt.Fprintln(writer, "name")
	fmt.Fprintln(writer, "----")
	for _, m := range mixinsList {
		// for now, do not print out any of the descriptions
		// m.Description = strings.Replace(m.Description, "\n", "", -1)
		// // maybe truncate the description if it's too long
		// if len(m.Description) <= 0 {
		// 	m.Description = "N/A"
		// }
		fmt.Fprintf(writer, "%s\n", color.GreenString(m.Name))
	}

	return writer.Flush()
}

// ParsesMixinJSON expects a top level key mixins which contains a list of mixins
func parseMixinJSON(body []byte) ([]mixin, error) {
	var mixins map[string][]mixin
	if err := json.Unmarshal(body, &mixins); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}
	mixinsList := mixins["mixins"]
	return mixinsList, nil
}

// if path is not specified, default to the mixins website
// otherwise, try parse as url
// otherwise, try look for a local json file
func listAction(c *cli.Context) error {
	path := c.String("path")
	var body []byte
	var err error
	if path == "" {
		body, err = queryWebsite(defaultWebsite)
		if err != nil {
			return err
		}
		mixins, err := parseMixinJSON(body)
		if err != nil {
			return err
		}
		return printMixins(mixins)
	}

	_, err = url.ParseRequestURI(path)
	if err != nil {
		// check if it's a local json file
		_, err = os.Stat(path)
		if err != nil {
			return err
		}
		body, err = ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		mixins, err := parseMixinJSON(body)
		if err != nil {
			return err
		}
		return printMixins(mixins)
	}

	body, err = queryWebsite(path)
	if err != nil {
		return err
	}

	mixins, err := parseMixinJSON(body)
	if err != nil {
		return err
	}

	return printMixins(mixins)
}
