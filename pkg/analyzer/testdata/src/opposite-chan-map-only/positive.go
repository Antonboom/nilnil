package examples

import (
	"io"
	"unsafe"
)

type User struct{}

func primitivePtr() (*int, error) {
	return new(int), io.EOF
}

func structPtr() (*User, error) {
	return new(User), io.EOF
}

func uintPtr0o() (uintptr, error) {
	return 0o000, io.EOF
}

func unsafePtr() (unsafe.Pointer, error) {
	return nil, nil
}

func chBi() (chan int, error) {
	return make(chan int), io.EOF // want "return both a non-nil error and a valid value: use separate returns instead"
}

func fun() (func(), error) {
	return func() {}, io.EOF
}

func anyType() (any, error) {
	return io.EOF, io.EOF
}

func m1() (map[int]int, error) {
	return make(map[int]int), io.EOF // want "return both a non-nil error and a valid value: use separate returns instead"
}

func structPtrInvalid() (*int, error) {
	return nil, nil
}

func structPtrValid() (*int, error) {
	return new(int), nil
}

func chInvalid() (chan int, error) {
	return nil, nil // want "return both a `nil` error and an invalid value: use a sentinel error instead"
}
