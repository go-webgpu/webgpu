// Package main demonstrates a colored triangle rendering using go-webgpu with vertex buffers.
// This example creates a window using Windows API and renders a triangle with red, green, and blue vertices.
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
	windowTitle  = "go-webgpu: Colored Triangle Example"
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
	width          uint32
	height         uint32
	running        bool
	needsRecreate  bool
	surfaceTex     *wgpu.SurfaceTexture
	surfaceTexView *wgpu.TextureView
}

// Shader source (WGSL) with vertex colors
const shaderSource = `
struct VertexInput {
    @location(0) position: vec2f,
    @location(1) color: vec3f,
}

struct VertexOutput {
    @builtin(position) position: vec4f,
    @location(0) color: vec3f,
}

@vertex
fn vs_main(in: VertexInput) -> VertexOutput {
    var out: VertexOutput;
    out.position = vec4f(in.position, 0.0, 1.0);
    out.color = in.color;
    return out;
}

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4f {
    return vec4f(in.color, 1.0);
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

	// Create vertex buffer
	if err := app.createVertexBuffer(); err != nil {
		return fmt.Errorf("create vertex buffer: %w", err)
	}

	// Create render pipeline
	if err := app.createPipeline(); err != nil {
		return fmt.Errorf("create pipeline: %w", err)
	}

	return nil
}

// createWindow creates the main window.
func (app *App) createWindow() error {
	className, err := windows.UTF16PtrFromString("GoWebGPUColoredTriangle")
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

// createVertexBuffer creates a vertex buffer with position and color data.
func (app *App) createVertexBuffer() error {
	// Vertex data: position (x, y) + color (r, g, b)
	vertices := []float32{
		// position (x, y)    color (r, g, b)
		0.0, 0.5, 1.0, 0.0, 0.0, // top - red
		-0.5, -0.5, 0.0, 1.0, 0.0, // bottom-left - green
		0.5, -0.5, 0.0, 0.0, 1.0, // bottom-right - blue
	}

	// nolint:gosec // len(vertices) is 15, * 4 = 60 bytes - no overflow risk
	vertexBufferSize := uint64(len(vertices) * 4) // 4 bytes per float32

	// Create buffer with MappedAtCreation = true for easy data upload
	app.vertexBuffer = app.device.CreateBuffer(&wgpu.BufferDescriptor{
		Label:            wgpu.StringView{},
		Usage:            gputypes.BufferUsageVertex | gputypes.BufferUsageCopyDst,
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

// createPipeline creates the render pipeline with vertex buffer layout.
func (app *App) createPipeline() error {
	// Create shader module
	shader := app.device.CreateShaderModuleWGSL(shaderSource)
	if shader == nil {
		return fmt.Errorf("failed to create shader module")
	}
	defer shader.Release()

	// Define vertex attributes
	attributes := []wgpu.VertexAttribute{
		{
			Format:         gputypes.VertexFormatFloat32x2, // position: vec2f
			Offset:         0,
			ShaderLocation: 0,
		},
		{
			Format:         gputypes.VertexFormatFloat32x3, // color: vec3f
			Offset:         8,                              // 2 floats * 4 bytes = 8 bytes offset
			ShaderLocation: 1,
		},
	}

	// Create render pipeline with vertex buffer layout
	pipeline := app.device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label:  "",
		Layout: nil, // auto layout
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
			Buffers: []wgpu.VertexBufferLayout{{
				ArrayStride:    20, // 5 floats * 4 bytes = 20 bytes per vertex
				StepMode:       gputypes.VertexStepModeVertex,
				AttributeCount: 2,
				Attributes:     &attributes[0],
			}},
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

// renderTriangle encodes the colored triangle rendering commands with vertex buffer.
func (app *App) renderTriangle(encoder *wgpu.CommandEncoder, view *wgpu.TextureView) error {
	pass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		Label: "Colored Triangle Render Pass",
		ColorAttachments: []wgpu.RenderPassColorAttachment{{
			View:    view,
			LoadOp:  gputypes.LoadOpClear,
			StoreOp: gputypes.StoreOpStore,
			ClearValue: wgpu.Color{
				R: 0.1,
				G: 0.2,
				B: 0.3,
				A: 1.0,
			},
		}},
	})
	if pass == nil {
		return fmt.Errorf("failed to begin render pass")
	}
	defer pass.Release()

	pass.SetPipeline(app.pipeline)
	// Set vertex buffer: slot 0, buffer, offset 0, size = entire buffer
	pass.SetVertexBuffer(0, app.vertexBuffer, 0, uint64(15*4)) // 15 floats * 4 bytes
	pass.Draw(3, 1, 0, 0)
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

	// Render triangle
	if err := app.renderTriangle(encoder, app.surfaceTexView); err != nil {
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
