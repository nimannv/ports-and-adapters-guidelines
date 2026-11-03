#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
export PYTHONPATH="$DIR"

if [ -f "$DIR/.venv/bin/pytest" ]; then
    PYTEST="$DIR/.venv/bin/pytest"
else
    PYTEST="pytest"
fi

echo "Running test suite (Unit & Integration)..."
exec "$PYTEST" "$DIR/tests" "$@"

