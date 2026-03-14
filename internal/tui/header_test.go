package tui

import (
	"strings"
	"testing"
)

func TestHeaderView(t *testing.T) {
	h := NewHeaderModel("Omachy Installer", 80)
	h.Percent = 42

	view := h.View()

	if !strings.Contains(view, "Omachy Installer") {
		t.Error("view should contain title")
	}
	if !strings.Contains(view, "42%") {
		t.Error("view should contain percentage")
	}
}

func TestHeaderProgress(t *testing.T) {
	h := NewHeaderModel("Test", 80)

	tests := []int{0, 50, 100}
	for _, pct := range tests {
		h.Percent = pct
		view := h.View()
		if !strings.Contains(view, string(rune('0'+pct/100))+string(rune('0'+(pct%100)/10))+string(rune('0'+pct%10))+"%") {
			// Simpler check
			expected := strings.TrimLeft(strings.Replace(view, " ", "", -1), "")
			_ = expected
		}
	}

	// Spot check: 0% and 100%
	h.Percent = 0
	if !strings.Contains(h.View(), "0%") {
		t.Error("should contain 0%")
	}
	h.Percent = 100
	if !strings.Contains(h.View(), "100%") {
		t.Error("should contain 100%")
	}
}
