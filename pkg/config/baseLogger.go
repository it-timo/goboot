package config

import (
	"fmt"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

// BaseLoggerConfig contains logger settings for generated project code.
type BaseLoggerConfig struct {
	// SourcePath is kept for backward-compatible config reads.
	// Logger rendering is owned by file-owning services, not by base_logger.
	SourcePath string `yaml:"sourcePath"`

	// LoggerType selects the logger implementation variant.
	LoggerType string `yaml:"loggerType"`

	// ProjectName is the project identifier.
	ProjectName string `yaml:"-"`

	// RepoPath is the repository URL without scheme.
	RepoPath string `yaml:"-"`

	// CapsProjectName is the uppercase variant of ProjectName.
	CapsProjectName string `yaml:"-"`

	// LowerProjectName is the lowercase variant of ProjectName.
	LowerProjectName string `yaml:"-"`
}

// newBaseLoggerConfig creates a BaseLoggerConfig with the given project name.
func newBaseLoggerConfig(projectName string) *BaseLoggerConfig {
	return &BaseLoggerConfig{
		ProjectName: projectName,
	}
}

// ID returns a stable identifier for this config.
func (bl *BaseLoggerConfig) ID() string {
	return goboottypes.ServiceNameBaseLogger
}

// ReadConfig loads base_logger config from confPath.
func (bl *BaseLoggerConfig) ReadConfig(confPath string, repoURL string, _ string) error {
	bl.RepoPath = strings.TrimPrefix(repoURL, "https://")
	bl.RepoPath = strings.TrimPrefix(bl.RepoPath, "http://")

	return readYMLConfig(confPath, bl)
}

// Validate checks required fields and derived values.
func (bl *BaseLoggerConfig) Validate() error {
	var missing []string

	if strings.TrimSpace(bl.LoggerType) == "" {
		missing = append(missing, "loggerType")
	}

	if strings.TrimSpace(bl.ProjectName) == "" {
		missing = append(missing, "projectName")
	}

	if strings.TrimSpace(bl.RepoPath) == "" {
		missing = append(missing, "repoPath")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config fields: %s", strings.Join(missing, ", "))
	}

	err := validateProjectName(bl.ProjectName)
	if err != nil {
		return err
	}

	bl.fillNeededInfos()

	return bl.validateValues()
}

// LoggerSettings returns validated logger settings for file-owning services.
func (bl *BaseLoggerConfig) LoggerSettings() goboottypes.LoggerSettings {
	return goboottypes.LoggerSettings{
		Enabled: true,
		Type:    bl.LoggerType,
	}
}

func (bl *BaseLoggerConfig) fillNeededInfos() {
	bl.CapsProjectName = strings.ToUpper(bl.ProjectName)
	bl.LowerProjectName = strings.ToLower(bl.ProjectName)
	bl.LoggerType = strings.ToLower(strings.TrimSpace(bl.LoggerType))
}

func (bl *BaseLoggerConfig) validateValues() error {
	switch bl.LoggerType {
	case goboottypes.LoggerTypeZerolog, goboottypes.LoggerTypeSlog:
		return nil
	default:
		return fmt.Errorf("loggerType must be '%s' or '%s'", goboottypes.LoggerTypeZerolog, goboottypes.LoggerTypeSlog)
	}
}
