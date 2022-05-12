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
	title string

	pty    *os.File
	buffer *Buffer

	currentFG color.Color
	currentBG color.Color
	bright    bool
}

func New(buffer *Buffer) *Terminal {
	t := &Terminal{}

	pty, err := startPty(getHomeDir())
	if err != nil {
		log.Fatal(err)
	}
	t.pty = pty

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
	_ = os.Chdir(homedir)

	os.Setenv("TERM", "xterm-256color")
	cmd := exec.Command("/bin/bash")

	pt, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	return pt, nil
}

func (t *Terminal) SetSize(rows, cols int) error {
	if err := pty.Setsize(t.pty, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	}); err != nil {
		return err
	}

	t.buffer.rows = rows
	t.buffer.cols = cols
	return nil
}

func (t *Terminal) setTitle(title string) {
	t.title = title
}

func (t *Terminal) Write(buf []byte) error {
	_, err := t.pty.Write(buf)
	return err
}

func (t *Terminal) Run(updateChan chan struct{}) {
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

		t.handleOutput(buf[:num])
		t.buffer.ScrollToBottom()

		updateChan <- struct{}{}

	}
}

func (t *Terminal) Clear() {
	t.buffer.clear()
}
