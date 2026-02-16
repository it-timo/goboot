package config

import (
	"fmt"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

// BaseTestConfig contains template inputs for test scaffolding.
type BaseTestConfig struct {
	// SourcePath points to test templates.
	SourcePath string `yaml:"sourcePath"`

	// UseStyle selects the test style.
	UseStyle string `yaml:"useStyle"`

	// TestCMD is the generated test command.
	TestCMD string `yaml:"testCmd"`

	// ProjectName is the project identifier.
	ProjectName string `yaml:"-"`

	// RepoImportPath is the import path used in generated tests.
	RepoImportPath string `yaml:"-"`

	// CapsProjectName is the uppercase variant of ProjectName.
	CapsProjectName string `yaml:"-"`

	// LowerProjectName is the lowercase variant of ProjectName.
	LowerProjectName string `yaml:"-"`
}

// newBaseTestConfig creates a BaseTestConfig with the given project name.
func newBaseTestConfig(projectName string) *BaseTestConfig {
	return &BaseTestConfig{
		ProjectName: projectName,
	}
}

// ID returns a stable identifier for this config.
func (bt *BaseTestConfig) ID() string {
	return goboottypes.ServiceNameBaseTest
}

// ReadConfig loads base_test config from confPath.
func (bt *BaseTestConfig) ReadConfig(confPath string, repoURL string, _ string) error {
	bt.RepoImportPath = strings.TrimPrefix(repoURL, "https://")
	bt.RepoImportPath = strings.TrimPrefix(bt.RepoImportPath, "http://")

	return readYMLConfig(confPath, bt)
}

// Validate checks required fields and derived values.
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

// fillNeededInfos populates derived fields used by templates.
func (bt *BaseTestConfig) fillNeededInfos() {
	bt.CapsProjectName = strings.ToUpper(bt.ProjectName)
	bt.LowerProjectName = strings.ToLower(bt.ProjectName)

	if strings.TrimSpace(bt.TestCMD) == "" {
		bt.TestCMD = goboottypes.DefaultGoTestCMD
	}
}

// validateValues validates enumerated config values.
func (bt *BaseTestConfig) validateValues() error {
	if strings.TrimSpace(bt.UseStyle) != goboottypes.TestStyleGinkgo &&
		strings.TrimSpace(bt.UseStyle) != goboottypes.TestStyleGo {
		return fmt.Errorf("useStyle must be '%s' or '%s'", goboottypes.TestStyleGinkgo, goboottypes.TestStyleGo)
	}

	return nil
}
