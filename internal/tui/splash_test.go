package tui

import (
	"strings"
	"testing"
)

func TestRenderSplash(t *testing.T) {
	view := renderSplash(80, 40, SplashOptions{}, "0.1.0")

	// Should contain logo (ASCII art spells it out)
	if !strings.Contains(view, "___") {
		t.Error("splash should contain logo ASCII art")
	}

	// Should contain tool names
	tools := []string{"AeroSpace", "SketchyBar", "Ghostty", "Neovim", "Tmux"}
	for _, tool := range tools {
		if !strings.Contains(view, tool) {
			t.Errorf("splash should contain tool name %q", tool)
		}
	}

	// Should contain prompt
	if !strings.Contains(view, "Enter") {
		t.Error("splash should contain Enter prompt")
	}

	// Should contain version
	if !strings.Contains(view, "0.1.0") {
		t.Error("splash should contain version")
	}
}

func TestRenderSplashWithFlags(t *testing.T) {
	view := renderSplash(80, 40, SplashOptions{
		DryRun:     true,
		Force:      true,
		SkipBackup: true,
	}, "0.1.0")

	if !strings.Contains(view, "--dry-run") {
		t.Error("splash should show --dry-run flag")
	}
	if !strings.Contains(view, "--force") {
		t.Error("splash should show --force flag")
	}
	if !strings.Contains(view, "--skip-backup") {
		t.Error("splash should show --skip-backup flag")
	}
}
