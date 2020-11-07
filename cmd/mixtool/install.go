// // Copyright 2018 mixtool authors
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //
// //     http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.

package main

// package main

// import (
// 	"fmt"
// 	"net/url"
// 	"os"
// 	"path"

// 	// _ "jsonnet-bundler-local/jsonnet-bundler/common"

// 	"github.com/urfave/cli"
// )

// // type mixin struct {
// // 	URL         string `json:"source"`
// // 	Description string `json:"description,omitempty"`
// // 	Name        string `json:"name"`
// // 	Subdir      string `json:"subdir"`
// // }

// func installCommand() cli.Command {
// 	return cli.Command{
// 		Name:        "install",
// 		Usage:       "Install a mixin",
// 		Description: "Install a mixin from a repository",
// 		Action:      installAction,
// 		Flags: []cli.Flag{
// 			cli.StringFlag{
// 				Name:  "bind-address",
// 				Usage: "Address to bind HTTP server to.",
// 				Value: "https://127.0.0.1:8080",
// 			},
// 			cli.StringFlag{
// 				Name:  "prometheus",
// 				Value: "http://127.0.0.1:9090/",
// 				Usage: "location of the prometheus server",
// 			},
// 			cli.StringFlag{
// 				Name:  "directory, d",
// 				Usage: "provide the path of where you would like the mixin to be downloaded",
// 			},
// 		},
// 	}
// }

// // Downloads a mixin from a given repository given by url and places into directory
// // by running jb init and jb install
// func downloadMixin(url string, directory string) error {
// 	os.Chdir(directory)
// 	fmt.Printf("downloading from %s\n", url)
// 	// intialize the jsonnet bundler library
// 	return nil
// }

// func installAction(c *cli.Context) error {
// 	directory := c.String("directory")
// 	if directory == "" {
// 		return fmt.Errorf("Expected directory which states where to put downloaded mixin into")
// 	}

// 	mixinPath := c.Args().First()
// 	if mixinPath == "" {
// 		return fmt.Errorf("Expected the url of mixin repository or name of the mixin. Show available mixins using mixtool list")
// 	}

// 	mixinsList, err := getMixins()
// 	if err != nil {
// 		return err
// 	}

// 	var mixinURL string
// 	if _, err := url.ParseRequestURI(mixinPath); err != nil {
// 		// check if the name exists in mixinsList
// 		found := false
// 		for _, m := range mixinsList {
// 			if m.Name == mixinPath {
// 				// join paths together
// 				u, err := url.Parse(m.URL)
// 				if err != nil {
// 					return err
// 				}
// 				u.Path = path.Join(u.Path, m.Subdir)
// 				mixinURL = u.String()
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			return fmt.Errorf("could not find mixin with name %s", mixinPath)
// 		}
// 	} else {
// 		mixinURL = mixinPath
// 	}

// 	// make mixinURL consistent
// 	if mixinURL != "" {
// 		err = downloadMixin(mixinURL, directory)
// 	}

// 	// run mixtool generate all

// 	// read files and run mixtool server

// 	// also need to reload grafana

// 	return nil
// }
