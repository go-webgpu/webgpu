// Package main demonstrates a textured quad rendering using go-webgpu.
// This example creates a procedural checkerboard texture and renders it on a quad.
package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/go-webgpu/webgpu/wgpu"
	"golang.org/x/sys/windows"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "go-webgpu: Textured Quad Example"
	textureSize  = 256 // 256x256 texture
)

// Win32 constants
const (
	csHRedraw                 = 0x0002
	csVRedraw                 = 0x0001
	wmDestroy                 = 0x0002
	wmSize                    = 0x0005
	idcArrow                  = 32512
	colorWindow               = 5
	swShowNormal              = 1
	pmRemove                  = 0x0001
	wsOverlappedWindow        = 0x00CF0000
	wsVisible                 = 0x10000000
	cwUseDefault       uint32 = 0x80000000
)

var (
	user32               = windows.NewLazyDLL("user32.dll")
	procRegisterClassExW = user32.NewProc("RegisterClassExW")
	procCreateWindowExW  = user32.NewProc("CreateWindowExW")
	procShowWindow       = user32.NewProc("ShowWindow")
	procUpdateWindow     = user32.NewProc("UpdateWindow")
	procPeekMessageW     = user32.NewProc("PeekMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessageW = user32.NewProc("DispatchMessageW")
	procDefWindowProcW   = user32.NewProc("DefWindowProcW")
	procPostQuitMessage  = user32.NewProc("PostQuitMessage")
	procLoadCursorW      = user32.NewProc("LoadCursorW")
	procGetModuleHandleW = user32.NewProc("GetModuleHandleW")
)

// WNDCLASSEXW represents the Win32 WNDCLASSEXW structure.
type WNDCLASSEXW struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     windows.Handle
	hIcon         windows.Handle
	hCursor       windows.Handle
	hbrBackground windows.Handle
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       windows.Handle
}

// MSG represents the Win32 MSG structure.
type MSG struct {
	hwnd    windows.HWND
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      struct{ x, y int32 }
}

// Application state
type App struct {
	hwnd           windows.HWND
	hinstance      windows.Handle
	instance       *wgpu.Instance
	adapter        *wgpu.Adapter
	device         *wgpu.Device
	queue          *wgpu.Queue
	surface        *wgpu.Surface
	pipeline       *wgpu.RenderPipeline
	vertexBuffer   *wgpu.Buffer
	indexBuffer    *wgpu.Buffer
	texture        *wgpu.Texture
	textureView    *wgpu.TextureView
	sampler        *wgpu.Sampler
	bindGroupLyt   *wgpu.BindGroupLayout
	bindGroup      *wgpu.BindGroup
	width          uint32
	height         uint32
	running        bool
	needsRecreate  bool
	surfaceTex     *wgpu.SurfaceTexture
	surfaceTexView *wgpu.TextureView
}

// Shader source (WGSL) with texture sampling
const shaderSource = `
struct VertexInput {
    @location(0) position: vec2f,
    @location(1) uv: vec2f,
}

struct VertexOutput {
    @builtin(position) position: vec4f,
    @location(0) uv: vec2f,
}

@group(0) @binding(0) var texSampler: sampler;
@group(0) @binding(1) var tex: texture_2d<f32>;

@vertex
fn vs_main(in: VertexInput) -> VertexOutput {
    var out: VertexOutput;
    out.position = vec4f(in.position, 0.0, 1.0);
    out.uv = in.uv;
    return out;
}

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4f {
    return textureSample(tex, texSampler, in.uv);
}
`

func main() {
	app := &App{
		width:   windowWidth,
		height:  windowHeight,
		running: true,
	}

	if err := app.init(); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	defer app.cleanup()

	app.run()
}

