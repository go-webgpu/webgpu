// Package main demonstrates using the Buffer Introspection API.
package main

import (
	"fmt"
	"log"

	"github.com/go-webgpu/webgpu/wgpu"
	"github.com/gogpu/gputypes"
)

func main() {
	// Create WebGPU instance
	instance, err := wgpu.CreateInstance(nil)
	if err != nil {
		log.Fatalf("Failed to create instance: %v", err)
	}
	defer instance.Release()

	// Request adapter
	adapter, err := instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		PowerPreference: gputypes.PowerPreferenceHighPerformance,
	})
	if err != nil {
		log.Fatalf("Failed to request adapter: %v", err)
	}
	defer adapter.Release()

	// Request device
	device, err := adapter.RequestDevice(nil)
	if err != nil {
		log.Fatalf("Failed to request device: %v", err)
	}
	defer device.Release()

	// Create a storage buffer
	bufferSize := uint64(1024 * 1024) // 1MB
	buffer := device.CreateBuffer(&wgpu.BufferDescriptor{
		Label: wgpu.EmptyStringView(),
		Usage: gputypes.BufferUsageStorage | gputypes.BufferUsageCopySrc | gputypes.BufferUsageCopyDst,
		Size:  bufferSize,
	})
	if buffer == nil {
		log.Fatal("Failed to create buffer")
	}
	defer buffer.Release()

	// Demonstrate buffer introspection
	fmt.Println("=== Buffer Introspection ===")

	// Get buffer size
	size := buffer.GetSize()
	fmt.Printf("Buffer size: %d bytes (%.2f MB)\n", size, float64(size)/(1024*1024))

	// Get buffer usage
	usage := buffer.GetUsage()
	fmt.Printf("Buffer usage: %s\n", usageToString(usage))

	// Get buffer map state
	mapState := buffer.GetMapState()
	fmt.Printf("Buffer map state: %s\n", mapStateToString(mapState))

	// Create a mappable buffer
	fmt.Println("\n=== Mappable Buffer Example ===")
	mappableBuffer := device.CreateBuffer(&wgpu.BufferDescriptor{
		Label:            wgpu.EmptyStringView(),
		Usage:            gputypes.BufferUsageMapRead | gputypes.BufferUsageCopyDst,
		Size:             1024,
		MappedAtCreation: wgpu.True,
	})
	if mappableBuffer == nil {
		log.Fatal("Failed to create mappable buffer")
	}
	defer mappableBuffer.Release()

	// Check state when mapped at creation
	mapState = mappableBuffer.GetMapState()
	fmt.Printf("Initial map state (MappedAtCreation): %s\n", mapStateToString(mapState))

	// Unmap the buffer
	mappableBuffer.Unmap()

	// Check state after unmap
	mapState = mappableBuffer.GetMapState()
	fmt.Printf("Map state after Unmap(): %s\n", mapStateToString(mapState))

	// Map async
	fmt.Println("\nMapping buffer asynchronously...")
	err = mappableBuffer.MapAsync(device, wgpu.MapModeRead, 0, 1024)
	if err != nil {
		log.Printf("MapAsync failed: %v", err)
	} else {
		mapState = mappableBuffer.GetMapState()
		fmt.Printf("Map state after MapAsync(): %s\n", mapStateToString(mapState))

		// Unmap again
		mappableBuffer.Unmap()
		mapState = mappableBuffer.GetMapState()
		fmt.Printf("Map state after final Unmap(): %s\n", mapStateToString(mapState))
	}

	fmt.Println("\n=== Buffer Lifecycle Demonstration ===")
	fmt.Println("Buffer introspection allows you to:")
	fmt.Println("- Query buffer size at runtime")
	fmt.Println("- Check which usage flags are set")
	fmt.Println("- Verify mapping state before operations")
	fmt.Println("- Debug buffer lifecycle issues")
}

func usageToString(usage gputypes.BufferUsage) string {
	var flags []string

	if usage&gputypes.BufferUsageMapRead != 0 {
		flags = append(flags, "MapRead")
	}
	if usage&gputypes.BufferUsageMapWrite != 0 {
		flags = append(flags, "MapWrite")
	}
	if usage&gputypes.BufferUsageCopySrc != 0 {
		flags = append(flags, "CopySrc")
	}
	if usage&gputypes.BufferUsageCopyDst != 0 {
		flags = append(flags, "CopyDst")
	}
	if usage&gputypes.BufferUsageIndex != 0 {
		flags = append(flags, "Index")
	}
	if usage&gputypes.BufferUsageVertex != 0 {
		flags = append(flags, "Vertex")
	}
	if usage&gputypes.BufferUsageUniform != 0 {
		flags = append(flags, "Uniform")
	}
	if usage&gputypes.BufferUsageStorage != 0 {
		flags = append(flags, "Storage")
	}
	if usage&gputypes.BufferUsageIndirect != 0 {
		flags = append(flags, "Indirect")
	}
	if usage&gputypes.BufferUsageQueryResolve != 0 {
		flags = append(flags, "QueryResolve")
	}

	if len(flags) == 0 {
		return "None"
	}

	result := flags[0]
	for i := 1; i < len(flags); i++ {
		result += " | " + flags[i]
	}
	return result
}

func mapStateToString(state wgpu.BufferMapState) string {
	switch state {
	case wgpu.BufferMapStateUnmapped:
		return "Unmapped"
	case wgpu.BufferMapStatePending:
		return "Pending"
	case wgpu.BufferMapStateMapped:
		return "Mapped"
	default:
		return fmt.Sprintf("Unknown (%d)", state)
	}
}
