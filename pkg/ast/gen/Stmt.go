package ast

// GENERATED CODE from cmd/tool/gen/ast/ast-gen.go

import "rhomel.com/crafting-interpreters-go/pkg/scanner"

type Stmt interface {
	isStmt() // private method to tag which structs are Stmt
}

var _ Stmt = (*Block)(nil)

type Block struct {
	Statements []Stmt
}

func (*Block) isStmt() {}

type BlockStringVisitor interface {
	VisitBlockStmtString(*Block) string
}

func (block *Block) AcceptString(visitor BlockStringVisitor) string {
	return visitor.VisitBlockStmtString(block)
}

type BlockVoidVisitor interface {
	VisitBlockStmtVoid(*Block)
}

func (block *Block) AcceptVoid(visitor BlockVoidVisitor) {
	visitor.VisitBlockStmtVoid(block)
}

type BlockVisitor interface {
	VisitBlockStmt(*Block) interface{}
}

func (block *Block) Accept(visitor BlockVisitor) interface{} {
	return visitor.VisitBlockStmt(block)
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

var _ Stmt = (*IfStmt)(nil)

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (*IfStmt) isStmt() {}

type IfStmtStringVisitor interface {
	VisitIfStmtStmtString(*IfStmt) string
}

func (ifstmt *IfStmt) AcceptString(visitor IfStmtStringVisitor) string {
	return visitor.VisitIfStmtStmtString(ifstmt)
}

type IfStmtVoidVisitor interface {
	VisitIfStmtStmtVoid(*IfStmt)
}

func (ifstmt *IfStmt) AcceptVoid(visitor IfStmtVoidVisitor) {
	visitor.VisitIfStmtStmtVoid(ifstmt)
}

type IfStmtVisitor interface {
	VisitIfStmtStmt(*IfStmt) interface{}
}

func (ifstmt *IfStmt) Accept(visitor IfStmtVisitor) interface{} {
	return visitor.VisitIfStmtStmt(ifstmt)
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
