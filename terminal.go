package main

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/creack/pty"
)

type Terminal struct {
	content *widget.TextGrid
	cursor  *canvas.Rectangle

	homedir string
	title   string

	pty io.Closer
	in  io.Writer
	out io.Reader
}

func NewTerminal() *Terminal {
	cursor := canvas.NewRectangle(color.White)
	cursor.Resize(fyne.NewSize(10, 18))

	t := &Terminal{
		content: widget.NewTextGrid(),
		homedir: homeDir(),
		cursor:  cursor,
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

func (t *Terminal) setPty() error {
	_ = os.Chdir(t.homedir)
	env := os.Environ()
	env = append(env, "TERM=xterm-256color")

	cmd := exec.Command(os.Getenv("SHELL"))
	cmd.Env = env
	pt, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	t.in = pt
	t.out = pt
	t.pty = pt

	return nil
}

func (t *Terminal) handleOutput(buf []byte) {
	str := string(buf)

	fmt.Println(str)

	t.content.Rows = append(t.content.Rows, widget.TextGridRow{})
	for _, r := range []rune(str) {

		t.content.Rows[len(t.content.Rows)-1].Cells = append(t.content.Rows[0].Cells, widget.TextGridCell{
			Rune: r,
		})
	}

}

func (t *Terminal) start() {
	buf := make([]byte, 4096)
	for {
		num, err := t.out.Read(buf)
		if err != nil {
			fyne.LogError("Failed to read pty", err)
		}

		t.handleOutput(buf[:num])
		if num < 4096 {
			t.content.Refresh()
		}
	}
}
