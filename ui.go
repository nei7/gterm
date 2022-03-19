package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func setupUI(t *Terminal) fyne.Window {
	a := app.New()
	w := a.NewWindow("gterm")

	w.SetContent(
		container.NewMax(t.content, container.NewWithoutLayout(t.cursor)),
	)

	w.Canvas().SetOnTypedRune(t.onTypedRune)
	w.Canvas().SetOnTypedKey(t.onTypedKey)

	return w
}
