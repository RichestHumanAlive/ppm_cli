package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/ppm/cmd"
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
		Long: `PPM is a platform-agnostic CLI tool that unifies package management across 
multiple package managers like npm, pip, scoop, and winget.`,
	}

	// Add commands
	rootCmd.AddCommand(
		cmd.NewInstallCmd(),
	)

	return rootCmd.Execute()
}
