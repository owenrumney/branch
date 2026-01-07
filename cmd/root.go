package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "branch",
	Short: "Create git branches with consistent naming patterns",
	Long:  `A CLI tool for creating git branches using a standardized pattern: <type>/<ticket>-<description>`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(newBranchCmd("feat", "Create a feature branch"))
	rootCmd.AddCommand(newBranchCmd("fix", "Create a bugfix branch"))
	rootCmd.AddCommand(newBranchCmd("tests", "Create a tests branch"))
	rootCmd.AddCommand(newBranchCmd("chore", "Create a chore branch"))
	rootCmd.AddCommand(newBranchCmd("docs", "Create a documentation branch"))
}
