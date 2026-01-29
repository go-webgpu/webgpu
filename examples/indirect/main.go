// Package main demonstrates GPU-driven rendering using DrawIndirect.
// The GPU reads draw parameters (vertex count, instance count) from a buffer.
package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/go-webgpu/webgpu/wgpu"
	"github.com/gogpu/gputypes"
	"golang.org/x/sys/windows"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "go-webgpu: Indirect Draw Example"
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

// Vertex with position and color
type Vertex struct {
	Position [2]float32
	Color    [3]float32
}

// Application state
type App struct {
	hwnd            windows.HWND
	hinstance       windows.Handle
	instance        *wgpu.Instance
	adapter         *wgpu.Adapter
	device          *wgpu.Device
	queue           *wgpu.Queue
	surface         *wgpu.Surface
	pipeline        *wgpu.RenderPipeline
	vertexBuffer    *wgpu.Buffer
	indirectBuffer  *wgpu.Buffer
	width           uint32
	height          uint32
	running         bool
	needsRecreate   bool
	surfaceTex      *wgpu.SurfaceTexture
	surfaceTexView  *wgpu.TextureView
	vertexBufSize   uint64
	indirectBufSize uint64
}

// Shader source (WGSL) with instancing
const shaderSource = `
struct VertexInput {
    @location(0) position: vec2<f32>,
    @location(1) color: vec3<f32>,
};

struct VertexOutput {
    @builtin(position) position: vec4<f32>,
    @location(0) color: vec3<f32>,
};

@vertex
fn vs_main(input: VertexInput, @builtin(instance_index) instance_idx: u32) -> VertexOutput {
    var output: VertexOutput;
    // Offset each instance horizontally
    let offset = f32(instance_idx) * 0.4 - 0.6;
    output.position = vec4<f32>(input.position.x * 0.3 + offset, input.position.y * 0.3, 0.0, 1.0);
    output.color = input.color;
    return output;
}

@fragment
fn fs_main(@location(0) color: vec3<f32>) -> @location(0) vec4<f32> {
    return vec4<f32>(color, 1.0);
}
`

func main() {
	fmt.Println("=== GPU-Driven Rendering (DrawIndirect) ===")
	fmt.Println()
	fmt.Println("This example demonstrates DrawIndirect where the GPU")
	fmt.Println("reads draw parameters from a buffer.")
	fmt.Println()
	fmt.Println("4 triangles are rendered using a single DrawIndirect call.")
	fmt.Println("In real applications, a compute shader could dynamically")
	fmt.Println("update the indirect buffer for GPU-driven culling/LOD.")
	fmt.Println()

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
	// Initialize WebGPU library
	if err := wgpu.Init(); err != nil {
		return fmt.Errorf("init wgpu: %w", err)
	}

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

	// Create buffers
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
	className, err := windows.UTF16PtrFromString("GoWebGPUIndirect")
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

	ret, _, _ := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wndClass)))
	if ret == 0 {
		return fmt.Errorf("RegisterClassExW failed")
	}

	titlePtr, err := windows.UTF16PtrFromString(windowTitle)
	if err != nil {
		return err
	}

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
		Format:      gputypes.TextureFormatBGRA8Unorm,
		Usage:       gputypes.TextureUsageRenderAttachment,
		Width:       app.width,
		Height:      app.height,
		AlphaMode:   gputypes.CompositeAlphaModeOpaque,
		PresentMode: gputypes.PresentModeFifo,
	})
	app.needsRecreate = false
	return nil
}

