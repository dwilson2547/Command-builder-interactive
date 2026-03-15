#!/bin/bash
set -e

VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v1.1.0")
LDFLAGS="-X github.com/dwilson2547/command-builder/internal/tui.AppVersion=${VERSION}"

go build -ldflags "${LDFLAGS}" -o command-builder .
echo "Built command-builder ${VERSION}"
