/*
Package baseProject implements the core bootstrapping logic for the "base_project" service.

This service is responsible for handling the foundational project generation tasks such as
- Setting up directory structures
- Applying project metadata
- Injecting base templates and naming conventions

It is one of the default services and required modules in the goboot system
and operates on the validated configuration of type config.BaseProjectConfig.
*/
package baseProject

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/types"
	"github.com/it-timo/goboot/pkg/utils"
)

const (
	// rootDir contains the template path for the project_path dirs and files.
	rootDir = "templates/project_base"

	// dirPerm is the directory permission level.
	dirPerm = 0755
)

// BaseProject implements the Service interface and represents the logic layer for generating the base project structure.
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

// Run executes the base project generation logic.
//
// It performs a type assertion to ensure the correct config type was passed.
//
// After storing the typed config, it begins the directory and file scaffolding process.
//
// This assumes config has been validated during initialization.
func (b *BaseProject) Run(cfg config.ServiceConfig) error {
	baseCfg, ok := cfg.(*config.BaseProjectConfig)
	if !ok {
		return fmt.Errorf("invalid config type for base_project")
	}

	b.cfg = baseCfg

	err := b.createRootDir()
	if err != nil {
		return fmt.Errorf("failed to create root dir: %w", err)
	}

	err = b.createNewProject()
	if err != nil {
		return fmt.Errorf("failed to create new project: %w", err)
	}

	return nil
}

// createRootDir creates the root directory for the new project.
//
// It uses the project name to create a new directory under the target directory.
//
// The directory is created with the correct permissions and ownership.
//
// It also opens the directory as a Root for further processing.
//
// This method is called before the template walker begins.
func (b *BaseProject) createRootDir() error {
	curPath := filepath.Join(b.targetDir, b.cfg.ProjectName)

	cleanPath, err := filepath.Abs(path.Clean(curPath))
	if err != nil {
		return fmt.Errorf("failed to clean path: %w", err)
	}

	err = os.MkdirAll(cleanPath, dirPerm)
	if err != nil {
		return fmt.Errorf("failed to create root dir: %w", err)
	}

	curRoot, err := os.OpenRoot(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to open root: %w", err)
	}

	b.root = curRoot

	return nil
}

// createNewProject initializes the project structure by rendering paths and file contents.
//
// It performs two passes over the template directory inside the secure root:
// 1. It renders file and directory paths using Go templates.
// 2. It renders the contents of template files using the BaseProject config.
//
// All operations are strictly contained within the `*os.Root` directory.
func (b *BaseProject) createNewProject() error {
	// Step 1: Copy from host template into root.
	err := b.walkAndApply(os.DirFS(rootDir), b.renderPath)
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
//   - fsys: the filesystem to walk (e.g., os.DirFS(rootDir), b.root.FS())
//   - fn: the handler to apply to each entry.
//
// Returns:
//   - An error if walking or handling fails.
func (b *BaseProject) walkAndApply(fsys fs.FS, fn func(path string, d fs.DirEntry) error) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		return fn(path, d)
	})
}

// renderPath processes directories and files from the template source,
// applying Go template rendering to all relative paths and replicating the structure inside the root.
//
// For files, the content is copied as-is (templating happens in renderContent).
//
// Parameters:
//   - path: The relative path within the root.
//   - d: The directory entry metadata.
//
// Returns:
//   - An error if path rendering, reading, or writing fails.
func (b *BaseProject) renderPath(relTemplatePath string, d fs.DirEntry) error {
	// Render the relative path using template logic (e.g. "cmd/{{project_name}}/main.go").
	renderedPath, err := b.executeTemplate("relpath", relTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to render path %q: %w", relTemplatePath, err)
	}

	// If it's a directory, create it inside the root.
	if d.IsDir() {
		return utils.EnsureDir(renderedPath, b.root, dirPerm)
	}

	// Read a template file from the rootDir.
	fullTemplatePath := filepath.Join(rootDir, relTemplatePath)

	content, err := os.ReadFile(fullTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %q: %w", fullTemplatePath, err)
	}

	// Ensure destination directory exists.
	err = utils.EnsureDir(filepath.Dir(renderedPath), b.root, dirPerm)
	if err != nil {
		return err
	}

	// Create and write a file into root.
	dstFile, err := b.root.Create(renderedPath)
	if err != nil {
		return fmt.Errorf("failed to create file %q in root: %w", renderedPath, err)
	}
	defer dstFile.Close()

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
//   - d: The directory entry metadata.
//
// Returns:
//   - An error if reading, rendering, or writing fails.
func (b *BaseProject) renderContent(path string, d fs.DirEntry) error {
	if d.IsDir() {
		return nil
	}

	// Read file content.
	file, err := b.root.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer file.Close()

	raw, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", path, err)
	}

	// Apply template rendering.
	rendered, err := b.executeTemplate("file", string(raw))
	if err != nil {
		return fmt.Errorf("failed to execute template on %q: %w", path, err)
	}

	// Save the rendered content back inside the os.Root-protected directory.
	outFile, err := b.root.Create(path)
	if err != nil {
		return fmt.Errorf("failed to open file for overwrite %q: %w", path, err)
	}
	defer outFile.Close()

	_, err = outFile.Write([]byte(rendered))
	if err != nil {
		return fmt.Errorf("failed to write rendered content to %q: %w", path, err)
	}

	return nil
}

// executeTemplate parses and renders a single Go template string using the current BaseProject configuration.
//
// The provided `name` is used as the identifier for the template (for debugging purposes),
// and `text` is the raw template content to be rendered.
//
// Returns the rendered string output or an error if parsing or execution fails.
func (b *BaseProject) executeTemplate(name, text string) (string, error) {
	tmpl, err := template.New(name).Parse(text)
	if err != nil {
		return "", fmt.Errorf("template parse failed: %w", err)
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, b.cfg)
	if err != nil {
		return "", fmt.Errorf("template execute failed: %w", err)
	}

	return buf.String(), nil
}
