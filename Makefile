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
#  Example:
#    make version  → Show current project version
#    make          → Default target
#
#  -----------------------------------------------------------------------------

# Project metadata (used in echo and version injection)
PROJECT := goboot
VERSION := $(shell cat .version)

# Lint tooling (containerized)
DOCKER_LINT_CMD := docker run --rm -v "$(PWD)":/workdir -w /workdir

GOLANGCI_LINT_IMAGE := golangci/golangci-lint:v2.7.2
GOLANGCI_LINT := $(DOCKER_LINT_CMD) $(GOLANGCI_LINT_IMAGE) golangci-lint run cmd/... pkg/...

MD_LINT_IMAGE := ghcr.io/igorshubovych/markdownlint-cli:v0.47.0
MD_LINT := $(DOCKER_LINT_CMD) $(MD_LINT_IMAGE) markdownlint $(shell find . -name '*.md')

YAMLLINT_IMAGE := pipelinecomponents/yamllint:0.35.9
YAML_LINT := $(DOCKER_LINT_CMD) $(YAMLLINT_IMAGE) yamllint .

CHECKMAKE_LINT_IMAGE := cytopia/checkmake:latest-0.5
CHECKMAKE_LINT := $(DOCKER_LINT_CMD) $(CHECKMAKE_LINT_IMAGE) Makefile

SHELL_FILES := $(shell find . -type f -name '*.sh')

SHELLCHECK_LINT_IMAGE := koalaman/shellcheck:v0.11.0
SHELLCHECK_LINT := $(DOCKER_LINT_CMD) $(SHELLCHECK_LINT_IMAGE) -x $(SHELL_FILES)

SHFMT_LINT_IMAGE := mvdan/shfmt:v3.12.0
SHFMT_LINT := $(DOCKER_LINT_CMD) $(SHFMT_LINT_IMAGE) -d -i 2 -ci $(SHELL_FILES)

# Test coverage
COVER_FILE := coverage.txt
TEST_PKGS := $$(go list ./... | grep -v '/test/noauto' | grep -v '/templates')

# .PHONY declares non-file targets to always run when invoked
.PHONY: all build clean test lint release version help lint_go lint_yaml lint_checkmake lint_md lint_sh fmtcheck_sh

#  ----------------------------------------
#  Default target (runs when `make` is called with no args)
#  ----------------------------------------
all: lint

#  ----------------------------------------
#  Build the project
#  ----------------------------------------
build:
	@echo "Building $(PROJECT)..."

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
	go test -race -timeout=5m -coverprofile="$(COVER_FILE)" $(TEST_PKGS); \
	go tool cover -func="$(COVER_FILE)"; \
	rm -f "$(COVER_FILE)"

#  ----------------------------------------
#  Run linters
#  ----------------------------------------
lint: lint_go lint_yaml lint_checkmake lint_md lint_sh fmtcheck_sh

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
	@echo "  make             Default target (run all)"
	@echo "  make help        Show this help message"
	@echo "  make clean       Remove build/test artifacts"

help_project:
	@echo "  make version     Show current project version"
	@echo "  make build       Build the project"
	@echo "  make release     Package the project using GoReleaser"

help_check:
	@echo "  make test        Run project tests"
	@echo "  make lint        Run static code analysis"

help_lint:
	@echo "  make lint_go         Run golangci-lint"
	@echo "  make lint_yaml       Run yamllint"
	@echo "  make lint_checkmake  Run checkmake"
	@echo "  make lint_md         Run markdownlint"
	@echo "  make lint_sh         Run ShellCheck"
	@echo "  make fmtcheck_sh     Run shfmt (check only)"