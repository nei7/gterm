package term

import (
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/nei7/gterm/term/text"
	"github.com/nei7/gterm/util"

	atlas "github.com/faiface/pixel/text"
)

type Terminal struct {
	window *pixelgl.Window
	text   *text.Text
	font   font.Face

	pty io.ReadWriter

	height float64
	width  float64

	title string

	config *Config
}

func New() *Terminal {
	t := &Terminal{}

	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fontPath := fmt.Sprintf("/usr/share/fonts/TTF/%s.ttf", config.Font.Family)

	face, err := loadTTF(fontPath, config.Font.Size)
	if err != nil {
		log.Fatal(err)
	}

	pty, err := startPty(util.GetHomeDir())
	if err != nil {
		log.Fatal(err)
	}

	t.pty = pty
	t.config = config
	t.font = face

	return t
}

func (t *Terminal) setupWindow(w *pixelgl.Window) {
	windowSize := w.Bounds().Size()
	t.width = windowSize.X - t.config.Window.Padding.X
	t.height = windowSize.Y - t.config.Window.Padding.Y
	t.window = w

	atlas := atlas.NewAtlas(t.font, atlas.ASCII)

	t.text = text.New(pixel.Vec{
		X: t.config.Window.Padding.X,
		Y: t.height - atlas.Ascent() - t.config.Window.Padding.Y,
	}, atlas)

	cols := int(t.width / t.text.LineHeight)
	rows := int(t.height / (t.text.LineHeight))

	t.text.SetSize(rows, cols)

}

func (t *Terminal) input() {
	scroll := t.window.MouseScroll()

	switch {
	case scroll.Y != 0:
		t.text.Scroll(int(scroll.Y))

	case t.window.JustPressed(pixelgl.KeyEnter):
		t.pty.Write([]byte{'\n'})

	case t.window.JustReleased(pixelgl.KeyTab):
		t.pty.Write([]byte{'\t'})

	case t.window.JustPressed(pixelgl.KeyBackspace):
		t.pty.Write([]byte{8})

	default:
		typed := t.window.Typed()

		if typed != "" {
			t.pty.Write([]byte(typed))

		}
	}
}

func (t *Terminal) readPty() {
	buf := make([]byte, 2048)
	for {
		num, err := t.pty.Read(buf)
		if err != nil {
			log.Printf("failed to read from pty: %v \n", err)
			break
		}

		t.text.Write(buf[:num])

	}
}

func (t *Terminal) Run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:     "gterm",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	t.setupWindow(win)

	go t.readPty()

	fps := time.Tick(time.Second / 60)
	for !win.Closed() {
		win.Clear(colornames.Black)

		t.input()

		t.text.Draw(t.window, pixel.IM)
		win.Update()

		<-fps
	}
}
