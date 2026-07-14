package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

var (
	supplyChainVersionPattern = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)
	licenseNamePattern        = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9.+-]*$`)
)

// BaseSupplyChainConfig contains inputs for supply-chain security automation.
type BaseSupplyChainConfig struct {
	// SourcePath points to supply-chain documentation templates.
	SourcePath string `yaml:"sourcePath"`
	// GovulncheckVersion pins the Go vulnerability scanner.
	GovulncheckVersion string `yaml:"govulncheckVersion"`
	// GoLicensesVersion pins the Go license scanner.
	GoLicensesVersion string `yaml:"goLicensesVersion"`
	// AllowedLicenses lists dependency license identifiers accepted by policy.
	AllowedLicenses []string `yaml:"allowedLicenses"`
	// ProjectName is the generated project identifier.
	ProjectName string `yaml:"-"`
	// GitProvider selects provider-specific CI automation.
	GitProvider string `yaml:"-"`
}

func newBaseSupplyChainConfig(projectName string) *BaseSupplyChainConfig {
	return &BaseSupplyChainConfig{ProjectName: projectName}
}

// ID returns the stable service identifier.
func (bs *BaseSupplyChainConfig) ID() string {
	return goboottypes.ServiceNameBaseSupplyChain
}

// ReadConfig loads base_supplychain config and provider metadata.
func (bs *BaseSupplyChainConfig) ReadConfig(confPath string, _ string, gitProvider string) error {
	bs.GitProvider = strings.ToLower(strings.TrimSpace(gitProvider))

	return readYMLConfig(confPath, bs)
}

// Validate checks security policy inputs and applies deterministic defaults.
func (bs *BaseSupplyChainConfig) Validate() error {
	if strings.TrimSpace(bs.SourcePath) == "" || strings.TrimSpace(bs.ProjectName) == "" || bs.GitProvider == "" {
		return errors.New("missing required config fields: sourcePath, projectName, or gitProvider")
	}

	if err := validateProjectName(bs.ProjectName); err != nil {
		return err
	}

	if bs.GitProvider != goboottypes.GitProviderGitHub && bs.GitProvider != goboottypes.GitProviderGitLab {
		return fmt.Errorf("invalid config: gitProvider %q is not supported", bs.GitProvider)
	}

	bs.fillDefaults()

	if !supplyChainVersionPattern.MatchString(bs.GovulncheckVersion) {
		return fmt.Errorf("invalid config: govulncheckVersion %q must be a pinned semantic version", bs.GovulncheckVersion)
	}

	if !supplyChainVersionPattern.MatchString(bs.GoLicensesVersion) {
		return fmt.Errorf("invalid config: goLicensesVersion %q must be a pinned semantic version", bs.GoLicensesVersion)
	}

	return bs.validateAllowedLicenses()
}

func (bs *BaseSupplyChainConfig) fillDefaults() {
	bs.GovulncheckVersion = strings.TrimSpace(bs.GovulncheckVersion)
	if bs.GovulncheckVersion == "" {
		bs.GovulncheckVersion = "v1.1.4"
	}

	bs.GoLicensesVersion = strings.TrimSpace(bs.GoLicensesVersion)
	if bs.GoLicensesVersion == "" {
		bs.GoLicensesVersion = "v2.0.1"
	}

	if len(bs.AllowedLicenses) == 0 {
		bs.AllowedLicenses = []string{"Apache-2.0", "BSD-2-Clause", "BSD-3-Clause", "ISC", "MIT"}
	}
}

func (bs *BaseSupplyChainConfig) validateAllowedLicenses() error {
	seen := make(map[string]bool, len(bs.AllowedLicenses))

	for index, licenseName := range bs.AllowedLicenses {
		licenseName = strings.TrimSpace(licenseName)
		if !licenseNamePattern.MatchString(licenseName) {
			return fmt.Errorf("invalid config: license %q is not a supported identifier", licenseName)
		}

		if seen[licenseName] {
			return fmt.Errorf("invalid config: duplicate license %q", licenseName)
		}

		seen[licenseName] = true
		bs.AllowedLicenses[index] = licenseName
	}

	return nil
}
