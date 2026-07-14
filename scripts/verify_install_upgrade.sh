#!/usr/bin/env bash
# Replace an installed baseline binary with the current candidate and smoke-test it.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BASELINE_REF="${1:-HEAD^}"
TEMP_ROOT="$(mktemp -d)"
BASELINE_ROOT="${TEMP_ROOT}/baseline"
INSTALL_ROOT="${TEMP_ROOT}/install"

cleanup() {
  git -C "${REPO_ROOT}" worktree remove --force "${BASELINE_ROOT}" >/dev/null 2>&1 || true
  rm -rf "${TEMP_ROOT}"
}
trap cleanup EXIT

mkdir -p "${INSTALL_ROOT}"
git -C "${REPO_ROOT}" worktree add --detach "${BASELINE_ROOT}" "${BASELINE_REF}" >/dev/null

go -C "${BASELINE_ROOT}" build -ldflags="-X main.version=baseline" -o "${INSTALL_ROOT}/goboot" ./cmd/goboot
test -x "${INSTALL_ROOT}/goboot"
BASELINE_DIGEST="$(sha256sum "${INSTALL_ROOT}/goboot" | cut -d ' ' -f 1)"

go -C "${REPO_ROOT}" build -ldflags="-X main.version=v0.9.0-rc" -o "${TEMP_ROOT}/goboot" ./cmd/goboot
install -m 0755 "${TEMP_ROOT}/goboot" "${INSTALL_ROOT}/goboot"
CANDIDATE_DIGEST="$(sha256sum "${INSTALL_ROOT}/goboot" | cut -d ' ' -f 1)"

test "${BASELINE_DIGEST}" != "${CANDIDATE_DIGEST}"
"${INSTALL_ROOT}/goboot" --version | grep -q '^goboot v0.9.0-rc$'
"${INSTALL_ROOT}/goboot" --help >/dev/null

cd "${REPO_ROOT}"
"${INSTALL_ROOT}/goboot" --config ./configs/goboot.yml --validate --output json \
  | grep -q '"status":"success"'

echo "Baseline-to-candidate installation upgrade passed."
