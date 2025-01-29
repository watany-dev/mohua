package main

import (
	"fmt"
	"mohua/cmd"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := cmd.Execute(); err != nil {
		// NoResourcesError because the message is already displayed,
		// skip displaying additional messages
		if _, ok := err.(*cmd.NoResourcesError); !ok {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
