package gui

import (
	"fmt"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/nei7/gterm/config"
	"github.com/nei7/gterm/term"
	"golang.org/x/image/font"
)

type GUI struct {
	window *pixelgl.Window

	height float64
	width  float64

	font font.Face

	text   *Text
	config *config.Config

	terminal *term.Terminal
}

func New(config *config.Config) *GUI {
	g := &GUI{}

	fontPath := fmt.Sprintf("/usr/share/fonts/TTF/%s.ttf", config.Font.Family)

	face, err := loadTTF(fontPath, config.Font.Size)
	if err != nil {
		log.Fatal(err)
	}
	g.font = face
	g.config = config
	g.terminal = term.New()

	return g
}

func (g *GUI) setupWindow(w *pixelgl.Window) {
	windowSize := w.Bounds().Size()
	g.width = windowSize.X - g.config.Window.Padding.X
	g.height = windowSize.Y - g.config.Window.Padding.Y
	g.window = w

	atlas := text.NewAtlas(g.font, text.ASCII)

	g.text = NewText(pixel.Vec{
		X: g.config.Window.Padding.X,
		Y: g.height - atlas.Ascent() - g.config.Window.Padding.Y,
	}, atlas)

	cols := int(g.width / g.text.LineHeight)
	rows := int(g.height / (g.text.LineHeight))

	g.terminal.SetSize(uint16(rows), uint16(cols))
}

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
