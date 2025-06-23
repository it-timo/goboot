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

# .PHONY declares non-file targets to always run when invoked
.PHONY: all build clean test lint release version help

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

#  ----------------------------------------
#  Run linters
#  ----------------------------------------
lint:
	@echo "Linting..."
	golangci-lint run ./...
	yamllint .
	checkmake Makefile
	docker run --rm -v "$(PWD)":/workdir -w /workdir ghcr.io/igorshubovych/markdownlint-cli:v0.45.0 markdownlint $(shell find . -name '*.md')

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
help: help_core help_project help_check

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
