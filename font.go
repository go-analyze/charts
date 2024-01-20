package charts

import (
	"errors"
	"sync"

	"github.com/golang/freetype/truetype"
	"github.com/wcharczuk/go-chart/v2/roboto"
)

var fonts = sync.Map{}
var ErrFontNotExists = errors.New("font is not exists")
var defaultFontFamily = "defaultFontFamily"

func init() {
	name := "roboto"
	_ = InstallFont(name, roboto.Roboto)
	font, _ := GetFont(name)
	SetDefaultFont(font)
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

// GetDefaultFont get default font
func GetDefaultFont() (*truetype.Font, error) {
	return GetFont(defaultFontFamily)
}

// SetDefaultFont set default font
func SetDefaultFont(font *truetype.Font) {
	if font == nil {
		return
	}
	fonts.Store(defaultFontFamily, font)
}

// GetFont get the font by font family
func GetFont(fontFamily string) (*truetype.Font, error) {
	value, ok := fonts.Load(fontFamily)
	if !ok {
		return nil, ErrFontNotExists
	}
	f, ok := value.(*truetype.Font)
	if !ok {
		return nil, ErrFontNotExists
	}
	return f, nil
}
