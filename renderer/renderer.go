package renderer

import (
	"log"

	"github.com/faiface/pixel/pixelgl"
	"github.com/nei7/gterm/config"
	"github.com/nei7/gterm/font"
	"github.com/nei7/gterm/renderer/text"
	"github.com/nei7/gterm/term"
)

type Renderer struct {
	window      *pixelgl.Window
	config      *config.Config
	terminal    *term.Terminal
	buffer      *term.Buffer
	fontManager *font.Manager
	text        *text.Text
	updateChan  chan struct{}
	pixelHeight float64
	pixelWidth  float64
}

func NewRenderer(config *config.Config, win *pixelgl.Window) *Renderer {
	fontManager := font.NewManager()
	err := fontManager.SetFont("SauceCodePro Nerd Font Mono")
	if err != nil {
		log.Fatal(err)
	}
	buffer := term.NewBuffer()

	w, h := win.Bounds().Size().XY()
	charW, charH := fontManager.CharSize().XY()

	t := term.New(buffer)
	t.SetSize(int(h/charH), int(w/charW))

	r := &Renderer{
		config:      config,
		buffer:      buffer,
		terminal:    t,
		updateChan:  make(chan struct{}),
		window:      win,
		fontManager: fontManager,
		text:        text.NewText(fontManager),
		pixelHeight: h - charH,
		pixelWidth:  w,
	}

	return r
}
