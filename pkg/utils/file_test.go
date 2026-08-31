package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	directory := t.TempDir()
	file := filepath.Join(directory, "present.txt")
	if err := os.WriteFile(file, []byte("present"), 0o600); err != nil {
		t.Fatal(err)
	}

	if !FileExists(file) {
		t.Fatal("FileExists returned false for an existing file")
	}
	if !FileExists(directory) {
		t.Fatal("FileExists returned false for an existing directory")
	}
	if FileExists(filepath.Join(directory, "missing.txt")) {
		t.Fatal("FileExists returned true for a missing path")
	}
}
