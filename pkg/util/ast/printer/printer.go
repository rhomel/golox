package printer

import (
	"fmt"
	"reflect"
	"strings"

	ast "github.com/rhomel/golox/pkg/ast/gen"
	"github.com/rhomel/golox/pkg/util/check"
	"github.com/rhomel/golox/pkg/util/exit"
)

type AstPrinter struct{}

func (a *AstPrinter) Accept(elem interface{}) string {
	// Go has no dynamic dispatch and inheritance so we have to resort to a type switch
	switch v := elem.(type) {
	case *ast.Binary:
		return v.AcceptString(a)
	case *ast.Grouping:
		return v.AcceptString(a)
	case *ast.Literal:
		return v.AcceptString(a)
	case *ast.Unary:
		return v.AcceptString(a)
	default:
		exit.Exitf(exit.ExitSyntaxError, "unsupported type: %s", reflect.TypeOf(elem).Name())
		return ""
	}
}

func (a *AstPrinter) VisitBinaryExprString(binary *ast.Binary) string {
	return a.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
}

func (a *AstPrinter) VisitGroupingExprString(grouping *ast.Grouping) string {
	return a.parenthesize("group", grouping.Expression)
}

func (a *AstPrinter) VisitLiteralExprString(literal *ast.Literal) string {
	if check.IsNil(literal) {
		return "nil"
	}
	return fmt.Sprintf("%v", literal.Value)
}

func (a *AstPrinter) VisitUnaryExprString(unary *ast.Unary) string {
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
