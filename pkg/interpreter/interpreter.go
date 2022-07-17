package interpreter

import ast "github.com/rhomel/golox/pkg/ast/gen"

type Interpreter interface {
	Interpret(statements []ast.Stmt)
	Resolve(expr ast.Expr, depth int)
}
