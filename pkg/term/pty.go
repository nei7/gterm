package term

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

func startPty(homedir string) (*os.File, error) {
	_ = os.Chdir(homedir)

	os.Setenv("TERM", "gterm")
	cmd := exec.Command("/bin/bash")

	pt, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	return pt, nil
}
