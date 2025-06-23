/*
Package baseproject implements the core bootstrapping logic for the "base_project" service.

This service is responsible for handling the foundational project generation tasks such as
- Setting up directory structures.
- Applying project metadata.
- Injecting base templates and naming conventions.

It is one of the default services and required modules in the goboot system
and operates on the validated configuration of type config.BaseProjectConfig.
*/
package baseproject

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/types"
	"github.com/it-timo/goboot/pkg/utils"
)

// BaseProject implements the Service interface and represents the logic layer
// for generating the base project structure.
//
// It holds a reference to the resolved config.BaseProjectConfig after validation.
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

// ID returns the static service identifier used to register and retrieve this service.
//
// It matches the constant defined in the type package and must align with
// the corresponding entry in the goboot config (e.g., "base_project").
func (b *BaseProject) ID() string {
	return types.ServiceNameBaseProject
}

// SetConfig assigns the base project configuration.
//
// It performs a type assertion to ensure the correct config type was passed.
//
// This assumes config has been validated during initialization.
func (b *BaseProject) SetConfig(cfg config.ServiceConfig) error {
	baseCfg, ok := cfg.(*config.BaseProjectConfig)
	if !ok {
		return errors.New("invalid config type for base_project")
	}

	b.cfg = baseCfg

	// Ensure source and target paths are different (prevent accidental overwrite).
	err := utils.ComparePaths(b.cfg.SourcePath, b.targetDir, true)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and target: %w", err)
	}

	return nil
}

// Run executes the base project generation logic.
//
// It performs a type assertion to ensure the correct config type was passed.
//
// After storing the typed config, it begins the directory and file scaffolding process.
//
// This assumes config has been validated during initialization.
func (b *BaseProject) Run() error {
	curRoot, err := utils.CreateRootDir(b.targetDir, b.cfg.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to create root dir: %w", err)
	}

	b.root = curRoot

	err = b.createNewProject()
	if err != nil {
		return fmt.Errorf("failed to create new project: %w", err)
	}

	return nil
}

// createNewProject initializes the project structure by rendering paths and file contents.
//
// It performs two passes over the template directory inside the secure root:
//   - 1. Renders file and directory paths using Go templates.
//   - 2. Renders the contents of template files using the BaseProject config.
//
// All operations are strictly contained within the `*os.Root` directory.
func (b *BaseProject) createNewProject() error {
	// Step 1: Copy from host template into root.
	err := b.walkAndApply(os.DirFS(b.cfg.SourcePath), b.renderPath)
	if err != nil {
		return fmt.Errorf("failed to render path: %w", err)
	}

	// 2. Render file contents inside an already-copied structure.
	err = b.walkAndApply(b.root.FS(), b.renderContent)
	if err != nil {
		return fmt.Errorf("failed to render content: %w", err)
	}

	return nil
}

// walkAndApply traverses the given fs.FS starting from the root ".", applying the handler function to each entry.
//
// Parameters:
//   - fsys: the filesystem to walk (e.g., os.DirFS(rootDir), b.root.FS()).
//   - handler: the function to apply to each entry.
//
// Returns an error if walking or handling fails.
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

// renderPath processes directories and files from the template source,
// applying Go template rendering to all relative paths and replicating the structure inside the root.
//
// For files, the content is copied as-is (templating happens in renderContent).
//
// Parameters:
//   - path: The relative target path within the root.
//   - dirEntry: The directory entry metadata.
//
// Returns an error if path rendering, reading, or writing fails.
func (b *BaseProject) renderPath(relTemplatePath string, dirEntry fs.DirEntry) error {
	// Render the target path using template logic (e.g. "cmd/{{project_name}}/main.go").
	renderedPath, err := utils.ExecuteTemplateText("relpath", relTemplatePath, b.cfg)
	if err != nil {
		return fmt.Errorf("failed to render path %q: %w", relTemplatePath, err)
	}

	// If it's a directory, create it inside the root.
	if dirEntry.IsDir() {
		err = utils.EnsureDir(renderedPath, b.root, types.DirPerm)
		if err != nil {
			return fmt.Errorf("failed to ensure directory %q: %w", renderedPath, err)
		}

		return nil
	}

	// Read a template file from the rootDir.
	fullTemplatePath := filepath.Join(b.cfg.SourcePath, relTemplatePath)

	// #nosec G304 -- the path is user-defined and expected to be dynamic.
	content, err := os.ReadFile(fullTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %q: %w", fullTemplatePath, err)
	}

	// Ensure destination directory exists.
	err = utils.EnsureDir(filepath.Dir(renderedPath), b.root, types.DirPerm)
	if err != nil {
		return fmt.Errorf("failed to ensure destination directory %q: %w", filepath.Dir(renderedPath), err)
	}

	// Create and write a file into root.
	dstFile, err := b.root.Create(renderedPath)
	if err != nil {
		return fmt.Errorf("failed to create file %q in root: %w", renderedPath, err)
	}
	defer utils.CloseFileWithErr(dstFile)

	_, err = dstFile.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write file %q: %w", renderedPath, err)
	}

	return nil
}

// renderContent applies Go template rendering to the content of all non-directory files
// within the root directory, using the BaseProject configuration as the template context.
//
// The rendered content is written back to the same file location.
//
// Parameters:
//   - path: The relative path to the file within root.
//   - dirEntry: The directory entry metadata.
//
// Returns an error if reading, rendering, or writing fails.
func (b *BaseProject) renderContent(path string, dirEntry fs.DirEntry) error {
	if dirEntry.IsDir() {
		return nil
	}

	err := utils.RenderTemplateToFile("project_file", b.root, path, b.cfg)
	if err != nil {
		return fmt.Errorf("failed to render template to file: %w", err)
	}

	return nil
}
