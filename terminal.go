package main

import (
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

	t := &Terminal{
		content: widget.NewTextGrid(),
		homedir: homeDir(),
		cursor:  createCursor(),
	}

	return t
}

func createCursor() *canvas.Rectangle {
	cursor := canvas.NewRectangle(color.White)
	cursor.Resize(fyne.Size{
		Width:  7,
		Height: 14,
	})
	cursor.FillColor = color.Transparent
	cursor.StrokeColor = color.White
	cursor.StrokeWidth = 1.7

	return cursor
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

func (t *Terminal) handleOutput(buf []byte) {

	for _, r := range []rune(string(buf)) {

		if r == '\r' {
			continue
		}

		if r == asciiBell {
			break
		}

		if len(t.content.Rows) == 0 {
			t.content.SetRow(t.cursorRow, widget.TextGridRow{})
			t.cursorCol = 0
		}

		if r == '\n' {
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

func (t *Terminal) Refresh() {
	t.cursor.Move(fyne.NewPos(9*float32(t.cursorCol)+5, 18*float32(t.cursorRow)+2))
	t.cursor.Refresh()
	t.content.Refresh()
}

func (t *Terminal) start() {
	buf := make([]byte, 2048)

	for {
		num, err := t.out.Read(buf)
		if err != nil {
			if err != nil {
				break
			}

			fyne.LogError("Failed to read pty", err)
		}
		t.handleOutput(buf[:num])
		t.Refresh()
	}
}
