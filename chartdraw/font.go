package chartdraw

import (
	"bytes"
	"compress/gzip"
	"embed"
	"errors"
	"io"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/golang/freetype/truetype"
)

//go:embed drawing/fonts/*.ttf.gz
var fontFiles embed.FS
var embeddedFontNames map[string]string // maps normalized font name -> filename
var fonts = sync.Map{}
var defaultFontFamily atomic.Value

func init() {
	embeddedFontNames = make(map[string]string)
	// Discover embedded font files
	if entries, err := fontFiles.ReadDir("drawing/fonts"); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".ttf.gz") {
				// Extract font name: NotoSans-Regular.ttf.gz -> notosans, NotoSans-Bold.ttf.gz -> notosans-bold
				name := strings.TrimSuffix(entry.Name(), ".ttf.gz")
				// Only remove -Regular and -Medium suffixes, keep others like -Bold
				if strings.HasSuffix(strings.ToLower(name), "-regular") {
					name = name[:len(name)-8] // Remove "-Regular"
				} else if strings.HasSuffix(strings.ToLower(name), "-medium") {
					name = name[:len(name)-7] // Remove "-Medium"
				}
				name = strings.ToLower(name)
				embeddedFontNames[name] = entry.Name()
			}
		}
	}
	defaultFontFamily.Store("roboto")
}

// InstallFont installs the font for charts
func InstallFont(fontFamily string, data []byte) error {
	fontFamily = strings.TrimSpace(strings.ToLower(fontFamily))

	font, err := truetype.Parse(data)
	if err != nil {
		return err
	}
	fonts.Store(fontFamily, font)
	return nil
}

// GetDefaultFont get default font.
func GetDefaultFont() *truetype.Font {
	return GetFont(defaultFontFamily.Load().(string))
}

// SetDefaultFont set default font by name.
func SetDefaultFont(fontFamily string) error {
	fontFamily = strings.TrimSpace(strings.ToLower(fontFamily))

	if _, ok := fonts.Load(fontFamily); ok {
		defaultFontFamily.Store(fontFamily)
		return nil
	}
	// check if it's an available font not yet loaded
	if _, exists := embeddedFontNames[fontFamily]; exists {
		defaultFontFamily.Store(fontFamily)
		return nil
	}
	return errors.New("font not found: " + fontFamily)
}

// GetFont get the font by font family or the default if the family is not installed.
func GetFont(fontFamily string) *truetype.Font {
	fontFamily = strings.TrimSpace(strings.ToLower(fontFamily))

	if value, ok := fonts.Load(fontFamily); ok {
		return value.(*truetype.Font)
	}

	// Check if it's an embedded font not yet loaded
	if _, exists := embeddedFontNames[fontFamily]; exists {
		if err := loadEmbeddedFont(fontFamily); err == nil {
			if value, ok := fonts.Load(fontFamily); ok {
				return value.(*truetype.Font)
			}
		}
	}

	return GetDefaultFont()
}

func gzipDecompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()
	return io.ReadAll(r)
}

// loadEmbeddedFont lazily loads and installs an embedded font
func loadEmbeddedFont(fontName string) error {
	filename, exists := embeddedFontNames[fontName]
	if !exists {
		return errors.New("embedded font not found: " + fontName)
	}

	// Load and decompress font
	compressed, err := fontFiles.ReadFile("drawing/fonts/" + filename)
	if err != nil {
		return err
	}
	decompressed, err := gzipDecompress(compressed)
	if err != nil {
		return err
	}

	return InstallFont(fontName, decompressed)
}
