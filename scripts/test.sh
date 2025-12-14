#!/usr/bin/env bash
# ==============================================================================
# Test Script — Comprehensive Project Test Runner
# ==============================================================================
#
# This script runs all Go tests with race detection and coverage reporting.
#
# USAGE:
#   ./scripts/test.sh
#
# EXIT CODES:
#   0 → All tests passed
#   1 → One or more tests failed
# ==============================================================================

set -euo pipefail

# ------------------------------------------------------------------------------
# Config
# ------------------------------------------------------------------------------

COVER_FILE="coverage.txt"

# ------------------------------------------------------------------------------
# Test: Go (go test)
# ------------------------------------------------------------------------------

echo "Running go test..."

TEST_PKGS="$(go list ./... | grep -v "/test/noauto" | grep -v "/templates")"

# shellcheck disable=SC2086 # Intended word splitting
go test -race -timeout=5m -coverprofile="${COVER_FILE}" ${TEST_PKGS}

go tool cover -func="${COVER_FILE}"
rm -f "${COVER_FILE}"

echo "go test passed"

# ------------------------------------------------------------------------------
# All Tests Passed
# ------------------------------------------------------------------------------

echo ""
echo "All tests completed successfully!"
