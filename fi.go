package sys

import "os"

// FakeFI represents a fake os.FileInfo struct
type FakeFI struct {
	os.FileInfo
	isDir bool
}

// NewFakeFI returns a FakeFI.
func NewFakeFI() *FakeFI {
	return &FakeFI{}
}

// IsDir returns the isDir bool value on a FakeFI instance
func (ffi *FakeFI) IsDir() bool {
	return ffi.isDir
}

// SetIsDir sets the isDir bool value on a FakeFI instance
func (ffi *FakeFI) SetIsDir(isDir bool) {
	ffi.isDir = isDir
}
