package vm

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

var debugPrintCode = false

var current *Compiler

func currentChunk() *Chunk {
	return current.function.chunk
}

func compile(source string) *ObjectFunction {
	scanner := InitScanner(source)
	parser.scanner = scanner
	parser.hadError = false
	parser.panicMode = false
	compiler := Compiler{}
	parser.InitCompiler(&compiler, TypeScript)

	parser.advance()

	for !parser.match(TOKEN_EOF) {
		parser.declaration()
	}

	function := endCompiler()
	if parser.hadError {
		return nil
	}
	return function
}

func endCompiler() *ObjectFunction {
	parser.emitReturn()
	function := current.function
	if debugPrintCode {
		if !parser.hadError {
			name := "<script>"
			if function.name != nil {
				name = function.name.String
			}
			currentChunk().Disassemble(name)
		}
	}
	return function
}

func (p *Parser) beginScope() {
	current.scopeDepth++
}

func (p *Parser) endScope() {
	current.scopeDepth--

	for current.localCount > 0 && current.locals[current.localCount-1].depth > current.scopeDepth {
		p.emitByte(OP_POP)
		current.localCount--
	}
}

func _chapter_16_compile(source string) {
	scanner := InitScanner(source)

	line := -1

	for {
		token := scanner.ScanToken()
		if token.Line != line {
			fmt.Printf("%4d ", token.Line)
			line = token.Line
		} else {
			fmt.Printf("   | ")
		}
		var str string
		if token.Type == TOKEN_ERROR {
			str = token.Error
		} else if token.Type == TOKEN_EOF {
			str = "EOF"
		} else {
			str = string(scanner.source[token.Start : token.Start+token.Length])
		}
		fmt.Printf("%2d '%s'\n", token.Type, str)
		if token.Type == TOKEN_EOF {
			break
		}
	}
}

type Parser struct {
	hadError  bool
	panicMode bool
	current   Token
	previous  Token
	scanner   *Scanner
}

func InitParser() *Parser {
	return &Parser{}
}

var parser = InitParser()

type Precedence int

const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

type ParseRule struct {
	prefix     ParseFn
	infix      ParseFn
	precedence Precedence
}

type Local struct {
	name  Token
	depth int
}

type FunctionType int

const (
	TypeFunction FunctionType = iota
	TypeScript
)

// AsString returns the local token's string from the parser source
func (l Local) AsString(source []rune) string {
	return l.name.StartAsString(source)
}

type Compiler struct {
	function *ObjectFunction
	typ      FunctionType

	locals     [UINT8_COUNT]Local
	localCount int
	scopeDepth int
}

// Return true if scopeDepth is greater than 0. This is not part of the book,
// it is a helper to make the other code easier to read.
func (c *Compiler) inLocalScope() bool {
	return c.scopeDepth > 0
}

type ParseFn func(bool)

var rules = make([]ParseRule, TOKEN_EOF+1)

