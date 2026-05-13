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
	"time"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboot"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	exitFunc               = os.Exit
	outputWriter io.Writer = os.Stdout
)

func writeOutputLine(msg string) {
	_, err := io.WriteString(outputWriter, msg+"\n")
	if err != nil {
		log.Error().Err(err).Msg("failed to write to CLI output")
	}
}

// run executes the whole goboot CLI with config load, app init, service registration and execution.
func run(args []string) error {
	// Step 0: Parse flags explicitly using a local FlagSet to avoid global state.
	flagSet := flag.NewFlagSet("goboot", flag.ContinueOnError)
	configPath := ""
	logLevel := "info"

	flagSet.StringVar(&configPath, "config", "./configs/goboot.yml", "Path to the goboot config file")
	flagSet.StringVar(&logLevel, "log-level", "info", "Log level: debug|info|warn|error")

	err := flagSet.Parse(args)
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("invalid --log-level value %q (allowed: debug, info, warn, error)", logLevel)
	}

	zerolog.SetGlobalLevel(level)

	logger := zerolog.New(os.Stderr).With().Timestamp().Str("app", "goboot").Logger()
	log.Logger = logger
	logger.Info().Str("log_level", logLevel).Msg("logger configured")

	runStart := time.Now()

	// Step 1: Load and validate goboot configuration from YAML.
	cfg := config.NewGoBoot(configPath)
	cfg.SetLogger(logger.With().Str("component", "config").Logger())

	err = cfg.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}

	// Step 2: Create a new goboot application instance.
	app := goboot.NewGoBootWithLogger(cfg, logger.With().Str("component", "orchestrator").Logger())

	// Step 3: Register all declared and enabled services.
	err = app.RegisterServices()
	if err != nil {
		return fmt.Errorf("service registration failed: %w", err)
	}

	// Step 4: Execute all registered services.
	err = app.RunServices()
	if err != nil {
		return fmt.Errorf("service execution failed: %w", err)
	}

	// Step 5: Run go mod tidy if the go.mod file exists.
	err = app.RunGoModTidy(true)
	if err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	logger.Info().
		Str("config_path", configPath).
		Str("target_path", cfg.TargetPath).
		Str("project_name", cfg.ProjectName).
		Int64("duration_ms", time.Since(runStart).Milliseconds()).
		Msg("goboot execution completed successfully")
	writeOutputLine("goboot execution completed successfully.")

	return nil
}

func main() {
	err := run(os.Args[1:])
	if err != nil {
		log.Error().Err(err).Msg("goboot failed")
		writeOutputLine(err.Error())

		exitFunc(1)
	}
}
