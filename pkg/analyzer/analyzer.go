package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const reportMsg = "return both the `nil` error and invalid value: use a sentinel error instead"

// New returns new nilnil analyzer.
func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "nilnil",
		Doc:      "Checks that there is no simultaneous return of `nil` error and an invalid value.",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

var (
	types = []ast.Node{(*ast.TypeSpec)(nil)}

	funcAndReturns = []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
		(*ast.ReturnStmt)(nil),
	}
)

type typeSpecByName map[string]*ast.TypeSpec

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	typeSpecs := typeSpecByName{}
	insp.Preorder(types, func(n ast.Node) {
		t := n.(*ast.TypeSpec)
		typeSpecs[t.Name.Name] = t
	})

	var fs funcTypeStack
	insp.Nodes(funcAndReturns, func(n ast.Node, push bool) (proceed bool) {
		switch v := n.(type) {
		case *ast.FuncLit:
			if push {
				fs.Push(v.Type)
			} else {
				fs.Pop()
			}

		case *ast.FuncDecl:
			if push {
				fs.Push(v.Type)
			} else {
				fs.Pop()
			}

		case *ast.ReturnStmt:
			ft := fs.Top() // Current function.

			if !push || len(v.Results) != 2 || ft == nil || ft.Results == nil || len(ft.Results.List) != 2 {
				return false
			}

			fRes1, fRes2 := ft.Results.List[0], ft.Results.List[1]
			if !(isDangerNilField(fRes1, typeSpecs) && isErrorField(fRes2)) {
				return
			}

			rRes1, rRes2 := v.Results[0], v.Results[1]
			if isNil(rRes1) && isNil(rRes2) {
				pass.Reportf(v.Pos(), reportMsg)
			}
		}

		return true
	})

	return nil, nil
}

func isDangerNilField(f *ast.Field, types typeSpecByName) bool {
	return isDangerNilType(f.Type, types)
}

func isDangerNilType(t ast.Expr, types typeSpecByName) bool {
	switch v := t.(type) {
	case *ast.StarExpr, *ast.ChanType, *ast.FuncType, *ast.InterfaceType, *ast.MapType:
		return true
	case *ast.Ident:
		if t, ok := types[v.Name]; ok {
			return isDangerNilType(t.Type, nil)
		}
	}
	return false
}

func isErrorField(f *ast.Field) bool {
	return isIdent(f.Type, "error")
}

func isNil(e ast.Expr) bool {
	return isIdent(e, "nil")
}

func isIdent(n ast.Node, name string) bool {
	i, ok := n.(*ast.Ident)
	if !ok {
		return false
	}
	return i.Name == name
}
