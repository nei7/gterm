package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"golang.org/x/image/colornames"

	"github.com/creack/pty"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type Terminal struct {
	mu     sync.Mutex
	window *pixelgl.Window
	text   *text.Text

	content [][]rune
	lines   int

	offsetY int

	homedir string
	title   string

	pty io.Closer
	in  io.Writer
	out io.Reader
}

func NewTerminal() *Terminal {

	home := homeDir()

	t := &Terminal{
		homedir: home,
	}

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

	os.Setenv("TERM", "gterm")
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

	startOffset, endOffset := calcOffset(len(t.content), t.offsetY)
	for _, line := range t.content[startOffset:endOffset] {
		for _, r := range line {
			t.text.WriteRune(r)

		}
		t.text.WriteRune('\n')
	}
}

func (t *Terminal) runUI() {
	buf := make([]byte, 4096)
	t.window = createWindow()

	t.loadFont()

	for !t.window.Closed() {

		t.window.Clear(colornames.Black)

		h := (int(screenHeight)/(fontSize+3) - 3)
		if scroll := int(t.window.MouseScroll().Y); t.offsetY-scroll >= 0 && t.offsetY-scroll <= len(t.content)-h {
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

		t.in.Write([]byte(t.window.Typed()))

		if t.window.JustPressed(pixelgl.KeyBackspace) {
			t.in.Write([]byte{asciiBackspace})
		}
		if t.window.JustPressed(pixelgl.KeyEnter) {
			t.in.Write([]byte{'\n'})

		}

		t.renderContent()
		t.text.Draw(t.window, pixel.IM)

		t.window.Update()
	}
}
