package backup

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRunBackup(t *testing.T) {
	// Create a fake home structure with files to back up
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)

	// Create some "config" files
	configDir := filepath.Join(fakeHome, ".config", "myapp")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "config.toml"), []byte("key = \"value\""), 0644)
	os.WriteFile(filepath.Join(fakeHome, ".myrc"), []byte("rc content"), 0644)

	// Paths use ~ syntax
	destPaths := []string{
		"~/.config/myapp",
		"~/.myrc",
	}

	var lines []string
	backupDir, err := Run(destPaths, func(line string) {
		lines = append(lines, line)
	})
	if err != nil {
		t.Fatal(err)
	}
	if backupDir == "" {
		t.Fatal("expected non-empty backup dir")
	}

	// Verify the backup contains the files
	backedUpRC, err := os.ReadFile(filepath.Join(backupDir, ".myrc"))
	if err != nil {
		t.Fatalf("backup file missing: %v", err)
	}
	if string(backedUpRC) != "rc content" {
		t.Errorf("backup content mismatch: got %q", string(backedUpRC))
	}

	backedUpConfig, err := os.ReadFile(filepath.Join(backupDir, ".config", "myapp", "config.toml"))
	if err != nil {
		t.Fatalf("backup dir file missing: %v", err)
	}
	if string(backedUpConfig) != "key = \"value\"" {
		t.Errorf("backup dir content mismatch: got %q", string(backedUpConfig))
	}
}

func TestRunBackupNothingToBackup(t *testing.T) {
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)

	destPaths := []string{"~/.does-not-exist", "~/.also-missing"}

	var lines []string
	backupDir, err := Run(destPaths, func(line string) {
		lines = append(lines, line)
	})
	if err != nil {
		t.Fatal(err)
	}
	if backupDir != "" {
		t.Errorf("expected empty backup dir, got %q", backupDir)
	}
}

func TestCopyFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")

	content := []byte("file content")
	os.WriteFile(src, content, 0755)

	if err := copyFile(src, dst); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("content mismatch: got %q", string(got))
	}

	// Check permissions preserved
	info, _ := os.Stat(dst)
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0755 {
		t.Errorf("permissions not preserved: got %o, want 0755", info.Mode().Perm())
	}
}

func TestCopyDir(t *testing.T) {
	tmp := t.TempDir()
	srcDir := filepath.Join(tmp, "src")
	dstDir := filepath.Join(tmp, "dst")

	// Create a directory tree
	os.MkdirAll(filepath.Join(srcDir, "sub"), 0755)
	os.WriteFile(filepath.Join(srcDir, "root.txt"), []byte("root"), 0644)
	os.WriteFile(filepath.Join(srcDir, "sub", "nested.txt"), []byte("nested"), 0644)

	if err := copyDir(srcDir, dstDir); err != nil {
		t.Fatal(err)
	}

	// Verify structure
	got, err := os.ReadFile(filepath.Join(dstDir, "root.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "root" {
		t.Errorf("root.txt content: got %q", string(got))
	}

	got, err = os.ReadFile(filepath.Join(dstDir, "sub", "nested.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "nested" {
		t.Errorf("nested.txt content: got %q", string(got))
	}
}

func TestExpandHome(t *testing.T) {
	home := "/Users/testuser"

	tests := []struct {
		input, want string
	}{
		{"~/foo", filepath.Join(home, "foo")},
		{"~/.config/bar", filepath.Join(home, ".config/bar")},
		{"/absolute/path", "/absolute/path"},
	}

	for _, tt := range tests {
		got := expandHome(tt.input, home)
		if got != tt.want {
			t.Errorf("expandHome(%q, %q) = %q, want %q", tt.input, home, got, tt.want)
		}
	}
}
