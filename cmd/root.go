package cmd

import (
	"fmt"

	"github.com/owenrumney/branch/internal/config"
	"github.com/spf13/cobra"
)

func NewRootCmd(cfg *config.Config, version string) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:     "branch",
		Short:   "Create git branches with consistent naming patterns",
		Long:    `A CLI tool for creating git branches using a standardized pattern: <type>/<ticket>-<description>`,
		Version: version,
	}

	for _, branchCommand := range cfg.BranchCommands {
		rootCmd.AddCommand(newBranchCmd(branchCommand, fmt.Sprintf("Create a %s branch", branchCommand)))
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	return rootCmd
}
