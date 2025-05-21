package goboot

import (
	"fmt"

	"github.com/it-timo/goboot/pkg/config"
)

// Service defines the contract for modular service units that goboot can orchestrate.
//
// Each service corresponds to a specific feature or generation logic.
//
// A service is identified by a stable ID and executing using a matching validated config.
type Service interface {
	// ID returns the stable identifier used to match this service with its configuration.
	ID() string

	// Run executes the service logic using the provided validated configuration.
	//
	// It returns an error if the operation fails.
	Run(cfg config.ServiceConfig) error
}

// serviceManager coordinates service registration and execution.
//
// It holds a registry of enabled service implementations and links them with their corresponding configurations
// via a central config.Manager.
type serviceManager struct {
	// services maps service IDs to their implementation.
	services map[string]Service

	// cfgMgr resolves and holds validated configuration instances by service ID.
	cfgMgr *config.Manager
}

// newServiceManager creates a new ServiceManager bound to the given config manager.
//
// The config manager is expected to be preloaded with validated configurations.
func newServiceManager(cfgMgr *config.Manager) *serviceManager {
	return &serviceManager{
		services: make(map[string]Service),
		cfgMgr:   cfgMgr,
	}
}

// register adds a service implementation to the manager's internal registry.
//
// The service must have a unique and consistent ID.
//
// If a service with the same ID is already registered, it will return an error.
func (sm *serviceManager) register(service Service) error {
	_, ok := sm.services[service.ID()]
	if ok {
		return fmt.Errorf("service %q already registered", service.ID())
	}

	sm.services[service.ID()] = service

	return nil
}

// runAll executes all registered services that have a matching configuration.
//
// For each service:
//   - If a config is available via cfgMgr, the service is executed with it
//   - If no config is found, the service is skipped with a warning
//
// It returns the first encountered execution error, if any. Skipped services do not fail the run.
func (sm *serviceManager) runAll() error {
	for id, svc := range sm.services {
		cfg, ok := sm.cfgMgr.Get(id)
		if !ok {
			fmt.Printf("Service %q skipped (no configuration loaded)\n", id)
			continue
		}

		err := svc.Run(cfg)
		if err != nil {
			return fmt.Errorf("service %q failed: %w", id, err)
		}
	}

	return nil
}
