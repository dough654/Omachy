package tui

import (
	"strings"
	"testing"
)

func TestHelpDefault(t *testing.T) {
	h := NewHelpModel(80)
	view := h.View()

	if !strings.Contains(view, "quit") {
		t.Error("default help should contain 'quit'")
	}
	if !strings.Contains(view, "scroll") {
		t.Error("default help should contain 'scroll'")
	}
}

func TestHelpFinished(t *testing.T) {
	h := NewHelpModel(80)
	h.Finished = true
	view := h.View()

	if !strings.Contains(view, "quit") {
		t.Error("finished help should contain 'quit'")
	}
	if !strings.Contains(view, "exit") {
		t.Error("finished help should contain 'exit'")
	}
	if strings.Contains(view, "scroll") {
		t.Error("finished help should not contain 'scroll'")
	}
}
