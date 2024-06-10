package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	for arg := range args {
		fmt.Println(arg)
	}
}
