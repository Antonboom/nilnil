package examples

import "unsafe"

type User struct{}

func primitivePtr() (*int, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func structPtr() (*User, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func emptyStructPtr() (*struct{}, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func anonymousStructPtr() (*struct{ ID string }, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

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

func chBi() (chan int, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func chIn() (chan<- int, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func chOut() (<-chan int, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func fun() (func(), error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func funWithArgsAndResults() (func(a, b, c int) (int, int), error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func iface() (interface{}, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func anyType() (any, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func m1() (map[int]int, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func m2() (map[int]*User, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

type mapAlias = map[int]*User

func m3() (mapAlias, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

type Storage struct{}

func (s *Storage) GetUser() (*User, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func ifReturn() (*User, error) {
	var s Storage
	if _, err := s.GetUser(); err != nil {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}
	return new(User), nil
}

func forReturn() (*User, error) {
	for {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}
}

func multipleReturn() (*User, error) {
	var s Storage

	if _, err := s.GetUser(); err != nil {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}

	if _, err := s.GetUser(); err != nil {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}

	if _, err := s.GetUser(); err != nil {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}

	return new(User), nil
}

func nested() {
	_ = func() (*User, error) {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}

	_, _ = func() (*User, error) {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}()
}

func deeplyNested() {
	_ = func() {
		_ = func() int {
			_ = func() {
				_ = func() (*User, error) {
					_ = func() {}
					_ = func() int { return 0 }
					return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
				}
			}
			return 0
		}
	}
}

type MyError interface {
	error
	Code() string
}

func myError() (*User, MyError) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}
