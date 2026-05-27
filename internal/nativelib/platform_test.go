package nativelib

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDetectPlatform(t *testing.T) {
	p, err := DetectPlatform()
	if err != nil {
		t.Fatalf("DetectPlatform() error: %v", err)
	}

	if p.OS != runtime.GOOS {
		t.Errorf("OS = %q, want %q", p.OS, runtime.GOOS)
	}
	if p.Arch != runtime.GOARCH {
		t.Errorf("Arch = %q, want %q", p.Arch, runtime.GOARCH)
	}
	if p.LibName == "" {
		t.Error("LibName is empty")
	}
}

func TestPlatformZipName(t *testing.T) {
	tests := []struct {
		os, arch string
		want     string
	}{
		{"windows", "amd64", "wgpu-windows-x86_64-msvc-release.zip"},
		{"windows", "arm64", "wgpu-windows-aarch64-msvc-release.zip"},
		{"darwin", "amd64", "wgpu-macos-x86_64-release.zip"},
		{"darwin", "arm64", "wgpu-macos-aarch64-release.zip"},
		{"linux", "amd64", "wgpu-linux-x86_64-release.zip"},
		{"linux", "arm64", "wgpu-linux-aarch64-release.zip"},
	}

	for _, tt := range tests {
		p := &Platform{OS: tt.os, Arch: tt.arch}
		got := p.ZipName()
		if got != tt.want {
			t.Errorf("ZipName(%s/%s) = %q, want %q", tt.os, tt.arch, got, tt.want)
		}
	}
}

func TestPlatformDownloadURL(t *testing.T) {
	p := &Platform{OS: "windows", Arch: "amd64"}
	url := p.DownloadURL("v29.0.0.0")

	expected := "https://github.com/gfx-rs/wgpu-native/releases/download/v29.0.0.0/wgpu-windows-x86_64-msvc-release.zip"
	if url != expected {
		t.Errorf("DownloadURL() = %q, want %q", url, expected)
	}
}

func TestLibraryName(t *testing.T) {
	name := LibraryName()
	if name == "" {
		t.Error("LibraryName() is empty")
	}

	switch runtime.GOOS {
	case "windows":
		if name != "wgpu_native.dll" {
			t.Errorf("LibraryName() = %q, want wgpu_native.dll", name)
		}
	case "darwin":
		if name != "libwgpu_native.dylib" {
			t.Errorf("LibraryName() = %q, want libwgpu_native.dylib", name)
		}
	case "linux":
		if name != "libwgpu_native.so" {
			t.Errorf("LibraryName() = %q, want libwgpu_native.so", name)
		}
	}
}

func TestDetectPlatformUnsupportedArch(t *testing.T) {
	// Can't directly test unsupported arch on current platform,
	// but verify the function returns valid values for current platform.
	p, err := DetectPlatform()
	if err != nil {
		t.Skip("current platform not supported:", err)
	}

	if p.archName() != "x86_64" && p.archName() != "aarch64" {
		t.Errorf("archName() = %q, want x86_64 or aarch64", p.archName())
	}
}

func TestFindLibraryEnvPath(t *testing.T) {
	// Create a real file to point WGPU_NATIVE_PATH at.
	tmpFile := filepath.Join(t.TempDir(), "test_lib")
	if err := os.WriteFile(tmpFile, []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("WGPU_NATIVE_PATH", tmpFile)

	found := FindLibrary()
	if found != tmpFile {
		t.Errorf("FindLibrary() = %q, want %q", found, tmpFile)
	}
}

func TestFindLibraryEnvPathMissing(t *testing.T) {
	// When WGPU_NATIVE_PATH points to a non-existent file, FindLibrary should
	// fall through to the other search paths rather than returning the invalid path.
	t.Setenv("WGPU_NATIVE_PATH", "/nonexistent/path/wgpu_native.so")

	// We cannot assert a specific result (it depends on what is installed),
	// but we can assert it does not return the broken env path.
	found := FindLibrary()
	if found == "/nonexistent/path/wgpu_native.so" {
		t.Error("FindLibrary() returned non-existent WGPU_NATIVE_PATH value")
	}
}

func TestFindLibraryLibDir(t *testing.T) {
	// Create a temp dir and populate it with a lib subdirectory containing the library.
	tmpDir := t.TempDir()
	libDir := filepath.Join(tmpDir, "lib")
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		t.Fatal(err)
	}

	libName := LibraryName()
	libFile := filepath.Join(libDir, libName)
	if err := os.WriteFile(libFile, []byte("fake"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Unset env override so we exercise the file-system search.
	t.Setenv("WGPU_NATIVE_PATH", "")

	// FindLibrary resolves relative paths, so we must run from tmpDir.
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir) //nolint:errcheck // test cleanup

	found := FindLibrary()
	if found == "" {
		t.Fatal("FindLibrary() returned empty — expected to find lib in ./lib/")
	}

	absLib, _ := filepath.Abs(libFile)
	if found != absLib {
		t.Errorf("FindLibrary() = %q, want %q", found, absLib)
	}
}
