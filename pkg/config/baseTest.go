package config

import (
	"fmt"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

// BaseTestConfig defines the metadata used by goboot to generate testing setup for a project.
// It injects values into templates (e.g., .golangci.yml) and governs how project-specific testing is rendered.
type BaseTestConfig struct {
	// SourcePath is the path to the template source directory (e.g., "./templates/test_base").
	SourcePath string `yaml:"sourcePath"`

	// UseStyle is the style to be used for testing.
	UseStyle string `yaml:"useStyle"`

	// TestCMD is the command to run tests.
	TestCMD string `yaml:"testCmd"`

	// ProjectName is the short identifier for the project (e.g., "goboot").
	// Used in headings, comments, and other rendered metadata.
	ProjectName string `yaml:"-"`

	// RepoImportPath is the full Go module import path (e.g., "github.com/org/project").
	// Used in test file imports.
	RepoImportPath string `yaml:"-"`

	// CapsProjectName is the uppercase variant of ProjectName (e.g., "GOBOOT").
	CapsProjectName string `yaml:"-"`

	// LowerProjectName is the lowercase variant (e.g., "goboot").
	LowerProjectName string `yaml:"-"`
}

// newBaseTestConfig returns a newly initialized BaseTestConfig with the project name.
func newBaseTestConfig(projectName string) *BaseTestConfig {
	return &BaseTestConfig{
		ProjectName: projectName,
	}
}

// ID returns a stable identifier for this config.
func (bt *BaseTestConfig) ID() string {
	return goboottypes.ServiceNameBaseTest
}

// ReadConfig loads the base test configuration from the provided YAML file path.
//
// It overwrites the current config values with the file contents.
func (bt *BaseTestConfig) ReadConfig(confPath string, repoURL string) error {
	bt.RepoImportPath = strings.TrimPrefix(repoURL, "https://")
	bt.RepoImportPath = strings.TrimPrefix(bt.RepoImportPath, "http://")

	return readYMLConfig(confPath, bt)
}

// Validate verifies the BaseTestConfig for use in scaffolding.
//
// It returns an error if required values are missing/invalid, or calls fillNeededInfos.
func (bt *BaseTestConfig) Validate() error {
	var missing []string

	if strings.TrimSpace(bt.SourcePath) == "" {
		missing = append(missing, "sourcePath")
	}

	if strings.TrimSpace(bt.ProjectName) == "" {
		missing = append(missing, "projectName")
	}

	if strings.TrimSpace(bt.RepoImportPath) == "" {
		missing = append(missing, "repoImportPath")
	}

	if strings.TrimSpace(bt.UseStyle) == "" {
		missing = append(missing, "useStyle")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config fields: %s", strings.Join(missing, ", "))
	}

	bt.fillNeededInfos()

	return bt.validateValues()
}

// fillNeededInfos fills any derived fields in the config.
func (bt *BaseTestConfig) fillNeededInfos() {
	bt.CapsProjectName = strings.ToUpper(bt.ProjectName)
	bt.LowerProjectName = strings.ToLower(bt.ProjectName)

	if strings.TrimSpace(bt.TestCMD) == "" {
		bt.TestCMD = goboottypes.DefaultGoTestCMD
	}
}

// validateValues validates the values in the config.
func (bt *BaseTestConfig) validateValues() error {
	if strings.TrimSpace(bt.UseStyle) != goboottypes.TestStyleGinkgo &&
		strings.TrimSpace(bt.UseStyle) != goboottypes.TestStyleGo {
		return fmt.Errorf("useStyle must be '%s' or '%s'", goboottypes.TestStyleGinkgo, goboottypes.TestStyleGo)
	}

	return nil
}
