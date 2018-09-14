#!/usr/bin/env bash

BINARY_NAME=mixtool
GOOS=$(uname -s | tr A-Z a-z)

HELP_FILE=$PWD/_output/help.txt
echo "$ $BINARY_NAME -h" > $HELP_FILE
PATH=$PATH:$PWD/_output/$GOOS/amd64 $BINARY_NAME -h &> $HELP_FILE
exit 0
