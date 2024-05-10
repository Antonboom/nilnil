package unsafe

import (
	"io"
	"unsafe"
)

func unsafePtr() (unsafe.Pointer, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func uintPtr() (uintptr, error) {
	return 0, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func uintPtr0b() (uintptr, error) {
	return 0b0, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func uintPtr0x() (uintptr, error) {
	return 0x00, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func uintPtr0o() (uintptr, error) {
	return 0o000, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
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
