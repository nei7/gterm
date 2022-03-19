package main

import (
	"image/color"
	"io"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/creack/pty"
)

type Terminal struct {
	content   *widget.TextGrid
	cursor    *canvas.Rectangle
	cursorRow int
	cursorCol int

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

func (t *Terminal) startPty() error {
	_ = os.Chdir(t.homedir)

	os.Setenv("TERM", "dumb")
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

func (t *Terminal) handleOutput(buf []byte) {
	str := string(buf)

	for _, r := range []rune(str) {

		if r == '\r' {
			continue
		}

		if r == asciiBell {
			break
		}

		if r == '\n' || len(t.content.Rows) == 0 {
			t.cursorRow++
			t.content.SetRow(t.cursorRow, widget.TextGridRow{})
			t.cursorCol = 0
			continue
		}

		if r == asciiBackspace {
			t.cursorCol--

			row := t.content.Rows[t.cursorRow]
			row.Cells = row.Cells[:t.cursorCol]

			t.content.SetRow(t.cursorRow, row)
			continue
		}
		t.content.SetCell(t.cursorRow, t.cursorCol, widget.TextGridCell{
			Rune: r,
		})

		t.cursorCol++
	}

}

func (t *Terminal) start() {
	buf := make([]byte, 4096)

	for {
		time.Sleep(50)

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
