package chartdraw

import (
	"bytes"
	"compress/gzip"
	"embed"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed drawing/fonts/Roboto-Medium.ttf.gz
var fontFiles embed.FS

// getTestFontData loads font data from embedded fonts for testing.
func getTestFontData(t *testing.T) []byte {
	t.Helper()

	// Load Roboto font data from embedded fonts
	compressed, err := fontFiles.ReadFile("drawing/fonts/Roboto-Medium.ttf.gz")
	require.NoError(t, err)

	r, err := gzip.NewReader(bytes.NewReader(compressed))
	require.NoError(t, err)
	defer func() { _ = r.Close() }()

	decompressed, err := io.ReadAll(r)
	require.NoError(t, err)
	return decompressed
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
		assert.Equal(t, GetDefaultFont(), font)
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

func TestGetDefaultFont(t *testing.T) {
	t.Parallel()

	defaultFont := GetDefaultFont()
	assert.NotNil(t, defaultFont)
}

func TestSetDefaultFont(t *testing.T) {
	t.Run("set_existing_font", func(t *testing.T) {
		t.Cleanup(func() {
			_ = SetDefaultFont("roboto")
		})
		fontName := "test-set-default"
		fontData := getTestFontData(t)
		err := InstallFont(fontName, fontData)
		require.NoError(t, err)

		err = SetDefaultFont(fontName)
		require.NoError(t, err)

		defaultFont := GetDefaultFont()
		expectedFont := GetFont(fontName)
		assert.Equal(t, expectedFont, defaultFont)
	})

	t.Run("set_nonexistent_font", func(t *testing.T) {
		t.Cleanup(func() {
			_ = SetDefaultFont("roboto")
		})
		err := SetDefaultFont("nonexistent-font")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "font not found")
	})

	t.Run("case_insensitive", func(t *testing.T) {
		t.Cleanup(func() {
			_ = SetDefaultFont("roboto")
		})
		fontName := "test-default-case"
		fontData := getTestFontData(t)
		err := InstallFont(fontName, fontData)
		require.NoError(t, err)

		err = SetDefaultFont("TEST-DEFAULT-CASE")
		require.NoError(t, err)
	})
}

func TestFontInitialization(t *testing.T) {
	defaultFamily := defaultFontFamily.Load()
	assert.NotNil(t, defaultFamily)
	assert.Equal(t, "roboto", defaultFamily)

	defaultFont := GetDefaultFont()
	assert.NotNil(t, defaultFont)
}
