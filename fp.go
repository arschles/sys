package sys

import "path/filepath"

// FP is the interface to a filepath.
type FP interface {

	// Walk walks the file tree rooted at root, calling walkFn for each file or directory in the tree, including root
	Walk(root string, walkFunc filepath.WalkFunc) error
}

// RealFP returns an FP object that behaves exactly like filepath (https://godoc.org/path/filepath)
func RealFP() FP {
	return &realFP{}
}

type realFP struct{}

func (r *realFP) Walk(root string, walkFunc filepath.WalkFunc) error {
	return filepath.Walk(root, walkFunc)
}

// FakeFP represents a fake filepath
type FakeFP struct {
	FakeFS
	walkInvoked bool
}

// NewFakeFP returns a FakeFP.
func NewFakeFP() *FakeFP {
	return &FakeFP{}
}

// Walk walks the file tree rooted at root, calling walkFn for each file or directory in the tree, including root.
// Additionally, it sets the f.walkInvoked bool as true
func (f *FakeFP) Walk(root string, walkFunc filepath.WalkFunc) error {
	f.walkInvoked = true
	return walkFunc(root, NewFakeFI(), nil)
}
