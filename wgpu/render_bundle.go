package wgpu

import (
	"unsafe"

	"github.com/gogpu/gputypes"
)

// RenderBundleEncoderDescriptor describes a render bundle encoder.
type RenderBundleEncoderDescriptor struct {
	Label              StringView
	ColorFormatCount   uintptr // size_t
	ColorFormats       *gputypes.TextureFormat
	DepthStencilFormat gputypes.TextureFormat
	SampleCount        uint32
	DepthReadOnly      Bool
	StencilReadOnly    Bool
}

// RenderBundleDescriptor describes a render bundle.
type RenderBundleDescriptor struct {
	NextInChain uintptr // *ChainedStruct
	Label       StringView
}

// CreateRenderBundleEncoder creates a render bundle encoder for pre-recording render commands.
// Render bundles allow you to pre-record a sequence of render commands that can be replayed
// multiple times, which is useful for static geometry.
func (d *Device) CreateRenderBundleEncoder(desc *RenderBundleEncoderDescriptor) *RenderBundleEncoder {
	mustInit()

	if desc == nil {
		return nil
	}

	// Build the native descriptor with converted format values
	type nativeDesc struct {
		nextInChain        uintptr
		label              StringView
		colorFormatCount   uintptr
		colorFormats       uintptr
		depthStencilFormat uint32 // converted from gputypes
		sampleCount        uint32
		depthReadOnly      Bool
		stencilReadOnly    Bool
	}

	nd := nativeDesc{
		label:              desc.Label,
		colorFormatCount:   desc.ColorFormatCount,
		depthStencilFormat: toWGPUTextureFormat(desc.DepthStencilFormat),
		sampleCount:        desc.SampleCount,
		depthReadOnly:      desc.DepthReadOnly,
		stencilReadOnly:    desc.StencilReadOnly,
	}

	// Convert color formats to wgpu-native values
	var convertedFormats []uint32
	if desc.ColorFormats != nil && desc.ColorFormatCount > 0 {
		formats := unsafe.Slice(desc.ColorFormats, desc.ColorFormatCount)
		convertedFormats = make([]uint32, desc.ColorFormatCount)
		for i, f := range formats {
			convertedFormats[i] = toWGPUTextureFormat(f)
		}
		nd.colorFormats = uintptr(unsafe.Pointer(&convertedFormats[0]))
	}

	handle, _, _ := procDeviceCreateRenderBundleEncoder.Call(
		d.handle,
		uintptr(unsafe.Pointer(&nd)),
	)
	if handle == 0 {
		return nil
	}
	trackResource(handle, "RenderBundleEncoder")
	return &RenderBundleEncoder{handle: handle}
}

// CreateRenderBundleEncoderSimple creates a render bundle encoder with common settings.
func (d *Device) CreateRenderBundleEncoderSimple(colorFormats []gputypes.TextureFormat, depthFormat gputypes.TextureFormat, sampleCount uint32) *RenderBundleEncoder {
	desc := &RenderBundleEncoderDescriptor{
		ColorFormatCount:   uintptr(len(colorFormats)),
		DepthStencilFormat: depthFormat,
		SampleCount:        sampleCount,
	}
	if len(colorFormats) > 0 {
		desc.ColorFormats = &colorFormats[0]
	}
	return d.CreateRenderBundleEncoder(desc)
}

// SetPipeline sets the render pipeline for subsequent draw calls.
func (rbe *RenderBundleEncoder) SetPipeline(pipeline *RenderPipeline) {
	mustInit()
	procRenderBundleEncoderSetPipeline.Call(rbe.handle, pipeline.handle) //nolint:errcheck
}

// SetBindGroup sets a bind group at the given index.
func (rbe *RenderBundleEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	mustInit()
	var offsetsPtr uintptr
	if len(dynamicOffsets) > 0 {
		offsetsPtr = uintptr(unsafe.Pointer(&dynamicOffsets[0]))
	}
	procRenderBundleEncoderSetBindGroup.Call( //nolint:errcheck
		rbe.handle,
		uintptr(groupIndex),
		group.handle,
		uintptr(len(dynamicOffsets)),
		offsetsPtr,
	)
}

