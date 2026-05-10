package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"

	stepdocs "github.com/pipectl/pipectl/website/docs/steps"
)

var docsCommand = &cobra.Command{
	Use:          "docs [step]",
	Short:        "Show documentation for pipeline steps",
	Long:         "Show documentation for a specific step, or list all available steps.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listSteps(cmd)
		}
		return showStep(args[0])
	},
}

func listSteps(cmd *cobra.Command) error {
	entries, err := fs.ReadDir(stepdocs.FS, ".")
	if err != nil {
		return err
	}

	type row struct {
		name string
		desc string
	}
	var steps []row
	maxLen := 0

	for _, e := range entries {
		fname := e.Name()
		if !strings.HasSuffix(fname, ".md") || fname == "index.md" {
			continue
		}
		name := strings.TrimSuffix(fname, ".md")
		data, err := fs.ReadFile(stepdocs.FS, fname)
		if err != nil {
			continue
		}
		steps = append(steps, row{name, extractDescription(string(data))})
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Available steps:")
	fmt.Fprintln(cmd.OutOrStdout())
	for _, s := range steps {
		fmt.Fprintf(cmd.OutOrStdout(), "  %-*s  %s\n", maxLen, s.name, s.desc)
	}
	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintln(cmd.OutOrStdout(), "Run 'pipectl docs <step>' for full documentation.")
	return nil
}

func showStep(name string) error {
	data, err := fs.ReadFile(stepdocs.FS, name+".md")
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("unknown step: %s\n\nRun 'pipectl docs' to list available steps.", name)
	}
	if err != nil {
		return err
	}
	if isTerminal(os.Stdout) {
		r, rerr := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(100))
		if rerr == nil {
			rendered, rerr := r.Render(string(data))
			if rerr == nil {
				fmt.Print(rendered)
				return nil
			}
		}
	}
	_, err = os.Stdout.Write(data)
	return err
}

func isTerminal(f *os.File) bool {
	stat, err := f.Stat()
	return err == nil && (stat.Mode()&os.ModeCharDevice) != 0
}

func extractDescription(content string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		return line
	}
	return ""
}

func init() {
	rootCommand.AddCommand(docsCommand)
}
