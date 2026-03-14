package tui

import "fmt"

// HelpModel renders the bottom keybinding hints.
type HelpModel struct {
	Width    int
	Finished bool
}

func NewHelpModel(width int) HelpModel {
	return HelpModel{Width: width}
}

func (m HelpModel) View() string {
	bindings := []struct{ key, desc string }{
		{"q", "quit"},
		{"↑↓", "scroll"},
	}

	if m.Finished {
		bindings = []struct{ key, desc string }{
			{"q", "quit"},
			{"enter", "exit"},
		}
	}

	var s string
	for i, b := range bindings {
		if i > 0 {
			s += "  "
		}
		s += fmt.Sprintf("%s %s",
			helpKeyStyle.Render("["+b.key+"]"),
			helpStyle.Render(b.desc),
		)
	}
	return " " + s
}
