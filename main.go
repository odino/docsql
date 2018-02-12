package main

import (
	"fmt"
	"github.com/odino/docsql/cmd"
	"os"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Unrecoverable error:", r)
			os.Exit(99) // I got 99 problems...
		}
	}()
	cmd.Execute()
}
