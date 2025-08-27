package drawing

import (
	"bytes"
	"compress/gzip"
	"embed"
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/golang/freetype/truetype"
)

// FallbackFonts is the ordered list of fallback fonts to try when a character is missing.
var FallbackFonts = []string{"notosans", "roboto"}

//go:embed fonts/*.ttf.gz
var fontFiles embed.FS
var embeddedFontNames map[string]string // maps normalized font name -> filename
var fonts = sync.Map{}

func init() {
	embeddedFontNames = make(map[string]string)
	// Discover embedded font files
	if entries, err := fontFiles.ReadDir("fonts"); err == nil {
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
}

// InstallFont installs the font in the global registry.
func InstallFont(fontFamily string, data []byte) error {
	fontFamily = strings.TrimSpace(strings.ToLower(fontFamily))

	font, err := truetype.Parse(data)
	if err != nil {
		return err
	}
	fonts.Store(fontFamily, font)
	return nil
}

// GetFont get the font by font family or nil if the family is not installed.
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

	return nil
}

func gzipDecompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()
	return io.ReadAll(r)
}

// loadEmbeddedFont lazily loads and installs an embedded font.
func loadEmbeddedFont(fontName string) error {
	filename, exists := embeddedFontNames[fontName]
	if !exists {
		return errors.New("embedded font not found: " + fontName)
	}

	// Load and decompress font
	compressed, err := fontFiles.ReadFile("fonts/" + filename)
	if err != nil {
		return err
	}
	decompressed, err := gzipDecompress(compressed)
	if err != nil {
		return err
	}

	return InstallFont(fontName, decompressed)
}