// init initializes the application.
func (app *App) init() error {
	// Get HINSTANCE
	ret, _, _ := procGetModuleHandleW.Call(0)
	app.hinstance = windows.Handle(ret)

	// Create window
	if err := app.createWindow(); err != nil {
		return fmt.Errorf("create window: %w", err)
	}

	// Initialize WebGPU
	if err := app.initWebGPU(); err != nil {
		return fmt.Errorf("init webgpu: %w", err)
	}

	// Configure surface
	if err := app.configureSurface(); err != nil {
		return fmt.Errorf("configure surface: %w", err)
	}

	// Create texture
	if err := app.createTexture(); err != nil {
		return fmt.Errorf("create texture: %w", err)
	}

	// Create sampler
	if err := app.createSampler(); err != nil {
		return fmt.Errorf("create sampler: %w", err)
	}

	// Create bind group
	if err := app.createBindGroup(); err != nil {
		return fmt.Errorf("create bind group: %w", err)
	}

	// Create vertex and index buffers
	if err := app.createBuffers(); err != nil {
		return fmt.Errorf("create buffers: %w", err)
	}

	// Create render pipeline
	if err := app.createPipeline(); err != nil {
		return fmt.Errorf("create pipeline: %w", err)
	}

	return nil
}

// createWindow creates the main window.
func (app *App) createWindow() error {
	className, err := windows.UTF16PtrFromString("GoWebGPUTexturedQuad")
	if err != nil {
		return err
	}

	wndClass := WNDCLASSEXW{
		cbSize:        uint32(unsafe.Sizeof(WNDCLASSEXW{})),
		style:         csHRedraw | csVRedraw,
		lpfnWndProc:   syscall.NewCallback(app.wndProc),
		hInstance:     app.hinstance,
		lpszClassName: className,
	}

	// Load default cursor
	cursor, _, _ := procLoadCursorW.Call(0, uintptr(idcArrow))
	wndClass.hCursor = windows.Handle(cursor)
	wndClass.hbrBackground = windows.Handle(colorWindow + 1)

	// nolint:gosec // Required for Win32 FFI - passing struct to Windows API
	ret, _, _ := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wndClass)))
	if ret == 0 {
		return fmt.Errorf("RegisterClassExW failed")
	}

	titlePtr, err := windows.UTF16PtrFromString(windowTitle)
	if err != nil {
		return err
	}

	// nolint:gosec // Required for Win32 FFI - passing string pointers to Windows API
	hwnd, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(wsOverlappedWindow|wsVisible),
		uintptr(cwUseDefault),
		uintptr(cwUseDefault),
		uintptr(app.width),
		uintptr(app.height),
		0,
		0,
		uintptr(app.hinstance),
		0,
	)

	if hwnd == 0 {
		return fmt.Errorf("CreateWindowExW failed")
	}

	app.hwnd = windows.HWND(hwnd)

	_, _, _ = procShowWindow.Call(uintptr(app.hwnd), swShowNormal)
	_, _, _ = procUpdateWindow.Call(uintptr(app.hwnd))

	return nil
}

// wndProc is the window procedure callback.
func (app *App) wndProc(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case wmDestroy:
		app.running = false
		_, _, _ = procPostQuitMessage.Call(0)
		return 0
	case wmSize:
		newWidth := uint32(lParam & 0xFFFF)
		newHeight := uint32((lParam >> 16) & 0xFFFF)
		if newWidth != app.width || newHeight != app.height {
			app.width = newWidth
			app.height = newHeight
			app.needsRecreate = true
		}
		return 0
	}
	ret, _, _ := procDefWindowProcW.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam,
	)
	return ret
}

// initWebGPU initializes WebGPU resources.
func (app *App) initWebGPU() error {
	// Create instance
	inst, err := wgpu.CreateInstance(nil)
	if err != nil {
		return fmt.Errorf("create instance: %w", err)
	}
	app.instance = inst

	// Request adapter
	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		return fmt.Errorf("request adapter: %w", err)
	}
	app.adapter = adapter

	// Request device
	device, err := adapter.RequestDevice(nil)
	if err != nil {
		return fmt.Errorf("request device: %w", err)
	}
	app.device = device

	// Get queue
	app.queue = device.GetQueue()

	// Create surface
	surface, err := inst.CreateSurfaceFromWindowsHWND(uintptr(app.hinstance), uintptr(app.hwnd))
	if err != nil {
		return fmt.Errorf("create surface: %w", err)
	}
	app.surface = surface

	return nil
}

