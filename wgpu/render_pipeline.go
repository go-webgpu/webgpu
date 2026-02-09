package wgpu

import (
	"unsafe"

	"github.com/gogpu/gputypes"
)

// VertexAttribute describes a vertex attribute.
type VertexAttribute struct {
	Format         gputypes.VertexFormat
	Offset         uint64
	ShaderLocation uint32
	_pad           [4]byte
}

// vertexAttributeWire is the FFI-compatible structure with converted Format.
// Field order matches webgpu.h: format, offset, shaderLocation
type vertexAttributeWire struct {
	Format         uint32 // converted from gputypes.VertexFormat
	_pad1          [4]byte
	Offset         uint64
	ShaderLocation uint32
	_pad2          [4]byte
}

// VertexBufferLayout describes how vertex data is laid out in a buffer.
type VertexBufferLayout struct {
	ArrayStride    uint64
	StepMode       gputypes.VertexStepMode
	_pad           [4]byte
	AttributeCount uintptr
	Attributes     *VertexAttribute
}

// vertexBufferLayoutWire is the FFI-compatible structure with converted StepMode.
// Field order matches webgpu.h: stepMode, arrayStride, attributeCount, attributes
type vertexBufferLayoutWire struct {
	StepMode       uint32  // converted from gputypes.VertexStepMode
	_pad           [4]byte // padding to align arrayStride to 8 bytes
	ArrayStride    uint64
	AttributeCount uintptr
	Attributes     uintptr // pointer to VertexAttribute array
}

// vertexState is the native structure for vertex stage.
type vertexState struct {
	nextInChain   uintptr    // 8 bytes
	module        uintptr    // 8 bytes (WGPUShaderModule)
	entryPoint    StringView // 16 bytes
	constantCount uintptr    // 8 bytes
	constants     uintptr    // 8 bytes
	bufferCount   uintptr    // 8 bytes
	buffers       uintptr    // 8 bytes (pointer to VertexBufferLayout array)
}

// primitiveState is the native structure for primitive assembly.
type primitiveState struct {
	nextInChain      uintptr                    // 8 bytes
	topology         gputypes.PrimitiveTopology // 4 bytes
	stripIndexFormat gputypes.IndexFormat       // 4 bytes
	frontFace        gputypes.FrontFace         // 4 bytes
	cullMode         gputypes.CullMode          // 4 bytes
	unclippedDepth   Bool                       // 4 bytes
	_pad             [4]byte                    // 4 bytes padding
}

// multisampleState is the native structure for multisample state.
type multisampleState struct {
	nextInChain            uintptr // 8 bytes
	count                  uint32  // 4 bytes
	mask                   uint32  // 4 bytes
	alphaToCoverageEnabled Bool    // 4 bytes
	_pad                   [4]byte // 4 bytes padding
}

// BlendComponent describes blend state for a color component.
type BlendComponent struct {
	Operation gputypes.BlendOperation
	SrcFactor gputypes.BlendFactor
	DstFactor gputypes.BlendFactor
}

// BlendState describes how colors are blended.
type BlendState struct {
	Color BlendComponent
	Alpha BlendComponent
}

// colorTargetStateWire is the native FFI-compatible structure for a color target.
// CRITICAL: writeMask is uint64 because WGPUColorWriteMaskFlags = WGPUFlags = uint64 in webgpu-headers!
type colorTargetStateWire struct {
	nextInChain uintptr // 8 bytes
	format      uint32  // 4 bytes (WGPUTextureFormat, converted)
	_pad1       [4]byte // 4 bytes padding (to align blend to 8)
	blend       uintptr // 8 bytes (pointer to BlendState, nullable)
	writeMask   uint64  // 8 bytes (WGPUColorWriteMaskFlags = uint64!)
}

// fragmentState is the native structure for fragment stage.
type fragmentState struct {
	nextInChain   uintptr    // 8 bytes
	module        uintptr    // 8 bytes (WGPUShaderModule)
	entryPoint    StringView // 16 bytes
	constantCount uintptr    // 8 bytes
	constants     uintptr    // 8 bytes
	targetCount   uintptr    // 8 bytes
	targets       uintptr    // 8 bytes (pointer to colorTargetState array)
}

