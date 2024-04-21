package examples

type data struct {
	a int
	b int
}

func nilnil(a, b int) (*data, error) {
	if a == 0 && b == 0 {
		return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
	}
	return &data{a: a, b: b}, nil
}

func main() {
	const a = 0
	const b = 0
	_, _ = nilnil(a, b)
}
