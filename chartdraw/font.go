package chartdraw

import (
	"errors"
	"fmt"
	"sync"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw/roboto"
)

var fonts = sync.Map{}
var defaultFontFamily = "default"

func init() {
	name := "roboto"
	if err := InstallFont(name, roboto.Roboto); err != nil {
		panic(fmt.Errorf("could not install default font - %w", err))
	} else if err = SetDefaultFont(name); err != nil {
		panic(fmt.Errorf("could not set default font - %w", err))
	}
}

// InstallFont installs the font for charts
func InstallFont(fontFamily string, data []byte) error {
	font, err := truetype.Parse(data)
	if err != nil {
		return err
	}
	fonts.Store(fontFamily, font)
	return nil
}

// GetDefaultFont get default font.
func GetDefaultFont() *truetype.Font {
	return GetFont(defaultFontFamily)
}

// SetDefaultFont set default font by name.
func SetDefaultFont(fontFamily string) error {
	if value, ok := fonts.Load(fontFamily); ok {
		fonts.Store(defaultFontFamily, value)
		return nil
	}
	return errors.New("font not found: " + fontFamily)
}

// GetFont get the font by font family or the default if the family is not installed.
func GetFont(fontFamily string) *truetype.Font {
	if value, ok := fonts.Load(fontFamily); ok {
		if f, ok := value.(*truetype.Font); ok {
			return f
		}
	}
	return GetDefaultFont()
}
