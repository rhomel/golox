package vm

import "fmt"

type TokenType int

const (
	// Single-character tokens.
	TOKEN_LEFT_PAREN TokenType = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR
	// One or two character tokens.
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL
	// Literals.
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER
	// Keywords.
	TOKEN_AND
	TOKEN_CLASS
	TOKEN_ELSE
	TOKEN_FALSE
	TOKEN_FOR
	TOKEN_FUN
	TOKEN_IF
	TOKEN_NIL
	TOKEN_OR
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_SUPER
	TOKEN_THIS
	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_WHILE
	TOKEN_ERROR
	TOKEN_EOF
)

type Token struct {
	Type   TokenType
	Start  int
	Length int
	Line   int

	// only applies to TOKEN_ERROR
	Error string
}

func (t Token) StartAsString(source []rune) string {
	return string(source[t.Start : t.Start+t.Length])
}

func (t Token) String() string {
	typ := ""
	switch t.Type {
	case TOKEN_LEFT_PAREN:
		typ = "TOKEN_LEFT_PAREN"
	case TOKEN_RIGHT_PAREN:
		typ = "TOKEN_RIGHT_PAREN"
	case TOKEN_LEFT_BRACE:
		typ = "TOKEN_LEFT_BRACE"
	case TOKEN_RIGHT_BRACE:
		typ = "TOKEN_RIGHT_BRACE"
	case TOKEN_COMMA:
		typ = "TOKEN_COMMA"
	case TOKEN_DOT:
		typ = "TOKEN_DOT"
	case TOKEN_MINUS:
		typ = "TOKEN_MINUS"
	case TOKEN_PLUS:
		typ = "TOKEN_PLUS"
	case TOKEN_SEMICOLON:
		typ = "TOKEN_SEMICOLON"
	case TOKEN_SLASH:
		typ = "TOKEN_SLASH"
	case TOKEN_STAR:
		typ = "TOKEN_STAR"
	case TOKEN_BANG:
		typ = "TOKEN_BANG"
	case TOKEN_BANG_EQUAL:
		typ = "TOKEN_BANG_EQUAL"
	case TOKEN_EQUAL:
		typ = "TOKEN_EQUAL"
	case TOKEN_EQUAL_EQUAL:
		typ = "TOKEN_EQUAL_EQUAL"
	case TOKEN_GREATER:
		typ = "TOKEN_GREATER"
	case TOKEN_GREATER_EQUAL:
		typ = "TOKEN_GREATER_EQUAL"
	case TOKEN_LESS:
		typ = "TOKEN_LESS"
	case TOKEN_LESS_EQUAL:
		typ = "TOKEN_LESS_EQUAL"
	case TOKEN_IDENTIFIER:
		typ = "TOKEN_IDENTIFIER"
	case TOKEN_STRING:
		typ = "TOKEN_STRING"
	case TOKEN_NUMBER:
		typ = "TOKEN_NUMBER"
	case TOKEN_AND:
		typ = "TOKEN_AND"
	case TOKEN_CLASS:
		typ = "TOKEN_CLASS"
	case TOKEN_ELSE:
		typ = "TOKEN_ELSE"
	case TOKEN_FALSE:
		typ = "TOKEN_FALSE"
	case TOKEN_FOR:
		typ = "TOKEN_FOR"
	case TOKEN_FUN:
		typ = "TOKEN_FUN"
	case TOKEN_IF:
		typ = "TOKEN_IF"
	case TOKEN_NIL:
		typ = "TOKEN_NIL"
	case TOKEN_OR:
		typ = "TOKEN_OR"
	case TOKEN_PRINT:
		typ = "TOKEN_PRINT"
	case TOKEN_RETURN:
		typ = "TOKEN_RETURN"
	case TOKEN_SUPER:
		typ = "TOKEN_SUPER"
	case TOKEN_THIS:
		typ = "TOKEN_THIS"
	case TOKEN_TRUE:
		typ = "TOKEN_TRUE"
	case TOKEN_VAR:
		typ = "TOKEN_VAR"
	case TOKEN_WHILE:
		typ = "TOKEN_WHILE"
	case TOKEN_ERROR:
		typ = "TOKEN_ERROR"
	case TOKEN_EOF:
		typ = "TOKEN_EOF"
	default:
		typ = "TOKEN_UNKONWN"
	}
	return fmt.Sprintf("%s %d %d %d", typ, t.Start, t.Length, t.Line)
}

type Scanner struct {
	start   int
	current int
	line    int

	source []rune
}

