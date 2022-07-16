package main

import (
	"fmt"
	"time"
)

// fibonacci in Go to compare against samples/14-fib-bench.lox

func main() {
	before := time.Now()
	fmt.Println(fib(35))
	after := time.Now()
	fmt.Println(after.Sub(before))
}

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
