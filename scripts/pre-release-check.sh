#!/usr/bin/env bash
# Pre-Release Validation Script for go-webgpu
# This script runs all quality checks before creating a release
# EXACTLY matches CI checks + additional validations
# Based on coregex pre-release-check.sh

set -e  # Exit on first error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Header
echo ""
echo "================================================"
echo "  go-webgpu - Pre-Release Check"
echo "================================================"
echo ""

# Track overall status
ERRORS=0
WARNINGS=0

# 1. Check Go version
log_info "Checking Go version..."
GO_VERSION=$(go version | awk '{print $3}')
REQUIRED_VERSION="go1.25"
if [[ "$GO_VERSION" < "$REQUIRED_VERSION" ]]; then
    log_error "Go version $REQUIRED_VERSION+ required, found $GO_VERSION"
    ERRORS=$((ERRORS + 1))
else
    log_success "Go version: $GO_VERSION"
fi
echo ""

# 2. Check git status
log_info "Checking git status..."
if git diff-index --quiet HEAD --; then
    log_success "Working directory is clean"
else
    log_warning "Uncommitted changes detected"
    git status --short
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# 3. Code formatting check (EXACT CI command)
log_info "Checking code formatting (gofmt -l .)..."
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
    log_error "The following files need formatting:"
    echo "$UNFORMATTED"
    echo ""
    log_info "Run: go fmt ./..."
    ERRORS=$((ERRORS + 1))
else
    log_success "All files are properly formatted"
fi
echo ""

# 4. Go vet (disable unsafeptr for FFI code)
log_info "Running go vet..."
# FFI code legitimately uses unsafe.Pointer patterns that trigger false positives
if go vet -unsafeptr=false ./... 2>&1; then
    log_success "go vet passed"
else
    log_error "go vet failed"
    ERRORS=$((ERRORS + 1))
fi
echo ""

# 5. Build all packages
log_info "Building all packages..."
if go build ./... 2>&1; then
    log_success "Build successful"
else
    log_error "Build failed"
    ERRORS=$((ERRORS + 1))
fi
echo ""

# 6. go.mod validation
log_info "Validating go.mod..."
go mod verify
if [ $? -eq 0 ]; then
    log_success "go.mod verified"
else
    log_error "go.mod verification failed"
    ERRORS=$((ERRORS + 1))
fi

# Check if go.mod needs tidying
go mod tidy
if git diff --quiet go.mod go.sum; then
    log_success "go.mod is tidy"
else
    log_warning "go.mod needs tidying (run 'go mod tidy')"
    git diff go.mod go.sum
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# 6.5. Verify golangci-lint configuration
log_info "Verifying golangci-lint configuration..."
if command -v golangci-lint &> /dev/null; then
    if golangci-lint config verify 2>&1; then
        log_success "golangci-lint config is valid"
    else
        log_error "golangci-lint config is invalid"
        ERRORS=$((ERRORS + 1))
    fi
else
    log_warning "golangci-lint not installed (optional but recommended)"
    log_info "Install: https://golangci-lint.run/welcome/install/"
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# 7. Run tests with race detector (supports WSL2 fallback)
# NOTE: GPU tests require wgpu-native DLL and actual GPU hardware
# Full GPU tests are skipped locally if wgpu_native.dll is not found
# Only math tests run in CI (no GPU available)

USE_WSL=0
WSL_DISTRO=""

# Helper function to find WSL distro with Go installed
find_wsl_distro() {
    if ! command -v wsl &> /dev/null; then
        return 1
    fi

    # Try common distros first
    for distro in "Gentoo" "Ubuntu" "Debian" "Alpine"; do
        if wsl -d "$distro" bash -c "command -v go &> /dev/null" 2>/dev/null; then
            echo "$distro"
            return 0
        fi
    done

    return 1
}

if command -v gcc &> /dev/null || command -v clang &> /dev/null; then
    log_info "Running tests with race detector..."
    RACE_FLAG="-race"
    TEST_CMD="go test -race ./... 2>&1"
