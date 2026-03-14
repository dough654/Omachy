package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PhaseStatus represents the current state of a phase.
type PhaseStatus int

const (
	StatusPending PhaseStatus = iota
	StatusActive
	StatusDone
	StatusFailed
)

// Phase holds the display state for a single install phase.
type Phase struct {
	Name   string
	Status PhaseStatus
}

// PhasesModel renders the left panel phase list.
type PhasesModel struct {
	Phases  []Phase
	spinner spinner.Model
}

func NewPhasesModel(names []string) PhasesModel {
	phases := make([]Phase, len(names))
	for i, n := range names {
		phases[i] = Phase{Name: n, Status: StatusPending}
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorPrimary)

	return PhasesModel{
		Phases:  phases,
		spinner: s,
	}
}

func (m PhasesModel) Update(msg tea.Msg) (PhasesModel, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m PhasesModel) View(height int) string {
	var b strings.Builder

	for _, p := range m.Phases {
		var icon, name string
		switch p.Status {
		case StatusPending:
			icon = phaseIconPending
			name = phaseNamePending.Render(p.Name)
		case StatusActive:
			icon = m.spinner.View()
			name = phaseNameActive.Render(p.Name)
		case StatusDone:
			icon = phaseIconDone
			name = phaseNameDone.Render(p.Name)
		case StatusFailed:
			icon = phaseIconFailed
			name = phaseNameFailed.Render(p.Name)
		}
		fmt.Fprintf(&b, " %s %s\n", icon, name)
	}

	// Pad to fill height
	lines := len(m.Phases)
	for i := lines; i < height; i++ {
		b.WriteString("\n")
	}

	return b.String()
}

func (m *PhasesModel) SetStatus(name string, status PhaseStatus) {
	for i := range m.Phases {
		if m.Phases[i].Name == name {
			m.Phases[i].Status = status
			return
		}
	}
}
