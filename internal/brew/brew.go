package brew

import (
	"fmt"
	"strings"

	"github.com/dough654/Omachy/internal/shell"
)

// Tap adds a Homebrew tap if not already tapped.
func Tap(name string, onLine func(string)) error {
	if IsTapped(name) {
		onLine(fmt.Sprintf("    Already tapped: %s", name))
		return nil
	}

	onLine(fmt.Sprintf("==> Tapping %s", name))
	return shell.RunStreaming("brew", []string{"tap", name}, onLine)
}

// IsTapped checks if a tap is already added.
func IsTapped(name string) bool {
	result, err := shell.Run("brew", "tap")
	if err != nil {
		return false
	}
	for _, line := range strings.Split(result.Stdout, "\n") {
		if strings.TrimSpace(line) == name {
			return true
		}
	}
	return false
}

// Install installs a formula or cask if not already installed.
func Install(name string, cask bool, onLine func(string)) error {
	if IsInstalled(name, cask) {
		onLine(fmt.Sprintf("    Already installed: %s", name))
		return nil
	}

	args := []string{"install"}
	if cask {
		// --adopt tells Homebrew to claim an existing .app that was installed
		// outside of Homebrew (e.g. via DMG) instead of failing with
		// "It seems there is already an App at ...".
		args = append(args, "--cask", "--adopt")
	}
	args = append(args, name)

	onLine(fmt.Sprintf("==> Installing %s", name))
	if err := shell.RunStreaming("brew", args, onLine); err != nil {
		// brew install exits 1 when the package installs but linking fails
		// (e.g. another formula owns the symlinks). The package is still
		// usable in the Cellar, so treat this as a warning, not a failure.
		if IsInstalled(name, cask) {
			onLine(fmt.Sprintf("    Warning: %s installed but not linked (conflicting package)", name))
			return nil
		}
		return err
	}
	return nil
}

// IsInstalled checks if a formula or cask is installed.
func IsInstalled(name string, cask bool) bool {
	args := []string{"list"}
	if cask {
		args = append(args, "--cask")
	}
	args = append(args, name)

	_, err := shell.Run("brew", args...)
	return err == nil
}

// StartService starts a brew service.
func StartService(name string, onLine func(string)) error {
	onLine(fmt.Sprintf("==> Starting service: %s", name))
	return shell.RunStreaming("brew", []string{"services", "start", name}, onLine)
}

// StopService stops a brew service.
func StopService(name string, onLine func(string)) error {
	onLine(fmt.Sprintf("==> Stopping service: %s", name))
	return shell.RunStreaming("brew", []string{"services", "stop", name}, onLine)
}

// Uninstall removes a formula or cask.
func Uninstall(name string, cask bool, onLine func(string)) error {
	if !IsInstalled(name, cask) {
		onLine(fmt.Sprintf("    Not installed: %s", name))
		return nil
	}

	args := []string{"uninstall"}
	if cask {
		args = append(args, "--cask")
	}
	args = append(args, name)

	onLine(fmt.Sprintf("==> Uninstalling %s", name))
	return shell.RunStreaming("brew", args, onLine)
}
