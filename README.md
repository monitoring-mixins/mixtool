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
   mixtool - Improves your jsonnet mixins workflows

USAGE:
   mixtool [global options] command [command options] [arguments...]

VERSION:
   v0.1.0-pre

DESCRIPTION:
   mixtool helps with generating, building and linting jsonnet mixins

COMMANDS:
     build     Build manifests from jsonnet input
     generate  Generate jsonnet mixin files
     lint      Lint jsonnet files
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
