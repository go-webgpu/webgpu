// Package setup provides wgpu-native library installation for go-webgpu.
//
// Usage:
//
//	go run github.com/go-webgpu/webgpu/cmd/setup@latest
//
// Or programmatically:
//
//	path, err := setup.Install("lib")
//	found := setup.FindLibrary()
package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-webgpu/webgpu/internal/nativelib"
)

const Version = "v29.0.0.0"

// Install downloads and installs the wgpu-native library to destDir.
// If destDir is empty, defaults to "./lib".
// Returns the absolute path to the installed library.
func Install(destDir string) (string, error) {
	platform, err := nativelib.DetectPlatform()
	if err != nil {
		return "", err
	}

	if destDir == "" {
		destDir = "lib"
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", fmt.Errorf("create directory %s: %w", destDir, err)
	}

	url := platform.DownloadURL(Version)
	fmt.Printf("Downloading wgpu-native %s for %s/%s...\n", Version, platform.OS, platform.Arch)
	fmt.Printf("URL: %s\n", url)

	zipPath, err := nativelib.Download(url)
	if err != nil {
		return "", err
	}
	defer os.Remove(zipPath)

	fmt.Printf("Extracting %s...\n", platform.LibName)
	libPath, err := nativelib.ExtractLibrary(zipPath, destDir, platform.LibName)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(libPath)
	if err != nil {
		absPath = libPath
	}
	fmt.Printf("Installed: %s\n\n", absPath)
	printUsage(platform, absPath, destDir)

	return absPath, nil
}

// FindLibrary searches common locations for the wgpu-native library.
// Returns the absolute path if found, empty string otherwise.
// Search order: WGPU_NATIVE_PATH env → ./lib/ → executable dir → PATH.
func FindLibrary() string {
	return nativelib.FindLibrary()
}

func printUsage(platform *nativelib.Platform, absPath, destDir string) {
	dir, err := filepath.Abs(destDir)
	if err != nil {
		dir = destDir
	}

	fmt.Println("To use, set environment variable:")
	switch platform.OS {
	case "windows":
		fmt.Printf("  set WGPU_NATIVE_PATH=%s\n", absPath)
	default:
		fmt.Printf("  export WGPU_NATIVE_PATH=%s\n", absPath)
	}

	fmt.Println("\nOr add directory to library path:")
	switch platform.OS {
	case "windows":
		fmt.Printf("  set PATH=%s;%%PATH%%\n", dir)
	case "darwin":
		fmt.Printf("  export DYLD_LIBRARY_PATH=%s:$DYLD_LIBRARY_PATH\n", dir)
	default: // linux and others
		fmt.Printf("  export LD_LIBRARY_PATH=%s:$LD_LIBRARY_PATH\n", dir)
	}
}
