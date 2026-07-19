//go:build windows && amd64

package wgpu

// Windows x64 passes a WGPUStringView callback argument indirectly because
// the aggregate is larger than one register. Normalize that pointer into the
// same value form used by the shared callback logic.

func adapterCallbackEntry(status, adapter, message, userdata1, _ uintptr) uintptr {
	return handleAdapterCallback(status, adapter, callbackStringView(message), userdata1)
}

func deviceCallbackEntry(status, device, message, userdata1, _ uintptr) uintptr {
	return handleDeviceCallback(status, device, callbackStringView(message), userdata1)
}

func mapCallbackEntry(status, message, userdata1, _ uintptr) uintptr {
	return handleMapCallback(status, callbackStringView(message), userdata1)
}

func errorScopeCallbackEntry(status, errType, message, userdata1, _ uintptr) uintptr {
	return handleErrorScopeCallback(status, errType, callbackStringView(message), userdata1)
}

func callbackStringView(message uintptr) StringView {
	if message == 0 {
		return StringView{}
	}
	return *(*StringView)(ptrFromUintptr(message))
}
