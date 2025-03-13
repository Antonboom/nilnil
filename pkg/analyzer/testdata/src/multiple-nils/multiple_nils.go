package examples

import "io"

func ex0() *User    { return new(User) }
func ex00() *User   { return nil }
func ex000() error  { return io.EOF }
func ex0000() error { return nil }

func ex1() (*User, *User, *User, *User) { return nil, nil, nil, nil }
func ex2() (*User, *User, *User, *User) { return nil, new(User), nil, nil }
func ex3() (*User, *User, *User, *User) { return new(User), nil, nil, nil }
func ex4() (*User, *User, *User, *User) { return nil, nil, new(User), nil }

func ex5() (*User, *User, *User, error) {
	return nil, nil, nil, io.EOF
}

func ex6() (*User, *User, *User, error) {
	return nil, nil, new(User), nil // want "return both a `nil` error and an invalid value: use a sentinel error instead"
}

func ex7() (*User, *User, *User, error) {
	return nil, nil, new(User), io.EOF
}

func ex8() (*User, *User, *User, error) {
	return nil, nil, nil, nil // want "return both a `nil` error and an invalid value: use a sentinel error instead"
}

func ex9() (*User, *User, *User, error) {
	return new(User), new(User), new(User), nil
}

func ex10() (*User, *User, *User, error) {
	return new(User), nil, new(User), io.EOF
}

func ex11() (*User, *User, *User, error) {
	return new(User), nil, new(User), nil // want "return both a `nil` error and an invalid value: use a sentinel error instead"
}

func ex12() (error, *User, *User, *User) {
	return nil, new(User), nil, new(User)
}

func ex13() (*User, bool, int, error) {
	return nil, false, 42, nil // want "return both a `nil` error and an invalid value: use a sentinel error instead"
}

type User struct{}
