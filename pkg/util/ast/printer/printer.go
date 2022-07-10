package printer

import (
	"fmt"
	"reflect"
	"strings"

	ast "rhomel.com/crafting-interpreters-go/pkg/ast/gen"
	"rhomel.com/crafting-interpreters-go/pkg/util/exit"
)

type AstPrinter struct{}

func (a *AstPrinter) Accept(elem interface{}) string {
	// Go has no dynamic dispatch and inheritance so we have to resort to a type switch
	switch v := elem.(type) {
	case *ast.Binary:
		return v.Accept(a)
	case *ast.Grouping:
		return v.Accept(a)
	case *ast.Literal:
		return v.Accept(a)
	case *ast.Unary:
		return v.Accept(a)
	default:
		exit.Exitf(exit.ExitSyntaxError, "unsupported type: %s", reflect.TypeOf(elem).Name())
		return ""
	}
}

func (a *AstPrinter) VisitBinaryExpr(binary *ast.Binary) string {
	return a.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
}

func (a *AstPrinter) VisitGroupingExpr(grouping *ast.Grouping) string {
	return a.parenthesize("group", grouping.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(literal *ast.Literal) string {
	if isNil(literal) {
		return "nil"
	}
	return fmt.Sprintf("%v", literal.Value)
}

func (a *AstPrinter) VisitUnaryExpr(unary *ast.Unary) string {
	return a.parenthesize(unary.Operator.Lexeme, unary.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...ast.Expr) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(a.Accept(expr))
	}
	builder.WriteString(")")
	return builder.String()
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
