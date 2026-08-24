package main

import (
	"os"

	"kck/run"
)

func main() {
	if err := run.Main(os.Args[1:]); err != nil {
		// run.MainWith already wrote user-facing Error: lines to stderr.
		os.Exit(1)
	}
}
