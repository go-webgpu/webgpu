// Package main demonstrates using RenderPass Debug Markers for GPU debugging.
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

	// Request adapter
	adapter, err := instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		PowerPreference: wgpu.PowerPreferenceHighPerformance,
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

	// Create a simple texture for rendering
	texture, err := device.CreateTexture(&wgpu.TextureDescriptor{
		Label: "",
		Size: wgpu.Extent3D{
			Width:              800,
			Height:             600,
			DepthOrArrayLayers: 1,
		},
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     wgpu.TextureDimension2D,
		Format:        wgpu.TextureFormatRGBA8Unorm,
		Usage:         wgpu.TextureUsageRenderAttachment,
	})
	if err != nil {
		log.Fatalf("create texture: %v", err)
	}
	defer texture.Release()

	// Create texture view
	view, err := texture.CreateView(nil)
	if err != nil {
		log.Fatalf("create texture view: %v", err)
	}
	defer view.Release()

	// Create command encoder
	encoder, err := device.CreateCommandEncoder(&wgpu.CommandEncoderDescriptor{
		Label: "",
	})
	if err != nil {
		log.Fatalf("create command encoder: %v", err)
	}
	defer encoder.Release()

	fmt.Println("=== Demonstrating RenderPass Debug Markers ===")
	fmt.Println("These markers will appear in GPU debugging tools like:")
	fmt.Println("- RenderDoc")
	fmt.Println("- PIX (Windows)")
	fmt.Println("- Xcode GPU Debugger (macOS)")
	fmt.Println("- Chrome DevTools (WebGPU)")
	fmt.Println()

	// Begin render pass with debug annotations
	renderPass, err := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		Label: "Main Render Pass",
		ColorAttachments: []wgpu.RenderPassColorAttachment{
			{
				View:   view,
				LoadOp: wgpu.LoadOpClear,
				ClearValue: wgpu.Color{
					R: 0.1,
					G: 0.2,
					B: 0.3,
					A: 1.0,
				},
				StoreOp: wgpu.StoreOpStore,
			},
		},
	})
	if err != nil {
		log.Fatalf("begin render pass: %v", err)
	}
	defer renderPass.Release()

	// Insert a single debug marker
	renderPass.InsertDebugMarker("Frame Start")
	fmt.Println("Inserted debug marker: 'Frame Start'")

	// Push a debug group for scene rendering
	renderPass.PushDebugGroup("Scene Rendering")
	fmt.Println("Pushed debug group: 'Scene Rendering'")

	// Nested debug group for geometry
	renderPass.PushDebugGroup("Geometry Pass")
	fmt.Println("  Pushed nested group: 'Geometry Pass'")

	// In a real application, you would draw geometry here
	renderPass.InsertDebugMarker("Draw Opaque Objects")
	fmt.Println("    Inserted marker: 'Draw Opaque Objects'")

	renderPass.InsertDebugMarker("Draw Alpha-Tested Objects")
	fmt.Println("    Inserted marker: 'Draw Alpha-Tested Objects'")

	// Pop geometry pass
	renderPass.PopDebugGroup()
	fmt.Println("  Popped debug group: 'Geometry Pass'")

	// Another nested group for lighting
	renderPass.PushDebugGroup("Lighting Pass")
	fmt.Println("  Pushed nested group: 'Lighting Pass'")

	renderPass.InsertDebugMarker("Compute Shadow Maps")
	fmt.Println("    Inserted marker: 'Compute Shadow Maps'")

	renderPass.InsertDebugMarker("Apply Lighting")
	fmt.Println("    Inserted marker: 'Apply Lighting'")

	renderPass.PopDebugGroup()
	fmt.Println("  Popped debug group: 'Lighting Pass'")

	// Pop scene rendering
	renderPass.PopDebugGroup()
	fmt.Println("Popped debug group: 'Scene Rendering'")

	// Post-processing group
	renderPass.PushDebugGroup("Post-Processing")
	fmt.Println("Pushed debug group: 'Post-Processing'")

	renderPass.InsertDebugMarker("Tone Mapping")
	fmt.Println("  Inserted marker: 'Tone Mapping'")

	renderPass.InsertDebugMarker("FXAA")
	fmt.Println("  Inserted marker: 'FXAA'")

	renderPass.PopDebugGroup()
	fmt.Println("Popped debug group: 'Post-Processing'")

	// Final marker
	renderPass.InsertDebugMarker("Frame End")
	fmt.Println("Inserted debug marker: 'Frame End'")

	// End render pass
	renderPass.End()
	fmt.Println("\nRender pass ended")

	// Finish encoding
	commandBuffer, err := encoder.Finish()
	if err != nil {
		log.Fatalf("finish encoder: %v", err)
	}
	defer commandBuffer.Release()

	// Submit commands
	queue := device.Queue()
	defer queue.Release()
	_ = queue.Submit(commandBuffer)

	fmt.Println("Commands submitted")
	fmt.Println("\n=== Debug Marker Hierarchy ===")
	fmt.Println("Frame Start [marker]")
	fmt.Println("Scene Rendering [group]")
	fmt.Println("   Geometry Pass [group]")
	fmt.Println("      Draw Opaque Objects [marker]")
	fmt.Println("      Draw Alpha-Tested Objects [marker]")
	fmt.Println("   Lighting Pass [group]")
	fmt.Println("      Compute Shadow Maps [marker]")
	fmt.Println("      Apply Lighting [marker]")
	fmt.Println("Post-Processing [group]")
	fmt.Println("   Tone Mapping [marker]")
	fmt.Println("   FXAA [marker]")
	fmt.Println("Frame End [marker]")
	fmt.Println("\nThis hierarchy will be visible in GPU debugging tools!")
}
