/*
Package config defines configuration models and loading helpers for goboot.
It handles scaffold-time config only, not runtime application config.
*/
package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/rs/zerolog"

	"gopkg.in/yaml.v3"
)

// GoBoot holds root config values and the config manager for a scaffold run.
type GoBoot struct {
	// configPath points to the main goboot YAML file.
	configPath string

	// ProjectName is the generated project identifier.
	ProjectName string `yaml:"projectName"`

	// RepoURL is the repository URL used for derived metadata.
	RepoURL string `yaml:"repoUrl"`

	// GitProvider selects provider-specific behavior.
	GitProvider string `yaml:"gitProvider"`

	// TargetPath is the directory that receives generated output.
	TargetPath string `yaml:"targetPath"`

	// Services lists service configs to load.
	Services []ServiceConfigMeta `yaml:"services"`

	// ConfManager stores validated service configs by ID.
	ConfManager *Manager

	logger zerolog.Logger
}

// NewGoBoot returns a GoBoot with an empty config manager.
func NewGoBoot(confPath string) *GoBoot {
	return &GoBoot{
		configPath:  confPath,
		ConfManager: NewConfigManager(),
		logger:      zerolog.Nop(),
	}
}

// SetLogger sets the logger used during config load and validation.
func (gb *GoBoot) SetLogger(logger zerolog.Logger) {
	gb.logger = logger
}

// Init loads the root config, validates it, then loads and registers enabled services.
func (gb *GoBoot) Init() error {
	err := gb.readConfig()
	if err != nil {
		return fmt.Errorf("failed to read goboot config: %w", err)
	}

	err = gb.validateBase()
	if err != nil {
		return fmt.Errorf("invalid goboot config: %w", err)
	}

	for _, svc := range gb.Services {
		if !svc.IsEnabled() {
			continue
		}

		gb.logger.Info().Str("service_id", svc.ID).Str("conf_path", svc.ConfPath).Msg("loading service config")

		cfg := createServiceConfig(svc.ID, gb.ProjectName)
		if cfg == nil {
			return fmt.Errorf("invalid or nil config returned for service ID: %q", svc.ID)
		}

		err = cfg.ReadConfig(svc.ConfPath, gb.RepoURL, gb.GitProvider)
		if err != nil {
			return fmt.Errorf("failed to read config for %q: %w", svc.ID, err)
		}

		err = gb.ConfManager.Register(cfg)
		if err != nil {
			return fmt.Errorf("failed to register config for %q: %w", svc.ID, err)
		}

		gb.logger.Debug().Str("service_id", svc.ID).Msg("service config loaded")
	}

	return nil
}

// readConfig loads the root goboot YAML into gb.
func (gb *GoBoot) readConfig() error {
	return readYMLConfig(gb.configPath, gb)
}

// validateBase validates required root fields and enabled service paths.
//
//nolint:cyclop // flat logic preferred for clarity and extensibility.
func (gb *GoBoot) validateBase() error {
	var missing []string

	if strings.TrimSpace(gb.ProjectName) == "" {
		missing = append(missing, "projectName")
	}

	if strings.TrimSpace(gb.TargetPath) == "" {
		missing = append(missing, "targetPath")
	}

	var importPathMissing bool

	for _, svc := range gb.Services {
		if !svc.IsEnabled() {
			continue
		}

		// Report repoUrl once, even if multiple dependent services are enabled.
		isExempt := svc.ID == goboottypes.ServiceNameBaseProject ||
			svc.ID == goboottypes.ServiceNameBaseLint ||
			svc.ID == goboottypes.ServiceNameBaseTest ||
			svc.ID == goboottypes.ServiceNameBaseCI ||
			svc.ID == goboottypes.ServiceNameBaseLogger ||
			svc.ID == goboottypes.ServiceNameBaseDocker ||
			svc.ID == goboottypes.ServiceNameBaseRelease

		if !importPathMissing && !isExempt {
			if strings.TrimSpace(gb.RepoURL) == "" {
				importPathMissing = true

				missing = append(missing, "repoUrl")

				continue
			}
		}

		if strings.TrimSpace(svc.ConfPath) == "" {
			missing = append(missing, fmt.Sprintf("services[%s].confPath", svc.ID))
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}

	err := validateProjectName(gb.ProjectName)
	if err != nil {
		return err
	}

	return nil
}

// createServiceConfig maps a service ID to its concrete config implementation.
func createServiceConfig(id, projectName string) ServiceConfig {
	switch id {
	case goboottypes.ServiceNameBaseProject:
		return newBaseProjectConfig(projectName)
	case goboottypes.ServiceNameBaseLint:
		return newBaseLintConfig(projectName)
	case goboottypes.ServiceNameBaseLocal:
		return newBaseLocalConfig(projectName)
	case goboottypes.ServiceNameBaseTest:
		return newBaseTestConfig(projectName)
	case goboottypes.ServiceNameBaseCI:
		return newBaseCIConfig(projectName)
	case goboottypes.ServiceNameBaseLogger:
		return newBaseLoggerConfig(projectName)
	case goboottypes.ServiceNameBaseDocker:
		return newBaseDockerConfig(projectName)
	case goboottypes.ServiceNameBaseRelease:
		return newBaseReleaseConfig(projectName)
	// Extend with more cases for additional service types.
	default:
		return nil
	}
}

// readYMLConfig reads confPath and unmarshals YAML into cfg.
func readYMLConfig(confPath string, cfg interface{}) error {
	curPath, err := filepath.Abs(path.Clean(confPath))
	if err != nil {
		return fmt.Errorf("failed to resolve config path: %w", err)
	}

	data, err := os.ReadFile(curPath) // #nosec G304 -- the path is user-defined and expected to be dynamic.
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	err = decoder.Decode(cfg)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
