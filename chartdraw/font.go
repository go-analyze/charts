package chartdraw

import (
	"errors"
	"sync/atomic"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

var defaultFontFamily atomic.Value

func init() {
	defaultFontFamily.Store("roboto")
}

// InstallFont installs the font for charts.
func InstallFont(fontFamily string, data []byte) error {
	return drawing.InstallFont(fontFamily, data)
}

// GetDefaultFont get default font.
func GetDefaultFont() *truetype.Font {
	return drawing.GetFont(defaultFontFamily.Load().(string))
}

// SetDefaultFont set default font by name.
func SetDefaultFont(fontFamily string) error {
	if font := drawing.GetFont(fontFamily); font == nil {
		return errors.New("font not found: " + fontFamily)
	}
	defaultFontFamily.Store(fontFamily)
	return nil
}

// GetFont get the font by font family or the default if the family is not installed.
func GetFont(fontFamily string) *truetype.Font {
	result := drawing.GetFont(fontFamily)
	if result != nil {
		return result
	}
	return drawing.GetFont(defaultFontFamily.Load().(string))
}
