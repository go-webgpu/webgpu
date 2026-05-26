package wgpu

import "unsafe"

// MappedRange provides safe access to a mapped buffer region.
// Obtained via [Buffer.MappedRange] after a successful [Buffer.Map] or
// [Buffer.MapAsync]. The data slice is invalidated by [Buffer.Unmap].
//
// Matches gogpu/wgpu MappedRange.
type MappedRange struct {
	data   unsafe.Pointer
	size   uint64
	offset uint64
	buf    *Buffer
}

// Bytes returns the mapped memory as a byte slice.
// Returns nil if the buffer is not mapped or the range has been invalidated.
//
// The slice is valid only while the buffer remains mapped. Calling
// [Buffer.Unmap] or [Buffer.Release] after this will cause undefined
// behavior if the slice is still accessed.
func (m *MappedRange) Bytes() []byte {
	if m == nil || m.data == nil {
		return nil
	}
	// Safety check: buffer must still exist and be mapped.
	if m.buf != nil && m.buf.MapState() != BufferMapStateMapped {
		return nil
	}
	return unsafe.Slice((*byte)(m.data), m.size)
}

// Len returns the size of the mapped range in bytes.
func (m *MappedRange) Len() int {
	if m == nil {
		return 0
	}
	return int(m.size)
}

// Offset returns the byte offset of this range within the buffer.
func (m *MappedRange) Offset() uint64 {
	if m == nil {
		return 0
	}
	return m.offset
}

// MappedRange returns a safe view over the mapped region [offset, offset+size).
//
// The buffer must be in the Mapped state ([Buffer.Map] or [Buffer.MapAsync]
// resolved to success). The returned MappedRange.Bytes() slice is invalidated
// by [Buffer.Unmap].
//
// Matches gogpu/wgpu Buffer.MappedRange(offset, size) (*MappedRange, error).
func (b *Buffer) MappedRange(offset, size uint64) (*MappedRange, error) {
	if b == nil || b.handle == 0 {
		return nil, &WGPUError{Op: "Buffer.MappedRange", Message: "buffer is nil or released"}
	}
	ptr := b.GetMappedRange(offset, size)
	if ptr == nil {
		return nil, &WGPUError{Op: "Buffer.MappedRange", Message: "buffer not mapped or invalid range"}
	}
	return &MappedRange{
		data:   ptr,
		size:   size,
		offset: offset,
		buf:    b,
	}, nil
}
