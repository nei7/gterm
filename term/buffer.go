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

	savedRows uint16
	savedCols uint16

	cursorPos struct{ X, Y int }
}

func NewBuffer() *Buffer {
	buf := &Buffer{
		lines: []Line{},
	}

	return buf
}

func (buf *Buffer) setSize(rows, cols uint16) {
	buf.rows = rows
	buf.cols = cols
}

func (buf *Buffer) insertLine() {
	buf.lines = append(buf.lines, Line{})
}

func (buf *Buffer) insertChar(char Char) {

	for len(buf.lines)-1 < buf.cursorPos.Y {
		buf.insertLine()
	}

	line := buf.cursorPos.Y
	if line < 0 || line > len(buf.lines) {
		return
	}

	for len(buf.lines[line].Chars) < buf.cursorPos.X {
		buf.lines[line].Chars = append(buf.lines[line].Chars, Char{
			R: ' ',
		})
	}

	buf.lines[line].Chars = append(buf.lines[line].Chars, char)
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

func (buf *Buffer) getLine(index int) *Line {
	if index < 0 || index >= len(buf.lines) {
		return nil
	}

	return &buf.lines[index]
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
