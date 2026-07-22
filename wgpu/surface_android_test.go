//go:build android

package wgpu

var _ func(*Instance, uintptr) (*Surface, error) = (*Instance).CreateSurfaceFromAndroidNativeWindow
