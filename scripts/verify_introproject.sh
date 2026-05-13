#!/usr/bin/env bash
# Regenerate representative Intro* projects and run root + generated checks.
# Usage:
#   ./scripts/verify_introproject.sh
#   ./scripts/verify_introproject.sh --skip-task

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
SKIP_TASK="false"
TEMP_DIR=""

PROJECT_CASES=(
  "IntroProject:ginkgo:slog:gitlab"
  "IntroGoSlogGitLab:go:slog:gitlab"
  "IntroGinkgoZerologGitLab:ginkgo:zerolog:gitlab"
  "IntroGoZerologGitLab:go:zerolog:gitlab"
)

cleanup() {
  if [[ -n "${TEMP_DIR}" ]]; then
    rm -rf "${TEMP_DIR}"
  fi
}

trap cleanup EXIT

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

write_base_test_config() {
  local path="$1"
  local test_style="$2"

  if [[ "${test_style}" == "ginkgo" ]]; then
    cat >"${path}" <<YAML
sourcePath: "${PROJECT_ROOT}/templates/test_base"
useStyle: "ginkgo"
testCmd: |-
  go test -race -timeout=5m -coverprofile=coverage.txt ./... &&
  go tool cover -func=coverage.txt &&
  rm -f coverage.txt
YAML

    return 0
  fi

  cat >"${path}" <<YAML
sourcePath: "${PROJECT_ROOT}/templates/test_base"
useStyle: "go"
testCmd: |-
  go test -race -timeout=5m -coverprofile=coverage.txt ./... &&
  go tool cover -func=coverage.txt &&
  rm -f coverage.txt
YAML
}

write_base_logger_config() {
  local path="$1"
  local logger_type="$2"

  cat >"${path}" <<YAML
loggerType: "${logger_type}"
YAML
}

write_goboot_config() {
  local path="$1"
  local project_name="$2"
  local git_provider="$3"
  local base_test_cfg="$4"
  local base_logger_cfg="$5"

  cat >"${path}" <<YAML
targetPath: "outputs"
projectName: "${project_name}"
repoUrl: "https://github.com/projects"
gitProvider: "${git_provider}"
services:
  - id: "base_project"
    confPath: "${PROJECT_ROOT}/configs/base_project.yml"
    enabled: true
  - id: "base_lint"
    confPath: "${PROJECT_ROOT}/configs/base_lint.yml"
    enabled: true
  - id: "base_test"
    confPath: "${base_test_cfg}"
    enabled: true
  - id: "base_logger"
    confPath: "${base_logger_cfg}"
    enabled: true
  - id: "base_local"
    confPath: "${PROJECT_ROOT}/configs/base_local.yml"
    enabled: true
  - id: "base_ci"
    confPath: "${PROJECT_ROOT}/configs/base_ci.yml"
    enabled: true
YAML
}

generate_project_case() {
  local case_spec="$1"
  local project_name
  local test_style
  local logger_type
  local git_provider
  local case_dir
  local base_test_cfg
  local base_logger_cfg
  local goboot_cfg

  IFS=":" read -r project_name test_style logger_type git_provider <<<"${case_spec}"

  case_dir="${TEMP_DIR}/${project_name}"
  mkdir -p "${case_dir}"

  base_test_cfg="${case_dir}/base_test.yml"
  base_logger_cfg="${case_dir}/base_logger.yml"
  goboot_cfg="${case_dir}/goboot.yml"

  write_base_test_config "${base_test_cfg}" "${test_style}"
  write_base_logger_config "${base_logger_cfg}" "${logger_type}"
  write_goboot_config "${goboot_cfg}" "${project_name}" "${git_provider}" "${base_test_cfg}" "${base_logger_cfg}"

  if [[ -d "${PROJECT_ROOT}/outputs/${project_name}" ]]; then
    run_step "Reset existing outputs/${project_name}" "rm -rf '${PROJECT_ROOT}/outputs/${project_name}'"
  fi

  run_step "Generate ${project_name} (${test_style}, ${logger_type}, ${git_provider})" "go run cmd/goboot/main.go --config '${goboot_cfg}'"
}

validate_project_case() {
  local case_spec="$1"
  local project_name
  local test_style
  local logger_type
  local git_provider
  local project_root

  IFS=":" read -r project_name test_style logger_type git_provider <<<"${case_spec}"
  project_root="${PROJECT_ROOT}/outputs/${project_name}"

  cd "${project_root}"

  run_step "${project_name}: make test" "make test"
  run_step "${project_name}: make lint" "make lint"

  if [[ "${SKIP_TASK}" != "true" ]]; then
    run_step "${project_name}: task test" "task test"
    run_step "${project_name}: task lint" "task lint"
  fi

  run_step "${project_name}: ./scripts/test.sh" "./scripts/test.sh"
  run_step "${project_name}: ./scripts/lint.sh" "./scripts/lint.sh"

  if [[ "${git_provider}" == "gitlab" ]]; then
    run_step "${project_name}: GitLab CI files exist" "test -f .gitlab-ci.yml && test -f .gitlab/ci/lint.yml && test -f .gitlab/ci/test.yml && test -f .gitlab/ci/build.yml"
  fi

  if [[ "${test_style}" == "go" ]]; then
    run_step "${project_name}: go-style tests generated" "! find . -name '*_suite_test.go' | grep -q ."
  else
    run_step "${project_name}: ginkgo suites generated" "find . -name '*_suite_test.go' | grep -q ."
  fi

  case "${logger_type}" in
    slog)
      run_step "${project_name}: slog logger generated" "grep -R 'log/slog' cmd pkg >/dev/null && ! grep -R 'github.com/rs/zerolog' cmd pkg >/dev/null"
      ;;
    zerolog)
      run_step "${project_name}: zerolog logger generated" "grep -R 'github.com/rs/zerolog' cmd pkg >/dev/null"
      ;;
    *)
      echo "Error: unknown logger type ${logger_type}"
      exit 1
      ;;
  esac

  cd "${PROJECT_ROOT}"
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
TEMP_DIR="$(mktemp -d)"

for case_spec in "${PROJECT_CASES[@]}"; do
  generate_project_case "${case_spec}"
done

run_step "Root: make test" "make test"
run_step "Root: make lint" "make lint"

if [[ "${SKIP_TASK}" != "true" ]]; then
  run_step "Root: task test" "task test"
  run_step "Root: task lint" "task lint"
fi

run_step "Root: ./scripts/test.sh" "./scripts/test.sh"
run_step "Root: ./scripts/lint.sh" "./scripts/lint.sh"

for case_spec in "${PROJECT_CASES[@]}"; do
  validate_project_case "${case_spec}"
done

echo ""
echo "Full verification workflow completed successfully (${#PROJECT_CASES[@]} generated project cases)."
