package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/it-timo/goboot/pkg/regeneration"
)

const (
	exitInternal   = 1
	exitUsage      = 2
	exitConfig     = 3
	exitGeneration = 4
	exitConflict   = 5
)

const (
	categoryInternal   = "internal"
	categoryUsage      = "usage"
	categoryConfig     = "config"
	categoryGeneration = "generation"
	categoryConflict   = "conflict"
)

const (
	operationGenerate = "generate"
	operationValidate = "validate"
	operationDryRun   = "dry-run"
	operationVersion  = "version"
)

type outputFormat string

const (
	outputHuman outputFormat = "human"
	outputJSON  outputFormat = "json"
)

type commandError struct {
	code      int
	category  string
	operation string
	plan      *regeneration.Plan
	err       error
}

func newCommandError(
	code int,
	category string,
	operation string,
	err error,
	plan *regeneration.Plan,
) *commandError {
	return &commandError{
		code:      code,
		category:  category,
		operation: operation,
		plan:      plan,
		err:       err,
	}
}

func (commandErr *commandError) Error() string {
	return commandErr.err.Error()
}

func (commandErr *commandError) Unwrap() error {
	return commandErr.err
}

type successResult struct {
	Status    string             `json:"status"`
	Operation string             `json:"operation"`
	Version   string             `json:"version"`
	Config    string             `json:"config,omitempty"`
	Project   string             `json:"project,omitempty"`
	Message   string             `json:"message"`
	Plan      *regeneration.Plan `json:"plan,omitempty"`
}

type errorResult struct {
	Status    string             `json:"status"`
	Operation string             `json:"operation"`
	Version   string             `json:"version"`
	ExitCode  int                `json:"exit_code"`
	Category  string             `json:"category"`
	Message   string             `json:"message"`
	Plan      *regeneration.Plan `json:"plan,omitempty"`
}

func parseOutputFormat(raw string) (outputFormat, error) {
	format := outputFormat(strings.ToLower(strings.TrimSpace(raw)))

	switch format {
	case outputHuman, outputJSON:
		return format, nil
	default:
		return "", fmt.Errorf("unsupported --output value %q (allowed: human, json)", raw)
	}
}

func operationForOptions(opts cliOptions) string {
	if opts.showVersion {
		return operationVersion
	}

	if opts.validateOnly {
		return operationValidate
	}

	if opts.dryRun {
		return operationDryRun
	}

	return operationGenerate
}

func requestedOperation(args []string) string {
	for _, arg := range args {
		switch arg {
		case "--version":
			return operationVersion
		case "--validate":
			return operationValidate
		case "--dry-run":
			return operationDryRun
		}
	}

	return operationGenerate
}

func requestedOutputFormat(args []string) outputFormat {
	for index, arg := range args {
		if strings.EqualFold(arg, "--output=json") {
			return outputJSON
		}

		if arg == "--output" && index+1 < len(args) && strings.EqualFold(args[index+1], string(outputJSON)) {
			return outputJSON
		}
	}

	return outputHuman
}

func commandExitCode(err error) int {
	var typedErr *commandError
	if errors.As(err, &typedErr) {
		return typedErr.code
	}

	return exitInternal
}

func writeSuccess(opts cliOptions, result successResult) {
	result.Status = "success"
	result.Version = version

	format, err := parseOutputFormat(opts.outputFormat)
	if err == nil && format == outputJSON {
		writeJSON(result)

		return
	}

	if result.Plan != nil {
		writeOutputLine(result.Plan.String())
	}

	writeOutputLine(result.Message)
}

func writeCommandError(args []string, err error) {
	result := errorResult{
		Status:    "error",
		Operation: requestedOperation(args),
		Version:   version,
		ExitCode:  commandExitCode(err),
		Category:  categoryInternal,
		Message:   err.Error(),
	}

	var typedErr *commandError
	if errors.As(err, &typedErr) {
		result.Operation = typedErr.operation
		result.Category = typedErr.category
		result.Plan = typedErr.plan
	}

	if requestedOutputFormat(args) == outputJSON {
		writeJSON(result)

		return
	}

	if result.Plan != nil {
		writeOutputLine(result.Plan.String())
	}

	writeOutputLine(result.Message)
}

func writeJSON(value any) {
	content, err := json.Marshal(value)
	if err != nil {
		writeOutputLine(`{"status":"error","exit_code":1,"category":"internal","message":"failed to encode output"}`)

		return
	}

	writeOutputLine(string(content))
}
