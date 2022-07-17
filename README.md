
# golox: Lox Implementation in Go

**[Test suite](https://github.com/munificent/craftinginterpreters#testing-your-implementation)
status: All 239 tests passed (556 expectations).**

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
  case will trigger and typically stop with an 'unsupported *ast node type*'
  panic.

*Since this implementation passes the test suite, I am tagging this with
'tree-walker-feature-complete' for future reference. I intend to do a
few things from here that may deviate from the book:*

- Expand base language functionality (array type, better IO support, AST printing, etc).
- Implement a faster VM version.
- Improve the REPL.

## Commands

- Start REPL: `go run cmd/golox/golox.go`
- Run file.lox: `go run cmd/golox/golox.go file.lox`
- Build golox binary: `go build cmd/golox/golox.go`

### Adding and Updating the AST types in `pkg/ast/gen`

- To Generate AST node types: `go run cmd/tool/gen/ast/ast-gen.go pkg/ast/gen`

`ast-gen.go` has the `defineAst` method that determines what types to generate.

After generation you should probably update the parser and interpreter to
utilize the new AST node types. Otherwise you may get an `unsupported
*new node type*` panic.

## REPL

`go run cmd/golox/golox.go`

Use `ctrl+d` to exit.

## Running Samples

The `samples` directory has many samples from
[Crafting Interpreters](https://craftinginterpreters.com/contents.html).

You can run them like any other lox file: `go run cmd/golox/golox.go samples/hello.lox`

## CPU Profile

You can run the interpreter with cpu profiling enabled.

Example:

```
# run `14-fib-bench.lox` file and output cpu profile to `cpuprofile`:
go run cmd/golox/golox.go samples/14-fib-bench.lox cpuprofile

# use pprof the examine `cpuprofile`
go tool pprof cpuprofile
```

## Benchmarks

As expected the tree-walk interpreter is very slow.

On a Macbook Air M1, `fib(35)` runs in:

- golox: 17 seconds
- go: 43 milliseconds

Commands:

```
go run cmd/golox/golox.go samples/14-fib-bench.lox
9227465
17
```

```
go run cmd/other/go-fib/main.go
9227465
43.213792ms
```

For comparison jlox in [Chunks of Bytecode](https://craftinginterpreters.com/chunks-of-bytecode.html)
runs `fib(40)` in 72 seconds. golox runs `fib(40)` in 197 seconds and plain Go
runs in 351ms on a Macbook Air M1. Since `fib(40)` takes too long to complete I
left the samples at `fib(35)`.