func init() {
	rules[TOKEN_LEFT_PAREN] = ParseRule{parser.grouping, nil, PREC_NONE}
	rules[TOKEN_RIGHT_PAREN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_LEFT_BRACE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_RIGHT_BRACE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_COMMA] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_DOT] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_MINUS] = ParseRule{parser.unary, parser.binary, PREC_TERM}
	rules[TOKEN_PLUS] = ParseRule{nil, parser.binary, PREC_TERM}
	rules[TOKEN_SEMICOLON] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_SLASH] = ParseRule{nil, parser.binary, PREC_FACTOR}
	rules[TOKEN_STAR] = ParseRule{nil, parser.binary, PREC_FACTOR}
	rules[TOKEN_BANG] = ParseRule{parser.unary, nil, PREC_NONE}
	rules[TOKEN_BANG_EQUAL] = ParseRule{nil, parser.binary, PREC_EQUALITY}
	rules[TOKEN_EQUAL] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_EQUAL_EQUAL] = ParseRule{nil, parser.binary, PREC_EQUALITY}
	rules[TOKEN_GREATER] = ParseRule{nil, parser.binary, PREC_COMPARISON}
	rules[TOKEN_GREATER_EQUAL] = ParseRule{nil, parser.binary, PREC_COMPARISON}
	rules[TOKEN_LESS] = ParseRule{nil, parser.binary, PREC_COMPARISON}
	rules[TOKEN_LESS_EQUAL] = ParseRule{nil, parser.binary, PREC_COMPARISON}
	rules[TOKEN_IDENTIFIER] = ParseRule{parser.variable, nil, PREC_NONE}
	rules[TOKEN_STRING] = ParseRule{parser.string, nil, PREC_NONE}
	rules[TOKEN_NUMBER] = ParseRule{parser.number, nil, PREC_NONE}
	rules[TOKEN_AND] = ParseRule{nil, parser.and_, PREC_AND}
	rules[TOKEN_CLASS] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_ELSE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FALSE] = ParseRule{parser.literal, nil, PREC_NONE}
	rules[TOKEN_FOR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FUN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_IF] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_NIL] = ParseRule{parser.literal, nil, PREC_NONE}
	rules[TOKEN_OR] = ParseRule{nil, parser.or_, PREC_OR}
	rules[TOKEN_PRINT] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_RETURN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_SUPER] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_THIS] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_TRUE] = ParseRule{parser.literal, nil, PREC_NONE}
	rules[TOKEN_VAR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_WHILE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_ERROR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_EOF] = ParseRule{nil, nil, PREC_NONE}
}

func (p *Parser) advance() {
	p.previous = p.current

	for {
		p.current = p.scanner.ScanToken()
		if p.current.Type != TOKEN_ERROR {
			break
		}
		p.errorAtCurrent(p.current.StartAsString(p.scanner.source)) // TODO: verify the purpose of using the token as the message
	}
}

func (p *Parser) expression() {
	p.parsePrecedence(PREC_ASSIGNMENT)
}

func (p *Parser) block() {
	for !p.check(TOKEN_RIGHT_BRACE) && !p.check(TOKEN_EOF) {
		p.declaration()
	}

	p.consume(TOKEN_RIGHT_BRACE, "Expect '}' after block.")
}

func (p *Parser) expressionStatement() {
	p.expression()
	p.consume(TOKEN_SEMICOLON, "Expect ';' after expression.")
	p.emitByte(OP_POP)
}

func (p *Parser) forStatement() {
	p.beginScope()
	p.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'for'.")
	if p.match(TOKEN_SEMICOLON) {
		// No initializer.
	} else if p.match(TOKEN_VAR) {
		p.varDeclaration()
	} else {
		p.expressionStatement()
	}

	loopStart := currentChunk().Count()
	exitJump := -1
	if !p.match(TOKEN_SEMICOLON) {
		p.expression()
		p.consume(TOKEN_SEMICOLON, "Expect ';' after loop condition.")

		// Jump out of the loop if the condition is false.
		exitJump = p.emitJump(OP_JUMP_IF_FALSE)
		p.emitByte(OP_POP)
	}

	if !p.match(TOKEN_RIGHT_PAREN) {
		bodyJump := p.emitJump(OP_JUMP)
		incrementStart := currentChunk().Count()
		p.expression()
		p.emitByte(OP_POP)
		p.consume(TOKEN_RIGHT_PAREN, "Expect ')' after for clauses.")

		p.emitLoop(uint8(loopStart))
		loopStart = incrementStart
		p.patchJump(bodyJump)
	}

	p.statement()
	p.emitLoop(uint8(loopStart))

	if exitJump != -1 {
		p.patchJump(exitJump)
		p.emitByte(OP_POP)
	}

	p.endScope()
}

func (p *Parser) ifStatement() {
	p.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'if'.")
	p.expression()
	p.consume(TOKEN_RIGHT_PAREN, "Expect ')' after condition.")

	var thenJump int = p.emitJump(OP_JUMP_IF_FALSE)
	p.emitByte(OP_POP)
	p.statement()

	var elseJump int = p.emitJump(OP_JUMP)

	p.patchJump(thenJump)
	p.emitByte(OP_POP)

	if p.match(TOKEN_ELSE) {
		p.statement()
	}
	p.patchJump(elseJump)
}

