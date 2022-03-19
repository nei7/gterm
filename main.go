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
	go t.start()

	ui := setupUI(t)

	ui.ShowAndRun()
}
