package tui

import "testing"

func TestSetStatus(t *testing.T) {
	m := NewPhasesModel([]string{"Preflight", "Packages", "Configs"})

	m.SetStatus("Packages", StatusActive)
	if m.Phases[1].Status != StatusActive {
		t.Errorf("expected StatusActive, got %d", m.Phases[1].Status)
	}

	m.SetStatus("Packages", StatusDone)
	if m.Phases[1].Status != StatusDone {
		t.Errorf("expected StatusDone, got %d", m.Phases[1].Status)
	}

	// Other phases should remain pending
	if m.Phases[0].Status != StatusPending {
		t.Errorf("Preflight should still be pending, got %d", m.Phases[0].Status)
	}
}

func TestSetStatusUnknown(t *testing.T) {
	m := NewPhasesModel([]string{"Preflight"})
	// Should not panic
	m.SetStatus("NonExistent", StatusDone)

	// Original phase unchanged
	if m.Phases[0].Status != StatusPending {
		t.Errorf("expected StatusPending, got %d", m.Phases[0].Status)
	}
}

func TestPhasesView(t *testing.T) {
	m := NewPhasesModel([]string{"Preflight", "Packages", "Configs"})
	m.SetStatus("Preflight", StatusDone)
	m.SetStatus("Packages", StatusActive)

	view := m.View(10)

	// Should contain phase names
	for _, name := range []string{"Preflight", "Packages", "Configs"} {
		if !containsStr(view, name) {
			t.Errorf("view should contain %q", name)
		}
	}

	// Should contain status icons
	if !containsStr(view, "✓") {
		t.Error("view should contain done icon ✓")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && searchStr(s, sub)
}

func searchStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
