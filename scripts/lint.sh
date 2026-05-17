#!/usr/bin/env bash
# Run repository linters with pinned Docker images.
# Usage: ./scripts/lint.sh

set -euo pipefail

# Load centralized tool versions.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
# shellcheck source=versions.env
source "${PROJECT_ROOT}/versions.env"

DOCKER_CMD=(docker run --rm --network host -v "$(pwd):/workdir" -w /workdir)

GOLANGCI_LINT_IMAGE="golangci/golangci-lint:${GOLANGCI_LINT_VERSION}"
MD_LINT_IMAGE="ghcr.io/igorshubovych/markdownlint-cli:${MARKDOWNLINT_VERSION}"
YAMLLINT_IMAGE="pipelinecomponents/yamllint:${YAMLLINT_VERSION}"
CHECKMAKE_IMAGE="cytopia/checkmake:${CHECKMAKE_VERSION}"
SHELLCHECK_IMAGE="koalaman/shellcheck:${SHELLCHECK_VERSION}"
SHFMT_IMAGE="mvdan/shfmt:${SHFMT_VERSION}"
EDITORCONFIG_CHECKER_IMAGE="mstruebing/editorconfig-checker:${EDITORCONFIG_CHECKER_VERSION}"

if ! command -v docker &>/dev/null; then
  echo "Error: Docker is required but not installed."
  exit 1
fi

echo "Running golangci-lint..."
"${DOCKER_CMD[@]}" "${GOLANGCI_LINT_IMAGE}" golangci-lint run ./...
echo "golangci-lint passed"

echo "Running yamllint..."
"${DOCKER_CMD[@]}" "${YAMLLINT_IMAGE}" yamllint .
echo "yamllint passed"

echo "Running checkmake..."
"${DOCKER_CMD[@]}" "${CHECKMAKE_IMAGE}" Makefile
echo "checkmake passed"

echo "Running markdownlint..."

# shellcheck disable=SC2312 # find failures are still surfaced by set -e in this script context.
mapfile -d '' -t MARKDOWN_FILES < <(
  find . \
    -path './.git' -prune -o \
    -path './.idea' -prune -o \
    -path './.gitlab-ci-local' -prune -o \
    -path './bin' -prune -o \
    -path './outputs' -prune -o \
    -type f -name "*.md" -print0
)

if ((${#MARKDOWN_FILES[@]} == 0)); then
  echo "No Markdown files found to lint."
else
  "${DOCKER_CMD[@]}" "${MD_LINT_IMAGE}" markdownlint "${MARKDOWN_FILES[@]}"
  echo "markdownlint passed"
fi

echo "Running shellcheck..."

# shellcheck disable=SC2312 # find failures are still surfaced by set -e in this script context.
mapfile -d '' -t SH_FILES < <(
  find . \
    -path './.git' -prune -o \
    -path './.idea' -prune -o \
    -path './.gitlab-ci-local' -prune -o \
    -path './bin' -prune -o \
    -path './outputs' -prune -o \
    -type f -name "*.sh" -print0
)

if ((${#SH_FILES[@]} == 0)); then
  echo "No shell scripts found to lint."
else
  "${DOCKER_CMD[@]}" "${SHELLCHECK_IMAGE}" -x "${SH_FILES[@]}"
  echo "shellcheck passed"
fi

echo "Running shfmt (check only)..."

if ((${#SH_FILES[@]} == 0)); then
  echo "No shell scripts found to format-check."
else
  "${DOCKER_CMD[@]}" "${SHFMT_IMAGE}" -d -i 2 -ci "${SH_FILES[@]}"
  echo "shfmt check passed"
fi

echo "Running editorconfig-checker..."
"${DOCKER_CMD[@]}" --entrypoint ec "${EDITORCONFIG_CHECKER_IMAGE}" -exclude '(\.git|\.idea|\.gitlab-ci-local|bin)'
echo "editorconfig-checker passed"

echo ""
echo "All linters completed successfully!"
