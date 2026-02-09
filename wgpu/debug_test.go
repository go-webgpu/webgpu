package wgpu

import "testing"

func TestLeakDetection(t *testing.T) {
	SetDebugMode(true)
	defer SetDebugMode(false)
	defer ResetLeakTracker()

	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Before release — should report leak
	report := ReportLeaks()
	if report == nil {
		t.Fatal("expected leak report, got nil")
	}
	if report.Count == 0 {
		t.Error("expected at least 1 tracked resource")
	}
	t.Logf("Before release: %s", report)

	inst.Release()

	// After release — should be clean
	report = ReportLeaks()
	if report != nil {
		t.Errorf("expected no leaks after release, got: %s", report)
	}
}

func TestLeakDetectionDisabled(t *testing.T) {
	SetDebugMode(false)
	defer ResetLeakTracker()

	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer inst.Release()

	// Debug mode off — ReportLeaks returns nil
	report := ReportLeaks()
	if report != nil {
		t.Errorf("expected nil report when debug disabled, got: %s", report)
	}
}
