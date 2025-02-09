// Howto is a humble command-line assistant.
// The user describes a task, and Howto suggests a command to solve it.
package main

import (
	"fmt"
	"os"

	"github.com/nalgeon/howto/ai"
)

var version string = "latest"

func main() {
	history, err := LoadHistory()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		printUsage(os.Stdout)
		os.Exit(1)
	}

	err = howto(os.Stdout, ai.Ask, os.Args[1:], history)

	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
