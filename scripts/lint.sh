#!/usr/bin/env bash
# ==============================================================================
# Lint Script — Comprehensive Project Linter
# ==============================================================================
#
# This script performs full-project linting using the following tools:
#
#   1. golangci-lint   → Static analysis for all Go source files
#   2. yamllint         → Syntax + style checking for YAML configuration files
#   3. checkmake        → Makefile validation and style guidance
#   4. markdownlint     → Markdown file linting via pinned Docker container
#
# Tools must be installed (except for markdownlint, which runs via Docker).
#
# USAGE:
#   ./scripts/lint.sh
#
# EXIT CODES:
#   0 → All linters passed
#   1 → One or more linters failed
#
# NOTE:
#   This script is standalone and **does not depend** on Makefile or Taskfile.
# ==============================================================================

set -euo pipefail

# ------------------------------------------------------------------------------
# Config
# ------------------------------------------------------------------------------

MARKDOWNLINT_IMAGE="ghcr.io/igorshubovych/markdownlint-cli:v0.45.0"

# ------------------------------------------------------------------------------
# Lint: Go (golangci-lint)
# ------------------------------------------------------------------------------

echo "Running golangci-lint..."
if ! command -v golangci-lint &>/dev/null; then
  echo "golangci-lint not found. Please install it: https://golangci-lint.run/usage/install/"
  exit 1
fi

golangci-lint run ./...
echo "golangci-lint passed"

# ------------------------------------------------------------------------------
# Lint: YAML (yamllint)
# ------------------------------------------------------------------------------

echo "Running yamllint..."
if ! command -v yamllint &>/dev/null; then
  echo "yamllint not found. Please install it: https://yamllint.readthedocs.io/en/stable/"
  exit 1
fi

yamllint .
echo "yamllint passed"

# ------------------------------------------------------------------------------
# Lint: Makefile (checkmake)
# ------------------------------------------------------------------------------

echo "Running checkmake..."
if ! command -v checkmake &>/dev/null; then
  echo "checkmake not found. Please install it: https://github.com/mrtazz/checkmake"
  exit 1
fi

checkmake Makefile
echo "checkmake passed"

# ------------------------------------------------------------------------------
# Lint: Markdown (markdownlint via Docker)
# ------------------------------------------------------------------------------

echo "Running markdownlint (via Docker)..."
if ! command -v docker &>/dev/null; then
  echo "Docker is required to run markdownlint (not found)"
  exit 1
fi

# Recursively find all Markdown files (*.md) from project root
MARKDOWN_FILES=$(find . -type f -name "*.md")

if [[ -z "$MARKDOWN_FILES" ]]; then
  echo "No Markdown files found to lint."
else
  docker run --rm \
    -v "$PWD":/workdir \
    -w /workdir \
    "$MARKDOWNLINT_IMAGE" \
    markdownlint $MARKDOWN_FILES
  echo "markdownlint passed"
fi

# ------------------------------------------------------------------------------
# All Linters Passed
# ------------------------------------------------------------------------------

echo ""
echo "All linters completed successfully!"