func InitScanner(source string) *Scanner {
	return &Scanner{
		start:   0,
		current: 0,
		line:    1,

		source: []rune(source),
	}
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'z') ||
		c == '_'
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) ScanToken() Token {
	s.skipWhitespace()
	s.start = s.current

	if s.isAtEnd() {
		return s.makeToken(TOKEN_EOF)
	}

	c := s.advance()

	if s.isAlpha(c) {
		return s.identifier()
	}
	if s.isDigit(c) {
		return s.number()
	}
	switch c {
	case '(':
		return s.makeToken(TOKEN_LEFT_PAREN)
	case ')':
		return s.makeToken(TOKEN_RIGHT_PAREN)
	case '{':
		return s.makeToken(TOKEN_LEFT_BRACE)
	case '}':
		return s.makeToken(TOKEN_RIGHT_BRACE)
	case ';':
		return s.makeToken(TOKEN_SEMICOLON)
	case ',':
		return s.makeToken(TOKEN_COMMA)
	case '.':
		return s.makeToken(TOKEN_DOT)
	case '-':
		return s.makeToken(TOKEN_MINUS)
	case '+':
		return s.makeToken(TOKEN_PLUS)
	case '/':
		return s.makeToken(TOKEN_SLASH)
	case '*':
		return s.makeToken(TOKEN_STAR)
	case '!':
		if s.match('=') {
			return s.makeToken(TOKEN_BANG_EQUAL)
		} else {
			return s.makeToken(TOKEN_BANG)
		}
	case '=':
		if s.match('=') {
			return s.makeToken(TOKEN_EQUAL_EQUAL)
		} else {
			return s.makeToken(TOKEN_EQUAL)
		}
	case '<':
		if s.match('=') {
			return s.makeToken(TOKEN_LESS_EQUAL)
		} else {
			return s.makeToken(TOKEN_LESS)
		}
	case '>':
		if s.match('=') {
			return s.makeToken(TOKEN_GREATER_EQUAL)
		} else {
			return s.makeToken(TOKEN_GREATER)
		}
	case '"':
		return s.stringToken()
	}

	return s.errorToken("Unexpected character.")
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	r := s.source[s.current]
	s.current++
	return r
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.isAtEnd() || s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) makeToken(Type TokenType) Token {
	return Token{
		Type:   Type,
		Start:  s.start,
		Length: s.current - s.start,
		Line:   s.line,
	}
}

func (s *Scanner) errorToken(message string) Token {
	return Token{
		Type:  TOKEN_ERROR,
		Line:  s.line,
		Error: message,
	}
}

func (s *Scanner) skipWhitespace() {
	for {
		c := s.peek()
		switch c {
		case ' ':
			s.advance()
			break
		case '\r':
			s.advance()
			break
		case '\t':
			s.advance()
			break
		case '\n':
			s.line++
			s.advance()
			break
		case '/':
			if s.peekNext() == '/' {
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				return
			}
			break
		default:
			return
		}
	}
}

func (s *Scanner) checkKeyword(start int, length int, rest string, Type TokenType) TokenType {
	rrest := []rune(rest)
	i := s.start + start
	j := 0
	for i < len(s.source) && j < length {
		if s.source[i] != rrest[j] {
			return TOKEN_IDENTIFIER
		}
		i++
		j++
	}
	if j != length {
		return TOKEN_IDENTIFIER
	}
	return Type
}

func (s *Scanner) identifierType() TokenType {
	switch s.source[s.start] {
	case 'a':
		return s.checkKeyword(1, 2, "nd", TOKEN_AND)
	case 'c':
		return s.checkKeyword(1, 4, "lass", TOKEN_CLASS)
	case 'e':
		return s.checkKeyword(1, 3, "lse", TOKEN_ELSE)
	case 'f':
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword(2, 3, "lse", TOKEN_FALSE)
			case 'o':
				return s.checkKeyword(2, 1, "r", TOKEN_FOR)
			case 'u':
				return s.checkKeyword(2, 1, "r", TOKEN_FUN)
			}
		}
	case 'i':
		return s.checkKeyword(1, 1, "f", TOKEN_IF)
	case 'n':
		return s.checkKeyword(1, 2, "il", TOKEN_NIL)
	case 'o':
		return s.checkKeyword(1, 1, "r", TOKEN_OR)
	case 'p':
		return s.checkKeyword(1, 4, "rint", TOKEN_PRINT)
	case 'r':
		return s.checkKeyword(1, 5, "eturn", TOKEN_RETURN)
	case 's':
		return s.checkKeyword(1, 4, "uper", TOKEN_SUPER)
	case 't':
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'h':
				return s.checkKeyword(2, 2, "is", TOKEN_THIS)
			case 'r':
				return s.checkKeyword(2, 2, "ue", TOKEN_TRUE)
			}
		}
	case 'v':
		return s.checkKeyword(1, 2, "ar", TOKEN_VAR)
	case 'w':
		return s.checkKeyword(1, 4, "hile", TOKEN_WHILE)
	}
	return TOKEN_IDENTIFIER
}

func (s *Scanner) identifier() Token {
	for s.isAlpha(s.peek()) || s.isDigit(s.peek()) {
		s.advance()
	}
	return s.makeToken(s.identifierType())
}

func (s *Scanner) number() Token {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance() // consume '.'
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	return s.makeToken(TOKEN_NUMBER)
}

func (s *Scanner) stringToken() Token {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return s.errorToken("Unterminated string.")
	}
	s.advance() // closing quote
	return s.makeToken(TOKEN_STRING)
}
