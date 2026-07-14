#!/usr/bin/env bash
# Run local CI simulation for goboot root and generated projects.
#
# This script executes the exact act/gitlab-ci-local commands used to validate
# generated CI behavior for both providers.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
OUTPUT_DIR="${PROJECT_ROOT}/outputs/IntroProject"
ARTIFACT_DIR="/tmp/act-artifacts"
TOOLCACHE_DIR="/tmp/act-toolcache"
ACT_PULL="false"
PREPULL_IMAGES="true"
PROVIDER_MODE="${CI_CANARY_PROVIDER:-config}"
GITLAB_CI_LOCAL_OPTS="--network host"
GITLAB_CI_LOCAL_SHELL_OPTS="${GITLAB_CI_LOCAL_OPTS} --force-shell-executor --concurrency 1"
GITHUB_CANARY_TEMP_DIR=""
GITHUB_CANARY_OUTPUT_DIR=""
TEMP_FILES=()

cleanup_temp_files() {
  local file
  for file in "${TEMP_FILES[@]}"; do
    rm -f "${file}"
  done

  if [[ -n "${GITHUB_CANARY_TEMP_DIR}" ]]; then
    rm -rf "${GITHUB_CANARY_TEMP_DIR}"
  fi
}

trap cleanup_temp_files EXIT

usage() {
  cat <<USAGE
Usage: ./scripts/verify_ci_canary.sh [options]

Options:
  --skip-generate  Skip regenerating outputs/IntroProject via goboot
  --refresh-images Pull fresh act runner images (less deterministic)
  --skip-prepull  Skip image pre-pull for both GitHub (act) and GitLab flows
  --provider       CI provider mode: github | gitlab | both | config (default)
  --help           Show this help
USAGE
}

SKIP_GENERATE="false"

for arg in "$@"; do
  case "${arg}" in
    --skip-generate)
      SKIP_GENERATE="true"
      ;;
    --refresh-images)
      ACT_PULL="true"
      ;;
    --skip-prepull)
      PREPULL_IMAGES="false"
      ;;
    --provider=*)
      PROVIDER_MODE="${arg#*=}"
      ;;
    --provider)
      echo "Error: --provider requires a value. Use --provider=<github|gitlab|both|config>."
      exit 1
      ;;
    --help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: ${arg}"
      usage
      exit 1
      ;;
  esac
done

resolve_provider_mode() {
  local mode="${PROVIDER_MODE}"
  local cfg_provider

  case "${mode}" in
    github | gitlab | both)
      echo "${mode}"
      return 0
      ;;
    config)
      cfg_provider="$(sed -n 's/^[[:space:]]*gitProvider:[[:space:]]*"\{0,1\}\([a-zA-Z]*\)"\{0,1\}[[:space:]]*$/\1/p' "${PROJECT_ROOT}/configs/goboot.yml" | head -n1 | tr '[:upper:]' '[:lower:]')"
      if [[ "${cfg_provider}" == "github" || "${cfg_provider}" == "gitlab" ]]; then
        echo "${cfg_provider}"
      else
        echo "both"
      fi
      return 0
      ;;
    *)
      echo "Error: invalid --provider mode '${mode}'. Expected github, gitlab, both, or config."
      exit 1
      ;;
  esac
}

run_step() {
  local title="$1"
  local cmd="$2"

  echo ""
  echo "==> ${title}"
  echo "    ${cmd}"
  bash -lc "${cmd}"
}

run_step_argv() {
  local title="$1"
  shift

  echo ""
  echo "==> ${title}"
  printf "    %q" "$@"
  echo ""
  "$@"
}

gitlab_jobs_for_stage() {
  local stage="$1"

  gitlab-ci-local --list | awk -v stage="${stage}" '
    index($0, "  " stage "  ") {
      name = substr($0, 1, index($0, "  " stage "  ") - 1)
      sub(/[[:space:]]+$/, "", name)
      print name
    }
  '
}

