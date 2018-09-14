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

const dashboard = `local grafana = import 'grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local graphPanel = grafana.graphPanel;

{
  grafanaDashboards+:: {
    'dashboard.json':
      dashboard.new(
        'New Dashboard',
        time_from='now-1h',
      ).addTemplate(
        {
          current: {
            text: 'Prometheus',
            value: 'Prometheus',
          },
          hide: 0,
          label: null,
          name: 'datasource',
          options: [],
          query: 'prometheus',
          refresh: 1,
          regex: '',
          type: 'datasource',
        },
      )
      .addRow(
        row.new()
        .addPanel(
          graphPanel.new(
            'Graph',
            datasource='$datasource',
            span=6,
            format='short',
          )
          .addTarget(prometheus.target(
            'node_cpu{%(nodeExporterSelector)s, mode!="idle"})' % $._config,
            legendFormat='{{cpu}}'
          ))
        )
      ),
  },
}
`

func GenerateGrafanaDashboard() ([]byte, error) {
	return []byte(dashboard), nil
}

const alerts = `{
  prometheusAlerts+:: {
    groups+: [
      {
        name: 'kubernetes-resources',
        rules: [
          {
            alert: 'KubeNodeNotReady',
            expr: |||
              kube_node_status_condition{%(kubeStateMetricsSelector)s,condition="Ready",status="true"} == 0
            ||| % $._config,
            labels: {
              severity: 'warning',
            },
            annotations: {
              message: 'Overcommited CPU resource requests on Pods, cannot tolerate node failure.',
            },
            'for': '1h',
          },
        ],
      },
    ],
  }, 
}
`

func GeneratePrometheusAlerts() ([]byte, error) {
	return []byte(alerts), nil
}

const rules = `{
  prometheusRules+:: {
    groups+: [
      {
        name: 'example.rules',
        rules: [
          {
            record: 'node:node_memory_utilisation:ratio',
            expr: |||
              (node:node_memory_bytes_total:sum - node:node_memory_bytes_available:sum)
              /
              scalar(sum(node:node_memory_bytes_total:sum))
            |||,
          },
        ],
      },
    ],
  }
}
`

func GeneratePrometheusRules() ([]byte, error) {
	return []byte(rules), nil
}
