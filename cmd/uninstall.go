package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dough654/Omachy/internal/tui"
	"github.com/dough654/Omachy/internal/uninstaller"
	"github.com/spf13/cobra"
)

var (
	uninstallDryRun      bool
	uninstallKeepConfigs bool
	uninstallKeepPkgs    bool
)

func init() {
	uninstallCmd.Flags().BoolVar(&uninstallDryRun, "dry-run", false, "Show what would be done without making changes")
	uninstallCmd.Flags().BoolVar(&uninstallKeepConfigs, "keep-configs", false, "Keep deployed config files")
	uninstallCmd.Flags().BoolVar(&uninstallKeepPkgs, "keep-packages", false, "Keep installed packages")
	rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove Omachy and restore original settings",
	Long:  "Stops services, removes configs, uninstalls packages, and restores macOS defaults from the saved state.",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := uninstaller.Options{
			DryRun:       uninstallDryRun,
			KeepConfigs:  uninstallKeepConfigs,
			KeepPackages: uninstallKeepPkgs,
		}

		splashOpts := tui.SplashOptions{
			DryRun: uninstallDryRun,
		}

		return tui.Run(uninstaller.PhaseNames(), func(p *tea.Program) {
			uninstaller.Run(p, opts)
		}, splashOpts, Version)
	},
}
