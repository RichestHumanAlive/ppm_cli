package main

import (
	"fmt"
	"os"

	"github.com/RichestHumanAlive/ppm_cli/cmd"
	"github.com/spf13/cobra"
)

func main() {
	if err := execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func execute() error {
	rootCmd := &cobra.Command{
		Use:   "ppm",
		Short: "Panoramic Package Manager - A unified package management CLI",
		Long: `PPM is a platform-agnostic CLI tool designed to unify the workflows of multiple package managers.
It supports npm, pip, and scoop, providing a consistent interface for managing packages across different ecosystems.`,
	}

	// Add commands
	rootCmd.AddCommand(
		cmd.NewInstallCmd(),
		cmd.NewSearchCmd(),
	)

	return rootCmd.Execute()
}
