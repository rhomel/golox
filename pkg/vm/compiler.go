package vm

import "fmt"

func compile(source string) {
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
