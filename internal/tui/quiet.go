package tui

import (
	"fmt"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// RunQuiet runs the installer/uninstaller without the Bubbletea TUI.
// It prints log lines to stdout and blocks until the work is done.
func RunQuiet(installer InstallerFunc) (RunResult, error) {
	app := &quietApp{}

	// Provide a dummy input and discard output so Bubbletea never tries
	// to open /dev/tty — this allows running over SSH without a PTY.
	p := tea.NewProgram(app,
		tea.WithInput(strings.NewReader("")),
		tea.WithOutput(io.Discard),
	)
	app.program = p

	go func() {
		installer(p)
	}()

	if _, err := p.Run(); err != nil {
		return RunResult{}, err
	}

	return RunResult{Finished: true, Err: app.lastErr}, nil
}

// quietApp is a minimal Bubbletea model that prints messages to stdout
// instead of rendering a TUI.
type quietApp struct {
	program *tea.Program
	lastErr error
}

func (q *quietApp) Init() tea.Cmd { return nil }

func (q *quietApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case PhaseStarted:
		fmt.Printf("[%s] starting...\n", msg.Name)
	case PhaseCompleted:
		fmt.Printf("[%s] done\n", msg.Name)
	case PhaseFailed:
		fmt.Printf("[%s] FAILED: %v\n", msg.Name, msg.Error)
	case LogLine:
		fmt.Println(msg.Text)
	case InstallFinished:
		if msg.Err != nil {
			fmt.Printf("\nFailed: %v\n", msg.Err)
		} else {
			fmt.Println("\nComplete!")
		}
		q.lastErr = msg.Err
		return q, tea.Quit
	case ProgressUpdate:
		// skip in quiet mode
	}
	return q, nil
}

func (q *quietApp) View() string { return "" }
