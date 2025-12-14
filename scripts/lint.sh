#!/usr/bin/env bash
# ==============================================================================
# Lint Script — Comprehensive Project Linter
# ==============================================================================
#
# This script performs full-project linting using the following tools (via Docker):
#
#   1. golangci-lint   → Static analysis for all Go source files
#   2. yamllint         → Syntax + style checking for YAML configuration files
#   3. checkmake        → Makefile validation and style guidance
#   4. markdownlint     → Markdown file linting
#   5. shellcheck       → Shell script linting
#   6. shfmt            → Shell script formatting check
#
# All tools run via Docker to ensure consistency and correct versions.
#
# USAGE:
#   ./scripts/lint.sh
#
# EXIT CODES:
#   0 → All linters passed
#   1 → One or more linters failed
# ==============================================================================

set -euo pipefail

# ------------------------------------------------------------------------------
# Config
# ------------------------------------------------------------------------------

DOCKER_CMD="docker run --rm -v $(pwd):/workdir -w /workdir"

GOLANGCI_LINT_IMAGE="golangci/golangci-lint:v2.7.1"
MD_LINT_IMAGE="ghcr.io/igorshubovych/markdownlint-cli:v0.46.0"
YAMLLINT_IMAGE="pipelinecomponents/yamllint:0.35.9"
CHECKMAKE_IMAGE="cytopia/checkmake:latest-0.5"
SHELLCHECK_IMAGE="koalaman/shellcheck:v0.11.0"
SHFMT_IMAGE="mvdan/shfmt:v3.12.0"

# ------------------------------------------------------------------------------
# Pre-flight Check
# ------------------------------------------------------------------------------

if ! command -v docker &>/dev/null; then
  echo "Error: Docker is required but not installed."
  exit 1
fi

# ------------------------------------------------------------------------------
# Lint: Go (golangci-lint)
# ------------------------------------------------------------------------------

echo "Running golangci-lint..."
${DOCKER_CMD} "${GOLANGCI_LINT_IMAGE}" golangci-lint run cmd/... pkg/...
echo "golangci-lint passed"

# ------------------------------------------------------------------------------
# Lint: YAML (yamllint)
# ------------------------------------------------------------------------------

echo "Running yamllint..."
${DOCKER_CMD} "${YAMLLINT_IMAGE}" yamllint .
echo "yamllint passed"

# ------------------------------------------------------------------------------
# Lint: Makefile (checkmake)
# ------------------------------------------------------------------------------

echo "Running checkmake..."
${DOCKER_CMD} "${CHECKMAKE_IMAGE}" Makefile
echo "checkmake passed"

# ------------------------------------------------------------------------------
# Lint: Markdown (markdownlint)
# ------------------------------------------------------------------------------

echo "Running markdownlint..."

MARKDOWN_FILES="$(find . -type f -name "*.md")"

if [[ -z "${MARKDOWN_FILES}" ]]; then
  echo "No Markdown files found to lint."
else
  # shellcheck disable=SC2086
  ${DOCKER_CMD} "${MD_LINT_IMAGE}" markdownlint ${MARKDOWN_FILES}
  echo "markdownlint passed"
fi

# ------------------------------------------------------------------------------
# Lint: Shell scripts (ShellCheck)
# ------------------------------------------------------------------------------

echo "Running shellcheck..."

SH_FILES="$(find . -type f -name "*.sh")"

if [[ -z "${SH_FILES}" ]]; then
  echo "No shell scripts found to lint."
else
  # shellcheck disable=SC2086
  ${DOCKER_CMD} "${SHELLCHECK_IMAGE}" -x ${SH_FILES}
  echo "shellcheck passed"
fi

# ------------------------------------------------------------------------------
# Format Check: Shell scripts (shfmt, check-only)
# ------------------------------------------------------------------------------

echo "Running shfmt (check only)..."

if [[ -z "${SH_FILES}" ]]; then
  echo "No shell scripts found to format-check."
else
  # shellcheck disable=SC2086
  ${DOCKER_CMD} "${SHFMT_IMAGE}" -d -i 2 -ci ${SH_FILES}
  echo "shfmt check passed"
fi

# ------------------------------------------------------------------------------
# All Linters Passed
# ------------------------------------------------------------------------------

echo ""
echo "All linters completed successfully!"
