package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pipectl/pipectl/internal/pipeline"
)

var inputPath string
var outputPath string
var verbose bool
var dryRun bool
var varFlags []string

var runCommand = &cobra.Command{
	Use:          "run pipeline.yaml",
	Short:        "Run a pipeline",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		vars, err := parseVars(varFlags)
		if err != nil {
			return err
		}

		var input []byte
		var outputFile *os.File
		output := io.Writer(os.Stdout)

		stat, err := os.Stdin.Stat()
		if err != nil {
			return err
		}
		stdinPiped := (stat.Mode() & os.ModeCharDevice) == 0

		if inputPath != "" {
			if stdinPiped {
				return fmt.Errorf("cannot use --input flag and piped stdin together")
			}
			input, err = os.ReadFile(inputPath)
			if err != nil {
				return fmt.Errorf("open input file: %w", err)
			}
		} else if stdinPiped {
			// TODO implement some kind of limit on data size
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
		}

		if outputPath != "" {
			outputFile, err = os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("open output file: %w", err)
			}
			output = outputFile
		}

		if err := pipeline.Run(path, input, output, verbose, dryRun, vars); err != nil {
			return fmt.Errorf("pipeline failed: %w", err)
		}

		if outputFile != nil {
			if err := outputFile.Close(); err != nil {
				return fmt.Errorf("close output file: %w", err)
			}
		}

		return nil
	},
}

func parseVars(flags []string) (map[string]string, error) {
	if len(flags) == 0 {
		return nil, nil
	}
	vars := make(map[string]string, len(flags))
	for _, f := range flags {
		k, v, ok := strings.Cut(f, "=")
		if !ok || k == "" {
			return nil, fmt.Errorf("--var %q: expected KEY=VALUE", f)
		}
		vars[k] = v
	}
	return vars, nil
}

func init() {
	runCommand.Flags().StringVarP(&inputPath, "input", "i", "", "Read pipeline input from file")
	runCommand.Flags().StringVarP(&outputPath, "output", "o", "", "Write pipeline output to file")
	runCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	runCommand.Flags().BoolVar(&dryRun, "dry-run", false, "Validate and print the pipeline plan without executing")
	runCommand.Flags().StringArrayVar(&varFlags, "var", nil, "Substitute ${VAR} tokens in pipeline YAML (repeatable, KEY=VALUE)")
	rootCommand.AddCommand(runCommand)
}
