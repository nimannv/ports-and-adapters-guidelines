#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR"

export CONFIGS_DIR="$DIR/configs"
go run ./bootstrap/cli/main.go "$@"

