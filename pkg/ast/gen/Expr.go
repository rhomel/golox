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

type BinaryVisitor interface {
	VisitBinaryExpr(*Binary) string
}

func (binary *Binary) Accept(visitor BinaryVisitor) string {
	return visitor.VisitBinaryExpr(binary)
}

var _ Expr = (*Grouping)(nil)

type Grouping struct {
	Expression Expr
}

func (*Grouping) isExpr() {}

type GroupingVisitor interface {
	VisitGroupingExpr(*Grouping) string
}

func (grouping *Grouping) Accept(visitor GroupingVisitor) string {
	return visitor.VisitGroupingExpr(grouping)
}

var _ Expr = (*Literal)(nil)

type Literal struct {
	Value interface{}
}

func (*Literal) isExpr() {}

type LiteralVisitor interface {
	VisitLiteralExpr(*Literal) string
}

func (literal *Literal) Accept(visitor LiteralVisitor) string {
	return visitor.VisitLiteralExpr(literal)
}

var _ Expr = (*Unary)(nil)

type Unary struct {
	Operator scanner.Token
	Right    Expr
}

func (*Unary) isExpr() {}

type UnaryVisitor interface {
	VisitUnaryExpr(*Unary) string
}

func (unary *Unary) Accept(visitor UnaryVisitor) string {
	return visitor.VisitUnaryExpr(unary)
}
