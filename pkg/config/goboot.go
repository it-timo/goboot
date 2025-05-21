/*
Package config defines and manages the configuration types used by goboot.

It provides structured representations of all YAML-based config files used during project generation,
validation, and templating.

This includes project metadata, generator behavior, feature toggles, and any future modular or plugin-based settings.

Configurations in this package are designed to be loaded from static sources, validated before use,
and injected into templates or internal processing logic.

This package does not deal with application runtime config; its scope is limited to bootstrapping
and template generation inputs.
*/
package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"

	"github.com/it-timo/goboot/pkg/types"

	"gopkg.in/yaml.v3"
)

// GoBoot holds the global context for a goboot run.
//
// It serves as the orchestrator for configuration loading, validation, and access.
//
// This includes:
//   - The base configuration path for the main goboot config
//   - A list of modular config declarations (ServiceConfigMeta)
//   - A central config manager (ConfManager) to register and resolve modules
type GoBoot struct {
	// configPath is the path to the main goboot YAML config file (e.g., ./configs/goboot.yml).
	configPath string

	// ConfManager holds validated and registered configuration modules.
	//
	// It provides access to modular service configs during generation.
	ConfManager *Manager

	// TargetPath is the path to the project target / output.
	TargetPath string `yaml:"target_path"`

	// Services is a list of external service config declarations to load (e.g., base_project, linting).
	Services []ServiceConfigMeta `yaml:"services"`
}

// NewGoBoot creates a new GoBoot instance with the given base configuration path.
//
// It initializes an empty ConfManager for later population.
func NewGoBoot(confPath string) *GoBoot {
	return &GoBoot{
		configPath:  confPath,
		ConfManager: NewConfigManager(),
	}
}

// Init loads and validates the goboot base configuration and all declared service modules.
//
// It performs the following:
//   - Reads and parses the main goboot YAML config
//   - Iterates over declared services and loads their configs via factory
//   - Validates and registers each service config with the ConfManager
func (gb *GoBoot) Init() error {
	err := gb.readConfig()
	if err != nil {
		return fmt.Errorf("failed to read goboot config: %w", err)
	}

	for _, svc := range gb.Services {
		if !svc.IsEnabled() {
			continue
		}

		fmt.Printf("loading service config for %q\n", svc.ID)

		cfg := createServiceConfig(svc.ID)
		if cfg == nil || (reflect.ValueOf(cfg).Kind() == reflect.Ptr && reflect.ValueOf(cfg).IsNil()) {
			return fmt.Errorf("invalid or nil config returned for service ID: %q", svc.ID)
		}

		err = cfg.ReadConfig(svc.ConfPath)
		if err != nil {
			return fmt.Errorf("failed to read config for %q: %w", svc.ID, err)
		}

		err = gb.ConfManager.Register(cfg)
		if err != nil {
			return fmt.Errorf("failed to register config for %q: %w", svc.ID, err)
		}
	}

	return nil
}

// readConfig reads the goboot base configuration from its YAML path
// and unmarshal the values into the current GoBoot struct instance.
func (gb *GoBoot) readConfig() error {
	return readYMLConfig(gb.configPath, gb)
}

// createServiceConfig acts as the central mapping point for service config IDs.
//
// Each known ServiceConfig type must be registered here explicitly.
//
// This maps string identifiers (e.g., "base_project") to their concrete implementations.
//
// Only configs listed here can be used during runtime.
func createServiceConfig(id string) ServiceConfig {
	switch id {
	case types.ServiceNameBaseProject:
		return newBaseProjectConfig()
	// Extend with more cases for additional service types
	default:
		return nil
	}
}

// readYMLConfig reads the given YAML file path and unmarshal it into the provided destination struct.
//
// It returns an error if the file cannot be read or the YAML is malformed.
func readYMLConfig(confPath string, cfg interface{}) error {
	curPath, err := filepath.Abs(path.Clean(confPath))
	if err != nil {
		return fmt.Errorf("failed to resolve config path: %w", err)
	}

	data, err := os.ReadFile(curPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
