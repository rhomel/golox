package main

// 5.4 A (Not Very) Pretty Printer
// ref:
//   https://craftinginterpreters.com/representing-code.html#a-not-very-pretty-printer

import (
	"fmt"

	ast "github.com/rhomel/golox/pkg/ast/gen"
	"github.com/rhomel/golox/pkg/scanner"
	"github.com/rhomel/golox/pkg/util/ast/printer"
)

func main() {
	expression := &ast.Binary{
		Left: &ast.Unary{
			Operator: scanner.Token{scanner.MINUS, "-", nil, 1},
			Right:    &ast.Literal{123},
		},
		Operator: scanner.Token{scanner.STAR, "*", nil, 1},
		Right: &ast.Grouping{
			Expression: &ast.Literal{45.67},
		},
	}

	printer := &printer.AstPrinter{}
	fmt.Println(printer.Accept(expression))
}
