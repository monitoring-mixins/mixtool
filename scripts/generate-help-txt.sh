#!/usr/bin/env bash

BINARY_NAME=mixtool
GOOS=$(uname -s | tr A-Z a-z)

$PWD/_output/$GOOS/amd64/$BINARY_NAME -h > $PWD/_output/help.txt
$PWD/_output/$GOOS/amd64/$BINARY_NAME build -h > $PWD/_output/help-build.txt
$PWD/_output/$GOOS/amd64/$BINARY_NAME lint -h > $PWD/_output/help-lint.txt
$PWD/_output/$GOOS/amd64/$BINARY_NAME new -h > $PWD/_output/help-new.txt
$PWD/_output/$GOOS/amd64/$BINARY_NAME runbook -h > $PWD/_output/help-runbook.txt
