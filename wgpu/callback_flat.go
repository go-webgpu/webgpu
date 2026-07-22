//go:build ((linux || darwin || freebsd) && (amd64 || arm64)) || (windows && arm64)

package wgpu

// Callback entry implementations support amd64 and arm64 on Linux, macOS,
// FreeBSD, and Windows. Windows amd64 uses callback_windows_amd64.go;
// wgpu-native does not support other architectures.
//
// Unix amd64/arm64 and Windows ARM64 ABIs pass the two-word WGPUStringView
// callback argument by value in integer registers. goffi callbacks expose
// those words as separate uintptr arguments, so each entry reconstructs the
// view before invoking shared logic.

func adapterCallbackEntry(status, adapter, messageData, messageLength, userdata1, _ uintptr) uintptr {
	return handleAdapterCallback(status, adapter, StringView{Data: messageData, Length: messageLength}, userdata1)
}

func deviceCallbackEntry(status, device, messageData, messageLength, userdata1, _ uintptr) uintptr {
	return handleDeviceCallback(status, device, StringView{Data: messageData, Length: messageLength}, userdata1)
}

func mapCallbackEntry(status, messageData, messageLength, userdata1, _ uintptr) uintptr {
	return handleMapCallback(status, StringView{Data: messageData, Length: messageLength}, userdata1)
}

func errorScopeCallbackEntry(status, errType, messageData, messageLength, userdata1, _ uintptr) uintptr {
	return handleErrorScopeCallback(status, errType, StringView{Data: messageData, Length: messageLength}, userdata1)
}
