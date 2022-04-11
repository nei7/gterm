package term

import (
	"fmt"
	"image/color"
	"sync"
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
	sync.RWMutex

	lines []Line

	rows uint16
	cols uint16

	cursorPos    struct{ X, Y int }
	scrollOffset int
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

func (buf *Buffer) ScrollDown() {
	if buf.scrollOffset < len(buf.lines)-int(buf.rows) {
		buf.scrollOffset++
	}
}

func (buf *Buffer) ScrollUp() {
	fmt.Println(buf.scrollOffset)
	if buf.scrollOffset > 0 {
		buf.scrollOffset--
	}
}

func (buf *Buffer) ScrollToBottom() {
	buf.scrollOffset = len(buf.lines) - int(buf.rows)
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
