package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rhomel/golox/pkg/util/exit"
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
		"Call":     "Callee Expr, Paren scanner.Token, Arguments []Expr",
		"Get":      "Object Expr, Name scanner.Token",
		"Grouping": "Expression Expr",
		"Literal":  "Value interface{}",
		"Logical":  "Left Expr, Operator scanner.Token, Right Expr",
		"Set":      "Object Expr, Name scanner.Token, Value Expr",
		"Super":    "Keyword scanner.Token, Method scanner.Token",
		"This":     "Keyword scanner.Token",
		"Unary":    "Operator scanner.Token, Right Expr",
		"Variable": "Name scanner.Token",
		"Assign":   "Name scanner.Token, Value Expr",
	}, "import \"github.com/rhomel/golox/pkg/scanner\"")
	defineAST(outputDirectory, "Stmt", map[string]string{
		"Block":      "Statements []Stmt",
		"Class":      "Name scanner.Token, Superclass *Variable, Methods []*Function",
		"Expression": "Expression Expr",
		"Function":   "Name scanner.Token, Params []scanner.Token, Body []Stmt",
		"IfStmt":     "Condition Expr, ThenBranch Stmt, ElseBranch Stmt",
		"Print":      "Expression Expr",
		"ReturnStmt": "Keyword scanner.Token, Value Expr",
		"VarStmt":    "Name scanner.Token, Initializer Expr",
		"While":      "Condition Expr, Body Stmt",
	}, "import \"github.com/rhomel/golox/pkg/scanner\"")
}

func defineAST(outputDirectory, baseName string, types map[string]string, imports string) {
	path := filepath.Join(outputDirectory, baseName+".go")
	file, err := os.Create(path)
	if err != nil {
		exit.Exitf(exit.ExitIOError, fmt.Sprintf("error creating file '%s': %v", path, err))
	}
	func() {
		defer file.Close()
		if err := writeFile(file, baseName, types, imports); err != nil {
			exit.Exitf(exit.ExitIOError, fmt.Sprintf("error writing to file '%s': %v", path, err))
		}
	}()
	gofmt(path)
}

func writeFile(file *os.File, baseName string, types map[string]string, imports string) error {
	header := strings.ReplaceAll(templateHeader, "<baseName>", baseName)
	header = strings.ReplaceAll(header, "<imports>", imports)
	if _, err := fmt.Fprint(file, header); err != nil {
		return err
	}
	keys := make([]string, 0, len(types))
	for k := range types {
		keys = append(keys, k)
	}
	sort.Strings(keys) // iterating Go maps is not stable so iterate by sorted key
	for _, key := range keys {
		name := key
		fields := types[key]
		generated := template
		generated = strings.ReplaceAll(generated, "<baseName>", baseName)
		generated = strings.ReplaceAll(generated, "<Name>", name)
		generated = strings.ReplaceAll(generated, "<lowercaseName>", strings.ToLower(name))
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

// GENERATED CODE from cmd/tool/gen/ast/ast-gen.go

<imports>

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

type <Name>StringVisitor interface {
	Visit<Name><baseName>String(*<Name>) string
}

func (<lowercaseName> *<Name>) AcceptString(visitor <Name>StringVisitor) string {
	return visitor.Visit<Name><baseName>String(<lowercaseName>)
}

type <Name>VoidVisitor interface {
	Visit<Name><baseName>Void(*<Name>)
}

func (<lowercaseName> *<Name>) AcceptVoid(visitor <Name>VoidVisitor) {
	visitor.Visit<Name><baseName>Void(<lowercaseName>)
}

type <Name>Visitor interface {
	Visit<Name><baseName>(*<Name>) interface{}
}

func (<lowercaseName> *<Name>) Accept(visitor <Name>Visitor) interface{} {
	return visitor.Visit<Name><baseName>(<lowercaseName>)
}
`

func gofmt(path string) {
	if err := exec.Command("go", "fmt", path).Run(); err != nil {
		exit.Exitf(exit.ExitIOError, "failed to run 'go fmt %s': %v", path, err)
	}
}
