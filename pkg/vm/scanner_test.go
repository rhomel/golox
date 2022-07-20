package vm

import "testing"

func Test_checkKeyword(t *testing.T) {
	matchAnd := func(source string, start int) {
		t.Helper()
		scanner := InitScanner(source)
		scanner.start = start
		token := scanner.checkKeyword(1, 2, "nd", TOKEN_AND)
		if want, got := TOKEN_AND, token; want != got {
			t.Errorf("want token %v, got: %v", want, got)
		}
	}

	matchAnd("and nd and", 0)
	matchAnd("and nd and", 3)
	matchAnd("and nd and", 7)
}

func Test_identifierType(t *testing.T) {
	matchKeyword := func(source string, start int, token TokenType) {
		t.Helper()
		scanner := InitScanner(source)
		scanner.start = start
		received := scanner.identifierType()
		if want, got := token, received; want != got {
			t.Errorf("want token %v, got: %v", want, got)
		}
	}

	matchKeyword("and", 0, TOKEN_AND)
	matchKeyword("and class while if else", 4, TOKEN_CLASS)
	matchKeyword("and class while if else", 10, TOKEN_WHILE)
}
