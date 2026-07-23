package wgpu

import "testing"

func TestABICallbackInitializers(t *testing.T) {
	tests := []struct {
		name   string
		init   func()
		target *uintptr
	}{
		{name: "adapter", init: initAdapterCallback, target: &adapterCallbackPtr},
		{name: "device", init: initDeviceCallback, target: &deviceCallbackPtr},
		{name: "buffer map", init: initMapCallback, target: &mapCallbackPtr},
		{name: "error scope", init: initErrorScopeCallback, target: &errorScopeCallbackPtr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := *test.target
			t.Cleanup(func() {
				*test.target = original
			})

			test.init()
			if *test.target == 0 {
				t.Fatal("callback pointer is zero")
			}
		})
	}
}

func registerTestAdapterRequest(t *testing.T, requestID uintptr) *adapterRequest {
	t.Helper()
	req := &adapterRequest{done: make(chan struct{})}
	adapterRequestsMu.Lock()
	adapterRequests[requestID] = req
	adapterRequestsMu.Unlock()
	t.Cleanup(func() {
		adapterRequestsMu.Lock()
		delete(adapterRequests, requestID)
		adapterRequestsMu.Unlock()
	})
	return req
}

func assertCallbackCompleted(t *testing.T, done <-chan struct{}, message string) {
	assertCallbackMessage(t, done, message, "callback message")
}

func assertCallbackMessage(t *testing.T, done <-chan struct{}, message, want string) {
	t.Helper()
	select {
	case <-done:
	default:
		t.Fatal("callback did not complete the registered request")
	}
	if message != want {
		t.Fatalf("message = %q, want %q", message, want)
	}
}
