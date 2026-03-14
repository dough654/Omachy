package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// HeaderModel renders the top bar with title and progress.
type HeaderModel struct {
	Title   string
	Percent int
	Width   int
}

func NewHeaderModel(title string, width int) HeaderModel {
	return HeaderModel{
		Title: title,
		Width: width,
	}
}

func (m HeaderModel) View() string {
	left := titleStyle.Render(m.Title)
	right := progressStyle.Render(fmt.Sprintf("%d%%", m.Percent))

	gap := m.Width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}

	line := lipgloss.NewStyle().
		Foreground(colorMuted).
		Render(repeatChar('─', gap))

	return fmt.Sprintf(" %s %s %s ", left, line, right)
}

func repeatChar(c rune, n int) string {
	if n < 0 {
		n = 0
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = c
	}
	return string(b)
}
