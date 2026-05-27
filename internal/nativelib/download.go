package nativelib

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Download(url string) (string, error) {
	resp, err := http.Get(url) //nolint:gosec // G107: URL constructed from constants
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: HTTP %d from %s", resp.StatusCode, url)
	}

	tmpFile, err := os.CreateTemp("", "wgpu-native-*.zip")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("download write: %w", err)
	}

	return tmpFile.Name(), nil
}

func ExtractLibrary(zipPath, destDir, libName string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		name := filepath.Base(f.Name)
		if name != libName {
			continue
		}

		src, err := f.Open()
		if err != nil {
			return "", fmt.Errorf("open %s in zip: %w", f.Name, err)
		}
		defer src.Close()

		destPath := filepath.Join(destDir, libName)
		dst, err := os.Create(destPath)
		if err != nil {
			return "", fmt.Errorf("create %s: %w", destPath, err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return "", fmt.Errorf("extract %s: %w", libName, err)
		}

		return destPath, nil
	}

	return "", fmt.Errorf("%s not found in archive", libName)
}
