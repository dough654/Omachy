package main

import (
	"fmt"
	"os"

	"github.com/dough654/Omachy/cmd"
	"github.com/dough654/Omachy/internal/installer"
)

// version is set via ldflags at build time.
var version = "dev"

func main() {
	cmd.Version = version
	installer.EmbeddedConfigs = Configs

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
