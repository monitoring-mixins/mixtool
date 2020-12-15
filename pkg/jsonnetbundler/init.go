package jsonnetbundler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	v1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
)

// InitCommand is basically the same as jb init
func InitCommand(dir string) error {
	exists, err := jsonnetfile.Exists(jsonnetfile.File)

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

	err = ioutil.WriteFile(filename, contents, 0644)
	if err != nil {
		return fmt.Errorf("Failed to write new jsonnetfile.json, %s", err.Error())
	}

	return nil
}
