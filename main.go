package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("gterm")

	t := NewTerminal()

	w.SetContent(
		container.NewMax(t.content, container.NewWithoutLayout(t.cursor)),
	)

	err := t.setPty()
	if err != nil {
		log.Fatal(err)
	}
	go t.start()

	w.ShowAndRun()
}
