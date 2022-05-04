package term

import (
	"image/color"
)

type Char struct {
	Id      int
	R       rune
	FgColor color.Color
	BgColor color.Color
}

type Line struct {
	Chars []Char
}

type Buffer struct {
	lines []Line

	rows uint16
	cols uint16

	savedCursorPos struct{ X, Y int }
	cursorPos      struct{ X, Y int }
}

func NewBuffer() *Buffer {
	buf := &Buffer{
		lines: []Line{},
	}

	return buf
}

func (buf *Buffer) insertLine() {
	buf.lines = append(buf.lines, Line{})
}

func (buf *Buffer) insertChar(char Char) {
	if buf.cursorPos.X < 0 || buf.cursorPos.Y < 0 {
		return
	}

	if buf.cursorPos.X >= int(buf.cols) {
		return
	}

	for len(buf.lines)-1 < buf.cursorPos.Y {
		buf.insertLine()
	}

	for len(buf.lines[buf.cursorPos.Y].Chars)-1 < buf.cursorPos.X {
		buf.lines[buf.cursorPos.Y].Chars = append(buf.lines[buf.cursorPos.Y].Chars, Char{
			R: ' ',
		})
	}

	cell := buf.lines[buf.cursorPos.Y].Chars[buf.cursorPos.X]
	if cell.R != char.R || char.FgColor != cell.FgColor || char.BgColor != cell.BgColor {
		cell.R = char.R
		cell.BgColor = char.BgColor
		cell.FgColor = char.FgColor

		for len(buf.lines) <= buf.cursorPos.Y {
			buf.insertLine()
		}
		data := buf.lines[buf.cursorPos.Y]

		for len(data.Chars) <= buf.savedCursorPos.X {
			data.Chars = append(data.Chars, Char{
				R: ' ',
			})
			buf.lines[buf.cursorPos.Y] = data
		}

		buf.lines[buf.cursorPos.Y].Chars[buf.cursorPos.X] = cell
	}

	buf.cursorPos.X++
}

func (buf *Buffer) clear() {
	buf.cursorPos = struct {
		X int
		Y int
	}{
		0, 0,
	}
	buf.cursorPos.X = 0
	buf.cursorPos.Y = 0

	buf.lines = []Line{}
}

func (buf *Buffer) SetRow(row int, content []Char) {
	if row < 0 {
		return
	}
	for len(buf.lines) <= row {
		buf.lines = append(buf.lines, Line{})
	}

	buf.lines[row] = Line{Chars: content}
}

func (buf *Buffer) Row(row int) Line {
	if row < 0 || row >= len(buf.lines) {
		return Line{}
	}

	return buf.lines[row]
}
