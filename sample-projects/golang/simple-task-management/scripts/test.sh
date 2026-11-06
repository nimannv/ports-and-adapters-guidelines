#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR"

echo "=== Running unit tests ==="
go test -v ./internal/core/... ./internal/adapters/inbound/cli/...

echo "=== Running integration tests ==="
go test -v -tags=integration ./test/integration/...

