package main

import (
	"log"
)

func main() {

	t := NewTerminal()

	err := t.startPty()
	if err != nil {
		log.Fatal(err)
	}

	ui := setupUI(t)
	go t.start()

	ui.ShowAndRun()
}
