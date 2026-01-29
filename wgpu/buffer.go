package wgpu

import (
	"errors"
	"sync"
	"unsafe"

	"github.com/go-webgpu/goffi/ffi"
	"github.com/gogpu/gputypes"
)

// MapMode specifies the mapping mode for MapAsync.
type MapMode uint64

const (
	MapModeNone  MapMode = 0x0000000000000000
	MapModeRead  MapMode = 0x0000000000000001
	MapModeWrite MapMode = 0x0000000000000002
)

// MapAsyncStatus is the status returned by MapAsync callback.
type MapAsyncStatus uint32

const (
	MapAsyncStatusSuccess         MapAsyncStatus = 0x00000001
	MapAsyncStatusInstanceDropped MapAsyncStatus = 0x00000002
	MapAsyncStatusError           MapAsyncStatus = 0x00000003
	MapAsyncStatusAborted         MapAsyncStatus = 0x00000004
	MapAsyncStatusUnknown         MapAsyncStatus = 0x00000005
)

// BufferMapCallbackInfo holds callback configuration for MapAsync.
type BufferMapCallbackInfo struct {
	NextInChain uintptr // *ChainedStruct
	Mode        CallbackMode
	Callback    uintptr // Function pointer
	Userdata1   uintptr
	Userdata2   uintptr
}

// mapRequest holds state for an async map request.
type mapRequest struct {
	done    chan struct{}
	status  MapAsyncStatus
	message string
}

var (
	// Global registry for pending map requests
	mapRequests   = make(map[uintptr]*mapRequest)
	mapRequestsMu sync.Mutex
	mapRequestID  uintptr

	// Callback function pointer (created once)
	mapCallbackPtr  uintptr
	mapCallbackOnce sync.Once
)

// mapCallbackHandler is the Go function called by C code via ffi.NewCallback.
// Signature: void(status uint32, message *StringView, userdata1 uintptr, userdata2 uintptr)
func mapCallbackHandler(status uintptr, message uintptr, userdata1, userdata2 uintptr) uintptr {
	// Extract message string
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
	mapRequestsMu.Lock()
	req, ok := mapRequests[userdata1]
	if ok {
		delete(mapRequests, userdata1)
	}
	mapRequestsMu.Unlock()

	if ok && req != nil {
		req.status = MapAsyncStatus(status)
		req.message = msg
		close(req.done)
	}
	return 0
}

// initMapCallback creates the C callback function pointer using goffi.
func initMapCallback() {
	mapCallbackPtr = ffi.NewCallback(mapCallbackHandler)
}

// BufferDescriptor describes a buffer to create.
type BufferDescriptor struct {
	NextInChain      uintptr              // *ChainedStruct
	Label            StringView           // Buffer label for debugging
	Usage            gputypes.BufferUsage // How the buffer will be used
	Size             uint64               // Size in bytes
	MappedAtCreation Bool                 // If true, buffer is mapped when created
}

// CreateBuffer creates a new GPU buffer.
func (d *Device) CreateBuffer(desc *BufferDescriptor) *Buffer {
	mustInit()
	if desc == nil {
		return nil
	}
	handle, _, _ := procDeviceCreateBuffer.Call(
		d.handle,
		uintptr(unsafe.Pointer(desc)),
	)
	if handle == 0 {
		return nil
	}
	return &Buffer{handle: handle}
}

// GetMappedRange returns a pointer to the mapped buffer data.
// The buffer must be mapped (either via MapAsync or MappedAtCreation).
// offset and size specify the range to access.
// Returns nil if the buffer is not mapped or the range is invalid.
func (b *Buffer) GetMappedRange(offset, size uint64) unsafe.Pointer {
	mustInit()
	ptr, _, _ := procBufferGetMappedRange.Call(
		b.handle,
		uintptr(offset),
		uintptr(size),
	)
	if ptr == 0 {
		return nil
	}
	// nolint:govet // ptr is uintptr from FFI call - returned immediately, GC safe
	return unsafe.Pointer(ptr)
}

