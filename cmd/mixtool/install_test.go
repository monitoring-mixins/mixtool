package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/monitoring-mixins/mixtool/pkg/mixer"
)

// Try to install every mixin from the mixin repository
// verify that each package generated has the yaml files
func TestInstallMixin(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "mixtool-install")
	if err != nil {
		t.Errorf("failed to make directory %v", err)
	}

	// defer os.RemoveAll(tmpdir)

	body, err := queryWebsite(defaultWebsite)
	if err != nil {
		t.Errorf("failed to query website %v", err)
	}
	mixins, err := parseMixinJSON(body)
	if err != nil {
		t.Errorf("failed to parse mixin body %v", err)
	}

	// download each mixin in turn
	for _, m := range mixins {

		generateCfg := mixer.GenerateOptions{
			AlertsFilename: "alerts.yaml",
			RulesFilename:  "rules.yaml",
			Directory:      "dashboards_out",
			JPaths:         []string{"vendor"},
			YAML:           true,
		}

		mixinURL := path.Join(m.URL, m.Subdir)

		fmt.Printf("installing %v\n", mixinURL)
		dldir := path.Join(tmpdir, m.Name)

		err = os.Mkdir(dldir, 0755)
		if err != nil {
			t.Errorf("failed to create directory %v", dldir)
		}

		err = downloadMixin(mixinURL, dldir)
		if err != nil {
			t.Errorf("failed to download mixin at %v. %v", mixinURL, err)
		}

		err = generateMixin(dldir, mixinURL, generateCfg)
		if err != nil {
			t.Errorf("failed to generate mixin yaml for %v. %v", mixinURL, err)
		}

		// // verify that the contents are correct
		// err = os.Chdir(dldir)
		// if err != nil {
		// 	t.Errorf("could not cd into %s", dldir)
		// }

		// if _, err := os.Stat("alerts.yaml"); os.IsNotExist(err) {
		// 	t.Errorf("expected alerts.yaml in %s", dldir)
		// }

		// if _, err := os.Stat("rules.yaml"); os.IsNotExist(err) {
		// 	t.Errorf("expected rules.yaml in %s", dldir)
		// }

		// if _, err := os.Stat("dashboards_out"); os.IsNotExist(err) {
		// 	t.Errorf("expected dashboards_out in %s", dldir)
		// }

	}

}
