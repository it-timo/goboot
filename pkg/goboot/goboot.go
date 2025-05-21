/*
Package goboot defines the core orchestration layer of the `goboot` CLI tool.

This package coordinates the loading, registration, and execution of modular generation services.
It acts as the top-level entry point, delegating control to a service manager that runs logic
based on the configuration declared in a YAML config file (see config.GoBoot).

Does not implement generation logic — instead, it dynamically wires together service modules (e.g., baseProject)
that encapsulate feature-specific behavior.
*/
package goboot

import (
	"fmt"
	"os"

	"github.com/it-timo/goboot/pkg/baseProject"
	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/types"
)

// GoBoot is the central controller struct for running goboot-based generation logic.
//
// It is initialized with the validated configuration (config.GoBoot)
// and uses a service manager to conditionally load and execute service implementations.
type GoBoot struct {
	// cfg holds the loaded and validated configuration contexts.
	cfg *config.GoBoot

	// ServiceMgr manages the lifecycle and execution of registered service modules.
	ServiceMgr *serviceManager
}

// NewGoBoot creates and returns a new GoBoot instance bound to the provided configuration.
//
// It wires the internal service manager to the configuration's ConfManager.
//
// Note: This does not yet load services — use RegisterServices() afterward.
func NewGoBoot(config *config.GoBoot) *GoBoot {
	return &GoBoot{
		cfg:        config,
		ServiceMgr: newServiceManager(config.ConfManager),
	}
}

// RegisterServices evaluates all declared services in the config and registers only those marked as enabled.
//
// Each service must be explicitly handled here by matching its ID.
//
// This avoids runtime registration logic and keeps service orchestration predictable.
//
// It performs the following for each declared and enabled service:
//   - Checks the service ID
//   - Instantiates the appropriate service implementation
//   - Registers the service with the service manager
//
// Returns an error if any declared service is unknown or registration fails.
func (gb *GoBoot) RegisterServices() error {
	if gb.cfg.Services == nil {
		return fmt.Errorf("no services declared in config")
	}

	// creates the target dir if not exist.
	err := os.MkdirAll(gb.cfg.TargetPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	for _, meta := range gb.cfg.Services {
		if !meta.IsEnabled() {
			continue
		}

		fmt.Printf("loading service %s\n", meta.ID)

		switch meta.ID {
		case types.ServiceNameBaseProject:
			err = gb.ServiceMgr.register(baseProject.NewBaseProject(gb.cfg.TargetPath))
			if err != nil {
				return fmt.Errorf("failed to register base_project service: %w", err)
			}
		// Future services can be added here.
		default:
			return fmt.Errorf("unknown service ID: %s", meta.ID)
		}
	}

	return nil
}

// RunServices executes all registered services in the order they were registered.
//
// It delegates to the internal service manager, which pulls the appropriate config
// for each service and invokes its logic.
//
// If a service has no config, it is skipped.
func (gb *GoBoot) RunServices() error {
	return gb.ServiceMgr.runAll()
}
