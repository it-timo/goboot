/*
Package main initializes and executes the goboot CLI.

It loads the main YAML configuration, registers all enabled services, and executes each one in order.

Errors during any stage cause early termination.
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboot"
)

// configPath holds the path to the goboot configuration file.
//
// Can be overridden using --config CLI flag.
var configPath string

// init parses CLI flags before the main execution begins.
func init() {
	flag.StringVar(&configPath, "config", "./configs/goboot.yml", "Path to the goboot config file")
	flag.Parse()
}

func main() {
	// Step 1: Load and validate goboot configuration from YAML.
	cfg := config.NewGoBoot(configPath)

	err := cfg.Init()
	if err != nil {
		fmt.Printf("Failed to initialize configuration: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded configuration: %+#v\n", cfg)

	// Step 2: Create a new goboot application instance.
	app := goboot.NewGoBoot(cfg)

	// Step 3: Register all declared and enabled services.
	err = app.RegisterServices()
	if err != nil {
		fmt.Printf("Service registration failed: %s\n", err)
		os.Exit(1)
	}

	// Step 4: Execute all registered services.
	err = app.RunServices()
	if err != nil {
		fmt.Printf("Service execution failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("goboot execution completed successfully.")
}
