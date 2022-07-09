package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"rhomel.com/crafting-interpreters-go/pkg/util/exit"
)

// ast.go code generator
//
// Expected usage:
//   go run cmd/tool/gen/ast/ast-gen.go pkg/ast/gen
//
// Note: although the files are placed in a directory called 'gen' the package
// of all generated files is 'ast'.
//
// reference:
//   Metaprogramming the trees 5.2.2
//   https://craftinginterpreters.com/representing-code.html#metaprogramming-the-trees

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		exit.Exitf(exit.ExitCodeUsageError, "Usage: ast-gen <output directory>")
	}
	outputDirectory := args[0]
	defineAST(outputDirectory, "Expr", map[string]string{
		"Binary":   "Left Expr, Operator scanner.Token, Right Expr",
		"Grouping": "Expression Expr",
		"Literal":  "Value interface{}",
		"Unary":    "operator scanner.Token, Right Expr",
	})
}

func defineAST(outputDirectory, baseName string, types map[string]string) {
	path := filepath.Join(outputDirectory, baseName+".go")
	file, err := os.Create(path)
	if err != nil {
		exit.Exitf(exit.ExitIOError, fmt.Sprintf("error creating file '%s': %v", path, err))
	}
	func() {
		defer file.Close()
		if err := writeFile(file, baseName, types); err != nil {
			exit.Exitf(exit.ExitIOError, fmt.Sprintf("error writing to file '%s': %v", path, err))
		}
	}()
	gofmt(path)
}

func writeFile(file *os.File, baseName string, types map[string]string) error {
	header := strings.ReplaceAll(templateHeader, "<baseName>", baseName)
	if _, err := fmt.Fprint(file, header); err != nil {
		return err
	}
	for name, fields := range types {
		generated := template
		generated = strings.ReplaceAll(generated, "<baseName>", baseName)
		generated = strings.ReplaceAll(generated, "<Name>", name)
		generated = strings.ReplaceAll(generated, "<Fields>", defineFields(fields))
		if _, err := fmt.Fprintf(file, generated); err != nil {
			return err
		}
	}
	return nil
}

func defineFields(fieldString string) string {
	var fields string
	for i, field := range strings.Split(fieldString, ",") {
		if i != 0 {
			fields += "\n"
		}
		fields += "\t" + field
	}
	return fields
}

var templateHeader = `
package ast

import "rhomel.com/crafting-interpreters-go/pkg/scanner"

type <baseName> interface {
	is<baseName>() // private method to tag which structs are <baseName>
}
`
var template = `

var _ <baseName> = (*<Name>)(nil)

type <Name> struct {
<Fields>
}

func (*<Name>) is<baseName>() {}
`

func gofmt(path string) {
	if err := exec.Command("go", "fmt", path).Run(); err != nil {
		exit.Exitf(exit.ExitIOError, "failed to run 'go fmt %s': %v", path, err)
	}
}
