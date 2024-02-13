// main.go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./lex filename")
		os.Exit(1)
	}

	filename := os.Args[1]
	editor := NewEditor(filename)

	if err := editor.InitScreen(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize screen: %s\n", err)
		os.Exit(1)
	}

	defer editor.FiniScreen()

	editor.Run()
}
