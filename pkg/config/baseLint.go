package config

import (
	"fmt"
	"strings"

	"github.com/it-timo/goboot/pkg/types"
)

// BaseLintConfig defines the metadata used by goboot to generate linting setup for a project.
// It injects values into templates (e.g., .golangci.yml) and governs how project-specific linting is rendered.
type BaseLintConfig struct {
	// SourcePath is the path to the template source directory (e.g., "./templates/lint_base").
	SourcePath string `yaml:"sourcePath"`

	// ProjectName is the short identifier for the project (e.g., "goboot").
	// Used in headings, comments, and other rendered metadata.
	ProjectName string `yaml:"projectName"`

	// RepoImportPath is the full Go module import path (e.g., "github.com/org/project").
	// Used in linter config like depguard to enforce proper import usage.
	RepoImportPath string `yaml:"repoImportPath"`

	// Linters is a map of named linter configs to be enabled for this project.
	Linters map[string]*Linter `yaml:"linters"`
}

// Linter defines an individual linter to be included in the generated linting setup.
type Linter struct {
	// Cmd is the shell command to run the linter (e.g., "golangci-lint run").
	Cmd string `yaml:"cmd"`

	// Enabled indicates whether this linter is active for the project.
	Enabled bool `yaml:"enabled"`
}

// lintCmds defines the default commands used if no custom `Cmd` is set in the config.
var lintCmds = map[string]string{
	types.LinterGo:   types.DefaultGoLintCmd,
	types.LinterYAML: types.DefaultYMLLintCmd,
	types.LinterMake: types.DefaultMakeLintCmd,
	types.LinterMD:   types.DefaultMDLintCmd,
}

// newBaseLintConfig returns a newly initialized BaseLintConfig with the project name.
func newBaseLintConfig(projectName string) *BaseLintConfig {
	return &BaseLintConfig{
		ProjectName: projectName,
	}
}

// ID returns a stable identifier for this config.
func (bl *BaseLintConfig) ID() string {
	return types.ServiceNameBaseLint
}

// ReadConfig loads the base lint configuration from the provided YAML file path.
//
// It overwrites the current config values with the file contents.
func (bl *BaseLintConfig) ReadConfig(confPath string) error {
	return readYMLConfig(confPath, bl)
}

// Validate verifies the BaseLintConfig for use in scaffolding.
//
// It returns an error if required values are missing/invalid, or calls fillNeededInfos.
func (bl *BaseLintConfig) Validate() error {
	var missing []string

	if strings.TrimSpace(bl.SourcePath) == "" {
		missing = append(missing, "sourcePath")
	}

	if strings.TrimSpace(bl.ProjectName) == "" {
		missing = append(missing, "projectName")
	}

	if strings.TrimSpace(bl.RepoImportPath) == "" {
		missing = append(missing, "repoImportPath")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config fields: %s", strings.Join(missing, ", "))
	}

	bl.fillNeededInfos()

	return nil
}

// fillNeededInfos fills any derived fields in the config.
func (bl *BaseLintConfig) fillNeededInfos() {
	// Apply default cmd if not set.
	for name, linter := range bl.Linters {
		if strings.TrimSpace(linter.Cmd) != "" {
			continue
		}

		if linter.Enabled {
			cmd, exist := lintCmds[name]
			if exist {
				linter.Cmd = cmd
			}

			fmt.Printf("[WARN] Unknown linter %q; no default command defined", name)
		}
	}
}
