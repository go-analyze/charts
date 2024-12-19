package charts

import (
	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

const defaultFontSize = 12.0

// InstallFont installs the font for charts
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

// GetDefaultFont get default font.
func GetDefaultFont() *truetype.Font {
	return chartdraw.GetDefaultFont()
}

// SetDefaultFont set default font by name.
func SetDefaultFont(fontFamily string) error {
	return chartdraw.SetDefaultFont(fontFamily)
}

// GetFont get the font by font family or the default if the family is not installed.
func GetFont(fontFamily string) *truetype.Font {
	return chartdraw.GetFont(fontFamily)
}
