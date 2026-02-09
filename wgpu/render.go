package wgpu

import (
	"math"
	"unsafe"

	"github.com/gogpu/gputypes"
)

// DepthSliceUndefined is used for 2D textures in color attachments.
const DepthSliceUndefined uint32 = 0xFFFFFFFF

// TimestampLocationUndefined indicates no timestamp write at this location.
const TimestampLocationUndefined uint32 = 0xFFFFFFFF

// Color represents RGBA color with double precision.
type Color struct {
	R, G, B, A float64
}

// renderPassColorAttachment is the native structure for color attachments.
// Uses uint32 for LoadOp/StoreOp with wgpu-native converted values.
type renderPassColorAttachment struct {
	nextInChain   uintptr // 8 bytes
	view          uintptr // 8 bytes (WGPUTextureView)
	depthSlice    uint32  // 4 bytes - MUST be DepthSliceUndefined for 2D!
	_pad1         [4]byte // 4 bytes padding
	resolveTarget uintptr // 8 bytes (WGPUTextureView, nullable)
	loadOp        uint32  // 4 bytes - wgpu-native converted value
	storeOp       uint32  // 4 bytes - wgpu-native converted value
	clearValue    Color   // 32 bytes (4 * float64)
}

// renderPassDescriptor is the native structure for render pass descriptor.
type renderPassDescriptor struct {
	nextInChain            uintptr    // 8 bytes
	label                  StringView // 16 bytes
	colorAttachmentCount   uintptr    // 8 bytes (size_t)
	colorAttachments       uintptr    // 8 bytes (pointer to array)
	depthStencilAttachment uintptr    // 8 bytes (nullable)
	occlusionQuerySet      uintptr    // 8 bytes (nullable)
	timestampWrites        uintptr    // 8 bytes (nullable)
}

// RenderPassColorAttachment describes a color attachment for a render pass.
type RenderPassColorAttachment struct {
	View          *TextureView
	ResolveTarget *TextureView // For MSAA, nil otherwise
	LoadOp        gputypes.LoadOp
	StoreOp       gputypes.StoreOp
	ClearValue    Color
}

// RenderPassDepthStencilAttachment describes a depth/stencil attachment (user API).
type RenderPassDepthStencilAttachment struct {
	View              *TextureView
	DepthLoadOp       gputypes.LoadOp
	DepthStoreOp      gputypes.StoreOp
	DepthClearValue   float32
	DepthReadOnly     bool
	StencilLoadOp     gputypes.LoadOp
	StencilStoreOp    gputypes.StoreOp
	StencilClearValue uint32
	StencilReadOnly   bool
}

// renderPassDepthStencilAttachment is the native structure (40 bytes).
// Uses uint32 for LoadOp/StoreOp with wgpu-native converted values.
type renderPassDepthStencilAttachment struct {
	view              uintptr
	depthLoadOp       uint32 // wgpu-native converted value
	depthStoreOp      uint32 // wgpu-native converted value
	depthClearValue   float32
	depthReadOnly     Bool
	stencilLoadOp     uint32 // wgpu-native converted value
	stencilStoreOp    uint32 // wgpu-native converted value
	stencilClearValue uint32
	stencilReadOnly   Bool
}

// renderPassTimestampWrites is the native structure for timestamp writes (24 bytes).
type renderPassTimestampWrites struct {
	querySet                  uintptr // 8 bytes (WGPUQuerySet)
	beginningOfPassWriteIndex uint32  // 4 bytes
	endOfPassWriteIndex       uint32  // 4 bytes
	_pad                      [8]byte // 8 bytes padding for alignment
}

// RenderPassTimestampWrites describes timestamp writes for a render pass.
type RenderPassTimestampWrites struct {
	QuerySet                  *QuerySet
	BeginningOfPassWriteIndex uint32 // Use TimestampLocationUndefined to disable
	EndOfPassWriteIndex       uint32 // Use TimestampLocationUndefined to disable
}

