package ast

// GENERATED CODE from cmd/tool/gen/ast/ast-gen.go

import "github.com/rhomel/golox/pkg/scanner"

type Expr interface {
	isExpr() // private method to tag which structs are Expr
}

var _ Expr = (*Assign)(nil)

type Assign struct {
	Name  scanner.Token
	Value Expr
}

func (*Assign) isExpr() {}

type AssignStringVisitor interface {
	VisitAssignExprString(*Assign) string
}

func (assign *Assign) AcceptString(visitor AssignStringVisitor) string {
	return visitor.VisitAssignExprString(assign)
}

type AssignVoidVisitor interface {
	VisitAssignExprVoid(*Assign)
}

func (assign *Assign) AcceptVoid(visitor AssignVoidVisitor) {
	visitor.VisitAssignExprVoid(assign)
}

type AssignVisitor interface {
	VisitAssignExpr(*Assign) interface{}
}

func (assign *Assign) Accept(visitor AssignVisitor) interface{} {
	return visitor.VisitAssignExpr(assign)
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

var _ Expr = (*Call)(nil)

type Call struct {
	Callee    Expr
	Paren     scanner.Token
	Arguments []Expr
}

func (*Call) isExpr() {}

type CallStringVisitor interface {
	VisitCallExprString(*Call) string
}

func (call *Call) AcceptString(visitor CallStringVisitor) string {
	return visitor.VisitCallExprString(call)
}

type CallVoidVisitor interface {
	VisitCallExprVoid(*Call)
}

func (call *Call) AcceptVoid(visitor CallVoidVisitor) {
	visitor.VisitCallExprVoid(call)
}

type CallVisitor interface {
	VisitCallExpr(*Call) interface{}
}

func (call *Call) Accept(visitor CallVisitor) interface{} {
	return visitor.VisitCallExpr(call)
}

var _ Expr = (*Get)(nil)

type Get struct {
	Object Expr
	Name   scanner.Token
}

func (*Get) isExpr() {}

type GetStringVisitor interface {
	VisitGetExprString(*Get) string
}

func (get *Get) AcceptString(visitor GetStringVisitor) string {
	return visitor.VisitGetExprString(get)
}

type GetVoidVisitor interface {
	VisitGetExprVoid(*Get)
}

func (get *Get) AcceptVoid(visitor GetVoidVisitor) {
	visitor.VisitGetExprVoid(get)
}

type GetVisitor interface {
	VisitGetExpr(*Get) interface{}
}

func (get *Get) Accept(visitor GetVisitor) interface{} {
	return visitor.VisitGetExpr(get)
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

var _ Expr = (*Logical)(nil)

type Logical struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (*Logical) isExpr() {}

type LogicalStringVisitor interface {
	VisitLogicalExprString(*Logical) string
}

func (logical *Logical) AcceptString(visitor LogicalStringVisitor) string {
	return visitor.VisitLogicalExprString(logical)
}

type LogicalVoidVisitor interface {
	VisitLogicalExprVoid(*Logical)
}

func (logical *Logical) AcceptVoid(visitor LogicalVoidVisitor) {
	visitor.VisitLogicalExprVoid(logical)
}

type LogicalVisitor interface {
	VisitLogicalExpr(*Logical) interface{}
}

func (logical *Logical) Accept(visitor LogicalVisitor) interface{} {
	return visitor.VisitLogicalExpr(logical)
}

var _ Expr = (*Set)(nil)

type Set struct {
	Object Expr
	Name   scanner.Token
	Value  Expr
}

func (*Set) isExpr() {}

type SetStringVisitor interface {
	VisitSetExprString(*Set) string
}

func (set *Set) AcceptString(visitor SetStringVisitor) string {
	return visitor.VisitSetExprString(set)
}

type SetVoidVisitor interface {
	VisitSetExprVoid(*Set)
}

func (set *Set) AcceptVoid(visitor SetVoidVisitor) {
	visitor.VisitSetExprVoid(set)
}

type SetVisitor interface {
	VisitSetExpr(*Set) interface{}
}

func (set *Set) Accept(visitor SetVisitor) interface{} {
	return visitor.VisitSetExpr(set)
}

var _ Expr = (*Super)(nil)

type Super struct {
	Keyword scanner.Token
	Method  scanner.Token
}

func (*Super) isExpr() {}

type SuperStringVisitor interface {
	VisitSuperExprString(*Super) string
}

func (super *Super) AcceptString(visitor SuperStringVisitor) string {
	return visitor.VisitSuperExprString(super)
}

type SuperVoidVisitor interface {
	VisitSuperExprVoid(*Super)
}

func (super *Super) AcceptVoid(visitor SuperVoidVisitor) {
	visitor.VisitSuperExprVoid(super)
}

type SuperVisitor interface {
	VisitSuperExpr(*Super) interface{}
}

func (super *Super) Accept(visitor SuperVisitor) interface{} {
	return visitor.VisitSuperExpr(super)
}

var _ Expr = (*This)(nil)

type This struct {
	Keyword scanner.Token
}

func (*This) isExpr() {}

type ThisStringVisitor interface {
	VisitThisExprString(*This) string
}

func (this *This) AcceptString(visitor ThisStringVisitor) string {
	return visitor.VisitThisExprString(this)
}

type ThisVoidVisitor interface {
	VisitThisExprVoid(*This)
}

func (this *This) AcceptVoid(visitor ThisVoidVisitor) {
	visitor.VisitThisExprVoid(this)
}

type ThisVisitor interface {
	VisitThisExpr(*This) interface{}
}

func (this *This) Accept(visitor ThisVisitor) interface{} {
	return visitor.VisitThisExpr(this)
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

var _ Expr = (*Variable)(nil)

type Variable struct {
	Name scanner.Token
}

func (*Variable) isExpr() {}

type VariableStringVisitor interface {
	VisitVariableExprString(*Variable) string
}

func (variable *Variable) AcceptString(visitor VariableStringVisitor) string {
	return visitor.VisitVariableExprString(variable)
}

type VariableVoidVisitor interface {
	VisitVariableExprVoid(*Variable)
}

func (variable *Variable) AcceptVoid(visitor VariableVoidVisitor) {
	visitor.VisitVariableExprVoid(variable)
}

type VariableVisitor interface {
	VisitVariableExpr(*Variable) interface{}
}

func (variable *Variable) Accept(visitor VariableVisitor) interface{} {
	return visitor.VisitVariableExpr(variable)
}
