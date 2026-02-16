/*
Package baselint implements the base_lint generation service.
*/
package baselint

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/it-timo/goboot/pkg/gobootutils"
)

// BaseLint renders linter config files and registers optional lint commands.
type BaseLint struct {
	cfg       *config.BaseLintConfig
	targetDir string
	root      *os.Root
	script    goboottypes.Registrar
	ci        goboottypes.Registrar
}

// NewBaseLint constructs BaseLint for a target directory.
func NewBaseLint(targetDir string) *BaseLint {
	return &BaseLint{
		targetDir: targetDir,
		script:    nil,
		ci:        nil,
	}
}

// SetScriptReceiver injects the registrar used for local script registration.
func (b *BaseLint) SetScriptReceiver(reg goboottypes.Registrar) {
	b.script = reg
}

// SetCIReceiver injects the registrar used for CI registration.
func (b *BaseLint) SetCIReceiver(reg goboottypes.Registrar) {
	b.ci = reg
}

// ID returns the service identifier.
func (b *BaseLint) ID() string {
	return goboottypes.ServiceNameBaseLint
}

// SetConfig assigns validated base_lint config and blocks source==target runs.
func (b *BaseLint) SetConfig(cfg config.ServiceConfig) error {
	baseCfg, ok := cfg.(*config.BaseLintConfig)
	if !ok {
		return errors.New("invalid config type for base_lint")
	}

	b.cfg = baseCfg

	// Ensure source and target paths are different (prevent accidental overwrite).
	err := gobootutils.ComparePaths(b.cfg.SourcePath, b.targetDir, true)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and target: %w", err)
	}

	return nil
}

// Run opens the target root and generates enabled linting files.
func (b *BaseLint) Run() error {
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

	err = b.copyFiles()
	if err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	return nil
}

// copyFiles copies and renders configuration files for enabled linters.
//nolint:cyclop // Flat switch is preferred for explicit control and traceability.
func (b *BaseLint) copyFiles() error {
	for name, info := range b.cfg.Linters {
		if !info.Enabled {
			continue
		}

		switch name {
		case goboottypes.LinterGo:
			err := b.handleLintFile(goboottypes.LinterGo, ".golangci.yml")
			if err != nil {
				return fmt.Errorf("failed to handle %s lint file: %w", goboottypes.LinterGo, err)
			}
		case goboottypes.LinterYAML:
			err := b.handleLintFile(goboottypes.LinterYAML, ".yamllint.yml")
			if err != nil {
				return fmt.Errorf("failed to handle %s lint file: %w", goboottypes.LinterYAML, err)
			}
		case goboottypes.LinterMake:
			// Skip for now — makefile linter uses flags, not a config file.
			continue
		case goboottypes.LinterMD:
			err := b.handleLintFile(goboottypes.LinterMD, ".markdownlint.yml")
			if err != nil {
				return fmt.Errorf("failed to handle %s lint file: %w", goboottypes.LinterMD, err)
			}
		case goboottypes.LinterShell:
			err := b.handleLintFile(goboottypes.LinterShell, ".shellcheckrc")
			if err != nil {
				return fmt.Errorf("failed to handle %s lint file: %w", goboottypes.LinterShell, err)
			}
		case goboottypes.LinterSHFMT:
			// Skip for now — shfmt linter uses flags, not a config file.
			continue
		case goboottypes.LinterEditor:
			// Skip for now — editor linter uses flags, not a config file.
			continue
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

	if b.ci != nil {
		err := b.registerCIJobs()
		if err != nil {
			return fmt.Errorf("failed to register ci jobs: %w", err)
		}
	}

	return nil
}

// handleLintFile copies a lint template file and renders it with config values.
func (b *BaseLint) handleLintFile(name, fileName string) error {
	err := b.copyFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to copy %s: %w", name, err)
	}

	err = gobootutils.RenderTemplateToFile("lint_file", b.root, fileName, b.cfg)
	if err != nil {
		return fmt.Errorf("failed to render template to file: %w", err)
	}

	return nil
}

// copyFile copies one lint template file into root.
func (b *BaseLint) copyFile(fileName string) error {
	src := path.Join(b.cfg.SourcePath, fileName+goboottypes.TemplateSuffix)

	// #nosec G304 -- path is safe and user-defined; used intentionally for scaffolding.
	content, err := os.ReadFile(src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("missing required template %q (expected %q)", fileName, src)
		}

		return fmt.Errorf("failed to read template file %q: %w", src, err)
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

	return nil
}

// gatherCommands collects all enabled linter commands.
func (b *BaseLint) gatherCommands() []string {
	cmds := make([]string, 0, len(b.cfg.Linters))

	for name, entry := range b.cfg.Linters {
		if !entry.Enabled {
			continue
		}

		if strings.TrimSpace(entry.Cmd) == "" {
			continue
		}

		switch name {
		case goboottypes.LinterGo, goboottypes.LinterYAML, goboottypes.LinterMake,
			goboottypes.LinterMD, goboottypes.LinterShell, goboottypes.LinterSHFMT,
			goboottypes.LinterEditor:
			cmds = append(cmds, entry.Cmd)
		default:
			continue
		}
	}

	return cmds
}

// registerScripts registers enabled lint commands with the local script registrar.
func (b *BaseLint) registerScripts() error {
	cmds := b.gatherCommands()

	if len(cmds) == 0 {
		return nil
	}

	err := b.script.RegisterLines(goboottypes.ServiceNameBaseLint, cmds)
	if err != nil {
		return fmt.Errorf("failed to register script commands: %w", err)
	}

	err = b.script.RegisterFile(goboottypes.ScriptFileLint, cmds)
	if err != nil {
		return fmt.Errorf("failed to register script file: %w", err)
	}

	return nil
}

// registerCIJobs registers all enabled linter commands with the attached CI registrar.
func (b *BaseLint) registerCIJobs() error {
	cmds := b.gatherCommands()

	if len(cmds) == 0 {
		return nil
	}

	err := b.ci.RegisterLines(goboottypes.ServiceNameBaseLint, cmds)
	if err != nil {
		return fmt.Errorf("failed to register ci commands: %w", err)
	}

	err = b.ci.RegisterFile(goboottypes.CIFileLint, cmds)
	if err != nil {
		return fmt.Errorf("failed to register ci job file: %w", err)
	}

	return nil
}
