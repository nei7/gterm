package gui

import "github.com/faiface/pixel"

func (g *GUI) drawText() {
	lines := g.terminal.Buffer.GetLines()
	g.text.DrawBuff(lines)
	g.text.Draw(g.window, pixel.IM)
}
