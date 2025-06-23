package config

import (
	"fmt"
	"strings"

	"github.com/it-timo/goboot/pkg/types"
)

// BaseLocalConfig defines the metadata used by goboot to generate the local setup for a project.
//
// It injects values into templates (e.g., Makefile) and governs how project-specific scripts are rendered.
type BaseLocalConfig struct {
	// SourcePath is the path to the template source directory (e.g., "./templates/local_base").
	SourcePath string `yaml:"sourcePath"`

	// ProjectName is the short identifier for the project (e.g., "goboot").
	// Used in headings, comments, and other rendered metadata.
	ProjectName string `yaml:"projectName"`

	// FileList is a list of files to be copied from the source path to the target path.
	FileList []string `yaml:"fileList"`
}

// newBaseLocalConfig returns a newly initialized BaseLocalConfig with the project name.
func newBaseLocalConfig(projectName string) *BaseLocalConfig {
	return &BaseLocalConfig{
		ProjectName: projectName,
	}
}

// ID returns a stable identifier for this config.
func (bl *BaseLocalConfig) ID() string {
	return types.ServiceNameBaseLocal
}

// ReadConfig loads the base local configuration from the provided YAML file path.
//
// It overwrites the current config values with the file contents.
func (bl *BaseLocalConfig) ReadConfig(confPath string) error {
	return readYMLConfig(confPath, bl)
}

// Validate verifies the BaseLocalConfig for use in scaffolding.
//
// It returns an error if required values are missing/invalid, or calls fillNeededInfos.
func (bl *BaseLocalConfig) Validate() error {
	var missing []string

	if strings.TrimSpace(bl.SourcePath) == "" {
		missing = append(missing, "sourcePath")
	}

	if strings.TrimSpace(bl.ProjectName) == "" {
		missing = append(missing, "projectName")
	}

	// do not allow empty file list - the service should be disabled if no files are copied.
	if len(bl.FileList) == 0 {
		missing = append(missing, "fileList")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config fields: %s", strings.Join(missing, ", "))
	}

	return nil
}
