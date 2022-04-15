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

func (line *Line) push(char Char) {
	line.Chars = append(line.Chars, char)
}

func (line *Line) pop() {
	l := len(line.Chars)
	if l > 0 {
		line.Chars = line.Chars[:l-1]
	}
}

type Buffer struct {
	lines []Line

	rows uint16
	cols uint16

	savedRows uint16
	savedCols uint16

	scrollOffset int
	cursorPos    struct{ X, Y int }
}

func NewBuffer() *Buffer {
	buf := &Buffer{
		lines: []Line{},
	}

	return buf
}

func (buf *Buffer) SetSize(rows, cols uint16) {
	buf.rows = rows
	buf.cols = cols
}

func (buf *Buffer) insertLine(line Line) {
	buf.lines = append(buf.lines, line)

}

func (buf *Buffer) appendToLine(line int, char Char) {
	if len(buf.lines) == 0 {
		buf.insertLine(Line{})
	}

	if line < 0 || line > len(buf.lines) {
		return
	}
	buf.lines[line].push(char)
	buf.cursorPos.X++
}

func (buf *Buffer) ScrollDown() {
	if buf.scrollOffset < len(buf.lines)-int(buf.rows) {
		buf.scrollOffset++
	}
}

func (buf *Buffer) ScrollUp() {
	if buf.scrollOffset > 0 {
		buf.scrollOffset--
	}
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

	buf.scrollOffset = 0

	buf.lines = []Line{}
}

func (buf *Buffer) ScrollToBottom() {
	if len(buf.lines)-int(buf.rows) > 0 {
		buf.scrollOffset = len(buf.lines) - int(buf.rows)
	}
}

func (buf *Buffer) GetLine(index int) *Line {
	if index < 0 || index > len(buf.lines) {
		return nil
	}

	return &buf.lines[index]
}

func (buf *Buffer) GetLines() []Line {

	if len(buf.lines) < int(buf.rows) {
		return buf.lines
	}

	offset := int(buf.rows) + buf.scrollOffset
	if length := len(buf.lines); offset >= length {
		return buf.lines[length-int(buf.rows) : length]
	}

	return buf.lines[buf.scrollOffset:offset]
}
