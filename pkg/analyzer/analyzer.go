package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "nilnil"
	doc  = "Checks that there is no simultaneous return of `nil` error and an invalid value."

	reportMsg = "return both the `nil` error and invalid value: use a sentinel error instead"
)

// New returns new nilnil analyzer.
func New() *analysis.Analyzer {
	n := newNilNil()

	a := &analysis.Analyzer{
		Name:     name,
		Doc:      doc,
		Run:      n.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	a.Flags.Var(&n.checkedTypes, "checked-types", "coma separated list")

	return a
}

type nilNil struct {
	checkedTypes checkedTypes
}

func newNilNil() *nilNil {
	return &nilNil{
		checkedTypes: newDefaultCheckedTypes(),
	}
}

var funcAndReturns = []ast.Node{
	(*ast.FuncDecl)(nil),
	(*ast.FuncLit)(nil),
	(*ast.ReturnStmt)(nil),
}

func (n *nilNil) run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	var fs funcTypeStack
	insp.Nodes(funcAndReturns, func(node ast.Node, push bool) (proceed bool) {
		switch v := node.(type) {
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
			if !(n.isDangerNilField(pass, fRes1) && n.isErrorField(pass, fRes2)) {
				return false
			}

			rRes1, rRes2 := v.Results[0], v.Results[1]
			if isNil(pass, rRes1) && isNil(pass, rRes2) {
				pass.Reportf(v.Pos(), reportMsg)
			}
		}

		return true
	})

	return nil, nil //nolint:nilnil
}

func (n *nilNil) isDangerNilField(pass *analysis.Pass, f *ast.Field) bool {
	return n.isDangerNilType(pass.TypesInfo.TypeOf(f.Type))
}

func (n *nilNil) isDangerNilType(t types.Type) bool {
	switch v := t.(type) {
	case *types.Pointer:
		return n.checkedTypes.Contains(ptrType)

	case *types.Signature:
		return n.checkedTypes.Contains(funcType)

	case *types.Interface:
		return n.checkedTypes.Contains(ifaceType)

	case *types.Map:
		return n.checkedTypes.Contains(mapType)

	case *types.Chan:
		return n.checkedTypes.Contains(chanType)

	case *types.Named:
		return n.isDangerNilType(v.Underlying())
	}
	return false
}

var errorIface = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func (n *nilNil) isErrorField(pass *analysis.Pass, f *ast.Field) bool {
	t := pass.TypesInfo.TypeOf(f.Type)
	if t == nil {
		return false
	}

	_, ok := t.Underlying().(*types.Interface)
	return ok && types.Implements(t, errorIface)
}

func isNil(pass *analysis.Pass, e ast.Expr) bool {
	i, ok := e.(*ast.Ident)
	if !ok {
		return false
	}

	_, ok = pass.TypesInfo.ObjectOf(i).(*types.Nil)
	return ok
}
