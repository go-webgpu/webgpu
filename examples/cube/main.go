// Package main demonstrates a rotating 3D cube with depth buffer using go-webgpu.
// This example creates a window using Windows API and renders a rotating colored cube
// with proper depth testing by updating MVP matrices in a uniform buffer each frame.
package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-webgpu/webgpu/wgpu"
	"golang.org/x/sys/windows"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "go-webgpu: Rotating Cube Example"
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
	hwnd             windows.HWND
	hinstance        windows.Handle
	instance         *wgpu.Instance
	adapter          *wgpu.Adapter
	device           *wgpu.Device
	queue            *wgpu.Queue
	surface          *wgpu.Surface
	pipeline         *wgpu.RenderPipeline
	vertexBuffer     *wgpu.Buffer
	uniformBuffer    *wgpu.Buffer
	bindGroupLayout  *wgpu.BindGroupLayout
	bindGroup        *wgpu.BindGroup
	depthTexture     *wgpu.Texture
	depthTextureView *wgpu.TextureView
	width            uint32
	height           uint32
	running          bool
	needsRecreate    bool
	surfaceTex       *wgpu.SurfaceTexture
	surfaceTexView   *wgpu.TextureView
	startTime        time.Time
}

// Shader source (WGSL) with uniform buffer for MVP matrix
const shaderSource = `
struct Uniforms {
    mvp: mat4x4f,
}

@group(0) @binding(0) var<uniform> uniforms: Uniforms;

struct VertexOutput {
    @builtin(position) position: vec4f,
    @location(0) color: vec3f,
}

@vertex
fn vs_main(@location(0) pos: vec3f, @location(1) color: vec3f) -> VertexOutput {
    var out: VertexOutput;
    out.position = uniforms.mvp * vec4f(pos, 1.0);
    out.color = color;
    return out;
}

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4f {
    return vec4f(in.color, 1.0);
}
`

