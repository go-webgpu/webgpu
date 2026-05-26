package wgpu

import (
	"context"
	"runtime"
	"unsafe"
)

// MapPending represents an in-flight buffer map request.
// Created by [Buffer.MapAsync]; poll Status() or call Wait() to resolve.
//
// The caller must not access the mapped buffer data until Status() returns
// ready=true with err=nil, or Wait() returns nil.
type MapPending struct {
	req  *mapRequest
	done bool
	err  error
}

// Status reports whether the map request has completed.
// Non-blocking — returns (false, nil) if still pending.
// Once it returns (true, ...), subsequent calls return the same value.
func (p *MapPending) Status() (ready bool, err error) {
	if p == nil {
		return true, nil
	}
	if p.done {
		return true, p.err
	}
	select {
	case <-p.req.done:
		p.done = true
		if p.req.status != MapAsyncStatusSuccess {
			msg := p.req.message
			if msg == "" {
				msg = "buffer map failed"
			}
			p.err = &WGPUError{Op: "Buffer.MapAsync", Message: msg}
		}
		return true, p.err
	default:
		return false, nil
	}
}

// Wait blocks until the map request completes or ctx is canceled.
// Returns nil on success, ctx.Err() if context was canceled before completion.
func (p *MapPending) Wait(ctx context.Context) error {
	if p == nil {
		return nil
	}
	if p.done {
		return p.err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-p.req.done:
		p.done = true
		if p.req.status != MapAsyncStatusSuccess {
			msg := p.req.message
			if msg == "" {
				msg = "buffer map failed"
			}
			p.err = &WGPUError{Op: "Buffer.MapAsync", Message: msg}
		}
		return p.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release discards the pending handle. Safe to call after Wait/Status resolved.
func (p *MapPending) Release() {}

// mapAsyncStart issues wgpuBufferMapAsync and returns the mapRequest.
// Shared by MapAsync (non-blocking) and Map (blocking with poll loop).
func (b *Buffer) mapAsyncStart(mode MapMode, offset, size uint64) (*mapRequest, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if b == nil || b.handle == 0 {
		return nil, &WGPUError{Op: "Buffer.MapAsync", Message: "buffer is nil or released"}
	}

	mapCallbackOnce.Do(initMapCallback)

	req := &mapRequest{
		done: make(chan struct{}),
	}

	mapRequestsMu.Lock()
	mapRequestID++
	reqID := mapRequestID
	mapRequests[reqID] = req
	mapRequestsMu.Unlock()

	callbackInfo := BufferMapCallbackInfo{
		NextInChain: 0,
		Mode:        CallbackModeAllowProcessEvents,
		Callback:    mapCallbackPtr,
		Userdata1:   reqID,
		Userdata2:   0,
	}

	procBufferMapAsync.Call( //nolint:errcheck
		b.handle,
		uintptr(mode),
		uintptr(offset),
		uintptr(size),
		uintptr(unsafe.Pointer(&callbackInfo)),
	)

	return req, nil
}

// MapAsync initiates an asynchronous buffer map without blocking.
// Returns a *MapPending that resolves once the GPU completes the operation.
//
// The caller must periodically drive Device.Poll(false) so the mapping resolves.
// For a blocking variant use [Buffer.Map].
//
// Matches gogpu/wgpu Buffer.MapAsync(mode, offset, size) (*MapPending, error).
func (b *Buffer) MapAsync(mode MapMode, offset, size uint64) (*MapPending, error) {
	req, err := b.mapAsyncStart(mode, offset, size)
	if err != nil {
		return nil, err
	}
	return &MapPending{req: req}, nil
}

// Map blocks until a CPU-visible mapping is established for the given byte
// range, or until ctx is canceled.
//
// The buffer must have been created with BufferUsageMapRead or
// BufferUsageMapWrite matching mode. offset must be a multiple of 8 and
// size must be a multiple of 4 (WebGPU MAP_ALIGNMENT).
//
// After Map succeeds, call MappedRange to obtain a byte view and Unmap
// when finished:
//
//	if err := buf.Map(ctx, wgpu.MapModeRead, 0, size); err != nil {
//	    return err
//	}
//	defer buf.Unmap()
//	rng, _ := buf.MappedRange(0, size)
//	data := rng.Bytes()
//
// Matches gogpu/wgpu Buffer.Map(ctx, mode, offset, size) error.
func (b *Buffer) Map(ctx context.Context, mode MapMode, offset, size uint64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	req, err := b.mapAsyncStart(mode, offset, size)
	if err != nil {
		return err
	}

	// Resolve device: prefer b.device, fall back to nil poll path.
	dev := b.device

	// Kick an initial synchronous poll — for immediate-complete backends.
	if dev != nil {
		dev.Poll(false)
	}

	// Check if already done (fast path for synchronous backends).
	select {
	case <-req.done:
		if req.status != MapAsyncStatusSuccess {
			msg := req.message
			if msg == "" {
				msg = "buffer map failed"
			}
			return &WGPUError{Op: "Buffer.Map", Message: msg}
		}
		return nil
	default:
	}

	// Start a polling goroutine so the mapping resolves even when the
	// caller does not drive Poll itself. This matches the gogpu/wgpu pattern.
	if dev != nil {
		go func() {
			for {
				select {
				case <-req.done:
					return
				default:
					dev.Poll(false)
					runtime.Gosched()
				}
			}
		}()
	}

	// Wait for completion or context cancellation.
	select {
	case <-req.done:
		if req.status != MapAsyncStatusSuccess {
			msg := req.message
			if msg == "" {
				msg = "buffer map failed"
			}
			return &WGPUError{Op: "Buffer.Map", Message: msg}
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
