package vm

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

var debugPrintCode = false

var compilingChunk *Chunk

func currentChunk() *Chunk {
	return compilingChunk
}

func compile(source string, chunk *Chunk) bool {
	scanner := InitScanner(source)
	parser.scanner = scanner
	parser.hadError = false
	parser.panicMode = false
	compilingChunk = chunk

	parser.advance()
	parser.expression()
	parser.consume(TOKEN_EOF, "Expect end of expression.")

	endCompiler()
	return !parser.hadError
}

func endCompiler() {
	parser.emitReturn()
	if debugPrintCode {
		if !parser.hadError {
			currentChunk().Disassemble("code")
		}
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

type ParseFn func()

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
	rules[TOKEN_BANG] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_BANG_EQUAL] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_EQUAL] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_EQUAL_EQUAL] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_GREATER] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_GREATER_EQUAL] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_LESS] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_LESS_EQUAL] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_IDENTIFIER] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_STRING] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_NUMBER] = ParseRule{parser.number, nil, PREC_NONE}
	rules[TOKEN_AND] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_CLASS] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_ELSE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FALSE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FOR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FUN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_IF] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_NIL] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_OR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_PRINT] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_RETURN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_SUPER] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_THIS] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_TRUE] = ParseRule{nil, nil, PREC_NONE}
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

func (p *Parser) binary() {
	operatorType := p.previous.Type
	rule := getRule(operatorType)
	p.parsePrecedence(Precedence(rule.precedence + 1))
	switch operatorType {
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

func (p *Parser) grouping() {
	p.expression()
	p.consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func (p *Parser) number() {
	value, err := strconv.ParseFloat(p.previous.StartAsString(p.scanner.source), 64)
	if err != nil {
		// should not happen
		panic("unable to parse float")
	}
	p.emitConstant(value)
}

func (p *Parser) unary() {
	operatorType := parser.previous.Type

	// compile the operand
	p.parsePrecedence(PREC_UNARY)

	// emit the operator instruction
	switch operatorType {
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
	prefixRule()

	for precedence <= getRule(p.current.Type).precedence {
		p.advance()
		infixRule := getRule(p.previous.Type).infix
		infixRule()
	}
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

func (p *Parser) emitByte(byt uint8) {
	chunk := currentChunk()
	chunk.Write(byt, p.previous.Line)
}

func (p *Parser) emitBytes(byte1, byte2 uint8) {
	p.emitByte(byte1)
	p.emitByte(byte2)
}

func (p *Parser) emitReturn() {
	p.emitByte(OP_RETURN)
}

func (p *Parser) emitConstant(value float64) {
	p.emitBytes(OP_CONSTANT, p.makeConstant(value))
}

func (p *Parser) makeConstant(value float64) uint8 {
	constant := currentChunk().AddConstant(Value(value))
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
