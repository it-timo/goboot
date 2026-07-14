#!/usr/bin/env bash
# Verify snapshot release archives, checksums, SBOMs, and an installed Linux binary.

set -euo pipefail

DIST_ROOT="${1:-dist}"
TEMP_ROOT="$(mktemp -d)"

cleanup() {
  rm -rf "${TEMP_ROOT}"
}
trap cleanup EXIT

test -f "${DIST_ROOT}/checksums.txt"
(
  cd "${DIST_ROOT}"
  sha256sum --check checksums.txt
)

ARCHIVE_COUNT="$(find "${DIST_ROOT}" -maxdepth 1 -type f \( -name '*.tar.gz' -o -name '*.zip' \) | wc -l)"
SBOM_COUNT="$(find "${DIST_ROOT}" -maxdepth 1 -type f -name '*.sbom.json' | wc -l)"
test "${ARCHIVE_COUNT}" -ge 6
test "${SBOM_COUNT}" -ge 6

LINUX_ARCHIVE="$(find "${DIST_ROOT}" -maxdepth 1 -type f -name '*_linux_amd64.tar.gz' | head -n 1)"
test -n "${LINUX_ARCHIVE}"
tar -xzf "${LINUX_ARCHIVE}" -C "${TEMP_ROOT}"
test -x "${TEMP_ROOT}/goboot"
"${TEMP_ROOT}/goboot" --version | grep -q '^goboot '
"${TEMP_ROOT}/goboot" --help >/dev/null

echo "Release archives, checksums, SBOMs, and installation passed."
