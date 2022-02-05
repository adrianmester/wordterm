package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"strconv"
)

func main() {

	var dayNumber int64 = -1
	var err error
	if len(os.Args) == 2 {
		dayNumber, err = strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			panic(err)
		}
	}

	wl, err := NewWordList()
	if err != nil {
		panic(err)
	}
	p := tea.NewProgram(initialModel(&wl, dayNumber), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