// SetVertexBuffer sets a vertex buffer at the given slot.
func (rbe *RenderBundleEncoder) SetVertexBuffer(slot uint32, buffer *Buffer, offset, size uint64) {
	mustInit()
	procRenderBundleEncoderSetVertexBuffer.Call( //nolint:errcheck
		rbe.handle,
		uintptr(slot),
		buffer.handle,
		uintptr(offset),
		uintptr(size),
	)
}

// SetIndexBuffer sets the index buffer.
func (rbe *RenderBundleEncoder) SetIndexBuffer(buffer *Buffer, format gputypes.IndexFormat, offset, size uint64) {
	mustInit()
	procRenderBundleEncoderSetIndexBuffer.Call( //nolint:errcheck
		rbe.handle,
		buffer.handle,
		uintptr(format),
		uintptr(offset),
		uintptr(size),
	)
}

// Draw records a non-indexed draw call.
func (rbe *RenderBundleEncoder) Draw(vertexCount, instanceCount, firstVertex, firstInstance uint32) {
	mustInit()
	procRenderBundleEncoderDraw.Call( //nolint:errcheck
		rbe.handle,
		uintptr(vertexCount),
		uintptr(instanceCount),
		uintptr(firstVertex),
		uintptr(firstInstance),
	)
}

// DrawIndexed records an indexed draw call.
func (rbe *RenderBundleEncoder) DrawIndexed(indexCount, instanceCount, firstIndex uint32, baseVertex int32, firstInstance uint32) {
	mustInit()
	procRenderBundleEncoderDrawIndexed.Call( //nolint:errcheck
		rbe.handle,
		uintptr(indexCount),
		uintptr(instanceCount),
		uintptr(firstIndex),
		uintptr(baseVertex),
		uintptr(firstInstance),
	)
}

// DrawIndirect records an indirect draw call.
func (rbe *RenderBundleEncoder) DrawIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	mustInit()
	procRenderBundleEncoderDrawIndirect.Call( //nolint:errcheck
		rbe.handle,
		indirectBuffer.handle,
		uintptr(indirectOffset),
	)
}

// DrawIndexedIndirect records an indirect indexed draw call.
func (rbe *RenderBundleEncoder) DrawIndexedIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	mustInit()
	procRenderBundleEncoderDrawIndexedIndirect.Call( //nolint:errcheck
		rbe.handle,
		indirectBuffer.handle,
		uintptr(indirectOffset),
	)
}

// Finish completes recording and returns the render bundle.
func (rbe *RenderBundleEncoder) Finish(desc *RenderBundleDescriptor) *RenderBundle {
	mustInit()

	var descPtr uintptr
	if desc != nil {
		descPtr = uintptr(unsafe.Pointer(desc))
	}

	handle, _, _ := procRenderBundleEncoderFinish.Call(rbe.handle, descPtr)
	if handle == 0 {
		return nil
	}
	trackResource(handle, "RenderBundle")
	return &RenderBundle{handle: handle}
}

// Release releases the render bundle encoder.
func (rbe *RenderBundleEncoder) Release() {
	if rbe.handle != 0 {
		untrackResource(rbe.handle)
		procRenderBundleEncoderRelease.Call(rbe.handle) //nolint:errcheck
		rbe.handle = 0
	}
}

// Handle returns the underlying handle.
func (rbe *RenderBundleEncoder) Handle() uintptr { return rbe.handle }

// Release releases the render bundle.
func (rb *RenderBundle) Release() {
	if rb.handle != 0 {
		untrackResource(rb.handle)
		procRenderBundleRelease.Call(rb.handle) //nolint:errcheck
		rb.handle = 0
	}
}

// Handle returns the underlying handle.
func (rb *RenderBundle) Handle() uintptr { return rb.handle }

// ExecuteBundles executes pre-recorded render bundles in the render pass.
// This is useful for replaying static geometry without re-recording commands.
func (rpe *RenderPassEncoder) ExecuteBundles(bundles []*RenderBundle) {
	mustInit()
	if len(bundles) == 0 {
		return
	}

	// Convert to handles
	handles := make([]uintptr, len(bundles))
	for i, b := range bundles {
		handles[i] = b.handle
	}

	procRenderPassEncoderExecuteBundles.Call( //nolint:errcheck
		rpe.handle,
		uintptr(len(handles)),
		uintptr(unsafe.Pointer(&handles[0])),
	)
}
