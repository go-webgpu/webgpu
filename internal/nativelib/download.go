// Package nativelib provides platform detection, download, and extraction
// of the wgpu-native binary library.
package nativelib

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// httpTimeout is the timeout for downloading the wgpu-native library.
// 120 seconds is generous for a ~15-20 MB file even on slow connections.
const httpTimeout = 120 * time.Second

// maxLibSize is the maximum decompressed size of a single zip entry.
// Protects against decompression bombs; wgpu-native is ~15-20 MB.
const maxLibSize = 200 * 1024 * 1024 // 200 MB

// TODO: add SHA256 checksum verification when wgpu-native starts publishing checksums.
// Currently relies on HTTPS transport security only.

// Download downloads the file at url to a temporary file and returns its path.
// The caller is responsible for removing the temporary file when done.
func Download(url string) (string, error) {
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url) //nolint:gosec // G107: URL constructed from constants
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // read-only close, error not actionable

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download %s: HTTP %d", url, resp.StatusCode)
	}

	if resp.ContentLength > 0 {
		fmt.Printf("Downloading %s (%.1f MB)...\n", url, float64(resp.ContentLength)/1024/1024)
	} else {
		fmt.Printf("Downloading %s...\n", url)
	}

	tmpFile, err := os.CreateTemp("", "wgpu-native-*.zip")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	if _, err := io.Copy(tmpFile, io.LimitReader(resp.Body, maxLibSize)); err != nil {
		tmpFile.Close() //nolint:errcheck // cleanup on write error
		_ = os.Remove(tmpPath)
		return "", fmt.Errorf("download write: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return "", fmt.Errorf("close %s: %w", tmpPath, err)
	}

	return tmpPath, nil
}

// ExtractLibrary extracts the file named libName from the zip archive at
// zipPath, writing it to destDir. Returns the path of the extracted file.
func ExtractLibrary(zipPath, destDir, libName string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("open zip: %w", err)
	}
	defer r.Close() //nolint:errcheck // read-only close

	for _, f := range r.File {
		name := filepath.Base(f.Name)
		if name != libName {
			continue
		}

		src, err := f.Open()
		if err != nil {
			return "", fmt.Errorf("open %s in zip: %w", f.Name, err)
		}
		defer src.Close() //nolint:errcheck // read-only close

		destPath := filepath.Join(destDir, libName)
		dst, err := os.Create(destPath)
		if err != nil {
			return "", fmt.Errorf("create %s: %w", destPath, err)
		}

		if _, err := io.Copy(dst, io.LimitReader(src, maxLibSize)); err != nil {
			dst.Close() //nolint:errcheck // cleanup on write error
			return "", fmt.Errorf("extract %s: %w", libName, err)
		}

		if err := dst.Close(); err != nil {
			return "", fmt.Errorf("close %s: %w", destPath, err)
		}

		return destPath, nil
	}

	return "", fmt.Errorf("%s not found in archive", libName)
}
