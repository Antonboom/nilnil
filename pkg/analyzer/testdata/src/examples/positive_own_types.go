package examples

type (
	StructPtrType    *User
	PrimitivePtrType *int
	ChannelType      chan int
	FuncType         func(int) int
	Checker          interface{ Check() }
)

func structPtrType() (StructPtrType, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func primitivePtrType() (PrimitivePtrType, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func channelType() (ChannelType, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func funcType() (FuncType, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

func ifaceType() (Checker, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

type checkerAlias = Checker

func ifaceTypeAliased() (checkerAlias, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}

type (
	IntegerType    int
	PtrIntegerType *IntegerType
)

func ptrIntegerType() (PtrIntegerType, error) {
	return nil, nil // want "return both the `nil` error and invalid value: use a sentinel error instead"
}
