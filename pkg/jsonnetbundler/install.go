// Copyright 2018 jsonnet-bundler authors
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

package jsonnetbundler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/pkg/errors"

	"github.com/jsonnet-bundler/jsonnet-bundler/pkg"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	v1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	"github.com/jsonnet-bundler/jsonnet-bundler/spec/v1/deps"
)

func InstallCommand(dir, jsonnetHome string, uris []string, single bool) error {
	if dir == "" {
		dir = "."
	}

	jbfilebytes, err := os.ReadFile(filepath.Join(dir, jsonnetfile.File))
	if err != nil {
		return fmt.Errorf("failed to load jsonnetfile %s", err.Error())
	}

	jsonnetFile, err := jsonnetfile.Unmarshal(jbfilebytes)
	if err != nil {
		return err
	}

	jblockfilebytes, err := os.ReadFile(filepath.Join(dir, jsonnetfile.LockFile))
	if !os.IsNotExist(err) {
		if err != nil {
			return fmt.Errorf("failed to load lockfile %s", err.Error())
		}
	}

	lockFile, err := jsonnetfile.Unmarshal(jblockfilebytes)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(dir, jsonnetHome, ".tmp"), os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating vendor folder %s", err.Error())
	}

	for _, u := range uris {
		d := deps.Parse(dir, u)
		if d == nil {
			return fmt.Errorf("Unable to parse package URI %s", u)
		}

		if single {
			d.Single = true
		}

		if !depEqual(jsonnetFile.Dependencies[d.Name()], *d) {
			// the dep passed on the cli is different from the jsonnetFile
			jsonnetFile.Dependencies[d.Name()] = *d

			// we want to install the passed version (ignore the lock)
			delete(lockFile.Dependencies, d.Name())
		}
	}

	jsonnetPkgHomeDir := filepath.Join(dir, jsonnetHome)
	locked, err := pkg.Ensure(jsonnetFile, jsonnetPkgHomeDir, lockFile.Dependencies)
	if err != nil {
		return fmt.Errorf("failed to install packages %s", err)
	}

	pkg.CleanLegacyName(jsonnetFile.Dependencies)

	err = writeChangedJsonnetFile(jbfilebytes, &jsonnetFile, filepath.Join(dir, jsonnetfile.File))
	if err != nil {
		return fmt.Errorf("updating jsonnetfile.json %s", err)
	}

	err = writeChangedJsonnetFile(jblockfilebytes, &v1.JsonnetFile{Dependencies: locked}, filepath.Join(dir, jsonnetfile.LockFile))
	if err != nil {
		return fmt.Errorf("updating jsonnetfile.lock.json %s", err)
	}

	return nil
}

func depEqual(d1, d2 deps.Dependency) bool {
	name := d1.Name() == d2.Name()
	version := d1.Version == d2.Version
	source := reflect.DeepEqual(d1.Source, d2.Source)

	return name && version && source
}

func writeJSONFile(name string, d interface{}) error {
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return errors.Wrap(err, "encoding json")
	}
	b = append(b, []byte("\n")...)

	return os.WriteFile(name, b, 0644)
}

func writeChangedJsonnetFile(originalBytes []byte, modified *v1.JsonnetFile, path string) error {
	origJsonnetFile, err := jsonnetfile.Unmarshal(originalBytes)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(origJsonnetFile, *modified) {
		return nil
	}

	return writeJSONFile(path, *modified)
}
