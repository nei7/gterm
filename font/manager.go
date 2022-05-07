package font

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/faiface/pixel"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Manager struct {
	family       string
	regular      font.Face
	bold         font.Face
	size         float64
	dpi          float64
	charSize     pixel.Vec
	fontDotDepth int
}

var styles = []string{"Regular", "Bold"}

const (
	Regular = 0
	Bold    = 1
)

func NewManager() *Manager {
	return &Manager{
		dpi:  72,
		size: 12,
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

	m.loadFontStyles()

	return m.calcMetrics()
}

func (m *Manager) loadFontStyles() error {
	for i, style := range styles {
		font, err := MatchFont(m.family, style)
		if err != nil {

			log.Fatal(err)
		}

		f, err := m.loadFontFace(font.Path)
		if err != nil {
			return err
		}

		switch i {
		case Regular:
			m.regular = f
		case Bold:
			m.bold = f
		}
	}

	return nil
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

	m.charSize.X = (math.Round(float64(prevAdvance) / m.dpi))
	m.charSize.Y = (math.Round(float64(metrics.Height) / m.dpi))

	return nil
}

func (m *Manager) RegularFontFace() font.Face {
	return m.regular
}

func (m *Manager) BoldFontFace() font.Face {
	return m.bold
}

func (m *Manager) CharSize() pixel.Vec {
	return m.charSize
}

func (m *Manager) DotDepth() int {
	return m.fontDotDepth
}
