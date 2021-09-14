package examples

import (
	"io"
	"unsafe"
)

func withoutArgs()                                {}
func withoutError1() *User                        { return nil }
func withoutError2() (*User, *User)               { return nil, nil }
func withoutError3() (*User, *User, *User)        { return nil, nil, nil }
func withoutError4() (*User, *User, *User, *User) { return nil, nil, nil, nil }

// Unsupported.

func invalidOrder() (error, *User)               { return nil, nil }
func withError3rd() (*User, bool, error)         { return nil, false, nil }
func withError4th() (*User, *User, *User, error) { return nil, nil, nil, nil }
func unsafePtr() (unsafe.Pointer, error)         { return nil, nil }
func uintPtr() (uintptr, error)                  { return 0, nil }
func slice() ([]int, error)                      { return nil, nil }
func ifaceExtPkg() (io.Closer, error)            { return nil, nil }

func implicitNil1() (*User, error) {
	err := (error)(nil)
	return nil, err
}

func implicitNil2() (*User, error) {
	err := io.EOF
	err = nil
	return nil, err
}

func implicitNil3() (*User, error) {
	return nil, wrap(nil)
}
func wrap(err error) error { return err }
