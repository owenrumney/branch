package cmd

import (
	"fmt"
	"os"

	"github.com/owenrumney/branch/internal/branch"
	"github.com/owenrumney/branch/internal/config"
	"github.com/owenrumney/branch/internal/git"
	"github.com/spf13/cobra"
)

func newBranchCmd(branchType, description string) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s [description...]", branchType),
		Short: description,
		Long: fmt.Sprintf(`%s

If the first word matches a known ticket pattern (e.g., PIP-1234, #123), it will be included in the branch name.

Examples:
  branch %s PIP-1234 implement new feature  ->  %s/pip-1234-implement-new-feature
  branch %s implement new feature           ->  %s/implement-new-feature`, description, branchType, branchType, branchType, branchType),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not load config: %v\n", err)
				cfg = config.Default()
			}

			ticket, descParts := parseArgs(args, cfg)
			branchName := branch.Generate(branchType, ticket, descParts)

			if err := git.CreateBranch(branchName); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating branch: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Created and switched to branch: %s\n", branchName)
		},
	}
}

func parseArgs(args []string, cfg *config.Config) (ticket string, description []string) {
	if len(args) == 0 {
		return "", nil
	}

	first := args[0]
	if cfg.IsTicket(first) {
		return first, args[1:]
	}

	return "", args
}
