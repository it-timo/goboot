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
	-path './outputs' -prune -o \
	-type f

SHELL_FILES := $(shell $(FIND_OWNED_FILES) -name '*.sh' -print0 | xargs -0 printf '%q ')
MD_FILES := $(shell $(FIND_OWNED_FILES) -name '*.md' -print0 | xargs -0 printf '%q ')

GOLANGCI_LINT := $(DOCKER_LINT_CMD) golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run ./...
MD_LINT := $(DOCKER_LINT_CMD) ghcr.io/igorshubovych/markdownlint-cli:$(MARKDOWNLINT_VERSION) markdownlint $(MD_FILES)
YAML_LINT := $(DOCKER_LINT_CMD) pipelinecomponents/yamllint:$(YAMLLINT_VERSION) yamllint .
CHECKMAKE_LINT := $(DOCKER_LINT_CMD) cytopia/checkmake:$(CHECKMAKE_VERSION) Makefile
SHELLCHECK_LINT := $(DOCKER_LINT_CMD) koalaman/shellcheck:$(SHELLCHECK_VERSION) -x $(SHELL_FILES)
SHFMT_LINT := $(DOCKER_LINT_CMD) mvdan/shfmt:$(SHFMT_VERSION) -d -i 2 -ci $(SHELL_FILES)
EDITORCONFIG_CHECKER_LINT := $(DOCKER_LINT_CMD) --entrypoint ec mstruebing/editorconfig-checker:$(EDITORCONFIG_CHECKER_VERSION) -exclude '(\.git|\.idea|\.gitlab-ci-local|bin)'
GORELEASER_CHECK := $(DOCKER_LINT_CMD) goreleaser/goreleaser:$(GORELEASER_VERSION) check
GOBOOT_IMAGE := goboot:$(VERSION)
COVER_FILE := coverage.txt
MIN_COVERAGE := 80
TEST_PKGS := $$(go list ./... | grep -v '/test/noauto' | grep -v '/templates')
BENCH_TIME ?= 1s
PROFILE_DIR ?= profiles

# .PHONY declares non-file targets to always run when invoked
.PHONY: all build docker_build docker_smoke clean test benchmark profile_generation lint release release_check goreleaser_check version help lint_go lint_yaml lint_checkmake lint_md lint_sh fmtcheck_sh lint_editorconfig verify_intro verify_ci_canary verify_dogfood verify_install_upgrade verify_docs

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

docker_build:
	@echo "Building $(PROJECT) container image..."
	docker build --build-arg VERSION="$(VERSION)" -t "$(GOBOOT_IMAGE)" .

docker_smoke: docker_build
	@echo "Running $(PROJECT) container smoke test..."
	docker run --rm --entrypoint sh "$(GOBOOT_IMAGE)" -c 'goboot -h >/tmp/goboot-help 2>&1 && test -s /tmp/goboot-help'

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
	./scripts/test.sh

benchmark:
	@echo "Running performance benchmarks..."
	go test ./... -run '^$$' -bench . -benchmem -benchtime="$(BENCH_TIME)"

profile_generation:
	@echo "Profiling large-project template generation..."
	mkdir -p "$(PROFILE_DIR)"
	go test ./pkg/gobootutils -run '^$$' -bench '^BenchmarkLargeProjectRendering$$' -benchtime="$(BENCH_TIME)" -cpuprofile="$(PROFILE_DIR)/generation-cpu.out" -memprofile="$(PROFILE_DIR)/generation-memory.out"

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
	./scripts/verify_ci_canary.sh --provider=both

verify_dogfood: build
	@echo "Running goboot dogfood generation twice..."
	./scripts/verify_dogfood.sh ./bin/goboot

verify_install_upgrade:
	@echo "Testing baseline-to-candidate binary replacement..."
	./scripts/verify_install_upgrade.sh

verify_docs:
	@echo "Auditing release-candidate documentation..."
	./scripts/verify_docs.sh

#  ----------------------------------------
#  Release the project
#  ----------------------------------------
goreleaser_check:
	@echo "Validating GoReleaser configuration..."
	$(GORELEASER_CHECK)

release_check: all docker_smoke goreleaser_check verify_dogfood verify_install_upgrade verify_docs

release: release_check
	@echo "Release checks passed for $(PROJECT) version: $(VERSION)"

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
	@echo "  make docker_build       Build the goboot container image"
	@echo "  make docker_smoke       Build and smoke-test the goboot container image"
	@echo "  make release            Run the full local release check sequence"
	@echo "  make goreleaser_check   Validate the GoReleaser configuration"
	@echo "  make verify_dogfood     Generate from committed configs twice and compare output"
	@echo "  make verify_install_upgrade Test replacement of a baseline binary"
	@echo "  make verify_docs        Audit release-candidate documentation coverage"

help_check:
	@echo "  make test               Run project tests"
	@echo "  make benchmark          Run performance benchmarks (BENCH_TIME=1s)"
	@echo "  make profile_generation Write CPU and memory profiles to profiles/"
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
	@echo "  make verify_ci_canary   Run local CI simulation for both providers"