// Unmap unmaps the buffer, making the mapped memory inaccessible.
// For buffers created with MappedAtCreation, this commits the data to the GPU.
func (b *Buffer) Unmap() {
	mustInit()
	procBufferUnmap.Call(b.handle) //nolint:errcheck
}

// GetSize returns the size of the buffer in bytes.
func (b *Buffer) GetSize() uint64 {
	mustInit()
	size, _, _ := procBufferGetSize.Call(b.handle)
	return uint64(size)
}

// MapAsync maps a buffer for reading or writing.
// This is a synchronous wrapper that blocks until the mapping is complete.
// The device parameter is used to poll for completion.
// After MapAsync succeeds, use GetMappedRange to access the data.
// Call Unmap when done to release the mapping.
func (b *Buffer) MapAsync(device *Device, mode MapMode, offset, size uint64) error {
	mustInit()

	// Initialize callback once
	mapCallbackOnce.Do(initMapCallback)

	// Create request state
	req := &mapRequest{
		done: make(chan struct{}),
	}

	// Register request
	mapRequestsMu.Lock()
	mapRequestID++
	reqID := mapRequestID
	mapRequests[reqID] = req
	mapRequestsMu.Unlock()

	// Prepare callback info
	callbackInfo := BufferMapCallbackInfo{
		NextInChain: 0,
		Mode:        CallbackModeAllowProcessEvents,
		Callback:    mapCallbackPtr,
		Userdata1:   reqID,
		Userdata2:   0,
	}

	// Call wgpuBufferMapAsync
	procBufferMapAsync.Call( //nolint:errcheck
		b.handle,
		uintptr(mode),
		uintptr(offset),
		uintptr(size),
		uintptr(unsafe.Pointer(&callbackInfo)),
	)

	// Poll device until callback fires
	for {
		select {
		case <-req.done:
			// Callback completed
			if req.status != MapAsyncStatusSuccess {
				msg := req.message
				if msg == "" {
					msg = "buffer map failed"
				}
				return errors.New("wgpu: " + msg)
			}
			return nil
		default:
			// Poll device to process callbacks
			device.Poll(false)
		}
	}
}

// Destroy destroys the buffer, making it invalid.
func (b *Buffer) Destroy() {
	mustInit()
	if b.handle != 0 {
		procBufferDestroy.Call(b.handle) //nolint:errcheck
	}
}

// Release releases the buffer reference.
func (b *Buffer) Release() {
	if b.handle != 0 {
		procBufferRelease.Call(b.handle) //nolint:errcheck
		b.handle = 0
	}
}

// WriteBuffer writes data to a buffer.
// This is a convenience method that stages data for upload to the GPU.
func (q *Queue) WriteBuffer(buffer *Buffer, offset uint64, data []byte) {
	mustInit()
	if len(data) == 0 {
		return
	}
	procQueueWriteBuffer.Call( //nolint:errcheck
		q.handle,
		buffer.handle,
		uintptr(offset),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
	)
}

// WriteBufferTyped writes typed data to a buffer.
// The data pointer should point to the first element, size is total byte size.
func (q *Queue) WriteBufferRaw(buffer *Buffer, offset uint64, data unsafe.Pointer, size uint64) {
	mustInit()
	if size == 0 {
		return
	}
	procQueueWriteBuffer.Call( //nolint:errcheck
		q.handle,
		buffer.handle,
		uintptr(offset),
		uintptr(data),
		uintptr(size),
	)
}

// GetUsage returns the usage flags of this buffer.
func (b *Buffer) GetUsage() gputypes.BufferUsage {
	mustInit()
	if b == nil || b.handle == 0 {
		return gputypes.BufferUsageNone
	}
	usage, _, _ := procBufferGetUsage.Call(b.handle)
	return gputypes.BufferUsage(usage)
}

// GetMapState returns the current mapping state of this buffer.
func (b *Buffer) GetMapState() BufferMapState {
	mustInit()
	if b == nil || b.handle == 0 {
		return BufferMapStateUnmapped
	}
	state, _, _ := procBufferGetMapState.Call(b.handle)
	return BufferMapState(state)
}

// Handle returns the underlying handle. For advanced use only.
func (b *Buffer) Handle() uintptr { return b.handle }
