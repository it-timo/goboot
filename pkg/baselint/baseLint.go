/*
Package baselint implements the core bootstrapping logic for the "base_lint" service.

This service is responsible for setting up the foundational linting infrastructure of a newly scaffolded Go project.

It handles:
  - Copying predefined linter configuration templates into the project.
  - Applying Go `text/template` rendering to inject project-specific metadata.
  - Skipping linters that are disabled in the config.
  - Respecting a strict separation of logic per linter (Go, YAML, Markdown, Make).

The service expects a validated configuration of type config.BaseLintConfig
and is one of the default built-in services within the `goboot` system.
*/
package baselint

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/types"
	"github.com/it-timo/goboot/pkg/utils"
)

// BaseLint implements the Service interface and encapsulates the execution logic
// for generating linter configuration files based on user-defined config.
//
// It holds a reference to the resolved config.BaseLintConfig and tracks the
// target directory and secure root for file operations.
type BaseLint struct {
	cfg       *config.BaseLintConfig // Validated service configuration.
	targetDir string                 // Destination path for rendered files.
	root      *os.Root               // Secure a root handle for a safe file writes.
	script    types.Registrar        // Contains the Methods to run in base_local.
}

// NewBaseLint constructs a new BaseLint instance for a given target directory and provided registrar.
func NewBaseLint(targetDir string) *BaseLint {
	return &BaseLint{
		targetDir: targetDir,
		script:    nil,
	}
}

// SetScriptReceiver sets the Registrar implementation used for registering script commands.
//
// This allows services like base_local to collect command definitions from this service.
//
// It must be called before Run if script registration is desired.
func (b *BaseLint) SetScriptReceiver(reg types.Registrar) {
	b.script = reg
}

// ID returns the static service identifier used to register and retrieve this service.
//
// It matches the constant defined in the type package and must align with
// the corresponding entry in the goboot config (e.g., "base_lint").
func (b *BaseLint) ID() string {
	return types.ServiceNameBaseLint
}

// SetConfig assigns the base lint configuration.
//
// It performs a type assertion to ensure the correct config type was passed.
//
// This assumes config has been validated during initialization.
func (b *BaseLint) SetConfig(cfg config.ServiceConfig) error {
	baseCfg, ok := cfg.(*config.BaseLintConfig)
	if !ok {
		return errors.New("invalid config type for base_lint")
	}

	b.cfg = baseCfg

	// Ensure source and target paths are different (prevent accidental overwrite).
	err := utils.ComparePaths(b.cfg.SourcePath, b.targetDir, true)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and target: %w", err)
	}

	return nil
}

// Run executes the base lint generation logic.
//
// It performs a type assertion to ensure the correct config type was passed.
//
// After storing the typed config, it begins the file scaffolding process.
//
// This assumes config has been validated during initialization.
func (b *BaseLint) Run() error {
	curRoot, err := utils.CreateRootDir(b.targetDir, b.cfg.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to create root dir: %w", err)
	}

	b.root = curRoot

	// Trigger the core logic to copy and render relevant linter files.
	err = b.copyFiles()
	if err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	return nil
}

// copyFiles iterates over all configured linters in the config, and for each enabled linter,
// it copies and renders the corresponding configuration file.
//
// Skipped:
//   - Linters that are disabled.
//   - Linters without associated files (e.g., "make" in this case).
//   - Unknown linters (to allow future-proofing).
//
// Returns an error if any file fails to copy or render.
//nolint:cyclop // Flat switch is preferred for explicit control and traceability.
func (b *BaseLint) copyFiles() error {
	for name, info := range b.cfg.Linters {
		if !info.Enabled {
			continue
		}

		switch name {
		case types.LinterGo:
			err := b.handleLintFile(types.LinterGo, ".golangci.yml")
			if err != nil {
				return fmt.Errorf("failed to handle %s lint file: %w", types.LinterGo, err)
			}
		case types.LinterYAML:
			err := b.handleLintFile(types.LinterYAML, ".yamllint.yml")
			if err != nil {
				return fmt.Errorf("failed to handle %s lint file: %w", types.LinterYAML, err)
			}
		case types.LinterMake:
			// Skip for now — makefile linter uses flags, not a config file.
			continue
		case types.LinterMD:
			err := b.handleLintFile(types.LinterMD, ".markdownlint.yml")
			if err != nil {
				return fmt.Errorf("failed to handle %s lint file: %w", types.LinterMD, err)
			}
		default:
			// Unknown linter — silently ignored for forward compatibility.
			continue
		}
	}

	if b.script != nil {
		err := b.registerScripts()
		if err != nil {
			return fmt.Errorf("failed to register scripts: %w", err)
		}
	}

	return nil
}

// handleLintFile is a helper that encapsulates the steps to:
//   - Copy a static linter config file from sourcePath to the project.
//   - Apply template rendering with project-specific values.
//
// Parameters:
//   - name: Linter name (used for log context).
//   - fileName: File to copy and render.
func (b *BaseLint) handleLintFile(name, fileName string) error {
	err := b.copyFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to copy %s: %w", name, err)
	}

	err = utils.RenderTemplateToFile("lint_file", b.root, fileName, b.cfg)
	if err != nil {
		return fmt.Errorf("failed to render template to file: %w", err)
	}

	return nil
}

// copyFile reads a single static file from the SourcePath and writes it
// into the target directory within the secure os.Root.
//
// Used for initial transfer before template rendering is applied.
//
// Expect a relative filename (e.g., ".golangci.yml").
//
// Returns an error if reading or writing fails.
func (b *BaseLint) copyFile(fileName string) error {
	// #nosec G304 -- this file path is safe and user-defined; used intentionally for scaffolding.
	content, err := os.ReadFile(path.Join(b.cfg.SourcePath, fileName))
	if err != nil {
		return fmt.Errorf("failed to read template file %q: %w", fileName, err)
	}

	// Create and write a file into the secured target root.
	dstFile, err := b.root.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file %q in root: %w", fileName, err)
	}
	defer utils.CloseFileWithErr(dstFile)

	_, err = dstFile.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write file %q: %w", fileName, err)
	}

	return nil
}

// registerScripts collects all enabled linter commands and registers them
// with the attached script registrar (typically base_local).
//
// This includes both line-based script registration (e.g., for Makefile/Taskfile)
// and file-based registration (e.g., scripts/lint.sh if applicable).
//
// Returns an error if script registration fails at any point.
func (b *BaseLint) registerScripts() error {
	cmds := make([]string, 0, len(b.cfg.Linters))

	for _, entry := range b.cfg.Linters {
		cmds = append(cmds, entry.Cmd)
	}

	err := b.script.RegisterLines(types.ServiceNameBaseLint, cmds)
	if err != nil {
		return fmt.Errorf("failed to register script commands: %w", err)
	}

	err = b.script.RegisterFile(types.ScriptFileLint, cmds)
	if err != nil {
		return fmt.Errorf("failed to register script file: %w", err)
	}

	return nil
}
