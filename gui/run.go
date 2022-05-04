package gui

import (
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func (g *GUI) Run() {

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:     "gterm",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	g.setupWindow(win)

	go g.terminal.Run()

	fps := time.Tick(time.Second / 60)
	for !win.Closed() {
		win.Clear(colornames.Black)

		g.handleInput()

		g.drawText()

		win.Update()

		time.Sleep(time.Millisecond * 50)

		<-fps
	}
}
