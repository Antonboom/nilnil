package examples

type User struct{}

func primitivePtr() (*int, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func structPtr() (*User, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func uintPtr0o() (uintptr, error) {
	return 0o000, nil
}

func chBi() (chan int, error) {
	return nil, nil
}

func fun() (func(), error) {
	return nil, nil
}

func anyType() (any, error) {
	return nil, nil
}

func m1() (map[int]int, error) {
	return nil, nil
}
