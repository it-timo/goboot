/*
Package baseci implements CI file generation for the base_ci service.
*/
package baseci

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/it-timo/goboot/pkg/gobootutils"
	"github.com/rs/zerolog"
)

// BaseCI generates CI files from templates plus registered job commands.
type BaseCI struct {
	cfg       *config.BaseCIConfig // Validated service configuration.
	targetDir string               // Destination path for rendered files.
	root      *os.Root             // Secure a root handle for safe file writes.
	log       zerolog.Logger
	ciRegistry
	staticFiles []string
}

// ciRegistry stores CI commands and metadata collected from services and config.
type ciRegistry struct {
	ProjectName      string
	CIDir            string
	GoVersions       []string
	AutoBranches     []string
	ImagePolicy      string
	JobScripts       map[string][]string // service → commands
	FileScripts      map[string][]string // fileName → commands
	FileAllowFailure map[string]bool     // fileName → allow failure
	EnabledJobFiles  []string
}

// NewBaseCI returns a BaseCI bound to targetDir.
func NewBaseCI(targetDir string) *BaseCI {
	return &BaseCI{
		targetDir: targetDir,
		log:       zerolog.Nop(),
		ciRegistry: ciRegistry{
			JobScripts:       make(map[string][]string),
			FileScripts:      make(map[string][]string),
			FileAllowFailure: make(map[string]bool),
		},
	}
}

// SetLogger injects the service-specific logger.
func (b *BaseCI) SetLogger(logger zerolog.Logger) {
	b.log = logger
}

// ID returns the service identifier.
func (b *BaseCI) ID() string {
	return goboottypes.ServiceNameBaseCI
}

// SetConfig assigns config and enforces path and template-source safety checks.
func (b *BaseCI) SetConfig(cfg config.ServiceConfig) error {
	baseCfg, ok := cfg.(*config.BaseCIConfig)
	if !ok {
		return errors.New("invalid config type for base_ci")
	}

	b.cfg = baseCfg

	// Prevent template source and target from pointing to the same location.
	err := gobootutils.ComparePaths(b.cfg.SourcePath, b.targetDir, true)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and target: %w", err)
	}

	err = gobootutils.EnforceTemplateSourceLimits(
		filepath.Join(b.cfg.SourcePath, strings.ToLower(strings.TrimSpace(b.cfg.GitProvider))),
		goboottypes.MaxTemplateSourceFiles,
		goboottypes.MaxTemplateSourceBytes,
	)
	if err != nil {
		return fmt.Errorf("failed template source guardrails for base_ci: %w", err)
	}

	return nil
}

// Run renders provider CI files into the generated project.
func (b *BaseCI) Run() error {
	b.log.Info().
		Str("project_name", b.cfg.ProjectName).
		Str("git_provider", b.cfg.GitProvider).
		Msg("running base_ci service")

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
	b.ProjectName = b.cfg.ProjectName
	b.GoVersions = append([]string(nil), b.cfg.GoVersion...)
	b.AutoBranches = append([]string(nil), b.cfg.AutoBranches...)
	b.ImagePolicy = b.cfg.ImagePolicy

	b.staticFiles, b.CIDir, err = b.resolveLayout()
	if err != nil {
		return err
	}

	err = b.registerConfigJobs()
	if err != nil {
		return fmt.Errorf("failed to register CI jobs: %w", err)
	}

	b.normalizeCommandsForPolicy()

	b.EnabledJobFiles = b.enabledJobFiles()

	err = b.copyFiles()
	if err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	b.log.Info().Int("job_file_count", len(b.EnabledJobFiles)).Msg("base_ci service completed")

	return nil
}

// RegisterLines stores service-level CI commands.
func (b *BaseCI) RegisterLines(name string, lines []string) error {
	if _, exist := b.JobScripts[name]; exist {
		return fmt.Errorf("service %q already registered in ci", name)
	}

	b.JobScripts[name] = lines

	return nil
}

// RegisterFile stores commands for a specific CI job file.
func (b *BaseCI) RegisterFile(name string, lines []string) error {
	if _, exist := b.FileScripts[name]; exist {
		return fmt.Errorf("file %q already registered in ci", name)
	}

	b.FileScripts[name] = lines
	b.FileAllowFailure[name] = false

	return nil
}

func (b *BaseCI) normalizeCommandsForPolicy() {
	if b.ImagePolicy != "strict" {
		return
	}

	for file, cmds := range b.FileScripts {
		normalized := make([]string, 0, len(cmds))

		for _, cmd := range cmds {
			normalized = append(normalized, strictCommandReplacement(cmd))
		}

		b.FileScripts[file] = normalized
	}
}

