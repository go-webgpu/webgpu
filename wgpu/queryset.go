package wgpu

import "unsafe"

// querySetDescriptor is the native structure for QuerySet descriptor (32 bytes).
type querySetDescriptor struct {
	nextInChain uintptr    // 8 bytes
	label       StringView // 16 bytes
	queryType   QueryType  // 4 bytes
	count       uint32     // 4 bytes
}

// QuerySetDescriptor describes a QuerySet to create.
type QuerySetDescriptor struct {
	Label string
	Type  QueryType
	Count uint32
}

// CreateQuerySet creates a new QuerySet for GPU profiling/timestamps.
func (d *Device) CreateQuerySet(desc *QuerySetDescriptor) *QuerySet {
	mustInit()
	if desc == nil {
		return nil
	}

	nativeDesc := querySetDescriptor{
		nextInChain: 0,
		label:       EmptyStringView(),
		queryType:   desc.Type,
		count:       desc.Count,
	}

	handle, _, _ := procDeviceCreateQuerySet.Call(
		d.handle,
		uintptr(unsafe.Pointer(&nativeDesc)),
	)
	if handle == 0 {
		return nil
	}
	trackResource(handle, "QuerySet")
	return &QuerySet{handle: handle}
}

// Destroy destroys the QuerySet, making it invalid.
func (qs *QuerySet) Destroy() {
	mustInit()
	if qs.handle != 0 {
		procQuerySetDestroy.Call(qs.handle) //nolint:errcheck
	}
}

// Release releases the QuerySet reference.
func (qs *QuerySet) Release() {
	if qs.handle != 0 {
		untrackResource(qs.handle)
		procQuerySetRelease.Call(qs.handle) //nolint:errcheck
		qs.handle = 0
	}
}

// Handle returns the underlying handle. For advanced use only.
func (qs *QuerySet) Handle() uintptr { return qs.handle }
