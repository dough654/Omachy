package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dough654/Omachy/internal/installer"
	"github.com/dough654/Omachy/internal/tui"
	"github.com/spf13/cobra"
)

var (
	flagDryRun     bool
	flagForce      bool
	flagVerbose    bool
	flagSkipBackup bool
)

func init() {
	installCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Show what would be done without making changes")
	installCmd.Flags().BoolVar(&flagForce, "force", false, "Overwrite existing configs without prompting")
	installCmd.Flags().BoolVar(&flagVerbose, "verbose", false, "Show detailed output")
	installCmd.Flags().BoolVar(&flagSkipBackup, "skip-backup", false, "Skip backing up existing configs")
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the full Omachy desktop environment",
	Long:  "Runs preflight checks, backs up existing configs, installs packages, deploys configs, and configures macOS system defaults.",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := installer.Options{
			DryRun:     flagDryRun,
			Force:      flagForce,
			Verbose:    flagVerbose,
			SkipBackup: flagSkipBackup,
		}

		splashOpts := tui.SplashOptions{
			DryRun:     flagDryRun,
			Force:      flagForce,
			SkipBackup: flagSkipBackup,
		}

		return tui.Run(installer.PhaseNames(), func(p *tea.Program) {
			installer.Run(p, opts)
		}, splashOpts, Version)
	},
}
