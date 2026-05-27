package nativelib

import (
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
