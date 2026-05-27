package setup

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version is empty")
	}
	if Version[0] != 'v' {
		t.Errorf("Version = %q, should start with 'v'", Version)
	}
}

func TestFindLibrary(t *testing.T) {
	// FindLibrary searches for the native library.
	// On dev machines it may or may not be found — just verify no panic.
	path := FindLibrary()
	t.Logf("FindLibrary() = %q", path)
}
