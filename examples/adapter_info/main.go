// Package main demonstrates using the Adapter Info API to query GPU capabilities.
package main

import (
	"fmt"
	"log"

	"github.com/go-webgpu/webgpu/wgpu"
)

func main() {
	// Create WebGPU instance
	instance, err := wgpu.CreateInstance(nil)
	if err != nil {
		log.Fatalf("Failed to create instance: %v", err)
	}
	defer instance.Release()

	// Request adapter (GPU)
	adapter, err := instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		PowerPreference: wgpu.PowerPreferenceHighPerformance,
	})
	if err != nil {
		log.Fatalf("Failed to request adapter: %v", err)
	}
	defer adapter.Release()

	// Get adapter information
	info, err := adapter.GetInfo()
	if err != nil {
		log.Fatalf("Failed to get adapter info: %v", err)
	}

	fmt.Println("=== Adapter Information ===")
	fmt.Printf("Vendor:       %s\n", info.Vendor)
	fmt.Printf("Device:       %s\n", info.Device)
	fmt.Printf("Description:  %s\n", info.Description)
	fmt.Printf("Architecture: %s\n", info.Architecture)
	fmt.Printf("Backend Type: %v\n", backendTypeToString(info.BackendType))
	fmt.Printf("Adapter Type: %v\n", adapterTypeToString(info.AdapterType))
	fmt.Printf("Vendor ID:    0x%04X\n", info.VendorID)
	fmt.Printf("Device ID:    0x%04X\n", info.DeviceID)

	// Get adapter limits
	limits, err := adapter.GetLimits()
	if err != nil {
		log.Fatalf("Failed to get adapter limits: %v", err)
	}

	fmt.Println("\n=== Key Adapter Limits ===")
	fmt.Printf("Max Texture 2D:              %d x %d\n", limits.Limits.MaxTextureDimension2D, limits.Limits.MaxTextureDimension2D)
	fmt.Printf("Max Texture 3D:              %d x %d x %d\n", limits.Limits.MaxTextureDimension3D, limits.Limits.MaxTextureDimension3D, limits.Limits.MaxTextureDimension3D)
	fmt.Printf("Max Bind Groups:             %d\n", limits.Limits.MaxBindGroups)
	fmt.Printf("Max Buffer Size:             %d bytes (%.2f GB)\n", limits.Limits.MaxBufferSize, float64(limits.Limits.MaxBufferSize)/(1024*1024*1024))
	fmt.Printf("Max Uniform Buffer Size:     %d bytes (%.2f MB)\n", limits.Limits.MaxUniformBufferBindingSize, float64(limits.Limits.MaxUniformBufferBindingSize)/(1024*1024))
	fmt.Printf("Max Storage Buffer Size:     %d bytes (%.2f GB)\n", limits.Limits.MaxStorageBufferBindingSize, float64(limits.Limits.MaxStorageBufferBindingSize)/(1024*1024*1024))
	fmt.Printf("Max Compute Workgroup Size:  %d x %d x %d\n", limits.Limits.MaxComputeWorkgroupSizeX, limits.Limits.MaxComputeWorkgroupSizeY, limits.Limits.MaxComputeWorkgroupSizeZ)
	fmt.Printf("Max Vertex Buffers:          %d\n", limits.Limits.MaxVertexBuffers)
	fmt.Printf("Max Color Attachments:       %d\n", limits.Limits.MaxColorAttachments)

	// Enumerate features
	features := adapter.EnumerateFeatures()
	fmt.Println("\n=== Supported Features ===")
	if len(features) == 0 {
		fmt.Println("No optional features supported")
	} else {
		for _, feature := range features {
			fmt.Printf("- %v\n", featureToString(feature))
		}
	}

	// Check for specific feature
	hasTimestamps := adapter.HasFeature(wgpu.FeatureNameTimestampQuery)
	fmt.Printf("\nTimestamp queries supported: %v\n", hasTimestamps)
}

func backendTypeToString(bt wgpu.BackendType) string {
	switch bt {
	case wgpu.BackendTypeD3D12:
		return "D3D12"
	case wgpu.BackendTypeVulkan:
		return "Vulkan"
	case wgpu.BackendTypeMetal:
		return "Metal"
	case wgpu.BackendTypeOpenGL:
		return "OpenGL"
	case wgpu.BackendTypeOpenGLES:
		return "OpenGL ES"
	case wgpu.BackendTypeWebGPU:
		return "WebGPU"
	default:
		return fmt.Sprintf("Unknown (%d)", bt)
	}
}

func adapterTypeToString(at wgpu.AdapterType) string {
	switch at {
	case wgpu.AdapterTypeDiscreteGPU:
		return "Discrete GPU"
	case wgpu.AdapterTypeIntegratedGPU:
		return "Integrated GPU"
	case wgpu.AdapterTypeCPU:
		return "CPU"
	default:
		return fmt.Sprintf("Unknown (%d)", at)
	}
}

func featureToString(f wgpu.FeatureName) string {
	switch f {
	case wgpu.FeatureNameTimestampQuery:
		return "Timestamp Query"
	default:
		return fmt.Sprintf("Feature %d", f)
	}
}
