package backup

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Run backs up existing config files to ~/.omachy/backups/<timestamp>/.
// Returns the backup directory path, or empty string if nothing was backed up.
func Run(destPaths []string, onLine func(string)) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Check which paths actually exist
	var toBackup []string
	for _, dest := range destPaths {
		expanded := expandHome(dest, home)
		if _, err := os.Stat(expanded); err == nil {
			toBackup = append(toBackup, expanded)
		}
	}

	if len(toBackup) == 0 {
		onLine("    No existing configs to back up")
		return "", nil
	}

	// Create backup directory
	ts := time.Now().Format("20060102-150405")
	backupDir := filepath.Join(home, ".omachy", "backups", ts)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("create backup dir: %w", err)
	}
	onLine(fmt.Sprintf("==> Backup directory: %s", backupDir))

	for _, src := range toBackup {
		// Preserve relative structure from home
		rel, _ := filepath.Rel(home, src)
		dest := filepath.Join(backupDir, rel)

		info, err := os.Stat(src)
		if err != nil {
			continue
		}

		if info.IsDir() {
			if err := copyDir(src, dest); err != nil {
				return "", fmt.Errorf("backup %s: %w", src, err)
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
				return "", err
			}
			if err := copyFile(src, dest); err != nil {
				return "", fmt.Errorf("backup %s: %w", src, err)
			}
		}
		onLine(fmt.Sprintf("    Backed up %s", rel))
	}

	return backupDir, nil
}

func expandHome(path, home string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	return path
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	info, err := in.Stat()
	if err == nil {
		os.Chmod(dst, info.Mode())
	}

	return nil
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}
