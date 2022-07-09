package exit

import (
	"fmt"
	"os"
)

const (
	ExitCodeOK         = 0
	ExitCodeUsageError = 1
	ExitSyntaxError    = 65
	ExitIOError        = 100
)

func Exitf(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
	os.Exit(code)
}
