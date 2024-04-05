# mixtool

> NOTE: This project is *alpha* stage. Flags, configuration, behavior and design may change significantly in following releases.

The mixtool is a helper for easily working with [jsonnet](http://jsonnet.org/) mixins.

## Install

Make sure you're using golang v1.21 or higher, and run:

```
go install github.com/monitoring-mixins/mixtool/cmd/mixtool@main
```

## Usage

All command line flags:

[embedmd]:# (_output/help.txt)
```txt
NAME:
   mixtool - Improves your jsonnet mixins workflow

USAGE:
   mixtool [global options] command [command options] [arguments...]

VERSION:
   v0.1.0-pre

DESCRIPTION:
   mixtool helps with generating, building and linting jsonnet mixins

COMMANDS:
   generate  Generate manifests from jsonnet input
   lint      Lint jsonnet files
   new       Create new jsonnet mixin files
   server    Start a server to provision Prometheus rule file(s) with.
   list      List all available mixins
   install   Install a mixin
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Generate
[embedmd]:# (_output/help-generate.txt)
```txt
NAME:
   mixtool generate - Generate manifests from jsonnet input

USAGE:
   mixtool generate command [command options] [arguments...]

COMMANDS:
   alerts      Generate Prometheus alerts based on the mixins
   rules       Generate Prometheus rules based on the mixins
   dashboards  Generate Grafana dashboards based on the mixins
   all         Generate all resources - Prometheus alerts, Prometheus rules and Grafana dashboards

OPTIONS:
   --help, -h  show help
   
```

### New

[embedmd]:# (_output/help-new.txt)
```txt
NAME:
   mixtool new - Create new files for Prometheus alerts & rules and Grafana dashboards as jsonnet mixin

USAGE:
   mixtool new command [command options] [arguments...]

COMMANDS:
   grafana-dashboard  Create a new file with a Grafana dashboard mixin inside
   prometheus-alerts  Create a new file with Prometheus alert mixins inside
   prometheus-rules   Create a new file with Prometheus rule mixins inside

OPTIONS:
   --help, -h  show help
   
```

#### New Examples

```bash
mixtool new grafana-dashboard > my-dashboard.jsonnet
mixtool new prometheus-alerts > my-alerts.jsonnet
mixtool new prometheus-rules > my-rules.jsonnet
```

### Lint

[embedmd]:# (_output/help-lint.txt)
```txt
NAME:
   mixtool lint - Lint jsonnet files

USAGE:
   mixtool lint [command options] [arguments...]

DESCRIPTION:
   Lint jsonnet files for correct structure of JSON objects

OPTIONS:
   --grafana                Lint Grafana dashboards against Grafana's schema
   --prometheus             Lint Prometheus alerts and rules and their given expressions
   --jpath value, -J value  Add folders to be used as vendor folders
   
```

#### Lint Examples

```bash
# This will lint the file for Prometheus alerts & rules and Grafana dashboards.
mixtool lint prometheus.jsonnet

# Don't lint Grafana dashboards.
mixtool lint --grafana=false prometheus.jsonnet

# Don't lint Prometheus alerts & rules.
mixtool lint --prometheus=false prometheus.jsonnet

# Lint multiple files sequentially.
mixtool lint prometheus.jsonnet grafana.jsonnet
```
