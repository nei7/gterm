package renderer

import (
	"github.com/faiface/pixel"
)

func (g *Renderer) drawText() {
	g.text.Draw(g.window, pixel.IM, g.terminal.GetLines())
}
