package wgpu

import (
	"unsafe"

	"github.com/gogpu/gputypes"
)

// CommandEncoderDescriptor describes a command encoder.
type CommandEncoderDescriptor struct {
	NextInChain uintptr // *ChainedStruct
	Label       StringView
}

// CommandBufferDescriptor describes a command buffer.
type CommandBufferDescriptor struct {
	NextInChain uintptr // *ChainedStruct
	Label       StringView
}

// ComputePassDescriptor describes a compute pass.
type ComputePassDescriptor struct {
	NextInChain     uintptr // *ChainedStruct
	Label           StringView
	TimestampWrites uintptr // *ComputePassTimestampWrites (nullable)
}

// CreateCommandEncoder creates a command encoder.
func (d *Device) CreateCommandEncoder(desc *CommandEncoderDescriptor) *CommandEncoder {
	mustInit()
	var descPtr uintptr
	if desc != nil {
		descPtr = uintptr(unsafe.Pointer(desc))
	}
	handle, _, _ := procDeviceCreateCommandEncoder.Call(
		d.handle,
		descPtr,
	)
	if handle == 0 {
		return nil
	}
	trackResource(handle, "CommandEncoder")
	return &CommandEncoder{handle: handle}
}

// BeginComputePass begins a compute pass.
func (enc *CommandEncoder) BeginComputePass(desc *ComputePassDescriptor) *ComputePassEncoder {
	mustInit()
	var descPtr uintptr
	if desc != nil {
		descPtr = uintptr(unsafe.Pointer(desc))
	}
	handle, _, _ := procCommandEncoderBeginComputePass.Call(
		enc.handle,
		descPtr,
	)
	if handle == 0 {
		return nil
	}
	trackResource(handle, "ComputePassEncoder")
	return &ComputePassEncoder{handle: handle}
}

// CopyBufferToBuffer copies data between buffers.
func (enc *CommandEncoder) CopyBufferToBuffer(src *Buffer, srcOffset uint64, dst *Buffer, dstOffset uint64, size uint64) {
	mustInit()
	procCommandEncoderCopyBufferToBuffer.Call( //nolint:errcheck
		enc.handle,
		src.handle,
		uintptr(srcOffset),
		dst.handle,
		uintptr(dstOffset),
		uintptr(size),
	)
}

// ClearBuffer clears a region of a buffer to zeros.
// size = 0 means clear from offset to end of buffer.
func (enc *CommandEncoder) ClearBuffer(buffer *Buffer, offset, size uint64) {
	mustInit()
	if buffer == nil {
		return
	}
	procCommandEncoderClearBuffer.Call( //nolint:errcheck
		enc.handle,
		buffer.handle,
		uintptr(offset),
		uintptr(size),
	)
}

// InsertDebugMarker inserts a single debug marker label.
// This is useful for GPU debugging tools to identify specific command points.
func (enc *CommandEncoder) InsertDebugMarker(markerLabel string) {
	mustInit()
	labelBytes := []byte(markerLabel)
	if len(labelBytes) == 0 {
		return
	}
	label := StringView{
		Data:   uintptr(unsafe.Pointer(&labelBytes[0])),
		Length: uintptr(len(labelBytes)),
	}
	procCommandEncoderInsertDebugMarker.Call( //nolint:errcheck
		enc.handle,
		uintptr(unsafe.Pointer(&label)),
	)
}

// PushDebugGroup begins a labeled debug group.
// Use PopDebugGroup to end the group. Groups can be nested.
func (enc *CommandEncoder) PushDebugGroup(groupLabel string) {
	mustInit()
	labelBytes := []byte(groupLabel)
	if len(labelBytes) == 0 {
		return
	}
	label := StringView{
		Data:   uintptr(unsafe.Pointer(&labelBytes[0])),
		Length: uintptr(len(labelBytes)),
	}
	procCommandEncoderPushDebugGroup.Call( //nolint:errcheck
		enc.handle,
		uintptr(unsafe.Pointer(&label)),
	)
}

// PopDebugGroup ends the current debug group.
// Must match a preceding PushDebugGroup call.
func (enc *CommandEncoder) PopDebugGroup() {
	mustInit()
	procCommandEncoderPopDebugGroup.Call(enc.handle) //nolint:errcheck
}

// CopyBufferToTexture copies data from a buffer to a texture.
// Errors are reported via Device error scopes, not as return values.
func (enc *CommandEncoder) CopyBufferToTexture(source *TexelCopyBufferInfo, destination *TexelCopyTextureInfo, copySize *gputypes.Extent3D) {
	mustInit()
	if source == nil || destination == nil || copySize == nil {
		return
	}
	procCommandEncoderCopyBufferToTexture.Call( //nolint:errcheck
		enc.handle,
		uintptr(unsafe.Pointer(source)),
		uintptr(unsafe.Pointer(destination)),
		uintptr(unsafe.Pointer(copySize)),
	)
}

// CopyTextureToBuffer copies data from a texture to a buffer.
// Errors are reported via Device error scopes, not as return values.
func (enc *CommandEncoder) CopyTextureToBuffer(source *TexelCopyTextureInfo, destination *TexelCopyBufferInfo, copySize *gputypes.Extent3D) {
	mustInit()
	if source == nil || destination == nil || copySize == nil {
		return
	}
	procCommandEncoderCopyTextureToBuffer.Call( //nolint:errcheck
		enc.handle,
		uintptr(unsafe.Pointer(source)),
		uintptr(unsafe.Pointer(destination)),
		uintptr(unsafe.Pointer(copySize)),
	)
}

