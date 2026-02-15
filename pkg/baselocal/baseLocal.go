/*
Package baselocal implements the base_local generation service.
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

// BaseLocal renders local tooling files (Makefile/Taskfile/pre-commit/scripts).
type BaseLocal struct {
	cfg       *config.BaseLocalConfig
	targetDir string
	root      *os.Root
	scriptRegistry
}

// scriptRegistry stores command registrations grouped by output type.
type scriptRegistry struct {
	ProjectName   string
	MakeScripts   map[string][]string // service → commands
	TaskScripts   map[string][]string // service → commands
	CommitScripts map[string][]string // service → commands
	ScriptFiles   map[string][]string // fileName → commands
}

// NewBaseLocal constructs BaseLocal for a target directory.
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

// ID returns the service identifier.
func (b *BaseLocal) ID() string {
	return goboottypes.ServiceNameBaseLocal
}

// SetConfig assigns validated base_local config and blocks source==target runs.
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

// Run opens the target root and generates enabled local tooling files.
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

	err = b.copyFiles()
	if err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	return nil
}

// RegisterLines stores commands for enabled output types (make/task/commit).
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

// RegisterFile stores script file commands when script output is enabled.
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

// copyFiles copies and renders all enabled local tooling outputs.
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

// copyFile copies one template file into root and renders it with scriptRegistry.
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
