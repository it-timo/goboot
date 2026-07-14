/* Package basesupplychain implements deterministic supply-chain security generation. */
package basesupplychain

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

const policyFile = "SUPPLY_CHAIN.md"

// BaseSupplyChain renders security policy documentation and provider automation.
type BaseSupplyChain struct {
	cfg       *config.BaseSupplyChainConfig
	targetDir string
	ci        goboottypes.Registrar
	log       zerolog.Logger
}

// NewBaseSupplyChain constructs a supply-chain service for targetDir.
func NewBaseSupplyChain(targetDir string) *BaseSupplyChain {
	return &BaseSupplyChain{targetDir: targetDir, ci: nil, log: zerolog.Nop()}
}

// SetCIReceiver injects the registrar used for security workflow registration.
func (b *BaseSupplyChain) SetCIReceiver(registrar goboottypes.Registrar) {
	b.ci = registrar
}

// SetLogger injects the service logger.
func (b *BaseSupplyChain) SetLogger(logger zerolog.Logger) {
	b.log = logger
}

// ID returns the stable service identifier.
func (b *BaseSupplyChain) ID() string {
	return goboottypes.ServiceNameBaseSupplyChain
}

// SetConfig assigns configuration and enforces source/target isolation.
func (b *BaseSupplyChain) SetConfig(cfg config.ServiceConfig) error {
	supplyChainConfig, ok := cfg.(*config.BaseSupplyChainConfig)
	if !ok {
		return errors.New("invalid config type for base_supplychain")
	}

	b.cfg = supplyChainConfig

	if err := gobootutils.ComparePaths(supplyChainConfig.SourcePath, b.targetDir, true); err != nil {
		return fmt.Errorf("failed path comparison of src and target: %w", err)
	}

	projectRoot := filepath.Join(b.targetDir, supplyChainConfig.ProjectName)
	if err := gobootutils.ComparePaths(supplyChainConfig.SourcePath, projectRoot, true); err != nil {
		return fmt.Errorf("failed path comparison of src and project root: %w", err)
	}

	if err := gobootutils.EnforceTemplateSourceLimits(
		supplyChainConfig.SourcePath,
		goboottypes.MaxTemplateSourceFiles,
		goboottypes.MaxTemplateSourceBytes,
	); err != nil {
		return fmt.Errorf("failed template source guardrails for base_supplychain: %w", err)
	}

	return nil
}

// Run writes supply-chain documentation and registers security CI jobs.
func (b *BaseSupplyChain) Run() error {
	root, err := gobootutils.CreateRootDir(b.targetDir, b.cfg.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to create root dir: %w", err)
	}

	defer func() {
		if closeErr := root.Close(); closeErr != nil {
			b.log.Error().Err(closeErr).Msg("failed to close root dir")
		}
	}()

	if err = b.renderPolicy(root); err != nil {
		return err
	}

	if err = b.registerCI(); err != nil {
		return err
	}

	b.log.Info().Str("provider", b.cfg.GitProvider).Msg("base_supplychain service completed")

	return nil
}

func (b *BaseSupplyChain) renderPolicy(root *os.Root) error {
	sourcePath := filepath.Join(b.cfg.SourcePath, policyFile+goboottypes.TemplateSuffix)

	// #nosec G304 -- the template path comes from validated scaffold configuration.
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read template %q: %w", sourcePath, err)
	}

	if err = gobootutils.WriteRootFile(root, policyFile, content, goboottypes.FilePerm); err != nil {
		return fmt.Errorf("failed to write supply-chain policy: %w", err)
	}

	if err = gobootutils.RenderTemplateToFile("supply_chain_policy", root, policyFile, b.cfg); err != nil {
		return fmt.Errorf("failed to render supply-chain policy: %w", err)
	}

	return nil
}

func (b *BaseSupplyChain) registerCI() error {
	if b.ci == nil {
		return nil
	}

	commands := []string{
		"go run golang.org/x/vuln/cmd/govulncheck@" + b.cfg.GovulncheckVersion + " ./...",
		"go run github.com/google/go-licenses/v2@" + b.cfg.GoLicensesVersion +
			" check ./... --allowed_licenses=" + strings.Join(b.cfg.AllowedLicenses, ","),
	}

	if err := b.ci.RegisterLines(goboottypes.ServiceNameBaseSupplyChain, commands); err != nil {
		return fmt.Errorf("failed to register supply-chain commands: %w", err)
	}

	if err := b.ci.RegisterFile(goboottypes.CIFileSecurity, commands); err != nil {
		return fmt.Errorf("failed to register security CI file: %w", err)
	}

	return nil
}
