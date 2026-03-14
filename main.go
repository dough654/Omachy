package main

import (
	"fmt"
	"os"

	"github.com/dough654/Omachy/cmd"
	"github.com/dough654/Omachy/internal/installer"
)

func main() {
	installer.EmbeddedConfigs = Configs

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
