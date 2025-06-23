package goboot

import (
	"fmt"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/types"
)

// Service defines the contract for modular service units that goboot can orchestrate.
//
// Each service corresponds to a specific feature or generation logic.
//
// A service is identified by a stable ID and executing using a matching validated config.
type Service interface {
	// ID returns the stable identifier used to match this service with its configuration.
	ID() string

	// SetConfig assigns the service config using the provided validated configuration.
	//
	// It returns an error if the operation fails.
	SetConfig(cfg config.ServiceConfig) error

	// Run executes the service logic.
	//
	// It returns an error if the operation fails.
	Run() error
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

	// priorServiceIDs defines the core services that must always run first,
	// such as directory structure or root project setup.
	priorServiceIDs []string

	// subsequentServiceIDs defines the care services that must always run last,
	// such as script injections for multiple services.
	subsequentServiceIDs []string
}

// newServiceManager creates a new ServiceManager bound to the given config manager.
//
// The config manager is expected to be preloaded with validated configurations.
func newServiceManager(cfgMgr *config.Manager) *serviceManager {
	return &serviceManager{
		services: make(map[string]Service),
		cfgMgr:   cfgMgr,
		priorServiceIDs: []string{
			types.ServiceNameBaseProject, // required to initialize the base project directory structure.
			// future pre-services can be added here.
		},
		subsequentServiceIDs: []string{
			types.ServiceNameBaseLocal, // required to handle the script files in fully.
			// future sub-services can be added here.
		},
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
//nolint:cyclop // branching required for service dispatch logic; each path reflects a distinct lifecycle phase.
func (sm *serviceManager) runAll() error {
	err := sm.assignConfigs()
	if err != nil {
		return fmt.Errorf("failed to assign configs: %w", err)
	}

	// Run priority bootstrap services first (e.g., base_project).
	err = sm.runPriorServices()
	if err != nil {
		return fmt.Errorf("failed to run prior services: %w", err)
	}

	for curID, svc := range sm.services {
		// Skip services already handled in runPriorServices or will be handled by runSubsequentServices.
		if sm.isPriorService(curID) || sm.isSubsequentService(curID) {
			continue
		}

		_, ok := sm.cfgMgr.GetRegistrar(curID)
		if !ok {
			_, ok = sm.cfgMgr.GetService(curID)
			if !ok {
				fmt.Printf("Service %q skipped (no configuration loaded)\n", curID)

				continue
			}
		}

		receiver, isScriptRec := svc.(types.ScriptReceiver)
		if isScriptRec {
			registrar, isRegistrar := sm.services[types.ServiceNameBaseLocal].(types.Registrar)
			if isRegistrar {
				fmt.Printf("Injecting script registrar into %q\n", curID)

				receiver.SetScriptReceiver(registrar)
			}
		}

		err = svc.Run()
		if err != nil {
			return fmt.Errorf("failed to run service %q: %w", curID, err)
		}
	}

	// Run later bootstrap services first (e.g., base_local).
	err = sm.runSubsequentServices()
	if err != nil {
		return fmt.Errorf("failed to run subsequent services: %w", err)
	}

	return nil
}

// assignConfigs calls on all registered services with a valid config the SetConfig.
//
// returns an error if any set fails.
func (sm *serviceManager) assignConfigs() error {
	for curID, svc := range sm.services {
		cfg, ok := sm.cfgMgr.GetRegistrar(curID)
		if !ok {
			cfg, ok = sm.cfgMgr.GetService(curID)
			if !ok {
				continue
			}
		}

		err := svc.SetConfig(cfg)
		if err != nil {
			return fmt.Errorf("failed to set config for %q: %w", curID, err)
		}
	}

	return nil
}

// isPriorService checks if the given service ID is in the list of priority services.
func (sm *serviceManager) isPriorService(id string) bool {
	for _, ps := range sm.priorServiceIDs {
		if ps == id {
			return true
		}
	}

	return false
}

// isSubsequentService checks if the given service ID is in the list of later services.
func (sm *serviceManager) isSubsequentService(id string) bool {
	for _, ps := range sm.subsequentServiceIDs {
		if ps == id {
			return true
		}
	}

	return false
}

// runPriorServices executes predefined core services that must run before any other services.
//
// This is typically used to guarantee foundational setup (e.g., base_project) is completed
// before rendering additional modules.
//
// The list of service IDs is hardcoded in a dedicated slice to allow future extension.
func (sm *serviceManager) runPriorServices() error {
	for _, serviceID := range sm.priorServiceIDs {
		svc, okay := sm.services[serviceID]
		if !okay {
			fmt.Printf("Service %q skipped (not registered)\n", serviceID)

			continue
		}

		_, ok := sm.cfgMgr.GetRegistrar(serviceID)
		if !ok {
			_, ok = sm.cfgMgr.GetService(serviceID)
			if !ok {
				fmt.Printf("Prior Service %q skipped (no configuration loaded)\n", serviceID)

				continue
			}
		}

		err := svc.Run()
		if err != nil {
			return fmt.Errorf("failed to run service %q: %w", serviceID, err)
		}
	}

	return nil
}

// runSubsequentServices executes predefined care services that must run after any other services.
//
// This is typically used to guarantee foundational setup (e.g., base_local) is completed
// before rendering this module.
//
// The list of service IDs is hardcoded in a dedicated slice to allow future extension.
func (sm *serviceManager) runSubsequentServices() error {
	for _, serviceID := range sm.subsequentServiceIDs {
		svc, okay := sm.services[serviceID]
		if !okay {
			fmt.Printf("Service %q skipped (not registered)\n", serviceID)

			continue
		}

		_, ok := sm.cfgMgr.GetRegistrar(serviceID)
		if !ok {
			_, ok = sm.cfgMgr.GetService(serviceID)
			if !ok {
				fmt.Printf("Subsequent Service %q skipped (no configuration loaded)\n", serviceID)

				continue
			}
		}

		err := svc.Run()
		if err != nil {
			return fmt.Errorf("failed to run service %q: %w", serviceID, err)
		}
	}

	return nil
}
