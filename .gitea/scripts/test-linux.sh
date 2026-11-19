#!/usr/bin/env bash
set -euo pipefail

useradd --create-home --user-group skillshare-ci
# Git honors the checkout umask. Match the fixture producer's file modes on both runners.
chmod -R go-w "$GITHUB_WORKSPACE"
chown -R skillshare-ci:skillshare-ci "$GITHUB_WORKSPACE"
trap 'cat /sys/fs/cgroup/pids.events 2>/dev/null || true' EXIT
runuser -u skillshare-ci -- env PATH="$PATH" GOMAXPROCS=2 SKILLSHARE_TEST_BINARY="$GITHUB_WORKSPACE/bin/skillshare" \
    bash -e -s -- "${1:?expected unit or integration}" <<'CHECKS'
test "$(id -u)" -ne 0
umask 022
case "$1" in
    unit)
        make docs-check test-herdr build fmt-check lint
        python3 -B -m unittest discover -s scripts/release -p '*_test.py'
        python3 -B -m unittest discover -s scripts/design -p '*_test.py'
        go test -race -p 2 -parallel 2 ./cmd/skillshare/... ./internal/...
        ;;
    integration)
        make build
        go test -race -p 2 -parallel 2 ./tests/integration/...
        ;;
    *) exit 2 ;;
esac
CHECKS
