/*
Package baseproject implements the base_project service.
*/
package baseproject

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/it-timo/goboot/pkg/gobootutils"
)

// BaseProject generates the base project structure from templates.
type BaseProject struct {
	cfg       *config.BaseProjectConfig
	targetDir string
	root      *os.Root
}

// NewBaseProject returns a new BaseProject with an associated target path.
func NewBaseProject(targetDir string) *BaseProject {
	return &BaseProject{
		targetDir: targetDir,
	}
}

// ID returns the service identifier used by config and orchestration.
func (b *BaseProject) ID() string {
	return goboottypes.ServiceNameBaseProject
}

// SetConfig assigns validated base project config and blocks source==target runs.
func (b *BaseProject) SetConfig(cfg config.ServiceConfig) error {
	baseCfg, ok := cfg.(*config.BaseProjectConfig)
	if !ok {
		return errors.New("invalid config type for base_project")
	}

	b.cfg = baseCfg

	// Ensure source and target paths are different (prevent accidental overwrite).
	err := gobootutils.ComparePaths(b.cfg.SourcePath, b.targetDir, true)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and target: %w", err)
	}

	err = gobootutils.EnforceTemplateSourceLimits(
		b.cfg.SourcePath,
		goboottypes.MaxTemplateSourceFiles,
		goboottypes.MaxTemplateSourceBytes,
	)
	if err != nil {
		return fmt.Errorf("failed template source guardrails for base_project: %w", err)
	}

	return nil
}

// Run creates the root and generates project paths and file contents.
func (b *BaseProject) Run() error {
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

	err = b.createNewProject()
	if err != nil {
		return fmt.Errorf("failed to create new project: %w", err)
	}

	return nil
}

// createNewProject renders paths first, then file contents, inside b.root.
func (b *BaseProject) createNewProject() error {
	// Step 1: Copy from host template into root.
	err := b.walkAndApply(os.DirFS(b.cfg.SourcePath), b.renderPath)
	if err != nil {
		return fmt.Errorf("failed to render path: %w", err)
	}

	// Step 2: Render file contents in the copied structure.
	err = b.walkAndApply(b.root.FS(), b.renderContent)
	if err != nil {
		return fmt.Errorf("failed to render content: %w", err)
	}

	return nil
}

// walkAndApply walks fsys and runs handler for every entry.
func (b *BaseProject) walkAndApply(fsys fs.FS, handler func(path string, d fs.DirEntry) error) error {
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		err = handler(path, d)
		if err != nil {
			return fmt.Errorf("failed to run function at %q: %w", path, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk dir: %w", err)
	}

	return nil
}

// renderPath renders a template-relative path and copies file bytes into b.root.
// File content templating is handled later by renderContent.
func (b *BaseProject) renderPath(relTemplatePath string, dirEntry fs.DirEntry) error {
	// Render the target path using template logic (e.g. "cmd/{{project_name}}/main.go").
	renderedPath, err := gobootutils.ExecuteTemplateText("relpath", relTemplatePath, b.cfg)
	if err != nil {
		return fmt.Errorf("failed to render path %q: %w", relTemplatePath, err)
	}

	// If it's a directory, create it inside the root.
	if dirEntry.IsDir() {
		err = gobootutils.EnsureDir(renderedPath, b.root, goboottypes.DirPerm)
		if err != nil {
			return fmt.Errorf("failed to ensure directory %q: %w", renderedPath, err)
		}

		return nil
	}

	// TemplateSuffix is a filename-only convention.
	// goboot renders all files equally; the suffix is stripped only for output naming.
	renderedPath = strings.TrimSuffix(renderedPath, goboottypes.TemplateSuffix)

	// Read a template file from the rootDir.
	fullTemplatePath := filepath.Join(b.cfg.SourcePath, relTemplatePath)

	// #nosec G304 -- the path is user-defined and expected to be dynamic.
	content, err := os.ReadFile(fullTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %q: %w", fullTemplatePath, err)
	}

	// Ensure destination directory exists.
	err = gobootutils.EnsureDir(filepath.Dir(renderedPath), b.root, goboottypes.DirPerm)
	if err != nil {
		return fmt.Errorf("failed to ensure destination directory %q: %w", filepath.Dir(renderedPath), err)
	}

	// Create and write a file into root.
	dstFile, err := b.root.Create(renderedPath)
	if err != nil {
		return fmt.Errorf("failed to create file %q in root: %w", renderedPath, err)
	}
	defer gobootutils.CloseFileWithErr(dstFile)

	_, err = dstFile.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write file %q: %w", renderedPath, err)
	}

	return nil
}

// renderContent renders non-directory file content in b.root with BaseProject config.
func (b *BaseProject) renderContent(path string, dirEntry fs.DirEntry) error {
	if dirEntry.IsDir() {
		return nil
	}

	err := gobootutils.RenderTemplateToFile("project_file", b.root, path, b.cfg)
	if err != nil {
		return fmt.Errorf("failed to render template to file: %w", err)
	}

	return nil
}