// RenderPassDescriptor describes a render pass.
type RenderPassDescriptor struct {
	Label                  string
	ColorAttachments       []RenderPassColorAttachment
	DepthStencilAttachment *RenderPassDepthStencilAttachment
	TimestampWrites        *RenderPassTimestampWrites
}

// BeginRenderPass begins a render pass.
func (enc *CommandEncoder) BeginRenderPass(desc *RenderPassDescriptor) *RenderPassEncoder {
	mustInit()

	if desc == nil || len(desc.ColorAttachments) == 0 {
		return nil
	}

	// Build native color attachments
	nativeColorAttachments := make([]renderPassColorAttachment, len(desc.ColorAttachments))
	for i, ca := range desc.ColorAttachments {
		var viewHandle uintptr
		if ca.View != nil {
			viewHandle = ca.View.handle
		}

		var resolveHandle uintptr
		if ca.ResolveTarget != nil {
			resolveHandle = ca.ResolveTarget.handle
		}

		nativeColorAttachments[i] = renderPassColorAttachment{
			nextInChain:   0,
			view:          viewHandle,
			depthSlice:    DepthSliceUndefined, // CRITICAL for 2D textures!
			resolveTarget: resolveHandle,
			loadOp:        toWGPULoadOp(ca.LoadOp),
			storeOp:       toWGPUStoreOp(ca.StoreOp),
			clearValue:    ca.ClearValue,
		}
	}

	// Build depth/stencil attachment if present
	var depthStencilPtr uintptr
	var nativeDepthStencil renderPassDepthStencilAttachment
	if desc.DepthStencilAttachment != nil {
		depthRO := False
		if desc.DepthStencilAttachment.DepthReadOnly {
			depthRO = True
		}
		stencilRO := False
		if desc.DepthStencilAttachment.StencilReadOnly {
			stencilRO = True
		}

		nativeDepthStencil = renderPassDepthStencilAttachment{
			view:              desc.DepthStencilAttachment.View.handle,
			depthLoadOp:       toWGPULoadOp(desc.DepthStencilAttachment.DepthLoadOp),
			depthStoreOp:      toWGPUStoreOp(desc.DepthStencilAttachment.DepthStoreOp),
			depthClearValue:   desc.DepthStencilAttachment.DepthClearValue,
			depthReadOnly:     depthRO,
			stencilLoadOp:     toWGPULoadOp(desc.DepthStencilAttachment.StencilLoadOp),
			stencilStoreOp:    toWGPUStoreOp(desc.DepthStencilAttachment.StencilStoreOp),
			stencilClearValue: desc.DepthStencilAttachment.StencilClearValue,
			stencilReadOnly:   stencilRO,
		}
		depthStencilPtr = uintptr(unsafe.Pointer(&nativeDepthStencil))
	}

	// Build timestamp writes if present
	var timestampWritesPtr uintptr
	var nativeTimestampWrites renderPassTimestampWrites
	if desc.TimestampWrites != nil {
		nativeTimestampWrites = renderPassTimestampWrites{
			querySet:                  desc.TimestampWrites.QuerySet.handle,
			beginningOfPassWriteIndex: desc.TimestampWrites.BeginningOfPassWriteIndex,
			endOfPassWriteIndex:       desc.TimestampWrites.EndOfPassWriteIndex,
		}
		timestampWritesPtr = uintptr(unsafe.Pointer(&nativeTimestampWrites))
	}

	nativeDesc := renderPassDescriptor{
		nextInChain:            0,
		label:                  EmptyStringView(),
		colorAttachmentCount:   uintptr(len(nativeColorAttachments)),
		colorAttachments:       uintptr(unsafe.Pointer(&nativeColorAttachments[0])),
		depthStencilAttachment: depthStencilPtr,
		occlusionQuerySet:      0,
		timestampWrites:        timestampWritesPtr,
	}

	handle, _, _ := procCommandEncoderBeginRenderPass.Call(
		enc.handle,
		uintptr(unsafe.Pointer(&nativeDesc)),
	)
	if handle == 0 {
		return nil
	}
	trackResource(handle, "RenderPassEncoder")
	return &RenderPassEncoder{handle: handle}
}

