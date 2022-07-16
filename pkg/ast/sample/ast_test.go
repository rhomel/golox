package ast_test

import "github.com/rhomel/golox/pkg/scanner"

// Manually created sample of code we want to generate with cmd/tool/gen/ast

type Expr interface {
	isExpr() // private method to tag which structs are Expr
}

var _ Expr = (*Binary)(nil)

type Binary struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (*Binary) isExpr() {}

type BinaryVisitor interface {
	VisitBinaryExpr(*Binary)
}

func (binary *Binary) Accept(visitor BinaryVisitor) {
	visitor.VisitBinaryExpr(binary)
}