func main() {
	app := &App{
		width:     windowWidth,
		height:    windowHeight,
		running:   true,
		startTime: time.Now(),
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

	// Create depth texture
	if err := app.createDepthTexture(); err != nil {
		return fmt.Errorf("create depth texture: %w", err)
	}

	// Create vertex buffer
	if err := app.createVertexBuffer(); err != nil {
		return fmt.Errorf("create vertex buffer: %w", err)
	}

	// Create uniform buffer
	if err := app.createUniformBuffer(); err != nil {
		return fmt.Errorf("create uniform buffer: %w", err)
	}

	// Create bind group layout
	if err := app.createBindGroupLayout(); err != nil {
		return fmt.Errorf("create bind group layout: %w", err)
	}

	// Create bind group
	if err := app.createBindGroup(); err != nil {
		return fmt.Errorf("create bind group: %w", err)
	}

	// Create render pipeline
	if err := app.createPipeline(); err != nil {
		return fmt.Errorf("create pipeline: %w", err)
	}

	return nil
}

// createWindow creates the main window.
func (app *App) createWindow() error {
	className, err := windows.UTF16PtrFromString("GoWebGPUCube")
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

// createDepthTexture creates the depth texture and view.
func (app *App) createDepthTexture() error {
	app.depthTexture = app.device.CreateDepthTexture(app.width, app.height, wgpu.TextureFormatDepth24Plus)
	if app.depthTexture == nil {
		return fmt.Errorf("failed to create depth texture")
	}

	app.depthTextureView = app.depthTexture.CreateView(nil)
	if app.depthTextureView == nil {
		return fmt.Errorf("failed to create depth texture view")
	}

	return nil
}

// recreateDepthTexture recreates the depth texture after resize.
func (app *App) recreateDepthTexture() error {
	// Release old depth texture
	if app.depthTextureView != nil {
		app.depthTextureView.Release()
		app.depthTextureView = nil
	}
	if app.depthTexture != nil {
		app.depthTexture.Release()
		app.depthTexture = nil
	}

	// Create new depth texture
	return app.createDepthTexture()
}

// createVertexBuffer creates a vertex buffer with cube geometry (36 vertices).
func (app *App) createVertexBuffer() error {
	// Cube vertices: 36 vertices (6 faces × 2 triangles × 3 vertices)
	// Each vertex: position (vec3f) + color (vec3f) = 6 floats = 24 bytes
	vertices := []float32{
		// Front face (Z+) - Red
		-0.5, -0.5, 0.5, 1.0, 0.0, 0.0, // 0
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0, // 1
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, // 2
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0, // 2
		-0.5, 0.5, 0.5, 1.0, 0.0, 0.0, // 3
		-0.5, -0.5, 0.5, 1.0, 0.0, 0.0, // 0

		// Back face (Z-) - Green
		-0.5, -0.5, -0.5, 0.0, 1.0, 0.0, // 4
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, // 7
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0, // 6
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0, // 6
		0.5, -0.5, -0.5, 0.0, 1.0, 0.0, // 5
		-0.5, -0.5, -0.5, 0.0, 1.0, 0.0, // 4

		// Top face (Y+) - Blue
		-0.5, 0.5, -0.5, 0.0, 0.0, 1.0, // 7
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, // 3
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, // 2
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0, // 2
		0.5, 0.5, -0.5, 0.0, 0.0, 1.0, // 6
		-0.5, 0.5, -0.5, 0.0, 0.0, 1.0, // 7

		// Bottom face (Y-) - Yellow
		-0.5, -0.5, -0.5, 1.0, 1.0, 0.0, // 4
		0.5, -0.5, -0.5, 1.0, 1.0, 0.0, // 5
		0.5, -0.5, 0.5, 1.0, 1.0, 0.0, // 1
		0.5, -0.5, 0.5, 1.0, 1.0, 0.0, // 1
		-0.5, -0.5, 0.5, 1.0, 1.0, 0.0, // 0
		-0.5, -0.5, -0.5, 1.0, 1.0, 0.0, // 4

		// Right face (X+) - Cyan
		0.5, -0.5, -0.5, 0.0, 1.0, 1.0, // 5
		0.5, 0.5, -0.5, 0.0, 1.0, 1.0, // 6
		0.5, 0.5, 0.5, 0.0, 1.0, 1.0, // 2
		0.5, 0.5, 0.5, 0.0, 1.0, 1.0, // 2
		0.5, -0.5, 0.5, 0.0, 1.0, 1.0, // 1
		0.5, -0.5, -0.5, 0.0, 1.0, 1.0, // 5

		// Left face (X-) - Magenta
		-0.5, -0.5, -0.5, 1.0, 0.0, 1.0, // 4
		-0.5, -0.5, 0.5, 1.0, 0.0, 1.0, // 0
		-0.5, 0.5, 0.5, 1.0, 0.0, 1.0, // 3
		-0.5, 0.5, 0.5, 1.0, 0.0, 1.0, // 3
		-0.5, 0.5, -0.5, 1.0, 0.0, 1.0, // 7
		-0.5, -0.5, -0.5, 1.0, 0.0, 1.0, // 4
	}

	// nolint:gosec // len(vertices) is 216, * 4 = 864 bytes - no overflow risk
	vertexBufferSize := uint64(len(vertices) * 4) // 4 bytes per float32

	// Create buffer with MappedAtCreation = true for easy data upload
	app.vertexBuffer = app.device.CreateBuffer(&wgpu.BufferDescriptor{
		Label:            wgpu.StringView{},
		Usage:            wgpu.BufferUsageVertex | wgpu.BufferUsageCopyDst,
		Size:             vertexBufferSize,
		MappedAtCreation: wgpu.True,
	})

	if app.vertexBuffer == nil {
		return fmt.Errorf("failed to create vertex buffer")
	}

	// Copy vertex data to buffer
	ptr := app.vertexBuffer.GetMappedRange(0, vertexBufferSize)
	if ptr == nil {
		return fmt.Errorf("failed to get mapped range")
	}

	// Copy data using unsafe slice conversion
	// nolint:gosec,govet // ptr is from GetMappedRange, validated non-nil, safe for slice conversion
	mappedSlice := unsafe.Slice((*float32)(ptr), len(vertices))
	copy(mappedSlice, vertices)

	// Unmap buffer to commit data to GPU
	app.vertexBuffer.Unmap()

	return nil
}

// createUniformBuffer creates a uniform buffer for the MVP matrix.
func (app *App) createUniformBuffer() error {
	// Size: mat4x4f = 16 floats * 4 bytes = 64 bytes
	const uniformBufferSize = 64

	// Create buffer with Uniform usage and CopyDst for updates
	app.uniformBuffer = app.device.CreateBuffer(&wgpu.BufferDescriptor{
		Label:            wgpu.StringView{},
		Usage:            wgpu.BufferUsageUniform | wgpu.BufferUsageCopyDst,
		Size:             uniformBufferSize,
		MappedAtCreation: wgpu.False,
	})

	if app.uniformBuffer == nil {
		return fmt.Errorf("failed to create uniform buffer")
	}

	return nil
}

// createBindGroupLayout creates the bind group layout for the uniform buffer.
func (app *App) createBindGroupLayout() error {
	entries := []wgpu.BindGroupLayoutEntry{
		{
			Binding:    0,
			Visibility: wgpu.ShaderStageVertex,
			Buffer: wgpu.BufferBindingLayout{
				Type:             wgpu.BufferBindingTypeUniform,
				HasDynamicOffset: wgpu.False,
				MinBindingSize:   64, // mat4x4f = 64 bytes
			},
		},
	}

	app.bindGroupLayout = app.device.CreateBindGroupLayoutSimple(entries)
	if app.bindGroupLayout == nil {
		return fmt.Errorf("failed to create bind group layout")
	}

	return nil
}

// createBindGroup creates the bind group with the uniform buffer.
func (app *App) createBindGroup() error {
	entries := []wgpu.BindGroupEntry{
		wgpu.BufferBindingEntry(0, app.uniformBuffer, 0, 64),
	}

	app.bindGroup = app.device.CreateBindGroupSimple(app.bindGroupLayout, entries)
	if app.bindGroup == nil {
		return fmt.Errorf("failed to create bind group")
	}

	return nil
}

// createPipelineLayout creates the pipeline layout with bind group layout.
func (app *App) createPipelineLayout() (*wgpu.PipelineLayout, error) {
	pipelineLayout := app.device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
		Label:                wgpu.EmptyStringView(),
		BindGroupLayoutCount: 1,
		BindGroupLayouts:     app.bindGroupLayout.Handle(),
	})
	if pipelineLayout == nil {
		return nil, fmt.Errorf("failed to create pipeline layout")
	}
	return pipelineLayout, nil
}

