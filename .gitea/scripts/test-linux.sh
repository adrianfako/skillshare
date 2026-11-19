#!/usr/bin/env bash
set -euo pipefail

useradd --create-home --user-group skillshare-ci
chown -R skillshare-ci:skillshare-ci "$GITHUB_WORKSPACE"
runuser -u skillshare-ci -- env PATH="$PATH" SKILLSHARE_TEST_BINARY="$GITHUB_WORKSPACE/bin/skillshare" \
    bash -e -s -- "${1:?expected unit or integration}" <<'CHECKS'
test "$(id -u)" -ne 0
case "$1" in
    unit)
        make docs-check test-herdr build fmt-check lint
        python3 -B -m unittest discover -s scripts/release -p '*_test.py'
        python3 -B -m unittest discover -s scripts/design -p '*_test.py'
        go test -race ./cmd/skillshare/... ./internal/...
        ;;
    integration)
        make build
        go test -race ./tests/integration/...
        ;;
    *) exit 2 ;;
esac
CHECKS
