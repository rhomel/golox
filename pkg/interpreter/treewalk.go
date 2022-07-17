package interpreter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	ast "github.com/rhomel/golox/pkg/ast/gen"
	"github.com/rhomel/golox/pkg/scanner"
	"github.com/rhomel/golox/pkg/util/check"
	"github.com/rhomel/golox/pkg/util/exit"
)

type TreeWalkInterpreter struct {
	reporter    RuntimeErrorReporter
	globals     *Environment
	environment *Environment
	locals      map[ast.Expr]int
}

var _ Interpreter = (*TreeWalkInterpreter)(nil)

type RuntimeErrorReporter interface {
	RuntimeError(token scanner.Token, message string)
}

func NewTreeWalkInterpreter(reporter RuntimeErrorReporter) *TreeWalkInterpreter {
	globals := NewEnvironment(nil)
	globals.Define("clock", &nativeClock{})
	return &TreeWalkInterpreter{reporter, globals, globals, make(map[ast.Expr]int)}
}

func (in *TreeWalkInterpreter) Interpret(statements []ast.Stmt) {
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

func (in *TreeWalkInterpreter) Accept(elem interface{}) interface{} {
	// Go has no dynamic dispatch and inheritance so we have to resort to a type switch
	switch v := elem.(type) {
	case *ast.Binary:
		return v.Accept(in)
	case *ast.Call:
		return v.Accept(in)
	case *ast.Get:
		return v.Accept(in)
	case *ast.Grouping:
		return v.Accept(in)
	case *ast.Literal:
		return v.Accept(in)
	case *ast.Logical:
		return v.Accept(in)
	case *ast.Set:
		return v.Accept(in)
	case *ast.Super:
		return v.Accept(in)
	case *ast.This:
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

func (in *TreeWalkInterpreter) VisitBinaryExpr(binary *ast.Binary) interface{} {
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
		if (!leftIsDouble && !leftIsString) || (!rightIsDouble && !rightIsString) {
			panic(&RuntimeError{binary.Operator, "Operands must be two numbers or two strings."})
		}
	}
	return nil // unreachable
}

func (in *TreeWalkInterpreter) VisitCallExpr(expr *ast.Call) interface{} {
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

func (in *TreeWalkInterpreter) VisitGetExpr(get *ast.Get) interface{} {
	object := in.evaluate(get.Object)
	if instance, ok := object.(*LoxInstance); ok {
		return instance.Get(get.Name)
	}
	panic(&RuntimeError{get.Name, "Only instances have properties."})
}

func (in *TreeWalkInterpreter) VisitGroupingExpr(grouping *ast.Grouping) interface{} {
	return in.evaluate(grouping.Expression)
}

func (in *TreeWalkInterpreter) VisitLiteralExpr(literal *ast.Literal) interface{} {
	return literal.Value
}

func (in *TreeWalkInterpreter) VisitLogicalExpr(logical *ast.Logical) interface{} {
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

func (in *TreeWalkInterpreter) VisitSetExpr(set *ast.Set) interface{} {
	object := in.evaluate(set.Object)
	if instance, ok := object.(*LoxInstance); ok {
		value := in.evaluate(set.Value)
		instance.Set(set.Name, value)
		return value
	}
	panic(&RuntimeError{set.Name, fmt.Sprintf("Only instances have fields.")})
}

func (in *TreeWalkInterpreter) VisitSuperExpr(super *ast.Super) interface{} {
	distance, ok := in.locals[super]
	if !ok {
		panic(&RuntimeError{super.Keyword, "no resolved local"})
	}
	superclass, ok := in.environment.GetAt(distance, "super").(*LoxClass)
	if !ok {
		panic(&RuntimeError{super.Keyword, "didn't find super class"})
	}
	object, ok := in.environment.GetAt(distance-1, "this").(*LoxInstance)
	if !ok {
		panic(&RuntimeError{super.Keyword, "didn't find super class instance"})
	}
	method := superclass.FindMethod(super.Method.Lexeme)
	if method == nil {
		panic(&RuntimeError{super.Method, fmt.Sprintf("Undefined property '%s'.", super.Method.Lexeme)})
	}
	return method.Bind(object)
}

func (in *TreeWalkInterpreter) VisitThisExpr(this *ast.This) interface{} {
	return in.lookUpVariable(this.Keyword, this)
}

func (in *TreeWalkInterpreter) VisitUnaryExpr(unary *ast.Unary) interface{} {
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

func (in *TreeWalkInterpreter) VisitVariableExpr(variable *ast.Variable) interface{} {
	return in.lookUpVariable(variable.Name, variable)
}

func (in *TreeWalkInterpreter) lookUpVariable(name scanner.Token, expr ast.Expr) interface{} {
	distance, ok := in.locals[expr]
	if ok {
		return in.environment.GetAt(distance, name.Lexeme)
	} else {
		return in.globals.Get(name)
	}
}

func (in *TreeWalkInterpreter) checkNumberOperand(operator scanner.Token, operand interface{}) {
	if _, ok := operand.(float64); ok {
		return
	}
	panic(&RuntimeError{operator, "Operand must be a number."})
}

func (in *TreeWalkInterpreter) checkNumberOperands(operator scanner.Token, left, right interface{}) {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)
	if leftOk && rightOk {
		return
	}
	panic(&RuntimeError{operator, "Operands must be numbers."})
}

// isTruthy returns false only for 'nil' and the boolean value false
func (in *TreeWalkInterpreter) isTruthy(it interface{}) bool {
	if check.IsNil(it) {
		return false
	}
	if b, ok := it.(bool); ok {
		return b
	}
	return true
}

func (in *TreeWalkInterpreter) isEqual(a, b interface{}) bool {
	if check.IsNil(a) && check.IsNil(b) {
		return true
	}
	if check.IsNil(a) {
		return false
	}
	return a == b // TODO: will this work?
}

func (in *TreeWalkInterpreter) stringify(it interface{}) string {
	if check.IsNil(it) {
		return "nil"
	}
	if double, ok := it.(float64); ok {
		text := strconv.FormatFloat(double, 'f', -1, 64)
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

func (in *TreeWalkInterpreter) evaluate(expr ast.Expr) interface{} {
	return in.Accept(expr)
}

func (in *TreeWalkInterpreter) execute(stmt ast.Stmt) {
	switch v := stmt.(type) {
	case *ast.Class:
		v.AcceptVoid(in)
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

func (in *TreeWalkInterpreter) executeBlock(statements []ast.Stmt, environment *Environment) {
	previous := in.environment // save the current environment scope before executing block scope
	defer func() {
		in.environment = previous
	}()
	in.environment = environment
	for _, statement := range statements {
		in.execute(statement)
	}
}

func (in *TreeWalkInterpreter) Resolve(expr ast.Expr, depth int) {
	in.locals[expr] = depth
}

func (in *TreeWalkInterpreter) VisitBlockStmtVoid(block *ast.Block) {
	in.executeBlock(block.Statements, NewEnvironment(in.environment))
}

func (in *TreeWalkInterpreter) VisitExpressionStmtVoid(stmt *ast.Expression) {
	in.evaluate(stmt.Expression)
}

func (in *TreeWalkInterpreter) VisitFunctionStmtVoid(stmt *ast.Function) {
	function := NewLoxFunction(stmt, in.environment, false)
	in.environment.Define(stmt.Name.Lexeme, function)
}

func (in *TreeWalkInterpreter) VisitClassStmtVoid(class *ast.Class) {
	var superklass *LoxClass
	if class.Superclass != nil {
		superclass := in.evaluate(class.Superclass)
		if v, ok := superclass.(*LoxClass); ok {
			superklass = v
		} else {
			panic(&RuntimeError{class.Superclass.Name, "Superclass must be a class."})
		}
	}
	in.environment.Define(class.Name.Lexeme, nil)
	if class.Superclass != nil {
		in.environment = NewEnvironment(in.environment)
		in.environment.Define("super", superklass)
	}
	methods := make(map[string]*LoxFunction)
	for _, method := range class.Methods {
		isInitializer := method.Name.Lexeme == "init"
		function := NewLoxFunction(method, in.environment, isInitializer)
		methods[method.Name.Lexeme] = function
	}
	klass := NewLoxClass(class.Name.Lexeme, superklass, methods)
	if superklass != nil {
		in.environment = in.environment.enclosing
	}
	in.environment.Assign(class.Name, klass)
}

func (in *TreeWalkInterpreter) VisitIfStmtStmtVoid(stmt *ast.IfStmt) {
	condition := in.evaluate(stmt.Condition)
	if in.isTruthy(condition) {
		in.execute(stmt.ThenBranch)
	} else {
		if stmt.ElseBranch != nil {
			in.execute(stmt.ElseBranch)
		}
	}
}

func (in *TreeWalkInterpreter) VisitPrintStmtVoid(stmt *ast.Print) {
	value := in.evaluate(stmt.Expression)
	fmt.Println(in.stringify(value))
}

func (in *TreeWalkInterpreter) VisitReturnStmtStmtVoid(stmt *ast.ReturnStmt) {
	var value interface{}
	if stmt.Value != nil {
		value = in.evaluate(stmt.Value)
	}
	panic(&Return{value})
}

func (in *TreeWalkInterpreter) VisitVarStmtStmtVoid(stmt *ast.VarStmt) {
	var value interface{}
	if stmt.Initializer != nil {
		value = in.evaluate(stmt.Initializer)
	}
	in.environment.Define(stmt.Name.Lexeme, value)
}

func (in *TreeWalkInterpreter) VisitWhileStmtVoid(while *ast.While) {
	for in.isTruthy(in.evaluate(while.Condition)) {
		in.execute(while.Body)
	}
}

func (in *TreeWalkInterpreter) VisitAssignExpr(assign *ast.Assign) interface{} {
	value := in.evaluate(assign.Value)
	distance, ok := in.locals[assign]
	if ok {
		in.environment.AssignAt(distance, assign.Name, value)
	} else {
		in.globals.Assign(assign.Name, value)
	}
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
	Call(*TreeWalkInterpreter, []interface{}) interface{}
	Arity() int
}

type nativeClock struct{}

var _ LoxCallable = (*nativeClock)(nil)

func (*nativeClock) Arity() int {
	return 0
}

func (*nativeClock) Call(in *TreeWalkInterpreter, arguments []interface{}) interface{} {
	return (float64)(time.Now().Unix())
}

func (*nativeClock) String() string {
	return "<native fn>"
}

type LoxFunction struct {
	declaration   *ast.Function
	closure       *Environment
	isInitializer bool
}

var _ LoxCallable = (*LoxFunction)(nil)

func NewLoxFunction(declaration *ast.Function, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{declaration, closure, isInitializer}
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) Call(in *TreeWalkInterpreter, arguments []interface{}) (ret interface{}) {
	environment := NewEnvironment(f.closure)
	for i := range f.declaration.Params {
		environment.Define(f.declaration.Params[i].Lexeme, arguments[i])
	}
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(*Return); ok {
				if f.isInitializer {
					ret = f.closure.GetAt(0, "this")
				} else {
					ret = re.value
				}
			} else {
				panic(r)
			}
		}
	}()
	in.executeBlock(f.declaration.Body, environment)
	if f.isInitializer {
		return f.closure.GetAt(0, "this")
	}
	return nil
}

func (f *LoxFunction) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnvironment(f.closure)
	environment.Define("this", instance)
	return NewLoxFunction(f.declaration, environment, f.isInitializer)
}

type Return struct {
	value interface{}
}

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

var _ LoxCallable = (*LoxClass)(nil)

func NewLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name, superclass, methods}
}

func (c *LoxClass) Arity() int {
	initializer := c.FindMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}

func (c *LoxClass) Call(in *TreeWalkInterpreter, args []interface{}) interface{} {
	instance := NewLoxInstance(c)
	initializer := c.FindMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(in, args)
	}
	return instance
}

func (c *LoxClass) String() string {
	return c.name
}

func (c *LoxClass) FindMethod(name string) *LoxFunction {
	if method, ok := c.methods[name]; ok {
		return method
	}
	if c.superclass != nil {
		return c.superclass.FindMethod(name)
	}
	return nil
}

type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class, make(map[string]interface{})}
}

func (i *LoxInstance) Get(name scanner.Token) interface{} {
	if field, ok := i.fields[name.Lexeme]; ok {
		return field
	}
	if method := i.class.FindMethod(name.Lexeme); method != nil {
		return method.Bind(i)
	}
	panic(&RuntimeError{name, fmt.Sprintf("Undefined property '%s'.", name.Lexeme)})
}

func (i *LoxInstance) Set(name scanner.Token, value interface{}) {
	i.fields[name.Lexeme] = value
}

func (i *LoxInstance) String() string {
	return i.class.name + " instance"
}
