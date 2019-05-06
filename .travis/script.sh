#!/usr/bin/env bash

function main {
    mkdir -p target
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o target/envly-darwin-amd64
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o target/envly-linux-amd64
}

main