run_gitlab_stage_jobs() {
  local scope="$1"
  local stage="$2"
  local excluded_job="${3:-}"
  local ran="false"
  local job
  local status
  local jobs_output
  local -a jobs=()

  jobs_output="$(gitlab_jobs_for_stage "${stage}")"
  mapfile -t jobs <<<"${jobs_output}"

  for job in "${jobs[@]}"; do
    if [[ -n "${excluded_job}" && "${job}" == "${excluded_job}" ]]; then
      continue
    fi

    cleanup_runtime_dirs
    ran="true"
    set +e
    run_step_argv "${scope}: gitlab-ci-local ${job} (docker)" gitlab-ci-local --network host "${job}"
    status="$?"
    set -e
    if [[ "${status}" -ne 0 ]]; then
      run_step_argv "${scope}: gitlab-ci-local ${job} (shell fallback)" \
        gitlab-ci-local --network host --force-shell-executor --concurrency 1 "${job}"
    fi
  done

  if [[ "${ran}" != "true" ]]; then
    if [[ -n "${excluded_job}" ]] && grep -qx "${excluded_job}" <<<"${jobs_output}"; then
      echo "${scope}: no non-${excluded_job} ${stage} jobs discovered; ${excluded_job} will be validated separately."
      return 0
    fi

    echo "Error: ${scope}: no ${stage} jobs discovered by gitlab-ci-local --list."
    exit 1
  fi
}

generate_github_canary_project() {
  local cfg_dir
  local goboot_cfg

  GITHUB_CANARY_TEMP_DIR="$(mktemp -d)"
  cfg_dir="${GITHUB_CANARY_TEMP_DIR}/configs"
  mkdir -p "${cfg_dir}"
  goboot_cfg="${cfg_dir}/goboot.yml"

  cat >"${goboot_cfg}" <<YAML
targetPath: "${GITHUB_CANARY_TEMP_DIR}/outputs"
projectName: "IntroGitHubCanary"
repoUrl: "https://github.com/projects"
gitProvider: "github"
services:
  - id: "base_project"
    confPath: "${PROJECT_ROOT}/configs/base_project.yml"
    enabled: true
  - id: "base_lint"
    confPath: "${PROJECT_ROOT}/configs/base_lint.yml"
    enabled: true
  - id: "base_test"
    confPath: "${PROJECT_ROOT}/configs/base_test.yml"
    enabled: true
  - id: "base_logger"
    confPath: "${PROJECT_ROOT}/configs/base_logger.yml"
    enabled: true
  - id: "base_docker"
    confPath: "${PROJECT_ROOT}/configs/base_docker.yml"
    enabled: true
  - id: "base_supplychain"
    confPath: "${PROJECT_ROOT}/configs/base_supplychain.yml"
    enabled: true
  - id: "base_local"
    confPath: "${PROJECT_ROOT}/configs/base_local.yml"
    enabled: true
  - id: "base_ci"
    confPath: "${PROJECT_ROOT}/configs/base_ci.yml"
    enabled: true
YAML

  run_step "Generate IntroGitHubCanary for generated GitHub CI canary" "go run cmd/goboot/main.go --config '${goboot_cfg}'"
  GITHUB_CANARY_OUTPUT_DIR="${GITHUB_CANARY_TEMP_DIR}/outputs/IntroGitHubCanary"

  if [[ ! -d "${GITHUB_CANARY_OUTPUT_DIR}/.github/workflows" ]]; then
    echo "Error: expected generated GitHub workflows at ${GITHUB_CANARY_OUTPUT_DIR}."
    exit 1
  fi
}

require_cmd() {
  local cmd="$1"
  if ! command -v "${cmd}" >/dev/null 2>&1; then
    echo "Error: ${cmd} is required but not installed."
    exit 1
  fi
}

require_cmd go
require_cmd make
require_cmd git

mkdir -p "${ARTIFACT_DIR}"
mkdir -p "${TOOLCACHE_DIR}"

cd "${PROJECT_ROOT}"
SELECTED_PROVIDER="$(resolve_provider_mode)"
echo "Provider mode: ${SELECTED_PROVIDER}"

