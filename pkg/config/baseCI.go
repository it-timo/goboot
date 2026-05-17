package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

var (
	goVersionPattern = regexp.MustCompile(`^[0-9]+(\.[0-9]+){1,2}$`)
	branchPattern    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]*$`)
)

// BaseCIConfig contains inputs for CI template generation.
type BaseCIConfig struct {
	// SourcePath points to provider CI templates.
	SourcePath string `yaml:"sourcePath"`

	// GoVersion lists Go versions used by generated CI jobs.
	GoVersion []string `yaml:"goVersions"`

	// AutoBranches lists branches that trigger CI automatically.
	AutoBranches []string `yaml:"autoBranches"`

	// ImagePolicy controls image pinning behavior (balanced, strict, simple).
	ImagePolicy string `yaml:"imagePolicy"`

	// Jobs defines CI-only jobs owned by base_ci.
	Jobs map[string]*CIJob `yaml:"jobs"`

	// ProjectName is the project identifier.
	ProjectName string `yaml:"-"`

	// GitProvider selects provider-specific CI layout.
	GitProvider string `yaml:"-"`
}

// CIJob defines one generated CI job.
type CIJob struct {
	// File is the rendered job filename.
	File string `yaml:"-"`

	// Commands are shell commands run by the job.
	Commands []string `yaml:"commands"`

	// AllowFailure marks a non-blocking job when supported.
	AllowFailure bool `yaml:"allowFailure"`
}

// newBaseCIConfig creates a BaseCIConfig with the given project name.
func newBaseCIConfig(projectName string) *BaseCIConfig {
	return &BaseCIConfig{
		ProjectName: projectName,
	}
}

// ID returns a stable identifier for this config.
func (bc *BaseCIConfig) ID() string {
	return goboottypes.ServiceNameBaseCI
}

// ReadConfig loads base_ci config from confPath.
func (bc *BaseCIConfig) ReadConfig(confPath string, _ string, gitProvider string) error {
	bc.GitProvider = strings.TrimSpace(gitProvider)

	return readYMLConfig(confPath, bc)
}

// Validate checks required fields and normalizes optional values.
func (bc *BaseCIConfig) Validate() error {
	if len(bc.GoVersion) == 0 {
		return errors.New("missing required config fields: goVersions")
	} else {
		err := bc.validateGoVersion()
		if err != nil {
			return err
		}
	}

	err := bc.validateAutoBranches()
	if err != nil {
		return err
	}

	err = bc.validateImagePolicy()
	if err != nil {
		return err
	}

	err = bc.validateRequiredFields()
	if err != nil {
		return err
	}

	err = validateProjectName(bc.ProjectName)
	if err != nil {
		return err
	}

	err = bc.validateJobs()
	if err != nil {
		return err
	}

	return nil
}

func (bc *BaseCIConfig) validateRequiredFields() error {
	var missing []string

	if strings.TrimSpace(bc.SourcePath) == "" {
		missing = append(missing, "sourcePath")
	}

	if strings.TrimSpace(bc.ProjectName) == "" {
		missing = append(missing, "projectName")
	}

	if strings.TrimSpace(bc.GitProvider) == "" {
		missing = append(missing, "gitProvider")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config fields: %s", strings.Join(missing, ", "))
	}

	return nil
}

// validateGoVersion rejects blank goVersions entries.
func (bc *BaseCIConfig) validateGoVersion() error {
	for _, version := range bc.GoVersion {
		trimmed := strings.TrimSpace(version)
		if trimmed == "" {
			return errors.New("invalid config: goVersions contains empty string")
		}

		if !goVersionPattern.MatchString(trimmed) {
			return fmt.Errorf("invalid config: goVersion %q is not supported", version)
		}
	}

	return nil
}

// validateAutoBranches validates branches or sets defaults.
func (bc *BaseCIConfig) validateAutoBranches() error {
	if len(bc.AutoBranches) == 0 {
		bc.AutoBranches = []string{"main", "master"}

		return nil
	}

	for _, branch := range bc.AutoBranches {
		trimmed := strings.TrimSpace(branch)
		if trimmed == "" {
			return errors.New("invalid config: autoBranches contains empty string")
		}

		if !branchPattern.MatchString(trimmed) {
			return fmt.Errorf("invalid config: autoBranch %q is not supported", branch)
		}
	}

	return nil
}

// validateImagePolicy validates and defaults imagePolicy.
func (bc *BaseCIConfig) validateImagePolicy() error {
	policy := strings.ToLower(strings.TrimSpace(bc.ImagePolicy))
	if policy == "" {
		bc.ImagePolicy = "balanced"

		return nil
	}

	switch policy {
	case "balanced", "strict", "simple":
		bc.ImagePolicy = policy
	default:
		return fmt.Errorf("invalid config: imagePolicy %q is not supported", bc.ImagePolicy)
	}

	return nil
}

// validateJobs validates configured CI job definitions.
func (bc *BaseCIConfig) validateJobs() error {
	for name, job := range bc.Jobs {
		if job == nil {
			return fmt.Errorf("invalid config: job %q is nil", name)
		}

		job.File = strings.ToLower(strings.TrimSpace(name)) + goboottypes.CIFileSuffix

		switch job.File {
		case goboottypes.CIFileLint, goboottypes.CIFileTest, goboottypes.CIFileBuild:
			// Valid predefined job file.
		default:
			return fmt.Errorf("invalid config: job %q has invalid job name", name)
		}

		if len(job.Commands) == 0 {
			return fmt.Errorf("invalid config: job %q missing commands", name)
		}

		for _, cmd := range job.Commands {
			if strings.TrimSpace(cmd) != "" {
				continue
			}

			return fmt.Errorf("invalid config: job %q contains blank commands", name)
		}
	}

	return nil
}
