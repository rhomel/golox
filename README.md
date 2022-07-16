
# golox: Lox Implementation in Go

This is a complete interpreter for the Lox programming language as defined in
[Crafting Interpreters](https://craftinginterpreters.com/a-tree-walk-interpreter.html).

This implementation closely follows the described Java implementation.

Notable differences:

- The Java implementation utilizes Java Exceptions. Go does not have true
  exceptions but Go panic/recover works well enough to model the same exception
  use cases.
- Java inheritance is also often used in the Java version when walking the AST
  with the visitor pattern. In this implementation I use a Go type switch to
  delegate to the proper struct type. The consequence is that after running the
  AST generator, any new AST node type must be added to all type switches. If
  you forget to do so, when the new node is encountered, the switch default
  case will trigger and typically stop with an 'unsupported <ast node type>'
  error.

## Commands

- Start REPL: `go run cmd/golox/golox.go`
- Run file.lox: `go run cmd/golox/golox.go file.lox`
- Build golox binary: `go build cmd/golox/golox.go`

### Adding and Updating the AST types in `pkg/ast/gen`

- To Generate AST node types: `go run cmd/tool/gen/ast/ast-gen.go pkg/ast/gen`

`ast-gen.go` has defineAst method calls which determine what types to generate.

After generation you should probably update the parser and interpreter to
utilize the new AST node types.

## REPL

`go run cmd/golox/golox.go`

Use `ctrl+d` to exit.

## CPU Profile

You can run the interpreter with cpu profiling enabled.

Example:

```
# run `14-fib-bench.lox` file and output cpu profile to `cpuprofile`:
go run cmd/golox/golox.go samples/14-fib-bench.lox cpuprofile

# use pprof the examine `cpuprofile`
go tool pprof cpuprofile
```

