package types

// Default linter commands.
//
// Can be overridden using the "linters.cmd" config field.
const (
	// DefaultGoLintCmd is the default command for the "go" linter.
	DefaultGoLintCmd = "golangci-lint run ./..."
	// DefaultYMLLintCmd is the default command for the "yaml" linter.
	DefaultYMLLintCmd = "yamllint ."
	// DefaultMakeLintCmd is the default command for the "make" linter.
	DefaultMakeLintCmd = "checkmake Makefile"
	// DefaultMDLintCmd is the default command for the "md" linter.
	DefaultMDLintCmd = "docker run -v \"$PWD\":/workdir -w /workdir \\\n" +
		"ghcr.io/igorshubovych/markdownlint-cli:v0.45.0 \"**/*.md\""
)

// Default linter identifiers.
const (
	// LinterGo is the default identifier for the "go" linter.
	LinterGo = "golang"
	// LinterYAML is the default identifier for the "yaml" linter.
	LinterYAML = "yaml"
	// LinterMake is the default identifier for the "make" linter.
	LinterMake = "make"
	// LinterMD is the default identifier for the "md" linter.
	LinterMD = "markdown"
)

// Default local script names.
const (
	// ScriptNameMake is the default name for the "make" script.
	ScriptNameMake = "make"
	// ScriptNameTask is the default name for the "taskfile" script.
	ScriptNameTask = "task"
	// ScriptNameScript is the default name for the "script" dir.
	ScriptNameScript = "script"
	// ScriptNameCommit is the default name for the "pre-commit" script.
	ScriptNameCommit = "commit"
)

// Default local file names.
const (
	// ScriptDirNameScript is the default name for the "script" dir.
	ScriptDirNameScript = "scripts"
)

// Default local script file names.
const (
	// ScriptFileLint is the default name for the "lint" script file in the "script" dir.
	ScriptFileLint = "lint.sh"
)
