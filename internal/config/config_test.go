package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg == nil {
		t.Fatal("Default() returned nil")
	}

	if len(cfg.TicketPatterns) == 0 {
		t.Error("Default() should have ticket patterns")
	}

	// Test that default patterns are present
	expectedPatterns := []string{
		`^#\d+$`,
		`^[A-Z]+-\d+$`,
		`^[A-Z]+_\d+$`,
	}

	if len(cfg.TicketPatterns) != len(expectedPatterns) {
		t.Errorf("Expected %d patterns, got %d", len(expectedPatterns), len(cfg.TicketPatterns))
	}

	// Test that default branch commands are present
	expectedCommands := []string{"feat", "fix", "tests", "chore", "docs"}
	if len(cfg.BranchCommands) != len(expectedCommands) {
		t.Errorf("Expected %d branch commands, got %d", len(expectedCommands), len(cfg.BranchCommands))
	}

	for i, expected := range expectedCommands {
		if i >= len(cfg.BranchCommands) {
			t.Errorf("Missing branch command: %q", expected)
			continue
		}
		if cfg.BranchCommands[i] != expected {
			t.Errorf("BranchCommands[%d] = %q, want %q", i, cfg.BranchCommands[i], expected)
		}
	}

	// Test that patterns are compiled
	if len(cfg.compiled) == 0 {
		t.Error("Default() should compile patterns")
	}
}

func TestIsTicket(t *testing.T) {
	cfg := Default()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// GitHub issues
		{"GitHub issue #123", "#123", true},
		{"GitHub issue #1", "#1", true},
		{"GitHub issue #9999", "#9999", true},
		{"not GitHub issue", "123", false},
		{"not GitHub issue with text", "#abc", false},

		// Jira/Linear style with hyphen
		{"Linear ticket PIP-1234", "PIP-1234", true},
		{"Linear ticket INFRA-124", "INFRA-124", true},
		{"Jira ticket PROJ-1", "PROJ-1", true},
		{"lowercase ticket", "pip-1234", false},
		{"ticket with text", "PIP-abc", false},

		// Underscore variant
		{"Underscore ticket PIP_1234", "PIP_1234", true},
		{"Underscore ticket ABC_999", "ABC_999", true},
		{"lowercase underscore", "pip_1234", false},

		// Invalid patterns
		{"empty string", "", false},
		{"plain text", "fix bug", false},
		{"just numbers", "1234", false},
		{"just letters", "PIP", false},
		{"mixed case", "Pip-1234", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.IsTicket(tt.input)
			if got != tt.expected {
				t.Errorf("IsTicket(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Test loading with non-existent file (should return default)
	t.Run("non-existent file returns default", func(t *testing.T) {
		// Temporarily set XDG_CONFIG_HOME to a non-existent path
		oldEnv := os.Getenv("XDG_CONFIG_HOME")
		defer func() { _ = os.Setenv("XDG_CONFIG_HOME", oldEnv) }()

		testDir := t.TempDir()
		_ = os.Setenv("XDG_CONFIG_HOME", filepath.Join(testDir, "nonexistent"))

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() with non-existent file should not error, got: %v", err)
		}
		if cfg == nil {
			t.Fatal("Load() should return default config")
		}
		if len(cfg.TicketPatterns) == 0 {
			t.Error("Load() should return config with patterns")
		}
	})

	// Test loading with valid config file
	t.Run("valid config file", func(t *testing.T) {
		testDir := t.TempDir()
		configDir := filepath.Join(testDir, "branch")
		configPath := filepath.Join(configDir, "config.json")

		// Create config file
		configData := `{
  "ticket_patterns": [
    "^CUSTOM-\\d+$",
    "^#\\d+$"
  ],
  "branch_commands": [
    "feat",
    "fix"
  ]
}`
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatalf("Failed to create config dir: %v", err)
		}
		if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		oldEnv := os.Getenv("XDG_CONFIG_HOME")
		defer func() { _ = os.Setenv("XDG_CONFIG_HOME", oldEnv) }()
		_ = os.Setenv("XDG_CONFIG_HOME", testDir)

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() with valid file should not error, got: %v", err)
		}
		if cfg == nil {
			t.Fatal("Load() should return config")
		}
		if len(cfg.TicketPatterns) != 2 {
			t.Errorf("Expected 2 patterns, got %d", len(cfg.TicketPatterns))
		}

		// Test that custom pattern works
		if !cfg.IsTicket("CUSTOM-123") {
			t.Error("Custom pattern should match CUSTOM-123")
		}
		if cfg.IsTicket("PIP-123") {
			t.Error("PIP-123 should not match custom patterns")
		}

		// Test that custom branch commands are loaded
		if len(cfg.BranchCommands) != 2 {
			t.Errorf("Expected 2 branch commands, got %d", len(cfg.BranchCommands))
		}
		if cfg.BranchCommands[0] != "feat" || cfg.BranchCommands[1] != "fix" {
			t.Errorf("Expected branch commands [feat, fix], got %v", cfg.BranchCommands)
		}
	})

	// Test loading with empty patterns (should merge with defaults)
	t.Run("empty patterns merges with defaults", func(t *testing.T) {
		testDir := t.TempDir()
		configDir := filepath.Join(testDir, "branch")
		configPath := filepath.Join(configDir, "config.json")

		configData := `{
  "ticket_patterns": []
}`
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatalf("Failed to create config dir: %v", err)
		}
		if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		oldEnv := os.Getenv("XDG_CONFIG_HOME")
		defer func() { _ = os.Setenv("XDG_CONFIG_HOME", oldEnv) }()
		_ = os.Setenv("XDG_CONFIG_HOME", testDir)

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() should not error, got: %v", err)
		}
		if len(cfg.TicketPatterns) == 0 {
			t.Error("Load() should merge empty patterns with defaults")
		}
	})

	// Test loading with invalid JSON
	t.Run("invalid JSON returns error", func(t *testing.T) {
		testDir := t.TempDir()
		configDir := filepath.Join(testDir, "branch")
		configPath := filepath.Join(configDir, "config.json")

		configData := `{ invalid json }`
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatalf("Failed to create config dir: %v", err)
		}
		if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		oldEnv := os.Getenv("XDG_CONFIG_HOME")
		defer func() { _ = os.Setenv("XDG_CONFIG_HOME", oldEnv) }()
		_ = os.Setenv("XDG_CONFIG_HOME", testDir)

		_, err := Load()
		if err == nil {
			t.Error("Load() with invalid JSON should return error")
		}
	})
}

