package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

// BaseLocalConfig contains inputs for local tooling generation.
type BaseLocalConfig struct {
	// SourcePath points to local tooling templates.
	SourcePath string `yaml:"sourcePath"`

	// ProjectName is the project identifier.
	ProjectName string `yaml:"-"`

	// FileList lists template files to generate.
	FileList []string `yaml:"fileList"`
}

// newBaseLocalConfig creates a BaseLocalConfig with the given project name.
func newBaseLocalConfig(projectName string) *BaseLocalConfig {
	return &BaseLocalConfig{
		ProjectName: projectName,
	}
}

// ID returns a stable identifier for this config.
func (bl *BaseLocalConfig) ID() string {
	return goboottypes.ServiceNameBaseLocal
}

// ReadConfig loads base_local config from confPath.
func (bl *BaseLocalConfig) ReadConfig(confPath string, _ string, _ string) error {
	return readYMLConfig(confPath, bl)
}

// Validate checks required fields and list validity.
func (bl *BaseLocalConfig) Validate() error {
	var missing []string

	if strings.TrimSpace(bl.SourcePath) == "" {
		missing = append(missing, "sourcePath")
	}

	if strings.TrimSpace(bl.ProjectName) == "" {
		missing = append(missing, "projectName")
	}

	// Empty file lists are invalid; disable the service instead.
	if len(bl.FileList) == 0 {
		missing = append(missing, "fileList")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config fields: %s", strings.Join(missing, ", "))
	}

	err := validateProjectName(bl.ProjectName)
	if err != nil {
		return err
	}

	return bl.validateFileList()
}

func (bl *BaseLocalConfig) validateFileList() error {
	seen := make(map[string]struct{})
	invalid := []string{}

	for _, file := range bl.FileList {
		trimmed := strings.TrimSpace(file)

		err := validateLocalFileListEntry(trimmed, seen)
		if err != nil {
			invalid = append(invalid, err.Error())

			continue
		}

		seen[trimmed] = struct{}{}
	}

	if len(invalid) > 0 {
		return fmt.Errorf("invalid config: %s", strings.Join(invalid, ", "))
	}

	return nil
}

func validateLocalFileListEntry(entry string, seen map[string]struct{}) error {
	if entry == "" {
		return errors.New("fileList contains blank entries")
	}

	if _, exists := seen[entry]; exists {
		return errors.New("fileList contains duplicates")
	}

	switch entry {
	case goboottypes.ScriptNameMake,
		goboottypes.ScriptNameTask,
		goboottypes.ScriptNameScript,
		goboottypes.ScriptNameCommit:
		return nil
	default:
		return fmt.Errorf("fileList contains unsupported entry %q", entry)
	}
}
