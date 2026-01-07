package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/owenrumney/branch/cmd"
)

func main() {
	// ensure that git is available on the PATH
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("git is not available on the PATH")
		os.Exit(1)
	}

	cmd.Execute()
}