// SetPipeline sets the render pipeline for this pass.
func (rpe *RenderPassEncoder) SetPipeline(pipeline *RenderPipeline) {
	mustInit()
	procRenderPassEncoderSetPipeline.Call(rpe.handle, pipeline.handle) //nolint:errcheck
}

// SetBindGroup sets a bind group for this pass.
func (rpe *RenderPassEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	mustInit()

	var offsetsPtr uintptr
	offsetCount := uintptr(0)
	if len(dynamicOffsets) > 0 {
		offsetsPtr = uintptr(unsafe.Pointer(&dynamicOffsets[0]))
		offsetCount = uintptr(len(dynamicOffsets))
	}

	procRenderPassEncoderSetBindGroup.Call( //nolint:errcheck
		rpe.handle,
		uintptr(groupIndex),
		group.handle,
		offsetCount,
		offsetsPtr,
	)
}

// SetVertexBuffer sets a vertex buffer for this pass.
func (rpe *RenderPassEncoder) SetVertexBuffer(slot uint32, buffer *Buffer, offset, size uint64) {
	mustInit()
	procRenderPassEncoderSetVertexBuffer.Call( //nolint:errcheck
		rpe.handle,
		uintptr(slot),
		buffer.handle,
		uintptr(offset),
		uintptr(size),
	)
}

// SetIndexBuffer sets the index buffer for this pass.
func (rpe *RenderPassEncoder) SetIndexBuffer(buffer *Buffer, format gputypes.IndexFormat, offset, size uint64) {
	mustInit()
	procRenderPassEncoderSetIndexBuffer.Call( //nolint:errcheck
		rpe.handle,
		buffer.handle,
		uintptr(format),
		uintptr(offset),
		uintptr(size),
	)
}

// Draw draws primitives.
func (rpe *RenderPassEncoder) Draw(vertexCount, instanceCount, firstVertex, firstInstance uint32) {
	mustInit()
	procRenderPassEncoderDraw.Call( //nolint:errcheck
		rpe.handle,
		uintptr(vertexCount),
		uintptr(instanceCount),
		uintptr(firstVertex),
		uintptr(firstInstance),
	)
}

// DrawIndexed draws indexed primitives.
func (rpe *RenderPassEncoder) DrawIndexed(indexCount, instanceCount, firstIndex uint32, baseVertex int32, firstInstance uint32) {
	mustInit()
	procRenderPassEncoderDrawIndexed.Call( //nolint:errcheck
		rpe.handle,
		uintptr(indexCount),
		uintptr(instanceCount),
		uintptr(firstIndex),
		uintptr(baseVertex),
		uintptr(firstInstance),
	)
}

// DrawIndirect draws primitives using parameters from a GPU buffer.
// indirectBuffer must contain a DrawIndirectArgs structure:
//   - vertexCount (uint32)
//   - instanceCount (uint32)
//   - firstVertex (uint32)
//   - firstInstance (uint32)
func (rpe *RenderPassEncoder) DrawIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	mustInit()
	procRenderPassEncoderDrawIndirect.Call( //nolint:errcheck
		rpe.handle,
		indirectBuffer.handle,
		uintptr(indirectOffset),
	)
}

// DrawIndexedIndirect draws indexed primitives using parameters from a GPU buffer.
// indirectBuffer must contain a DrawIndexedIndirectArgs structure:
//   - indexCount (uint32)
//   - instanceCount (uint32)
//   - firstIndex (uint32)
//   - baseVertex (int32)
//   - firstInstance (uint32)
func (rpe *RenderPassEncoder) DrawIndexedIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	mustInit()
	procRenderPassEncoderDrawIndexedIndirect.Call( //nolint:errcheck
		rpe.handle,
		indirectBuffer.handle,
		uintptr(indirectOffset),
	)
}

