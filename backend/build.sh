#!/usr/bin/env bash

set -euo pipefail

echo ">> Installing Go dependencies"
go mod download

echo ">> Building GrowTrade backend binary"
go build -o growtrade-server main.go

echo "Build complete: $(pwd)/growtrade-server"

