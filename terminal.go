package main

import (
	"io"
	"log"
	"os"
	"os/exec"

	"golang.org/x/image/colornames"

	"github.com/creack/pty"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Terminal struct {
	cursorRow, cursorCol int

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

func (t *Terminal) run() {
	buf := make([]byte, 4096)
	ui := NewUI()

	for !ui.window.Closed() {
		ui.window.Clear(colornames.Black)

		go func() {
			ui.mu.Lock()
			num, err := t.out.Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			ui.text.Write(buf[:num])

			ui.mu.Unlock()
		}()

		t.in.Write([]byte(ui.window.Typed()))

		if ui.window.JustPressed(pixelgl.KeyEnter) {

			t.in.Write([]byte{'\n'})

		}
		if ui.window.JustPressed(pixelgl.KeyBackspace) {
			ui.text.Clear()
		}

		ui.text.Draw(ui.window, pixel.IM)

		ui.window.Update()

	}
}
