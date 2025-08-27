package drawing

import (
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getTestFontData loads font data from embedded fonts for testing.
func getTestFontData(t *testing.T) []byte {
	t.Helper()

	// Load Roboto font data from embedded fonts
	compressed, err := fontFiles.ReadFile("fonts/Roboto-Medium.ttf.gz")
	require.NoError(t, err)

	decompressed, err := gzipDecompress(compressed)
	require.NoError(t, err)

	return decompressed
}

func getTestFont(t *testing.T) *truetype.Font {
	t.Helper()

	return GetFont("roboto")
}

func TestInstallFont(t *testing.T) {
	t.Parallel()

	t.Run("valid_font", func(t *testing.T) {
		fontData := getTestFontData(t)
		err := InstallFont("test-valid", fontData)
		require.NoError(t, err)

		font := GetFont("test-valid")
		assert.NotNil(t, font)
	})

	t.Run("invalid_font_data", func(t *testing.T) {
		err := InstallFont("test-invalid", []byte("invalid data"))
		assert.Error(t, err)
	})

	t.Run("case_insensitive_retrieval", func(t *testing.T) {
		fontData := getTestFontData(t)
		err := InstallFont("test-case", fontData)
		require.NoError(t, err)

		font1 := GetFont("test-case")
		font2 := GetFont("TEST-CASE")
		font3 := GetFont("Test-Case")

		assert.NotNil(t, font1)
		assert.Equal(t, font1, font2)
		assert.Equal(t, font1, font3)
	})
}

func TestGetFont(t *testing.T) {
	t.Parallel()

	t.Run("existing_font", func(t *testing.T) {
		fontName := "test-existing"
		fontData := getTestFontData(t)
		err := InstallFont(fontName, fontData)
		require.NoError(t, err)

		font := GetFont(fontName)
		assert.NotNil(t, font)
	})

	t.Run("nonexistent_font_returns_default", func(t *testing.T) {
		font := GetFont("nonexistent-font")
		assert.Nil(t, font)
	})

	t.Run("case_insensitive", func(t *testing.T) {
		fontName := "test-case-insensitive"
		fontData := getTestFontData(t)
		err := InstallFont(fontName, fontData)
		require.NoError(t, err)

		font1 := GetFont(fontName)
		font2 := GetFont("TEST-CASE-INSENSITIVE")
		font3 := GetFont("Test-Case-Insensitive")

		assert.Equal(t, font1, font2)
		assert.Equal(t, font1, font3)
	})
}

func TestLoadEmbeddedFont(t *testing.T) {
	t.Parallel()

	t.Run("nonexistent_embedded_font", func(t *testing.T) {
		err := loadEmbeddedFont("nonexistent-font")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "embedded font not found")
	})
}