func strictCommandReplacement(cmd string) string {
	replacer := strings.NewReplacer(
		"golangci/golangci-lint:v2.7.2", "$GOLANGCI_LINT_IMAGE",
		"pipelinecomponents/yamllint:0.35.9", "$YAMLLINT_IMAGE",
		"ghcr.io/igorshubovych/markdownlint-cli:v0.47.0", "$MARKDOWNLINT_IMAGE",
		"cytopia/checkmake:latest-0.5", "$CHECKMAKE_IMAGE",
		"koalaman/shellcheck:v0.11.0", "$SHELLCHECK_IMAGE",
		"cytopia/shellcheck:latest-0.8.0", "$SHELLCHECK_IMAGE",
		"mvdan/shfmt:v3.12.0", "$SHFMT_IMAGE",
		"cytopia/shfmt:latest-1.10", "$SHFMT_IMAGE",
		"mstruebing/editorconfig-checker:v3.6.0", "$EDITORCONFIG_CHECKER_IMAGE",
	)

	return replacer.Replace(cmd)
}

func (b *BaseCI) registerConfigJobs() error {
	if len(b.cfg.Jobs) == 0 {
		return nil
	}

	for name, job := range b.cfg.Jobs {
		if job == nil {
			continue
		}

		err := b.RegisterLines(name, job.Commands)
		if err != nil {
			return err
		}

		err = b.RegisterFile(job.File, job.Commands)
		if err != nil {
			return err
		}

		b.FileAllowFailure[job.File] = job.AllowFailure
	}

	return nil
}

func (b *BaseCI) enabledJobFiles() []string {
	enabled := make([]string, 0, len(b.FileScripts))

	for file := range b.FileScripts {
		enabled = append(enabled, file)
	}

	sort.Strings(enabled)

	return enabled
}

// copyFiles copies and renders static and registered job files.
func (b *BaseCI) copyFiles() error {
	for _, file := range b.staticFiles {
		err := b.copyFile(strings.TrimSpace(file))
		if err != nil {
			return fmt.Errorf("failed to copy %q: %w", file, err)
		}
	}

	for _, file := range b.EnabledJobFiles {
		relPath := b.jobFilePath(file)

		err := b.copyFile(relPath)
		if err != nil {
			return fmt.Errorf("failed to copy %q: %w", relPath, err)
		}
	}

	return nil
}

func (b *BaseCI) copyFile(relPath string) error {
	src := path.Join(
		b.cfg.SourcePath,
		strings.ToLower(strings.TrimSpace(b.cfg.GitProvider)),
		relPath+goboottypes.TemplateSuffix,
	)

	// #nosec G304 -- this file path is safe and user-defined; used intentionally for scaffolding.
	content, err := os.ReadFile(src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("missing required template %q (expected %q)", relPath, src)
		}

		return fmt.Errorf("failed to read template file %q: %w", src, err)
	}

	err = gobootutils.EnsureDir(filepath.Dir(relPath), b.root, goboottypes.DirPerm)
	if err != nil {
		return fmt.Errorf("failed to ensure directory %q: %w", filepath.Dir(relPath), err)
	}

	dstFile, err := b.root.Create(relPath)
	if err != nil {
		return fmt.Errorf("failed to create file %q in root: %w", relPath, err)
	}
	defer gobootutils.CloseFileWithErr(dstFile)

	_, err = dstFile.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write file %q: %w", relPath, err)
	}

	err = gobootutils.RenderTemplateToFile("ci_file", b.root, relPath, b.ciRegistry)
	if err != nil {
		return fmt.Errorf("failed to render template to file: %w", err)
	}

	return nil
}

func (b *BaseCI) jobFilePath(file string) string {
	if b.CIDir == "" {
		return file
	}

	return path.Join(b.CIDir, file)
}

func (b *BaseCI) resolveLayout() ([]string, string, error) {
	switch strings.ToLower(strings.TrimSpace(b.cfg.GitProvider)) {
	case "gitlab":
		return []string{
			".gitlab-ci.yml",
			".gitlab/ci/commands.yml",
			".gitlab/ci/versions.yml",
		}, ".gitlab/ci", nil
	case "github":
		return nil, ".github/workflows", nil
	default:
		return nil, "", fmt.Errorf("unsupported gitProvider: %s", b.cfg.GitProvider)
	}
}
