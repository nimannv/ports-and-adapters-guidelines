#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR"

if [ -f "$DIR/.venv/bin/lint-imports" ]; then
    LINTER="$DIR/.venv/bin/lint-imports"
else
    LINTER="lint-imports"
fi

echo "Verifying architectural contracts with import-linter..."
exec "$LINTER" "$@"