func TestSave(t *testing.T) {
	testDir := t.TempDir()
	configDir := filepath.Join(testDir, "branch")
	configPath := filepath.Join(configDir, "config.json")

	oldEnv := os.Getenv("XDG_CONFIG_HOME")
	defer func() { _ = os.Setenv("XDG_CONFIG_HOME", oldEnv) }()
	_ = os.Setenv("XDG_CONFIG_HOME", testDir)

	cfg := Default()
	cfg.TicketPatterns = []string{`^TEST-\d+$`}
	cfg.BranchCommands = []string{"custom", "command"}

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() should not error, got: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Save() should create config file")
	}

	// Load and verify
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if len(loaded.TicketPatterns) != 1 {
		t.Errorf("Expected 1 pattern, got %d", len(loaded.TicketPatterns))
	}

	if loaded.TicketPatterns[0] != `^TEST-\d+$` {
		t.Errorf("Expected pattern ^TEST-\\d+$, got %q", loaded.TicketPatterns[0])
	}

	// Verify branch commands are saved and loaded
	if len(loaded.BranchCommands) != 2 {
		t.Errorf("Expected 2 branch commands, got %d", len(loaded.BranchCommands))
	}
	if loaded.BranchCommands[0] != "custom" || loaded.BranchCommands[1] != "command" {
		t.Errorf("Expected branch commands [custom, command], got %v", loaded.BranchCommands)
	}
}

func TestCompile(t *testing.T) {
	cfg := &Config{
		TicketPatterns: []string{
			`^#\d+$`,
			`invalid[regex`, // This should be skipped
			`^[A-Z]+-\d+$`,
		},
	}

	cfg.compile()

	// Should have 2 compiled patterns (invalid one skipped)
	if len(cfg.compiled) != 2 {
		t.Errorf("Expected 2 compiled patterns, got %d", len(cfg.compiled))
	}

	// Test that valid patterns work
	if !cfg.IsTicket("#123") {
		t.Error("Pattern should match #123")
	}
	if !cfg.IsTicket("PIP-123") {
		t.Error("Pattern should match PIP-123")
	}
}
