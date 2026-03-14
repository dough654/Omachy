package cmd

import (
	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "omachy",
	Short: "Tiling window manager setup for macOS",
	Long:  "Omachy automates setting up an opinionated, tiling-WM-driven desktop experience on macOS.",
}

func Execute() error {
	return rootCmd.Execute()
}
