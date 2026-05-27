package nativelib

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	libWindows = "wgpu_native.dll"
	libDarwin  = "libwgpu_native.dylib"
	libLinux   = "libwgpu_native.so"
)

type Platform struct {
	OS      string
	Arch    string
	LibName string
}

func DetectPlatform() (*Platform, error) {
	p := &Platform{OS: runtime.GOOS, Arch: runtime.GOARCH}

	switch p.OS {
	case "windows":
		p.LibName = libWindows
	case "darwin":
		p.LibName = libDarwin
	case "linux":
		p.LibName = libLinux
	default:
		return nil, fmt.Errorf("unsupported OS: %s", p.OS)
	}

	if p.Arch != "amd64" && p.Arch != "arm64" {
		return nil, fmt.Errorf("unsupported architecture: %s", p.Arch)
	}

	return p, nil
}

func (p *Platform) archName() string {
	if p.Arch == "amd64" {
		return "x86_64"
	}
	return "aarch64"
}

func (p *Platform) ZipName() string {
	arch := p.archName()
	switch p.OS {
	case "windows":
		return fmt.Sprintf("wgpu-%s-%s-msvc-release.zip", p.OS, arch)
	case "darwin":
		return fmt.Sprintf("wgpu-macos-%s-release.zip", arch)
	case "linux":
		return fmt.Sprintf("wgpu-%s-%s-release.zip", p.OS, arch)
	}
	return ""
}

func (p *Platform) DownloadURL(version string) string {
	return fmt.Sprintf("https://github.com/gfx-rs/wgpu-native/releases/download/%s/%s", version, p.ZipName())
}

func LibraryName() string {
	switch runtime.GOOS {
	case "windows":
		return libWindows
	case "darwin":
		return libDarwin
	default:
		return libLinux
	}
}

func FindLibrary() string {
	if p := os.Getenv("WGPU_NATIVE_PATH"); p != "" {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	libName := LibraryName()
	searchPaths := []string{
		filepath.Join(".", libName),
		filepath.Join("lib", libName),
	}

	if exe, err := os.Executable(); err == nil {
		searchPaths = append(searchPaths, filepath.Join(filepath.Dir(exe), libName))
	}

	for _, p := range searchPaths {
		if _, err := os.Stat(p); err == nil {
			abs, _ := filepath.Abs(p)
			return abs
		}
	}

	for _, dir := range strings.Split(os.Getenv("PATH"), string(os.PathListSeparator)) {
		p := filepath.Join(dir, libName)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}
