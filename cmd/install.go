package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [package]",
		Short: "Install a package using the appropriate package manager",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement package installation logic
			fmt.Printf("Installing package: %s\n", args[0])
			return nil
		},
	}

	return cmd
}
