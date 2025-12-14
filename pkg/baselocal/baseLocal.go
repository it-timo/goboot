/*
Package baselocal implements the core bootstrapping logic for the "base_local" service.

This service is responsible for setting up the foundational local script infrastructure
of a newly scaffolded Go project.

It handles:
- Copying predefined local script templates into the project.
- Applying Go `text/template` rendering to inject project-specific metadata.
- Respecting a strict separation of logic per script (Linter, Test, etc.).

The service expects a validated configuration of type config.BaseLocalConfig
and is one of the default built-in services within the `goboot` system.
*/
package baselocal

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/it-timo/goboot/pkg/gobootutils"
)

// BaseLocal implements the Service interface and encapsulates the execution logic
// for generating local scripts.
//
// It holds a reference to the resolved config.BaseLocalConfig and tracks the
// target directory and secure root for file operations.
type BaseLocal struct {
	cfg       *config.BaseLocalConfig // Validated service configuration.
	targetDir string                  // Destination path for rendered files.
	root      *os.Root                // Secure a root handle for a safe file writes.
	scriptRegistry
}

// scriptRegistry holds collected command-line scripts registered by other services.
//
// These scripts are grouped by output format (Makefile, Taskfile, or script directory).
type scriptRegistry struct {
	ProjectName   string
	MakeScripts   map[string][]string // service → commands
	TaskScripts   map[string][]string // service → commands
	CommitScripts map[string][]string // service → commands
	ScriptFiles   map[string][]string // fileName → commands
}

// NewBaseLocal constructs a new BaseLocal instance for a given target directory.
func NewBaseLocal(targetDir string) *BaseLocal {
	return &BaseLocal{
		targetDir: targetDir,
		scriptRegistry: scriptRegistry{
			MakeScripts:   make(map[string][]string),
			TaskScripts:   make(map[string][]string),
			CommitScripts: make(map[string][]string),
			ScriptFiles:   make(map[string][]string),
		},
	}
}

// ID returns the static service identifier used to register and retrieve this service.
//
// It matches the constant defined in the type package and must align with
// the corresponding entry in the goboot config (e.g., "base_local").
func (b *BaseLocal) ID() string {
	return goboottypes.ServiceNameBaseLocal
}

// SetConfig assigns the base local configuration.
//
// It performs a type assertion to ensure the correct config type was passed.
//
// This assumes config has been validated during initialization.
func (b *BaseLocal) SetConfig(cfg config.ServiceConfig) error {
	baseCfg, ok := cfg.(*config.BaseLocalConfig)
	if !ok {
		return errors.New("invalid config type for base_local")
	}

	b.cfg = baseCfg

	// Ensure source and target paths are different (prevent accidental overwrite).
	err := gobootutils.ComparePaths(b.cfg.SourcePath, b.targetDir, true)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and target: %w", err)
	}

	return nil
}

// Run executes the base local generation logic.
//
// It performs a type assertion to ensure the correct config type was passed.
//
// After storing the typed config, it begins the file scaffolding process.
//
// This assumes config has been validated during initialization.
func (b *BaseLocal) Run() error {
	curRoot, err := gobootutils.CreateRootDir(b.targetDir, b.cfg.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to create root dir: %w", err)
	}

	defer func() {
		err := curRoot.Close()
		if err != nil {
			fmt.Println("Failed to close root dir:", err)
		}
	}()

	b.root = curRoot
	b.ProjectName = b.cfg.ProjectName

	// Trigger the core logic to copy and render relevant script files.
	err = b.copyFiles()
	if err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	return nil
}

// RegisterLines implements goboottypes.Registrar by accepting script lines from another service.
//
// Depending on which script formats are enabled (make/task), the lines are grouped by service name.
//
// It ensures that no service registers the same script type more than once.
func (b *BaseLocal) RegisterLines(name string, lines []string) error {
	for _, entry := range b.cfg.FileList {
		switch entry {
		case goboottypes.ScriptNameMake:
			_, exist := b.MakeScripts[name]
			if exist {
				return fmt.Errorf("service %q already registered in make", name)
			}

			b.MakeScripts[name] = lines
		case goboottypes.ScriptNameTask:
			_, exist := b.TaskScripts[name]
			if exist {
				return fmt.Errorf("service %q already registered in task", name)
			}

			b.TaskScripts[name] = lines
		case goboottypes.ScriptNameCommit:
			_, exist := b.CommitScripts[name]
			if exist {
				return fmt.Errorf("service %q already registered in commit", name)
			}

			b.CommitScripts[name] = lines
		default: // do nothing if no script files are not enabled.
		}
	}

	return nil
}

