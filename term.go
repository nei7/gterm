package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/creack/pty"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type Terminal struct {
	mu   sync.Mutex
	text *text.Text

	content      [][]rune
	lines        int
	offsetY      int
	screenHeight float64

	homedir string
	title   string

	pty io.Closer
	in  io.Writer
	out io.Reader

	config *Config
}

func NewTerminal() *Terminal {
	t := &Terminal{}

	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %w", err)
	}

	t.config = config

	t.homedir = homeDir()

	return t
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return home
}

func (t *Terminal) startPty() error {
	_ = os.Chdir(t.homedir)

	os.Setenv("TERM", "xterm-256color")
	cmd := exec.Command("/bin/bash")

	pt, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	t.in = pt
	t.out = pt
	t.pty = pt

	return nil
}

func (t *Terminal) renderContent() {
	t.text.Clear()
	startOffset, endOffset := calcOffset(len(t.content), t.offsetY, int(t.screenHeight), int(t.config.Font.Size))
	for _, line := range t.content[startOffset:endOffset] {
		for _, r := range line {
			t.text.WriteRune(r)
		}
		t.text.WriteRune('\n')
	}
}

func (t *Terminal) runUI() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:     "gterm",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	face, err := loadTTF(fmt.Sprintf("/usr/share/fonts/TTF/%s.ttf", t.config.Font.Family), t.config.Font.Size)
	if err != nil {
		face = basicfont.Face7x13
	}
	atlas := text.NewAtlas(face, text.ASCII)

	t.screenHeight = win.Bounds().H()
	t.text = text.New(pixel.V(5, t.screenHeight-20), atlas)

	buf := make([]byte, 4096)
	for !win.Closed() {
		win.Clear(colornames.Black)

		h := (int(t.screenHeight)/(int(t.config.Font.Size)+3) - 3)
		if scroll := int(win.MouseScroll().Y); t.offsetY-scroll >= 0 && t.offsetY-scroll <= len(t.content)-h {
			t.offsetY -= scroll
		}

		go func() {
			t.mu.Lock()
			num, err := t.out.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			t.handleOutput(buf[:num])
			if len(t.content)-h > 0 {
				t.offsetY = len(t.content) - h
			}
			t.mu.Unlock()
		}()

		t.in.Write([]byte(win.Typed()))

		if win.JustPressed(pixelgl.KeyBackspace) {
			t.in.Write([]byte{asciiBackspace})
		}
		if win.JustPressed(pixelgl.KeyEnter) {
			t.in.Write([]byte{'\n'})
		}

		t.renderContent()
		t.text.Draw(win, pixel.IM)

		win.Update()
	}
}
