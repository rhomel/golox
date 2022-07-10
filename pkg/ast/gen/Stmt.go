package ast

// GENERATED CODE from cmd/tool/gen/ast/ast-gen.go

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

type PrintVisitor interface {
	VisitPrintStmt(*Print) interface{}
}

func (print *Print) Accept(visitor PrintVisitor) interface{} {
	return visitor.VisitPrintStmt(print)
}
