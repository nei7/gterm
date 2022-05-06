package font

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Manager struct {
	family   string
	regular  font.Face
	bold     font.Face
	size     float64
	dpi      float64
	charSize image.Point
}

func NewManager() *Manager {
	return &Manager{
		dpi:  72,
		size: 20,
	}
}

func (m *Manager) loadFontFace(path string) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fnt, err := opentype.ParseReaderAt(file)
	if err != nil {
		return nil, err
	}

	return m.createFace(fnt)
}

func (m *Manager) createFace(f *opentype.Font) (font.Face, error) {
	return opentype.NewFace(f, &opentype.FaceOptions{
		Size:    m.size,
		DPI:     m.dpi,
		Hinting: font.HintingFull,
	})
}

func (m *Manager) SetFont(name string) error {
	m.family = name

	font, err := MatchFont("SauceCodePro Nerd Font Mono", "Regular")
	if err != nil {

		log.Fatal(err)
	}

	f, err := m.loadFontFace(font.Path)
	if err != nil {
		return err
	}
	m.regular = f

	return m.calcMetrics()
}

func (m *Manager) calcMetrics() error {

	face := m.regular

	var prevAdvance int
	for ch := rune(32); ch <= 127; ch++ {
		width, ok := face.GlyphAdvance(ch)
		if ok && width > 0 {
			advance := int(width)
			if prevAdvance > 0 && prevAdvance != advance {
				return fmt.Errorf("the specified font is not monospaced: %d 0x%X=%d", prevAdvance, ch, advance)
			}
			prevAdvance = advance
		}
	}

	if prevAdvance == 0 {
		return fmt.Errorf("failed to calculate advance width for font face")
	}

	metrics := face.Metrics()

	m.charSize.X = int(math.Round(float64(prevAdvance) / m.dpi))
	m.charSize.Y = int(math.Round(float64(metrics.Height) / m.dpi))
	return nil
}

func (m *Manager) RegularFontFace() font.Face {
	return m.regular
}

func (m *Manager) CharSize() image.Point {
	return m.charSize
}
