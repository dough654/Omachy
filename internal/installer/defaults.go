package installer

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dough654/Omachy/internal/brew"
	"github.com/dough654/Omachy/internal/manifest"
	"github.com/dough654/Omachy/internal/shell"
	"github.com/dough654/Omachy/internal/tui"
)

// macOSDefault represents a defaults write operation.
type macOSDefault struct {
	Domain string
	Key    string
	Type   string // -bool, -int, -float, -string
	Value  string
	Label  string // human-readable description
}

var macOSDefaults = []macOSDefault{
	{"com.apple.dock", "autohide", "-bool", "true", "Auto-hide Dock"},
	{"com.apple.dock", "autohide-delay", "-float", "0", "Remove Dock auto-hide delay"},
	{"com.apple.dock", "autohide-time-modifier", "-float", "0.25", "Fast Dock hide animation"},
	{"com.apple.dock", "mru-spaces", "-bool", "false", "Disable MRU Spaces reordering"},
	{"com.apple.dock", "tilesize", "-int", "48", "Set Dock icon size"},
	{"com.apple.dock", "mineffect", "-string", "scale", "Use scale minimize effect"},
	{"com.apple.dock", "show-recents", "-bool", "false", "Hide recent apps in Dock"},
	{"NSGlobalDomain", "NSAutomaticWindowAnimationsEnabled", "-bool", "false", "Disable window open/close animations"},
	{"NSGlobalDomain", "AppleShowAllExtensions", "-bool", "true", "Show all file extensions"},
	{"NSGlobalDomain", "KeyRepeat", "-int", "2", "Fast key repeat rate"},
	{"NSGlobalDomain", "InitialKeyRepeat", "-int", "15", "Short key repeat delay"},
	{"-g", "ApplePressAndHoldEnabled", "-bool", "false", "Disable press-and-hold for key repeat"},
}

func runSystem(p *tea.Program, opts Options) error {
	log := func(text string) {
		p.Send(tui.LogLine{Text: text})
	}

	state, err := LoadState()
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	// Apply macOS defaults
	log("==> Applying macOS defaults")
	for _, d := range macOSDefaults {
		if opts.DryRun {
			log(fmt.Sprintf("    Would set: %s", d.Label))
			continue
		}

		// Read current value first for undo
		stateKey := fmt.Sprintf("%s:%s", d.Domain, d.Key)
		if _, exists := state.OriginalDefaults[stateKey]; !exists {
			current, err := readDefault(d.Domain, d.Key)
			if err == nil {
				state.OriginalDefaults[stateKey] = current
			}
		}

		// Write new value
		if err := writeDefault(d.Domain, d.Key, d.Type, d.Value); err != nil {
			return fmt.Errorf("defaults write %s %s: %w", d.Domain, d.Key, err)
		}
		log(fmt.Sprintf("    %s", d.Label))
	}

	// Restart Dock to apply changes
	if !opts.DryRun {
		log("==> Restarting Dock")
		shell.Run("killall", "Dock")
	}

	p.Send(tui.ProgressUpdate{Percent: 90})

	// Start brew services
	log("==> Starting brew services")
	for _, svc := range manifest.Services() {
		if opts.DryRun {
			log(fmt.Sprintf("    Would start service: %s", svc.Name))
			continue
		}
		if err := brew.StartService(svc.Name, log); err != nil {
			log(fmt.Sprintf("    Warning: failed to start %s: %v", svc.Name, err))
		}
		state.Services = appendUnique(state.Services, svc.Name)
	}

	// Prompt for AeroSpace accessibility permissions
	log("==> AeroSpace requires Accessibility permissions")
	log("    Opening System Settings...")
	if !opts.DryRun {
		shell.Run("open", "x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility")
		log("    Please grant AeroSpace access, then restart it.")
	}

	// Save state
	if !opts.DryRun {
		if err := SaveState(state); err != nil {
			return fmt.Errorf("save state: %w", err)
		}
	}

	return nil
}

func readDefault(domain, key string) (string, error) {
	result, err := shell.Run("defaults", "read", domain, key)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Stdout), nil
}

func writeDefault(domain, key, typ, value string) error {
	_, err := shell.Run("defaults", "write", domain, key, typ, value)
	return err
}

func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}
