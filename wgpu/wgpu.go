package wgpu

import (
	"errors"
	"runtime"
	"sync"
)

var (
	wgpuLib  Library
	initOnce sync.Once
	initErr  error

	// Function pointers - Instance
	procCreateInstance        Proc
	procInstanceRelease       Proc
	procInstanceProcessEvents Proc

	// Function pointers - Adapter
	procAdapterRelease           Proc
	procInstanceRequestAdapter   Proc
	procAdapterRequestDevice     Proc
	procAdapterGetLimits         Proc
	procAdapterEnumerateFeatures Proc
	procAdapterHasFeature        Proc
	procAdapterGetInfo           Proc
	procAdapterInfoFreeMembers   Proc

	// Function pointers - Device
	procDeviceRelease        Proc
	procDeviceGetQueue       Proc
	procDeviceCreateBuffer   Proc
	procDevicePoll           Proc // wgpu-native extension
	procDevicePushErrorScope Proc
	procDevicePopErrorScope  Proc
	procDeviceGetFeatures    Proc
	procDeviceHasFeature     Proc
	procDeviceGetLimits      Proc

	// Function pointers - Queue
	procQueueRelease     Proc
	procQueueWriteBuffer Proc

	// Function pointers - Buffer
	procBufferRelease        Proc
	procBufferDestroy        Proc
	procBufferGetMappedRange Proc
	procBufferUnmap          Proc
	procBufferGetSize        Proc
	procBufferMapAsync       Proc
	procBufferGetUsage       Proc
	procBufferGetMapState    Proc

	// Function pointers - ShaderModule
	procDeviceCreateShaderModule Proc
	procShaderModuleRelease      Proc

	// Function pointers - BindGroupLayout
	procDeviceCreateBindGroupLayout Proc
	procBindGroupLayoutRelease      Proc

	// Function pointers - BindGroup
	procDeviceCreateBindGroup Proc
	procBindGroupRelease      Proc

	// Function pointers - PipelineLayout
	procDeviceCreatePipelineLayout Proc
	procPipelineLayoutRelease      Proc

	// Function pointers - ComputePipeline
	procDeviceCreateComputePipeline       Proc
	procComputePipelineGetBindGroupLayout Proc
	procComputePipelineRelease            Proc

	// Function pointers - CommandEncoder
	procDeviceCreateCommandEncoder         Proc
	procCommandEncoderBeginComputePass     Proc
	procCommandEncoderCopyBufferToBuffer   Proc
	procCommandEncoderCopyBufferToTexture  Proc
	procCommandEncoderCopyTextureToBuffer  Proc
	procCommandEncoderCopyTextureToTexture Proc
	procCommandEncoderClearBuffer          Proc
	procCommandEncoderInsertDebugMarker    Proc
	procCommandEncoderPushDebugGroup       Proc
	procCommandEncoderPopDebugGroup        Proc
	procCommandEncoderFinish               Proc
	procCommandEncoderRelease              Proc

	// Function pointers - ComputePassEncoder
	procComputePassEncoderSetPipeline                Proc
	procComputePassEncoderSetBindGroup               Proc
	procComputePassEncoderDispatchWorkgroups         Proc
	procComputePassEncoderDispatchWorkgroupsIndirect Proc
	procComputePassEncoderEnd                        Proc
	procComputePassEncoderRelease                    Proc

	// Function pointers - CommandBuffer
	procCommandBufferRelease Proc

	// Function pointers - Queue (additional)
	procQueueSubmit Proc

	// Function pointers - Surface
	procInstanceCreateSurface          Proc
	procSurfaceRelease                 Proc
	procSurfaceConfigure               Proc
	procSurfaceUnconfigure             Proc
	procSurfaceGetCapabilities         Proc
	procSurfaceCapabilitiesFreeMembers Proc
	procSurfaceGetCurrentTexture       Proc
	procSurfacePresent                 Proc

	// Function pointers - Texture
	procDeviceCreateTexture          Proc
	procTextureRelease               Proc
	procTextureDestroy               Proc
	procTextureCreateView            Proc
	procTextureViewRelease           Proc
	procTextureGetWidth              Proc
	procTextureGetHeight             Proc
	procTextureGetDepthOrArrayLayers Proc
	procTextureGetMipLevelCount      Proc
	procTextureGetFormat             Proc

	// Function pointers - Sampler
	procDeviceCreateSampler Proc
	procSamplerRelease      Proc

	// Function pointers - Queue (texture operations)
	procQueueWriteTexture Proc

	// Function pointers - RenderPass
	procCommandEncoderBeginRenderPass        Proc
	procRenderPassEncoderSetPipeline         Proc
	procRenderPassEncoderSetBindGroup        Proc
	procRenderPassEncoderSetVertexBuffer     Proc
	procRenderPassEncoderSetIndexBuffer      Proc
	procRenderPassEncoderDraw                Proc
	procRenderPassEncoderDrawIndexed         Proc
	procRenderPassEncoderDrawIndirect        Proc
	procRenderPassEncoderDrawIndexedIndirect Proc
	procRenderPassEncoderEnd                 Proc
	procRenderPassEncoderRelease             Proc
	procRenderPassEncoderSetViewport         Proc
	procRenderPassEncoderSetScissorRect      Proc
	procRenderPassEncoderSetBlendConstant    Proc
	procRenderPassEncoderSetStencilReference Proc
	procRenderPassEncoderInsertDebugMarker   Proc
	procRenderPassEncoderPushDebugGroup      Proc
	procRenderPassEncoderPopDebugGroup       Proc

	// Function pointers - RenderPipeline
	procDeviceCreateRenderPipeline       Proc
	procRenderPipelineRelease            Proc
	procRenderPipelineGetBindGroupLayout Proc

	// Function pointers - QuerySet
	procDeviceCreateQuerySet          Proc
	procQuerySetDestroy               Proc
	procQuerySetRelease               Proc
	procCommandEncoderWriteTimestamp  Proc
	procCommandEncoderResolveQuerySet Proc

	// Function pointers - RenderBundle
	procDeviceCreateRenderBundleEncoder        Proc
	procRenderBundleEncoderSetPipeline         Proc
	procRenderBundleEncoderSetBindGroup        Proc
	procRenderBundleEncoderSetVertexBuffer     Proc
	procRenderBundleEncoderSetIndexBuffer      Proc
	procRenderBundleEncoderDraw                Proc
	procRenderBundleEncoderDrawIndexed         Proc
	procRenderBundleEncoderDrawIndirect        Proc
	procRenderBundleEncoderDrawIndexedIndirect Proc
	procRenderBundleEncoderFinish              Proc
	procRenderBundleEncoderRelease             Proc
	procRenderBundleRelease                    Proc
	procRenderPassEncoderExecuteBundles        Proc
)

