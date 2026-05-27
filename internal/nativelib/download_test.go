package nativelib

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractLibrary(t *testing.T) {
	zipPath := createTestZip(t, "test.dll", []byte("fake dll content"))
	destDir := t.TempDir()

	path, err := ExtractLibrary(zipPath, destDir, "test.dll")
	if err != nil {
		t.Fatalf("ExtractLibrary() error: %v", err)
	}

	if filepath.Base(path) != "test.dll" {
		t.Errorf("extracted path = %q, want test.dll basename", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if string(data) != "fake dll content" {
		t.Errorf("content = %q, want %q", string(data), "fake dll content")
	}
}

func TestExtractLibraryNestedPath(t *testing.T) {
	zipPath := createTestZipNested(t, "lib/wgpu_native.dll", []byte("nested content"))
	destDir := t.TempDir()

	path, err := ExtractLibrary(zipPath, destDir, "wgpu_native.dll")
	if err != nil {
		t.Fatalf("ExtractLibrary() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != "nested content" {
		t.Errorf("content = %q, want %q", string(data), "nested content")
	}
}

func TestExtractLibraryNotFound(t *testing.T) {
	zipPath := createTestZip(t, "other.dll", []byte("data"))
	destDir := t.TempDir()

	_, err := ExtractLibrary(zipPath, destDir, "missing.dll")
	if err == nil {
		t.Fatal("expected error for missing file in archive")
	}
}

func TestExtractLibraryInvalidZip(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "bad.zip")
	os.WriteFile(tmpFile, []byte("not a zip"), 0o644)

	_, err := ExtractLibrary(tmpFile, t.TempDir(), "test.dll")
	if err == nil {
		t.Fatal("expected error for invalid zip")
	}
}

func createTestZip(t *testing.T, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.zip")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	w := zip.NewWriter(f)
	entry, err := w.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	entry.Write(content)
	w.Close()
	return path
}

func createTestZipNested(t *testing.T, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.zip")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	w := zip.NewWriter(f)
	entry, err := w.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	entry.Write(content)
	w.Close()
	return path
}
