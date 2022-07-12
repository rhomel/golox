package interpreter

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	ast "rhomel.com/crafting-interpreters-go/pkg/ast/gen"
	"rhomel.com/crafting-interpreters-go/pkg/scanner"
	"rhomel.com/crafting-interpreters-go/pkg/util/check"
	"rhomel.com/crafting-interpreters-go/pkg/util/exit"
)

type Interpreter struct {
	reporter    RuntimeErrorReporter
	globals     *Environment
	environment *Environment
}

type RuntimeErrorReporter interface {
	RuntimeError(token scanner.Token, message string)
}

func NewInterpreter(reporter RuntimeErrorReporter) *Interpreter {
	globals := NewEnvironment(nil)
	globals.Define("clock", &nativeClock{})
	return &Interpreter{reporter, globals, globals}
}

func (in *Interpreter) Interpret(statements []ast.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			// TODO
			if err, ok := r.(*RuntimeError); ok {
				in.reporter.RuntimeError(err.token, err.message)
			} else {
				in.reporter.RuntimeError(scanner.Token{}, fmt.Sprintf("%v", r)) // TODO
			}
		}
	}()
	for _, stmt := range statements {
		if stmt != nil { // TODO: on parse errors a nil entry will be inserted
			in.execute(stmt)
		}
	}
}

func (in *Interpreter) Accept(elem interface{}) interface{} {
	// Go has no dynamic dispatch and inheritance so we have to resort to a type switch
	switch v := elem.(type) {
	case *ast.Binary:
		return v.Accept(in)
	case *ast.Call:
		return v.Accept(in)
	case *ast.Grouping:
		return v.Accept(in)
	case *ast.Literal:
		return v.Accept(in)
	case *ast.Logical:
		return v.Accept(in)
	case *ast.Unary:
		return v.Accept(in)
	case *ast.Variable:
		return v.Accept(in)
	case *ast.Assign:
		return v.Accept(in)
	default:
		exit.Exitf(exit.ExitSyntaxError, "unsupported expression type: %s", check.TypeOf(elem))
		return nil
	}
}

func (in *Interpreter) VisitBinaryExpr(binary *ast.Binary) interface{} {
	left := in.evaluate(binary.Left)
	right := in.evaluate(binary.Right)

	switch binary.Operator.Typ {
	case scanner.GREATER:
		in.checkNumberOperands(binary.Operator, left, right)
		return mustDouble(left) > mustDouble(right)
	case scanner.GREATER_EQUAL:
		in.checkNumberOperands(binary.Operator, left, right)
		return mustDouble(left) >= mustDouble(right)
	case scanner.LESS:
		in.checkNumberOperands(binary.Operator, left, right)
		return mustDouble(left) < mustDouble(right)
	case scanner.LESS_EQUAL:
		in.checkNumberOperands(binary.Operator, left, right)
		return mustDouble(left) <= mustDouble(right)
	case scanner.BANG_EQUAL:
		return !in.isEqual(left, right)
	case scanner.EQUAL_EQUAL:
		return in.isEqual(left, right)
	case scanner.MINUS:
		in.checkNumberOperands(binary.Operator, left, right)
		return mustDouble(left) - mustDouble(right)
	case scanner.SLASH:
		in.checkNumberOperands(binary.Operator, left, right)
		return mustDouble(left) / mustDouble(right)
	case scanner.STAR:
		in.checkNumberOperands(binary.Operator, left, right)
		return mustDouble(left) * mustDouble(right)
	case scanner.PLUS:
		leftDouble, leftIsDouble := left.(float64)
		rightDouble, rightIsDouble := right.(float64)
		if leftIsDouble && rightIsDouble {
			return leftDouble + rightDouble
		}
		leftString, leftIsString := left.(string)
		rightString, rightIsString := right.(string)
		if leftIsString && rightIsString {
			return leftString + rightString
		}
		if leftIsDouble && rightIsString {
			// TODO: improve this to reference the actual token literal
			panic(&RuntimeError{binary.Operator, fmt.Sprintf("Left operand '%f' is double but right operand '%s' is string.", leftDouble, rightString)})
		}
		if leftIsString && rightIsDouble {
			// TODO: improve this to reference the actual token literal
			panic(&RuntimeError{binary.Operator, fmt.Sprintf("Left operand '%s' is string but right operand '%f' is double.", leftString, rightDouble)})
		}
		if !leftIsDouble && !leftIsString {
			panic(fmt.Sprintf("expected either double or string value for left expression: %v", left)) // TODO
		}
		if !rightIsDouble && !rightIsString {
			panic(fmt.Sprintf("expected either double or string value for right expression: %v", right)) // TODO
		}
	}
	return nil // unreachable
}