// Init initializes the wgpu library. Called automatically on first use.
// Can be called explicitly to check for initialization errors early.
func Init() error {
	initOnce.Do(func() {
		libPath := getLibraryPath()
		wgpuLib = loadLibrary(libPath)

		initSymbols()
	})
	return initErr
}

func getLibraryPath() string {
	switch runtime.GOOS {
	case "windows":
		return "wgpu_native.dll"
	case "darwin":
		return "libwgpu_native.dylib"
	default: // linux, freebsd, etc.
		return "libwgpu_native.so"
	}
}

func initSymbols() {
	// Instance
	procCreateInstance = wgpuLib.NewProc("wgpuCreateInstance")
	procInstanceRelease = wgpuLib.NewProc("wgpuInstanceRelease")
	procInstanceProcessEvents = wgpuLib.NewProc("wgpuInstanceProcessEvents")

	// Adapter
	procAdapterRelease = wgpuLib.NewProc("wgpuAdapterRelease")
	procInstanceRequestAdapter = wgpuLib.NewProc("wgpuInstanceRequestAdapter")
	procAdapterRequestDevice = wgpuLib.NewProc("wgpuAdapterRequestDevice")
	procAdapterGetLimits = wgpuLib.NewProc("wgpuAdapterGetLimits")
	procAdapterEnumerateFeatures = wgpuLib.NewProc("wgpuAdapterEnumerateFeatures")
	procAdapterHasFeature = wgpuLib.NewProc("wgpuAdapterHasFeature")
	procAdapterGetInfo = wgpuLib.NewProc("wgpuAdapterGetInfo")
	procAdapterInfoFreeMembers = wgpuLib.NewProc("wgpuAdapterInfoFreeMembers")

	// Device
	procDeviceRelease = wgpuLib.NewProc("wgpuDeviceRelease")
	procDeviceGetQueue = wgpuLib.NewProc("wgpuDeviceGetQueue")
	procDeviceCreateBuffer = wgpuLib.NewProc("wgpuDeviceCreateBuffer")
	procDevicePoll = wgpuLib.NewProc("wgpuDevicePoll") // wgpu-native extension
	procDevicePushErrorScope = wgpuLib.NewProc("wgpuDevicePushErrorScope")
	procDevicePopErrorScope = wgpuLib.NewProc("wgpuDevicePopErrorScope")
	procDeviceGetFeatures = wgpuLib.NewProc("wgpuDeviceGetFeatures")
	procDeviceHasFeature = wgpuLib.NewProc("wgpuDeviceHasFeature")
	procDeviceGetLimits = wgpuLib.NewProc("wgpuDeviceGetLimits")

	// Queue
	procQueueRelease = wgpuLib.NewProc("wgpuQueueRelease")
	procQueueWriteBuffer = wgpuLib.NewProc("wgpuQueueWriteBuffer")

	// Buffer
	procBufferRelease = wgpuLib.NewProc("wgpuBufferRelease")
	procBufferDestroy = wgpuLib.NewProc("wgpuBufferDestroy")
	procBufferGetMappedRange = wgpuLib.NewProc("wgpuBufferGetMappedRange")
	procBufferUnmap = wgpuLib.NewProc("wgpuBufferUnmap")
	procBufferGetSize = wgpuLib.NewProc("wgpuBufferGetSize")
	procBufferMapAsync = wgpuLib.NewProc("wgpuBufferMapAsync")
	procBufferGetUsage = wgpuLib.NewProc("wgpuBufferGetUsage")
	procBufferGetMapState = wgpuLib.NewProc("wgpuBufferGetMapState")

	// ShaderModule
	procDeviceCreateShaderModule = wgpuLib.NewProc("wgpuDeviceCreateShaderModule")
	procShaderModuleRelease = wgpuLib.NewProc("wgpuShaderModuleRelease")

	// BindGroupLayout
	procDeviceCreateBindGroupLayout = wgpuLib.NewProc("wgpuDeviceCreateBindGroupLayout")
	procBindGroupLayoutRelease = wgpuLib.NewProc("wgpuBindGroupLayoutRelease")

	// BindGroup
	procDeviceCreateBindGroup = wgpuLib.NewProc("wgpuDeviceCreateBindGroup")
	procBindGroupRelease = wgpuLib.NewProc("wgpuBindGroupRelease")

	// PipelineLayout
	procDeviceCreatePipelineLayout = wgpuLib.NewProc("wgpuDeviceCreatePipelineLayout")
	procPipelineLayoutRelease = wgpuLib.NewProc("wgpuPipelineLayoutRelease")

	// ComputePipeline
	procDeviceCreateComputePipeline = wgpuLib.NewProc("wgpuDeviceCreateComputePipeline")
	procComputePipelineGetBindGroupLayout = wgpuLib.NewProc("wgpuComputePipelineGetBindGroupLayout")
	procComputePipelineRelease = wgpuLib.NewProc("wgpuComputePipelineRelease")

	// CommandEncoder
	procDeviceCreateCommandEncoder = wgpuLib.NewProc("wgpuDeviceCreateCommandEncoder")
	procCommandEncoderBeginComputePass = wgpuLib.NewProc("wgpuCommandEncoderBeginComputePass")
	procCommandEncoderCopyBufferToBuffer = wgpuLib.NewProc("wgpuCommandEncoderCopyBufferToBuffer")
	procCommandEncoderCopyBufferToTexture = wgpuLib.NewProc("wgpuCommandEncoderCopyBufferToTexture")
	procCommandEncoderCopyTextureToBuffer = wgpuLib.NewProc("wgpuCommandEncoderCopyTextureToBuffer")
	procCommandEncoderCopyTextureToTexture = wgpuLib.NewProc("wgpuCommandEncoderCopyTextureToTexture")
	procCommandEncoderClearBuffer = wgpuLib.NewProc("wgpuCommandEncoderClearBuffer")
	procCommandEncoderInsertDebugMarker = wgpuLib.NewProc("wgpuCommandEncoderInsertDebugMarker")
	procCommandEncoderPushDebugGroup = wgpuLib.NewProc("wgpuCommandEncoderPushDebugGroup")
	procCommandEncoderPopDebugGroup = wgpuLib.NewProc("wgpuCommandEncoderPopDebugGroup")
	procCommandEncoderFinish = wgpuLib.NewProc("wgpuCommandEncoderFinish")
	procCommandEncoderRelease = wgpuLib.NewProc("wgpuCommandEncoderRelease")

	// ComputePassEncoder
	procComputePassEncoderSetPipeline = wgpuLib.NewProc("wgpuComputePassEncoderSetPipeline")
	procComputePassEncoderSetBindGroup = wgpuLib.NewProc("wgpuComputePassEncoderSetBindGroup")
	procComputePassEncoderDispatchWorkgroups = wgpuLib.NewProc("wgpuComputePassEncoderDispatchWorkgroups")
	procComputePassEncoderDispatchWorkgroupsIndirect = wgpuLib.NewProc("wgpuComputePassEncoderDispatchWorkgroupsIndirect")
	procComputePassEncoderEnd = wgpuLib.NewProc("wgpuComputePassEncoderEnd")
	procComputePassEncoderRelease = wgpuLib.NewProc("wgpuComputePassEncoderRelease")

	// CommandBuffer
	procCommandBufferRelease = wgpuLib.NewProc("wgpuCommandBufferRelease")

	// Queue (additional)
	procQueueSubmit = wgpuLib.NewProc("wgpuQueueSubmit")

	// Surface
	procInstanceCreateSurface = wgpuLib.NewProc("wgpuInstanceCreateSurface")
	procSurfaceRelease = wgpuLib.NewProc("wgpuSurfaceRelease")
	procSurfaceConfigure = wgpuLib.NewProc("wgpuSurfaceConfigure")
	procSurfaceUnconfigure = wgpuLib.NewProc("wgpuSurfaceUnconfigure")
	procSurfaceGetCapabilities = wgpuLib.NewProc("wgpuSurfaceGetCapabilities")
	procSurfaceCapabilitiesFreeMembers = wgpuLib.NewProc("wgpuSurfaceCapabilitiesFreeMembers")
	procSurfaceGetCurrentTexture = wgpuLib.NewProc("wgpuSurfaceGetCurrentTexture")
	procSurfacePresent = wgpuLib.NewProc("wgpuSurfacePresent")

	// Texture
	procDeviceCreateTexture = wgpuLib.NewProc("wgpuDeviceCreateTexture")
	procTextureRelease = wgpuLib.NewProc("wgpuTextureRelease")
	procTextureDestroy = wgpuLib.NewProc("wgpuTextureDestroy")
	procTextureCreateView = wgpuLib.NewProc("wgpuTextureCreateView")
	procTextureViewRelease = wgpuLib.NewProc("wgpuTextureViewRelease")
	procTextureGetWidth = wgpuLib.NewProc("wgpuTextureGetWidth")
	procTextureGetHeight = wgpuLib.NewProc("wgpuTextureGetHeight")
	procTextureGetDepthOrArrayLayers = wgpuLib.NewProc("wgpuTextureGetDepthOrArrayLayers")
	procTextureGetMipLevelCount = wgpuLib.NewProc("wgpuTextureGetMipLevelCount")
	procTextureGetFormat = wgpuLib.NewProc("wgpuTextureGetFormat")

	// Sampler
	procDeviceCreateSampler = wgpuLib.NewProc("wgpuDeviceCreateSampler")
	procSamplerRelease = wgpuLib.NewProc("wgpuSamplerRelease")

	// Queue (texture operations)
	procQueueWriteTexture = wgpuLib.NewProc("wgpuQueueWriteTexture")

	// RenderPass
	procCommandEncoderBeginRenderPass = wgpuLib.NewProc("wgpuCommandEncoderBeginRenderPass")
	procRenderPassEncoderSetPipeline = wgpuLib.NewProc("wgpuRenderPassEncoderSetPipeline")
	procRenderPassEncoderSetBindGroup = wgpuLib.NewProc("wgpuRenderPassEncoderSetBindGroup")
	procRenderPassEncoderSetVertexBuffer = wgpuLib.NewProc("wgpuRenderPassEncoderSetVertexBuffer")
	procRenderPassEncoderSetIndexBuffer = wgpuLib.NewProc("wgpuRenderPassEncoderSetIndexBuffer")
	procRenderPassEncoderDraw = wgpuLib.NewProc("wgpuRenderPassEncoderDraw")
	procRenderPassEncoderDrawIndexed = wgpuLib.NewProc("wgpuRenderPassEncoderDrawIndexed")
	procRenderPassEncoderDrawIndirect = wgpuLib.NewProc("wgpuRenderPassEncoderDrawIndirect")
	procRenderPassEncoderDrawIndexedIndirect = wgpuLib.NewProc("wgpuRenderPassEncoderDrawIndexedIndirect")
	procRenderPassEncoderEnd = wgpuLib.NewProc("wgpuRenderPassEncoderEnd")
	procRenderPassEncoderRelease = wgpuLib.NewProc("wgpuRenderPassEncoderRelease")
	procRenderPassEncoderSetViewport = wgpuLib.NewProc("wgpuRenderPassEncoderSetViewport")
	procRenderPassEncoderSetScissorRect = wgpuLib.NewProc("wgpuRenderPassEncoderSetScissorRect")
	procRenderPassEncoderSetBlendConstant = wgpuLib.NewProc("wgpuRenderPassEncoderSetBlendConstant")
	procRenderPassEncoderSetStencilReference = wgpuLib.NewProc("wgpuRenderPassEncoderSetStencilReference")
	procRenderPassEncoderInsertDebugMarker = wgpuLib.NewProc("wgpuRenderPassEncoderInsertDebugMarker")
	procRenderPassEncoderPushDebugGroup = wgpuLib.NewProc("wgpuRenderPassEncoderPushDebugGroup")
	procRenderPassEncoderPopDebugGroup = wgpuLib.NewProc("wgpuRenderPassEncoderPopDebugGroup")

	// RenderPipeline
	procDeviceCreateRenderPipeline = wgpuLib.NewProc("wgpuDeviceCreateRenderPipeline")
	procRenderPipelineRelease = wgpuLib.NewProc("wgpuRenderPipelineRelease")
	procRenderPipelineGetBindGroupLayout = wgpuLib.NewProc("wgpuRenderPipelineGetBindGroupLayout")

	// QuerySet
	procDeviceCreateQuerySet = wgpuLib.NewProc("wgpuDeviceCreateQuerySet")
	procQuerySetDestroy = wgpuLib.NewProc("wgpuQuerySetDestroy")
	procQuerySetRelease = wgpuLib.NewProc("wgpuQuerySetRelease")
	procCommandEncoderWriteTimestamp = wgpuLib.NewProc("wgpuCommandEncoderWriteTimestamp")
	procCommandEncoderResolveQuerySet = wgpuLib.NewProc("wgpuCommandEncoderResolveQuerySet")

	// RenderBundle
	procDeviceCreateRenderBundleEncoder = wgpuLib.NewProc("wgpuDeviceCreateRenderBundleEncoder")
	procRenderBundleEncoderSetPipeline = wgpuLib.NewProc("wgpuRenderBundleEncoderSetPipeline")
	procRenderBundleEncoderSetBindGroup = wgpuLib.NewProc("wgpuRenderBundleEncoderSetBindGroup")
	procRenderBundleEncoderSetVertexBuffer = wgpuLib.NewProc("wgpuRenderBundleEncoderSetVertexBuffer")
	procRenderBundleEncoderSetIndexBuffer = wgpuLib.NewProc("wgpuRenderBundleEncoderSetIndexBuffer")
	procRenderBundleEncoderDraw = wgpuLib.NewProc("wgpuRenderBundleEncoderDraw")
	procRenderBundleEncoderDrawIndexed = wgpuLib.NewProc("wgpuRenderBundleEncoderDrawIndexed")
	procRenderBundleEncoderDrawIndirect = wgpuLib.NewProc("wgpuRenderBundleEncoderDrawIndirect")
	procRenderBundleEncoderDrawIndexedIndirect = wgpuLib.NewProc("wgpuRenderBundleEncoderDrawIndexedIndirect")
	procRenderBundleEncoderFinish = wgpuLib.NewProc("wgpuRenderBundleEncoderFinish")
	procRenderBundleEncoderRelease = wgpuLib.NewProc("wgpuRenderBundleEncoderRelease")
	procRenderBundleRelease = wgpuLib.NewProc("wgpuRenderBundleRelease")
	procRenderPassEncoderExecuteBundles = wgpuLib.NewProc("wgpuRenderPassEncoderExecuteBundles")
}

// ErrLibraryNotLoaded is returned when wgpu-native library is not loaded or failed to initialize.
var ErrLibraryNotLoaded = errors.New("wgpu: native library not loaded or failed to initialize")

// checkInit checks that the library is initialized, returning error if not.
func checkInit() error {
	if err := Init(); err != nil {
		return ErrLibraryNotLoaded
	}
	return nil
}

func mustInit() {
	if err := Init(); err != nil {
		panic(err)
	}
}
