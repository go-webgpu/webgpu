// Command setup downloads and installs the wgpu-native library.
//
// Usage:
//
//	go run github.com/go-webgpu/webgpu/cmd/setup@latest
//	go run github.com/go-webgpu/webgpu/cmd/setup@latest ./path/to/lib
package main

import (
	"fmt"
	"os"

	"github.com/go-webgpu/webgpu/setup"
)

func main() {
	destDir := "lib"
	if len(os.Args) > 1 {
		destDir = os.Args[1]
	}

	if _, err := setup.Install(destDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
