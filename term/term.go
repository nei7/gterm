package term

import (
	"image/color"
	"log"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"golang.org/x/image/colornames"
)

type Terminal struct {
	pty    *os.File
	buffer *Buffer

	title string
	debug bool

	scrollOffset int

	currentFG color.Color
	currentBG color.Color
	bright    bool
}

func New() *Terminal {
	t := &Terminal{}

	pty, err := startPty(getHomeDir())
	if err != nil {
		log.Fatal(err)
	}
	t.pty = pty

	buffer := NewBuffer()
	t.buffer = buffer

	t.currentFG = colornames.White
	t.currentBG = colornames.Black

	return t
}

func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return home
}

func startPty(homedir string) (*os.File, error) {
	_ = os.Chdir(homedir + "Projects/gterm")

	os.Setenv("TERM", "xterm-256color")
	cmd := exec.Command("/bin/bash")

	pt, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	return pt, nil
}

func (t *Terminal) SetSize(rows, cols uint16) error {
	t.buffer.setSize(rows, cols)
	if err := pty.Setsize(t.pty, &pty.Winsize{
		Rows: rows,
		Cols: cols,
	}); err != nil {
		return err
	}
	return nil
}

func (t *Terminal) setTitle(title string) {
	t.title = title
}

func (t *Terminal) Write(buf []byte) error {
	_, err := t.pty.Write(buf)
	return err
}

func (t *Terminal) Run() {
	buf := make([]byte, 2048)
	for {
		num, err := t.pty.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				os.Exit(1)
			} else if err, ok := err.(*os.PathError); ok && err.Err.Error() == "input/output error" {
				os.Exit(1)
			}
			log.Printf("failed to read from pty: %v \n", err)
			break
		}

		t.Print(buf[:num])
		t.ScrollToBottom()
	}
}

func (t *Terminal) Backspace() {
	last := &t.buffer.lines[t.buffer.cursorPos.Y]
	last.Chars = last.Chars[:t.buffer.cursorPos.X-1]
	t.moveCursor(t.buffer.cursorPos.Y, t.buffer.cursorPos.X-1)
}

func (t *Terminal) Clear() {
	t.buffer.clear()
}

func (t *Terminal) GetLines() []Line {
	if len(t.buffer.lines) < int(t.buffer.rows) {
		return t.buffer.lines
	}

	offset := int(t.buffer.rows) + t.scrollOffset
	if length := len(t.buffer.lines); offset >= length {
		return t.buffer.lines[t.scrollOffset:length]
	}

	return t.buffer.lines[t.scrollOffset:offset]
}

// Handle scroll
func (t *Terminal) ScrollDown() {
	if t.scrollOffset < len(t.buffer.lines)-int(t.buffer.rows) {
		t.scrollOffset++
	}
}

func (t *Terminal) ScrollUp() {
	if t.scrollOffset > 0 {
		t.scrollOffset--
	}
}

func (t *Terminal) ScrollToBottom() {

	if len(t.buffer.lines)-int(t.buffer.rows) > 0 {
		t.scrollOffset = len(t.buffer.lines) - int(t.buffer.rows)
	}
}
