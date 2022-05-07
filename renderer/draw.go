package renderer

func (g *Renderer) drawText() {
	g.text.Draw(g.window, g.terminal.GetLines())
}
