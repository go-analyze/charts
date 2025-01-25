package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkParseColor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ParseColor("#333")
		_ = ParseColor("#313233")
		_ = ParseColor("rgb(31,32,33)")
		_ = ParseColor("rgba(50,51,52,250)")
	}
}

func BenchmarkColorString(b *testing.B) {
	c := ParseColor("rgb(31,32,33)")
	for i := 0; i < b.N; i++ {
		_ = c.String()
	}
}

const makeColorShadeSamples = false

func testColorShades(t *testing.T, colors ...Color) {
	if !makeColorShadeSamples {
		return // color samples are generated through a test failure
	}
	t.Helper()

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	})

	sampleWidth := p.Width() / len(colors)
	for i, c := range colors {
		p.FilledRect(i*sampleWidth, 0, (i+1)*sampleWidth, p.Height(),
			c, ColorTransparent, 0.0)
	}

	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "", data)
}

func TestGrayColors(t *testing.T) {
	testColorShades(t, ColorDarkGray, ColorGray, ColorLightGray)
}

func TestBlueColors(t *testing.T) {
	testColorShades(t, ColorBlue, ColorBlueAlt1, ColorBlueAlt2)
}

func TestGreenColors(t *testing.T) {
	testColorShades(t, ColorGreen, ColorGreenAlt1, ColorGreenAlt2, ColorGreenAlt3, ColorGreenAlt4)
}

func TestRedColors(t *testing.T) {
	testColorShades(t, ColorRed, ColorRedAlt1, ColorRedAlt2)
}

func TestOrangeColors(t *testing.T) {
	testColorShades(t, ColorOrange, ColorOrangeAlt1, ColorOrangeAlt2, ColorOrangeAlt3)
}

func TestAquaColors(t *testing.T) {
	testColorShades(t, ColorAqua, ColorAquaAlt1)
}

func TestYellowColors(t *testing.T) {
	testColorShades(t, ColorYellow, ColorYellowAlt1)
}

func TestPurpleColors(t *testing.T) {
	testColorShades(t, ColorPurple, ColorViolet, ColorPlum, ColorFuchsia)
}

func TestIsLightColor(t *testing.T) {
	t.Parallel()

	assert.True(t, isLightColor(Color{R: 255, G: 255, B: 255}))
	assert.True(t, isLightColor(Color{R: 145, G: 204, B: 117}))

	assert.False(t, isLightColor(Color{R: 88, G: 112, B: 198}))
	assert.False(t, isLightColor(Color{R: 0, G: 0, B: 0}))
	assert.False(t, isLightColor(Color{R: 16, G: 12, B: 42}))
}

func TestParseColor(t *testing.T) {
	t.Parallel()

	c := ParseColor("")
	assert.True(t, c.IsZero())

	c = ParseColor("#333")
	assert.Equal(t, Color{R: 51, G: 51, B: 51, A: 255}, c)

	c = ParseColor("#313233")
	assert.Equal(t, Color{R: 49, G: 50, B: 51, A: 255}, c)

	c = ParseColor("rgb(31,32,33)")
	assert.Equal(t, Color{R: 31, G: 32, B: 33, A: 255}, c)

	c = ParseColor("rgba(50,51,52,.981)")
	assert.Equal(t, Color{R: 50, G: 51, B: 52, A: 250}, c)

	c = ParseColor("rgba(50,51,52,250)")
	assert.Equal(t, Color{R: 50, G: 51, B: 52, A: 250}, c)
}