// configureSurface configures the surface for rendering.
func (app *App) configureSurface() error {
	app.surface.Configure(&wgpu.SurfaceConfiguration{
		Device:      app.device,
		Format:      wgpu.TextureFormatBGRA8Unorm,
		Usage:       wgpu.TextureUsageRenderAttachment,
		Width:       app.width,
		Height:      app.height,
		AlphaMode:   wgpu.CompositeAlphaModeOpaque,
		PresentMode: wgpu.PresentModeFifo,
	})
	app.needsRecreate = false
	return nil
}

// createTexture creates a procedural checkerboard texture.
// nolint:funlen // FFI descriptor builders and texture data generation require more lines
func (app *App) createTexture() error {
	// Create 256x256 RGBA8 checkerboard texture
	const size = textureSize
	const bytesPerPixel = 4 // RGBA8
	textureData := make([]byte, size*size*bytesPerPixel)

	// Generate checkerboard pattern (8x8 squares)
	const squareSize = 32 // 256 / 8 = 32 pixels per square
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Determine if this pixel is in a light or dark square
			squareX := x / squareSize
			squareY := y / squareSize
			isLight := (squareX+squareY)%2 == 0

			offset := (y*size + x) * bytesPerPixel
			if isLight {
				// Light yellow
				textureData[offset] = 255   // R
				textureData[offset+1] = 255 // G
				textureData[offset+2] = 200 // B
				textureData[offset+3] = 255 // A
			} else {
				// Dark blue
				textureData[offset] = 50    // R
				textureData[offset+1] = 50  // G
				textureData[offset+2] = 150 // B
				textureData[offset+3] = 255 // A
			}
		}
	}

	// Create texture
	textureDesc := wgpu.TextureDescriptor{
		Label:     wgpu.EmptyStringView(),
		Usage:     wgpu.TextureUsageTextureBinding | wgpu.TextureUsageCopyDst,
		Dimension: wgpu.TextureDimension2D,
		Size: wgpu.Extent3D{
			Width:              size,
			Height:             size,
			DepthOrArrayLayers: 1,
		},
		Format:        wgpu.TextureFormatRGBA8Unorm,
		MipLevelCount: 1,
		SampleCount:   1,
	}

	app.texture = app.device.CreateTexture(&textureDesc)
	if app.texture == nil {
		return fmt.Errorf("failed to create texture")
	}

	// Upload texture data
	// BytesPerRow must be aligned to 256 bytes
	const bytesPerRow = size * bytesPerPixel
	const alignedBytesPerRow = ((bytesPerRow + 255) / 256) * 256

	destInfo := wgpu.TexelCopyTextureInfo{
		Texture:  app.texture.Handle(),
		MipLevel: 0,
		Origin:   wgpu.Origin3D{X: 0, Y: 0, Z: 0},
		Aspect:   wgpu.TextureAspectAll,
	}

	layout := wgpu.TexelCopyBufferLayout{
		Offset:       0,
		BytesPerRow:  alignedBytesPerRow,
		RowsPerImage: size,
	}

	extent := wgpu.Extent3D{
		Width:              size,
		Height:             size,
		DepthOrArrayLayers: 1,
	}

	app.queue.WriteTexture(&destInfo, textureData, &layout, &extent)

	// Create texture view
	app.textureView = app.texture.CreateView(nil)
	if app.textureView == nil {
		return fmt.Errorf("failed to create texture view")
	}

	return nil
}

// createSampler creates a linear sampler.
func (app *App) createSampler() error {
	app.sampler = app.device.CreateLinearSampler()
	if app.sampler == nil {
		return fmt.Errorf("failed to create sampler")
	}
	return nil
}

