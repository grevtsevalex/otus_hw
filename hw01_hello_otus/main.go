package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	const msg string = "Hello, OTUS!"
	fmt.Println(reverse.String(msg))
}