// renderPipelineDescriptor is the native structure for creating a render pipeline.
type renderPipelineDescriptor struct {
	nextInChain  uintptr          // 8 bytes
	label        StringView       // 16 bytes
	layout       uintptr          // 8 bytes (WGPUPipelineLayout, nullable)
	vertex       vertexState      // vertex stage
	primitive    primitiveState   // primitive assembly
	depthStencil uintptr          // 8 bytes (nullable)
	multisample  multisampleState // multisample state
	fragment     uintptr          // 8 bytes (nullable, pointer to fragmentState)
}

// ColorTargetState describes a render target in a render pipeline.
type ColorTargetState struct {
	Format    gputypes.TextureFormat
	Blend     *BlendState // nil for no blending
	WriteMask gputypes.ColorWriteMask
}

// VertexState describes the vertex stage of a render pipeline.
type VertexState struct {
	Module     *ShaderModule
	EntryPoint string
	Buffers    []VertexBufferLayout
}

// FragmentState describes the fragment stage of a render pipeline.
type FragmentState struct {
	Module     *ShaderModule
	EntryPoint string
	Targets    []ColorTargetState
}

// PrimitiveState describes how primitives are assembled.
type PrimitiveState struct {
	Topology         gputypes.PrimitiveTopology
	StripIndexFormat gputypes.IndexFormat
	FrontFace        gputypes.FrontFace
	CullMode         gputypes.CullMode
}

// MultisampleState describes multisampling.
type MultisampleState struct {
	Count                  uint32
	Mask                   uint32
	AlphaToCoverageEnabled bool
}

// StencilFaceState describes stencil operations for a face.
type StencilFaceState struct {
	Compare     gputypes.CompareFunction
	FailOp      gputypes.StencilOperation
	DepthFailOp gputypes.StencilOperation
	PassOp      gputypes.StencilOperation
}

// DepthStencilState describes depth and stencil test state (user API).
type DepthStencilState struct {
	Format              gputypes.TextureFormat
	DepthWriteEnabled   bool
	DepthCompare        gputypes.CompareFunction
	StencilFront        StencilFaceState
	StencilBack         StencilFaceState
	StencilReadMask     uint32
	StencilWriteMask    uint32
	DepthBias           int32
	DepthBiasSlopeScale float32
	DepthBiasClamp      float32
}

// depthStencilStateWire is the native FFI-compatible structure for depth/stencil state.
// Uses uint32 for format (converted from gputypes).
type depthStencilStateWire struct {
	nextInChain         uintptr
	format              uint32 // converted from gputypes.TextureFormat
	depthWriteEnabled   OptionalBool
	depthCompare        gputypes.CompareFunction
	stencilFront        StencilFaceState
	stencilBack         StencilFaceState
	stencilReadMask     uint32
	stencilWriteMask    uint32
	depthBias           int32
	depthBiasSlopeScale float32
	depthBiasClamp      float32
}

// RenderPipelineDescriptor describes a render pipeline to create.
type RenderPipelineDescriptor struct {
	Label        string
	Layout       *PipelineLayout // nil for auto layout
	Vertex       VertexState
	Primitive    PrimitiveState
	DepthStencil *DepthStencilState // nil for no depth/stencil
	Multisample  MultisampleState
	Fragment     *FragmentState // nil for no fragment stage (depth-only)
}

