package config

import (
	"fmt"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

// ServiceConfig is the common contract for loadable service configs.
type ServiceConfig interface {
	// ID returns the config identifier.
	ID() string

	// ReadConfig loads config from file and repository metadata.
	ReadConfig(confPath string, repoURL string, gitProvider string) error

	// Validate returns an error when config is incomplete or invalid.
	Validate() error
}

// ServiceConfigMeta declares a service config file to load.
type ServiceConfigMeta struct {
	ID       string `yaml:"id"`       // e.g., "base_project"
	ConfPath string `yaml:"confPath"` // e.g., "./configs/base_project.yml"
	Enabled  bool   `yaml:"enabled"`
}

// IsEnabled returns the enabled state.
func (scm *ServiceConfigMeta) IsEnabled() bool {
	return scm.Enabled
}

// Manager stores validated service configs for service and registrar lookups.
type Manager struct {
	services   map[string]ServiceConfig
	registrars map[string]ServiceConfig
}

// NewConfigManager constructs an empty config manager.
func NewConfigManager() *Manager {
	return &Manager{
		services:   make(map[string]ServiceConfig),
		registrars: make(map[string]ServiceConfig),
	}
}

// Register validates cfg and stores it in service or registrar map by ID.
func (cm *Manager) Register(cfg ServiceConfig) error {
	err := cfg.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	switch cfg.ID() {
	case goboottypes.ServiceNameBaseProject:
		cm.registrars[cfg.ID()] = cfg
	default:
		cm.services[cfg.ID()] = cfg
	}

	return nil
}

// UnregisterService removes a service config by ID.
func (cm *Manager) UnregisterService(id string) {
	delete(cm.services, id)
}

// UnregisterRegistrar removes a registrar config by ID.
func (cm *Manager) UnregisterRegistrar(id string) {
	delete(cm.registrars, id)
}

// GetService returns a service config and a found flag.
func (cm *Manager) GetService(id string) (ServiceConfig, bool) {
	cfg, ok := cm.services[id]

	return cfg, ok
}

// GetRegistrar returns a registrar config and a found flag.
func (cm *Manager) GetRegistrar(id string) (ServiceConfig, bool) {
	cfg, ok := cm.registrars[id]

	return cfg, ok
}
