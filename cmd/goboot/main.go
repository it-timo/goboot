/*
Package main initializes and executes the goboot CLI.

It loads the main YAML configuration, registers all enabled services, and executes each one in order.

Errors during any stage cause early termination.
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboot"
	"github.com/it-timo/goboot/pkg/regeneration"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const stagingProjectMode = 0o750

var (
	exitFunc               = os.Exit
	outputWriter io.Writer = os.Stdout
	version                = "dev"
)

type cliOptions struct {
	configPath         string
	logLevel           string
	regenerationPolicy string
	dryRun             bool
	skipGoModTidy      bool
}

func writeOutputLine(msg string) {
	_, err := io.WriteString(outputWriter, msg+"\n")
	if err != nil {
		log.Error().Err(err).Msg("failed to write to CLI output")
	}
}

func parseCLIOptions(args []string) (cliOptions, error) {
	flagSet := flag.NewFlagSet("goboot", flag.ContinueOnError)
	opts := cliOptions{
		configPath: "./configs/goboot.yml",
		logLevel:   "info",
	}

	flagSet.StringVar(&opts.configPath, "config", opts.configPath, "Path to the goboot config file")
	flagSet.StringVar(&opts.logLevel, "log-level", opts.logLevel, "Log level: debug|info|warn|error")
	flagSet.StringVar(
		&opts.regenerationPolicy,
		"regeneration-policy",
		"",
		"Override regeneration policy: managed|replace|preserve",
	)
	flagSet.BoolVar(&opts.dryRun, "dry-run", false, "Plan generation without changing the target project")
	flagSet.BoolVar(&opts.skipGoModTidy, "skip-go-mod-tidy", false, "Skip running go mod tidy after generation")

	err := flagSet.Parse(args)
	if err != nil {
		return cliOptions{}, fmt.Errorf("failed to parse flags: %w", err)
	}

	return opts, nil
}

// run executes the whole goboot CLI with config load, app init, service registration and execution.
func run(args []string) error {
	// Step 0: Parse flags explicitly using a local FlagSet to avoid global state.
	opts, err := parseCLIOptions(args)
	if err != nil {
		return err
	}

	level, err := zerolog.ParseLevel(opts.logLevel)
	if err != nil {
		return fmt.Errorf("invalid --log-level value %q (allowed: debug, info, warn, error)", opts.logLevel)
	}

	zerolog.SetGlobalLevel(level)

	logger := zerolog.New(os.Stderr).With().Timestamp().Str("app", "goboot").Logger()
	log.Logger = logger
	logger.Info().Str("log_level", opts.logLevel).Msg("logger configured")

	runStart := time.Now()

	// Step 1: Load and validate goboot configuration from YAML.
	cfg := config.NewGoBoot(opts.configPath)
	cfg.SetLogger(logger.With().Str("component", "config").Logger())

	err = cfg.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}

	err = executeGeneration(opts, cfg, logger)
	if err != nil {
		return err
	}

	logger.Info().
		Str("config_path", opts.configPath).
		Str("target_path", cfg.TargetPath).
		Str("project_name", cfg.ProjectName).
		Int64("duration_ms", time.Since(runStart).Milliseconds()).
		Msg("goboot execution completed successfully")
	writeOutputLine("goboot execution completed successfully.")

	return nil
}

func executeGeneration(opts cliOptions, cfg *config.GoBoot, logger zerolog.Logger) error {
	policy, err := requestedPolicy(opts.regenerationPolicy, cfg.RegenerationPolicy)
	if err != nil {
		return err
	}

	targetPath := cfg.TargetPath

	stagingRoot, err := os.MkdirTemp("", "goboot-generation-*")
	if err != nil {
		return fmt.Errorf("failed to create generation staging directory: %w", err)
	}

	defer func() {
		if removeErr := os.RemoveAll(stagingRoot); removeErr != nil {
			logger.Error().Err(removeErr).Msg("failed to remove generation staging directory")
		}
	}()

	cfg.TargetPath = stagingRoot
	defer func() {
		cfg.TargetPath = targetPath
	}()

	err = generateStagedProject(opts, cfg, logger, stagingRoot)
	if err != nil {
		return err
	}

	return applyStagedProject(opts, cfg, policy, stagingRoot, targetPath)
}

func requestedPolicy(override, configured string) (regeneration.Policy, error) {
	policyRaw := configured
	if override != "" {
		policyRaw = override
	}

	policy, err := regeneration.ParsePolicy(policyRaw)
	if err != nil {
		return "", fmt.Errorf("invalid regeneration policy: %w", err)
	}

	return policy, nil
}

func generateStagedProject(opts cliOptions, cfg *config.GoBoot, logger zerolog.Logger, stagingRoot string) error {
	// Create a new goboot application instance for the isolated staging tree.
	app := goboot.NewGoBootWithLogger(cfg, logger.With().Str("component", "orchestrator").Logger())

	err := app.RegisterServices()
	if err != nil {
		return fmt.Errorf("service registration failed: %w", err)
	}

	err = app.RunServices()
	if err != nil {
		return fmt.Errorf("service execution failed: %w", err)
	}

	err = os.MkdirAll(filepath.Join(stagingRoot, cfg.ProjectName), stagingProjectMode)
	if err != nil {
		return fmt.Errorf("failed to ensure staged project root: %w", err)
	}

	err = app.RunGoModTidy(!opts.skipGoModTidy)
	if err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	return nil
}

func applyStagedProject(
	opts cliOptions,
	cfg *config.GoBoot,
	policy regeneration.Policy,
	stagingRoot string,
	targetPath string,
) error {
	plan, err := regeneration.Apply(regeneration.Request{
		StagedProject:    filepath.Join(stagingRoot, cfg.ProjectName),
		TargetProject:    filepath.Join(targetPath, cfg.ProjectName),
		GeneratorVersion: version,
		Profile:          cfg.Profile,
		Services:         enabledServiceIDs(cfg.Services),
		Policy:           policy,
		DryRun:           opts.dryRun,
	})
	if opts.dryRun || err != nil {
		writeOutputLine(plan.String())
	}

	if err != nil {
		return fmt.Errorf("failed to apply generation transaction: %w", err)
	}

	if opts.dryRun {
		writeOutputLine("goboot dry run completed without changing the target project.")
	}

	return nil
}

func enabledServiceIDs(services []config.ServiceConfigMeta) []string {
	serviceIDs := make([]string, 0, len(services))

	for _, service := range services {
		if service.IsEnabled() {
			serviceIDs = append(serviceIDs, service.ID)
		}
	}

	return serviceIDs
}

func main() {
	err := run(os.Args[1:])
	if err != nil {
		log.Error().Err(err).Msg("goboot failed")
		writeOutputLine(err.Error())

		exitFunc(1)
	}
}
