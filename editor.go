package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"io/ioutil"
	"os"
	"strings"
)

type Editor struct {
	content  []string
	cursorX  int
	cursorY  int
	screen   tcell.Screen
	filename string
}

func NewEditor(filename string) *Editor {
	return &Editor{content: []string{""}, cursorX: 0, cursorY: 0, filename: filename}
}

func (e *Editor) InitScreen() error {
	var err error
	e.screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}

	err = e.screen.Init()
	if err != nil {
		return err
	}

	return nil
}

func (e *Editor) FiniScreen() {
	e.screen.Fini()
}

func (e *Editor) Run() {
	if e.filename != "" {
		err := e.LoadFile(e.filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading file: %s\n", err)
		}
	}

	for {
		e.drawScreen()

		ev := e.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			e.handleKey(ev)
		}
	}
}

func (e *Editor) LoadFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	e.content = strings.Split(string(content), "\n")
	return nil
}

func (e *Editor) SaveToFile(filename string) error {
	content := strings.Join(e.content, "\n")
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func (e *Editor) saveAndExit() {
	if e.filename != "" {
		if e.confirmSave() {
			err := e.SaveToFile(e.filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving file: %s\n", err)
			} else {
				fmt.Println("File saved successfully.")
			}
		}
	}

	e.FiniScreen()
	fmt.Println("Exiting the text editor.")
	os.Exit(0)
}

func (e *Editor) drawScreen() {
	e.screen.Clear()

	for y, line := range e.content {
		for x, char := range line {
			style := tcell.StyleDefault

			if y == e.cursorY && x == e.cursorX {
				style = style.Reverse(true)
			}

			e.screen.SetContent(x, y, char, nil, style)
		}
	}

	e.screen.Show()
}

func (e *Editor) handleKey(event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyCtrlL:
		e.drawScreen()
	case tcell.KeyCtrlQ:
		e.FiniScreen()
		e.saveAndExit()
	case tcell.KeyCtrlS:
		e.saveToFile()
	case tcell.KeyEnter:
		e.handleEnter()
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		e.handleBackspace()
	default:
		char := event.Rune()
		e.handleChar(char)
	}
}

func (e *Editor) saveToFile() {
	if e.filename != "" {
		err := e.SaveToFile(e.filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving file: %s\n", err)
		} else {
			fmt.Println("File saved successfully.")
		}
	} else {
		fmt.Println("Please provide a filename to save the content.")
	}
}

func (e *Editor) handleEnter() {
	e.content = append(e.content[:e.cursorY+1], "")
	copy(e.content[e.cursorY+2:], e.content[e.cursorY+1:])
	e.cursorY++
	e.cursorX = 0
}

func (e *Editor) handleBackspace() {
	if e.cursorX > 0 {
		e.content[e.cursorY] = e.content[e.cursorY][:e.cursorX-1] + e.content[e.cursorY][e.cursorX:]
		e.cursorX--
	} else if e.cursorY > 0 {
		prevLineLen := len(e.content[e.cursorY-1])
		e.content[e.cursorY-1] += e.content[e.cursorY]
		copy(e.content[e.cursorY:], e.content[e.cursorY+1:])
		e.content = e.content[:len(e.content)-1]
		e.cursorY--
		e.cursorX = prevLineLen
	}
}

func (e *Editor) handleChar(char rune) {
	e.content[e.cursorY] = e.content[e.cursorY][:e.cursorX] + string(char) + e.content[e.cursorY][e.cursorX:]
	e.cursorX++
}

func (e *Editor) confirmSave() bool {
	if len(e.content) == 1 && e.content[0] == "" {
		// No changes made, no need to confirm save
		return true
	}

	saveConfirmation := "Save modified content (y/n)? "
	e.drawMessage(saveConfirmation)

	for {
		ev := e.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'y', 'Y':
					return true
				case 'n', 'N':
					return false
				}
			}
		}
	}
}

func (e *Editor) drawMessage(message string) {
	e.screen.Clear()
	for y, char := range message {
		e.screen.SetContent(y, 0, char, nil, tcell.StyleDefault)
	}
	e.screen.Show()
}

func (e *Editor) printMessageAndExit(message string) {
	e.screen.Fini()
	println(message)
}
