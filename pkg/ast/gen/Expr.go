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

type BinaryVoidVisitor interface {
	VisitBinaryExprVoid(*Binary)
}

func (binary *Binary) AcceptVoid(visitor BinaryVoidVisitor) {
	visitor.VisitBinaryExprVoid(binary)
}

type BinaryVisitor interface {
	VisitBinaryExpr(*Binary) interface{}
}

func (binary *Binary) Accept(visitor BinaryVisitor) interface{} {
	return visitor.VisitBinaryExpr(binary)
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

type GroupingVoidVisitor interface {
	VisitGroupingExprVoid(*Grouping)
}

func (grouping *Grouping) AcceptVoid(visitor GroupingVoidVisitor) {
	visitor.VisitGroupingExprVoid(grouping)
}

type GroupingVisitor interface {
	VisitGroupingExpr(*Grouping) interface{}
}

func (grouping *Grouping) Accept(visitor GroupingVisitor) interface{} {
	return visitor.VisitGroupingExpr(grouping)
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

type LiteralVoidVisitor interface {
	VisitLiteralExprVoid(*Literal)
}

func (literal *Literal) AcceptVoid(visitor LiteralVoidVisitor) {
	visitor.VisitLiteralExprVoid(literal)
}

type LiteralVisitor interface {
	VisitLiteralExpr(*Literal) interface{}
}

func (literal *Literal) Accept(visitor LiteralVisitor) interface{} {
	return visitor.VisitLiteralExpr(literal)
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

type UnaryVoidVisitor interface {
	VisitUnaryExprVoid(*Unary)
}

func (unary *Unary) AcceptVoid(visitor UnaryVoidVisitor) {
	visitor.VisitUnaryExprVoid(unary)
}

type UnaryVisitor interface {
	VisitUnaryExpr(*Unary) interface{}
}

func (unary *Unary) Accept(visitor UnaryVisitor) interface{} {
	return visitor.VisitUnaryExpr(unary)
}
