package config

import "fmt"

// ServiceConfig represents a modular configuration component used by goboot.
//
// Each configuration module (such as base project settings, linting, or Docker support)
// should implement this interface to allow validation and identification within the config manager.
type ServiceConfig interface {
	// ID returns a stable identifier for the configuration module.
	//
	// Example: "base_project", "linting", "docker".
	ID() string

	// ReadConfig loads the configuration from a source file.
	ReadConfig(confPath string) error

	// Validate verifies that the configuration is complete and semantically correct.
	//
	// It returns an error if the configuration is invalid.
	Validate() error
}

// ServiceConfigMeta represents a declaration of a modular config block to load.
//
// It is not the config itself, but a registry of what to load from which file.
type ServiceConfigMeta struct {
	// ID is the stable identifier for the config module.
	ID string `yaml:"id"` // e.g., "base_project"
	// ConfPath is the path to the config file to load.
	ConfPath string `yaml:"conf_path"` // e.g., "./configs/base_project.yml"
	// Enabled indicates whether the service should be enabled.
	Enabled bool `yaml:"enabled"`
}

// IsEnabled returns the enabled state.
func (scm *ServiceConfigMeta) IsEnabled() bool {
	return scm.Enabled
}

// Manager provides centralized registration and retrieval of modular ServiceConfig implementations.
//
// It allows goboot to dynamically register, validate, and access multiple configuration modules
// without hard-coding their types or structure.
//
// This supports future extensibility as new config types are introduced.
type Manager struct {
	services map[string]ServiceConfig
}

// NewConfigManager returns a new instance of Manager with an initialized internal registry.
func NewConfigManager() *Manager {
	return &Manager{services: make(map[string]ServiceConfig)}
}

// Register adds a ServiceConfig to the manager after validating it.
//
// If validation fails, the configuration is not registered and an error is returned.
func (cm *Manager) Register(cfg ServiceConfig) error {
	err := cfg.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	cm.services[cfg.ID()] = cfg

	return nil
}

// Unregister removes a registered ServiceConfig by its ID.
//
// No-op if the ID is not found.
func (cm *Manager) Unregister(id string) {
	delete(cm.services, id)
}

// Get retrieves a registered ServiceConfig by its ID.
//
// The second return value indicates whether the config was found.
func (cm *Manager) Get(id string) (ServiceConfig, bool) {
	cfg, ok := cm.services[id]

	return cfg, ok
}
