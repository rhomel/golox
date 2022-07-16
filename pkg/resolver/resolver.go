package resolver

import (
	ast "rhomel.com/crafting-interpreters-go/pkg/ast/gen"
	"rhomel.com/crafting-interpreters-go/pkg/interpreter"
	"rhomel.com/crafting-interpreters-go/pkg/scanner"
	"rhomel.com/crafting-interpreters-go/pkg/util/check"
	"rhomel.com/crafting-interpreters-go/pkg/util/exit"
)

type ErrorReporter interface {
	ResolveError(token scanner.Token, message string)
}

type FunctionType int

const (
	NONE        FunctionType = 0
	FUNCTION    FunctionType = 1
	INITIALIZER FunctionType = 2
	METHOD      FunctionType = 3
)

type ClassType int

const (
	NOCLASS  ClassType = 0
	CLASS    ClassType = 1
	SUBCLASS ClassType = 2
)

type Resolver struct {
	in              *interpreter.Interpreter
	reporter        ErrorReporter
	scopes          *stack
	currentFunction FunctionType
	curentClass     ClassType
}

func NewResolver(in *interpreter.Interpreter, reporter ErrorReporter) *Resolver {
	return &Resolver{
		in:              in,
		reporter:        reporter,
		scopes:          &stack{},
		currentFunction: NONE,
		curentClass:     NOCLASS,
	}
}

func (re *Resolver) ResolveStmts(statements []ast.Stmt) {
	for _, stmt := range statements {
		re.resolve(stmt)
	}
}

func (re *Resolver) resolve(elem interface{}) {
	switch v := elem.(type) {
	case *ast.Binary:
		v.AcceptVoid(re)
	case *ast.Call:
		v.AcceptVoid(re)
	case *ast.Get:
		v.AcceptVoid(re)
	case *ast.Grouping:
		v.AcceptVoid(re)
	case *ast.Literal:
		v.AcceptVoid(re)
	case *ast.Logical:
		v.AcceptVoid(re)
	case *ast.Set:
		v.AcceptVoid(re)
	case *ast.Super:
		v.AcceptVoid(re)
	case *ast.This:
		v.AcceptVoid(re)
	case *ast.Unary:
		v.AcceptVoid(re)
	case *ast.Variable:
		v.AcceptVoid(re)
	case *ast.Assign:
		v.AcceptVoid(re)
	case *ast.Class:
		v.AcceptVoid(re)
	case *ast.IfStmt:
		v.AcceptVoid(re)
	case *ast.Block:
		v.AcceptVoid(re)
	case *ast.Expression:
		v.AcceptVoid(re)
	case *ast.Function:
		v.AcceptVoid(re)
	case *ast.Print:
		v.AcceptVoid(re)
	case *ast.ReturnStmt:
		v.AcceptVoid(re)
	case *ast.VarStmt:
		v.AcceptVoid(re)
	case *ast.While:
		v.AcceptVoid(re)
	default:
		exit.Exitf(exit.ExitSyntaxError, "unsupported expression/statement: %s", check.TypeOf(elem))
	}
}

func (re *Resolver) resolveLocal(expr ast.Expr, name scanner.Token) {
	for i := re.scopes.size() - 1; i >= 0; i-- {
		_, containsKey := re.scopes.get(i)[name.Lexeme]
		if containsKey {
			re.in.Resolve(expr, re.scopes.size()-1-i)
			return
		}
	}
}

func (re *Resolver) resolveFunction(function *ast.Function, typ FunctionType) {
	enclosingFunction := re.currentFunction
	re.currentFunction = typ
	re.beginScope()
	for _, param := range function.Params {
		re.declare(param)
		re.define(param)
	}
	re.ResolveStmts(function.Body)
	re.endScope()
	re.currentFunction = enclosingFunction
}

func (re *Resolver) beginScope() {
	re.scopes.push(make(map[string]bool))
}

func (re *Resolver) endScope() {
	re.scopes.pop()
}

func (re *Resolver) declare(name scanner.Token) {
	if re.scopes.isEmpty() {
		return
	}
	scope := re.scopes.peek()
	if _, ok := scope[name.Lexeme]; ok {
		re.reporter.ResolveError(name, "Already a variable with this name in this scope.")
	}
	scope[name.Lexeme] = false
}

func (re *Resolver) define(name scanner.Token) {
	if re.scopes.isEmpty() {
		return
	}
	scope := re.scopes.peek()
	scope[name.Lexeme] = true
}

func (re *Resolver) VisitBlockStmtVoid(block *ast.Block) {
	re.beginScope()
	re.ResolveStmts(block.Statements)
	re.endScope()
}

func (re *Resolver) VisitExpressionStmtVoid(stmt *ast.Expression) {
	re.resolve(stmt.Expression)
}

func (re *Resolver) VisitFunctionStmtVoid(stmt *ast.Function) {
	re.declare(stmt.Name)
	re.define(stmt.Name)
	re.resolveFunction(stmt, FUNCTION)
}

