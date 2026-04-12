package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/pipectl/pipectl/internal/pipeline"
)

var outputPath string
var verbose bool
var dryRun bool

var runCommand = &cobra.Command{
	Use:   "run pipeline.yaml",
	Short: "Run a pipeline",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		var input []byte
		output := io.Writer(os.Stdout)

		stat, err := os.Stdin.Stat()
		if err != nil {
			return err
		}

		// If stdin is NOT a terminal, read from it
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// TODO implement some kind of limit on data size
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
		}

		if outputPath != "" {
			file, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("open output file: %w", err)
			}
			defer file.Close()
			output = file
		}

		if err := pipeline.Run(path, input, output, verbose, dryRun); err != nil {
			return fmt.Errorf("pipeline failed: %w", err)
		}

		return nil
	},
}

func init() {
	runCommand.Flags().StringVarP(&outputPath, "output", "o", "", "Write pipeline output to file")
	runCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	runCommand.Flags().BoolVar(&dryRun, "dry-run", false, "Validate and print the pipeline plan without executing")
	rootCommand.AddCommand(runCommand)
}
