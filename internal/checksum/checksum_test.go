package checksum

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPathFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.txt")
	content := []byte("hello world\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	got, err := Path(path)
	if err != nil {
		t.Fatal(err)
	}

	want := fmt.Sprintf("%x", sha256.Sum256(content))
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestPathDirectory(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "a.txt"), []byte("aaa"), 0644)
	os.WriteFile(filepath.Join(tmp, "b.txt"), []byte("bbb"), 0644)

	hash1, err := Path(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if hash1 == "" {
		t.Fatal("expected non-empty hash")
	}

	// Same content should produce same hash
	hash2, err := Path(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if hash1 != hash2 {
		t.Error("same directory content produced different hashes")
	}
}

func TestPathFileChanged(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.txt")
	os.WriteFile(path, []byte("original"), 0644)

	hash1, err := Path(path)
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile(path, []byte("modified"), 0644)

	hash2, err := Path(path)
	if err != nil {
		t.Fatal(err)
	}

	if hash1 == hash2 {
		t.Error("checksum did not change after file modification")
	}
}

func TestPathMissing(t *testing.T) {
	_, err := Path("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("expected error for non-existent path")
	}
}
