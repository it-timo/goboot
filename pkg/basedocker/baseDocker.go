/*
Package basedocker implements the base_docker generation service.
*/
package basedocker

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/it-timo/goboot/pkg/gobootutils"
	"github.com/rs/zerolog"
)

const (
	dockerfileFileEntry = "dockerfile"
	composeFileEntry    = "compose"
)

// BaseDocker renders Dockerfile, compose, and dockerignore templates.
type BaseDocker struct {
	cfg       *config.BaseDockerConfig
	targetDir string
	root      *os.Root
	script    goboottypes.Registrar
	ci        goboottypes.Registrar
	log       zerolog.Logger
}

// NewBaseDocker constructs BaseDocker for a target directory.
func NewBaseDocker(targetDir string) *BaseDocker {
	return &BaseDocker{
		targetDir: targetDir,
		script:    nil,
		ci:        nil,
		log:       zerolog.Nop(),
	}
}

// SetLogger injects the service-specific logger.
func (b *BaseDocker) SetLogger(logger zerolog.Logger) {
	b.log = logger
}

// SetScriptReceiver injects the registrar used for local script registration.
func (b *BaseDocker) SetScriptReceiver(reg goboottypes.Registrar) {
	b.script = reg
}

// SetCIReceiver injects the registrar used for CI registration.
func (b *BaseDocker) SetCIReceiver(reg goboottypes.Registrar) {
	b.ci = reg
}

// ID returns the service identifier.
func (b *BaseDocker) ID() string {
	return goboottypes.ServiceNameBaseDocker
}

// SetConfig assigns validated base_docker config and blocks source==target runs.
func (b *BaseDocker) SetConfig(cfg config.ServiceConfig) error {
	baseCfg, ok := cfg.(*config.BaseDockerConfig)
	if !ok {
		return errors.New("invalid config type for base_docker")
	}

	b.cfg = baseCfg

	err := gobootutils.ComparePaths(b.cfg.SourcePath, b.targetDir, true)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and target: %w", err)
	}

	err = gobootutils.ComparePaths(b.cfg.SourcePath, filepath.Join(b.targetDir, b.cfg.ProjectName), true)
	if err != nil {
		return fmt.Errorf("failed path comparison of src and project root: %w", err)
	}

	err = gobootutils.EnforceTemplateSourceLimits(
		b.cfg.SourcePath,
		goboottypes.MaxTemplateSourceFiles,
		goboottypes.MaxTemplateSourceBytes,
	)
	if err != nil {
		return fmt.Errorf("failed template source guardrails for base_docker: %w", err)
	}

	return nil
}

// Run renders containerization files and registers optional local/CI commands.
func (b *BaseDocker) Run() error {
	b.log.Info().Str("project_name", b.cfg.ProjectName).Msg("running base_docker service")

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

	err = b.copyFiles()
	if err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	if b.script != nil {
		err = b.registerScripts()
		if err != nil {
			return fmt.Errorf("failed to register scripts: %w", err)
		}
	}

	if b.ci != nil {
		err = b.registerCIJobs()
		if err != nil {
			return fmt.Errorf("failed to register ci jobs: %w", err)
		}
	}

	b.log.Info().Msg("base_docker service completed")

	return nil
}

func (b *BaseDocker) copyFiles() error {
	for _, entry := range b.cfg.FileList {
		switch strings.TrimSpace(entry) {
		case dockerfileFileEntry:
			err := b.copyFile("Dockerfile")
			if err != nil {
				return fmt.Errorf("failed to copy Dockerfile: %w", err)
			}
		case composeFileEntry:
			err := b.copyFile("docker-compose.yml")
			if err != nil {
				return fmt.Errorf("failed to copy docker-compose.yml: %w", err)
			}
		case "dockerignore":
			err := b.copyFile(".dockerignore")
			if err != nil {
				return fmt.Errorf("failed to copy .dockerignore: %w", err)
			}
		}
	}

	return nil
}

func (b *BaseDocker) copyFile(fileName string) error {
	src := path.Join(b.cfg.SourcePath, fileName+goboottypes.TemplateSuffix)

	// #nosec G304 -- this file path is safe and user-defined; used intentionally for scaffolding.
	content, err := os.ReadFile(src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("missing required template %q (expected %q)", fileName, src)
		}

		return fmt.Errorf("failed to read template file %q: %w", src, err)
	}

	err = gobootutils.WriteRootFile(b.root, fileName, content, goboottypes.FilePerm)
	if err != nil {
		return fmt.Errorf("failed to write template file %q: %w", fileName, err)
	}

	err = gobootutils.RenderTemplateToFile("docker_file", b.root, fileName, b.cfg)
	if err != nil {
		return fmt.Errorf("failed to render template to file: %w", err)
	}

	return nil
}

func (b *BaseDocker) registerScripts() error {
	hasDockerfile := b.hasEnabledFile(dockerfileFileEntry)
	hasCompose := b.hasEnabledFile(composeFileEntry)
	containerCheck := b.containerCheckCommand(hasDockerfile, hasCompose)
	cmds := []string{
		b.dockerBuildCommand(hasDockerfile),
		b.composeUpCommand(hasCompose),
		containerCheck,
	}

	err := b.script.RegisterLines(goboottypes.ServiceNameBaseDocker, cmds)
	if err != nil {
		return fmt.Errorf("failed to register script commands: %w", err)
	}

	err = b.script.RegisterFile(goboottypes.ScriptFileDocker, cmds)
	if err != nil {
		return fmt.Errorf("failed to register script file: %w", err)
	}

	return nil
}

func (b *BaseDocker) registerCIJobs() error {
	cmds := b.containerCheckCommands(
		b.hasEnabledFile(dockerfileFileEntry),
		b.hasEnabledFile(composeFileEntry),
	)

	err := b.ci.RegisterLines(goboottypes.ServiceNameBaseDocker, cmds)
	if err != nil {
		return fmt.Errorf("failed to register ci commands: %w", err)
	}

	err = b.ci.RegisterFile(goboottypes.CIFileContainer, cmds)
	if err != nil {
		return fmt.Errorf("failed to register ci job file: %w", err)
	}

	return nil
}

func (b *BaseDocker) hasEnabledFile(entry string) bool {
	for _, file := range b.cfg.FileList {
		if strings.TrimSpace(file) == entry {
			return true
		}
	}

	return false
}

func (b *BaseDocker) dockerBuildCommand(enabled bool) string {
	if !enabled {
		return "echo \"No Dockerfile output enabled.\""
	}

	return "docker build -t " + b.cfg.ImageName + " ."
}

func (b *BaseDocker) composeUpCommand(enabled bool) string {
	if !enabled {
		return "echo \"No compose output enabled.\""
	}

	return "docker compose up --build"
}

func (b *BaseDocker) containerCheckCommand(hasDockerfile, hasCompose bool) string {
	cmds := b.containerCheckCommands(hasDockerfile, hasCompose)
	if len(cmds) == 0 {
		return "echo \"No container checks registered.\""
	}

	return strings.Join(cmds, " && ")
}

func (b *BaseDocker) containerCheckCommands(hasDockerfile, hasCompose bool) []string {
	var cmds []string

	if hasCompose {
		cmds = append(cmds, "docker compose config")
	}

	if hasDockerfile {
		cmds = append(
			cmds,
			"docker build -t "+b.cfg.ImageName+" .",
			"docker run --rm "+b.cfg.ImageName+" -h",
		)
	}

	return cmds
}
