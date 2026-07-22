package wgpu

import "testing"

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
