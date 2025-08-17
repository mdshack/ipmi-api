#!/bin/bash
set -e

echo "Building IPMI API binary..."

# Build the binary
go build -o ipmi-api main.go

echo "Build complete: ipmi-api"