package analyzer

import (
	"go/ast"
	"testing"
)

func Test_funcTypeStack(t *testing.T) {
	var fs funcTypeStack

	if fs.Pop() != nil {
		t.FailNow()
	}
	if fs.Top() != nil {
		t.FailNow()
	}

	funcs := []string{
		"func1",
		"func2",
		"func3",
		"func4",
		"func5",
	}

	for _, f := range funcs {
		fs.Push(&ast.FuncType{Params: &ast.FieldList{
			List: []*ast.Field{{Names: []*ast.Ident{{Name: f}}}},
		}})
		assertFuncFirstParamName(t, fs.Top(), f)
	}

	if len(fs) != len(funcs) {
		t.FailNow()
	}

	for i := len(funcs) - 1; i > 0; i-- {
		assertFuncFirstParamName(t, fs.Top(), funcs[i])
		assertFuncFirstParamName(t, fs.Pop(), funcs[i])
	}
}

func assertFuncFirstParamName(t *testing.T, ft *ast.FuncType, name string) {
	t.Helper()

	if ft == nil {
		t.FailNow()
	}

	if ft.Params.List[0].Names[0].Name != name {
		t.FailNow()
	}
}
