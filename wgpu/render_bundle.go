package wgpu

import (
	"unsafe"

	"github.com/gogpu/gputypes"
)

// RenderBundleEncoderDescriptor describes a render bundle encoder to create.
type RenderBundleEncoderDescriptor struct {
	Label              string
	ColorFormats       []gputypes.TextureFormat
	DepthStencilFormat gputypes.TextureFormat
	SampleCount        uint32
	DepthReadOnly      bool
	StencilReadOnly    bool
}

// RenderBundleDescriptor describes a render bundle.
type RenderBundleDescriptor struct {
	NextInChain uintptr // *ChainedStruct
	Label       StringView
}

// renderBundleEncoderDescriptorWire is the FFI-compatible C-layout struct.
// nextInChain(8)+label(16)+colorFormatCount(8)+colorFormats(8)+
// depthStencilFormat(4)+sampleCount(4)+depthReadOnly(4)+stencilReadOnly(4) = 56 bytes.
type renderBundleEncoderDescriptorWire struct {
	nextInChain        uintptr
	label              StringView
	colorFormatCount   uintptr
	colorFormats       uintptr // *uint32 (converted from gputypes.TextureFormat)
	depthStencilFormat uint32
	sampleCount        uint32
	depthReadOnly      Bool
	stencilReadOnly    Bool
}

// CreateRenderBundleEncoder creates a render bundle encoder for pre-recording render commands.
// Render bundles allow you to pre-record a sequence of render commands that can be replayed
// multiple times, which is useful for static geometry.
func (d *Device) CreateRenderBundleEncoder(desc *RenderBundleEncoderDescriptor) (*RenderBundleEncoder, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if d == nil || d.handle == 0 {
		return nil, &WGPUError{Op: "CreateRenderBundleEncoder", Message: "device is nil or released"}
	}
	if desc == nil {
		return nil, &WGPUError{Op: "CreateRenderBundleEncoder", Message: "descriptor is nil"}
	}

	wire := renderBundleEncoderDescriptorWire{
		label:              stringToStringView(desc.Label),
		colorFormatCount:   uintptr(len(desc.ColorFormats)),
		depthStencilFormat: uint32(desc.DepthStencilFormat),
		sampleCount:        desc.SampleCount,
		depthReadOnly:      boolToWGPU(desc.DepthReadOnly),
		stencilReadOnly:    boolToWGPU(desc.StencilReadOnly),
	}

	// Convert color formats to uint32 (gputypes v0.3.0 values equal wgpu-native v29 values)
	var convertedFormats []uint32
	if len(desc.ColorFormats) > 0 {
		convertedFormats = make([]uint32, len(desc.ColorFormats))
		for i, f := range desc.ColorFormats {
			convertedFormats[i] = uint32(f)
		}
		wire.colorFormats = uintptr(unsafe.Pointer(&convertedFormats[0]))
	}

	handle, _, _ := procDeviceCreateRenderBundleEncoder.Call(
		d.handle,
		uintptr(unsafe.Pointer(&wire)),
	)
	if handle == 0 {
		return nil, &WGPUError{Op: "CreateRenderBundleEncoder", Message: "wgpu returned null handle"}
	}
	trackResource(handle, "RenderBundleEncoder")
	return &RenderBundleEncoder{handle: handle}, nil
}

// CreateRenderBundleEncoderSimple creates a render bundle encoder with common settings.
func (d *Device) CreateRenderBundleEncoderSimple(colorFormats []gputypes.TextureFormat, depthFormat gputypes.TextureFormat, sampleCount uint32) *RenderBundleEncoder {
	enc, _ := d.CreateRenderBundleEncoder(&RenderBundleEncoderDescriptor{
		ColorFormats:       colorFormats,
		DepthStencilFormat: depthFormat,
		SampleCount:        sampleCount,
	})
	return enc
}

// SetPipeline sets the render pipeline for subsequent draw calls.
func (rbe *RenderBundleEncoder) SetPipeline(pipeline *RenderPipeline) {
	mustInit()
	if rbe == nil || rbe.handle == 0 || pipeline == nil || pipeline.handle == 0 {
		return
	}
	procRenderBundleEncoderSetPipeline.Call(rbe.handle, pipeline.handle) //nolint:errcheck
}

// SetBindGroup sets a bind group at the given index.
func (rbe *RenderBundleEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	mustInit()
	if rbe == nil || rbe.handle == 0 || group == nil || group.handle == 0 {
		return
	}
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
	if rbe == nil || rbe.handle == 0 || buffer == nil || buffer.handle == 0 {
		return
	}
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
	if rbe == nil || rbe.handle == 0 || buffer == nil || buffer.handle == 0 {
		return
	}
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
	if rbe == nil || rbe.handle == 0 {
		return
	}
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
	if rbe == nil || rbe.handle == 0 {
		return
	}
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
	if rbe == nil || rbe.handle == 0 || indirectBuffer == nil || indirectBuffer.handle == 0 {
		return
	}
	procRenderBundleEncoderDrawIndirect.Call( //nolint:errcheck
		rbe.handle,
		indirectBuffer.handle,
		uintptr(indirectOffset),
	)
}

// DrawIndexedIndirect records an indirect indexed draw call.
func (rbe *RenderBundleEncoder) DrawIndexedIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	mustInit()
	if rbe == nil || rbe.handle == 0 || indirectBuffer == nil || indirectBuffer.handle == 0 {
		return
	}
	procRenderBundleEncoderDrawIndexedIndirect.Call( //nolint:errcheck
		rbe.handle,
		indirectBuffer.handle,
		uintptr(indirectOffset),
	)
}

// Finish completes recording and returns the render bundle.
// The optional desc parameter allows specifying a label; if omitted, nil is used.
func (rbe *RenderBundleEncoder) Finish(desc ...*RenderBundleDescriptor) *RenderBundle {
	mustInit()
	if rbe == nil || rbe.handle == 0 {
		return nil
	}

	var descPtr uintptr
	if len(desc) > 0 && desc[0] != nil {
		descPtr = uintptr(unsafe.Pointer(desc[0]))
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
	if rpe == nil || rpe.handle == 0 || len(bundles) == 0 {
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
