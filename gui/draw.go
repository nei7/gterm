package gui

import (
	"github.com/faiface/pixel"
)

func (g *GUI) drawText() {
	lines := g.terminal.GetLines()
	g.text.Draw(g.window, pixel.IM, lines)
}
