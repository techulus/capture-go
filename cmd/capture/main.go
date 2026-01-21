package main

import (
	"os"

	"github.com/techulus/capture-go/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
