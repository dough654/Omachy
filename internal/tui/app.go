package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const phasePanelWidth = 20

// InstallerFunc is the function signature for the installer goroutine.
// It receives the tea.Program to send messages back to the TUI.
type InstallerFunc func(p *tea.Program)

// App is the root Bubbletea model.
type App struct {
	header       HeaderModel
	phases       PhasesModel
	output       OutputModel
	help         HelpModel
	width        int
	height       int
	started      bool // false = showing splash, true = installer running
	finished     bool
	err          error
	installer    InstallerFunc
	splashOpts   SplashOptions
	version      string
	program      *tea.Program // set after Run() creates the program
}

func NewApp(phaseNames []string, installer InstallerFunc, splashOpts SplashOptions, version string) App {
	return App{
		header:     NewHeaderModel("Omachy Installer", 80),
		phases:     NewPhasesModel(phaseNames),
		output:     NewOutputModel(56, 15),
		help:       NewHelpModel(80),
		installer:  installer,
		splashOpts: splashOpts,
		version:    version,
	}
}

func (a App) Init() tea.Cmd {
	return a.phases.spinner.Tick
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.layout()
		return a, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "enter":
			if !a.started {
				a.started = true
				// Launch the installer goroutine now
				return a, func() tea.Msg {
					go a.installer(a.program)
					return nil
				}
			}
			if a.finished {
				return a, tea.Quit
			}
		}

	case spinner.TickMsg:
		if a.started {
			var cmd tea.Cmd
			a.phases, cmd = a.phases.Update(msg)
			cmds = append(cmds, cmd)
		}

	case PhaseStarted:
		a.phases.SetStatus(msg.Name, StatusActive)

	case PhaseCompleted:
		a.phases.SetStatus(msg.Name, StatusDone)

	case PhaseFailed:
		a.phases.SetStatus(msg.Name, StatusFailed)

	case LogLine:
		a.output.AppendLine(msg.Text)

	case ProgressUpdate:
		a.header.Percent = msg.Percent

	case InstallFinished:
		a.finished = true
		a.help.Finished = true
		a.err = msg.Err
		if msg.Err != nil {
			a.output.AppendLine("")
			a.output.AppendLine(lipgloss.NewStyle().Foreground(colorError).Bold(true).Render("Installation failed: " + msg.Err.Error()))
		} else {
			a.header.Percent = 100
			a.output.AppendLine("")
			a.output.AppendLine(lipgloss.NewStyle().Foreground(colorSuccess).Bold(true).Render("Installation complete!"))
		}
	}

	// Pass through to viewport for scroll handling
	if a.started {
		var cmd tea.Cmd
		a.output, cmd = a.output.Update(msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	if a.width == 0 {
		return "Loading..."
	}

	// Splash screen before installation starts
	if !a.started {
		return renderSplash(a.width, a.height, a.splashOpts, a.version)
	}

	header := a.header.View()

	// Phase panel
	contentHeight := a.height - 4 // header + help + borders
	if contentHeight < 3 {
		contentHeight = 3
	}

	phaseContent := a.phases.View(contentHeight - 2) // panel border
	phasePanel := panelStyle.
		Width(phasePanelWidth).
		Height(contentHeight).
		Render(panelTitleStyle.Render("Phases") + "\n" + phaseContent)

	// Output panel
	outputWidth := a.width - phasePanelWidth - 5
	if outputWidth < 10 {
		outputWidth = 10
	}
	outputPanel := panelStyle.
		Width(outputWidth).
		Height(contentHeight).
		Render(panelTitleStyle.Render("Output") + "\n" + a.output.View())

	body := lipgloss.JoinHorizontal(lipgloss.Top, phasePanel, outputPanel)
	help := a.help.View()

	return lipgloss.JoinVertical(lipgloss.Left, header, body, help)
}

func (a *App) layout() {
	a.header.Width = a.width
	a.help.Width = a.width

	contentHeight := a.height - 4
	if contentHeight < 3 {
		contentHeight = 3
	}
	outputWidth := a.width - phasePanelWidth - 7
	if outputWidth < 10 {
		outputWidth = 10
	}
	outputHeight := contentHeight - 3
	if outputHeight < 1 {
		outputHeight = 1
	}
	a.output.SetSize(outputWidth, outputHeight)
}

// Run starts the Bubbletea program with the given installer function.
func Run(phaseNames []string, installer InstallerFunc, splashOpts SplashOptions, version string) error {
	app := NewApp(phaseNames, installer, splashOpts, version)
	p := tea.NewProgram(&app, tea.WithAltScreen())
	app.program = p

	_, err := p.Run()
	return err
}
