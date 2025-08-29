#!/usr/bin/env sh

echo "[ ! -f ${PWD}/.env ] || export \$(grep -v '^#' ${PWD}/.env | xargs)" >>/home/vscode/.bashrc
go install github.com/joerdav/xc/cmd/xc@latest && xc -complete
go mod tidy
