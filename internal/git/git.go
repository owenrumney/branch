package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func CreateBranch(name string) error {
	// Check if we're in a git repository
	if err := exec.Command("git", "rev-parse", "--git-dir").Run(); err != nil {
		return fmt.Errorf("not a git repository")
	}

	// Check if branch already exists
	out, _ := exec.Command("git", "branch", "--list", name).Output()
	if strings.TrimSpace(string(out)) != "" {
		return fmt.Errorf("branch %q already exists", name)
	}

	// Create and switch to the new branch
	cmd := exec.Command("git", "checkout", "-b", name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(output)))
	}

	return nil
}