func (in *Interpreter) VisitCallExpr(expr *ast.Call) interface{} {
	callee := in.evaluate(expr.Callee)
	var arguments []interface{}
	for _, argument := range expr.Arguments {
		arguments = append(arguments, in.evaluate(argument))
	}
	function, ok := callee.(LoxCallable)
	if !ok {
		panic(&RuntimeError{expr.Paren, "Can only call functions and classes."})
	}
	if expected, got := function.Arity(), len(arguments); expected != got {
		panic(&RuntimeError{expr.Paren, fmt.Sprintf("Expected %d arguments but got %d.", expected, got)})
	}
	return function.Call(in, arguments)
}

func (in *Interpreter) VisitGroupingExpr(grouping *ast.Grouping) interface{} {
	return in.evaluate(grouping)
}

func (in *Interpreter) VisitLiteralExpr(literal *ast.Literal) interface{} {
	return literal.Value
}

func (in *Interpreter) VisitLogicalExpr(logical *ast.Logical) interface{} {
	left := in.evaluate(logical.Left)
	leftIsTruthy := in.isTruthy(left)
	switch logical.Operator.Typ {
	case scanner.OR:
		if leftIsTruthy {
			return left
		}
	case scanner.AND:
		if !leftIsTruthy {
			return left
		}
	default:
		panic(fmt.Sprintf("unsupported logical operator: %s", logical.Operator.Typ))
	}
	return in.evaluate(logical.Right)
}

func (in *Interpreter) VisitUnaryExpr(unary *ast.Unary) interface{} {
	right := in.evaluate(unary.Right)

	switch unary.Operator.Typ {
	case scanner.MINUS:
		in.checkNumberOperand(unary.Operator, right)
		return -mustDouble(right)
	case scanner.BANG:
		return !in.isTruthy(right)
	}
	return nil // unreachable
}

func (in *Interpreter) VisitVariableExpr(variable *ast.Variable) interface{} {
	return in.environment.Get(variable.Name)
}

func (in *Interpreter) checkNumberOperand(operator scanner.Token, operand interface{}) {
	if _, ok := operand.(float64); ok {
		return
	}
	panic(&RuntimeError{operator, "Operand must be a number."})
}

func (in *Interpreter) checkNumberOperands(operator scanner.Token, left, right interface{}) {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)
	if leftOk && rightOk {
		return
	}
	panic(&RuntimeError{operator, "Operands must be a numbers."})
}

// isTruthy returns false only for 'nil' and the boolean value false
func (in *Interpreter) isTruthy(it interface{}) bool {
	if check.IsNil(it) {
		return false
	}
	if b, ok := it.(bool); ok {
		return b
	}
	return true
}

func (in *Interpreter) isEqual(a, b interface{}) bool {
	if check.IsNil(a) && check.IsNil(b) {
		return true
	}
	if check.IsNil(a) {
		return false
	}
	return a == b // TODO: will this work?
}

func (in *Interpreter) stringify(it interface{}) string {
	if check.IsNil(it) {
		return "nil"
	}
	if double, ok := it.(float64); ok {
		text := fmt.Sprintf("%f", double)
		text = maybeInteger(text)
		return text
	}
	if str, ok := it.(string); ok {
		return str
	}
	if str, ok := it.(fmt.Stringer); ok {
		return str.String()
	}
	if b, ok := it.(bool); ok {
		if b {
			return "true"
		}
		return "false"
	}
	panic(fmt.Sprintf("unhandled type in stringify: %v", it)) // TODO
}

func (in *Interpreter) evaluate(expr ast.Expr) interface{} {
	return in.Accept(expr)
}

func (in *Interpreter) execute(stmt ast.Stmt) {
	switch v := stmt.(type) {
	case *ast.IfStmt:
		v.AcceptVoid(in)
	case *ast.Block:
		v.AcceptVoid(in)
	case *ast.Expression:
		v.AcceptVoid(in)
	case *ast.Function:
		v.AcceptVoid(in)
	case *ast.Print:
		v.AcceptVoid(in)
	case *ast.ReturnStmt:
		v.AcceptVoid(in)
	case *ast.VarStmt:
		v.AcceptVoid(in)
	case *ast.While:
		v.AcceptVoid(in)
	default:
		exit.Exitf(exit.ExitSyntaxError, "unsupported statement: %s", check.TypeOf(stmt))
	}
}

