// Package cmd package cmd provides CLI capabilities
package cmd

import (
	"context"

	"github.com/thesoulless/bf/internal/cmd/run"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Example: "run [-s \"bf commands\"] | [-f file_path]",
	Short:   "BF is a Brainfuck interpreter.",
	Long:    `bf is a CLI tool for running Brainfuck commands`,
}

// Execute runs the app with the given context
func Execute(ctx context.Context) int {
	err := runCmd(ctx)
	if err == nil {
		return 0
	}

	return 1
}

func runCmd(ctx context.Context) error {
	rootCmd.Flags().Bool("version", false, "show bf version")

	rootCmd.AddCommand(run.Cmd())

	return rootCmd.ExecuteContext(ctx)
}
