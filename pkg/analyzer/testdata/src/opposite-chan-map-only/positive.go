package examples

import "io"

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
