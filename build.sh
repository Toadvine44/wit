#!/bin/bash

set -e

cd "$(dirname "$0")"

echo "Building wit..."
go build -o wit .

echo "Build complete: ./wit"
