package installer

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dough654/Omachy/internal/checksum"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dough654/Omachy/internal/backup"
	"github.com/dough654/Omachy/internal/manifest"
	"github.com/dough654/Omachy/internal/tui"
)

// EmbeddedConfigs is set by main to provide access to the embedded filesystem.
var EmbeddedConfigs embed.FS

func runBackup(p *tea.Program, opts Options) error {
	log := func(text string) {
		p.Send(tui.LogLine{Text: text})
	}

	if opts.SkipBackup {
		log("==> Skipping backup (--skip-backup)")
		return nil
	}

	log("==> Checking for existing configs to back up")

	// Collect destination paths from manifest
	var destPaths []string
	for _, cfg := range manifest.Configs() {
		destPaths = append(destPaths, cfg.Dest)
	}

	if opts.DryRun {
		for _, d := range destPaths {
			log(fmt.Sprintf("    Would back up %s (if exists)", d))
		}
		return nil
	}

	state, err := LoadState()
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	backupPath, err := backup.Run(destPaths, log)
	if err != nil {
		return err
	}

	if backupPath != "" {
		state.BackupPath = backupPath
		if err := SaveState(state); err != nil {
			return fmt.Errorf("save state: %w", err)
		}
	}

	return nil
}

func runConfigs(p *tea.Program, opts Options) error {
	log := func(text string) {
		p.Send(tui.LogLine{Text: text})
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	state, err := LoadState()
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	configs := manifest.Configs()
	for i, cfg := range configs {
		dest := expandHome(cfg.Dest, home)

		if opts.DryRun {
			log(fmt.Sprintf("==> Would deploy %s → %s", cfg.Source, cfg.Dest))
			continue
		}

		log(fmt.Sprintf("==> Deploying %s → %s", cfg.Source, cfg.Dest))

		if cfg.IsDir {
			if err := deployDir(cfg.Source, dest, cfg.Mode); err != nil {
				return fmt.Errorf("deploy %s: %w", cfg.Source, err)
			}
		} else {
			if err := deployFile(cfg.Source, dest, fs.FileMode(cfg.Mode)); err != nil {
				return fmt.Errorf("deploy %s: %w", cfg.Source, err)
			}
		}

		// Compute checksum for state tracking
		hash, _ := checksum.Path(dest)
		state.DeployedConfigs[dest] = hash

		pct := 60 + ((i+1)*20)/len(configs) // configs phase covers 60-80%
		p.Send(tui.ProgressUpdate{Percent: pct})
	}

	if !opts.DryRun {
		log("    Writing checksums to state file")
		if err := SaveState(state); err != nil {
			return fmt.Errorf("save state: %w", err)
		}
	}

	return nil
}

func deployFile(source, dest string, mode fs.FileMode) error {
	data, err := EmbeddedConfigs.ReadFile(filepath.Join("configs", source))
	if err != nil {
		return fmt.Errorf("read embedded %s: %w", source, err)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	return os.WriteFile(dest, data, mode)
}

func deployDir(source, dest string, mode uint32) error {
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	srcDir := filepath.Join("configs", source)
	return fs.WalkDir(EmbeddedConfigs, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip .gitkeep files
		if d.Name() == ".gitkeep" {
			return nil
		}

		rel, _ := filepath.Rel(srcDir, path)
		target := filepath.Join(dest, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		data, err := EmbeddedConfigs.ReadFile(path)
		if err != nil {
			return err
		}

		fileMode := fs.FileMode(mode)
		// Scripts get 0755
		if strings.HasSuffix(d.Name(), ".sh") || strings.HasPrefix(d.Name(), "plugin.") {
			fileMode = 0755
		}

		return os.WriteFile(target, data, fileMode)
	})
}

func expandHome(path, home string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	return path
}
