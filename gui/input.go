package gui

import (
	"github.com/faiface/pixel/pixelgl"
)

func (g *GUI) handleInput() {
	scroll := g.window.MouseScroll()

	switch {
	case scroll.Y != 0:
		switch {
		case scroll.Y < 0:
			g.terminal.ScrollDown()
		case scroll.Y > 0:
			g.terminal.ScrollUp()
		}

	case g.window.Pressed(pixelgl.KeyLeftControl):
		switch {
		case g.window.JustPressed(pixelgl.KeyC):
			g.terminal.Write([]byte{3})
		}

	case g.window.JustPressed(pixelgl.KeyEnter):
		g.terminal.Write([]byte{'\n'})

	case g.window.JustReleased(pixelgl.KeyTab):
		g.terminal.Write([]byte{'\t'})

	case g.window.JustPressed(pixelgl.KeyBackspace):
		g.terminal.Write([]byte{8})

	default:
		g.terminal.Write([]byte(g.window.Typed()))

	}
}
