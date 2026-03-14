package uninstaller

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dough654/Omachy/internal/brew"
	"github.com/dough654/Omachy/internal/installer"
	"github.com/dough654/Omachy/internal/manifest"
	"github.com/dough654/Omachy/internal/shell"
	"github.com/dough654/Omachy/internal/tui"
)

// PhaseNames returns the uninstall phase names.
func PhaseNames() []string {
	return []string{
		"Services",
		"Configs",
		"Packages",
		"Defaults",
	}
}

// Options holds uninstall configuration.
type Options struct {
	DryRun      bool
	KeepConfigs bool
	KeepPackages bool
}

// Run executes the full uninstall flow.
func Run(p *tea.Program, opts Options) {
	phases := []struct {
		name string
		fn   func(p *tea.Program, opts Options) error
	}{
		{"Services", stopServices},
		{"Configs", removeConfigs},
		{"Packages", removePackages},
		{"Defaults", restoreDefaults},
	}

	for i, phase := range phases {
		p.Send(tui.PhaseStarted{Name: phase.name})

		err := phase.fn(p, opts)
		if err != nil {
			p.Send(tui.PhaseFailed{Name: phase.name, Error: err})
			p.Send(tui.InstallFinished{Err: fmt.Errorf("phase %q failed: %w", phase.name, err)})
			return
		}

		p.Send(tui.PhaseCompleted{Name: phase.name})
		pct := ((i + 1) * 100) / len(phases)
		p.Send(tui.ProgressUpdate{Percent: pct})
	}

	p.Send(tui.InstallFinished{})
}

func stopServices(p *tea.Program, opts Options) error {
	log := func(text string) { p.Send(tui.LogLine{Text: text}) }

	state, err := installer.LoadState()
	if err != nil {
		return err
	}

	for _, svc := range state.Services {
		if opts.DryRun {
			log(fmt.Sprintf("==> Would stop service: %s", svc))
			continue
		}
		if err := brew.StopService(svc, log); err != nil {
			log(fmt.Sprintf("    Warning: %v", err))
		}
	}
	return nil
}

func removeConfigs(p *tea.Program, opts Options) error {
	log := func(text string) { p.Send(tui.LogLine{Text: text}) }

	if opts.KeepConfigs {
		log("==> Keeping configs (--keep-configs)")
		return nil
	}

	state, err := installer.LoadState()
	if err != nil {
		return err
	}

	for dest := range state.DeployedConfigs {
		if opts.DryRun {
			log(fmt.Sprintf("==> Would remove %s", shortPath(dest)))
			continue
		}

		log(fmt.Sprintf("==> Removing %s", shortPath(dest)))
		os.RemoveAll(dest)
	}

	// Restore from backup if available
	if state.BackupPath != "" {
		log(fmt.Sprintf("==> Restoring backup from %s", state.BackupPath))
		if !opts.DryRun {
			if err := restoreBackup(state.BackupPath, log); err != nil {
				log(fmt.Sprintf("    Warning: backup restore failed: %v", err))
			}
		}
	}

	return nil
}

func removePackages(p *tea.Program, opts Options) error {
	log := func(text string) { p.Send(tui.LogLine{Text: text}) }

	if opts.KeepPackages {
		log("==> Keeping packages (--keep-packages)")
		return nil
	}

	// Remove in reverse order
	pkgs := manifest.Packages()
	for i := len(pkgs) - 1; i >= 0; i-- {
		pkg := pkgs[i]
		if opts.DryRun {
			log(fmt.Sprintf("==> Would uninstall %s", pkg.Name))
			continue
		}
		if err := brew.Uninstall(pkg.Name, pkg.Cask, log); err != nil {
			log(fmt.Sprintf("    Warning: %v", err))
		}
	}
	return nil
}

func restoreDefaults(p *tea.Program, opts Options) error {
	log := func(text string) { p.Send(tui.LogLine{Text: text}) }

	state, err := installer.LoadState()
	if err != nil {
		return err
	}

	log("==> Restoring macOS defaults")
	for key, value := range state.OriginalDefaults {
		parts := strings.SplitN(key, ":", 2)
		if len(parts) != 2 {
			continue
		}
		domain, defKey := parts[0], parts[1]

		if opts.DryRun {
			log(fmt.Sprintf("    Would restore %s %s → %s", domain, defKey, value))
			continue
		}

		_, err := shell.Run("defaults", "write", domain, defKey, value)
		if err != nil {
			log(fmt.Sprintf("    Warning: could not restore %s: %v", defKey, err))
		} else {
			log(fmt.Sprintf("    Restored %s", defKey))
		}
	}

	if !opts.DryRun {
		log("==> Restarting Dock")
		shell.Run("killall", "Dock")

		// Clean up state file
		log("==> Cleaning up state file")
		home, _ := os.UserHomeDir()
		os.Remove(filepath.Join(home, ".omachy", "state.json"))
	}

	return nil
}

func restoreBackup(backupDir string, onLine func(string)) error {
	home, _ := os.UserHomeDir()

	return filepath.WalkDir(backupDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		rel, _ := filepath.Rel(backupDir, path)
		dest := filepath.Join(home, rel)

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		info, _ := d.Info()
		mode := info.Mode()

		onLine(fmt.Sprintf("    Restoring %s", rel))
		return os.WriteFile(dest, data, mode)
	})
}

func shortPath(path string) string {
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}
