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

package jsonnetbundler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	v1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
)

// InitCommand is basically the same as jb init
func InitCommand(dir string) error {
	exists, err := jsonnetfile.Exists(jsonnetfile.File)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("jsonnetfile.json already exists")
	}

	s := v1.New()
	// TODO: disable them by default eventually
	// s.LegacyImports = false

	contents, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("formatting jsonnetfile contents as json, %s", err.Error())
	}
	contents = append(contents, []byte("\n")...)

	filename := filepath.Join(dir, jsonnetfile.File)

	err = os.WriteFile(filename, contents, 0644)
	if err != nil {
		return fmt.Errorf("failed to write new jsonnetfile.json, %s", err.Error())
	}

	return nil
}
