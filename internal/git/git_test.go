package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateBranch(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not available, skipping test")
	}

	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Initialize git repo using Dir field for reliable directory handling
	initCmd := exec.Command("git", "init")
	initCmd.Dir = testDir
	if output, err := initCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to initialize git repo: %v\nOutput: %s", err, string(output))
	}

	// Configure git user (required for commits) - use --local to ensure it's set for this repo
	emailCmd := exec.Command("git", "config", "--local", "user.email", "test@example.com")
	emailCmd.Dir = testDir
	if output, err := emailCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to set git user email: %v\nOutput: %s", err, string(output))
	}

	nameCmd := exec.Command("git", "config", "--local", "user.name", "Test User")
	nameCmd.Dir = testDir
	if output, err := nameCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to set git user name: %v\nOutput: %s", err, string(output))
	}

	// Disable GPG signing for test commits
	gpgSignCmd := exec.Command("git", "config", "--local", "commit.gpgsign", "false")
	gpgSignCmd.Dir = testDir
	if output, err := gpgSignCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to disable GPG signing: %v\nOutput: %s", err, string(output))
	}

	// Create an initial commit (needed for checkout)
	readmePath := filepath.Join(testDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	addCmd := exec.Command("git", "add", "README.md")
	addCmd.Dir = testDir
	if output, err := addCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to add file: %v\nOutput: %s", err, string(output))
	}

	// Use --no-gpg-sign flag as an additional safeguard
	commitCmd := exec.Command("git", "commit", "--no-gpg-sign", "-m", "Initial commit")
	commitCmd.Dir = testDir
	if output, err := commitCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to create initial commit: %v\nOutput: %s", err, string(output))
	}

	t.Run("create new branch successfully", func(t *testing.T) {
		// Change to test directory so CreateBranch can find the git repo
		if err := os.Chdir(testDir); err != nil {
			t.Fatalf("Failed to change to test directory: %v", err)
		}

		branchName := "feat/test-branch"
		err := CreateBranch(branchName)
		if err != nil {
			t.Fatalf("CreateBranch() should succeed, got error: %v", err)
		}

		// Verify branch was created
		branchCmd := exec.Command("git", "branch", "--list", branchName)
		branchCmd.Dir = testDir
		out, err := branchCmd.Output()
		if err != nil {
			t.Fatalf("Failed to check branch: %v", err)
		}
		if strings.TrimSpace(string(out)) == "" {
			t.Error("Branch should exist after creation")
		}

		// Clean up: switch back to main/master
		checkoutCmd := exec.Command("git", "checkout", "-")
		checkoutCmd.Dir = testDir
		_ = checkoutCmd.Run()
	})

	t.Run("error when branch already exists", func(t *testing.T) {
		// Ensure we're in the test directory
		if err := os.Chdir(testDir); err != nil {
			t.Fatalf("Failed to change to test directory: %v", err)
		}

		branchName := "fix/existing-branch"

		// Create branch first time
		if err := CreateBranch(branchName); err != nil {
			t.Fatalf("First CreateBranch() should succeed: %v", err)
		}

		// Switch back to main/master
		checkoutCmd := exec.Command("git", "checkout", "-")
		checkoutCmd.Dir = testDir
		_ = checkoutCmd.Run()

		// Try to create again
		err := CreateBranch(branchName)
		if err == nil {
			t.Error("CreateBranch() should fail when branch already exists")
		}
		if err != nil && err.Error() != `branch "fix/existing-branch" already exists` {
			t.Errorf("CreateBranch() error message = %q, want %q", err.Error(), `branch "fix/existing-branch" already exists`)
		}
	})

	t.Run("error when not in git repository", func(t *testing.T) {
		// Change to a non-git directory
		nonGitDir := t.TempDir()
		if err := os.Chdir(nonGitDir); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}

		err := CreateBranch("feat/test")
		if err == nil {
			t.Error("CreateBranch() should fail when not in git repository")
		}
		if err != nil && err.Error() != "not a git repository" {
			t.Errorf("CreateBranch() error = %q, want %q", err.Error(), "not a git repository")
		}

		// Return to test directory
		if err := os.Chdir(testDir); err != nil {
			t.Fatalf("Failed to return to test directory: %v", err)
		}
	})
}
