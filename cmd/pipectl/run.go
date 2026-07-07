package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pipectl/pipectl/internal/pipeline"
)

var inputPath string
var outputPath string
var verbose bool
var quiet bool
var dryRun bool
var timing bool
var varFlags []string
var maxInputSizeStr string

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

		maxInputSize, err := parseByteSize(maxInputSizeStr)
		if err != nil {
			return fmt.Errorf("--max-input-size: %w", err)
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
			fi, err := os.Stat(inputPath)
			if err != nil {
				return fmt.Errorf("open input file: %w", err)
			}
			if fi.Size() > maxInputSize {
				return errInputTooLarge(fmt.Sprintf("input file %q", inputPath), maxInputSizeStr)
			}
			input, err = os.ReadFile(inputPath)
			if err != nil {
				return fmt.Errorf("open input file: %w", err)
			}
		} else if stdinPiped {
			input, err = io.ReadAll(io.LimitReader(os.Stdin, limitPlusOne(maxInputSize)))
			if err != nil {
				return err
			}
			if int64(len(input)) > maxInputSize {
				return errInputTooLarge("stdin input", maxInputSizeStr)
			}
		}

		if outputPath != "" {
			outputFile, err = os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("open output file: %w", err)
			}
			output = outputFile
		}

		if err := pipeline.Run(path, input, output, verbose, dryRun, quiet, timing, vars); err != nil {
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

// errInputTooLarge reports that a pipeline input source exceeded --max-input-size.
func errInputTooLarge(source, limit string) error {
	return fmt.Errorf("%s exceeds maximum input size %s (adjust with --max-input-size)", source, limit)
}

// limitPlusOne returns limit+1, saturating at math.MaxInt64 instead of
// overflowing to a negative number that io.LimitReader would treat as EOF.
func limitPlusOne(limit int64) int64 {
	if limit == math.MaxInt64 {
		return limit
	}
	return limit + 1
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
	runCommand.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress all diagnostic output")
	runCommand.Flags().BoolVar(&dryRun, "dry-run", false, "Validate and print the pipeline plan without executing")
	runCommand.Flags().BoolVar(&timing, "timing", false, "Print per-step timing table to stderr after execution")
	runCommand.Flags().StringArrayVar(&varFlags, "var", nil, "Substitute ${VAR} tokens in pipeline YAML (repeatable, KEY=VALUE)")
	runCommand.Flags().StringVar(&maxInputSizeStr, "max-input-size", "256MB", "Maximum size of pipeline input read from stdin or --input (e.g. 64KB, 256MB, 1GB); reject input exceeding this size")
	rootCommand.AddCommand(runCommand)
}
