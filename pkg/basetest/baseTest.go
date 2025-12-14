/*
Package basetest implements the core bootstrapping logic for the "base_test" service.

This service is responsible for setting up the foundational testing infrastructure of a newly scaffolded Go project.

It handles:
  - Copying predefined testing configuration templates into the project.
  - Applying Go `text/template` rendering to inject project-specific metadata into both filenames and content.

The service expects a validated configuration of type config.BaseTestConfig
and is one of the default built-in services within the `goboot` system.
*/
package basetest

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

// BaseTest implements the Service interface and encapsulates the execution logic
// for generating testing configuration files based on user-defined config.
//
// It holds a reference to the resolved config.BaseTestConfig and tracks the
// target directory and secure root for file operations.
type BaseTest struct {
	cfg       *config.BaseTestConfig // Validated service configuration.
	targetDir string                 // Destination path for rendered files.
	root      *os.Root               // Secure a root handle for a safe file writes.
	script    goboottypes.Registrar  // Contains the Methods to run in base_local.
}

// NewBaseTest constructs a new BaseTest instance for a given target directory.
func NewBaseTest(targetDir string) *BaseTest {
	return &BaseTest{
		targetDir: targetDir,
		script:    nil,
	}
}

// SetScriptReceiver sets the Registrar implementation used for registering script commands.
//
// This allows services like base_local to collect command definitions from this service.
//
// It must be called before Run if script registration is desired.
func (b *BaseTest) SetScriptReceiver(reg goboottypes.Registrar) {
	b.script = reg
}

// ID returns the static service identifier used to register and retrieve this service.
func (b *BaseTest) ID() string {
	return goboottypes.ServiceNameBaseTest
}

// SetConfig assigns the base test configuration.
//
// It performs a type assertion to ensure the correct config type was passed.
//
// This assumes config has been validated during initialization.
func (b *BaseTest) SetConfig(cfg config.ServiceConfig) error {
	baseCfg, ok := cfg.(*config.BaseTestConfig)
	if !ok {
		return errors.New("invalid config type for base_test")
	}

	b.cfg = baseCfg

	// Ensure source and target paths are different (prevent accidental overwrite).
	err := gobootutils.ComparePaths(b.cfg.SourcePath, b.targetDir, true)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and target: %w", err)
	}

	return nil
}

// Run executes the base test generation logic.
//
// It recursively walks the configured SourcePath (templates), copies files to the target,
// and applies template rendering to both file paths and file content.
func (b *BaseTest) Run() error {
	curRoot, err := gobootutils.CreateRootDir(b.targetDir, b.cfg.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to create root dir: %w", err)
	}

	b.root = curRoot

	err = b.createNewTestSetup()
	if err != nil {
		return fmt.Errorf("failed to create new test setup: %w", err)
	}

	if b.script != nil {
		err := b.registerScripts()
		if err != nil {
			return fmt.Errorf("failed to register scripts: %w", err)
		}
	}

	return nil
}

// createNewTestSetup initializes the test structure by rendering paths and file contents.
//
// It performs two passes over the template directory inside the secure root:
//   - 1. Renders file and directory paths using Go templates.
//   - 2. Renders the contents of template files using the BaseTest config.
//
// All operations are strictly contained within the `*os.Root` directory.
func (b *BaseTest) createNewTestSetup() error {
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
func (b *BaseTest) walkAndApply(fsys fs.FS, handler func(path string, d fs.DirEntry) error) error {
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
func (b *BaseTest) renderPath(relTemplatePath string, dirEntry fs.DirEntry) error {
	// Render the target path using template logic (e.g. "cmd/{{project_name}}/main.go").
	renderedPath, err := gobootutils.ExecuteTemplateText("relpath", relTemplatePath, b.cfg)
	if err != nil {
		return fmt.Errorf("failed to render path %q: %w", relTemplatePath, err)
	}

	if strings.Contains(renderedPath, "suite_test.go") && b.cfg.UseStyle != goboottypes.TestStyleGinkgo {
		// skip suite files for ginkgo
		return nil
	}

	// Remove the template suffix for the test files.
	renderedPath = strings.TrimSuffix(renderedPath, goboottypes.TemplateSuffix)

	// If it's a directory, create it inside the root.
	if dirEntry.IsDir() {
		err = gobootutils.EnsureDir(renderedPath, b.root, goboottypes.DirPerm)
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

// renderContent applies Go template rendering to the content of all non-directory files
// within the root directory, using the BaseTest configuration as the template context.
//
// The rendered content is written back to the same file location.
//
// Parameters:
//   - path: The relative path to the file within root.
//   - dirEntry: The directory entry metadata.
//
// Returns an error if reading, rendering, or writing fails.
func (b *BaseTest) renderContent(path string, dirEntry fs.DirEntry) error {
	if dirEntry.IsDir() {
		return nil
	}

	err := gobootutils.RenderTemplateToFile("test_file", b.root, path, b.cfg)
	if err != nil {
		return fmt.Errorf("failed to render template to file: %w", err)
	}

	return nil
}

// registerScripts registers the standard test command.
func (b *BaseTest) registerScripts() error {
	err := b.script.RegisterLines(goboottypes.ServiceNameBaseTest, []string{b.cfg.TestCMD})
	if err != nil {
		return fmt.Errorf("failed to register script commands: %w", err)
	}

	err = b.script.RegisterFile(goboottypes.ScriptFileTest, []string{b.cfg.TestCMD})
	if err != nil {
		return fmt.Errorf("failed to register script file: %w", err)
	}

	return nil
}
