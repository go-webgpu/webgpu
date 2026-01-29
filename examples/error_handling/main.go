// Package main demonstrates WebGPU error handling using error scopes.
package main

import (
	"fmt"
	"log"

	"github.com/go-webgpu/webgpu/wgpu"
)

func main() {
	// Initialize WebGPU
	instance, err := wgpu.CreateInstance(nil)
	if err != nil {
		log.Fatalf("Failed to create instance: %v", err)
	}
	defer instance.Release()

	// Request adapter
	adapter, err := instance.RequestAdapter(nil)
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

	fmt.Println("WebGPU device created successfully")

	// Example 1: Basic error scope usage
	fmt.Println("\n=== Example 1: Basic Error Scope ===")
	device.PushErrorScope(wgpu.ErrorFilterValidation)

	// Perform some GPU operations (valid ones)
	queue := device.GetQueue()
	if queue == nil {
		log.Fatal("Failed to get queue")
	}
	defer queue.Release()
	fmt.Println("Queue obtained successfully")

	// Pop error scope and check for errors
	errType, message := device.PopErrorScope(instance)
	if errType != wgpu.ErrorTypeNoError {
		fmt.Printf("ERROR: Validation error occurred: %s\n", message)
	} else {
		fmt.Println("No errors captured")
	}

	// Example 2: Nested error scopes
	fmt.Println("\n=== Example 2: Nested Error Scopes ===")
	device.PushErrorScope(wgpu.ErrorFilterValidation) // Outer scope
	fmt.Println("Pushed outer validation scope")

	device.PushErrorScope(wgpu.ErrorFilterValidation) // Inner scope
	fmt.Println("Pushed inner validation scope")

	// Some operations...
	fmt.Println("Performing GPU operations...")

	// Pop inner scope first (LIFO)
	errType, message = device.PopErrorScope(instance)
	fmt.Printf("Inner scope: errType=%v, message=%q\n", errType, message)

	// Pop outer scope
	errType, message = device.PopErrorScope(instance)
	fmt.Printf("Outer scope: errType=%v, message=%q\n", errType, message)

	// Example 3: Different error filters
	fmt.Println("\n=== Example 3: Error Filters ===")

	// Catch validation errors
	device.PushErrorScope(wgpu.ErrorFilterValidation)
	fmt.Println("Monitoring for validation errors...")
	errType, _ = device.PopErrorScope(instance)
	fmt.Printf("Validation filter: errType=%v\n", errType)

	// Catch out-of-memory errors
	device.PushErrorScope(wgpu.ErrorFilterOutOfMemory)
	fmt.Println("Monitoring for out-of-memory errors...")
	errType, _ = device.PopErrorScope(instance)
	fmt.Printf("Out-of-memory filter: errType=%v\n", errType)

	// Catch internal errors
	device.PushErrorScope(wgpu.ErrorFilterInternal)
	fmt.Println("Monitoring for internal errors...")
	errType, _ = device.PopErrorScope(instance)
	fmt.Printf("Internal filter: errType=%v\n", errType)

	fmt.Println("\n=== Error Handling Examples Completed ===")

	// IMPORTANT NOTES:
	// 1. Always pop every pushed error scope
	// 2. Error scopes are LIFO (stack-based)
	// 3. Popping an empty stack will cause a panic (wgpu-native limitation)
	// 4. Use PopErrorScopeAsync if you need error handling for stack underflow
}