// CopyTextureToTexture copies data from one texture to another.
// Errors are reported via Device error scopes, not as return values.
func (enc *CommandEncoder) CopyTextureToTexture(source *TexelCopyTextureInfo, destination *TexelCopyTextureInfo, copySize *gputypes.Extent3D) {
	mustInit()
	if source == nil || destination == nil || copySize == nil {
		return
	}
	procCommandEncoderCopyTextureToTexture.Call( //nolint:errcheck
		enc.handle,
		uintptr(unsafe.Pointer(source)),
		uintptr(unsafe.Pointer(destination)),
		uintptr(unsafe.Pointer(copySize)),
	)
}

// Finish finishes recording and returns a command buffer.
func (enc *CommandEncoder) Finish(desc *CommandBufferDescriptor) *CommandBuffer {
	mustInit()
	var descPtr uintptr
	if desc != nil {
		descPtr = uintptr(unsafe.Pointer(desc))
	}
	handle, _, _ := procCommandEncoderFinish.Call(
		enc.handle,
		descPtr,
	)
	if handle == 0 {
		return nil
	}
	trackResource(handle, "CommandBuffer")
	return &CommandBuffer{handle: handle}
}

// Release releases the command encoder.
func (enc *CommandEncoder) Release() {
	if enc.handle != 0 {
		untrackResource(enc.handle)
		procCommandEncoderRelease.Call(enc.handle) //nolint:errcheck
		enc.handle = 0
	}
}

// WriteTimestamp writes a timestamp to a query.
// Note: This is a wgpu-native extension. Prefer pass-level timestamps
// via RenderPassTimestampWrites or ComputePassTimestampWrites when possible.
func (enc *CommandEncoder) WriteTimestamp(querySet *QuerySet, queryIndex uint32) {
	mustInit()
	procCommandEncoderWriteTimestamp.Call( //nolint:errcheck
		enc.handle,
		querySet.handle,
		uintptr(queryIndex),
	)
}

// ResolveQuerySet resolves query results to a buffer.
// The buffer must have BufferUsageQueryResolve usage.
func (enc *CommandEncoder) ResolveQuerySet(querySet *QuerySet, firstQuery, queryCount uint32, destination *Buffer, destinationOffset uint64) {
	mustInit()
	procCommandEncoderResolveQuerySet.Call( //nolint:errcheck
		enc.handle,
		querySet.handle,
		uintptr(firstQuery),
		uintptr(queryCount),
		destination.handle,
		uintptr(destinationOffset),
	)
}

// Handle returns the underlying handle.
func (enc *CommandEncoder) Handle() uintptr { return enc.handle }

// SetPipeline sets the compute pipeline.
func (cpe *ComputePassEncoder) SetPipeline(pipeline *ComputePipeline) {
	mustInit()
	procComputePassEncoderSetPipeline.Call( //nolint:errcheck
		cpe.handle,
		pipeline.handle,
	)
}

// SetBindGroup sets a bind group.
func (cpe *ComputePassEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	mustInit()
	var offsetsPtr uintptr
	offsetCount := uintptr(0)
	if len(dynamicOffsets) > 0 {
		offsetsPtr = uintptr(unsafe.Pointer(&dynamicOffsets[0]))
		offsetCount = uintptr(len(dynamicOffsets))
	}
	procComputePassEncoderSetBindGroup.Call( //nolint:errcheck
		cpe.handle,
		uintptr(groupIndex),
		group.handle,
		offsetCount,
		offsetsPtr,
	)
}

// DispatchWorkgroups dispatches compute work.
func (cpe *ComputePassEncoder) DispatchWorkgroups(x, y, z uint32) {
	mustInit()
	procComputePassEncoderDispatchWorkgroups.Call( //nolint:errcheck
		cpe.handle,
		uintptr(x),
		uintptr(y),
		uintptr(z),
	)
}

// DispatchWorkgroupsIndirect dispatches compute work using parameters from a GPU buffer.
// indirectBuffer must contain a DispatchIndirectArgs structure:
//   - workgroupCountX (uint32)
//   - workgroupCountY (uint32)
//   - workgroupCountZ (uint32)
func (cpe *ComputePassEncoder) DispatchWorkgroupsIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	mustInit()
	procComputePassEncoderDispatchWorkgroupsIndirect.Call( //nolint:errcheck
		cpe.handle,
		indirectBuffer.handle,
		uintptr(indirectOffset),
	)
}

// End ends the compute pass.
func (cpe *ComputePassEncoder) End() {
	mustInit()
	procComputePassEncoderEnd.Call(cpe.handle) //nolint:errcheck
}

// Release releases the compute pass encoder.
func (cpe *ComputePassEncoder) Release() {
	if cpe.handle != 0 {
		untrackResource(cpe.handle)
		procComputePassEncoderRelease.Call(cpe.handle) //nolint:errcheck
		cpe.handle = 0
	}
}

// Handle returns the underlying handle.
func (cpe *ComputePassEncoder) Handle() uintptr { return cpe.handle }

// Submit submits command buffers for execution.
func (q *Queue) Submit(commands ...*CommandBuffer) {
	mustInit()
	if len(commands) == 0 {
		return
	}
	handles := make([]uintptr, len(commands))
	for i, cmd := range commands {
		handles[i] = cmd.handle
	}
	procQueueSubmit.Call( //nolint:errcheck
		q.handle,
		uintptr(len(handles)),
		uintptr(unsafe.Pointer(&handles[0])),
	)
}

// Release releases the command buffer.
func (cb *CommandBuffer) Release() {
	if cb.handle != 0 {
		untrackResource(cb.handle)
		procCommandBufferRelease.Call(cb.handle) //nolint:errcheck
		cb.handle = 0
	}
}

// Handle returns the underlying handle.
func (cb *CommandBuffer) Handle() uintptr { return cb.handle }
