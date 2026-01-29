package wgpu

import (
	"errors"
	"sync"
	"unsafe"

	"github.com/go-webgpu/goffi/ffi"
	"github.com/gogpu/gputypes"
)

// RequestDeviceCallbackInfo holds callback configuration for RequestDevice.
type RequestDeviceCallbackInfo struct {
	NextInChain uintptr // *ChainedStruct
	Mode        CallbackMode
	Callback    uintptr // Function pointer
	Userdata1   uintptr
	Userdata2   uintptr
}

// deviceRequest holds state for an async device request.
type deviceRequest struct {
	done    chan struct{}
	device  *Device
	status  RequestDeviceStatus
	message string
}

var (
	// Global registry for pending device requests
	deviceRequests   = make(map[uintptr]*deviceRequest)
	deviceRequestsMu sync.Mutex
	deviceRequestID  uintptr

	// Callback function pointer (created once)
	deviceCallbackPtr  uintptr
	deviceCallbackOnce sync.Once
)

// deviceCallbackHandler is the Go function called by C code via ffi.NewCallback.
// Signature: void(status uint32, device uintptr, message *StringView, userdata1 uintptr, userdata2 uintptr)
func deviceCallbackHandler(status uintptr, device uintptr, message uintptr, userdata1, userdata2 uintptr) uintptr {
	// Extract message string (message is pointer to StringView on Windows)
	var msg string
	if message != 0 {
		// nolint:govet // message is uintptr from FFI callback - GC safe
		sv := (*StringView)(unsafe.Pointer(message))
		if sv.Data != 0 && sv.Length > 0 && sv.Length < 1<<20 {
			// nolint:govet // sv.Data is uintptr from C memory - GC safe
			msg = unsafe.String((*byte)(unsafe.Pointer(sv.Data)), int(sv.Length))
		}
	}

	// Find and complete the request
	deviceRequestsMu.Lock()
	req, ok := deviceRequests[userdata1]
	if ok {
		delete(deviceRequests, userdata1)
	}
	deviceRequestsMu.Unlock()

	if ok && req != nil {
		req.status = RequestDeviceStatus(status)
		if device != 0 {
			req.device = &Device{handle: device}
		}
		req.message = msg
		close(req.done)
	}
	return 0
}

// initDeviceCallback creates the C callback function pointer using goffi.
func initDeviceCallback() {
	deviceCallbackPtr = ffi.NewCallback(deviceCallbackHandler)
}

// RequestDevice requests a GPU device from the adapter.
// This is a synchronous wrapper that blocks until the device is available.
func (a *Adapter) RequestDevice(options *DeviceDescriptor) (*Device, error) {
	mustInit()

	// Initialize callback once
	deviceCallbackOnce.Do(initDeviceCallback)

	// Create request state
	req := &deviceRequest{
		done: make(chan struct{}),
	}

	// Register request
	deviceRequestsMu.Lock()
	deviceRequestID++
	reqID := deviceRequestID
	deviceRequests[reqID] = req
	deviceRequestsMu.Unlock()

	// Prepare options
	var optionsPtr uintptr
	if options != nil {
		optionsPtr = uintptr(unsafe.Pointer(options))
	}

	// Prepare callback info
	callbackInfo := RequestDeviceCallbackInfo{
		NextInChain: 0,
		Mode:        CallbackModeAllowProcessEvents,
		Callback:    deviceCallbackPtr,
		Userdata1:   reqID,
		Userdata2:   0,
	}

	// Call wgpuAdapterRequestDevice
	procAdapterRequestDevice.Call( //nolint:errcheck
		a.handle,
		optionsPtr,
		uintptr(unsafe.Pointer(&callbackInfo)),
	)

	// Process events until callback fires
	// We need an instance to call ProcessEvents - get it from global or use a busy loop
	for {
		select {
		case <-req.done:
			// Callback completed
			if req.status != RequestDeviceStatusSuccess {
				msg := req.message
				if msg == "" {
					msg = "device request failed"
				}
				return nil, errors.New("wgpu: " + msg)
			}
			return req.device, nil
		default:
			// Brief pause to avoid busy spinning
			// In real usage, you'd call instance.ProcessEvents()
		}
	}
}

// GetQueue returns the default queue for the device.
func (d *Device) GetQueue() *Queue {
	mustInit()
	handle, _, _ := procDeviceGetQueue.Call(d.handle)
	if handle == 0 {
		return nil
	}
	return &Queue{handle: handle}
}

// Poll polls the device for completed work.
// If wait is true, blocks until there is work to process.
// Returns true if the queue is empty.
// This is a wgpu-native extension.
func (d *Device) Poll(wait bool) bool {
	mustInit()
	var waitArg uintptr
	if wait {
		waitArg = 1
	}
	result, _, _ := procDevicePoll.Call(d.handle, waitArg, 0)
	return result != 0
}

// Release releases the device resources.
func (d *Device) Release() {
	if d.handle != 0 {
		procDeviceRelease.Call(d.handle) //nolint:errcheck
		d.handle = 0
	}
}

// Release releases the queue resources.
func (q *Queue) Release() {
	if q.handle != 0 {
		procQueueRelease.Call(q.handle) //nolint:errcheck
		q.handle = 0
	}
}

// DeviceDescriptor configures device creation.
// For now, passing nil uses default settings.
type DeviceDescriptor struct {
	NextInChain uintptr // *ChainedStruct
	// Additional fields can be added as needed
}

// CreateDepthTexture creates a depth texture with the specified dimensions and format.
// This is a convenience function for creating depth buffers for render passes.
func (d *Device) CreateDepthTexture(width, height uint32, format gputypes.TextureFormat) *Texture {
	mustInit()

	desc := TextureDescriptor{
		NextInChain:     0,
		Label:           EmptyStringView(),
		Usage:           gputypes.TextureUsageRenderAttachment,
		Dimension:       gputypes.TextureDimension2D,
		Size:            gputypes.Extent3D{Width: width, Height: height, DepthOrArrayLayers: 1},
		Format:          format,
		MipLevelCount:   1,
		SampleCount:     1,
		ViewFormatCount: 0,
		ViewFormats:     0,
	}

	return d.CreateTexture(&desc)
}
