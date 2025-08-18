package charts

import (
	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

const defaultFontSize = 12.0

// InstallFont installs a font for chart rendering.
func InstallFont(fontFamily string, data []byte) error {
	return chartdraw.InstallFont(fontFamily, data)
}

func getPreferredFont(fonts ...*truetype.Font) *truetype.Font {
	for _, font := range fonts {
		if font != nil {
			return font
		}
	}
	return GetDefaultFont()
}

// GetDefaultFont returns the default font.
func GetDefaultFont() *truetype.Font {
	return chartdraw.GetDefaultFont()
}

// SetDefaultFont sets the default font by name.
func SetDefaultFont(fontFamily string) error {
	return chartdraw.SetDefaultFont(fontFamily)
}

// GetFont returns the font by family name, or the default if the family is not installed.
func GetFont(fontFamily string) *truetype.Font {
	return chartdraw.GetFont(fontFamily)
}