// createBindGroup creates bind group layout and bind group for texture and sampler.
func (app *App) createBindGroup() error {
	// Create bind group layout
	layoutEntries := []wgpu.BindGroupLayoutEntry{
		{
			Binding:    0,
			Visibility: wgpu.ShaderStageFragment,
			Sampler: wgpu.SamplerBindingLayout{
				Type: wgpu.SamplerBindingTypeFiltering,
			},
		},
		{
			Binding:    1,
			Visibility: wgpu.ShaderStageFragment,
			Texture: wgpu.TextureBindingLayout{
				SampleType:    wgpu.TextureSampleTypeFloat,
				ViewDimension: wgpu.TextureViewDimension2D,
				Multisampled:  wgpu.False,
			},
		},
	}

	app.bindGroupLyt = app.device.CreateBindGroupLayoutSimple(layoutEntries)
	if app.bindGroupLyt == nil {
		return fmt.Errorf("failed to create bind group layout")
	}

	// Create bind group
	entries := []wgpu.BindGroupEntry{
		wgpu.SamplerBindingEntry(0, app.sampler),
		wgpu.TextureBindingEntry(1, app.textureView),
	}

	app.bindGroup = app.device.CreateBindGroupSimple(app.bindGroupLyt, entries)
	if app.bindGroup == nil {
		return fmt.Errorf("failed to create bind group")
	}

	return nil
}

// createBuffers creates vertex and index buffers for a quad.
func (app *App) createBuffers() error {
	// Vertex data: position (x, y) + uv (u, v)
	// Quad covers most of the screen
	vertices := []float32{
		// position      uv
		-0.8, 0.8, 0.0, 0.0, // top-left
		0.8, 0.8, 1.0, 0.0, // top-right
		-0.8, -0.8, 0.0, 1.0, // bottom-left
		0.8, -0.8, 1.0, 1.0, // bottom-right
	}

	// Index data (2 triangles)
	indices := []uint16{
		0, 1, 2, // first triangle
		2, 1, 3, // second triangle
	}

	// Create vertex buffer
	// nolint:gosec // len(vertices) is 16, * 4 = 64 bytes - no overflow risk
	vertexBufferSize := uint64(len(vertices) * 4)
	app.vertexBuffer = app.device.CreateBuffer(&wgpu.BufferDescriptor{
		Label:            wgpu.EmptyStringView(),
		Usage:            wgpu.BufferUsageVertex | wgpu.BufferUsageCopyDst,
		Size:             vertexBufferSize,
		MappedAtCreation: wgpu.True,
	})
	if app.vertexBuffer == nil {
		return fmt.Errorf("failed to create vertex buffer")
	}

	// Copy vertex data
	vPtr := app.vertexBuffer.GetMappedRange(0, vertexBufferSize)
	if vPtr == nil {
		return fmt.Errorf("failed to get mapped range for vertex buffer")
	}
	// nolint:gosec,govet // vPtr is from GetMappedRange, validated non-nil, safe for slice conversion
	vMappedSlice := unsafe.Slice((*float32)(vPtr), len(vertices))
	copy(vMappedSlice, vertices)
	app.vertexBuffer.Unmap()

	// Create index buffer
	// nolint:gosec // len(indices) is 6, * 2 = 12 bytes - no overflow risk
	indexBufferSize := uint64(len(indices) * 2)
	app.indexBuffer = app.device.CreateBuffer(&wgpu.BufferDescriptor{
		Label:            wgpu.EmptyStringView(),
		Usage:            wgpu.BufferUsageIndex | wgpu.BufferUsageCopyDst,
		Size:             indexBufferSize,
		MappedAtCreation: wgpu.True,
	})
	if app.indexBuffer == nil {
		return fmt.Errorf("failed to create index buffer")
	}

	// Copy index data
	iPtr := app.indexBuffer.GetMappedRange(0, indexBufferSize)
	if iPtr == nil {
		return fmt.Errorf("failed to get mapped range for index buffer")
	}
	// nolint:gosec,govet // iPtr is from GetMappedRange, validated non-nil, safe for slice conversion
	iMappedSlice := unsafe.Slice((*uint16)(iPtr), len(indices))
	copy(iMappedSlice, indices)
	app.indexBuffer.Unmap()

	return nil
}

