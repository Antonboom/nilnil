package analyzer

import (
	"go/ast"
)

type FuncTypeStack []*ast.FuncType

func (s *FuncTypeStack) Push(f *ast.FuncType) {
	*s = append(*s, f)
}

func (s *FuncTypeStack) Pop() *ast.FuncType {
	if len(*s) == 0 {
		return nil
	}

	last := len(*s) - 1
	f := (*s)[last]
	*s = (*s)[:last]
	return f
}

func (s *FuncTypeStack) Top() *ast.FuncType {
	if len(*s) == 0 {
		return nil
	}
	return (*s)[len(*s)-1]
}
