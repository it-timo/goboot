/* Package baserelease implements deterministic release automation generation. */
package baserelease

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/it-timo/goboot/pkg/gobootutils"
	"github.com/rs/zerolog"
)

// BaseRelease renders GoReleaser configuration and provider automation.
type BaseRelease struct {
	cfg       *config.BaseReleaseConfig
	targetDir string
	ci        goboottypes.Registrar
	log       zerolog.Logger
}

// NewBaseRelease constructs a release service for targetDir.
func NewBaseRelease(targetDir string) *BaseRelease {
	return &BaseRelease{targetDir: targetDir, ci: nil, log: zerolog.Nop()}
}

// SetCIReceiver injects the registrar used for release workflow registration.
func (b *BaseRelease) SetCIReceiver(registrar goboottypes.Registrar) {
	b.ci = registrar
}

// SetLogger injects the service logger.
func (b *BaseRelease) SetLogger(logger zerolog.Logger) {
	b.log = logger
}

// ID returns the stable service identifier.
func (b *BaseRelease) ID() string {
	return goboottypes.ServiceNameBaseRelease
}

// SetConfig assigns configuration and enforces source/target isolation.
func (b *BaseRelease) SetConfig(cfg config.ServiceConfig) error {
	releaseConfig, ok := cfg.(*config.BaseReleaseConfig)
	if !ok {
		return errors.New("invalid config type for base_release")
	}

	b.cfg = releaseConfig

	err := gobootutils.ComparePaths(releaseConfig.SourcePath, b.targetDir, true)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and target: %w", err)
	}

	err = gobootutils.ComparePaths(
		releaseConfig.SourcePath,
		filepath.Join(b.targetDir, releaseConfig.ProjectName),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and project root: %w", err)
	}

	err = gobootutils.EnforceTemplateSourceLimits(
		releaseConfig.SourcePath,
		goboottypes.MaxTemplateSourceFiles,
		goboottypes.MaxTemplateSourceBytes,
	)
	if err != nil {
		return fmt.Errorf("failed template source guardrails for base_release: %w", err)
	}

	return nil
}

// Run writes release files into the generated project.
func (b *BaseRelease) Run() error {
	root, err := gobootutils.CreateRootDir(b.targetDir, b.cfg.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to create root dir: %w", err)
	}

	defer func() {
		if closeErr := root.Close(); closeErr != nil {
			b.log.Error().Err(closeErr).Msg("failed to close root dir")
		}
	}()

	err = b.renderFiles(root)
	if err != nil {
		return err
	}

	err = b.registerCI()
	if err != nil {
		return err
	}

	b.log.Info().Str("provider", strings.ToLower(b.cfg.GitProvider)).Msg("base_release service completed")

	return nil
}

func (b *BaseRelease) renderFiles(root *os.Root) error {
	files := []string{".goreleaser.yml", "RELEASE.md"}

	for _, file := range files {
		err := b.renderFile(root, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BaseRelease) renderFile(root *os.Root, file string) error {
	src := filepath.Join(b.cfg.SourcePath, file+goboottypes.TemplateSuffix)

	// #nosec G304 -- release template paths come from validated, user-defined scaffold configuration.
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read template %q: %w", src, err)
	}

	err = gobootutils.EnsureDir(filepath.Dir(file), root, goboottypes.DirPerm)
	if err != nil {
		return fmt.Errorf("failed to ensure release directory: %w", err)
	}

	err = gobootutils.WriteRootFile(root, file, content, goboottypes.FilePerm)
	if err != nil {
		return fmt.Errorf("failed to write release file %q: %w", file, err)
	}

	err = gobootutils.RenderTemplateToFile("release_file", root, file, b.cfg)
	if err != nil {
		return fmt.Errorf("failed to render release file %q: %w", file, err)
	}

	return nil
}

func (b *BaseRelease) registerCI() error {
	if b.ci == nil {
		return nil
	}

	commands := []string{"goreleaser release --clean"}

	err := b.ci.RegisterLines(goboottypes.ServiceNameBaseRelease, commands)
	if err != nil {
		return fmt.Errorf("failed to register release commands: %w", err)
	}

	err = b.ci.RegisterFile(goboottypes.CIFileRelease, commands)
	if err != nil {
		return fmt.Errorf("failed to register release CI file: %w", err)
	}

	return nil
}