// RegisterFile implements goboottypes.Registrar by accepting a list of commands to be written into a script file.
//
// Only store the file if the script directory target is enabled and not already used.
func (b *BaseLocal) RegisterFile(name string, lines []string) error {
	for _, entry := range b.cfg.FileList {
		switch entry {
		case goboottypes.ScriptNameScript:
			_, exist := b.ScriptFiles[name]
			if exist {
				return fmt.Errorf("file %q already registered in scripts", name)
			}

			b.ScriptFiles[name] = lines
		default: // Skip unknown or disabled script directory entries.
		}
	}

	return nil
}

// copyFiles performs the actual copy and render operation for all configured script file types.
//
// It handles Makefiles, Taskfiles, pre-commit config, and scripts/ directory as needed,
// based on what the user enabled in the config.
//
// For script files, only files that were previously registered will be copied.
//
//nolint:cyclop // flat logic preferred for clarity and extensibility.
func (b *BaseLocal) copyFiles() error {
	for _, entry := range b.cfg.FileList {
		switch entry {
		case goboottypes.ScriptNameMake:
			err := b.copyFile(b.cfg.SourcePath, "", "Makefile")
			if err != nil {
				return fmt.Errorf("failed to copy Makefile: %w", err)
			}
		case goboottypes.ScriptNameTask:
			err := b.copyFile(b.cfg.SourcePath, "", "Taskfile.yml")
			if err != nil {
				return fmt.Errorf("failed to copy Taskfile: %w", err)
			}
		case goboottypes.ScriptNameCommit:
			err := b.copyFile(b.cfg.SourcePath, "", ".pre-commit-config.yaml")
			if err != nil {
				return fmt.Errorf("failed to copy Pre-Commit: %w", err)
			}
		case goboottypes.ScriptNameScript:
			if len(b.ScriptFiles) > 0 {
				err := gobootutils.EnsureDir(goboottypes.ScriptDirNameScript, b.root, goboottypes.DirPerm)
				if err != nil {
					return fmt.Errorf("failed to create scripts dir: %w", err)
				}

				scriptsSrcPath := path.Join(b.cfg.SourcePath, goboottypes.ScriptDirNameScript)

				for fileName := range b.ScriptFiles {
					err = b.copyFile(scriptsSrcPath, goboottypes.ScriptDirNameScript, fileName)
					if err != nil {
						return fmt.Errorf("failed to copy %q: %w", fileName, err)
					}
				}
			}
		}
	}

	return nil
}

// copyFile reads a single static file from the SourcePath and writes it
// into the target directory within the secure os.Root.
//
// Used for initial transfer before template rendering is applied.
//
// Expect the target path and a relative filename (e.g., "Makefile").
//
// Returns an error if reading or writing fails.
func (b *BaseLocal) copyFile(srcPath, targetPath, fileName string) error {
	src := path.Join(srcPath, fileName+goboottypes.TemplateSuffix)

	// #nosec G304 -- this file path is safe and user-defined; used intentionally for scaffolding.
	content, err := os.ReadFile(src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("missing required template %q (expected %q)", fileName, src)
		}

		return fmt.Errorf("failed to read template file %q: %w", src, err)
	}

	if targetPath == goboottypes.ScriptDirNameScript {
		fileName = path.Join(targetPath, fileName)
	}

	// Create and write a file into the secured target root.
	dstFile, err := b.root.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file %q in root: %w", fileName, err)
	}
	defer gobootutils.CloseFileWithErr(dstFile)

	_, err = dstFile.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write file %q: %w", fileName, err)
	}

	if targetPath == goboottypes.ScriptDirNameScript {
		err = dstFile.Chmod(goboottypes.ScriptPerm)
		if err != nil {
			return fmt.Errorf("failed to set executable permissions on %q: %w", fileName, err)
		}
	}

	err = gobootutils.RenderTemplateToFile("script_file", b.root, fileName, b.scriptRegistry)
	if err != nil {
		return fmt.Errorf("failed to render template to file: %w", err)
	}

	return nil
}
