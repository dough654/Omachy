package installer

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dough654/Omachy/internal/preflight"
	"github.com/dough654/Omachy/internal/tui"
)

func runPreflight(p *tea.Program, opts Options) error {
	checks := preflight.RunAll()

	for _, c := range checks {
		if c.Passed {
			p.Send(tui.LogLine{Text: fmt.Sprintf("==> %s: %s", c.Name, c.Detail)})
		} else if c.Warning {
			p.Send(tui.LogLine{Text: fmt.Sprintf("==> %s: ⚠ %s", c.Name, c.Detail)})
		} else {
			p.Send(tui.LogLine{Text: fmt.Sprintf("==> %s: ✗ %s", c.Name, c.Detail)})
		}
	}

	if !preflight.AllPassed(checks) {
		return fmt.Errorf("preflight checks failed — run 'omachy doctor' for details")
	}

	return nil
}
