#!/usr/bin/env bash
# Enforce release-candidate documentation coverage and reject stale status claims.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

REQUIRED_DOCS=(
  "doc/cli-config-contract.md"
  "doc/performance.md"
  "doc/regeneration.md"
  "doc/release-candidate.md"
  "doc/releasing.md"
  "doc/supply-chain-security.md"
  "doc/threat-model.md"
)

for relative_path in "${REQUIRED_DOCS[@]}"; do
  test -s "${REPO_ROOT}/${relative_path}"
  grep -q "$(basename "${relative_path}")" "${REPO_ROOT}/doc/README.md"
done

if grep -R -n --include='*.md' 'pre-alpha' "${REPO_ROOT}/README.md" "${REPO_ROOT}/doc"; then
  echo "Stale pre-alpha status remains in public documentation." >&2
  exit 1
fi

grep -q 'doc/release-candidate.md' "${REPO_ROOT}/README.md"
grep -q 'doc/threat-model.md' "${REPO_ROOT}/README.md"
grep -q 'v0.9.0 — Release Candidate Hardening' "${REPO_ROOT}/ROADMAP.md"

echo "Release-candidate documentation audit passed."