func (re *Resolver) VisitClassStmtVoid(class *ast.Class) {
	enclosingClass := re.curentClass
	re.curentClass = CLASS
	re.declare(class.Name)
	re.define(class.Name)
	if class.Superclass != nil && class.Name.Lexeme == class.Superclass.Name.Lexeme {
		re.reporter.ResolveError(class.Superclass.Name, "A class can't inherit from itself.")
	}
	if class.Superclass != nil {
		re.curentClass = SUBCLASS
		re.beginScope()
		re.scopes.peek()["super"] = true
		re.resolve(class.Superclass)
	}
	re.beginScope()
	re.scopes.peek()["this"] = true
	for _, method := range class.Methods {
		declaration := METHOD
		if method.Name.Lexeme == "init" {
			declaration = INITIALIZER
		}
		re.resolveFunction(method, declaration)
	}
	re.endScope()
	if class.Superclass != nil {
		re.endScope()
	}
	re.curentClass = enclosingClass
}

func (re *Resolver) VisitIfStmtStmtVoid(stmt *ast.IfStmt) {
	re.resolve(stmt.Condition)
	re.resolve(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		re.resolve(stmt.ElseBranch)
	}
}

func (re *Resolver) VisitPrintStmtVoid(stmt *ast.Print) {
	re.resolve(stmt.Expression)
}

func (re *Resolver) VisitReturnStmtStmtVoid(stmt *ast.ReturnStmt) {
	if re.currentFunction == NONE {
		re.reporter.ResolveError(stmt.Keyword, "Can't return from top-level code.")
	}
	if stmt.Value != nil {
		if re.currentFunction == INITIALIZER {
			re.reporter.ResolveError(stmt.Keyword, "Can't return a value from an initializer.")
		}
		re.resolve(stmt.Value)
	}
}

func (re *Resolver) VisitVarStmtStmtVoid(stmt *ast.VarStmt) {
	re.declare(stmt.Name)
	if stmt.Initializer != nil {
		re.resolve(stmt.Initializer)
	}
	re.define(stmt.Name)
}

func (re *Resolver) VisitWhileStmtVoid(while *ast.While) {
	re.resolve(while.Condition)
	re.resolve(while.Body)
}

func (re *Resolver) VisitAssignExprVoid(assign *ast.Assign) {
	re.resolve(assign.Value)
	re.resolveLocal(assign, assign.Name)
}

func (re *Resolver) VisitBinaryExprVoid(binary *ast.Binary) {
	re.resolve(binary.Left)
	re.resolve(binary.Right)
}

func (re *Resolver) VisitCallExprVoid(expr *ast.Call) {
	re.resolve(expr.Callee)
	for _, arg := range expr.Arguments {
		re.resolve(arg)
	}
}

func (re *Resolver) VisitGetExprVoid(get *ast.Get) {
	re.resolve(get.Object)
}

func (re *Resolver) VisitGroupingExprVoid(grouping *ast.Grouping) {
	re.resolve(grouping.Expression)
}

func (re *Resolver) VisitLiteralExprVoid(literal *ast.Literal) {
	// no-op
}

func (re *Resolver) VisitLogicalExprVoid(logical *ast.Logical) {
	re.resolve(logical.Left)
	re.resolve(logical.Right)
}

func (re *Resolver) VisitSetExprVoid(set *ast.Set) {
	re.resolve(set.Value)
	re.resolve(set.Object)
}

func (re *Resolver) VisitSuperExprVoid(super *ast.Super) {
	if re.curentClass == NOCLASS {
		re.reporter.ResolveError(super.Keyword, "Can't use 'super' outside of a class.")
	} else if re.curentClass != SUBCLASS {
		re.reporter.ResolveError(super.Keyword, "Can't use 'super' in a class with no superclass.")
	}
	re.resolveLocal(super, super.Keyword)
}

func (re *Resolver) VisitThisExprVoid(this *ast.This) {
	if re.curentClass == NOCLASS {
		re.reporter.ResolveError(this.Keyword, "Can't use 'this' outside of a class.")
		return
	}
	re.resolveLocal(this, this.Keyword)
}

func (re *Resolver) VisitUnaryExprVoid(unary *ast.Unary) {
	re.resolve(unary.Right)
}

func (re *Resolver) VisitVariableExprVoid(variable *ast.Variable) {
	isEmpty := re.scopes.isEmpty()
	v, ok := re.scopes.peek()[variable.Name.Lexeme]
	if !isEmpty && v == false && ok {
		re.reporter.ResolveError(variable.Name, "Can't read local variable in its own initializer.")
	}
	re.resolveLocal(variable, variable.Name)
}

type stack struct {
	elems []map[string]bool
}

func (s *stack) size() int {
	return len(s.elems)
}

func (s *stack) isEmpty() bool {
	return s.size() == 0
}

func (s *stack) push(elem map[string]bool) {
	s.elems = append(s.elems, elem)
}

func (s *stack) pop() map[string]bool {
	if s.size() == 0 {
		return nil
	}
	last := s.size() - 1
	elem := s.elems[last]
	s.elems = s.elems[0:last]
	return elem
}

func (s *stack) peek() map[string]bool {
	if s.size() == 0 {
		return nil
	}
	last := s.size() - 1
	return s.elems[last]
}

func (s *stack) get(i int) map[string]bool {
	// TODO: check bounds
	return s.elems[i]
}
