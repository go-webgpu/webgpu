package wgpu

import (
	"sync"
	"testing"
)

func TestConcurrentInit(t *testing.T) {
	// Verify Init() is safe for concurrent calls
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = Init()
		}()
	}
	wg.Wait()
}

func TestConcurrentAdapterRequests(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer inst.Release()

	// Request multiple adapters concurrently
	var wg sync.WaitGroup
	errs := make([]error, 3)
	adapters := make([]*Adapter, 3)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			a, e := inst.RequestAdapter(nil)
			adapters[idx] = a
			errs[idx] = e
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("adapter %d failed: %v", i, err)
		}
		if adapters[i] != nil {
			adapters[i].Release()
		}
	}
}
