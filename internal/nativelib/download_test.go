package nativelib

import (
	"archive/zip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fake content")) //nolint:errcheck // test response
	}))
	defer ts.Close()

	path, err := Download(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "fake content" {
		t.Errorf("content = %q, want %q", string(data), "fake content")
	}
}

func TestDownloadHTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	_, err := Download(ts.URL)
	if err == nil {
		t.Fatal("expected error for HTTP 404")
	}
}

func TestDownloadNetworkError(t *testing.T) {
	// Use a server that immediately closes the connection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Error("ResponseWriter does not support Hijacker")
			return
		}
		conn, _, _ := hj.Hijack()
		conn.Close()
	}))
	defer ts.Close()

	_, err := Download(ts.URL)
	if err == nil {
		t.Fatal("expected error for connection drop")
	}
}

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
	os.WriteFile(tmpFile, []byte("not a zip"), 0o644) //nolint:errcheck // test setup

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
	defer f.Close() //nolint:errcheck // test helper

	w := zip.NewWriter(f)
	entry, err := w.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	entry.Write(content) //nolint:errcheck // test helper
	w.Close()            //nolint:errcheck // test helper
	return path
}

func createTestZipNested(t *testing.T, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.zip")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close() //nolint:errcheck // test helper

	w := zip.NewWriter(f)
	entry, err := w.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	entry.Write(content) //nolint:errcheck // test helper
	w.Close()            //nolint:errcheck // test helper
	return path
}