// createPipeline creates the render pipeline with texture bind group.
// nolint:funlen // FFI descriptor builders with detailed pipeline configuration
func (app *App) createPipeline() error {
	// Create shader module
	shader := app.device.CreateShaderModuleWGSL(shaderSource)
	if shader == nil {
		return fmt.Errorf("failed to create shader module")
	}
	defer shader.Release()

	// Create pipeline layout with bind group layout
	pipelineLayout := app.device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		Label:                wgpu.EmptyStringView(),
		BindGroupLayoutCount: 1,
		BindGroupLayouts:     app.bindGroupLyt.Handle(),
	})
	if pipelineLayout == nil {
		return fmt.Errorf("failed to create pipeline layout")
	}
	// Note: pipelineLayout is passed to pipeline descriptor and will be cleaned up with pipeline
	defer pipelineLayout.Release()

	// Define vertex attributes
	attributes := []wgpu.VertexAttribute{
		{
			Format:         wgpu.VertexFormatFloat32x2, // position: vec2f
			Offset:         0,
			ShaderLocation: 0,
		},
		{
			Format:         wgpu.VertexFormatFloat32x2, // uv: vec2f
			Offset:         8,                          // 2 floats * 4 bytes = 8 bytes offset
			ShaderLocation: 1,
		},
	}

	// Create render pipeline
	pipeline := app.device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label:  "",
		Layout: pipelineLayout,
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
			Buffers: []wgpu.VertexBufferLayout{{
				ArrayStride:    16, // 4 floats * 4 bytes = 16 bytes per vertex
				StepMode:       wgpu.VertexStepModeVertex,
				AttributeCount: 2,
				Attributes:     &attributes[0],
			}},
		},
		Primitive: wgpu.PrimitiveState{
			Topology:  wgpu.PrimitiveTopologyTriangleList,
			FrontFace: wgpu.FrontFaceCCW,
			CullMode:  wgpu.CullModeNone,
		},
		Multisample: wgpu.MultisampleState{
			Count: 1,
			Mask:  0xFFFFFFFF,
		},
		Fragment: &wgpu.FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{{
				Format:    wgpu.TextureFormatBGRA8Unorm,
				WriteMask: wgpu.ColorWriteMaskAll,
			}},
		},
	})

	if pipeline == nil {
		return fmt.Errorf("failed to create render pipeline")
	}

	app.pipeline = pipeline
	return nil
}

// releasePreviousFrame releases resources from the previous frame.
func (app *App) releasePreviousFrame() {
	if app.surfaceTexView != nil {
		app.surfaceTexView.Release()
		app.surfaceTexView = nil
	}
	if app.surfaceTex != nil && app.surfaceTex.Texture != nil {
		app.surfaceTex.Texture.Release()
		app.surfaceTex = nil
	}
}

// acquireSurfaceTexture gets the current surface texture.
func (app *App) acquireSurfaceTexture() error {
	surfaceTex, err := app.surface.GetCurrentTexture()
	if err != nil {
		// Handle common surface errors
		if err == wgpu.ErrSurfaceLost || err == wgpu.ErrSurfaceNeedsReconfigure {
			app.needsRecreate = true
			return nil
		}
		return fmt.Errorf("get current texture: %w", err)
	}
	app.surfaceTex = surfaceTex

	// Create texture view
	view := surfaceTex.Texture.CreateView(nil)
	if view == nil {
		return fmt.Errorf("failed to create texture view")
	}
	app.surfaceTexView = view
	return nil
}

