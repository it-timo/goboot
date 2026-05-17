#!/usr/bin/env bash
# Enforce total Go statement coverage from a coverprofile.
# Usage: ./scripts/check_coverage.sh coverage.txt 80

set -euo pipefail

COVER_FILE="${1:-coverage.txt}"
MIN_COVERAGE="${2:-80}"

if [[ ! -f "${COVER_FILE}" ]]; then
  echo "Coverage file not found: ${COVER_FILE}" >&2
  exit 1
fi

total_coverage="$(
  go tool cover -func="${COVER_FILE}" |
    awk '/^total:/ { gsub("%", "", $3); print $3 }'
)"

if [[ -z "${total_coverage}" ]]; then
  echo "Could not read total coverage from ${COVER_FILE}" >&2
  exit 1
fi

awk -v actual="${total_coverage}" -v minimum="${MIN_COVERAGE}" '
  BEGIN {
    if (actual + 0 < minimum + 0) {
      printf "coverage %.1f%% is below required %.1f%%\n", actual, minimum > "/dev/stderr"
      exit 1
    }

    printf "coverage %.1f%% meets required %.1f%%\n", actual, minimum
  }
'