if [[ "${SKIP_GENERATE}" != "true" ]]; then
  if [[ -d "${OUTPUT_DIR}" ]]; then
    run_step "Reset existing outputs/IntroProject" "rm -rf '${OUTPUT_DIR}'"
  fi
  run_step "Generate IntroProject from default config" "go run cmd/goboot/main.go"
fi

if [[ ! -d "${OUTPUT_DIR}" ]]; then
  echo "Error: expected generated output at ${OUTPUT_DIR}."
  exit 1
fi

if [[ "${SELECTED_PROVIDER}" =~ ^(github|both)$ && ! -d "${OUTPUT_DIR}/.github/workflows" ]]; then
  generate_github_canary_project
fi

cleanup_runtime_dirs() {
  # gitlab-ci-local writes transient build state under .gitlab-ci-local.
  # It must not leak into lint scope.
  rm -rf .gitlab-ci-local
  rm -f .gitlab-ci-local-canary.*.yml
}

create_gitlab_lint_canary_file() {
  local file

  mkdir -p .gitlab-ci-local
  file="${PWD}/.gitlab-ci-local/canary-lint.yml"
  TEMP_FILES+=("${file}")

  cat >"${file}" <<'YAML'
stages:
  - lint

include:
  - local: .gitlab/ci/versions.yml

variables:
  DOCKER_HOST: unix:///var/run/docker.sock
  GOLANGCI_LINT_IMAGE: golangci/golangci-lint:v2.12.2
  YAMLLINT_IMAGE: pipelinecomponents/yamllint:0.35.9
  MARKDOWNLINT_IMAGE: ghcr.io/igorshubovych/markdownlint-cli:v0.48.0
  CHECKMAKE_IMAGE: cytopia/checkmake:latest-0.5
  SHELLCHECK_IMAGE: koalaman/shellcheck:v0.11.0
  SHFMT_IMAGE: mvdan/shfmt:v3.13.1
  EDITORCONFIG_CHECKER_IMAGE: mstruebing/editorconfig-checker:v3.6.1

.lint-local:
  stage: lint
  script:
    - ROOT_DIR="${CI_PROJECT_DIR:-$PWD}"
    - DOCKER_RUN_CMD="docker run --rm --network host -v ${ROOT_DIR}:/workdir -w /workdir"
    - SH_FILES="$(find . -type f -name '*.sh' -not -path './templates/*')"
    - $DOCKER_RUN_CMD $GOLANGCI_LINT_IMAGE golangci-lint run ./...
    - $DOCKER_RUN_CMD $YAMLLINT_IMAGE yamllint .
    - $DOCKER_RUN_CMD $MARKDOWNLINT_IMAGE "**/*.md"
    - $DOCKER_RUN_CMD $CHECKMAKE_IMAGE Makefile
    - $DOCKER_RUN_CMD $SHELLCHECK_IMAGE -x $SH_FILES
    - $DOCKER_RUN_CMD $SHFMT_IMAGE -d -i 2 -ci $SH_FILES
    - $DOCKER_RUN_CMD --entrypoint ec $EDITORCONFIG_CHECKER_IMAGE -exclude "(\.git|\.idea|\.gitlab-ci-local|\.gitlab-ci-local-canary.*\.yml|bin)"

lint-branch:
  extends: .lint-local

lint-auto:
  extends: .lint-local

lint-tag:
  extends: .lint-local
YAML

  echo ".gitlab-ci-local/canary-lint.yml"
}

