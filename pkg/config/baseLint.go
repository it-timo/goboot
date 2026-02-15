package config

import (
	"fmt"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

// BaseLintConfig configures lint template rendering and lint command defaults.
type BaseLintConfig struct {
	// SourcePath points to lint template files.
	SourcePath string `yaml:"sourcePath"`

	// ProjectName is injected into rendered output.
	ProjectName string `yaml:"-"`

	// RepoImportPath is used by lint configs that check import paths.
	RepoImportPath string `yaml:"-"`

	// Linters maps linter IDs to config.
	Linters map[string]*Linter `yaml:"linters"`

	// AllowedPackages lists import exceptions for lint policies.
	AllowedPackages []string `yaml:"allowedPackages"`
}

// Linter configures one linter entry.
type Linter struct {
	// Cmd is the linter command.
	Cmd string `yaml:"cmd"`

	// Enabled controls whether the linter is active.
	Enabled bool `yaml:"enabled"`
}

// lintCmds defines defaults used when Cmd is empty.
var lintCmds = map[string]string{
	goboottypes.LinterGo:     goboottypes.DefaultGoLintCmd,
	goboottypes.LinterYAML:   goboottypes.DefaultYMLLintCmd,
	goboottypes.LinterMake:   goboottypes.DefaultMakeLintCmd,
	goboottypes.LinterMD:     goboottypes.DefaultMDLintCmd,
	goboottypes.LinterShell:  goboottypes.DefaultShellLintCmd,
	goboottypes.LinterSHFMT:  goboottypes.DefaultSHFMTCmd,
	goboottypes.LinterEditor: goboottypes.DefaultEditorLintCmd,
}

// newBaseLintConfig constructs BaseLintConfig with derived project name.
func newBaseLintConfig(projectName string) *BaseLintConfig {
	return &BaseLintConfig{
		ProjectName: projectName,
	}
}

// ID returns a stable identifier for this config.
func (bl *BaseLintConfig) ID() string {
	return goboottypes.ServiceNameBaseLint
}

// ReadConfig loads base_lint YAML and derives RepoImportPath from repoURL.
func (bl *BaseLintConfig) ReadConfig(confPath string, repoURL string, _ string) error {
	bl.RepoImportPath = strings.TrimPrefix(repoURL, "https://")
	bl.RepoImportPath = strings.TrimPrefix(bl.RepoImportPath, "http://")

	return readYMLConfig(confPath, bl)
}

// Validate checks required fields and applies derived defaults.
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

// fillNeededInfos assigns default lint commands for enabled linters.
func (bl *BaseLintConfig) fillNeededInfos() {
	for name, linter := range bl.Linters {
		if strings.TrimSpace(linter.Cmd) != "" {
			continue
		}

		if linter.Enabled {
			cmd, exist := lintCmds[name]
			if exist {
				linter.Cmd = cmd

				continue
			}

			fmt.Printf("[WARN] Unknown linter %q; no default command defined", name)
		}
	}
}
