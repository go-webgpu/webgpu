package wgpu

import (
	"testing"

	"github.com/gogpu/gputypes"
)

// TestNullGuard_Device_Creation tests nil device guards on creation methods.
func TestNullGuard_Device_Creation(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var d *Device

	t.Run("CreateCommandEncoder", func(t *testing.T) {
		result := d.CreateCommandEncoder(nil)
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("CreateBuffer", func(t *testing.T) {
		result := d.CreateBuffer(&BufferDescriptor{})
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("CreateTexture", func(t *testing.T) {
		result := d.CreateTexture(&TextureDescriptor{})
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("CreateShaderModuleWGSL", func(t *testing.T) {
		result := d.CreateShaderModuleWGSL("@vertex fn main() {}")
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("CreateSampler", func(t *testing.T) {
		result := d.CreateSampler(&SamplerDescriptor{})
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("CreateBindGroupLayout", func(t *testing.T) {
		result := d.CreateBindGroupLayout(&BindGroupLayoutDescriptor{})
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("CreateBindGroup", func(t *testing.T) {
		result := d.CreateBindGroup(&BindGroupDescriptor{})
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("CreatePipelineLayout", func(t *testing.T) {
		result := d.CreatePipelineLayout(&PipelineLayoutDescriptor{})
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("CreateComputePipeline", func(t *testing.T) {
		result := d.CreateComputePipeline(&ComputePipelineDescriptor{})
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("CreateRenderPipeline", func(t *testing.T) {
		result := d.CreateRenderPipeline(&RenderPipelineDescriptor{})
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("CreateQuerySet", func(t *testing.T) {
		result := d.CreateQuerySet(&QuerySetDescriptor{})
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("CreateRenderBundleEncoder", func(t *testing.T) {
		result := d.CreateRenderBundleEncoder(&RenderBundleEncoderDescriptor{})
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("GetQueue", func(t *testing.T) {
		result := d.GetQueue()
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})

	t.Run("CreateDepthTexture", func(t *testing.T) {
		result := d.CreateDepthTexture(100, 100, gputypes.TextureFormatDepth24Plus)
		if result != nil {
			t.Error("expected nil for nil device")
		}
	})
}

// TestNullGuard_Device_Void tests nil device guards on void methods.
func TestNullGuard_Device_Void(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var d *Device

	t.Run("Poll", func(t *testing.T) {
		d.Poll(true) // should not panic
	})

	t.Run("PushErrorScope", func(t *testing.T) {
		d.PushErrorScope(ErrorFilterValidation) // should not panic
	})
}

// TestNullGuard_Device_ZeroHandle tests zero-handle device guards.
func TestNullGuard_Device_ZeroHandle(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	d := &Device{handle: 0}

	t.Run("CreateCommandEncoder", func(t *testing.T) {
		if d.CreateCommandEncoder(nil) != nil {
			t.Error("expected nil for zero-handle device")
		}
	})

	t.Run("CreateBuffer", func(t *testing.T) {
		if d.CreateBuffer(&BufferDescriptor{}) != nil {
			t.Error("expected nil for zero-handle device")
		}
	})

	t.Run("GetQueue", func(t *testing.T) {
		if d.GetQueue() != nil {
			t.Error("expected nil for zero-handle device")
		}
	})

	t.Run("Poll", func(t *testing.T) {
		d.Poll(true) // should not panic
	})

	t.Run("PushErrorScope", func(t *testing.T) {
		d.PushErrorScope(ErrorFilterValidation) // should not panic
	})
}

// TestNullGuard_Instance tests nil instance guards.
func TestNullGuard_Instance(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var i *Instance

	t.Run("RequestAdapter", func(t *testing.T) {
		result, err := i.RequestAdapter(nil)
		if result != nil {
			t.Error("expected nil adapter for nil instance")
		}
		if err == nil {
			t.Error("expected error for nil instance")
		}
	})

	t.Run("ProcessEvents", func(t *testing.T) {
		i.ProcessEvents() // should not panic
	})
}

// TestNullGuard_Adapter tests nil adapter guards.
func TestNullGuard_Adapter(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var a *Adapter

	t.Run("RequestDevice", func(t *testing.T) {
		result, err := a.RequestDevice(nil)
		if result != nil {
			t.Error("expected nil device for nil adapter")
		}
		if err == nil {
			t.Error("expected error for nil adapter")
		}
	})
}

// TestNullGuard_CommandEncoder tests nil command encoder guards.
func TestNullGuard_CommandEncoder(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var enc *CommandEncoder

	t.Run("BeginComputePass", func(t *testing.T) {
		if enc.BeginComputePass(nil) != nil {
			t.Error("expected nil for nil encoder")
		}
	})

	t.Run("BeginRenderPass", func(t *testing.T) {
		if enc.BeginRenderPass(&RenderPassDescriptor{
			ColorAttachments: []RenderPassColorAttachment{{}},
		}) != nil {
			t.Error("expected nil for nil encoder")
		}
	})

	t.Run("Finish", func(t *testing.T) {
		if enc.Finish(nil) != nil {
			t.Error("expected nil for nil encoder")
		}
	})

	t.Run("CopyBufferToBuffer", func(t *testing.T) {
		enc.CopyBufferToBuffer(nil, 0, nil, 0, 0) // should not panic
	})

	t.Run("ClearBuffer", func(t *testing.T) {
		enc.ClearBuffer(nil, 0, 0) // should not panic
	})

	t.Run("CopyBufferToTexture", func(t *testing.T) {
		enc.CopyBufferToTexture(nil, nil, nil) // should not panic
	})

	t.Run("CopyTextureToBuffer", func(t *testing.T) {
		enc.CopyTextureToBuffer(nil, nil, nil) // should not panic
	})

	t.Run("CopyTextureToTexture", func(t *testing.T) {
		enc.CopyTextureToTexture(nil, nil, nil) // should not panic
	})

	t.Run("InsertDebugMarker", func(t *testing.T) {
		enc.InsertDebugMarker("test") // should not panic
	})

	t.Run("PushDebugGroup", func(t *testing.T) {
		enc.PushDebugGroup("test") // should not panic
	})

	t.Run("PopDebugGroup", func(t *testing.T) {
		enc.PopDebugGroup() // should not panic
	})

	t.Run("WriteTimestamp", func(t *testing.T) {
		enc.WriteTimestamp(nil, 0) // should not panic
	})

	t.Run("ResolveQuerySet", func(t *testing.T) {
		enc.ResolveQuerySet(nil, 0, 0, nil, 0) // should not panic
	})
}

// TestNullGuard_ComputePassEncoder tests nil compute pass encoder guards.
func TestNullGuard_ComputePassEncoder(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var cpe *ComputePassEncoder

	t.Run("SetPipeline", func(t *testing.T) {
		cpe.SetPipeline(nil) // should not panic
	})

	t.Run("SetBindGroup", func(t *testing.T) {
		cpe.SetBindGroup(0, nil, nil) // should not panic
	})

	t.Run("DispatchWorkgroups", func(t *testing.T) {
		cpe.DispatchWorkgroups(1, 1, 1) // should not panic
	})

	t.Run("DispatchWorkgroupsIndirect", func(t *testing.T) {
		cpe.DispatchWorkgroupsIndirect(nil, 0) // should not panic
	})

	t.Run("End", func(t *testing.T) {
		cpe.End() // should not panic
	})
}

// TestNullGuard_RenderPassEncoder tests nil render pass encoder guards.
func TestNullGuard_RenderPassEncoder(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var rpe *RenderPassEncoder

	t.Run("SetPipeline", func(t *testing.T) {
		rpe.SetPipeline(nil) // should not panic
	})

	t.Run("SetBindGroup", func(t *testing.T) {
		rpe.SetBindGroup(0, nil, nil) // should not panic
	})

	t.Run("SetVertexBuffer", func(t *testing.T) {
		rpe.SetVertexBuffer(0, nil, 0, 0) // should not panic
	})

	t.Run("SetIndexBuffer", func(t *testing.T) {
		rpe.SetIndexBuffer(nil, gputypes.IndexFormatUint16, 0, 0) // should not panic
	})

	t.Run("Draw", func(t *testing.T) {
		rpe.Draw(0, 0, 0, 0) // should not panic
	})

	t.Run("DrawIndexed", func(t *testing.T) {
		rpe.DrawIndexed(0, 0, 0, 0, 0) // should not panic
	})

	t.Run("DrawIndirect", func(t *testing.T) {
		rpe.DrawIndirect(nil, 0) // should not panic
	})

	t.Run("DrawIndexedIndirect", func(t *testing.T) {
		rpe.DrawIndexedIndirect(nil, 0) // should not panic
	})

	t.Run("SetViewport", func(t *testing.T) {
		rpe.SetViewport(0, 0, 100, 100, 0, 1) // should not panic
	})

	t.Run("SetScissorRect", func(t *testing.T) {
		rpe.SetScissorRect(0, 0, 100, 100) // should not panic
	})

	t.Run("SetBlendConstant", func(t *testing.T) {
		rpe.SetBlendConstant(&Color{1, 1, 1, 1}) // should not panic
	})

	t.Run("SetStencilReference", func(t *testing.T) {
		rpe.SetStencilReference(0) // should not panic
	})

	t.Run("InsertDebugMarker", func(t *testing.T) {
		rpe.InsertDebugMarker("test") // should not panic
	})

	t.Run("PushDebugGroup", func(t *testing.T) {
		rpe.PushDebugGroup("test") // should not panic
	})

	t.Run("PopDebugGroup", func(t *testing.T) {
		rpe.PopDebugGroup() // should not panic
	})

	t.Run("End", func(t *testing.T) {
		rpe.End() // should not panic
	})

	t.Run("ExecuteBundles", func(t *testing.T) {
		rpe.ExecuteBundles(nil) // should not panic
	})
}

// TestNullGuard_RenderBundleEncoder tests nil render bundle encoder guards.
func TestNullGuard_RenderBundleEncoder(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var rbe *RenderBundleEncoder

	t.Run("SetPipeline", func(t *testing.T) {
		rbe.SetPipeline(nil) // should not panic
	})

	t.Run("SetBindGroup", func(t *testing.T) {
		rbe.SetBindGroup(0, nil, nil) // should not panic
	})

	t.Run("SetVertexBuffer", func(t *testing.T) {
		rbe.SetVertexBuffer(0, nil, 0, 0) // should not panic
	})

	t.Run("SetIndexBuffer", func(t *testing.T) {
		rbe.SetIndexBuffer(nil, gputypes.IndexFormatUint16, 0, 0) // should not panic
	})

	t.Run("Draw", func(t *testing.T) {
		rbe.Draw(0, 0, 0, 0) // should not panic
	})

	t.Run("DrawIndexed", func(t *testing.T) {
		rbe.DrawIndexed(0, 0, 0, 0, 0) // should not panic
	})

	t.Run("DrawIndirect", func(t *testing.T) {
		rbe.DrawIndirect(nil, 0) // should not panic
	})

	t.Run("DrawIndexedIndirect", func(t *testing.T) {
		rbe.DrawIndexedIndirect(nil, 0) // should not panic
	})

	t.Run("Finish", func(t *testing.T) {
		if rbe.Finish(nil) != nil {
			t.Error("expected nil for nil encoder")
		}
	})
}

// TestNullGuard_Buffer tests nil buffer guards.
func TestNullGuard_Buffer(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var buf *Buffer

	t.Run("GetMappedRange", func(t *testing.T) {
		if buf.GetMappedRange(0, 0) != nil {
			t.Error("expected nil for nil buffer")
		}
	})

	t.Run("GetSize", func(t *testing.T) {
		if buf.GetSize() != 0 {
			t.Error("expected 0 for nil buffer")
		}
	})

	t.Run("Unmap", func(t *testing.T) {
		buf.Unmap() // should not panic
	})
}

// TestNullGuard_Texture tests nil texture guards.
func TestNullGuard_Texture(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var tex *Texture

	t.Run("CreateView", func(t *testing.T) {
		if tex.CreateView(nil) != nil {
			t.Error("expected nil for nil texture")
		}
	})
}

// TestNullGuard_QuerySet tests nil queryset guards.
func TestNullGuard_QuerySet(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var qs *QuerySet

	t.Run("Destroy", func(t *testing.T) {
		qs.Destroy() // should not panic
	})
}

// TestNullGuard_Queue tests nil queue guards.
func TestNullGuard_Queue(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var q *Queue

	t.Run("Submit", func(t *testing.T) {
		q.Submit(&CommandBuffer{}) // should not panic
	})

	t.Run("WriteBuffer", func(t *testing.T) {
		q.WriteBuffer(nil, 0, []byte{1, 2, 3}) // should not panic
	})

	t.Run("WriteBufferRaw", func(t *testing.T) {
		q.WriteBufferRaw(nil, 0, nil, 0) // should not panic
	})
}

// TestNullGuard_Surface tests nil surface guards.
func TestNullGuard_Surface(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var s *Surface

	t.Run("Configure", func(t *testing.T) {
		s.Configure(nil) // should not panic
	})

	t.Run("Unconfigure", func(t *testing.T) {
		s.Unconfigure() // should not panic
	})

	t.Run("GetCurrentTexture", func(t *testing.T) {
		result, _ := s.GetCurrentTexture()
		if result != nil {
			t.Error("expected nil for nil surface")
		}
	})

	t.Run("Present", func(t *testing.T) {
		s.Present() // should not panic
	})
}

// TestNullGuard_ComputePipeline tests nil compute pipeline guards.
func TestNullGuard_ComputePipeline(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var cp *ComputePipeline

	t.Run("GetBindGroupLayout", func(t *testing.T) {
		if cp.GetBindGroupLayout(0) != nil {
			t.Error("expected nil for nil pipeline")
		}
	})
}

// TestNullGuard_RenderPipeline tests nil render pipeline guards.
func TestNullGuard_RenderPipeline(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var rp *RenderPipeline

	t.Run("GetBindGroupLayout", func(t *testing.T) {
		if rp.GetBindGroupLayout(0) != nil {
			t.Error("expected nil for nil pipeline")
		}
	})
}

// TestNullGuard_PopErrorScopeAsync tests nil device in PopErrorScopeAsync.
func TestNullGuard_PopErrorScopeAsync(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	var d *Device
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	errType, msg, errResult := d.PopErrorScopeAsync(inst)
	if errResult == nil {
		t.Error("expected error for nil device")
	}
	if errType != ErrorTypeNoError {
		t.Errorf("expected ErrorTypeNoError, got %d", errType)
	}
	if msg != "" {
		t.Errorf("expected empty message, got %q", msg)
	}
}

// TestNullGuard_NilDesc tests nil descriptor guards on creation methods.
func TestNullGuard_NilDesc(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Use zero-handle device to test desc-nil paths that guard before FFI
	d := &Device{handle: 1} // fake non-zero handle

	t.Run("CreateBuffer_NilDesc", func(t *testing.T) {
		if d.CreateBuffer(nil) != nil {
			t.Error("expected nil for nil desc")
		}
	})

	t.Run("CreateTexture_NilDesc", func(t *testing.T) {
		if d.CreateTexture(nil) != nil {
			t.Error("expected nil for nil desc")
		}
	})

	t.Run("CreateSampler_NilDesc", func(t *testing.T) {
		if d.CreateSampler(nil) != nil {
			t.Error("expected nil for nil desc")
		}
	})

	t.Run("CreateBindGroupLayout_NilDesc", func(t *testing.T) {
		if d.CreateBindGroupLayout(nil) != nil {
			t.Error("expected nil for nil desc")
		}
	})

	t.Run("CreateBindGroup_NilDesc", func(t *testing.T) {
		if d.CreateBindGroup(nil) != nil {
			t.Error("expected nil for nil desc")
		}
	})

	t.Run("CreatePipelineLayout_NilDesc", func(t *testing.T) {
		if d.CreatePipelineLayout(nil) != nil {
			t.Error("expected nil for nil desc")
		}
	})

	t.Run("CreateComputePipeline_NilDesc", func(t *testing.T) {
		if d.CreateComputePipeline(nil) != nil {
			t.Error("expected nil for nil desc")
		}
	})

	t.Run("CreateRenderPipeline_NilDesc", func(t *testing.T) {
		if d.CreateRenderPipeline(nil) != nil {
			t.Error("expected nil for nil desc")
		}
	})

	t.Run("CreateQuerySet_NilDesc", func(t *testing.T) {
		if d.CreateQuerySet(nil) != nil {
			t.Error("expected nil for nil desc")
		}
	})

	t.Run("CreateRenderBundleEncoder_NilDesc", func(t *testing.T) {
		if d.CreateRenderBundleEncoder(nil) != nil {
			t.Error("expected nil for nil desc")
		}
	})
}
