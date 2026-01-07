package cmd

import (
	"fmt"

	"github.com/owenrumney/branch/internal/config"
	"github.com/spf13/cobra"
)

func NewRootCmd(cfg *config.Config) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "branch",
		Short: "Create git branches with consistent naming patterns",
		Long:  `A CLI tool for creating git branches using a standardized pattern: <type>/<ticket>-<description>`,
	}

	for _, branchCommand := range cfg.BranchCommands {
		rootCmd.AddCommand(newBranchCmd(branchCommand, fmt.Sprintf("Create a %s branch", branchCommand)))
	}

	// don't need completions so lets remove that from the help output
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	return rootCmd
}
