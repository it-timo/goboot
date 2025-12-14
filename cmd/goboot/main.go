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

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboot"
)

var (
	exitFunc               = os.Exit
	outputWriter io.Writer = os.Stdout
)

// run executes the whole goboot CLI with config load, app init, service registration and execution.
func run(args []string) error {
	// Step 0: Parse flags explicitly using a local FlagSet to avoid global state.
	fs := flag.NewFlagSet("goboot", flag.ContinueOnError)
	configPath := ""

	fs.StringVar(&configPath, "config", "./configs/goboot.yml", "Path to the goboot config file")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	// Step 1: Load and validate goboot configuration from YAML.
	cfg := config.NewGoBoot(configPath)

	err := cfg.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}

	_, err = fmt.Fprintf(outputWriter, "Loaded configuration: %+#v\n", cfg)
	if err != nil {
		fmt.Println("Failed to write error to output:", err)
	}

	// Step 2: Create a new goboot application instance.
	app := goboot.NewGoBoot(cfg)

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

	_, err = fmt.Fprintln(outputWriter, "goboot execution completed successfully.")
	if err != nil {
		fmt.Println("Failed to write error to output:", err)
	}

	return nil
}

func main() {
	err := run(os.Args[1:])
	if err != nil {
		_, err = fmt.Fprintf(outputWriter, "%s\n", err)
		if err != nil {
			fmt.Println("Failed to write error to output:", err)
		}

		exitFunc(1)
	}
}
