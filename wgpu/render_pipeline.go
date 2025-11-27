package wgpu

import "unsafe"

// VertexAttribute describes a vertex attribute.
type VertexAttribute struct {
	Format         VertexFormat
	Offset         uint64
	ShaderLocation uint32
	_pad           [4]byte
}

// VertexBufferLayout describes how vertex data is laid out in a buffer.
type VertexBufferLayout struct {
	ArrayStride    uint64
	StepMode       VertexStepMode
	_pad           [4]byte
	AttributeCount uintptr
	Attributes     *VertexAttribute
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
	nextInChain      uintptr           // 8 bytes
	topology         PrimitiveTopology // 4 bytes
	stripIndexFormat IndexFormat       // 4 bytes
	frontFace        FrontFace         // 4 bytes
	cullMode         CullMode          // 4 bytes
	unclippedDepth   Bool              // 4 bytes
	_pad             [4]byte           // 4 bytes padding
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
	Operation BlendOperation
	SrcFactor BlendFactor
	DstFactor BlendFactor
}

// BlendState describes how colors are blended.
type BlendState struct {
	Color BlendComponent
	Alpha BlendComponent
}

// colorTargetState is the native structure for a color target.
type colorTargetState struct {
	nextInChain uintptr        // 8 bytes
	format      TextureFormat  // 4 bytes
	_pad1       [4]byte        // 4 bytes padding
	blend       uintptr        // 8 bytes (pointer to BlendState, nullable)
	writeMask   ColorWriteMask // 4 bytes
	_pad2       [4]byte        // 4 bytes padding
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

// BlendOperation describes a blend operation.
type BlendOperation uint32

const (
	BlendOperationUndefined       BlendOperation = 0x00
	BlendOperationAdd             BlendOperation = 0x01
	BlendOperationSubtract        BlendOperation = 0x02
	BlendOperationReverseSubtract BlendOperation = 0x03
	BlendOperationMin             BlendOperation = 0x04
	BlendOperationMax             BlendOperation = 0x05
)

// BlendFactor describes a blend factor.
type BlendFactor uint32

const (
	BlendFactorUndefined         BlendFactor = 0x00
	BlendFactorZero              BlendFactor = 0x01
	BlendFactorOne               BlendFactor = 0x02
	BlendFactorSrc               BlendFactor = 0x03
	BlendFactorOneMinusSrc       BlendFactor = 0x04
	BlendFactorSrcAlpha          BlendFactor = 0x05
	BlendFactorOneMinusSrcAlpha  BlendFactor = 0x06
	BlendFactorDst               BlendFactor = 0x07
	BlendFactorOneMinusDst       BlendFactor = 0x08
	BlendFactorDstAlpha          BlendFactor = 0x09
	BlendFactorOneMinusDstAlpha  BlendFactor = 0x0A
	BlendFactorSrcAlphaSaturated BlendFactor = 0x0B
	BlendFactorConstant          BlendFactor = 0x0C
	BlendFactorOneMinusConstant  BlendFactor = 0x0D
	BlendFactorSrc1              BlendFactor = 0x0E
	BlendFactorOneMinusSrc1      BlendFactor = 0x0F
	BlendFactorSrc1Alpha         BlendFactor = 0x10
	BlendFactorOneMinusSrc1Alpha BlendFactor = 0x11
)

// ColorTargetState describes a render target in a render pipeline.
type ColorTargetState struct {
	Format    TextureFormat
	Blend     *BlendState // nil for no blending
	WriteMask ColorWriteMask
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
	Topology         PrimitiveTopology
	StripIndexFormat IndexFormat
	FrontFace        FrontFace
	CullMode         CullMode
}

// MultisampleState describes multisampling.
type MultisampleState struct {
	Count                  uint32
	Mask                   uint32
	AlphaToCoverageEnabled bool
}

// StencilFaceState describes stencil operations for a face.
type StencilFaceState struct {
	Compare     CompareFunction
	FailOp      StencilOperation
	DepthFailOp StencilOperation
	PassOp      StencilOperation
}

// DepthStencilState describes depth and stencil test state (user API).
type DepthStencilState struct {
	Format              TextureFormat
	DepthWriteEnabled   bool
	DepthCompare        CompareFunction
	StencilFront        StencilFaceState
	StencilBack         StencilFaceState
	StencilReadMask     uint32
	StencilWriteMask    uint32
	DepthBias           int32
	DepthBiasSlopeScale float32
	DepthBiasClamp      float32
}

// depthStencilState is the native structure for depth/stencil state (72 bytes).
type depthStencilState struct {
	nextInChain         uintptr
	format              TextureFormat
	depthWriteEnabled   OptionalBool
	depthCompare        CompareFunction
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

	if len(desc.Vertex.Buffers) > 0 {
		nativeVertex.buffers = uintptr(unsafe.Pointer(&desc.Vertex.Buffers[0]))
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

	// Build depth/stencil state if present
	var depthStencilPtr uintptr
	var nativeDepthStencil depthStencilState
	if desc.DepthStencil != nil {
		depthWriteOpt := OptionalBoolFalse
		if desc.DepthStencil.DepthWriteEnabled {
			depthWriteOpt = OptionalBoolTrue
		}

		nativeDepthStencil = depthStencilState{
			nextInChain:         0,
			format:              desc.DepthStencil.Format,
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
	var nativeTargets []colorTargetState
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

		// Build color targets
		nativeTargets = make([]colorTargetState, len(desc.Fragment.Targets))
		for i, target := range desc.Fragment.Targets {
			nativeTargets[i] = colorTargetState{
				nextInChain: 0,
				format:      target.Format,
				writeMask:   target.WriteMask,
			}
			if target.Blend != nil {
				nativeTargets[i].blend = uintptr(unsafe.Pointer(target.Blend))
			}
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

	return &RenderPipeline{handle: handle}
}

// CreateRenderPipelineSimple creates a simple render pipeline with common defaults.
func (d *Device) CreateRenderPipelineSimple(
	layout *PipelineLayout,
	vertexShader *ShaderModule,
	vertexEntryPoint string,
	fragmentShader *ShaderModule,
	fragmentEntryPoint string,
	targetFormat TextureFormat,
) *RenderPipeline {
	return d.CreateRenderPipeline(&RenderPipelineDescriptor{
		Layout: layout,
		Vertex: VertexState{
			Module:     vertexShader,
			EntryPoint: vertexEntryPoint,
		},
		Primitive: PrimitiveState{
			Topology:  PrimitiveTopologyTriangleList,
			FrontFace: FrontFaceCCW,
			CullMode:  CullModeNone,
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
				WriteMask: ColorWriteMaskAll,
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
	return &BindGroupLayout{handle: handle}
}

// Release releases the render pipeline.
func (rp *RenderPipeline) Release() {
	if rp.handle != 0 {
		procRenderPipelineRelease.Call(rp.handle) //nolint:errcheck
		rp.handle = 0
	}
}

// Handle returns the underlying handle. For advanced use only.
func (rp *RenderPipeline) Handle() uintptr { return rp.handle }
