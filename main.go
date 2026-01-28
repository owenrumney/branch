package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/owenrumney/branch/cmd"
	"github.com/owenrumney/branch/internal/config"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// ensure that git is available on the PATH
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("git is not available on the PATH")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	rootCmd := cmd.NewRootCmd(cfg, fmt.Sprintf("%s (%s) built on %s", version, commit, date))
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
