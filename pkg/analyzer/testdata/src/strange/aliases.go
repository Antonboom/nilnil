package strange

type error = *int

func myOwnError() (*int, error) {
	return nil, nil
}

func myOwnNil() (*int, error) {
	nil := new(int)
	return nil, nil
}
