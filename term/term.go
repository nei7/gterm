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
	Buffer *Buffer

	title string
	debug bool

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
	t.Buffer = buffer

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
	_ = os.Chdir(homedir)

	os.Setenv("TERM", "xterm-256color")
	cmd := exec.Command("/bin/bash")

	pt, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	return pt, nil
}

func (t *Terminal) SetSize(rows, cols uint16) error {
	t.Buffer.SetSize(rows, cols)
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
			log.Printf("failed to read from pty: %v \n", err)
			break
		}
		t.Print(buf[:num])

		t.Buffer.ScrollToBottom()
	}
}

func (t *Terminal) Backspace() {
	last := &t.Buffer.lines[t.Buffer.cursorPos.Y]

	last.Chars = last.Chars[:t.Buffer.cursorPos.X-1]
	t.moveCursor(t.Buffer.cursorPos.Y, t.Buffer.cursorPos.X-1)
}

func (t *Terminal) Clear() {
	t.Buffer.clear()
}
