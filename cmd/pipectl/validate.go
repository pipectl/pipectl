package main

import (
	"github.com/spf13/cobra"

	"github.com/pipectl/pipectl/internal/pipeline/plan"
	"github.com/pipectl/pipectl/internal/pipeline/spec"
)

var validateCommand = &cobra.Command{
	Use:          "validate pipeline.yaml",
	Short:        "Validate a pipeline file",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := spec.Load(args[0], nil)
		if err != nil {
			return err
		}
		_, err = plan.Build(p)
		return err
	},
}

func init() {
	rootCommand.AddCommand(validateCommand)
}
