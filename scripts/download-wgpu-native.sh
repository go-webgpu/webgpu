#!/usr/bin/env bash
# Download wgpu-native library for the current platform
# Usage: ./scripts/download-wgpu-native.sh [version]
#
# This script downloads pre-built wgpu-native binaries from GitHub releases
# and places them in the lib/ directory.

set -e

# Configuration
WGPU_VERSION="${1:-v24.0.0.2}"
GITHUB_REPO="gfx-rs/wgpu-native"
LIB_DIR="lib"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Detect platform
detect_platform() {
    local os arch

    case "$(uname -s)" in
        Linux*)  os="linux" ;;
        Darwin*) os="macos" ;;
        MINGW*|MSYS*|CYGWIN*) os="windows" ;;
        *)
            log_error "Unsupported OS: $(uname -s)"
            exit 1
            ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64) arch="x86_64" ;;
        arm64|aarch64) arch="aarch64" ;;
        *)
            log_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac

    echo "${os}-${arch}"
}

# Get download URL for platform
get_download_url() {
    local platform="$1"
    local version="$2"
    local filename

    case "$platform" in
        linux-x86_64)
            filename="wgpu-linux-x86_64-release.zip"
            ;;
        linux-aarch64)
            filename="wgpu-linux-aarch64-release.zip"
            ;;
        macos-x86_64)
            filename="wgpu-macos-x86_64-release.zip"
            ;;
        macos-aarch64)
            filename="wgpu-macos-aarch64-release.zip"
            ;;
        windows-x86_64)
            filename="wgpu-windows-x86_64-release.zip"
            ;;
        *)
            log_error "Unknown platform: $platform"
            exit 1
            ;;
    esac

    echo "https://github.com/${GITHUB_REPO}/releases/download/${version}/${filename}"
}

# Get library filename for platform
get_lib_filename() {
    local platform="$1"

    case "$platform" in
        linux-*)
            echo "libwgpu_native.so"
            ;;
        macos-*)
            echo "libwgpu_native.dylib"
            ;;
        windows-*)
            echo "wgpu_native.dll"
            ;;
    esac
}

# Main
main() {
    log_info "wgpu-native downloader"
    log_info "Version: ${WGPU_VERSION}"
    echo ""

    # Detect platform
    local platform
    platform=$(detect_platform)
    log_info "Detected platform: ${platform}"

    # Create lib directory
    mkdir -p "${LIB_DIR}"

    # Get URLs and filenames
    local url lib_filename
    url=$(get_download_url "$platform" "$WGPU_VERSION")
    lib_filename=$(get_lib_filename "$platform")

    log_info "Download URL: ${url}"
    log_info "Library file: ${lib_filename}"
    echo ""

    # Check if already exists
    if [ -f "${LIB_DIR}/${lib_filename}" ]; then
        log_warning "Library already exists: ${LIB_DIR}/${lib_filename}"
        read -p "Overwrite? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Skipping download"
            exit 0
        fi
    fi

    # Download
    local tmp_dir tmp_zip
    tmp_dir=$(mktemp -d)
    tmp_zip="${tmp_dir}/wgpu-native.zip"

    log_info "Downloading..."
    if command -v curl &> /dev/null; then
        curl -L -o "${tmp_zip}" "${url}"
    elif command -v wget &> /dev/null; then
        wget -O "${tmp_zip}" "${url}"
    else
        log_error "Neither curl nor wget found"
        exit 1
    fi

    # Extract
    log_info "Extracting..."
    unzip -q -o "${tmp_zip}" -d "${tmp_dir}"

    # Find and copy library
    local found_lib
    found_lib=$(find "${tmp_dir}" -name "${lib_filename}" -type f | head -1)

    if [ -z "$found_lib" ]; then
        # Try alternative names
        found_lib=$(find "${tmp_dir}" -name "*.so" -o -name "*.dylib" -o -name "*.dll" | head -1)
    fi

    if [ -z "$found_lib" ]; then
        log_error "Library not found in archive"
        log_info "Archive contents:"
        ls -la "${tmp_dir}"
        rm -rf "${tmp_dir}"
        exit 1
    fi

    cp "${found_lib}" "${LIB_DIR}/${lib_filename}"

    # Cleanup
    rm -rf "${tmp_dir}"

    log_success "Downloaded: ${LIB_DIR}/${lib_filename}"
    echo ""

    # Show usage instructions
    log_info "Usage instructions:"
    case "$platform" in
        linux-*)
            echo "  export LD_LIBRARY_PATH=\$PWD/lib:\$LD_LIBRARY_PATH"
            echo "  # Or copy to system: sudo cp lib/${lib_filename} /usr/local/lib/"
            ;;
        macos-*)
            echo "  export DYLD_LIBRARY_PATH=\$PWD/lib:\$DYLD_LIBRARY_PATH"
            echo "  # Or copy to system: sudo cp lib/${lib_filename} /usr/local/lib/"
            ;;
        windows-*)
            echo "  # Copy lib/${lib_filename} to your PATH or project directory"
            echo "  # Or add lib/ to PATH: set PATH=%PATH%;%CD%\\lib"
            ;;
    esac
    echo ""

    log_success "Done!"
}

main "$@"
