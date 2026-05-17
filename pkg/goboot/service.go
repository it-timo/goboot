package goboot

import (
	"fmt"
	"time"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/rs/zerolog"
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

type loggerAware interface {
	SetLogger(logger zerolog.Logger)
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

	log zerolog.Logger
}

// newServiceManager creates a service manager bound to cfgMgr.
func newServiceManager(cfgMgr *config.Manager, logger zerolog.Logger) *serviceManager {
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
		log: logger,
	}
}

// register adds a service implementation and rejects duplicate IDs.
func (sm *serviceManager) register(service Service) error {
	_, ok := sm.services[service.ID()]
	if ok {
		return fmt.Errorf("service %q already registered", service.ID())
	}

	if receiver, ok := service.(loggerAware); ok {
		receiver.SetLogger(sm.log.With().Str("service_id", service.ID()).Logger())
	}

	sm.services[service.ID()] = service
	sm.log.Info().Str("service_id", service.ID()).Msg("service attached to manager")

	return nil
}

// runAll assigns configs and executes services in prior -> regular -> subsequent order.
// Services without config are skipped.
func (sm *serviceManager) runAll() error {
	sm.log.Info().Int("service_count", len(sm.services)).Msg("starting service manager run")

	err := sm.assignConfigs()
	if err != nil {
		return fmt.Errorf("failed to assign configs: %w", err)
	}

	err = sm.runPriorServices()
	if err != nil {
		return fmt.Errorf("failed to run prior services: %w", err)
	}

	err = sm.runRegularServices()
	if err != nil {
		return fmt.Errorf("failed to run regular services: %w", err)
	}

	err = sm.runSubsequentServices()
	if err != nil {
		return fmt.Errorf("failed to run subsequent services: %w", err)
	}

	return nil
}

func (sm *serviceManager) runRegularServices() error {
	for curID, svc := range sm.services {
		if sm.isPriorService(curID) || sm.isSubsequentService(curID) {
			continue
		}

		if !sm.hasConfig(curID) {
			sm.log.Debug().Str("service_id", curID).Msg("service skipped because no configuration was loaded")

			continue
		}

		sm.injectRegistrars(curID, svc)

		err := sm.runService(curID, svc, "service")
		if err != nil {
			return err
		}
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
				sm.log.Debug().Str("service_id", curID).Msg("no config found while assigning service configs")

				continue
			}
		}

		err := svc.SetConfig(cfg)
		if err != nil {
			return fmt.Errorf("failed to set config for %q: %w", curID, err)
		}

		if receiver, ok := svc.(goboottypes.LoggerSettingsReceiver); ok {
			loggerCfg, found := sm.cfgMgr.GetService(goboottypes.ServiceNameBaseLogger)
			if found {
				provider, providerOK := loggerCfg.(goboottypes.LoggerSettingsProvider)
				if !providerOK {
					return fmt.Errorf("logger config %q does not provide logger settings", goboottypes.ServiceNameBaseLogger)
				}

				receiver.SetLoggerSettings(provider.LoggerSettings())
			}
		}

		sm.log.Debug().Str("service_id", curID).Msg("service configuration assigned")
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
	return sm.runOrderedServices(sm.priorServiceIDs, "prior")
}

// runSubsequentServices executes configured subsequent-phase services.
func (sm *serviceManager) runSubsequentServices() error {
	return sm.runOrderedServices(sm.subsequentServiceIDs, "subsequent")
}

func (sm *serviceManager) runOrderedServices(serviceIDs []string, phase string) error {
	for _, serviceID := range serviceIDs {
		svc, okay := sm.services[serviceID]
		if !okay {
			sm.log.Debug().Str("service_id", serviceID).Str("phase", phase).Msg("service skipped because it is not registered")

			continue
		}

		if !sm.hasConfig(serviceID) {
			sm.log.Debug().
				Str("service_id", serviceID).
				Str("phase", phase).
				Msg("service skipped because no configuration was loaded")

			continue
		}

		err := sm.runService(serviceID, svc, phase+" service")
		if err != nil {
			return err
		}
	}

	return nil
}

func (sm *serviceManager) hasConfig(serviceID string) bool {
	_, registrarFound := sm.cfgMgr.GetRegistrar(serviceID)
	if registrarFound {
		return true
	}

	_, serviceFound := sm.cfgMgr.GetService(serviceID)

	return serviceFound
}

func (sm *serviceManager) injectRegistrars(serviceID string, svc Service) {
	receiver, receivesScripts := svc.(goboottypes.ScriptReceiver)
	if receivesScripts {
		registrar, hasRegistrar := sm.services[goboottypes.ServiceNameBaseLocal].(goboottypes.Registrar)
		if hasRegistrar {
			sm.log.Debug().Str("service_id", serviceID).Msg("injecting script registrar")

			receiver.SetScriptReceiver(registrar)
		}
	}

	ciReceiver, receivesCI := svc.(goboottypes.CIReceiver)
	if receivesCI {
		registrar, hasRegistrar := sm.services[goboottypes.ServiceNameBaseCI].(goboottypes.Registrar)
		if hasRegistrar {
			sm.log.Debug().Str("service_id", serviceID).Msg("injecting CI registrar")

			ciReceiver.SetCIReceiver(registrar)
		}
	}
}

func (sm *serviceManager) runService(serviceID string, svc Service, label string) error {
	sm.log.Info().Str("service_id", serviceID).Msg("running " + label)

	start := time.Now()

	err := svc.Run()
	if err != nil {
		return fmt.Errorf("failed to run service %q: %w", serviceID, err)
	}

	sm.log.Info().
		Str("service_id", serviceID).
		Int64("duration_ms", time.Since(start).Milliseconds()).
		Msg(label + " completed")

	return nil
}
