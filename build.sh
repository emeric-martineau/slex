#!/bin/env bash
go test ./...
cp cmd/command_line_options.go cmd/command_line_options.go.old
sed -e '/@@@@data@@@@/{r lexer/lexer.go' -e 'd}' -i cmd/command_line_options.go
go build .
cp cmd/command_line_options.go.old cmd/command_line_options.go
rm -f cp cmd/command_line_options.go.old
