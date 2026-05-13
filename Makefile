#  -----------------------------------------------------------------------------
#  Makefile — Developer Targets for `goboot`
#  -----------------------------------------------------------------------------
#
#  This Makefile defines reproducible and documented developer commands for working with the `goboot` project.
#
#  Targets are intended to be simple and transparent.
#
#  Usage:
#    make [target]
#
#  -----------------------------------------------------------------------------

# Project metadata (used in echo and version injection)
PROJECT := goboot
VERSION := $(shell cat .version)

# Include centralized versions
include versions.env

# Lint tooling (containerized)
DOCKER_LINT_CMD := docker run --rm --network host -v "$(PWD)":/workdir -w /workdir

FIND_OWNED_FILES := find . \
	-path './.git' -prune -o \
	-path './.idea' -prune -o \
	-path './.gitlab-ci-local' -prune -o \
	-path './bin' -prune -o \
	-type f

SHELL_FILES := $(shell $(FIND_OWNED_FILES) -name '*.sh' -print)

GOLANGCI_LINT := $(DOCKER_LINT_CMD) golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run ./...
MD_LINT := $(DOCKER_LINT_CMD) ghcr.io/igorshubovych/markdownlint-cli:$(MARKDOWNLINT_VERSION) markdownlint $(shell $(FIND_OWNED_FILES) -name '*.md' -print)
YAML_LINT := $(DOCKER_LINT_CMD) pipelinecomponents/yamllint:$(YAMLLINT_VERSION) yamllint .
CHECKMAKE_LINT := $(DOCKER_LINT_CMD) cytopia/checkmake:$(CHECKMAKE_VERSION) Makefile
SHELLCHECK_LINT := $(DOCKER_LINT_CMD) koalaman/shellcheck:$(SHELLCHECK_VERSION) -x $(SHELL_FILES)
SHFMT_LINT := $(DOCKER_LINT_CMD) mvdan/shfmt:$(SHFMT_VERSION) -d -i 2 -ci $(SHELL_FILES)
EDITORCONFIG_CHECKER_LINT := $(DOCKER_LINT_CMD) --entrypoint ec mstruebing/editorconfig-checker:$(EDITORCONFIG_CHECKER_VERSION) -exclude '(\.git|\.idea|\.gitlab-ci-local|bin)'
COVER_FILE := coverage.txt
TEST_PKGS := $$(go list ./... | grep -v '/test/noauto' | grep -v '/templates')

# .PHONY declares non-file targets to always run when invoked
.PHONY: all build clean test lint release version help lint_go lint_yaml lint_checkmake lint_md lint_sh fmtcheck_sh lint_editorconfig verify_intro verify_ci_canary

#  ----------------------------------------
#  Default target (runs when `make` is called with no args)
#  ----------------------------------------
all: build lint test verify_intro verify_ci_canary

#  ----------------------------------------
#  Build the project
#  ----------------------------------------
build:
	@echo "Building $(PROJECT)..."
	go build -ldflags="-X main.version=$(VERSION)" -o bin/goboot ./cmd/goboot

#  ----------------------------------------
#  Clean build/test artifacts
#  ----------------------------------------
clean:
	@echo "Cleaning artifacts..."

#  ----------------------------------------
#  Run project tests
#  ----------------------------------------
test:
	@echo "Running tests..."
	set -e; \
	go test -race -timeout=9m -coverprofile="$(COVER_FILE)" $(TEST_PKGS); \
	test -f "$(COVER_FILE)" && go tool cover -func="$(COVER_FILE)"; \
	rm -f "$(COVER_FILE)"

#  ----------------------------------------
#  Run linters
#  ----------------------------------------
lint: lint_go lint_yaml lint_checkmake lint_md lint_sh fmtcheck_sh lint_editorconfig

lint_go:
	@echo "golangci-lint..."
	$(GOLANGCI_LINT)

lint_yaml:
	@echo "yamllint..."
	$(YAML_LINT)

lint_checkmake:
	@echo "checkmake..."
	$(CHECKMAKE_LINT)

lint_md:
	@echo "markdownlint..."
	$(MD_LINT)

lint_sh:
	@echo "ShellCheck..."
	$(SHELLCHECK_LINT)

# Formatting check only
fmtcheck_sh:
	@echo "shfmt (check only)..."
	$(SHFMT_LINT)

lint_editorconfig:
	@echo "editorconfig-checker..."
	$(EDITORCONFIG_CHECKER_LINT)

#  ----------------------------------------
#  Full local verification flow
#  ----------------------------------------
verify_intro:
	@echo "Running full Intro project matrix verification flow..."
	./scripts/verify_introproject.sh

verify_ci_canary:
	@echo "Running canary CI verification flow..."
	./scripts/verify_ci_canary.sh

#  ----------------------------------------
#  Release the project
#  ----------------------------------------
release:
	@echo "Running release (version: $(VERSION))..."

#  ----------------------------------------
#  Show the current version
#  ----------------------------------------
version:
	@echo "$(PROJECT) version: $(VERSION)"

#  ----------------------------------------
#  Print usage help
#  Keep in sync with available targets
#  ----------------------------------------
help: help_core help_project help_check help_lint

help_core:
	@echo "Usage:"
	@echo "  make                    Default target (run build, lint, test, verify_intro, verify_ci_canary)"
	@echo "  make help               Show this help message"
	@echo "  make clean              Remove build/test artifacts"

help_project:
	@echo "  make version            Show current project version"
	@echo "  make build              Build the project"
	@echo "  make release            Package the project using GoReleaser"

help_check:
	@echo "  make test               Run project tests"
	@echo "  make lint               Run static code analysis"

help_lint:
	@echo "  make lint_go            Run golangci-lint"
	@echo "  make lint_yaml          Run yamllint"
	@echo "  make lint_checkmake     Run checkmake"
	@echo "  make lint_md            Run markdownlint"
	@echo "  make lint_sh            Run ShellCheck"
	@echo "  make fmtcheck_sh        Run shfmt (check only)"
	@echo "  make lint_editorconfig  Run editorconfig-checker"
	@echo "  make verify_intro       Regenerate Intro project matrix and run full root+output checks"
	@echo "  make verify_ci_canary   Run local CI simulation (act + gitlab-ci-local) for root and IntroProject"
