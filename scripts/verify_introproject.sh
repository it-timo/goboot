#!/usr/bin/env bash
# Regenerate IntroProject and run the full root + generated verification flow.
# Usage:
#   ./scripts/verify_introproject.sh
#   ./scripts/verify_introproject.sh --skip-task

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
SKIP_TASK="false"

for arg in "$@"; do
  case "${arg}" in
    --skip-task)
      SKIP_TASK="true"
      ;;
    *)
      echo "Unknown argument: ${arg}"
      echo "Usage: ./scripts/verify_introproject.sh [--skip-task]"
      exit 1
      ;;
  esac
done

run_step() {
  local title="$1"
  local cmd="$2"

  echo ""
  echo "==> ${title}"
  echo "    ${cmd}"
  bash -lc "${cmd}"
}

if ! command -v go >/dev/null 2>&1; then
  echo "Error: go is required but not installed."
  exit 1
fi

if ! command -v make >/dev/null 2>&1; then
  echo "Error: make is required but not installed."
  exit 1
fi

if [[ "${SKIP_TASK}" != "true" ]] && ! command -v task >/dev/null 2>&1; then
  echo "Error: task is required for this flow. Re-run with --skip-task to skip task commands."
  exit 1
fi

cd "${PROJECT_ROOT}"

if [[ -d "${PROJECT_ROOT}/outputs/IntroProject" ]]; then
  run_step "Reset existing outputs/IntroProject" "rm -rf '${PROJECT_ROOT}/outputs/IntroProject'"
fi

run_step "Generate IntroProject from default config" "go run cmd/goboot/main.go"

run_step "Root: make test" "make test"
run_step "Root: make lint" "make lint"

if [[ "${SKIP_TASK}" != "true" ]]; then
  run_step "Root: task test" "task test"
  run_step "Root: task lint" "task lint"
fi

run_step "Root: ./scripts/test.sh" "./scripts/test.sh"
run_step "Root: ./scripts/lint.sh" "./scripts/lint.sh"

cd "${PROJECT_ROOT}/outputs/IntroProject"

run_step "IntroProject: make test" "make test"
run_step "IntroProject: make lint" "make lint"

if [[ "${SKIP_TASK}" != "true" ]]; then
  run_step "IntroProject: task test" "task test"
  run_step "IntroProject: task lint" "task lint"
fi

run_step "IntroProject: ./scripts/test.sh" "./scripts/test.sh"
run_step "IntroProject: ./scripts/lint.sh" "./scripts/lint.sh"

echo ""
echo "Full verification workflow completed successfully."