prepull_act_images() {
  local scope="$1"
  local workflow_file=".github/workflows/lint.yml"
  local -a images
  local image

  require_cmd docker

  # Default act runner image used by this repository.
  images=("ghcr.io/catthehacker/ubuntu:act-24.04")

  if [[ -f "${workflow_file}" ]]; then
    while IFS= read -r image; do
      if [[ -n "${image}" ]]; then
        images+=("${image}")
      fi
    done < <(grep -Eo '[a-z0-9]+([._-][a-z0-9]+)*/[a-z0-9./_-]+(@sha256:[a-f0-9]{64}|:[A-Za-z0-9._-]+)' "${workflow_file}" | sort -u || true)
  fi

  local -A seen=()
  local -a unique_images=()
  for image in "${images[@]}"; do
    if [[ -z "${seen[${image}]:-}" ]]; then
      unique_images+=("${image}")
      seen["${image}"]=1
    fi
  done

  for image in "${unique_images[@]}"; do
    run_step "${scope}: pre-pull image ${image}" "docker pull '${image}'"
  done
}

prepull_gitlab_images() {
  local scope="$1"
  local -a files
  local -a images
  local file
  local image
  local var_image
  local version_images_raw

  require_cmd docker

  files=(".gitlab/ci/lint.yml" ".gitlab/ci/build.yml" ".gitlab/ci/test.yml" ".gitlab/ci/security.yml" ".gitlab/ci/versions.yml")

  for file in "${files[@]}"; do
    if [[ -f "${file}" ]]; then
      while IFS= read -r image; do
        if [[ -n "${image}" ]]; then
          images+=("${image}")
        fi
      done < <(grep -Eo '[a-z0-9]+([._-][a-z0-9]+)*/[a-z0-9./_-]+(@sha256:[a-f0-9]{64}|:[A-Za-z0-9._-]+)' "${file}" | sort -u || true)
    fi
  done

  # Include plain-image variables from versions.yml (for example DOCKER_IMAGE: "docker:24.0.5-dind"),
  # but avoid non-image values like DOCKER_HOST=tcp://docker:2375.
  if [[ -f ".gitlab/ci/versions.yml" ]]; then
    version_images_raw="$(sed -n 's/^[[:space:]]*[A-Z0-9_]*_IMAGE:[[:space:]]*"\{0,1\}\([^"[:space:]]\+\)"\{0,1\}[[:space:]]*$/\1/p' ".gitlab/ci/versions.yml")"
    while IFS= read -r var_image; do
      if [[ -n "${var_image}" && "${var_image}" != *"://"* ]]; then
        images+=("${var_image}")
      fi
    done <<<"${version_images_raw}"
  fi

  local -A seen=()
  local -a unique_images=()
  for image in "${images[@]}"; do
    if [[ -z "${seen[${image}]:-}" ]]; then
      unique_images+=("${image}")
      seen["${image}"]=1
    fi
  done

  for image in "${unique_images[@]}"; do
    run_step "${scope}: pre-pull image ${image}" "docker pull '${image}'"
  done
}

run_ci_suite() {
  local scope="$1"
  local has_github_ci="false"
  local has_gitlab_ci="false"
  local temp_git_repo="false"
  local uid
  local gid
  local docker_sock_gid
  local user_opts
  local lint_opts
  local act_opts
  local gitlab_build_jobs

  if [[ -d ".github/workflows" ]]; then
    has_github_ci="true"
  fi
  if [[ -f ".gitlab-ci.yml" ]]; then
    has_gitlab_ci="true"
  fi
  if [[ "${SELECTED_PROVIDER}" == "github" ]]; then
    has_gitlab_ci="false"
  elif [[ "${SELECTED_PROVIDER}" == "gitlab" ]]; then
    has_github_ci="false"
  fi

  if [[ "${has_github_ci}" != "true" && "${has_gitlab_ci}" != "true" ]]; then
    echo "Warning: ${scope}: no GitHub or GitLab CI config found; skipping CI simulation."
    return 0
  fi

  uid="$(id -u)"
  gid="$(id -g)"
  user_opts="--user ${uid}:${gid}"

  if [[ "${has_github_ci}" == "true" ]]; then
    require_cmd act
    require_cmd stat
    if [[ ! -S /var/run/docker.sock ]]; then
      echo "Error: ${scope}: /var/run/docker.sock not found (required for act lint group-add)."
      exit 1
    fi
    docker_sock_gid="$(stat -c '%g' /var/run/docker.sock)"
    lint_opts="--user ${uid}:${gid} --group-add ${docker_sock_gid}"
    act_opts="--pull=${ACT_PULL} --action-offline-mode --use-gitignore=false --artifact-server-path ${ARTIFACT_DIR} --env RUNNER_TOOL_CACHE=${TOOLCACHE_DIR} --env AGENT_TOOLSDIRECTORY=${TOOLCACHE_DIR}"

    if [[ "${PREPULL_IMAGES}" == "true" ]]; then
      prepull_act_images "${scope}"
    fi

    cleanup_runtime_dirs
    run_step "${scope}: act build" "act -j build --container-options \"${user_opts}\" ${act_opts}"
    cleanup_runtime_dirs
    run_step "${scope}: act test" "act -j test --container-options \"${user_opts}\" ${act_opts}"
    cleanup_runtime_dirs
    run_step "${scope}: act lint" "act -j lint --container-options \"${lint_opts}\" ${act_opts}"
    if [[ -f ".github/workflows/container.yml" ]]; then
      cleanup_runtime_dirs
      run_step "${scope}: act container" "act -j container --container-options \"${lint_opts}\" ${act_opts}"
    fi
  else
    echo "Skipping ${scope}: act checks (.github/workflows not found)."
  fi

  if [[ "${has_gitlab_ci}" == "true" ]]; then
    require_cmd gitlab-ci-local
    if [[ "${PREPULL_IMAGES}" == "true" ]]; then
      prepull_gitlab_images "${scope}"
    fi
    if [[ ! -d ".git" ]]; then
      # gitlab-ci-local resolves local includes more reliably in a git-indexed workspace.
      run_step "${scope}: prepare temporary git workspace for gitlab-ci-local" "git init -q && git add -A -f && git -c user.name='goboot-ci-local' -c user.email='goboot-ci-local@local' commit -qm 'temp: ci-local workspace'"
      temp_git_repo="true"
    fi
    run_gitlab_stage_jobs "${scope}" "build" "container"
    gitlab_build_jobs="$(gitlab_jobs_for_stage "build")"
    if [[ -f ".gitlab/ci/container.yml" ]]; then
      if ! grep -qx "container" <<<"${gitlab_build_jobs}"; then
        echo "Error: ${scope}: .gitlab/ci/container.yml exists, but gitlab-ci-local --list did not discover build job 'container'."
        exit 1
      fi

      cleanup_runtime_dirs
      run_step_argv "${scope}: gitlab-ci-local container" gitlab-ci-local --privileged container
    fi
    cleanup_runtime_dirs
    run_step "${scope}: gitlab-ci-local test" "gitlab-ci-local ${GITLAB_CI_LOCAL_SHELL_OPTS} --stage test"
    cleanup_runtime_dirs
    lint_canary_file="$(create_gitlab_lint_canary_file)"
    run_step "${scope}: gitlab-ci-local lint (local shell Docker)" "gitlab-ci-local --file '${lint_canary_file}' --force-shell-executor --concurrency 1 --stage lint"
    cleanup_runtime_dirs
    if [[ "${temp_git_repo}" == "true" ]]; then
      run_step "${scope}: cleanup temporary git workspace" "rm -rf .git"
    fi
  else
    echo "Skipping ${scope}: gitlab-ci-local checks (.gitlab-ci.yml not found)."
  fi
}

run_ci_suite "Root"

cd "${OUTPUT_DIR}"
run_ci_suite "IntroProject"

if [[ -n "${GITHUB_CANARY_OUTPUT_DIR}" ]]; then
  cd "${GITHUB_CANARY_OUTPUT_DIR}"
  run_ci_suite "IntroGitHubCanary"
fi

echo ""
if [[ -n "${GITHUB_CANARY_OUTPUT_DIR}" ]]; then
  echo "CI canary verification completed successfully (root + IntroProject + IntroGitHubCanary)."
else
  echo "CI canary verification completed successfully (root + IntroProject)."
fi
