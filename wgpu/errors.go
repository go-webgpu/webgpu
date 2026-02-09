package wgpu

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/go-webgpu/goffi/ffi"
)

// PushErrorScope pushes an error scope for catching GPU errors.
// Errors of the specified filter type will be caught until PopErrorScope is called.
// Error scopes are LIFO (stack-based) - last pushed scope is popped first.
//
// IMPORTANT: You must call PopErrorScope for each PushErrorScope.
// Popping an empty stack will cause a panic in wgpu-native (known limitation).
// Users should track push/pop calls manually to avoid stack underflow.
//
// Example usage:
//
//	device.PushErrorScope(ErrorFilterValidation)
//	// ... GPU operations that might produce validation errors
//	errType, message := device.PopErrorScope(instance)
//	if errType != ErrorTypeNoError {
//	    log.Printf("Validation error: %s", message)
//	}
func (d *Device) PushErrorScope(filter ErrorFilter) {
	mustInit()
	// nolint:errcheck // PushErrorScope has no meaningful return value to check
	procDevicePushErrorScope.Call(d.handle, uintptr(filter))
}

// popErrorScopeCallbackInfo matches WGPUPopErrorScopeCallbackInfo C struct.
type popErrorScopeCallbackInfo struct {
	nextInChain uintptr // *ChainedStruct
	mode        CallbackMode
	callback    uintptr // function pointer
	userdata1   uintptr
	userdata2   uintptr
}

// errorScopeResult holds the result of a PopErrorScope operation.
type errorScopeResult struct {
	done    chan struct{}
	status  PopErrorScopeStatus
	errType ErrorType
	message string
}

var (
	// Global registry for pending error scope operations
	errorScopeResults   = make(map[uintptr]*errorScopeResult)
	errorScopeResultsMu sync.Mutex
	errorScopeResultID  uintptr

	// Callback function pointer (created once)
	errorScopeCallbackPtr  uintptr
	errorScopeCallbackOnce sync.Once
)

// errorScopeCallbackHandler is the Go function called by C code via ffi.NewCallback.
// Signature matches: void callback(WGPUPopErrorScopeStatus status, WGPUErrorType type,
//
//	WGPUStringView message, void* userdata1, void* userdata2)
func errorScopeCallbackHandler(status uintptr, errType uintptr, message uintptr, userdata1, _ uintptr) uintptr {
	// Extract message string (message is pointer to StringView)
	var msg string
	if message != 0 {
		// nolint:govet,gosec // message is uintptr from FFI callback - GC safe
		sv := (*StringView)(unsafe.Pointer(message))
		if sv.Data != 0 && sv.Length > 0 && sv.Length < 1<<20 {
			// nolint:govet,gosec // sv.Data is uintptr from C memory - GC safe
			msg = unsafe.String((*byte)(unsafe.Pointer(sv.Data)), int(sv.Length))
		}
	}

	// Find and complete the operation
	errorScopeResultsMu.Lock()
	result, ok := errorScopeResults[userdata1]
	if ok {
		delete(errorScopeResults, userdata1)
	}
	errorScopeResultsMu.Unlock()

	if ok && result != nil {
		result.status = PopErrorScopeStatus(status)
		result.errType = ErrorType(errType)
		result.message = msg
		close(result.done)
	}

	return 0 // void return
}

// initErrorScopeCallback creates the C callback function pointer using goffi.
func initErrorScopeCallback() {
	errorScopeCallbackPtr = ffi.NewCallback(errorScopeCallbackHandler)
}

// Deprecated: PopErrorScope panics on failure. Use PopErrorScopeAsync instead.
//
// PopErrorScope pops the current error scope and returns the first error caught.
// This is a synchronous wrapper that blocks until the result is available.
//
// IMPORTANT: You must have pushed an error scope before calling this.
// Calling PopErrorScope on an empty stack will cause a panic in wgpu-native.
// Use PopErrorScopeAsync if you need to handle empty stack gracefully.
//
// Returns:
//   - ErrorType: The type of error that occurred (ErrorTypeNoError if no error)
//   - string: Error message (empty if no error)
//
// Note: Error scopes are LIFO - the last pushed scope is popped first.
func (d *Device) PopErrorScope(instance *Instance) (ErrorType, string) {
	errType, message, err := d.PopErrorScopeAsync(instance)
	if err != nil {
		// Panic on error since this is the "unsafe" simple version
		panic(fmt.Sprintf("PopErrorScope failed: %v", err))
	}
	return errType, message
}

// PopErrorScopeAsync pops the current error scope and returns the first error caught.
// This version returns an error instead of panicking if the operation fails.
//
// Returns:
//   - ErrorType: The type of error that occurred (ErrorTypeNoError if no error)
//   - string: Error message (empty if no error)
//   - error: An error if the operation failed (e.g., empty stack)
//
// Note: Error scopes are LIFO - the last pushed scope is popped first.
// If the error scope stack is empty, returns an error instead of panicking.
func (d *Device) PopErrorScopeAsync(instance *Instance) (ErrorType, string, error) {
	if err := checkInit(); err != nil {
		return ErrorTypeNoError, "", err
	}

	if instance == nil {
		return ErrorTypeNoError, "", &WGPUError{Op: "PopErrorScopeAsync", Message: "instance is required for PopErrorScope"}
	}

	// Initialize callback once
	errorScopeCallbackOnce.Do(initErrorScopeCallback)

	// Create result holder
	result := &errorScopeResult{
		done: make(chan struct{}),
	}

	// Register result
	errorScopeResultsMu.Lock()
	errorScopeResultID++
	resultID := errorScopeResultID
	errorScopeResults[resultID] = result
	errorScopeResultsMu.Unlock()

	// Prepare callback info
	callbackInfo := popErrorScopeCallbackInfo{
		nextInChain: 0,
		mode:        CallbackModeAllowProcessEvents,
		callback:    errorScopeCallbackPtr,
		userdata1:   resultID,
		userdata2:   0,
	}

	// Call wgpuDevicePopErrorScope (returns WGPUFuture)
	// nolint:errcheck // Error handling is done via callback, not return value
	// nolint:gosec // FFI requires unsafe.Pointer conversion for struct passing
	procDevicePopErrorScope.Call(
		d.handle,
		uintptr(unsafe.Pointer(&callbackInfo)),
	)

	// Process events until callback fires
	// With CallbackModeAllowProcessEvents, we need to call ProcessEvents
	for {
		select {
		case <-result.done:
			// Callback completed
			if result.status != PopErrorScopeStatusSuccess {
				switch result.status {
				case PopErrorScopeStatusEmptyStack:
					return ErrorTypeNoError, "", &WGPUError{Op: "PopErrorScopeAsync", Message: "error scope stack is empty"}
				case PopErrorScopeStatusInstanceDropped:
					return ErrorTypeNoError, "", &WGPUError{Op: "PopErrorScopeAsync", Message: "instance was dropped"}
				default:
					return ErrorTypeNoError, "", &WGPUError{Op: "PopErrorScopeAsync", Message: fmt.Sprintf("pop error scope failed with status %d", result.status)}
				}
			}
			return result.errType, result.message, nil
		default:
			// Process events to fire callbacks
			instance.ProcessEvents()
		}
	}
}