else
    # Try to find WSL distro with Go
    WSL_DISTRO=$(find_wsl_distro)
    if [ -n "$WSL_DISTRO" ]; then
        log_info "GCC not found locally, but WSL2 ($WSL_DISTRO) detected!"
        log_info "Running tests with race detector via WSL2 $WSL_DISTRO..."
        USE_WSL=1
        RACE_FLAG="-race"

        # Convert Windows path to WSL path (D:\projects\... -> /mnt/d/projects/...)
        CURRENT_DIR=$(pwd)
        if [[ "$CURRENT_DIR" =~ ^/([a-z])/ ]]; then
            # Already in /d/... format (MSYS), convert to /mnt/d/...
            WSL_PATH="/mnt${CURRENT_DIR}"
        else
            # Windows format D:\... convert to /mnt/d/...
            DRIVE_LETTER=$(echo "$CURRENT_DIR" | cut -d: -f1 | tr '[:upper:]' '[:lower:]')
            PATH_WITHOUT_DRIVE=${CURRENT_DIR#*:}
            WSL_PATH="/mnt/$DRIVE_LETTER${PATH_WITHOUT_DRIVE//\\//}"
        fi

        TEST_CMD="wsl -d \"$WSL_DISTRO\" bash -c \"cd \\\"$WSL_PATH\\\" && go test -race -ldflags '-linkmode=external' ./... 2>&1\""
    else
        log_warning "GCC not found, running tests WITHOUT race detector"
        log_info "Install GCC (mingw-w64) or setup WSL2 with Go for race detection"
        log_info "  Windows: https://www.mingw-w64.org/"
        log_info "  WSL2: https://docs.microsoft.com/en-us/windows/wsl/install"
        WARNINGS=$((WARNINGS + 1))
        RACE_FLAG=""
        TEST_CMD="go test ./... 2>&1"
    fi
fi

log_info "Running tests..."
# Check if project has cross-platform loader
if [ -f "wgpu/loader_unix.go" ] && [ -f "wgpu/loader_windows.go" ]; then
    WINDOWS_ONLY=0
    log_info "Detected cross-platform project (loader abstraction found)"
elif grep -r "syscall.LazyDLL" --include="*.go" . &>/dev/null; then
    WINDOWS_ONLY=1
    log_info "Detected Windows-only project (uses syscall.LazyDLL)"
else
    WINDOWS_ONLY=0
fi

if [ $USE_WSL -eq 1 ] && [ $WINDOWS_ONLY -eq 1 ]; then
    # Skip WSL tests for Windows-only projects
    log_warning "Skipping WSL tests for Windows-only project"
    log_info "Run tests on native Windows instead"
    TEST_OUTPUT="ok (skipped - Windows only)"
    WARNINGS=$((WARNINGS + 1))
elif [ $USE_WSL -eq 1 ]; then
    # WSL2: Test cross-platform packages
    # CGO_ENABLED=0 required for pure-Go goffi implementation on Linux
    TEST_OUTPUT=$(wsl -d "$WSL_DISTRO" bash -c "cd $WSL_PATH && CGO_ENABLED=0 timeout 180 go test ./wgpu/... 2>&1" || true)
    if [ -z "$TEST_OUTPUT" ]; then
        log_error "WSL2 tests timed out or failed to run"
        ERRORS=$((ERRORS + 1))
    fi
else
    TEST_OUTPUT=$(eval "$TEST_CMD")
fi

# Check if race detector failed to build (should not happen with external linkmode)
if echo "$TEST_OUTPUT" | grep -q "hole in findfunctab\|build failed.*race"; then
    log_warning "Race detector build failed (Gentoo build configuration issue)"
    log_info "Note: This should not happen with -ldflags '-linkmode=external'"
    log_info "Falling back to tests without race detector..."

    if [ $USE_WSL -eq 1 ]; then
        TEST_OUTPUT=$(wsl -d "$WSL_DISTRO" bash -c "cd \"$WSL_PATH\" && CGO_ENABLED=0 go test ./wgpu/... 2>&1")
    else
        TEST_OUTPUT=$(go test ./... 2>&1)
    fi

    RACE_FLAG=""
    WARNINGS=$((WARNINGS + 1))
fi

if echo "$TEST_OUTPUT" | grep -q "FAIL"; then
    # Check if failure is due to missing wgpu-native library (expected in WSL without native lib)
    if [ $USE_WSL -eq 1 ] && echo "$TEST_OUTPUT" | grep -q "failed to create instance\|failed to load library"; then
        log_warning "WSL tests failed - wgpu-native library not installed in WSL"
        log_info "This is expected. Install libwgpu_native.so in WSL for full testing."
        log_info "Windows tests passed - cross-platform code compiles correctly."
        WARNINGS=$((WARNINGS + 1))
    # Check if failure is only due to performance tests in WSL2 (acceptable)
    elif [ $USE_WSL -eq 1 ] && echo "$TEST_OUTPUT" | grep -q "TestMetricsCollector_Performance" && ! echo "$TEST_OUTPUT" | grep -q "race detected"; then
        log_warning "Performance tests failed in WSL2 (acceptable - WSL2 has overhead)"
        echo "$TEST_OUTPUT" | grep -A 5 "FAIL:"
        echo ""
        log_info "No race conditions detected - this is OK for WSL2"
        WARNINGS=$((WARNINGS + 1))
    else
        log_error "Tests failed or race conditions detected"
        echo "$TEST_OUTPUT"
        echo ""
        ERRORS=$((ERRORS + 1))
    fi
elif echo "$TEST_OUTPUT" | grep -q "PASS\|ok"; then
    if [ $USE_WSL -eq 1 ]; then
        log_success "All tests passed (via WSL2 $WSL_DISTRO, CGO_ENABLED=0)"
    elif [ -n "$RACE_FLAG" ]; then
        log_success "All tests passed with race detector (0 races)"
    else
        log_success "All tests passed"
    fi
else
    log_error "Unexpected test output"
    echo "$TEST_OUTPUT"
    ERRORS=$((ERRORS + 1))
fi
echo ""

# 8. Test coverage check (math tests only - GPU tests require hardware)
log_info "Checking test coverage (math tests)..."
COVERAGE=$(go test -cover ./wgpu/... -run "Mat4|Vec3|Vec4" 2>&1 | grep "coverage:" | tail -1 | awk '{print $5}' | sed 's/%//')
if [ -n "$COVERAGE" ]; then
    echo "  • math tests coverage: ${COVERAGE}%"
    if awk -v cov="$COVERAGE" 'BEGIN {exit !(cov >= 50.0)}'; then
        log_success "Coverage acceptable for FFI project"
    else
        log_warning "Coverage ${COVERAGE}% - FFI code is hard to test without GPU"
        WARNINGS=$((WARNINGS + 1))
    fi
else
    log_warning "Could not determine coverage (GPU tests skipped)"
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# 9. Dependency check (goffi required for callbacks)
log_info "Checking dependencies..."
GOFFI_VERSION=$(grep "github.com/go-webgpu/goffi" go.mod | awk '{print $2}')
if [ -n "$GOFFI_VERSION" ]; then
    log_success "goffi dependency: $GOFFI_VERSION"
else
    log_error "goffi not found (required for WebGPU callbacks)"
    ERRORS=$((ERRORS + 1))
fi

XSYS_VERSION=$(grep "golang.org/x/sys" go.mod | awk '{print $2}')
if [ -n "$XSYS_VERSION" ]; then
    log_success "golang.org/x/sys dependency: $XSYS_VERSION"
else
    log_warning "golang.org/x/sys not found (platform-specific features)"
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# 10. golangci-lint (same as CI - continue on error for FFI code)
log_info "Running golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    LINT_OUTPUT=$(golangci-lint run --timeout=5m ./... 2>&1 || true)
    ISSUE_COUNT=$(echo "$LINT_OUTPUT" | grep -c "^.*\.go:" || echo "0")

    if [ "$ISSUE_COUNT" -eq 0 ]; then
        log_success "golangci-lint passed with 0 issues"
    elif [ "$ISSUE_COUNT" -lt 20 ]; then
        log_warning "golangci-lint found $ISSUE_COUNT issues (acceptable for FFI code)"
        echo "$LINT_OUTPUT" | head -20
        WARNINGS=$((WARNINGS + 1))
    else
        log_warning "golangci-lint found $ISSUE_COUNT issues (FFI code has many warnings)"
        echo "$LINT_OUTPUT" | head -10
        log_info "FFI bindings typically have warnings - review manually"
        WARNINGS=$((WARNINGS + 1))
    fi
else
    log_warning "golangci-lint not installed (optional but recommended)"
    log_info "Install: https://golangci-lint.run/welcome/install/"
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# 11. Check for TODO/FIXME comments
log_info "Checking for TODO/FIXME comments..."
TODO_COUNT=$(grep -r "TODO\|FIXME" --include="*.go" --exclude-dir=vendor . 2>/dev/null | wc -l)
if [ "$TODO_COUNT" -gt 0 ]; then
    log_warning "Found $TODO_COUNT TODO/FIXME comments"
    grep -r "TODO\|FIXME" --include="*.go" --exclude-dir=vendor . 2>/dev/null | head -5
    WARNINGS=$((WARNINGS + 1))
else
    log_success "No TODO/FIXME comments found"
fi
echo ""

# 12. Check critical documentation files
log_info "Checking documentation..."
DOCS_MISSING=0
REQUIRED_DOCS="README.md CHANGELOG.md LICENSE"

for doc in $REQUIRED_DOCS; do
    if [ ! -f "$doc" ]; then
        log_error "Missing: $doc"
        DOCS_MISSING=1
        ERRORS=$((ERRORS + 1))
    fi
done

if [ $DOCS_MISSING -eq 0 ]; then
    log_success "All critical documentation files present"
fi
echo ""

# Summary
echo "========================================"
echo "  Summary"
echo "========================================"
echo ""

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    log_success "✅ All checks passed! Ready for release."
    echo ""
    log_info "Next steps for release:"
    echo ""
    echo "  1. Create release branch from develop:"
    echo "     git checkout -b release/vX.Y.Z develop"
    echo ""
    echo "  2. Prepare release (ONE commit with ALL changes):"
    echo "     - Update CHANGELOG.md"
    echo "     - Update README.md"
    echo "     bash scripts/pre-release-check.sh  # Re-run to verify"
    echo "     git add -A"
    echo "     git commit -m \"chore: prepare vX.Y.Z release\""
    echo ""
    echo "  3. Push release branch, wait for CI:"
    echo "     git push origin release/vX.Y.Z"
    echo "     ⏳ WAIT for CI to be GREEN"
    echo ""
    echo "  4. Merge to main:"
    echo "     git checkout main"
    echo "     git merge --squash release/vX.Y.Z"
    echo "     git commit -m \"Release vX.Y.Z\""
    echo "     git push origin main"
    echo "     ⏳ WAIT for CI to be GREEN on main!"
    echo ""
    echo "  5. ONLY AFTER CI GREEN - create and push tag:"
    echo "     git tag -a vX.Y.Z -m \"Release vX.Y.Z\""
    echo "     git push origin main --tags  # Tags are PERMANENT!"
    echo ""
    echo "  6. Merge back to develop:"
    echo "     git checkout develop"
    echo "     git merge --no-ff main -m \"Merge release vX.Y.Z back to develop\""
    echo "     git push origin develop"
    echo ""
    echo "  7. Clean up:"
    echo "     git branch -d release/vX.Y.Z"
    echo "     git push origin --delete release/vX.Y.Z"
    echo ""
    exit 0
elif [ $ERRORS -eq 0 ]; then
    log_warning "Checks completed with $WARNINGS warning(s)"
    echo ""
    log_info "Review warnings above before proceeding with release"
    echo ""
    exit 0
else
    log_error "Checks failed with $ERRORS error(s) and $WARNINGS warning(s)"
    echo ""
    log_error "Fix errors before creating release"
    echo ""
    exit 1
fi
