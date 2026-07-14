/* Package basegovernance implements deterministic repository governance generation. */
package basegovernance

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/it-timo/goboot/pkg/gobootutils"
	"github.com/rs/zerolog"
)

const (
	githubTemplateDir = "github"
	gitlabTemplateDir = "gitlab"
)

// BaseGovernance renders common and provider-specific governance files.
type BaseGovernance struct {
	cfg       *config.BaseGovernanceConfig
	targetDir string
	log       zerolog.Logger
}

// NewBaseGovernance constructs a governance service for targetDir.
func NewBaseGovernance(targetDir string) *BaseGovernance {
	return &BaseGovernance{targetDir: targetDir, log: zerolog.Nop()}
}

// SetLogger injects the service logger.
func (b *BaseGovernance) SetLogger(logger zerolog.Logger) {
	b.log = logger
}

// ID returns the stable service identifier.
func (b *BaseGovernance) ID() string {
	return goboottypes.ServiceNameBaseGovernance
}

// SetConfig assigns configuration and enforces source/target isolation.
func (b *BaseGovernance) SetConfig(cfg config.ServiceConfig) error {
	governanceConfig, ok := cfg.(*config.BaseGovernanceConfig)
	if !ok {
		return errors.New("invalid config type for base_governance")
	}

	b.cfg = governanceConfig

	err := gobootutils.ComparePaths(governanceConfig.SourcePath, b.targetDir, true)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and target: %w", err)
	}

	err = gobootutils.ComparePaths(
		governanceConfig.SourcePath,
		filepath.Join(b.targetDir, governanceConfig.ProjectName),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and project root: %w", err)
	}

	err = gobootutils.EnforceTemplateSourceLimits(
		governanceConfig.SourcePath,
		goboottypes.MaxTemplateSourceFiles,
		goboottypes.MaxTemplateSourceBytes,
	)
	if err != nil {
		return fmt.Errorf("failed template source guardrails for base_governance: %w", err)
	}

	return nil
}

// Run writes governance files into the generated project.
func (b *BaseGovernance) Run() error {
	root, err := gobootutils.CreateRootDir(b.targetDir, b.cfg.ProjectName)
	if err != nil {
		return fmt.Errorf("failed to create root dir: %w", err)
	}

	defer func() {
		if closeErr := root.Close(); closeErr != nil {
			b.log.Error().Err(closeErr).Msg("failed to close root dir")
		}
	}()

	for _, spec := range b.fileSpecs() {
		err = b.renderFile(root, spec.source, spec.target)
		if err != nil {
			return err
		}
	}

	b.log.Info().Str("provider", b.cfg.GitProvider).Str("profile", b.cfg.Profile).
		Msg("base_governance service completed")

	return nil
}

type fileSpec struct {
	source string
	target string
}

func (b *BaseGovernance) fileSpecs() []fileSpec {
	files := []fileSpec{
		{source: "CODEOWNERS", target: "CODEOWNERS"},
		{source: "CONTRIBUTING.md", target: "CONTRIBUTING.md"},
		{source: "SECURITY.md", target: "SECURITY.md"},
	}

	if b.cfg.GitProvider == goboottypes.GitProviderGitHub {
		return append(files, b.githubFiles()...)
	}

	return append(files, b.gitlabFiles()...)
}

func (b *BaseGovernance) githubFiles() []fileSpec {
	files := []fileSpec{
		{source: filepath.Join(githubTemplateDir, "PULL_REQUEST_TEMPLATE.md"), target: ".github/PULL_REQUEST_TEMPLATE.md"},
		{source: filepath.Join(githubTemplateDir, "bug_report.yml"), target: ".github/ISSUE_TEMPLATE/bug_report.yml"},
	}

	if b.cfg.Profile != goboottypes.ProfileMinimal {
		files = append(files, fileSpec{
			source: filepath.Join(githubTemplateDir, "feature_request.yml"),
			target: ".github/ISSUE_TEMPLATE/feature_request.yml",
		})
	}

	if b.cfg.Profile == goboottypes.ProfileOSS {
		files = append(files,
			fileSpec{
				source: filepath.Join(githubTemplateDir, "documentation.yml"),
				target: ".github/ISSUE_TEMPLATE/documentation.yml",
			},
			fileSpec{source: filepath.Join(githubTemplateDir, "config.yml"), target: ".github/ISSUE_TEMPLATE/config.yml"},
		)
	}

	if b.cfg.Profile == goboottypes.ProfileEnterprise {
		files = append(files, fileSpec{
			source: filepath.Join(githubTemplateDir, "change_request.yml"),
			target: ".github/ISSUE_TEMPLATE/change_request.yml",
		})
	}

	return files
}

func (b *BaseGovernance) gitlabFiles() []fileSpec {
	files := []fileSpec{
		{source: filepath.Join(gitlabTemplateDir, "Default.md"), target: ".gitlab/merge_request_templates/Default.md"},
		{source: filepath.Join(gitlabTemplateDir, "Bug.md"), target: ".gitlab/issue_templates/Bug.md"},
	}

	if b.cfg.Profile != goboottypes.ProfileMinimal {
		files = append(files, fileSpec{
			source: filepath.Join(gitlabTemplateDir, "Feature.md"),
			target: ".gitlab/issue_templates/Feature.md",
		})
	}

	if b.cfg.Profile == goboottypes.ProfileOSS {
		files = append(files, fileSpec{
			source: filepath.Join(gitlabTemplateDir, "Documentation.md"),
			target: ".gitlab/issue_templates/Documentation.md",
		})
	}

	if b.cfg.Profile == goboottypes.ProfileEnterprise {
		files = append(files, fileSpec{
			source: filepath.Join(gitlabTemplateDir, "Change.md"),
			target: ".gitlab/issue_templates/Change.md",
		})
	}

	return files
}

func (b *BaseGovernance) renderFile(root *os.Root, source string, target string) error {
	sourcePath := filepath.Join(b.cfg.SourcePath, source+goboottypes.TemplateSuffix)

	// #nosec G304 -- governance template paths come from validated scaffold configuration.
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read template %q: %w", sourcePath, err)
	}

	err = gobootutils.EnsureDir(filepath.Dir(target), root, goboottypes.DirPerm)
	if err != nil {
		return fmt.Errorf("failed to ensure governance directory: %w", err)
	}

	err = gobootutils.WriteRootFile(root, target, content, goboottypes.FilePerm)
	if err != nil {
		return fmt.Errorf("failed to write governance file %q: %w", target, err)
	}

	err = gobootutils.RenderTemplateToFile("governance_file", root, target, b.cfg)
	if err != nil {
		return fmt.Errorf("failed to render governance file %q: %w", target, err)
	}

	return nil
}
