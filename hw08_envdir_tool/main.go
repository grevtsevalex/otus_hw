package main

import (
	"fmt"
	"os"
)

func main() {
	envFile := os.Args[1]
	envs, err := ReadDir(envFile)
	if err != nil {
		fmt.Println(err)
	}

	RunCmd(os.Args[2:], envs)
}
