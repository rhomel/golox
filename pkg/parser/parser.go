package parser

// # recursive descent parser
//
// ## grammar reference:
//   [https://craftinginterpreters.com/parsing-expressions.html#ambiguity-and-the-parsing-game]
//   [https://craftinginterpreters.com/statements-and-state.html#assignment-syntax]
//
// expression     → assignment ;
// assignment     → IDENTIFIER "=" assignment
//                | equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary
//                | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")"
//                | IDENTIFIER ;

// ## statement rules
//   [https://craftinginterpreters.com/statements-and-state.html#statements]
//   [https://craftinginterpreters.com/statements-and-state.html#variable-syntax]
//
// program        → declaration* EOF ;
//
// declaration    → varDecl
//                | statement ;
//
// varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
//
// ## block syntax
//   [https://craftinginterpreters.com/statements-and-state.html#block-syntax-and-semantics]
// ## conditional
//   [https://craftinginterpreters.com/control-flow.html#conditional-execution]
// statement      → exprStmt
//                | ifStmt ;
//                | printStmt ;
//                | block ;
//
// ifStmt         → "if" "(" expression ")" statement
//                ( "else" statement )? ;
// block          → "{" declaration* "}" ;
// exprStmt       → expression ";" ;
// printStmt      → "print" expression ";" ;

import (
	ast "rhomel.com/crafting-interpreters-go/pkg/ast/gen"
	"rhomel.com/crafting-interpreters-go/pkg/scanner"
)

type Parser struct {
	reporter ParseErrorReporter
	tokens   []*scanner.Token
	current  int
}

type ParseErrorReporter interface {
	ParseError(token scanner.Token, message string)
}

func NewParser(tokens []*scanner.Token, reporter ParseErrorReporter) *Parser {
	parser := &Parser{
		reporter: reporter,
		tokens:   tokens,
	}
	return parser
}

func (p *Parser) Parse() (statements []ast.Stmt) {
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) assignment() ast.Expr {
	// this is no the standard recursive descent pattern:
	// ref: https://craftinginterpreters.com/statements-and-state.html#assignment-syntax
	// assume there's an equality expression for now, even if it is an
	// IDENTIFER it will be parsed as an IDENTIFIER
	expr := p.equality()

	// see if there's an EQUAL token
	if p.match(scanner.EQUAL) {
		equals := p.previous()
		// see if the matched expr was an IDENTIFIER
		variable, ok := expr.(*ast.Variable)
		if ok {
			value := p.assignment()
			return &ast.Assign{variable.Name, value}
		}
		p.err(equals, "Invalid assignment target.")
	}
	return expr
}

func (p *Parser) declaration() (stmt ast.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			p.syncrhonize()
			stmt = nil
		}
	}()
	if p.match(scanner.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() ast.Stmt {
	identifier := p.consume(scanner.IDENTIFIER, "Expect variable name.")
	var initializer ast.Expr
	if p.match(scanner.EQUAL) {
		initializer = p.expression()
	}
	p.consume(scanner.SEMICOLON, "Expect ';' after variable declaration.")
	return &ast.VarStmt{identifier, initializer}
}

func (p *Parser) statement() ast.Stmt {
	if p.match(scanner.IF) {
		return p.ifStatement()
	}
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}
	if p.match(scanner.LEFT_BRACE) {
		return &ast.Block{p.block()}
	}
	return p.expressionStatement()
}

func (p *Parser) ifStatement() ast.Stmt {
	p.consume(scanner.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(scanner.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch ast.Stmt
	// greedily match 'else' token so it is always attached to the nearest 'if'
	if p.match(scanner.ELSE) {
		elseBranch = p.statement()
	}
	return &ast.IfStmt{condition, thenBranch, elseBranch}
}

func (p *Parser) printStatement() ast.Stmt {
	expr := p.expression()
	p.consume(scanner.SEMICOLON, "Expect ';' after value.")
	return &ast.Print{expr}
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.consume(scanner.SEMICOLON, "Expect ';' after expression.")
	return &ast.Expression{expr}
}

func (p *Parser) block() []ast.Stmt {
	var statements []ast.Stmt
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(scanner.RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) equality() ast.Expr {
	var expr ast.Expr = p.comparison()
	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &ast.Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) comparison() ast.Expr {
	var expr ast.Expr = p.term()
	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &ast.Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) term() ast.Expr {
	var expr ast.Expr = p.factor()
	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &ast.Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) factor() ast.Expr {
	var expr ast.Expr = p.unary()
	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &ast.Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right := p.unary()
		return &ast.Unary{operator, right}
	}
	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.match(scanner.FALSE) {
		return &ast.Literal{false}
	}
	if p.match(scanner.TRUE) {
		return &ast.Literal{true}
	}
	if p.match(scanner.NIL) {
		return &ast.Literal{nil}
	}
	if p.match(scanner.NUMBER, scanner.STRING) {
		return &ast.Literal{p.previous().Literal}
	}
	if p.match(scanner.IDENTIFIER) {
		name := p.previous()
		return &ast.Variable{name}
	}
	if p.match(scanner.LEFT_PAREN) {
		expr := p.expression()
		p.consume(scanner.RIGHT_PAREN, "Expect ')' after expression.")
		return &ast.Grouping{expr}
	}
	panic(p.err(p.peek(), "Expect expression."))
}

func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(typ scanner.TokenType, message string) scanner.Token {
	if p.check(typ) {
		return p.advance()
	}
	panic(p.err(p.peek(), message))
}

func (p *Parser) check(typ scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Typ == typ
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Typ == scanner.EOF
}

func (p *Parser) peek() scanner.Token {
	return *p.tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return *p.tokens[p.current-1]
}

func (p *Parser) err(token scanner.Token, message string) *ParseError {
	p.reporter.ParseError(token, message)
	return &ParseError{}
}

func (p *Parser) syncrhonize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Typ == scanner.SEMICOLON {
			return
		}

		switch p.peek().Typ {
		case scanner.CLASS, scanner.FUN, scanner.VAR, scanner.FOR, scanner.IF, scanner.WHILE, scanner.PRINT, scanner.RETURN:
			return
		}

		p.advance()
	}
}

type ParseError struct{}

func (*ParseError) Error() string {
	return ""
}

var _ error = (*ParseError)(nil)
