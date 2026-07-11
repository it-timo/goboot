package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

var releaseNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// BaseReleaseConfig contains inputs for GoReleaser and provider release automation.
type BaseReleaseConfig struct {
	// SourcePath points to release templates.
	SourcePath string   `yaml:"sourcePath"`
	// BinaryName is the archive binary name.
	BinaryName string   `yaml:"binaryName"`
	// MainPackage is the package GoReleaser builds.
	MainPackage string   `yaml:"mainPackage"`
	// Formats lists enabled archive formats.
	Formats     []string `yaml:"formats"`
	// ProjectName is the generated project identifier.
	ProjectName string   `yaml:"-"`
	// GitProvider selects the provider release template.
	GitProvider string   `yaml:"-"`
}

func newBaseReleaseConfig(projectName string) *BaseReleaseConfig {
	return &BaseReleaseConfig{ProjectName: projectName}
}

// ID returns the stable service identifier.
func (br *BaseReleaseConfig) ID() string {
	return goboottypes.ServiceNameBaseRelease
}

// ReadConfig loads base_release config and provider metadata.
func (br *BaseReleaseConfig) ReadConfig(confPath string, _ string, gitProvider string) error {
	br.GitProvider = strings.ToLower(strings.TrimSpace(gitProvider))

	return readYMLConfig(confPath, br)
}

// Validate checks release inputs and applies deterministic defaults.
func (br *BaseReleaseConfig) Validate() error {
	if strings.TrimSpace(br.SourcePath) == "" || strings.TrimSpace(br.ProjectName) == "" || br.GitProvider == "" {
		return errors.New("missing required config fields: sourcePath, projectName, or gitProvider")
	}
	err := validateProjectName(br.ProjectName)
	if err != nil {
		return err
	}

	if br.GitProvider != goboottypes.GitProviderGitHub && br.GitProvider != goboottypes.GitProviderGitLab {
		return fmt.Errorf("invalid config: gitProvider %q is not supported", br.GitProvider)
	}

	br.BinaryName = strings.TrimSpace(br.BinaryName)
	if br.BinaryName == "" {
		br.BinaryName = strings.ToLower(br.ProjectName)
	}

	if !releaseNamePattern.MatchString(br.BinaryName) {
		return fmt.Errorf("invalid config: binaryName %q is not supported", br.BinaryName)
	}

	br.MainPackage = strings.TrimSpace(br.MainPackage)
	if br.MainPackage == "" {
		br.MainPackage = "./cmd/" + strings.ToLower(br.ProjectName)
	}

	if strings.Contains(br.MainPackage, "..") || !dockerMainPackagePattern.MatchString(br.MainPackage) {
		return fmt.Errorf("invalid config: mainPackage %q is not supported", br.MainPackage)
	}

	if len(br.Formats) == 0 {
		br.Formats = []string{"tar.gz", "zip"}
	}

	seen := map[string]bool{}
	for _, format := range br.Formats {
		format = strings.TrimSpace(format)
		if format != "tar.gz" && format != "zip" {
			return fmt.Errorf("invalid config: format %q is not supported", format)
		}

		if seen[format] {
			return fmt.Errorf("invalid config: duplicate format %q", format)
		}

		seen[format] = true
	}

	return nil
}
