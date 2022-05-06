package font

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
)

var validExtensions = []string{
	".ttf",
	".otf",
}

var fontDirs = []string{
	"~/.fonts",
	"~/.local/share/fonts",
	"/usr/local/share/fonts",
	"/usr/share/fonts",
}

var home string

func init() {
	if user, _ := user.Current(); user != nil {
		home = user.HomeDir
	}
}

type Font struct {
	Family string
	Style  string
	Path   string
}

func MatchFont(fontFamily string, fontStyle string) (*Font, error) {
	var font *Font
	for _, dir := range fontDirs {

		if strings.HasPrefix(dir, "~/") {
			dir = filepath.Join(home, dir[2:])
		}

		if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
			continue
		}

		if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			ext := filepath.Ext(path)
			for _, valid := range validExtensions {
				if strings.EqualFold(ext, valid) {
					f, err := os.Open(path)
					if err != nil {
						return err
					}
					defer f.Close()

					bytes, err := ioutil.ReadAll(f)
					meta, err := freetype.ParseFont(bytes)
					if err != nil {
						continue
					}
					family, style := meta.Name(1), meta.Name(2)

					if strings.EqualFold(family, fontFamily) && strings.EqualFold(style, fontStyle) {
						font = &Font{
							Family: family,
							Style:  style,
							Path:   path,
						}
						return nil
					}

				}
			}

			return nil
		}); err != nil {
			return nil, err
		}
	}

	if font == nil {
		return nil, errors.New("not matched")
	}

	return font, nil
}
