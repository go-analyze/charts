package charts

import (
	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

// FontFamilyRoboto is the default chart font (Roboto Medium), it provides a well spaced Sans style font with good latin character support.
const FontFamilyRoboto = "roboto"

// FontFamilyNotoSans provides Noto Sans Display Medium, a slightly more condensed Sans variant compared to FontFamilyRoboto.
// This font offers better internal character and some symbol and emoji support.
const FontFamilyNotoSans = "notosans"

// FontFamilyNotoSansBold provides Noto Sans Display Extra Bold, a bold version of FontFamilyNotoSans.
const FontFamilyNotoSansBold = "notosans-bold"

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
