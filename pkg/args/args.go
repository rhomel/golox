package args

import "flag"

type Args struct {
	args []string
}

func New() *Args {
	return &Args{flag.Args()}
}

func (a *Args) Len() int {
	return len(a.args)
}

func (a *Args) Get() []string {
	return a.args
}
