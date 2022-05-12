package term

import (
	"image/color"
)

type Char struct {
	R       rune
	FgColor color.Color
	BgColor color.Color
}

type Line struct {
	Chars []Char
}

type Cord struct{ X, Y int }

type Buffer struct {
	lines []Line

	rows int
	cols int

	savedCursorPos Cord
	cursorPos      Cord

	scrollY                 int
	scrollTop, scrollBottom int
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

func (buf *Buffer) moveCursor(row, col int) {
	buf.cursorPos.X = col
	buf.cursorPos.Y = row
}

func (buf *Buffer) backspace() {
	last := buf.Row(buf.cursorPos.Y).Chars
	if len(last) > 0 {
		last = last[:buf.cursorPos.X-1]
	}
	buf.moveCursor(buf.cursorPos.Y, buf.cursorPos.X-1)
}

func (buf *Buffer) clear() {
	pos := Cord{0, 0}
	buf.cursorPos = pos
	buf.savedCursorPos = pos
	buf.lines = []Line{}
	buf.scrollY = 0
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

func (buf *Buffer) Size() Cord {
	return Cord{
		X: buf.cols,
		Y: buf.rows,
	}
}

func (buf *Buffer) Len() int {
	return len(buf.lines)
}

func (buf *Buffer) ScrollDown() {
	if buf.scrollY < len(buf.lines)-int(buf.rows) {
		buf.scrollY++
	}
}

func (buf *Buffer) ScrollUp() {
	if buf.scrollY > 0 {
		buf.scrollY--
	}
}

func (buf *Buffer) ScrollToBottom() {
	if len(buf.lines)-int(buf.rows) > 0 {
		buf.scrollY = len(buf.lines) - int(buf.rows)
	}
}

func (buf *Buffer) ScrollY() int {
	return buf.scrollY
}
