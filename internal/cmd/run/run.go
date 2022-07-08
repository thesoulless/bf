// Package run provides functionality for the `run` command of the CLI
package run

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thesoulless/bf"
)

// Cmd is the command for running the BF commands
func Cmd() *cobra.Command {
	var file string
	var s string

	cmd := &cobra.Command{
		Use:   "run [-s \"bf commands\"] | [-f file_path]",
		Args:  cobra.ExactArgs(0),
		Short: "Runs BF commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), s, file)
		},
	}

	cmd.Flags().StringVarP(&file,
		"file", "f", "", "BF file path")

	cmd.Flags().StringVarP(&s,
		"string", "s", "", "BF commands")

	return cmd
}

func run(ctx context.Context, s string, file string) error {
	if s != "" {
		err := bf.Run(strings.NewReader(s))

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	fb, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("faild to read file: %w", err)
	}

	err = bf.Run(bytes.NewReader(fb))

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return err
}