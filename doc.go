// Package sys is a small utility library for interacting with the system from your programs. It exists to make your system level code more unit testable. For example, instead of reading a file with os.Open or ioutil.ReadAll, you can call fs.ReadAll, where fs is a an FS implementation. Then, in unit tests, you can substitute a FakeFS to test your code without any disk dependencies.
//
// Example usage:
//
//  func readFile(fs FS, fName string) (string, error) {
//    return fs.ReadFile(fName)
//  }
//
//  //  in tests
//  func TestReadFile(t *testing.T) {
//    fakeFS := NewFakeFS()
//    contents, err := readFile(fs, "someFile")
//    // contents should be "" and there should be an error, since the fake filesystem had no files in it
//  }
package sys
