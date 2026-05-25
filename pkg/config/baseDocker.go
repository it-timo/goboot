package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

var (
	dockerBinaryNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
	dockerImageRefPattern   = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]*` +
		`(?::[A-Za-z0-9_][A-Za-z0-9._-]{0,127})?(?:@sha256:[A-Fa-f0-9]{64})?\z`)
	dockerMainPackagePattern = regexp.MustCompile(`^(\.|./[A-Za-z0-9._/-]+)$`)
	dockerPortPattern        = regexp.MustCompile(`^[0-9]+(:[0-9]+)?$`)
)

// BaseDockerConfig contains inputs for containerization template generation.
type BaseDockerConfig struct {
	// SourcePath points to Docker templates.
	SourcePath string `yaml:"sourcePath"`

	// GoVersion is the Go image tag used by the build stage.
	GoVersion string `yaml:"goVersion"`

	// RuntimeImage is the final-stage image.
	RuntimeImage string `yaml:"runtimeImage"`

	// BinaryName is the produced binary name.
	BinaryName string `yaml:"binaryName"`

	// MainPackage is the Go package built into the container.
	MainPackage string `yaml:"mainPackage"`

	// FileList lists container files to generate.
	FileList []string `yaml:"fileList"`

	// Ports lists optional docker-compose port mappings.
	Ports []string `yaml:"ports"`

	// ProjectName is the project identifier.
	ProjectName string `yaml:"-"`

	// LowerProjectName is the lowercase variant of ProjectName.
	LowerProjectName string `yaml:"-"`

	// ImageName is the local container image name.
	ImageName string `yaml:"-"`
}

// newBaseDockerConfig creates a BaseDockerConfig with the given project name.
func newBaseDockerConfig(projectName string) *BaseDockerConfig {
	return &BaseDockerConfig{
		ProjectName: projectName,
	}
}

// ID returns a stable identifier for this config.
func (bd *BaseDockerConfig) ID() string {
	return goboottypes.ServiceNameBaseDocker
}

// ReadConfig loads base_docker config from confPath.
func (bd *BaseDockerConfig) ReadConfig(confPath string, _ string, _ string) error {
	return readYMLConfig(confPath, bd)
}

// Validate checks required fields and normalizes defaults.
func (bd *BaseDockerConfig) Validate() error {
	var missing []string

	if strings.TrimSpace(bd.SourcePath) == "" {
		missing = append(missing, "sourcePath")
	}

	if strings.TrimSpace(bd.ProjectName) == "" {
		missing = append(missing, "projectName")
	}

	if len(bd.FileList) == 0 {
		missing = append(missing, "fileList")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config fields: %s", strings.Join(missing, ", "))
	}

	err := validateProjectName(bd.ProjectName)
	if err != nil {
		return err
	}

	bd.fillDefaults()

	err = bd.validateFileList()
	if err != nil {
		return err
	}

	err = bd.validateValues()
	if err != nil {
		return err
	}

	return bd.validatePorts()
}

func (bd *BaseDockerConfig) fillDefaults() {
	bd.LowerProjectName = strings.ToLower(bd.ProjectName)
	bd.GoVersion = strings.TrimSpace(bd.GoVersion)
	bd.RuntimeImage = strings.TrimSpace(bd.RuntimeImage)
	bd.BinaryName = strings.TrimSpace(bd.BinaryName)
	bd.MainPackage = strings.TrimSpace(bd.MainPackage)

	if strings.TrimSpace(bd.GoVersion) == "" {
		bd.GoVersion = "1.26.3"
	}

	if strings.TrimSpace(bd.RuntimeImage) == "" {
		bd.RuntimeImage = "gcr.io/distroless/static-debian12:nonroot"
	}

	if strings.TrimSpace(bd.BinaryName) == "" {
		bd.BinaryName = bd.LowerProjectName
	}

	if strings.TrimSpace(bd.MainPackage) == "" {
		bd.MainPackage = "./cmd/" + bd.LowerProjectName
	}

	bd.ImageName = bd.LowerProjectName + ":local"
}

func (bd *BaseDockerConfig) validateValues() error {
	if !goVersionPattern.MatchString(bd.GoVersion) {
		return fmt.Errorf("invalid config: goVersion %q is not supported", bd.GoVersion)
	}

	if !dockerImageRefPattern.MatchString(bd.RuntimeImage) {
		return fmt.Errorf("invalid config: runtimeImage %q is not supported", bd.RuntimeImage)
	}

	if !dockerBinaryNamePattern.MatchString(bd.BinaryName) {
		return fmt.Errorf("invalid config: binaryName %q is not supported", bd.BinaryName)
	}

	if strings.Contains(bd.MainPackage, "..") ||
		strings.Contains(bd.MainPackage, "//") ||
		!dockerMainPackagePattern.MatchString(bd.MainPackage) {
		return fmt.Errorf("invalid config: mainPackage %q is not supported", bd.MainPackage)
	}

	return nil
}

func (bd *BaseDockerConfig) validateFileList() error {
	seen := make(map[string]struct{})
	invalid := []string{}

	for _, file := range bd.FileList {
		trimmed := strings.TrimSpace(file)

		err := validateDockerFileListEntry(trimmed, seen)
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

func validateDockerFileListEntry(entry string, seen map[string]struct{}) error {
	if entry == "" {
		return errors.New("fileList contains blank entries")
	}

	if _, exists := seen[entry]; exists {
		return errors.New("fileList contains duplicates")
	}

	switch entry {
	case "dockerfile", "compose", "dockerignore":
		return nil
	default:
		return fmt.Errorf("fileList contains unsupported entry %q", entry)
	}
}

func (bd *BaseDockerConfig) validatePorts() error {
	for _, port := range bd.Ports {
		trimmed := strings.TrimSpace(port)
		if trimmed == "" {
			return errors.New("invalid config: ports contains empty string")
		}

		if !dockerPortPattern.MatchString(trimmed) {
			return fmt.Errorf("invalid config: port %q is not supported", port)
		}
	}

	return nil
}
