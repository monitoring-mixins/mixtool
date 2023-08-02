#!/usr/bin/env bash

BINARY_NAME=mixtool
GOOS=$(uname -s | tr A-Z a-z)
GOARCH=$(uname -m | sed 's/amd64/x86_64/' | sed 's/i.86/386/')

$PWD/_output/$GOOS/$GOARCH/$BINARY_NAME -h > $PWD/_output/help.txt
$PWD/_output/$GOOS/$GOARCH/$BINARY_NAME generate -h > $PWD/_output/help-generate.txt
$PWD/_output/$GOOS/$GOARCH/$BINARY_NAME lint -h > $PWD/_output/help-lint.txt
$PWD/_output/$GOOS/$GOARCH/$BINARY_NAME new -h > $PWD/_output/help-new.txt
# $PWD/_output/$GOOS/amd64/$BINARY_NAME runbook -h > $PWD/_output/help-runbook.txt
