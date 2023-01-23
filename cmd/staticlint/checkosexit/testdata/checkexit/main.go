package main

import (
	"fmt"
	"os"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
	os.Exit(1) // want "func os.Exit in main"
}