func (p *Parser) printStatement() {
	p.expression()
	p.consume(TOKEN_SEMICOLON, "Expect ';' after value.")
	p.emitByte(OP_PRINT)
}

func (p *Parser) whileStatement() {
	loopStart := currentChunk().Count()
	p.consume(TOKEN_LEFT_PAREN, "Expect '(' after 'while'.")
	p.expression()
	p.consume(TOKEN_RIGHT_PAREN, "Expect ')' after 'while'.")

	exitJump := p.emitJump(OP_JUMP_IF_FALSE)
	p.emitByte(OP_POP)
	p.statement()
	p.emitLoop(uint8(loopStart))

	p.patchJump(exitJump)
	p.emitByte(OP_POP)
}

func (p *Parser) synchronize() {
	p.panicMode = false

	for p.current.Type != TOKEN_EOF {
		if p.previous.Type == TOKEN_SEMICOLON {
			return
		}
		switch p.current.Type {
		case TOKEN_CLASS:
			return
		case TOKEN_FUN:
			return
		case TOKEN_VAR:
			return
		case TOKEN_FOR:
			return
		case TOKEN_IF:
			return
		case TOKEN_WHILE:
			return
		case TOKEN_PRINT:
			return
		case TOKEN_RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) declaration() {
	if p.match(TOKEN_VAR) {
		p.varDeclaration()
	} else {
		p.statement()
	}

	if p.panicMode {
		p.synchronize()
	}
}

func (p *Parser) varDeclaration() {
	var global uint8 = p.parseVariable("Expect variable name.")

	if p.match(TOKEN_EQUAL) {
		p.expression()
	} else {
		p.emitByte(OP_NIL)
	}

	p.consume(TOKEN_SEMICOLON, "Expect ';' after variable declaration.")
	p.defineVariable(global)
}

func (p *Parser) statement() {
	if p.match(TOKEN_PRINT) {
		p.printStatement()
	} else if p.match(TOKEN_FOR) {
		p.forStatement()
	} else if p.match(TOKEN_IF) {
		p.ifStatement()
	} else if p.match(TOKEN_WHILE) {
		p.whileStatement()
	} else if p.match(TOKEN_LEFT_BRACE) {
		p.beginScope()
		p.block()
		p.endScope()
	} else {
		p.expressionStatement()
	}
}

func (p *Parser) binary(canAssign bool) {
	operatorType := p.previous.Type
	rule := getRule(operatorType)
	p.parsePrecedence(Precedence(rule.precedence + 1))
	switch operatorType {
	case TOKEN_BANG_EQUAL:
		p.emitBytes(OP_EQUAL, OP_NOT)
	case TOKEN_EQUAL_EQUAL:
		p.emitByte(OP_EQUAL)
	case TOKEN_GREATER:
		p.emitByte(OP_GREATER)
	case TOKEN_GREATER_EQUAL:
		p.emitBytes(OP_LESS, OP_NOT)
	case TOKEN_LESS:
		p.emitByte(OP_LESS)
	case TOKEN_LESS_EQUAL:
		p.emitBytes(OP_GREATER, OP_NOT)
	case TOKEN_PLUS:
		p.emitByte(OP_ADD)
	case TOKEN_MINUS:
		p.emitByte(OP_SUBTRACT)
	case TOKEN_STAR:
		p.emitByte(OP_MULTIPLY)
	case TOKEN_SLASH:
		p.emitByte(OP_DIVIDE)
	default:
		return // unreachable
	}
}

func (p *Parser) literal(canAssign bool) {
	switch parser.previous.Type {
	case TOKEN_FALSE:
		p.emitByte(OP_FALSE)
	case TOKEN_NIL:
		p.emitByte(OP_NIL)
	case TOKEN_TRUE:
		p.emitByte(OP_TRUE)
	default:
		return // unreachable
	}
}

func (p *Parser) grouping(canAssign bool) {
	p.expression()
	p.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func (p *Parser) number(canAssign bool) {
	value, err := strconv.ParseFloat(p.previous.StartAsString(p.scanner.source), 64)
	if err != nil {
		// should not happen
		panic("unable to parse float")
	}
	p.emitConstant(NumberValue(value))
}

func (p *Parser) or_(canAssign bool) {
	elseJump := p.emitJump(OP_JUMP_IF_FALSE)
	endJump := p.emitJump(OP_JUMP)

	p.patchJump(elseJump)
	p.emitByte(OP_POP)

	p.parsePrecedence(PREC_OR)
	p.patchJump(endJump)
}

func (p *Parser) string(canAssign bool) {
	prev := p.previous
	// copy only the string (avoid copying the quote characters)
	p.emitConstant(ObjVal(copyString(string(p.scanner.source[prev.Start+1 : prev.Start+prev.Length-1]))))
}

func (p *Parser) namedVariable(name Token, canAssign bool) {
	var getOp, setOp uint8
	arg := p.resolveLocal(current, &name)
	if arg != -1 {
		getOp = OP_GET_LOCAL
		setOp = OP_SET_LOCAL
	} else {
		arg = int(p.identifierConstant(&name))
		getOp = OP_GET_GLOBAL
		setOp = OP_SET_GLOBAL
	}

	if canAssign && p.match(TOKEN_EQUAL) {
		p.expression()
		p.emitBytes(setOp, uint8(arg))
	} else {
		p.emitBytes(getOp, uint8(arg))
	}
}

func (p *Parser) variable(canAssign bool) {
	p.namedVariable(p.previous, canAssign)
}

func (p *Parser) unary(canAssign bool) {
	operatorType := parser.previous.Type

	// compile the operand
	p.parsePrecedence(PREC_UNARY)

	// emit the operator instruction
	switch operatorType {
	case TOKEN_BANG:
		p.emitByte(OP_NOT)
	case TOKEN_MINUS:
		p.emitByte(OP_NEGATE)
	default:
		return // unreachable
	}
}

func (p *Parser) parsePrecedence(precedence Precedence) {
	p.advance()
	prefixRule := getRule(p.previous.Type).prefix
	if prefixRule == nil {
		p.error("Expect expression")
		return
	}

	canAssign := precedence <= PREC_ASSIGNMENT
	prefixRule(canAssign)

	for precedence <= getRule(p.current.Type).precedence {
		p.advance()
		infixRule := getRule(p.previous.Type).infix
		infixRule(canAssign)
	}

	if canAssign && p.match(TOKEN_EQUAL) {
		p.error("Invalid assignment target.")
	}
}

func (p *Parser) identifierConstant(name *Token) uint8 {
	s := string(p.scanner.source[name.Start : name.Start+name.Length])
	return p.makeConstant(ObjVal(copyString(s)))
}

func (p *Parser) identifiersEqual(a, b *Token) bool {
	if a.Length != b.Length {
		return false
	}
	return a.StartAsString(p.scanner.source) == b.StartAsString(p.scanner.source)
}

func (p *Parser) resolveLocal(compiler *Compiler, name *Token) int {
	for i := compiler.localCount - 1; i >= 0; i-- {
		local := &compiler.locals[i]
		if p.identifiersEqual(name, &local.name) {
			if local.depth == -1 {
				p.error("Can't read local variable in its own initializer.")
			}
			return i
		}
	}
	return -1
}

func (p *Parser) addLocal(name Token) {
	if current.localCount == UINT8_COUNT {
		p.error("Too many local variables in function.")
		return
	}
	local := &current.locals[current.localCount]
	current.localCount++
	local.name = name
	local.depth = -1
}

func (p *Parser) declareVariable() {
	if current.scopeDepth == 0 {
		return
	}
	name := &p.previous
	for i := current.localCount - 1; i >= 0; i-- {
		local := &current.locals[i]
		if local.depth != -1 && local.depth < current.scopeDepth {
			break
		}
		if p.identifiersEqual(name, &local.name) {
			p.error("Already a variable with this name in this scope.")
		}
	}
	p.addLocal(*name)
}

func (p *Parser) parseVariable(errorMessage string) uint8 {
	p.consume(TOKEN_IDENTIFIER, errorMessage)

	p.declareVariable()
	if current.inLocalScope() {
		return 0
	}

	return p.identifierConstant(&p.previous)
}

func (p *Parser) markInitialized() {
	current.locals[current.localCount-1].depth = current.scopeDepth
}

func (p *Parser) defineVariable(global uint8) {
	if current.inLocalScope() {
		p.markInitialized()
		return
	}
	p.emitBytes(OP_DEFINE_GLOBAL, global)
}

func (p *Parser) and_(canAssign bool) {
	endJump := p.emitJump(OP_JUMP_IF_FALSE)

	p.emitByte(OP_POP)
	p.parsePrecedence(PREC_AND)

	p.patchJump(endJump)
}

func getRule(typ TokenType) ParseRule {
	return rules[typ]
}

func (p *Parser) consume(typ TokenType, message string) {
	if p.current.Type == typ {
		p.advance()
		return
	}
	p.errorAtCurrent(message)
}

func (p *Parser) match(typ TokenType) bool {
	if !p.check(typ) {
		return false
	}
	p.advance()
	return true
}

func (p *Parser) check(typ TokenType) bool {
	return p.current.Type == typ
}

func (p *Parser) emitByte(byt uint8) {
	chunk := currentChunk()
	chunk.Write(byt, p.previous.Line)
}

func (p *Parser) emitBytes(byte1, byte2 uint8) {
	p.emitByte(byte1)
	p.emitByte(byte2)
}

func (p *Parser) emitLoop(loopStart uint8) {
	p.emitByte(OP_LOOP)

	offset := currentChunk().Count() - int(loopStart) + 2
	if offset > math.MaxUint16 {
		p.error("Loop body too large.")
	}
	high := uint8((offset >> 8) & 0xff)
	low := uint8(offset & 0xff)
	p.emitByte(high)
	p.emitByte(low)
}

func (p *Parser) emitJump(instruction uint8) int {
	p.emitByte(instruction)
	p.emitByte(0xff)
	p.emitByte(0xff)
	return currentChunk().Count() - 2
}

func (p *Parser) emitReturn() {
	p.emitByte(OP_RETURN)
}

func (p *Parser) emitConstant(value Value) {
	p.emitBytes(OP_CONSTANT, p.makeConstant(value))
}

func (p *Parser) patchJump(offset int) {
	var jump int = currentChunk().Count() - offset - 2

	if jump > math.MaxUint16 {
		p.error("Too much code to jump over.")
	}

	high := (jump >> 8) & 0xff
	low := jump & 0xff
	currentChunk().Code[offset] = (uint8)(high)
	currentChunk().Code[offset+1] = (uint8)(low)
}

func (p *Parser) InitCompiler(compiler *Compiler, typ FunctionType) {
	compiler.function = nil
	compiler.typ = typ
	compiler.localCount = 0
	compiler.scopeDepth = 0
	compiler.function = newFunction()
	current = compiler

	// the compiler claims stack slot 0
	local := current.locals[current.localCount]
	current.localCount++
	local.depth = 0
	local.name.Start = 0 // the c implementation (24.2.1) uses an empty string (may have an effect later)
	local.name.Length = 0
}

func (p *Parser) makeConstant(value Value) uint8 {
	constant := currentChunk().AddConstant(value)
	if constant > math.MaxUint8 {
		p.error("Too many constants in one chunk.")
		return 0
	}
	return uint8(constant)
}

func (p *Parser) errorAtCurrent(message string) {
	p.errorAt(p.current, message)
}

func (p *Parser) error(message string) {
	p.errorAt(p.previous, message)
}

func (p *Parser) errorAt(token Token, message string) {
	if p.panicMode {
		return
	}
	p.panicMode = true
	fmt.Fprintf(os.Stderr, "[line %d] Error", token.Line)

	if token.Type == TOKEN_EOF {
		fmt.Fprintf(os.Stderr, "at end")
	} else if token.Type == TOKEN_ERROR {
		// nothing
	} else {
		str := token.StartAsString(p.scanner.source)
		fmt.Fprintf(os.Stderr, " at '%s'", str)
	}
	fmt.Fprintf(os.Stderr, ": %s\n", message)
	p.hadError = true
}
