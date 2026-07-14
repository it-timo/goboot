package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

var (
	governanceOwnerPattern  = regexp.MustCompile(`^@[A-Za-z0-9][A-Za-z0-9_.-]*(/[A-Za-z0-9][A-Za-z0-9_.-]*)*$`)
	governanceBranchPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]*$`)
)

// BaseGovernanceConfig contains inputs for repository governance files.
type BaseGovernanceConfig struct {
	// SourcePath points to governance templates.
	SourcePath string `yaml:"sourcePath"`
	// Maintainers lists provider handles written to CODEOWNERS.
	Maintainers []string `yaml:"maintainers"`
	// DefaultBranch is the contribution target branch.
	DefaultBranch string `yaml:"defaultBranch"`
	// ProjectName is the generated project identifier.
	ProjectName string `yaml:"-"`
	// ProjectURL is the generated repository URL.
	ProjectURL string `yaml:"-"`
	// GitProvider selects provider-specific contribution templates.
	GitProvider string `yaml:"-"`
	// Profile selects the governance baseline.
	Profile string `yaml:"-"`
}

func newBaseGovernanceConfig(projectName string) *BaseGovernanceConfig {
	return &BaseGovernanceConfig{ProjectName: projectName}
}

// ID returns the stable service identifier.
func (bg *BaseGovernanceConfig) ID() string {
	return goboottypes.ServiceNameBaseGovernance
}

// SetProfile injects the validated template profile.
func (bg *BaseGovernanceConfig) SetProfile(profile string) {
	bg.Profile = profile
}

// ReadConfig loads base_governance config and repository metadata.
func (bg *BaseGovernanceConfig) ReadConfig(confPath string, repoURL string, gitProvider string) error {
	bg.ProjectURL = strings.TrimSuffix(strings.TrimSpace(repoURL), "/")
	bg.GitProvider = strings.ToLower(strings.TrimSpace(gitProvider))

	return readYMLConfig(confPath, bg)
}

// Validate checks governance inputs and applies deterministic defaults.
func (bg *BaseGovernanceConfig) Validate() error {
	err := bg.validateRequiredFields()
	if err != nil {
		return err
	}

	err = validateProjectName(bg.ProjectName)
	if err != nil {
		return err
	}

	err = bg.validateProvider()
	if err != nil {
		return err
	}

	err = bg.validateDefaultBranch()
	if err != nil {
		return err
	}

	return bg.validateMaintainers()
}

func (bg *BaseGovernanceConfig) validateRequiredFields() error {
	if strings.TrimSpace(bg.SourcePath) == "" || strings.TrimSpace(bg.ProjectName) == "" ||
		bg.ProjectURL == "" || bg.GitProvider == "" || bg.Profile == "" {
		return errors.New("missing required config fields: sourcePath, projectName, projectUrl, gitProvider, or profile")
	}

	return nil
}

func (bg *BaseGovernanceConfig) validateProvider() error {
	if bg.GitProvider != goboottypes.GitProviderGitHub && bg.GitProvider != goboottypes.GitProviderGitLab {
		return fmt.Errorf("invalid config: gitProvider %q is not supported", bg.GitProvider)
	}

	return nil
}

func (bg *BaseGovernanceConfig) validateDefaultBranch() error {
	bg.DefaultBranch = strings.TrimSpace(bg.DefaultBranch)
	if bg.DefaultBranch == "" {
		bg.DefaultBranch = "main"
	}

	if !governanceBranchPattern.MatchString(bg.DefaultBranch) || strings.Contains(bg.DefaultBranch, "..") ||
		strings.Contains(bg.DefaultBranch, "//") || strings.HasSuffix(bg.DefaultBranch, "/") {
		return fmt.Errorf("invalid config: defaultBranch %q is not supported", bg.DefaultBranch)
	}

	return nil
}

func (bg *BaseGovernanceConfig) validateMaintainers() error {
	if len(bg.Maintainers) == 0 {
		return errors.New("invalid config: at least one maintainer is required")
	}

	seen := make(map[string]bool, len(bg.Maintainers))

	for index, maintainer := range bg.Maintainers {
		maintainer = strings.TrimSpace(maintainer)
		if !governanceOwnerPattern.MatchString(maintainer) {
			return fmt.Errorf("invalid config: maintainer %q is not a provider handle", maintainer)
		}

		if seen[maintainer] {
			return fmt.Errorf("invalid config: duplicate maintainer %q", maintainer)
		}

		seen[maintainer] = true
		bg.Maintainers[index] = maintainer
	}

	return nil
}
