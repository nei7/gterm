package utils

import (
	"os"
)

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return home
}
