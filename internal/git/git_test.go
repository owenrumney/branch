package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCreateBranch(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Initialize a git repository
	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}
	defer func() {
		// Try to return to original directory, but don't fail if it doesn't work
		_ = os.Chdir("/")
	}()

	// Initialize git repo
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Fatalf("Failed to initialize git repo: %v", err)
	}

	// Configure git user (required for commits)
	_ = exec.Command("git", "config", "user.email", "test@example.com").Run()
	_ = exec.Command("git", "config", "user.name", "Test User").Run()

	// Create an initial commit (needed for checkout)
	readmePath := filepath.Join(testDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	_ = exec.Command("git", "add", "README.md").Run()
	_ = exec.Command("git", "commit", "-m", "Initial commit").Run()

	t.Run("create new branch successfully", func(t *testing.T) {
		branchName := "feat/test-branch"
		err := CreateBranch(branchName)
		if err != nil {
			t.Fatalf("CreateBranch() should succeed, got error: %v", err)
		}

		// Verify branch was created
		out, err := exec.Command("git", "branch", "--list", branchName).Output()
		if err != nil {
			t.Fatalf("Failed to check branch: %v", err)
		}
		if string(out) == "" {
			t.Error("Branch should exist after creation")
		}

		// Clean up: switch back to main/master
		_ = exec.Command("git", "checkout", "-").Run()
	})

	t.Run("error when branch already exists", func(t *testing.T) {
		branchName := "fix/existing-branch"

		// Create branch first time
		if err := CreateBranch(branchName); err != nil {
			t.Fatalf("First CreateBranch() should succeed: %v", err)
		}
		_ = exec.Command("git", "checkout", "-").Run()

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
