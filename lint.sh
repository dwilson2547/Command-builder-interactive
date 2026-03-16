#!/usr/bin/env bash
set -euo pipefail

LINT_TIMEOUT="5m"
GOLANGCI_LINT_CMD="golangci-lint"

# Check if Go is installed
if ! command -v go &>/dev/null; then
  echo "Error: go is not installed or not in PATH." >&2
  exit 1
fi

echo "Using Go $(go version | awk '{print $3}')"

# Install golangci-lint if not present
if ! command -v "$GOLANGCI_LINT_CMD" &>/dev/null; then
  echo "golangci-lint not found — installing latest..."
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  # Ensure GOPATH/bin is on PATH for the rest of the script
  export PATH="$(go env GOPATH)/bin:$PATH"
  if ! command -v "$GOLANGCI_LINT_CMD" &>/dev/null; then
    echo "Error: golangci-lint installation failed." >&2
    exit 1
  fi
fi

echo "Using $(golangci-lint --version)"
echo "Running golangci-lint..."

golangci-lint run --timeout="$LINT_TIMEOUT" "$@"

echo "Lint passed."