// SetViewport sets the viewport used during the rasterization stage.
// x, y: top-left corner of the viewport in pixels
// width, height: dimensions of the viewport in pixels
// minDepth, maxDepth: depth range for the viewport (typically 0.0 to 1.0)
func (rpe *RenderPassEncoder) SetViewport(x, y, width, height, minDepth, maxDepth float32) {
	mustInit()
	procRenderPassEncoderSetViewport.Call( //nolint:errcheck
		rpe.handle,
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(math.Float32bits(minDepth)),
		uintptr(math.Float32bits(maxDepth)),
	)
}

// SetScissorRect sets the scissor rectangle used during the rasterization stage.
// Pixels outside the scissor rectangle will be discarded.
// x, y: top-left corner of the scissor rectangle in pixels
// width, height: dimensions of the scissor rectangle in pixels
func (rpe *RenderPassEncoder) SetScissorRect(x, y, width, height uint32) {
	mustInit()
	procRenderPassEncoderSetScissorRect.Call( //nolint:errcheck
		rpe.handle,
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
	)
}

// SetBlendConstant sets the blend constant color used by blend operations.
// Errors are reported via Device error scopes.
func (rpe *RenderPassEncoder) SetBlendConstant(color *Color) {
	mustInit()
	if color == nil {
		return
	}
	procRenderPassEncoderSetBlendConstant.Call( //nolint:errcheck
		rpe.handle,
		uintptr(unsafe.Pointer(color)),
	)
}

// SetStencilReference sets the stencil reference value used by stencil operations.
func (rpe *RenderPassEncoder) SetStencilReference(reference uint32) {
	mustInit()
	procRenderPassEncoderSetStencilReference.Call( //nolint:errcheck
		rpe.handle,
		uintptr(reference),
	)
}

// InsertDebugMarker inserts a single debug marker label into the render pass.
// This is useful for GPU debugging tools to identify specific command points.
func (rpe *RenderPassEncoder) InsertDebugMarker(markerLabel string) {
	mustInit()
	if rpe == nil || rpe.handle == 0 {
		return
	}
	labelBytes := []byte(markerLabel)
	if len(labelBytes) == 0 {
		return
	}
	label := StringView{
		Data:   uintptr(unsafe.Pointer(&labelBytes[0])),
		Length: uintptr(len(labelBytes)),
	}
	procRenderPassEncoderInsertDebugMarker.Call( //nolint:errcheck
		rpe.handle,
		uintptr(unsafe.Pointer(&label)),
	)
}

// PushDebugGroup begins a labeled debug group in the render pass.
// Use PopDebugGroup to end the group. Groups can be nested.
func (rpe *RenderPassEncoder) PushDebugGroup(groupLabel string) {
	mustInit()
	if rpe == nil || rpe.handle == 0 {
		return
	}
	labelBytes := []byte(groupLabel)
	if len(labelBytes) == 0 {
		return
	}
	label := StringView{
		Data:   uintptr(unsafe.Pointer(&labelBytes[0])),
		Length: uintptr(len(labelBytes)),
	}
	procRenderPassEncoderPushDebugGroup.Call( //nolint:errcheck
		rpe.handle,
		uintptr(unsafe.Pointer(&label)),
	)
}

// PopDebugGroup ends the current debug group in the render pass.
// Must match a preceding PushDebugGroup call.
func (rpe *RenderPassEncoder) PopDebugGroup() {
	mustInit()
	if rpe == nil || rpe.handle == 0 {
		return
	}
	procRenderPassEncoderPopDebugGroup.Call(rpe.handle) //nolint:errcheck
}

// End ends the render pass.
func (rpe *RenderPassEncoder) End() {
	mustInit()
	procRenderPassEncoderEnd.Call(rpe.handle) //nolint:errcheck
}

// Release releases the render pass encoder.
func (rpe *RenderPassEncoder) Release() {
	if rpe.handle != 0 {
		untrackResource(rpe.handle)
		procRenderPassEncoderRelease.Call(rpe.handle) //nolint:errcheck
		rpe.handle = 0
	}
}

// Handle returns the underlying handle. For advanced use only.
func (rpe *RenderPassEncoder) Handle() uintptr { return rpe.handle }
