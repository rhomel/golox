package ast

// GENERATED CODE from cmd/tool/gen/ast/ast-gen.go

import "rhomel.com/crafting-interpreters-go/pkg/scanner"

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

var _ Expr = (*Grouping)(nil)

type Grouping struct {
	Expression Expr
}

func (*Grouping) isExpr() {}

var _ Expr = (*Literal)(nil)

type Literal struct {
	Value interface{}
}

func (*Literal) isExpr() {}

var _ Expr = (*Unary)(nil)

type Unary struct {
	Operator scanner.Token
	Right    Expr
}

func (*Unary) isExpr() {}
