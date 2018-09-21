# mixtool

> NOTE: This project is *alpha* stage. Flags, configuration, behavior and design may change significantly in following releases.

The mixtool is a helper for easily working with [jsonnet](http://jsonnet.org/) mixins.

## Install

```
go get -u github.com/metalmatze/mixtool/cmd/mixtool
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
     build    Build manifests from jsonnet input
     lint     Lint jsonnet files
     new      Create new jsonnet mixin files
     runbook  Generate a runbook markdown file
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Build

[embedmd]:# (_output/help-build.txt)
```txt
NAME:
   mixtool build - Build manifests from jsonnet input

USAGE:
   mixtool build [command options] [arguments...]

OPTIONS:
   --jpath value, -J value        
   --multi value, -m value        
   --output-file value, -o value  
   --yaml, -y                     
   
```

#### Build Examples

```bash
# The simplest example. It will use ./vendor if available and print YAML.
mixtool build prometheus.jsonnet

# These next lines are equivalent and both write to file.
mixtool build prometheus.jsonnet > prometheus.yaml
mixtool build --output-file prometheus.yaml prometheus.jsonnet

# Change the folder for imports.
mixtool build --jpath /some/path/vendor prometheus.jsonnet

# Instead of writing YAML this will simply output JSON.
mixtool build --yaml=false prometheus.jsonnet > prometheus.json

# Write the output to multiple files in a give directory.
mixtool build --multi prometheus/ prometheus.jsonnet
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

### Runbook

[embedmd]:# (_output/help-runbook.txt)
```txt
NAME:
   mixtool runbook - Generate a runbook markdown file

USAGE:
   mixtool runbook [command options] [arguments...]

DESCRIPTION:
   Generate a runbook markdown file from the jsonnet mixins

OPTIONS:
   --jpath value, -J value        
   --output-file value, -o value  
   
```

#### Runbook Examples

```bash
# The simplest example. It will use ./vendor if available and write the runbook as markdown to stdout.
mixtool runbook alerts.libsonnet

# These next lines are equivalent and both write to file.
mixtool runbook alerts.libsonnet > runbook.md
mixtool runbook -o runbook.md alerts.libsonnet

# Change the folder for imports.
mixtool runbook --jpath /some/path/vendor alerts.libsonnet
```
