package examples

import (
	"bytes"
	"io"
	"unsafe"
)

// Not checked at all.

func withoutArgs()                                {}
func withoutError1() *User                        { return nil }
func withoutError2() (*User, *User)               { return nil, nil }
func withoutError3() (*User, *User, *User)        { return nil, nil, nil }
func withoutError4() (*User, *User, *User, *User) { return nil, nil, nil, nil }

func invalidOrder() (error, *User)               { return nil, nil }
func withError3rd() (*User, bool, error)         { return nil, false, nil }
func withError4th() (*User, *User, *User, error) { return nil, nil, nil, nil }

func slice() ([]int, error) { return nil, nil }

func strNil() (string, error)   { return "nil", nil }
func strEmpty() (string, error) { return "", nil }

// Valid.

func primitivePtrTypeValid() (*int, error) {
	if false {
		return nil, io.EOF
	}
	return new(int), nil
}

func structPtrTypeValid() (*User, error) {
	if false {
		return nil, io.EOF
	}
	return new(User), nil
}

func unsafePtrValid() (unsafe.Pointer, error) {
	if false {
		return nil, io.EOF
	}
	var i int
	return unsafe.Pointer(&i), nil
}

func uintPtrValid() (uintptr, error) {
	if false {
		return 0, io.EOF
	}
	return 0xc82000c290, nil
}

func channelTypeValid() (ChannelType, error) {
	if false {
		return nil, io.EOF
	}
	return make(ChannelType), nil
}

func funcTypeValid() (FuncType, error) {
	if false {
		return nil, io.EOF
	}
	return func(i int) int {
		return 0
	}, nil
}

func ifaceTypeValid() (io.Reader, error) {
	if false {
		return nil, io.EOF
	}
	return new(bytes.Buffer), nil
}

// Unsupported.

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
