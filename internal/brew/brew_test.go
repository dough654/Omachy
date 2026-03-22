package brew

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFakeBrew creates a fake brew script in a temp directory and prepends it
// to PATH so that exec.Command("brew", ...) finds it instead of the real brew.
// The script body should be a valid sh script (without the shebang).
func writeFakeBrew(t *testing.T, scriptBody string) {
	t.Helper()
	dir := t.TempDir()
	script := filepath.Join(dir, "brew")
	err := os.WriteFile(script, []byte("#!/bin/sh\n"+scriptBody), 0755)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+":"+os.Getenv("PATH"))
}

func TestInstall_CaskPassesAdopt(t *testing.T) {
	// Verify that cask installs include --adopt so pre-existing apps
	// are claimed by Homebrew instead of causing "already an App at" errors.
	writeFakeBrew(t, `
case "$1" in
  list)
    exit 1
    ;;
  install)
    # Log all args so the test can inspect them
    echo "args: $*"
    exit 0
    ;;
esac
`)

	var lines []string
	err := Install("ghostty", true, func(line string) {
		lines = append(lines, line)
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := strings.Join(lines, "\n")
	if !strings.Contains(output, "--adopt") {
		t.Errorf("expected --adopt in brew args, got:\n%s", output)
	}
	if !strings.Contains(output, "--cask") {
		t.Errorf("expected --cask in brew args, got:\n%s", output)
	}
}

func TestInstall_FormulaDoesNotPassAdopt(t *testing.T) {
	// --adopt is cask-only; formulae should not include it.
	writeFakeBrew(t, `
case "$1" in
  list)
    exit 1
    ;;
  install)
    echo "args: $*"
    exit 0
    ;;
esac
`)

	var lines []string
	err := Install("neovim", false, func(line string) {
		lines = append(lines, line)
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := strings.Join(lines, "\n")
	if strings.Contains(output, "--adopt") {
		t.Errorf("formula install should not include --adopt, got:\n%s", output)
	}
	if strings.Contains(output, "--cask") {
		t.Errorf("formula install should not include --cask, got:\n%s", output)
	}
}

func TestInstall_CaskAlreadyInstalledViaBrew(t *testing.T) {
	// Simulate: brew list --cask ghostty → exit 0 (already installed)
	writeFakeBrew(t, `
case "$1" in
  list)
    echo "ghostty"
    exit 0
    ;;
esac
`)

	var lines []string
	err := Install("ghostty", true, func(line string) {
		lines = append(lines, line)
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := strings.Join(lines, "\n")
	if !strings.Contains(output, "Already installed: ghostty") {
		t.Errorf("expected 'Already installed' message, got:\n%s", output)
	}
}

func TestInstall_CaskRealFailure(t *testing.T) {
	// Simulate: brew list --cask badpkg → exit 1
	//           brew install --cask --adopt badpkg → generic error, exit 1
	writeFakeBrew(t, `
case "$1" in
  list)
    exit 1
    ;;
  install)
    echo "Error: No available cask for badpkg" >&2
    exit 1
    ;;
esac
`)

	err := Install("badpkg", true, func(string) {})

	if err == nil {
		t.Fatal("expected error for genuine install failure, got nil")
	}
}

func TestInstall_FormulaSuccess(t *testing.T) {
	// Simulate: brew list neovim → exit 1, brew install neovim → exit 0
	writeFakeBrew(t, `
case "$1" in
  list)
    exit 1
    ;;
  install)
    echo "Installing neovim..."
    exit 0
    ;;
esac
`)

	var lines []string
	err := Install("neovim", false, func(line string) {
		lines = append(lines, line)
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := strings.Join(lines, "\n")
	if !strings.Contains(output, "==> Installing neovim") {
		t.Errorf("expected install message, got:\n%s", output)
	}
}
