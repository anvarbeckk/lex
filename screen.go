package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

func initScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing screen: %s\n", err)
		os.Exit(1)
	}

	err = screen.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing screen: %s\n", err)
		os.Exit(1)
	}

	return screen
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
