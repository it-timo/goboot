/*
Package goboot provides top-level service orchestration.
*/
package goboot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/it-timo/goboot/pkg/baseci"
	"github.com/it-timo/goboot/pkg/baselint"
	"github.com/it-timo/goboot/pkg/baselocal"
	"github.com/it-timo/goboot/pkg/baseproject"
	"github.com/it-timo/goboot/pkg/basetest"
	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

// GoBoot orchestrates service registration and execution for one scaffold run.
type GoBoot struct {
	cfg *config.GoBoot

	// ServiceMgr manages service order and lifecycle.
	ServiceMgr *serviceManager
}

// NewGoBoot returns a GoBoot wired to the provided config.
func NewGoBoot(config *config.GoBoot) *GoBoot {
	return &GoBoot{
		cfg:        config,
		ServiceMgr: newServiceManager(config.ConfManager),
	}
}

// RegisterServices ensures target root exists and registers enabled services.
func (gb *GoBoot) RegisterServices() error {
	if gb.cfg.Services == nil {
		return errors.New("no services declared in config")
	}

	// Ensure target root exists.
	err := os.MkdirAll(gb.cfg.TargetPath, goboottypes.DirPerm)
	if err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	err = gb.registerPreServices()
	if err != nil {
		return fmt.Errorf("failed to register pre services: %w", err)
	}

	err = gb.registerMainServices()
	if err != nil {
		return fmt.Errorf("failed to register main services: %w", err)
	}

	return nil
}

// RunServices executes registered services in manager-defined order.
func (gb *GoBoot) RunServices() error {
	return gb.ServiceMgr.runAll()
}

// RunGoModTidy runs `go mod tidy` when enabled and go.mod exists.
func (gb *GoBoot) RunGoModTidy(execute bool) error {
	if !execute {
		return nil
	}

	projectRoot := filepath.Join(gb.cfg.TargetPath, gb.cfg.ProjectName)
	goModPath := filepath.Join(projectRoot, "go.mod")

	_, err := os.Stat(goModPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("failed to stat go.mod: %w", err)
	}

	// Run in generated project root.
	cmd := exec.CommandContext(context.Background(), "go", "mod", "tidy")
	cmd.Dir = projectRoot

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	return nil
}

// registerPreServices registers services that other services depend on.
func (gb *GoBoot) registerPreServices() error {
	for _, meta := range gb.cfg.Services {
		if !meta.IsEnabled() {
			continue
		}

		switch meta.ID {
		case goboottypes.ServiceNameBaseCI:
			baseCI := baseci.NewBaseCI(gb.cfg.TargetPath)

			err := gb.ServiceMgr.register(baseCI)
			if err != nil {
				return fmt.Errorf("failed to register %s service: %w", goboottypes.ServiceNameBaseCI, err)
			}
		case goboottypes.ServiceNameBaseLocal:
			baseLocal := baselocal.NewBaseLocal(gb.cfg.TargetPath)

			err := gb.ServiceMgr.register(baseLocal)
			if err != nil {
				return fmt.Errorf("failed to register %s service: %w", goboottypes.ServiceNameBaseLocal, err)
			}
		default:
			// Ignore non-pre services here.
			continue
		}

		fmt.Printf("loaded pre-service %s\n", meta.ID)
	}

	return nil
}

// registerMainServices registers all non-pre services.
//nolint:cyclop // Flat switch is preferred for explicit control and traceability.
func (gb *GoBoot) registerMainServices() error {
	for _, meta := range gb.cfg.Services {
		if !meta.IsEnabled() {
			continue
		}

		switch meta.ID {
		case goboottypes.ServiceNameBaseProject:
			err := gb.ServiceMgr.register(baseproject.NewBaseProject(gb.cfg.TargetPath))
			if err != nil {
				return fmt.Errorf("failed to register %s service: %w", goboottypes.ServiceNameBaseProject, err)
			}
		case goboottypes.ServiceNameBaseLint:
			err := gb.ServiceMgr.register(baselint.NewBaseLint(gb.cfg.TargetPath))
			if err != nil {
				return fmt.Errorf("failed to register %s service: %w", goboottypes.ServiceNameBaseLint, err)
			}
		case goboottypes.ServiceNameBaseLocal:
			// Registered earlier as a pre service.
			continue
		case goboottypes.ServiceNameBaseCI:
			// Registered earlier as a pre service.
			continue
		case goboottypes.ServiceNameBaseTest:
			err := gb.ServiceMgr.register(basetest.NewBaseTest(gb.cfg.TargetPath))
			if err != nil {
				return fmt.Errorf("failed to register %s service: %w", goboottypes.ServiceNameBaseTest, err)
			}
		// Future services can be added here.
		default:
			return fmt.Errorf("unknown service ID: %s", meta.ID)
		}

		fmt.Printf("loaded service %s\n", meta.ID)
	}

	return nil
}