func (in *Interpreter) executeBlock(statements []ast.Stmt, environment *Environment) {
	previous := in.environment // save the current environment scope before executing block scope
	defer func() {
		in.environment = previous
	}()
	in.environment = environment
	for _, statement := range statements {
		in.execute(statement)
	}
}

func (in *Interpreter) VisitBlockStmtVoid(block *ast.Block) {
	in.executeBlock(block.Statements, NewEnvironment(in.environment))
}

func (in *Interpreter) VisitExpressionStmtVoid(stmt *ast.Expression) {
	in.evaluate(stmt.Expression)
}

func (in *Interpreter) VisitFunctionStmtVoid(stmt *ast.Function) {
	function := NewLoxFunction(stmt, in.environment)
	in.environment.Define(stmt.Name.Lexeme, function)
}

func (in *Interpreter) VisitIfStmtStmtVoid(stmt *ast.IfStmt) {
	condition := in.evaluate(stmt.Condition)
	if in.isTruthy(condition) {
		in.execute(stmt.ThenBranch)
	} else {
		if stmt.ElseBranch != nil {
			in.execute(stmt.ElseBranch)
		}
	}
}

func (in *Interpreter) VisitPrintStmtVoid(stmt *ast.Print) {
	value := in.evaluate(stmt.Expression)
	fmt.Println(in.stringify(value))
}

func (in *Interpreter) VisitReturnStmtStmtVoid(stmt *ast.ReturnStmt) {
	var value interface{}
	if stmt.Value != nil {
		value = in.evaluate(stmt.Value)
	}
	panic(&Return{value})
}

func (in *Interpreter) VisitVarStmtStmtVoid(stmt *ast.VarStmt) {
	var value interface{}
	if stmt.Initializer != nil {
		value = in.evaluate(stmt.Initializer)
	}
	in.environment.Define(stmt.Name.Lexeme, value)
}

func (in *Interpreter) VisitWhileStmtVoid(while *ast.While) {
	for in.isTruthy(in.evaluate(while.Condition)) {
		in.execute(while.Body)
	}
}

func (in *Interpreter) VisitAssignExpr(assign *ast.Assign) interface{} {
	value := in.evaluate(assign.Value)
	in.environment.Assign(assign.Name, value)
	return value
}

func mustDouble(it interface{}) float64 {
	if double, ok := it.(float64); ok {
		return double
	}
	panic(fmt.Sprintf("expression value did not evaluate to a float64: %v", it)) // TODO
}

var matchTrailingZeros = regexp.MustCompile("\\.0+$")

func maybeInteger(num string) string {
	if matchTrailingZeros.Match([]byte(num)) {
		return num[:strings.Index(num, ".")]
	}
	return num
}

type RuntimeError struct {
	token   scanner.Token
	message string
}

func (e *RuntimeError) Error() string {
	return e.message
}

var _ error = (*RuntimeError)(nil)

type LoxCallable interface {
	Call(*Interpreter, []interface{}) interface{}
	Arity() int
}

type nativeClock struct{}

var _ LoxCallable = (*nativeClock)(nil)

func (*nativeClock) Arity() int {
	return 0
}

func (*nativeClock) Call(in *Interpreter, arguments []interface{}) interface{} {
	return (float64)(time.Now().Unix())
}

func (*nativeClock) String() string {
	return "<native fn>"
}

type LoxFunction struct {
	declaration *ast.Function
	closure     *Environment
}

var _ LoxCallable = (*LoxFunction)(nil)

func NewLoxFunction(declaration *ast.Function, closure *Environment) *LoxFunction {
	return &LoxFunction{declaration, closure}
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) Call(in *Interpreter, arguments []interface{}) (ret interface{}) {
	environment := NewEnvironment(f.closure)
	for i := range f.declaration.Params {
		environment.Define(f.declaration.Params[i].Lexeme, arguments[i])
	}
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(*Return); ok {
				ret = re.value
			} else {
				panic(r)
			}
		}
	}()
	in.executeBlock(f.declaration.Body, environment)
	return nil
}

func (f *LoxFunction) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}

type Return struct {
	value interface{}
}
