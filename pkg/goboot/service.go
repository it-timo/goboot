package goboot

import (
	"fmt"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

// Service defines the lifecycle contract for generator services.
type Service interface {
	// ID returns the stable service identifier.
	ID() string

	// SetConfig assigns validated service config.
	SetConfig(cfg config.ServiceConfig) error

	// Run executes the service.
	Run() error
}

// serviceManager coordinates registration, config assignment, and run order.
type serviceManager struct {
	// services maps service IDs to implementations.
	services map[string]Service

	// cfgMgr resolves validated service configs by ID.
	cfgMgr *config.Manager

	// priorServiceIDs run before regular services.
	priorServiceIDs []string

	// subsequentServiceIDs run after regular services.
	subsequentServiceIDs []string
}

// newServiceManager creates a service manager bound to cfgMgr.
func newServiceManager(cfgMgr *config.Manager) *serviceManager {
	return &serviceManager{
		services: make(map[string]Service),
		cfgMgr:   cfgMgr,
		priorServiceIDs: []string{
			goboottypes.ServiceNameBaseProject, // required to initialize the base project directory structure.
			// Future pre services can be added here.
		},
		subsequentServiceIDs: []string{
			goboottypes.ServiceNameBaseCI,    // required to render aggregated CI files.
			goboottypes.ServiceNameBaseLocal, // required to render aggregated local scripts.
			// Future subsequent services can be added here.
		},
	}
}

// register adds a service implementation and rejects duplicate IDs.
func (sm *serviceManager) register(service Service) error {
	_, ok := sm.services[service.ID()]
	if ok {
		return fmt.Errorf("service %q already registered", service.ID())
	}

	sm.services[service.ID()] = service

	return nil
}

// runAll assigns configs and executes services in prior -> regular -> subsequent order.
// Services without config are skipped.
//
//nolint:cyclop // branching required for service dispatch logic; each path reflects a distinct lifecycle phase.
func (sm *serviceManager) runAll() error {
	err := sm.assignConfigs()
	if err != nil {
		return fmt.Errorf("failed to assign configs: %w", err)
	}

	// Run priority services first (for example base_project).
	err = sm.runPriorServices()
	if err != nil {
		return fmt.Errorf("failed to run prior services: %w", err)
	}

	for curID, svc := range sm.services {
		// Skip services handled by prior/subsequent phases.
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

		receiver, isScriptRec := svc.(goboottypes.ScriptReceiver)
		if isScriptRec {
			registrar, isRegistrar := sm.services[goboottypes.ServiceNameBaseLocal].(goboottypes.Registrar)
			if isRegistrar {
				fmt.Printf("Injecting script registrar into %q\n", curID)

				receiver.SetScriptReceiver(registrar)
			}
		}

		ciReceiver, isCIReceiver := svc.(goboottypes.CIReceiver)
		if isCIReceiver {
			registrar, isRegistrar := sm.services[goboottypes.ServiceNameBaseCI].(goboottypes.Registrar)
			if isRegistrar {
				fmt.Printf("Injecting ci registrar into %q\n", curID)

				ciReceiver.SetCIReceiver(registrar)
			}
		}

		err = svc.Run()
		if err != nil {
			return fmt.Errorf("failed to run service %q: %w", curID, err)
		}
	}

	// Run subsequent services last (for example base_local).
	err = sm.runSubsequentServices()
	if err != nil {
		return fmt.Errorf("failed to run subsequent services: %w", err)
	}

	return nil
}

// assignConfigs calls SetConfig for each registered service with loaded config.
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

// isPriorService reports whether id is in the prior phase.
func (sm *serviceManager) isPriorService(id string) bool {
	for _, ps := range sm.priorServiceIDs {
		if ps == id {
			return true
		}
	}

	return false
}

// isSubsequentService reports whether id is in the subsequent phase.
func (sm *serviceManager) isSubsequentService(id string) bool {
	for _, ps := range sm.subsequentServiceIDs {
		if ps == id {
			return true
		}
	}

	return false
}

// runPriorServices executes configured prior-phase services.
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

// runSubsequentServices executes configured subsequent-phase services.
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
