package main

import (
	"os"

	"github.com/spf13/cobra"
)

// Set by GoReleaser via ldflags at build time.
var version = "dev"

var rootCommand = &cobra.Command{
	Use:     "pipectl [command]",
	Short:   "pipectl runs deterministic data pipelines",
	Version: version,
	Long: `pipectl executes declarative pipelines defined in YAML or JSON.

Each pipeline runs step-by-step, passing data from one stage to the next.`,
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
