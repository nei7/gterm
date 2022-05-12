package renderer

import (
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/nei7/gterm/config"

	"golang.org/x/image/colornames"
)

func Run() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:     "gterm",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	r := NewRenderer(config, win)

	go func() {
		r.terminal.Run(r.updateChan)
	}()

	r.Draw(win)
}

func (g *Renderer) Draw(win *pixelgl.Window) {
	for !win.Closed() {

		select {
		case <-g.updateChan:
			win.Clear(colornames.Black)
			g.drawText()
		default:
			g.handleInput()
			g.handleScroll()
		}

		win.Update()
	}

}
