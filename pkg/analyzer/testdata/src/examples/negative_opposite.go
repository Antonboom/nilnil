package examples

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"unsafe"
)

func primitivePtrTypeOpposite() (*int, error) {
	if false {
		return nil, io.EOF
	}
	return new(int), errors.New("validation failed")
}

func structPtrTypeOpposite() (*User, error) {
	if false {
		return nil, io.EOF
	}
	return new(User), fmt.Errorf("invalid %v", 42)
}

func unsafePtrOpposite() (unsafe.Pointer, error) {
	if false {
		return nil, io.EOF
	}
	var i int
	return unsafe.Pointer(&i), io.EOF
}

func uintPtrOpposite() (uintptr, error) {
	if false {
		return 0, io.EOF
	}
	return 0x1000, wrap(io.EOF)
}

func channelTypeOpposite() (ChannelType, error) {
	if false {
		return nil, io.EOF
	}
	return make(ChannelType), fmt.Errorf("wrapped: %w", io.EOF)
}

func funcTypeOpposite() (FuncType, error) {
	if false {
		return nil, io.EOF
	}
	return func(i int) int {
		return 0
	}, errors.New("no func type, please")
}

func ifaceTypeOpposite() (io.Reader, error) {
	if false {
		return nil, io.EOF
	}
	return new(bytes.Buffer), new(net.AddrError)
}
