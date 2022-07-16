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

var _ Stmt = (*Class)(nil)

type Class struct {
	Name       scanner.Token
	Superclass *Variable
	Methods    []*Function
}

func (*Class) isStmt() {}

type ClassStringVisitor interface {
	VisitClassStmtString(*Class) string
}

func (class *Class) AcceptString(visitor ClassStringVisitor) string {
	return visitor.VisitClassStmtString(class)
}

type ClassVoidVisitor interface {
	VisitClassStmtVoid(*Class)
}

func (class *Class) AcceptVoid(visitor ClassVoidVisitor) {
	visitor.VisitClassStmtVoid(class)
}

type ClassVisitor interface {
	VisitClassStmt(*Class) interface{}
}

func (class *Class) Accept(visitor ClassVisitor) interface{} {
	return visitor.VisitClassStmt(class)
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

var _ Stmt = (*Function)(nil)

type Function struct {
	Name   scanner.Token
	Params []scanner.Token
	Body   []Stmt
}

func (*Function) isStmt() {}

type FunctionStringVisitor interface {
	VisitFunctionStmtString(*Function) string
}

func (function *Function) AcceptString(visitor FunctionStringVisitor) string {
	return visitor.VisitFunctionStmtString(function)
}

type FunctionVoidVisitor interface {
	VisitFunctionStmtVoid(*Function)
}

func (function *Function) AcceptVoid(visitor FunctionVoidVisitor) {
	visitor.VisitFunctionStmtVoid(function)
}

type FunctionVisitor interface {
	VisitFunctionStmt(*Function) interface{}
}

func (function *Function) Accept(visitor FunctionVisitor) interface{} {
	return visitor.VisitFunctionStmt(function)
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

var _ Stmt = (*ReturnStmt)(nil)

type ReturnStmt struct {
	Keyword scanner.Token
	Value   Expr
}

func (*ReturnStmt) isStmt() {}

type ReturnStmtStringVisitor interface {
	VisitReturnStmtStmtString(*ReturnStmt) string
}

func (returnstmt *ReturnStmt) AcceptString(visitor ReturnStmtStringVisitor) string {
	return visitor.VisitReturnStmtStmtString(returnstmt)
}

type ReturnStmtVoidVisitor interface {
	VisitReturnStmtStmtVoid(*ReturnStmt)
}

func (returnstmt *ReturnStmt) AcceptVoid(visitor ReturnStmtVoidVisitor) {
	visitor.VisitReturnStmtStmtVoid(returnstmt)
}

type ReturnStmtVisitor interface {
	VisitReturnStmtStmt(*ReturnStmt) interface{}
}

func (returnstmt *ReturnStmt) Accept(visitor ReturnStmtVisitor) interface{} {
	return visitor.VisitReturnStmtStmt(returnstmt)
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

var _ Stmt = (*While)(nil)

type While struct {
	Condition Expr
	Body      Stmt
}

func (*While) isStmt() {}

type WhileStringVisitor interface {
	VisitWhileStmtString(*While) string
}

func (while *While) AcceptString(visitor WhileStringVisitor) string {
	return visitor.VisitWhileStmtString(while)
}

type WhileVoidVisitor interface {
	VisitWhileStmtVoid(*While)
}

func (while *While) AcceptVoid(visitor WhileVoidVisitor) {
	visitor.VisitWhileStmtVoid(while)
}

type WhileVisitor interface {
	VisitWhileStmt(*While) interface{}
}

func (while *While) Accept(visitor WhileVisitor) interface{} {
	return visitor.VisitWhileStmt(while)
}
