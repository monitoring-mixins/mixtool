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

package mixer

import (
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	"github.com/google/go-jsonnet"
)

const markdownTemplate = `# Alert Runbook
{{range .Groups}}
### {{.Name}}
{{range .Rules}}
##### {{.Alert}}
+ *Severity*: {{.Labels.Severity}}
+ *Message*: ` + "`" + `{{.Annotations.Message}}` + "`" + `

{{.RunbookOutput}}
{{end}}{{end}}
`

type runbookJSON struct {
	Groups []struct {
		Name  string `json:"name"`
		Rules []struct {
			Alert  string `json:"alert"`
			Labels struct {
				Severity string `json:"severity"`
			} `json:"labels"`
			Annotations struct {
				Message string `json:"message"`
			} `json:"annotations"`
			RunbookOutput string `json:"runbookOutput"`
		} `json:"rules"`
	} `json:"groups"`
}

type RunbookOptions struct {
	JPaths       []string
	TemplateFile string
}

func Runbook(w io.Writer, filename string, opts RunbookOptions) error {
	vm := jsonnet.MakeVM()

	vm.Importer(&jsonnet.FileImporter{
		JPaths: opts.JPaths,
	})

	snippetTemplate := `
(import '%s') {
  prometheusAlerts+::
    local mapRuleGroups(f) = {
      groups: [
        group {
          rules: [
            f(rule)
            for rule in super.rules
          ],
        }
        for group in super.groups
      ],
    };
    local removeRunbookURL(rule) = rule {
      [if 'alert' in rule then 'runbookOutput']+:
        if 'runbook' in rule then super.runbook,
    };

    mapRuleGroups(removeRunbookURL),
}.prometheusAlerts
`
	snippet := fmt.Sprintf(snippetTemplate, filename)

	j, err := vm.EvaluateSnippet("", snippet)
	if err != nil {
		return fmt.Errorf("failed to evaluate snippet: %v", err)
	}

	var rj runbookJSON
	if err := json.Unmarshal([]byte(j), &rj); err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}

	var tmpl *template.Template
	if opts.TemplateFile != "" {
		tmpl, err = template.ParseFiles(opts.TemplateFile)
		if err != nil {
			return fmt.Errorf("failed to parse template from file: %v", err)
		}
	} else {
		tmpl, err = template.New("test").Parse(markdownTemplate)
		if err != nil {
			return fmt.Errorf("failed to parse template: %v", err)
		}
	}

	if err := tmpl.Execute(w, rj); err != nil {
		return fmt.Errorf("failed to execute template with data: %v", err)
	}

	return nil
}
