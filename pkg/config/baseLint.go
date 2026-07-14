package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/rs/zerolog/log"
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

	// Profile selects the generated Go lint baseline.
	Profile string `yaml:"-"`

	// CyclomaticComplexity is the profile-derived cyclop threshold.
	CyclomaticComplexity int `yaml:"-"`
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

// SetProfile injects the validated template profile.
func (bl *BaseLintConfig) SetProfile(profile string) {
	bl.Profile = profile
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

	err := validateProjectName(bl.ProjectName)
	if err != nil {
		return err
	}

	err = bl.validateLinters()
	if err != nil {
		return err
	}

	bl.fillNeededInfos()
	bl.fillProfileSettings()

	return nil
}

func (bl *BaseLintConfig) fillProfileSettings() {
	switch bl.Profile {
	case goboottypes.ProfileMinimal:
		bl.CyclomaticComplexity = 20
	case goboottypes.ProfileEnterprise:
		bl.CyclomaticComplexity = 8
	case goboottypes.ProfileOSS:
		bl.CyclomaticComplexity = 10
	default:
		bl.CyclomaticComplexity = 10
	}
}

func (bl *BaseLintConfig) validateLinters() error {
	if len(bl.Linters) == 0 {
		return errors.New("invalid config: linters must not be empty")
	}

	for name, linter := range bl.Linters {
		if _, exists := lintCmds[name]; !exists {
			return fmt.Errorf("invalid config: linter %q is not supported", name)
		}

		if linter == nil {
			return fmt.Errorf("invalid config: linter %q is nil", name)
		}
	}

	return nil
}

// fillNeededInfos assigns default lint commands for enabled linters.
func (bl *BaseLintConfig) fillNeededInfos() {
	for name, linter := range bl.Linters {
		if linter == nil {
			continue
		}

		if strings.TrimSpace(linter.Cmd) != "" {
			continue
		}

		if linter.Enabled {
			cmd, exist := lintCmds[name]
			if exist {
				linter.Cmd = cmd

				continue
			}

			log.Warn().Str("linter", name).Msg("unknown linter; no default command defined")
		}
	}
}
