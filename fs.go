package sys

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// FS is the interface to a file system.
type FS interface {
	// ReadAll gets the contents of filename, or an error if the file didn't exist or there was an
	// error reading it.
	ReadFile(filename string) ([]byte, error)

	// RemoveAll removes all of the files under the directory at name. It behaves similarly to the
	// func of the same name in the os package (https://godoc.org/os#RemoveAll).
	RemoveAll(name string) error

	// Create invokes the func of the same name in the os package (https://godoc.org/os#Create).
	Create(string) (io.WriteCloser, error)

	// Stat invokes the func of the same name in the os package (https://godoc.org/os#Stat).
	Stat(string) (os.FileInfo, error)

	// MkdirAll invokes the func of the same name in the os package (https://godoc.org/os#MkdirAll).
	MkdirAll(string, os.FileMode) error

	// WriteFile invokes the func of the same name in the os package (https://godoc.org/io/ioutil#WriteFile).
	WriteFile(string, []byte, os.FileMode) (int, error)
}

// RealFS returns an FS object that interacts with the real local filesystem.
func RealFS() FS {
	return &realFS{}
}

type realFS struct{}

// ReadFile is the interface implementation for FS.
func (r *realFS) ReadFile(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}

// RemoveAll is the interface implementation for FS.
func (r *realFS) RemoveAll(name string) error {
	return os.RemoveAll(name)
}

// Stat is the interface implementation for FS.
func (r *realFS) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// MkdirAll is the interface implementation for FS.
func (r *realFS) MkdirAll(dirName string, perm os.FileMode) error {
	return os.MkdirAll(dirName, perm)
}

// Create is the interface implementation for FS.
func (r *realFS) Create(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

// WriteFile is the interface implementation for FS.
func (r *realFS) WriteFile(filename string, data []byte, perm os.FileMode) (int, error) {
	return len(data), ioutil.WriteFile(filename, data, perm)
}

// FakeFileNotFound is the error returned by FakeFS when a requested file isn't found.
type FakeFileNotFound struct {
	Filename string
}

// Error is the error interface implementation.
func (f FakeFileNotFound) Error() string {
	return fmt.Sprintf("Fake file %s not found", f.Filename)
}

// FakeFS is an in-memory FS implementation.
type FakeFS struct {
	Files map[string]*bytes.Buffer
}

// NewFakeFS returns a FakeFS with no files.
func NewFakeFS() *FakeFS {
	return &FakeFS{Files: make(map[string]*bytes.Buffer)}
}

type inMemoryCloser struct {
	buf *bytes.Buffer
}

func (i inMemoryCloser) Write(b []byte) (int, error) {
	return i.buf.Write(b)
}

func (i inMemoryCloser) Close() error {
	return nil
}

// ReadFile is the FS interface implementation. It returns FakeFileNotFound if the file was not
// found in the in-memory 'filesystem' of f.
func (f *FakeFS) ReadFile(name string) ([]byte, error) {
	buf, ok := f.Files[name]
	if !ok {
		return nil, FakeFileNotFound{Filename: name}
	}
	return buf.Bytes(), nil
}

// RemoveAll is the interface implementation for FS.
func (f *FakeFS) RemoveAll(name string) error {
	_, ok := f.Files[name]
	if !ok {
		return FakeFileNotFound{Filename: name}
	}
	delete(f.Files, name)
	return nil
}

// Stat is the interface implementation for FS.  It returns os.ErrNotExist if the file was not
// found in the in-memory 'filesystem' of f
func (f *FakeFS) Stat(path string) (os.FileInfo, error) {
	_, err := f.ReadFile(path)
	if err != nil {
		return nil, os.ErrNotExist
	}
	return NewFakeFI(), nil
}

// MkdirAll is the interface implementation for FS.
func (f *FakeFS) MkdirAll(dirName string, perm os.FileMode) error {
	_, err := f.Create(dirName)
	return err
}

// Create is the interface implementation for FS.  It populates an entry in f.Files for path
// with an empty byte array and returns an empty os.File struct.
func (f *FakeFS) Create(path string) (io.WriteCloser, error) {
	buf := new(bytes.Buffer)
	f.Files[path] = buf
	return inMemoryCloser{buf: buf}, nil
}

// WriteFile is the interface implementation for FS.  To properly emulate WriteFile, it
// creates a new bytes.Buffer for the value Files[filename] references, then writes data
// to the this buffer.
func (f *FakeFS) WriteFile(filename string, data []byte, perm os.FileMode) (int, error) {
	f.Files[filename] = new(bytes.Buffer)
	return f.Files[filename].Write(data)
}
