package main

import (
	"fmt"
	"os"

	"github.com/danielmiessler/fabric/core"
	"github.com/danielmiessler/fabric/cli"
)

// Version information - set via ldflags during build
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Initialize the CLI application
	app := cli.NewApp()
	app.Version = version

	// Set up the core fabric instance
	fabric, err := core.NewFabric()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing Fabric: %v\n", err)
		os.Exit(1)
	}

	// Register the fabric instance with the CLI
	app.SetFabric(fabric)

	// Run the CLI application
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
