package goboottypes

// Default lint commands (overridable via config).
const (
	// DefaultGoLintCmd is the default command for the "go" linter.
	DefaultGoLintCmd = "{{DOCKER_RUN}} golangci/golangci-lint:v2.7.2 golangci-lint run ./..."
	// DefaultYMLLintCmd is the default command for the "yaml" linter.
	DefaultYMLLintCmd = "{{DOCKER_RUN}} pipelinecomponents/yamllint:0.35.9 yamllint ."
	// DefaultMakeLintCmd is the default command for the "make" linter.
	DefaultMakeLintCmd = "{{DOCKER_RUN}} cytopia/checkmake:latest-0.5 Makefile"
	// DefaultMDLintCmd is the default command for the "md" linter.
	DefaultMDLintCmd = "{{DOCKER_RUN}} ghcr.io/igorshubovych/markdownlint-cli:v0.47.0 markdownlint \"**/*.md\""
	// DefaultShellLintCmd is the default command for the "shell" linter.
	DefaultShellLintCmd = "{{DOCKER_RUN}} cytopia/shellcheck:latest-0.8.0 shellcheck {{SH_FILES}}"
	// DefaultSHFMTCmd is the default command for the "shfmt" linter.
	DefaultSHFMTCmd = "{{DOCKER_RUN}} cytopia/shfmt:latest-1.10 shfmt -d {{SH_FILES}}"
	// DefaultEditorLintCmd is the default command for the "editor" linter.
	DefaultEditorLintCmd = "{{DOCKER_RUN}} --entrypoint ec " +
		"mstruebing/editorconfig-checker:v3.6.0 -exclude '(\\.git|\\.idea|\\.vscode)'"
)

// Default test commands.
const (
	// DefaultGoTestCMD is the default command for running tests.
	DefaultGoTestCMD = "go test -race -timeout=5m -coverprofile=coverage.txt " +
		"&& go tool cover -func=coverage.txt; rm -f coverage.txt"
)

// Default CI commands.
const (
	// DefaultGoBuildCMD is the default command for running builds in CI.
	DefaultGoBuildCMD = "go build ./..."
)

// Linter IDs.
const (
	// LinterGo is the default identifier for the "go" linter.
	LinterGo = "golang"
	// LinterYAML is the default identifier for the "yaml" linter.
	LinterYAML = "yaml"
	// LinterMake is the default identifier for the "make" linter.
	LinterMake = "make"
	// LinterMD is the default identifier for the "md" linter.
	LinterMD = "markdown"
	// LinterShell is the default identifier for the "shell" linter.
	LinterShell = "shellcheck"
	// LinterSHFMT is the default identifier for the "shfmt" linter.
	LinterSHFMT = "shfmt"
	// LinterEditor is the default identifier for the "editor" linter.
	LinterEditor = "editor"
)

// Test style IDs.
const (
	// TestStyleGinkgo is the default identifier for the "ginkgo" test style.
	TestStyleGinkgo = "ginkgo"
	// TestStyleGo is the default identifier for the "go" test style.
	TestStyleGo = "go"
)

// Git provider IDs.
const (
	// GitProviderGitLab is the default identifier for the "gitlab" git provider.
	GitProviderGitLab = "gitlab"
	// GitProviderGitHub is the default identifier for the "github" git provider.
	GitProviderGitHub = "github"
)

// Local output selector names.
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

// Local directory names.
const (
	// ScriptDirNameScript is the default name for the "script" dir.
	ScriptDirNameScript = "scripts"
)

// Local script filenames.
const (
	// ScriptFileLint is the default name for the "lint" script file in the "script" dir.
	ScriptFileLint = "lint.sh"
	// ScriptFileTest is the default name for the "test" script file in the "script" dir.
	ScriptFileTest = "test.sh"
)

// CI job filenames.
const (
	// CIFileSuffix is the default suffix for CI job files.
	CIFileSuffix = ".yml"
	// CIFileLint is the default name for the lint CI job file.
	CIFileLint = "lint.yml"
	// CIFileTest is the default name for the test CI job file.
	CIFileTest = "test.yml"
	// CIFileBuild is the default name for the build CI job file.
	CIFileBuild = "build.yml"
)

// Template source guardrails.
const (
	// MaxTemplateSourceFiles limits non-directory files processed from template sources.
	MaxTemplateSourceFiles = 5000
	// MaxTemplateSourceBytes limits total bytes processed from template sources.
	MaxTemplateSourceBytes = 64 * 1024 * 1024
)
