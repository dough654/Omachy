package tui

import tea "github.com/charmbracelet/bubbletea"

// PhaseStarted signals a phase has begun.
type PhaseStarted struct {
	Name string
}

// PhaseCompleted signals a phase finished successfully.
type PhaseCompleted struct {
	Name string
}

// PhaseFailed signals a phase failed.
type PhaseFailed struct {
	Name  string
	Error error
}

// LogLine is a line of output to append to the viewport.
type LogLine struct {
	Text string
}

// ProgressUpdate updates the overall progress percentage.
type ProgressUpdate struct {
	Percent int
}

// InstallFinished signals the installer goroutine is done.
type InstallFinished struct {
	Err error
}

// WaitForUser asks the user to confirm before continuing.
// The installer goroutine blocks on Done until the TUI signals it.
type WaitForUser struct {
	Prompt string
	Done   chan struct{}
}

// All event types implement tea.Msg.
var (
	_ tea.Msg = PhaseStarted{}
	_ tea.Msg = PhaseCompleted{}
	_ tea.Msg = PhaseFailed{}
	_ tea.Msg = LogLine{}
	_ tea.Msg = ProgressUpdate{}
	_ tea.Msg = InstallFinished{}
	_ tea.Msg = WaitForUser{}
)
