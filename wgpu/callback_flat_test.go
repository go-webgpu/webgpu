//go:build ((linux || darwin || freebsd) && (amd64 || arm64)) || (windows && arm64)

package wgpu

import (
	"testing"
	"unsafe"
)

func TestABICallbackEntriesPreserveStringViewAndUserdata(t *testing.T) {
	message := []byte("callback message")
	messageData := uintptr(unsafe.Pointer(&message[0]))
	messageLength := uintptr(len(message))

	t.Run("adapter", func(t *testing.T) {
		const requestID = uintptr(101)
		req := &adapterRequest{done: make(chan struct{})}
		adapterRequestsMu.Lock()
		adapterRequests[requestID] = req
		adapterRequestsMu.Unlock()
		t.Cleanup(func() {
			adapterRequestsMu.Lock()
			delete(adapterRequests, requestID)
			adapterRequestsMu.Unlock()
		})

		adapterCallbackEntry(7, 0, messageData, messageLength, requestID, 0)

		assertCallbackCompleted(t, req.done, req.message)
		if req.status != RequestAdapterStatus(7) {
			t.Fatalf("status = %d, want 7", req.status)
		}
	})

	t.Run("device", func(t *testing.T) {
		const requestID = uintptr(102)
		req := &deviceRequest{done: make(chan struct{})}
		deviceRequestsMu.Lock()
		deviceRequests[requestID] = req
		deviceRequestsMu.Unlock()
		t.Cleanup(func() {
			deviceRequestsMu.Lock()
			delete(deviceRequests, requestID)
			deviceRequestsMu.Unlock()
		})

		deviceCallbackEntry(8, 0, messageData, messageLength, requestID, 0)

		assertCallbackCompleted(t, req.done, req.message)
		if req.status != RequestDeviceStatus(8) {
			t.Fatalf("status = %d, want 8", req.status)
		}
	})

	t.Run("buffer map", func(t *testing.T) {
		const requestID = uintptr(103)
		req := &mapRequest{done: make(chan struct{})}
		mapRequestsMu.Lock()
		mapRequests[requestID] = req
		mapRequestsMu.Unlock()
		t.Cleanup(func() {
			mapRequestsMu.Lock()
			delete(mapRequests, requestID)
			mapRequestsMu.Unlock()
		})

		mapCallbackEntry(9, messageData, messageLength, requestID, 0)

		assertCallbackCompleted(t, req.done, req.message)
		if req.status != MapAsyncStatus(9) {
			t.Fatalf("status = %d, want 9", req.status)
		}
	})

	t.Run("error scope", func(t *testing.T) {
		const requestID = uintptr(104)
		result := &errorScopeResult{done: make(chan struct{})}
		errorScopeResultsMu.Lock()
		errorScopeResults[requestID] = result
		errorScopeResultsMu.Unlock()
		t.Cleanup(func() {
			errorScopeResultsMu.Lock()
			delete(errorScopeResults, requestID)
			errorScopeResultsMu.Unlock()
		})

		errorScopeCallbackEntry(10, 11, messageData, messageLength, requestID, 0)

		assertCallbackCompleted(t, result.done, result.message)
		if result.status != PopErrorScopeStatus(10) {
			t.Fatalf("status = %d, want 10", result.status)
		}
		if result.errType != ErrorType(11) {
			t.Fatalf("error type = %d, want 11", result.errType)
		}
	})
}

func assertCallbackCompleted(t *testing.T, done <-chan struct{}, message string) {
	t.Helper()
	select {
	case <-done:
	default:
		t.Fatal("callback did not complete the registered request")
	}
	if message != "callback message" {
		t.Fatalf("message = %q, want %q", message, "callback message")
	}
}
