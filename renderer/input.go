package renderer

import (
	"github.com/faiface/pixel/pixelgl"
)

func (g *Renderer) handleInput() error {

	switch {
	case g.window.Pressed(pixelgl.KeyLeftControl):
		switch {
		case g.window.JustPressed(pixelgl.KeyC):
			return g.terminal.Write([]byte{3})
		}
	case g.window.JustPressed(pixelgl.KeyEnter):
		return g.terminal.Write([]byte{'\n'})

	case g.window.JustReleased(pixelgl.KeyTab):
		return g.terminal.Write([]byte{'\t'})

	case g.window.JustPressed(pixelgl.KeyBackspace):
		return g.terminal.Write([]byte{8})

	case g.window.JustPressed(pixelgl.KeyEscape):
		return g.terminal.Write([]byte{27})

	default:
		return g.terminal.Write([]byte(g.window.Typed()))
	}

	return nil
}

func (g *Renderer) handleScroll() {

	scroll := g.window.MouseScroll()

	if scroll.Y < 0 {
		g.buffer.ScrollDown()
	} else {
		g.buffer.ScrollUp()
	}

}