// createBuffers creates the vertex and indirect buffers.
func (app *App) createBuffers() error {
	// Triangle vertices with colors
	vertices := []Vertex{
		{Position: [2]float32{0.0, 0.5}, Color: [3]float32{1.0, 0.0, 0.0}},   // Top - Red
		{Position: [2]float32{-0.5, -0.5}, Color: [3]float32{0.0, 1.0, 0.0}}, // Bottom left - Green
		{Position: [2]float32{0.5, -0.5}, Color: [3]float32{0.0, 0.0, 1.0}},  // Bottom right - Blue
	}

	// Create vertex buffer
	app.vertexBufSize = uint64(len(vertices) * int(unsafe.Sizeof(Vertex{})))
	app.vertexBuffer = app.device.CreateBuffer(&wgpu.BufferDescriptor{
		Usage:            gputypes.BufferUsageVertex | gputypes.BufferUsageCopyDst,
		Size:             app.vertexBufSize,
		MappedAtCreation: wgpu.True,
	})
	if app.vertexBuffer == nil {
		return fmt.Errorf("failed to create vertex buffer")
	}

	// Copy vertex data
	ptr := app.vertexBuffer.GetMappedRange(0, app.vertexBufSize)
	if ptr != nil {
		mappedSlice := unsafe.Slice((*Vertex)(ptr), len(vertices))
		copy(mappedSlice, vertices)
	}
	app.vertexBuffer.Unmap()

	// Create indirect buffer with draw arguments
	// This is the key part - the GPU reads these parameters!
	indirectArgs := wgpu.DrawIndirectArgs{
		VertexCount:   3, // Draw 3 vertices (triangle)
		InstanceCount: 4, // Draw 4 instances
		FirstVertex:   0,
		FirstInstance: 0,
	}

	fmt.Printf("Indirect buffer contents:\n")
	fmt.Printf("  VertexCount:   %d\n", indirectArgs.VertexCount)
	fmt.Printf("  InstanceCount: %d\n", indirectArgs.InstanceCount)
	fmt.Printf("  FirstVertex:   %d\n", indirectArgs.FirstVertex)
	fmt.Printf("  FirstInstance: %d\n", indirectArgs.FirstInstance)
	fmt.Println()

	app.indirectBufSize = uint64(unsafe.Sizeof(indirectArgs))
	app.indirectBuffer = app.device.CreateBuffer(&wgpu.BufferDescriptor{
		Usage:            gputypes.BufferUsageIndirect | gputypes.BufferUsageCopyDst,
		Size:             app.indirectBufSize,
		MappedAtCreation: wgpu.True,
	})
	if app.indirectBuffer == nil {
		return fmt.Errorf("failed to create indirect buffer")
	}

	// Copy indirect args to buffer
	indirectPtr := app.indirectBuffer.GetMappedRange(0, app.indirectBufSize)
	if indirectPtr != nil {
		*(*wgpu.DrawIndirectArgs)(indirectPtr) = indirectArgs
	}
	app.indirectBuffer.Unmap()

	return nil
}

// createPipeline creates the render pipeline.
func (app *App) createPipeline() error {
	// Create shader module
	shader := app.device.CreateShaderModuleWGSL(shaderSource)
	if shader == nil {
		return fmt.Errorf("failed to create shader module")
	}
	defer shader.Release()

	// Vertex attributes
	vertexAttributes := []wgpu.VertexAttribute{
		{
			Format:         gputypes.VertexFormatFloat32x2,
			Offset:         0,
			ShaderLocation: 0, // position
		},
		{
			Format:         gputypes.VertexFormatFloat32x3,
			Offset:         8, // After position (2 * 4 bytes)
			ShaderLocation: 1, // color
		},
	}

	vertexBufferLayout := wgpu.VertexBufferLayout{
		ArrayStride:    uint64(unsafe.Sizeof(Vertex{})),
		StepMode:       gputypes.VertexStepModeVertex,
		AttributeCount: uintptr(len(vertexAttributes)),
		Attributes:     &vertexAttributes[0],
	}

	// Create render pipeline
	pipeline := app.device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
			Buffers:    []wgpu.VertexBufferLayout{vertexBufferLayout},
		},
		Primitive: wgpu.PrimitiveState{
			Topology:  gputypes.PrimitiveTopologyTriangleList,
			FrontFace: gputypes.FrontFaceCCW,
			CullMode:  gputypes.CullModeNone,
		},
		Multisample: wgpu.MultisampleState{
			Count: 1,
			Mask:  0xFFFFFFFF,
		},
		Fragment: &wgpu.FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{{
				Format:    gputypes.TextureFormatBGRA8Unorm,
				WriteMask: gputypes.ColorWriteMaskAll,
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

// render draws a frame.
func (app *App) render() error {
	// Recreate surface if needed (e.g., after resize)
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
	if app.surfaceTexView == nil {
		return nil // Surface was lost, will retry next frame
	}

	// Create command encoder
	encoder := app.device.CreateCommandEncoder(nil)
	if encoder == nil {
		return fmt.Errorf("failed to create command encoder")
	}
	defer encoder.Release()

	// Begin render pass
	pass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		Label: "Indirect Draw Pass",
		ColorAttachments: []wgpu.RenderPassColorAttachment{{
			View:    app.surfaceTexView,
			LoadOp:  gputypes.LoadOpClear,
			StoreOp: gputypes.StoreOpStore,
			ClearValue: wgpu.Color{
				R: 0.1,
				G: 0.1,
				B: 0.15,
				A: 1.0,
			},
		}},
	})
	if pass == nil {
		return fmt.Errorf("failed to begin render pass")
	}
	defer pass.Release()

	pass.SetPipeline(app.pipeline)
	pass.SetVertexBuffer(0, app.vertexBuffer, 0, app.vertexBufSize)

	// GPU-driven draw call - parameters read from buffer!
	pass.DrawIndirect(app.indirectBuffer, 0)

	pass.End()

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
			_, _, _ = procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
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
func (app *App) cleanup() {
	if app.surfaceTexView != nil {
		app.surfaceTexView.Release()
	}
	if app.surfaceTex != nil && app.surfaceTex.Texture != nil {
		app.surfaceTex.Texture.Release()
	}
	if app.indirectBuffer != nil {
		app.indirectBuffer.Release()
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
