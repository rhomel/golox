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

type BinaryStringVisitor interface {
	VisitBinaryExprString(*Binary) string
}

func (binary *Binary) AcceptString(visitor BinaryStringVisitor) string {
	return visitor.VisitBinaryExprString(binary)
}

var _ Expr = (*Grouping)(nil)

type Grouping struct {
	Expression Expr
}

func (*Grouping) isExpr() {}

type GroupingStringVisitor interface {
	VisitGroupingExprString(*Grouping) string
}

func (grouping *Grouping) AcceptString(visitor GroupingStringVisitor) string {
	return visitor.VisitGroupingExprString(grouping)
}

var _ Expr = (*Literal)(nil)

type Literal struct {
	Value interface{}
}

func (*Literal) isExpr() {}

type LiteralStringVisitor interface {
	VisitLiteralExprString(*Literal) string
}

func (literal *Literal) AcceptString(visitor LiteralStringVisitor) string {
	return visitor.VisitLiteralExprString(literal)
}

var _ Expr = (*Unary)(nil)

type Unary struct {
	Operator scanner.Token
	Right    Expr
}

func (*Unary) isExpr() {}

type UnaryStringVisitor interface {
	VisitUnaryExprString(*Unary) string
}

func (unary *Unary) AcceptString(visitor UnaryStringVisitor) string {
	return visitor.VisitUnaryExprString(unary)
}