// CreateRenderPipeline creates a render pipeline.
func (d *Device) CreateRenderPipeline(desc *RenderPipelineDescriptor) *RenderPipeline {
	mustInit()

	if desc == nil {
		return nil
	}

	// Build vertex state
	var entryPointBytes []byte
	if desc.Vertex.EntryPoint != "" {
		entryPointBytes = append([]byte(desc.Vertex.EntryPoint), 0)
	}

	nativeVertex := vertexState{
		nextInChain:   0,
		module:        desc.Vertex.Module.handle,
		constantCount: 0,
		constants:     0,
		bufferCount:   uintptr(len(desc.Vertex.Buffers)),
	}

	if len(entryPointBytes) > 0 {
		nativeVertex.entryPoint = StringView{
			Data:   uintptr(unsafe.Pointer(&entryPointBytes[0])),
			Length: uintptr(len(entryPointBytes) - 1),
		}
	} else {
		nativeVertex.entryPoint = EmptyStringView()
	}

	// Convert vertex buffer layouts with StepMode and VertexFormat conversion
	var nativeBuffers []vertexBufferLayoutWire
	var allNativeAttrs [][]vertexAttributeWire // keep alive during FFI call
	if len(desc.Vertex.Buffers) > 0 {
		nativeBuffers = make([]vertexBufferLayoutWire, len(desc.Vertex.Buffers))
		allNativeAttrs = make([][]vertexAttributeWire, len(desc.Vertex.Buffers))
		for i, buf := range desc.Vertex.Buffers {
			var attrsPtr uintptr
			if buf.Attributes != nil && buf.AttributeCount > 0 {
				// Convert attributes with format conversion
				attrs := unsafe.Slice(buf.Attributes, buf.AttributeCount)
				allNativeAttrs[i] = make([]vertexAttributeWire, len(attrs))
				for j, attr := range attrs {
					allNativeAttrs[i][j] = vertexAttributeWire{
						Format:         toWGPUVertexFormat(attr.Format),
						Offset:         attr.Offset,
						ShaderLocation: attr.ShaderLocation,
					}
				}
				attrsPtr = uintptr(unsafe.Pointer(&allNativeAttrs[i][0]))
			}
			nativeBuffers[i] = vertexBufferLayoutWire{
				StepMode:       toWGPUVertexStepMode(buf.StepMode),
				ArrayStride:    buf.ArrayStride,
				AttributeCount: buf.AttributeCount,
				Attributes:     attrsPtr,
			}
		}
		nativeVertex.buffers = uintptr(unsafe.Pointer(&nativeBuffers[0]))
	}

	// Build primitive state
	nativePrimitive := primitiveState{
		nextInChain:      0,
		topology:         desc.Primitive.Topology,
		stripIndexFormat: desc.Primitive.StripIndexFormat,
		frontFace:        desc.Primitive.FrontFace,
		cullMode:         desc.Primitive.CullMode,
		unclippedDepth:   False,
	}

	// Build multisample state
	count := desc.Multisample.Count
	if count == 0 {
		count = 1
	}
	mask := desc.Multisample.Mask
	if mask == 0 {
		mask = 0xFFFFFFFF
	}
	alphaToCov := False
	if desc.Multisample.AlphaToCoverageEnabled {
		alphaToCov = True
	}

	nativeMultisample := multisampleState{
		nextInChain:            0,
		count:                  count,
		mask:                   mask,
		alphaToCoverageEnabled: alphaToCov,
	}

	// Build depth/stencil state if present (with format conversion)
	var depthStencilPtr uintptr
	var nativeDepthStencil depthStencilStateWire
	if desc.DepthStencil != nil {
		depthWriteOpt := OptionalBoolFalse
		if desc.DepthStencil.DepthWriteEnabled {
			depthWriteOpt = OptionalBoolTrue
		}

		nativeDepthStencil = depthStencilStateWire{
			nextInChain:         0,
			format:              toWGPUTextureFormat(desc.DepthStencil.Format),
			depthWriteEnabled:   depthWriteOpt,
			depthCompare:        desc.DepthStencil.DepthCompare,
			stencilFront:        desc.DepthStencil.StencilFront,
			stencilBack:         desc.DepthStencil.StencilBack,
			stencilReadMask:     desc.DepthStencil.StencilReadMask,
			stencilWriteMask:    desc.DepthStencil.StencilWriteMask,
			depthBias:           desc.DepthStencil.DepthBias,
			depthBiasSlopeScale: desc.DepthStencil.DepthBiasSlopeScale,
			depthBiasClamp:      desc.DepthStencil.DepthBiasClamp,
		}
		depthStencilPtr = uintptr(unsafe.Pointer(&nativeDepthStencil))
	}

	// Build fragment state if present
	var fragmentPtr uintptr
	var nativeFragment fragmentState
	var nativeTargets []colorTargetStateWire
	var fragEntryPointBytes []byte

	if desc.Fragment != nil {
		if desc.Fragment.EntryPoint != "" {
			fragEntryPointBytes = append([]byte(desc.Fragment.EntryPoint), 0)
		}

		nativeFragment = fragmentState{
			nextInChain:   0,
			module:        desc.Fragment.Module.handle,
			constantCount: 0,
			constants:     0,
			targetCount:   uintptr(len(desc.Fragment.Targets)),
		}

		if len(fragEntryPointBytes) > 0 {
			nativeFragment.entryPoint = StringView{
				Data:   uintptr(unsafe.Pointer(&fragEntryPointBytes[0])),
				Length: uintptr(len(fragEntryPointBytes) - 1),
			}
		} else {
			nativeFragment.entryPoint = EmptyStringView()
		}

		// Build color targets with wire format (uint64 writeMask!)
		nativeTargets = make([]colorTargetStateWire, len(desc.Fragment.Targets))
		for i, target := range desc.Fragment.Targets {
			convertedFormat := toWGPUTextureFormat(target.Format)
			nativeTargets[i] = colorTargetStateWire{
				nextInChain: 0,
				format:      convertedFormat,
				writeMask:   uint64(target.WriteMask), // widen to uint64
			}
			if target.Blend != nil {
				nativeTargets[i].blend = uintptr(unsafe.Pointer(target.Blend))
			}
			// DEBUG: print the target bytes
			_ = convertedFormat // silence unused warning
		}

		if len(nativeTargets) > 0 {
			nativeFragment.targets = uintptr(unsafe.Pointer(&nativeTargets[0]))
		}

		fragmentPtr = uintptr(unsafe.Pointer(&nativeFragment))
	}

	// Build pipeline layout
	var layoutHandle uintptr
	if desc.Layout != nil {
		layoutHandle = desc.Layout.handle
	}

	// Build the full descriptor
	nativeDesc := renderPipelineDescriptor{
		nextInChain:  0,
		label:        EmptyStringView(),
		layout:       layoutHandle,
		vertex:       nativeVertex,
		primitive:    nativePrimitive,
		depthStencil: depthStencilPtr,
		multisample:  nativeMultisample,
		fragment:     fragmentPtr,
	}

	handle, _, _ := procDeviceCreateRenderPipeline.Call(
		d.handle,
		uintptr(unsafe.Pointer(&nativeDesc)),
	)
	if handle == 0 {
		return nil
	}

	trackResource(handle, "RenderPipeline")
	return &RenderPipeline{handle: handle}
}

