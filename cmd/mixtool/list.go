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
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"

	"github.com/urfave/cli"
)

type mixin struct {
	URL         string `json:"source"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Subdir      string `json:"subdir"`
}

type mixins struct {
	d map[string][]mixin
}

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
		Timeout: time.Second * 3,
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

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func listAction(c *cli.Context) error {

	// if path is not specified, default to the mixins website
	// otherwise, try parse as url
	// otherwise, try look for a local json file

	defaultWebsite := "https://monitoring.mixins.dev/mixins.json"

	path := c.String("path")
	var body []byte
	var err error
	if path == "" {
		body, err = queryWebsite(defaultWebsite)
		if err != nil {
			return err
		}
	} else {
		// check if it's a url
		_, err := url.ParseRequestURI(path)
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
		} else {
			body, err = queryWebsite(path)
			if err != nil {
				return err
			}
		}
	}

	var mixins map[string][]mixin
	jsonErr := json.Unmarshal(body, &mixins)

	if jsonErr != nil {
		return errors.New("failed to unmarshal json. Maybe your json schema is incorrect?")
	}

	mixinsList := mixins["mixins"]

	writer := tabwriter.NewWriter(os.Stdout, 4, 8, 0, '\t', tabwriter.TabIndent)
	fmt.Fprintln(writer, "name\tdescription")
	fmt.Fprintln(writer, "----\t----------")
	// print all the mixins out
	for _, m := range mixinsList {
		// chop off the description if it's too long
		m.Description = strings.Replace(m.Description, "\n", "", -1)
		if len(m.Description) <= 0 {
			m.Description = "N/A"
		}
		fmt.Fprintf(writer, "%s\t%s\n", color.GreenString(m.Name), m.Description)
	}

	return writer.Flush()
}
