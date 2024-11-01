package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTextWrapWord(t *testing.T) {
	r, err := PNG(1024, 1024)
	require.NoError(t, err)
	f, err := GetDefaultFont()
	require.NoError(t, err)

	basicTextStyle := Style{FontStyle: FontStyle{Font: f, FontSize: 24}}

	output := Text.WrapFitWord(r, "this is a test string", 100, basicTextStyle)
	assert.NotEmpty(t, output)
	assert.Len(t, output, 3)

	for _, line := range output {
		basicTextStyle.WriteToRenderer(r)
		lineBox := r.MeasureText(line)
		assert.True(t, lineBox.Width() < 100, line+": "+lineBox.String())
	}
	assert.Equal(t, "this is", output[0])
	assert.Equal(t, "a test", output[1])
	assert.Equal(t, "string", output[2])

	output = Text.WrapFitWord(r, "foo", 100, basicTextStyle)
	assert.Len(t, output, 1)
	assert.Equal(t, "foo", output[0])

	// test that it handles newlines.
	output = Text.WrapFitWord(r, "this\nis\na\ntest\nstring", 100, basicTextStyle)
	assert.Len(t, output, 5)

	// test that it handles newlines and long lines.
	output = Text.WrapFitWord(r, "this\nis\na\ntest\nstring that is very long", 100, basicTextStyle)
	assert.Len(t, output, 8)
}

func TestTextWrapRune(t *testing.T) {
	r, err := PNG(1024, 1024)
	require.NoError(t, err)
	f, err := GetDefaultFont()
	require.NoError(t, err)

	basicTextStyle := Style{FontStyle: FontStyle{Font: f, FontSize: 24}}

	output := Text.WrapFitRune(r, "this is a test string", 150, basicTextStyle)
	assert.NotEmpty(t, output)
	assert.Len(t, output, 2)
	assert.Equal(t, "this is a t", output[0])
	assert.Equal(t, "est string", output[1])
}
