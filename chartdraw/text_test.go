package chartdraw

import (
	"fmt"
	"testing"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func TestTextWrapWord(t *testing.T) {
	t.Parallel()

	r := PNG(1024, 1024)

	basicTextStyle := Style{FontStyle: FontStyle{Font: GetDefaultFont(), FontSize: 24}}

	output := Text.WrapFitWord(r, "this is a test string", 100, basicTextStyle)
	assert.NotEmpty(t, output)
	require.Len(t, output, 3)

	for _, line := range output {
		basicTextStyle.WriteToRenderer(r)
		lineBox := r.MeasureText(line)
		assert.Less(t, lineBox.Width(), 100)
	}
	assert.Equal(t, "this is", output[0])
	assert.Equal(t, "a test", output[1])
	assert.Equal(t, "string", output[2])

	output = Text.WrapFitWord(r, "foo", 100, basicTextStyle)
	require.Len(t, output, 1)
	assert.Equal(t, "foo", output[0])

	// test that it handles newlines.
	output = Text.WrapFitWord(r, "this\nis\na\ntest\nstring", 100, basicTextStyle)
	assert.Len(t, output, 5)

	// test that it handles newlines and long lines.
	output = Text.WrapFitWord(r, "this\nis\na\ntest\nstring that is very long", 100, basicTextStyle)
	assert.Len(t, output, 8)
}

func TestTextWrapRune(t *testing.T) {
	t.Parallel()

	r := PNG(1024, 1024)

	basicTextStyle := Style{FontStyle: FontStyle{Font: GetDefaultFont(), FontSize: 24}}

	output := Text.WrapFitRune(r, "this is a test string", 150, basicTextStyle)
	assert.NotEmpty(t, output)
	require.Len(t, output, 2)
	assert.Equal(t, "this is a t", output[0])
	assert.Equal(t, "est string", output[1])
}

func TestExtents(t *testing.T) {
	t.Parallel()

	font, err := truetype.Parse(goregular.TTF)
	require.NoError(t, err)

	tests := []struct {
		size float64
	}{
		{8}, {16}, {32},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("size_%g", tt.size), func(t *testing.T) {
			got := drawing.Extents(font, tt.size)

			bounds := font.Bounds(fixed.Int26_6(font.FUnitsPerEm()))
			scale := tt.size / float64(font.FUnitsPerEm())
			want := drawing.FontExtents{
				Ascent:  float64(bounds.Max.Y) * scale,
				Descent: float64(bounds.Min.Y) * scale,
				Height:  float64(bounds.Max.Y-bounds.Min.Y) * scale,
			}

			assert.InDelta(t, want.Ascent, got.Ascent, 1e-5)
			assert.InDelta(t, want.Descent, got.Descent, 1e-5)
			assert.InDelta(t, want.Height, got.Height, 1e-5)
		})
	}
}