// renderQuad encodes the textured quad rendering commands.
func (app *App) renderQuad(encoder *wgpu.CommandEncoder, view *wgpu.TextureView) error {
	pass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		Label: "Textured Quad Render Pass",
		ColorAttachments: []wgpu.RenderPassColorAttachment{{
			View:    view,
			LoadOp:  wgpu.LoadOpClear,
			StoreOp: wgpu.StoreOpStore,
			ClearValue: wgpu.Color{
				R: 0.05,
				G: 0.05,
				B: 0.1,
				A: 1.0,
			},
		}},
	})
	if pass == nil {
		return fmt.Errorf("failed to begin render pass")
	}
	defer pass.Release()

	pass.SetPipeline(app.pipeline)
	pass.SetBindGroup(0, app.bindGroup, nil)
	pass.SetVertexBuffer(0, app.vertexBuffer, 0, uint64(16*4))                   // 16 floats * 4 bytes
	pass.SetIndexBuffer(app.indexBuffer, wgpu.IndexFormatUint16, 0, uint64(6*2)) // 6 indices * 2 bytes
	pass.DrawIndexed(6, 1, 0, 0, 0)
	pass.End()
	return nil
}

// render draws a frame.
func (app *App) render() error {
	// Recreate surface if needed
	if app.needsRecreate {
		if err := app.configureSurface(); err != nil {
			return fmt.Errorf("reconfigure surface: %w", err)
		}
	}

	// Release previous frame resources
	app.releasePreviousFrame()

	// Acquire surface texture
	if err := app.acquireSurfaceTexture(); err != nil {
		return err
	}

	// Create command encoder
	encoder := app.device.CreateCommandEncoder(nil)
	if encoder == nil {
		return fmt.Errorf("failed to create command encoder")
	}
	defer encoder.Release()

	// Render quad
	if err := app.renderQuad(encoder, app.surfaceTexView); err != nil {
		return err
	}

	// Finish encoding
	cmdBuffer := encoder.Finish(nil)
	if cmdBuffer == nil {
		return fmt.Errorf("failed to finish command encoder")
	}
	defer cmdBuffer.Release()

	// Submit commands and present
	app.queue.Submit(cmdBuffer)
	app.surface.Present()

	return nil
}

// run is the main application loop.
func (app *App) run() {
	for app.running {
		// Process Windows messages
		var msg MSG
		for {
			// nolint:gosec // Required for Win32 FFI - passing MSG struct to Windows API
			ret, _, _ := procPeekMessageW.Call(
				uintptr(unsafe.Pointer(&msg)),
				0,
				0,
				0,
				pmRemove,
			)
			if ret == 0 {
				break
			}
			// nolint:gosec // Required for Win32 FFI - passing MSG struct to Windows API
			_, _, _ = procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			// nolint:gosec // Required for Win32 FFI - passing MSG struct to Windows API
			_, _, _ = procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
		}

		// Render frame
		if err := app.render(); err != nil {
			fmt.Fprintf(os.Stderr, "Render error: %v\n", err)
			app.running = false
		}
	}
}

// cleanup releases all resources.
// nolint:cyclop // Resource cleanup naturally has many branches
func (app *App) cleanup() {
	if app.surfaceTexView != nil {
		app.surfaceTexView.Release()
	}
	if app.surfaceTex != nil && app.surfaceTex.Texture != nil {
		app.surfaceTex.Texture.Release()
	}
	if app.bindGroup != nil {
		app.bindGroup.Release()
	}
	if app.bindGroupLyt != nil {
		app.bindGroupLyt.Release()
	}
	if app.sampler != nil {
		app.sampler.Release()
	}
	if app.textureView != nil {
		app.textureView.Release()
	}
	if app.texture != nil {
		app.texture.Release()
	}
	if app.indexBuffer != nil {
		app.indexBuffer.Release()
	}
	if app.vertexBuffer != nil {
		app.vertexBuffer.Release()
	}
	if app.pipeline != nil {
		app.pipeline.Release()
	}
	if app.surface != nil {
		app.surface.Release()
	}
	if app.queue != nil {
		app.queue.Release()
	}
	if app.device != nil {
		app.device.Release()
	}
	if app.adapter != nil {
		app.adapter.Release()
	}
	if app.instance != nil {
		app.instance.Release()
	}
}
