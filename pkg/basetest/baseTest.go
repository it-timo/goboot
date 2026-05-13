/*
Package basetest implements the base_test generation service.
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
	"github.com/rs/zerolog"
)

// BaseTest renders test scaffolding templates and registers optional test commands.
type BaseTest struct {
	cfg       *config.BaseTestConfig
	targetDir string
	root      *os.Root
	script    goboottypes.Registrar
	ci        goboottypes.Registrar
	log       zerolog.Logger
}

// NewBaseTest constructs BaseTest for a target directory.
func NewBaseTest(targetDir string) *BaseTest {
	return &BaseTest{
		targetDir: targetDir,
		script:    nil,
		ci:        nil,
		log:       zerolog.Nop(),
	}
}

// SetLogger injects the service-specific logger.
func (b *BaseTest) SetLogger(logger zerolog.Logger) {
	b.log = logger
}

// SetScriptReceiver injects the registrar used for local script registration.
func (b *BaseTest) SetScriptReceiver(reg goboottypes.Registrar) {
	b.script = reg
}

// SetCIReceiver injects the registrar used for CI registration.
func (b *BaseTest) SetCIReceiver(reg goboottypes.Registrar) {
	b.ci = reg
}

// ID returns the service identifier.
func (b *BaseTest) ID() string {
	return goboottypes.ServiceNameBaseTest
}

// SetConfig assigns validated base_test config and blocks source==target runs.
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

	err = gobootutils.EnforceTemplateSourceLimits(
		b.cfg.SourcePath,
		goboottypes.MaxTemplateSourceFiles,
		goboottypes.MaxTemplateSourceBytes,
	)
	if err != nil {
		return fmt.Errorf("failed template source guardrails for base_test: %w", err)
	}

	return nil
}

// Run creates test files and optionally registers script/CI commands.
func (b *BaseTest) Run() error {
	b.log.Info().Str("project_name", b.cfg.ProjectName).Msg("running base_test service")

	curRoot, err := gobootutils.CreateRootDir(b.targetDir, b.cfg.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to create root dir: %w", err)
	}
	defer func() {
		err := curRoot.Close()
		if err != nil {
			b.log.Error().Err(err).Msg("failed to close root dir")
		}
	}()

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

		b.log.Debug().Msg("registered test command in script registrar")
	}

	if b.ci != nil {
		err := b.registerCIJobs()
		if err != nil {
			return fmt.Errorf("failed to register ci jobs: %w", err)
		}

		b.log.Debug().Msg("registered test command in CI registrar")
	}

	b.log.Info().Msg("base_test service completed")

	return nil
}

// createNewTestSetup renders paths first, then file contents, inside b.root.
func (b *BaseTest) createNewTestSetup() error {
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

// renderPath renders a template-relative path and copies file bytes into b.root.
// File content templating is handled later by renderContent.
func (b *BaseTest) renderPath(relTemplatePath string, dirEntry fs.DirEntry) error {
	// Render the target path using template logic (e.g. "cmd/{{project_name}}/main.go").
	renderedPath, err := gobootutils.ExecuteTemplateText("relpath", relTemplatePath, b.cfg)
	if err != nil {
		return fmt.Errorf("failed to render path %q: %w", relTemplatePath, err)
	}

	if strings.Contains(renderedPath, "suite_test.go") && b.cfg.UseStyle != goboottypes.TestStyleGinkgo {
		// Skip suite files when stdlib style is selected.
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

// renderContent renders non-directory file content in b.root with BaseTest config.
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

// registerCIJobs registers the standard test command with the attached CI registrar.
func (b *BaseTest) registerCIJobs() error {
	err := b.ci.RegisterLines(goboottypes.ServiceNameBaseTest, []string{b.cfg.TestCMD})
	if err != nil {
		return fmt.Errorf("failed to register ci commands: %w", err)
	}

	err = b.ci.RegisterFile(goboottypes.CIFileTest, []string{b.cfg.TestCMD})
	if err != nil {
		return fmt.Errorf("failed to register ci job file: %w", err)
	}

	return nil
}
