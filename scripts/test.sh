#!/usr/bin/env bash
# Run repository tests with race detection and coverage.
# Usage: ./scripts/test.sh

set -euo pipefail

COVER_FILE="coverage.txt"

echo "Running go test..."

TEST_PKGS="$(go list ./... | grep -v "/test/noauto" | grep -v "/templates")"

# shellcheck disable=SC2086 # Intended word splitting
go test -race -timeout=5m -coverprofile="${COVER_FILE}" ${TEST_PKGS}

go tool cover -func="${COVER_FILE}"
rm -f "${COVER_FILE}"

echo "go test passed"

echo ""
echo "All tests completed successfully!"
