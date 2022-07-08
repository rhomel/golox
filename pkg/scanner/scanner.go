package scanner

import "fmt"

type TokenType string

const (
	// Single-character tokens.
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	COMMA       TokenType = "COMMA"
	DOT         TokenType = "DOT"
	MINUS       TokenType = "MINUS"
	PLUS        TokenType = "PLUS"
	SEMICOLON   TokenType = "SEMICOLON"
	SLASH       TokenType = "SLASH"
	STAR        TokenType = "STAR"

	// One or two character tokens.
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"

	// Literals.
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"

	// Keywords.
	AND    TokenType = "AND"
	CLASS  TokenType = "CLASS"
	ELSE   TokenType = "ELSE"
	FALSE  TokenType = "FALSE"
	FUN    TokenType = "FUN"
	FOR    TokenType = "FOR"
	IF     TokenType = "IF"
	NIL    TokenType = "NIL"
	OR     TokenType = "OR"
	PRINT  TokenType = "PRINT"
	RETURN TokenType = "RETURN"
	SUPER  TokenType = "SUPER"
	THIS   TokenType = "THIS"
	TRUE   TokenType = "TRUE"
	VAR    TokenType = "VAR"
	WHILE  TokenType = "WHILE"

	EOF TokenType = "EOF"
)

type Token struct {
	Typ     TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %v %v", t.Typ, t.Lexeme, t.Literal)
}

type Scanner struct {
	reporter ErrorReporter

	source []rune
	tokens []*Token

	start   int
	current int
	line    int
}

type ErrorReporter interface {
	Error(line int, message string)
}

func NewScanner(source string, reporter ErrorReporter) *Scanner {
	return &Scanner{
		reporter: reporter,
		source:   []rune(source),
		start:    0,
		current:  0,
		line:     1,
	}
}

func (s *Scanner) ScanTokens() []*Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, &Token{EOF, "", nil, s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	default:
		s.reporter.Error(s.line, fmt.Sprintf("Unexpected character '%s'.", runeToReadableString(c)))
	}
}

func runeToReadableString(r rune) string {
	switch r {
	case '\n':
		return "\\n"
	case '\r':
		return "\\r"
	case '\t':
		return "\\t"
	default:
		return string(r)
	}
}

func (s *Scanner) addToken(typ TokenType) {
	s.addTokenLiteral(typ, nil)
}

func (s *Scanner) addTokenLiteral(typ TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, &Token{typ, string(text), literal, s.line})
}

func (s *Scanner) advance() rune {
	current := s.current
	s.current++
	return s.source[current]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
