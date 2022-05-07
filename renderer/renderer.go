package renderer

import (
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/nei7/gterm/config"
	"github.com/nei7/gterm/font"
	"github.com/nei7/gterm/renderer/text"
	"github.com/nei7/gterm/term"
)

type Renderer struct {
	window *pixelgl.Window
	config *config.Config

	terminal *term.Terminal
	buffer   *term.Buffer

	fontManager *font.Manager
	text        *text.Text

	windowSize pixel.Vec
}

func NewRenderer(config *config.Config) *Renderer {
	r := &Renderer{}
	r.config = config
	r.buffer = term.NewBuffer()
	r.terminal = term.New(r.buffer)
	return r
}

func (r *Renderer) setupWindow(w *pixelgl.Window) {
	r.window = w
	r.windowSize = w.Bounds().Size()

	fontManager := font.NewManager()
	err := fontManager.SetFont("SauceCodePro Nerd Font Mono")
	if err != nil {
		log.Fatal(err)
	}
	r.fontManager = fontManager

	txtPos := pixel.V(0, r.windowSize.Y-float64(fontManager.CharSize().Y))
	r.text = text.NewText(txtPos, fontManager)

	r.ResizeTerminal()
}

func (g *Renderer) ResizeTerminal() {
	windowSize := g.window.Bounds().Size()

	cols := int(g.windowSize.X / float64(g.fontManager.CharSize().X))
	rows := int(g.windowSize.Y / float64(g.fontManager.CharSize().Y))

	g.terminal.SetSize(uint16(rows), uint16(cols))
	g.terminal.SetPtySize(windowSize)
}
