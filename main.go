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

	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("branch version %s (commit %s built on %s)\n", version, commit, date)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "generate-config" {
		if err := cfg.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
		}
		return
	}

	rootCmd := cmd.NewRootCmd(cfg, fmt.Sprintf("%s (%s) built on %s", version, commit, date))
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
