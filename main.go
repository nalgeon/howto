// Howto is a humble command-line assistant.
// The user describes a task, and Howto suggests a command to solve it.
package main

import (
	"fmt"
	"os"

	"github.com/nalgeon/howto/internal"
	"github.com/nalgeon/howto/internal/ai"
)

var (
	version = "dev"
	commit  = "head"
	date    = "now"
)

func main() {
	history, err := internal.LoadHistory()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		internal.PrintUsage(os.Stdout)
		os.Exit(1)
	}

	ver := internal.NewVersion(version, commit, date)
	err = internal.Howto(os.Stdout, ai.Ask, ver, os.Args[1:], history)

	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