// CreateRenderPipelineSimple creates a simple render pipeline with common defaults.
func (d *Device) CreateRenderPipelineSimple(
	layout *PipelineLayout,
	vertexShader *ShaderModule,
	vertexEntryPoint string,
	fragmentShader *ShaderModule,
	fragmentEntryPoint string,
	targetFormat gputypes.TextureFormat,
) *RenderPipeline {
	return d.CreateRenderPipeline(&RenderPipelineDescriptor{
		Layout: layout,
		Vertex: VertexState{
			Module:     vertexShader,
			EntryPoint: vertexEntryPoint,
		},
		Primitive: PrimitiveState{
			Topology:  gputypes.PrimitiveTopologyTriangleList,
			FrontFace: gputypes.FrontFaceCCW,
			CullMode:  gputypes.CullModeNone,
		},
		Multisample: MultisampleState{
			Count: 1,
			Mask:  0xFFFFFFFF,
		},
		Fragment: &FragmentState{
			Module:     fragmentShader,
			EntryPoint: fragmentEntryPoint,
			Targets: []ColorTargetState{{
				Format:    targetFormat,
				WriteMask: gputypes.ColorWriteMaskAll,
			}},
		},
	})
}

// GetBindGroupLayout returns the bind group layout for the given index.
func (rp *RenderPipeline) GetBindGroupLayout(groupIndex uint32) *BindGroupLayout {
	mustInit()
	handle, _, _ := procRenderPipelineGetBindGroupLayout.Call(
		rp.handle,
		uintptr(groupIndex),
	)
	if handle == 0 {
		return nil
	}
	trackResource(handle, "BindGroupLayout")
	return &BindGroupLayout{handle: handle}
}

// Release releases the render pipeline.
func (rp *RenderPipeline) Release() {
	if rp.handle != 0 {
		untrackResource(rp.handle)
		procRenderPipelineRelease.Call(rp.handle) //nolint:errcheck
		rp.handle = 0
	}
}

// Handle returns the underlying handle. For advanced use only.
func (rp *RenderPipeline) Handle() uintptr { return rp.handle }
