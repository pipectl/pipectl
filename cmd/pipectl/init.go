package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/pipectl/pipectl/internal/wizard"
)

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "Create a new pipeline YAML interactively",
	Long: `Start an interactive wizard that asks a few questions and writes a
pipeline YAML file with placeholder values ready to customise.`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := wizard.Run()
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		if err != nil {
			return err
		}

		yaml := wizard.Render(result)

		if result.OutputFile == "" {
			fmt.Print(yaml)
			return nil
		}

		if err := os.WriteFile(result.OutputFile, []byte(yaml), 0644); err != nil {
			return fmt.Errorf("write %s: %w", result.OutputFile, err)
		}

		fmt.Fprintf(os.Stderr, "Wrote %s\n", result.OutputFile)
		return nil
	},
}

func init() {
	rootCommand.AddCommand(initCommand)
}
