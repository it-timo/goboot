/*
Package main initializes and executes the goboot CLI.

It loads the main YAML configuration, registers all enabled services, and executes each one in order.

Errors during any stage cause early termination.
*/
package main

import (
	"errors"
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
	outputFormat       string
	regenerationPolicy string
	dryRun             bool
	skipGoModTidy      bool
	validateOnly       bool
	showVersion        bool
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
	flagSet.StringVar(&opts.outputFormat, "output", "human", "Output format: human|json")
	flagSet.StringVar(
		&opts.regenerationPolicy,
		"regeneration-policy",
		"",
		"Override regeneration policy: managed|replace|preserve",
	)
	flagSet.BoolVar(&opts.dryRun, "dry-run", false, "Plan generation without changing the target project")
	flagSet.BoolVar(&opts.skipGoModTidy, "skip-go-mod-tidy", false, "Skip running go mod tidy after generation")
	flagSet.BoolVar(&opts.validateOnly, "validate", false, "Validate root and service configs without generation")
	flagSet.BoolVar(&opts.showVersion, "version", false, "Print the goboot version and exit")

	err := flagSet.Parse(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			flagSet.SetOutput(outputWriter)
			flagSet.Usage()

			return cliOptions{}, flag.ErrHelp
		}

		return cliOptions{}, fmt.Errorf("failed to parse flags: %w", err)
	}

	if flagSet.NArg() != 0 {
		return cliOptions{}, fmt.Errorf("unexpected positional arguments: %v", flagSet.Args())
	}

	err = validateCLIOptions(opts)
	if err != nil {
		return cliOptions{}, err
	}

	return opts, nil
}

func validateCLIOptions(opts cliOptions) error {
	_, err := parseOutputFormat(opts.outputFormat)
	if err != nil {
		return err
	}

	if opts.showVersion && (opts.validateOnly || opts.dryRun || opts.skipGoModTidy || opts.regenerationPolicy != "") {
		return errors.New("--version cannot be combined with generation or validation flags")
	}

	if opts.validateOnly && (opts.dryRun || opts.skipGoModTidy || opts.regenerationPolicy != "") {
		return errors.New("--validate cannot be combined with generation-only flags")
	}

	return nil
}

// run executes the whole goboot CLI with config load, app init, service registration and execution.
func run(args []string) error {
	// Step 0: Parse flags explicitly using a local FlagSet to avoid global state.
	opts, err := parseCLIOptions(args)
	if errors.Is(err, flag.ErrHelp) {
		return nil
	}

	if err != nil {
		return newCommandError(exitUsage, categoryUsage, requestedOperation(args), err, nil)
	}

	if opts.showVersion {
		writeSuccess(opts, successResult{
			Operation: operationVersion,
			Message:   "goboot " + version,
		})

		return nil
	}

	level, err := zerolog.ParseLevel(opts.logLevel)
	if err != nil {
		parseErr := fmt.Errorf("invalid --log-level value %q (allowed: debug, info, warn, error)", opts.logLevel)

		return newCommandError(exitUsage, categoryUsage, operationForOptions(opts), parseErr, nil)
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
		configErr := fmt.Errorf("failed to initialize configuration: %w", err)

		return newCommandError(exitConfig, categoryConfig, operationForOptions(opts), configErr, nil)
	}

	if opts.validateOnly {
		writeSuccess(opts, successResult{
			Operation: operationValidate,
			Config:    opts.configPath,
			Project:   cfg.ProjectName,
			Message:   "configuration is valid.",
		})

		return nil
	}

	plan, err := executeGeneration(opts, cfg, logger)
	if err != nil {
		return err
	}

	logger.Info().
		Str("config_path", opts.configPath).
		Str("target_path", cfg.TargetPath).
		Str("project_name", cfg.ProjectName).
		Int64("duration_ms", time.Since(runStart).Milliseconds()).
		Msg("goboot execution completed successfully")

	result := successResult{
		Operation: operationGenerate,
		Config:    opts.configPath,
		Project:   cfg.ProjectName,
		Message:   "goboot execution completed successfully.",
	}
	if opts.dryRun {
		result.Operation = operationDryRun
		result.Message = "goboot dry run completed without changing the target project."
		result.Plan = &plan
	}

	writeSuccess(opts, result)

	return nil
}

func executeGeneration(opts cliOptions, cfg *config.GoBoot, logger zerolog.Logger) (regeneration.Plan, error) {
	policy, err := requestedPolicy(opts.regenerationPolicy, cfg.RegenerationPolicy)
	if err != nil {
		return regeneration.Plan{}, newCommandError(
			exitUsage,
			categoryUsage,
			operationForOptions(opts),
			err,
			nil,
		)
	}

	targetPath := cfg.TargetPath

	stagingRoot, err := os.MkdirTemp("", "goboot-generation-*")
	if err != nil {
		stagingErr := fmt.Errorf("failed to create generation staging directory: %w", err)

		return regeneration.Plan{}, newCommandError(
			exitGeneration,
			categoryGeneration,
			operationForOptions(opts),
			stagingErr,
			nil,
		)
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
		return regeneration.Plan{}, newCommandError(
			exitGeneration,
			categoryGeneration,
			operationForOptions(opts),
			err,
			nil,
		)
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
) (regeneration.Plan, error) {
	plan, err := regeneration.Apply(regeneration.Request{
		StagedProject:    filepath.Join(stagingRoot, cfg.ProjectName),
		TargetProject:    filepath.Join(targetPath, cfg.ProjectName),
		GeneratorVersion: version,
		Profile:          cfg.Profile,
		Services:         enabledServiceIDs(cfg.Services),
		Policy:           policy,
		DryRun:           opts.dryRun,
	})
	if err != nil {
		applyErr := fmt.Errorf("failed to apply generation transaction: %w", err)
		code := exitGeneration
		category := categoryGeneration
		if regeneration.IsConflict(err) {
			code = exitConflict
			category = categoryConflict
		}

		return plan, newCommandError(code, category, operationForOptions(opts), applyErr, &plan)
	}

	return plan, nil
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
		writeCommandError(os.Args[1:], err)

		exitFunc(commandExitCode(err))
	}
}