// createPipeline creates the render pipeline with depth testing.
// nolint:funlen // Pipeline creation requires many configuration steps
func (app *App) createPipeline() error {
	// Create shader module
	shader := app.device.CreateShaderModuleWGSL(shaderSource)
	if shader == nil {
		return fmt.Errorf("failed to create shader module")
	}
	defer shader.Release()

	// Create pipeline layout
	pipelineLayout, err := app.createPipelineLayout()
	if err != nil {
		return err
	}
	defer pipelineLayout.Release()

	// Define vertex attributes
	attributes := []wgpu.VertexAttribute{
		{
			Format:         wgpu.VertexFormatFloat32x3, // position: vec3f
			Offset:         0,
			ShaderLocation: 0,
		},
		{
			Format:         wgpu.VertexFormatFloat32x3, // color: vec3f
			Offset:         12,                         // 3 floats * 4 bytes = 12 bytes offset
			ShaderLocation: 1,
		},
	}

	// Create render pipeline with vertex buffer layout and depth testing
	pipeline := app.device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label:  "",
		Layout: pipelineLayout,
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
			Buffers: []wgpu.VertexBufferLayout{{
				ArrayStride:    24, // 6 floats * 4 bytes = 24 bytes per vertex
				StepMode:       wgpu.VertexStepModeVertex,
				AttributeCount: 2,
				Attributes:     &attributes[0],
			}},
		},
		Primitive: wgpu.PrimitiveState{
			Topology:  wgpu.PrimitiveTopologyTriangleList,
			FrontFace: wgpu.FrontFaceCCW,
			CullMode:  wgpu.CullModeBack, // Enable back-face culling for cube
		},
		DepthStencil: &wgpu.DepthStencilState{
			Format:            wgpu.TextureFormatDepth24Plus,
			DepthWriteEnabled: true,
			DepthCompare:      wgpu.CompareFunctionLess,
			StencilFront: wgpu.StencilFaceState{
				Compare:     wgpu.CompareFunctionAlways,
				FailOp:      wgpu.StencilOperationKeep,
				DepthFailOp: wgpu.StencilOperationKeep,
				PassOp:      wgpu.StencilOperationKeep,
			},
			StencilBack: wgpu.StencilFaceState{
				Compare:     wgpu.CompareFunctionAlways,
				FailOp:      wgpu.StencilOperationKeep,
				DepthFailOp: wgpu.StencilOperationKeep,
				PassOp:      wgpu.StencilOperationKeep,
			},
			StencilReadMask:     0xFFFFFFFF,
			StencilWriteMask:    0xFFFFFFFF,
			DepthBias:           0,
			DepthBiasSlopeScale: 0.0,
			DepthBiasClamp:      0.0,
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

// updateUniformBuffer updates the MVP matrix for rotation.
func (app *App) updateUniformBuffer() {
	// Calculate rotation angles based on time
	elapsed := time.Since(app.startTime).Seconds()
	angleY := float32(elapsed * 0.5) // Rotate around Y axis
	angleX := float32(elapsed * 0.3) // Rotate around X axis (slower)

	// Build MVP matrix
	// Model matrix: rotation around Y and X axes
	modelY := wgpu.Mat4RotateY(angleY)
	modelX := wgpu.Mat4RotateX(angleX)
	model := modelY.Mul(modelX)

	// View matrix: camera at (0, 0, 3) looking at origin
	view := wgpu.Mat4LookAt(
		wgpu.Vec3{X: 0, Y: 0, Z: 3},
		wgpu.Vec3{X: 0, Y: 0, Z: 0},
		wgpu.Vec3{X: 0, Y: 1, Z: 0},
	)

	// Projection matrix: perspective with 45° FOV
	aspect := float32(app.width) / float32(app.height)
	projection := wgpu.Mat4Perspective(45.0*math.Pi/180.0, aspect, 0.1, 100.0)

	// Combine: MVP = Projection * View * Model
	mvp := projection.Mul(view).Mul(model)

	// Upload matrix to uniform buffer
	// nolint:gosec // mvp is a fixed 16-element array, size calculation is safe
	app.queue.WriteBufferRaw(
		app.uniformBuffer,
		0,
		unsafe.Pointer(&mvp[0]),
		64, // 16 floats * 4 bytes = 64 bytes
	)
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

// renderCube encodes the rotating cube rendering commands.
func (app *App) renderCube(encoder *wgpu.CommandEncoder, view *wgpu.TextureView) error {
	pass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		Label: "Cube Render Pass",
		ColorAttachments: []wgpu.RenderPassColorAttachment{{
			View:    view,
			LoadOp:  wgpu.LoadOpClear,
			StoreOp: wgpu.StoreOpStore,
			ClearValue: wgpu.Color{
				R: 0.1,
				G: 0.2,
				B: 0.3,
				A: 1.0,
			},
		}},
		DepthStencilAttachment: &wgpu.RenderPassDepthStencilAttachment{
			View:              app.depthTextureView,
			DepthLoadOp:       wgpu.LoadOpClear,
			DepthStoreOp:      wgpu.StoreOpStore,
			DepthClearValue:   1.0,
			DepthReadOnly:     false,
			StencilLoadOp:     wgpu.LoadOpClear,
			StencilStoreOp:    wgpu.StoreOpStore,
			StencilClearValue: 0,
			StencilReadOnly:   false,
		},
	})
	if pass == nil {
		return fmt.Errorf("failed to begin render pass")
	}
	defer pass.Release()

	pass.SetPipeline(app.pipeline)
	pass.SetBindGroup(0, app.bindGroup, nil)
	// Set vertex buffer: slot 0, buffer, offset 0, size = entire buffer
	pass.SetVertexBuffer(0, app.vertexBuffer, 0, uint64(36*6*4)) // 36 vertices * 6 floats * 4 bytes
	pass.Draw(36, 1, 0, 0)                                       // Draw 36 vertices (12 triangles)
	pass.End()
	return nil
}

// render draws a frame.
func (app *App) render() error {
	// Recreate surface if needed (e.g., after resize)
	if app.needsRecreate {
		if err := app.configureSurface(); err != nil {
			return fmt.Errorf("reconfigure surface: %w", err)
		}
		if err := app.recreateDepthTexture(); err != nil {
			return fmt.Errorf("recreate depth texture: %w", err)
		}
	}

	// Update uniform buffer with new rotation
	app.updateUniformBuffer()

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

	// Render cube
	if err := app.renderCube(encoder, app.surfaceTexView); err != nil {
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
	if app.depthTextureView != nil {
		app.depthTextureView.Release()
	}
	if app.depthTexture != nil {
		app.depthTexture.Release()
	}
	if app.bindGroup != nil {
		app.bindGroup.Release()
	}
	if app.bindGroupLayout != nil {
		app.bindGroupLayout.Release()
	}
	if app.uniformBuffer != nil {
		app.uniformBuffer.Release()
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
