package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/dough654/Omachy/internal/preflight"
	"github.com/spf13/cobra"
)

var (
	iconPass = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")).Render("✓")
	iconFail = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Render("✗")
	iconWarn = lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B")).Render("!")

	doctorTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			MarginBottom(1)

	doctorName = lipgloss.NewStyle().Width(20)

	doctorDetail = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF"))
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system readiness for Omachy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(doctorTitle.Render("Omachy Doctor"))
		fmt.Println()

		checks := preflight.RunAll()
		allOk := true

		for _, c := range checks {
			var icon string
			switch {
			case c.Passed:
				icon = iconPass
			case c.Warning:
				icon = iconWarn
			default:
				icon = iconFail
				allOk = false
			}

			fmt.Printf("  %s %s %s\n", icon, doctorName.Render(c.Name), doctorDetail.Render(c.Detail))
		}

		fmt.Println()
		if allOk {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")).Bold(true).Render("  System is ready for Omachy!"))
		} else {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Bold(true).Render("  Some checks failed. Fix the issues above before installing."))
		}
	},
}
