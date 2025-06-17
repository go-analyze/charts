package drawing

import (
	"image"
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/roboto"
)

func TestRasterGraphicContext(t *testing.T) {
	t.Run("matrix operations", func(t *testing.T) {
		t.Parallel()

		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		ctx := NewRasterGraphicContext(img)

		originalMatrix := ctx.GetMatrixTransform()
		assert.NotNil(t, originalMatrix)

		identityMatrix := NewIdentityMatrix()
		ctx.SetMatrixTransform(identityMatrix)
		currentMatrix := ctx.GetMatrixTransform()
		assert.Equal(t, identityMatrix, currentMatrix)

		ctx.Translate(10, 20)
		ctx.Scale(2, 3)
		ctx.Rotate(0.5)

		// The matrix should have changed
		transformedMatrix := ctx.GetMatrixTransform()
		assert.NotEqual(t, identityMatrix, transformedMatrix)
	})

	t.Run("font operations", func(t *testing.T) {
		t.Parallel()

		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		ctx := NewRasterGraphicContext(img)

		ctx.SetFontSize(12)
		fontSize := ctx.GetFontSize()
		assert.InDelta(t, 12.0, fontSize, 0.001)

		ctx.SetFontSize(24)
		fontSize = ctx.GetFontSize()
		assert.InDelta(t, 24.0, fontSize, 0.001)

		originalFont := ctx.GetFont()
		ctx.SetFont(nil)
		currentFont := ctx.GetFont()
		assert.Nil(t, currentFont)

		// Restore original font if it existed
		if originalFont != nil {
			ctx.SetFont(originalFont)
			restoredFont := ctx.GetFont()
			assert.Equal(t, originalFont, restoredFont)
		}
	})

	t.Run("DPI operations", func(t *testing.T) {
		t.Parallel()

		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		ctx := NewRasterGraphicContext(img)

		ctx.SetDPI(72.0)
		dpi := ctx.GetDPI()
		assert.InDelta(t, 72.0, dpi, 0.001)

		ctx.SetDPI(300.0)
		dpi = ctx.GetDPI()
		assert.InDelta(t, 300.0, dpi, 0.001)
	})

	t.Run("save and restore", func(t *testing.T) {
		t.Parallel()

		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		ctx := NewRasterGraphicContext(img)

		// Set initial state
		ctx.SetLineWidth(5)
		ctx.SetFontSize(16)
		ctx.Translate(10, 20)

		// Capture the state before saving
		expectedMatrix := ctx.GetMatrixTransform()
		expectedFontSize := ctx.GetFontSize()

		ctx.Save()

		// Modify state
		ctx.SetLineWidth(10)
		ctx.SetFontSize(32)
		ctx.Translate(30, 40)

		// Verify state was modified
		modifiedMatrix := ctx.GetMatrixTransform()
		modifiedFontSize := ctx.GetFontSize()
		assert.NotEqual(t, expectedMatrix, modifiedMatrix)
		assert.NotEqual(t, expectedFontSize, modifiedFontSize)

		// Restore the state
		ctx.Restore()

		// Validate the restore state - should exactly match the saved state
		restoredMatrix := ctx.GetMatrixTransform()
		restoredFontSize := ctx.GetFontSize()
		assert.Equal(t, expectedMatrix, restoredMatrix)
		assert.InDelta(t, expectedFontSize, restoredFontSize, 0.001)
	})

	t.Run("text operations", func(t *testing.T) {
		t.Parallel()

		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		ctx := NewRasterGraphicContext(img)
		font, err := truetype.Parse(roboto.Roboto)
		require.NoError(t, err)
		ctx.SetFont(font)
		ctx.SetFontSize(12)

		left, top, right, bottom, err := ctx.GetStringBounds("Hello")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, right, left)
		assert.GreaterOrEqual(t, bottom, top)

		cursor, err := ctx.CreateStringPath("Test", 10, 20)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, cursor, 0.0)

		cursor, err = ctx.FillStringAt("World", 30, 40)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, cursor, 0.0)

		cursor, err = ctx.StrokeStringAt("StrokeAt", 50, 60)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, cursor, 0.0)
	})
}
