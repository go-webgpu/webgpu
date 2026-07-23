package wgpu

import (
	"os"
	"testing"
)

type timestampPeriodProcStub struct {
	handle uintptr
	period float32
}

func (p *timestampPeriodProcStub) Call(args ...uintptr) (uintptr, uintptr, error) {
	return 0, 0, nil
}

func (p *timestampPeriodProcStub) CallFloat32(args ...uintptr) (float32, error) {
	if len(args) == 1 {
		p.handle = args[0]
	}
	return p.period, nil
}

func TestABIQueueGetTimestampPeriodNullGuard(t *testing.T) {
	var nilQueue *Queue
	if got := nilQueue.GetTimestampPeriod(); got != 0 {
		t.Fatalf("nil queue timestamp period = %v, want 0", got)
	}

	releasedQueue := &Queue{}
	if got := releasedQueue.GetTimestampPeriod(); got != 0 {
		t.Fatalf("released queue timestamp period = %v, want 0", got)
	}
}

type integerOnlyTimestampPeriodProc struct{}

func (*integerOnlyTimestampPeriodProc) Call(args ...uintptr) (uintptr, uintptr, error) {
	return 0, 0, nil
}

func TestABIQueueGetTimestampPeriodRequiresFloat32Proc(t *testing.T) {
	original := procQueueGetTimestampPeriod
	procQueueGetTimestampPeriod = &integerOnlyTimestampPeriodProc{}
	defer func() { procQueueGetTimestampPeriod = original }()

	if got := (&Queue{handle: 0x1234}).GetTimestampPeriod(); got != 0 {
		t.Fatalf("queue timestamp period = %v, want 0 for integer-only proc", got)
	}
}

func TestABIQueueGetTimestampPeriodUnavailable(t *testing.T) {
	original := procQueueGetTimestampPeriod
	procQueueGetTimestampPeriod = nil
	defer func() { procQueueGetTimestampPeriod = original }()

	if got := (&Queue{handle: 0x1234}).GetTimestampPeriod(); got != 0 {
		t.Fatalf("queue timestamp period = %v, want 0 for unavailable proc", got)
	}
}

func TestABIQueueGetTimestampPeriodUsesNativeFloat32(t *testing.T) {
	stub := &timestampPeriodProcStub{period: 0.125}
	original := procQueueGetTimestampPeriod
	procQueueGetTimestampPeriod = stub
	defer func() { procQueueGetTimestampPeriod = original }()

	got := (&Queue{handle: 0x1234}).GetTimestampPeriod()
	if got != stub.period {
		t.Fatalf("queue timestamp period = %v, want %v", got, stub.period)
	}
	if stub.handle != 0x1234 {
		t.Fatalf("queue handle = %#x, want %#x", stub.handle, uintptr(0x1234))
	}
}

func TestABIQueueGetTimestampPeriodDynamicLibrary(t *testing.T) {
	path := os.Getenv("WGPU_TIMESTAMP_PERIOD_ABI_STUB_LIBRARY")
	if path == "" {
		t.Skip("set WGPU_TIMESTAMP_PERIOD_ABI_STUB_LIBRARY to a shared library exporting the test symbol")
	}
	library, err := loadLibrary(path)
	if err != nil {
		t.Fatal(err)
	}
	proc, ok := library.NewProc("wgpuQueueGetTimestampPeriod").(float32Proc)
	if !ok {
		t.Fatal("platform loader does not implement float32 return calls")
	}
	got, err := proc.CallFloat32(0x1234)
	if err != nil {
		t.Fatal(err)
	}
	if got != 0.125 {
		t.Fatalf("dynamic library timestamp period = %v, want 0.125", got)
	}
}
