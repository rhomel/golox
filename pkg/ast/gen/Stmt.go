package ast

// GENERATED CODE from cmd/tool/gen/ast/ast-gen.go

import "rhomel.com/crafting-interpreters-go/pkg/scanner"

type Stmt interface {
	isStmt() // private method to tag which structs are Stmt
}

var _ Stmt = (*Expression)(nil)

type Expression struct {
	Expression Expr
}

func (*Expression) isStmt() {}

type ExpressionStringVisitor interface {
	VisitExpressionStmtString(*Expression) string
}

func (expression *Expression) AcceptString(visitor ExpressionStringVisitor) string {
	return visitor.VisitExpressionStmtString(expression)
}

type ExpressionVoidVisitor interface {
	VisitExpressionStmtVoid(*Expression)
}

func (expression *Expression) AcceptVoid(visitor ExpressionVoidVisitor) {
	visitor.VisitExpressionStmtVoid(expression)
}

type ExpressionVisitor interface {
	VisitExpressionStmt(*Expression) interface{}
}

func (expression *Expression) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitExpressionStmt(expression)
}

var _ Stmt = (*Print)(nil)

type Print struct {
	Expression Expr
}

func (*Print) isStmt() {}

type PrintStringVisitor interface {
	VisitPrintStmtString(*Print) string
}

func (print *Print) AcceptString(visitor PrintStringVisitor) string {
	return visitor.VisitPrintStmtString(print)
}

type PrintVoidVisitor interface {
	VisitPrintStmtVoid(*Print)
}

func (print *Print) AcceptVoid(visitor PrintVoidVisitor) {
	visitor.VisitPrintStmtVoid(print)
}

type PrintVisitor interface {
	VisitPrintStmt(*Print) interface{}
}

func (print *Print) Accept(visitor PrintVisitor) interface{} {
	return visitor.VisitPrintStmt(print)
}

var _ Stmt = (*VarStmt)(nil)

type VarStmt struct {
	Name        scanner.Token
	Initializer Expr
}

func (*VarStmt) isStmt() {}

type VarStmtStringVisitor interface {
	VisitVarStmtStmtString(*VarStmt) string
}

func (varstmt *VarStmt) AcceptString(visitor VarStmtStringVisitor) string {
	return visitor.VisitVarStmtStmtString(varstmt)
}

type VarStmtVoidVisitor interface {
	VisitVarStmtStmtVoid(*VarStmt)
}

func (varstmt *VarStmt) AcceptVoid(visitor VarStmtVoidVisitor) {
	visitor.VisitVarStmtStmtVoid(varstmt)
}

type VarStmtVisitor interface {
	VisitVarStmtStmt(*VarStmt) interface{}
}

func (varstmt *VarStmt) Accept(visitor VarStmtVisitor) interface{} {
	return visitor.VisitVarStmtStmt(varstmt)
}
