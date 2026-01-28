package cmd

import (
	"testing"

	"github.com/owenrumney/branch/internal/config"
	"github.com/spf13/cobra"
)

func TestNewRootCmd(t *testing.T) {
	// Helper to find command by name
	findCommand := func(rootCmd *cobra.Command, name string) *cobra.Command {
		for _, c := range rootCmd.Commands() {
			if c.Use == name+" [description...]" {
				return c
			}
		}
		return nil
	}

	t.Run("creates commands from default config", func(t *testing.T) {
		cfg := config.Default()
		rootCmd := NewRootCmd(cfg, "test")

		if rootCmd == nil {
			t.Fatal("NewRootCmd() should not return nil")
		}

		// Verify all default commands are present
		expectedCommands := []string{"feat", "fix", "tests", "chore", "docs"}
		for _, cmdName := range expectedCommands {
			if findCommand(rootCmd, cmdName) == nil {
				t.Errorf("Expected command %q to be created", cmdName)
			}
		}
	})

	t.Run("creates commands from custom config", func(t *testing.T) {
		cfg := &config.Config{
			BranchCommands: []string{"custom1", "custom2", "custom3"},
		}

		rootCmd := NewRootCmd(cfg, "test")

		if rootCmd == nil {
			t.Fatal("NewRootCmd() should not return nil")
		}

		// Verify custom commands are present
		expectedCommands := []string{"custom1", "custom2", "custom3"}
		for _, cmdName := range expectedCommands {
			if findCommand(rootCmd, cmdName) == nil {
				t.Errorf("Expected command %q to be created", cmdName)
			}
		}
	})

	t.Run("hides completion command", func(t *testing.T) {
		cfg := config.Default()
		rootCmd := NewRootCmd(cfg, "test")

		if !rootCmd.CompletionOptions.HiddenDefaultCmd {
			t.Error("Completion command should be hidden")
		}
	})

	t.Run("root command has correct metadata", func(t *testing.T) {
		cfg := config.Default()
		rootCmd := NewRootCmd(cfg, "test")

		if rootCmd.Use != "branch" {
			t.Errorf("Expected Use to be 'branch', got %q", rootCmd.Use)
		}

		if rootCmd.Short == "" {
			t.Error("Root command should have a Short description")
		}

		if rootCmd.Long == "" {
			t.Error("Root command should have a Long description")
		}
	})
}

func TestNewRootCmdCommandStructure(t *testing.T) {
	cfg := config.Default()
	rootCmd := NewRootCmd(cfg, "test")

	// Helper to find command by name
	findCommand := func(name string) *cobra.Command {
		for _, c := range rootCmd.Commands() {
			if c.Use == name+" [description...]" {
				return c
			}
		}
		return nil
	}

	// Test that each command has the correct structure
	for _, branchCmd := range cfg.BranchCommands {
		foundCmd := findCommand(branchCmd)

		if foundCmd == nil {
			t.Errorf("Command for %q not found", branchCmd)
			continue
		}

		// Verify command structure
		if foundCmd.Short == "" {
			t.Errorf("Command %q should have a Short description", branchCmd)
		}

		if foundCmd.Long == "" {
			t.Errorf("Command %q should have a Long description", branchCmd)
		}

		// Verify minimum args requirement
		if foundCmd.Args == nil {
			t.Errorf("Command %q should have Args validation", branchCmd)
		}
	}
}
