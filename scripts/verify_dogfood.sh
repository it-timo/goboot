#!/usr/bin/env bash
# Generate goboot from its committed configuration twice and require byte-identical output.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BINARY="${1:-${REPO_ROOT}/bin/goboot}"
TEMP_ROOT="$(mktemp -d)"

cleanup() {
  rm -rf "${TEMP_ROOT}"
}
trap cleanup EXIT

if [[ "${BINARY}" != /* ]]; then
  BINARY="$(cd "$(dirname "${BINARY}")" && pwd)/$(basename "${BINARY}")"
fi

CONFIG_FILE="${TEMP_ROOT}/goboot.yml"
TARGET_ROOT="${TEMP_ROOT}/output"

awk -v target="${TARGET_ROOT}" '
  $1 == "targetPath:" {
    print "targetPath: \"" target "\""
    next
  }
  { print }
' "${REPO_ROOT}/configs/goboot.yml" >"${CONFIG_FILE}"

cd "${REPO_ROOT}"

"${BINARY}" --config "${CONFIG_FILE}" --validate --output json >"${TEMP_ROOT}/validation.json"
grep -q '"status":"success"' "${TEMP_ROOT}/validation.json"

"${BINARY}" --config "${CONFIG_FILE}" --skip-go-mod-tidy --output json >"${TEMP_ROOT}/first.json"
find "${TARGET_ROOT}" -type f -print0 | sort -z | xargs -0 sha256sum >"${TEMP_ROOT}/first.sha256"

"${BINARY}" --config "${CONFIG_FILE}" --skip-go-mod-tidy --output json >"${TEMP_ROOT}/second.json"
find "${TARGET_ROOT}" -type f -print0 | sort -z | xargs -0 sha256sum >"${TEMP_ROOT}/second.sha256"

diff -u "${TEMP_ROOT}/first.sha256" "${TEMP_ROOT}/second.sha256"
grep -q '"status":"success"' "${TEMP_ROOT}/second.json"

echo "Dogfood generation is deterministic."
