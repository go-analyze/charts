package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/colornames"
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
	t.Helper()
	if !makeColorShadeSamples {
		return // color samples are generated through a test failure
	}

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	})

	sampleWidth := p.Width() / len(colors)
	for i, c := range colors {
		endX := (i + 1) * sampleWidth
		if i == len(colors)-1 {
			endX = p.Width() // ensure edge is painted
		}
		p.FilledRect(i*sampleWidth, 0, endX, p.Height(),
			c, ColorTransparent, 0.0)
	}

	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, nil, data)
}

func TestGrayColors(t *testing.T) {
	testColorShades(t, ColorDarkGray, ColorGray, ColorSilver, ColorLightGray, ColorAzure,
		ColorSlateGray, ColorLightSlateGray)
}

func TestBlueColors(t *testing.T) {
	testColorShades(t, ColorSkyBlue, ColorBlue, ColorNavy, ColorBlueAlt1, ColorBlueAlt2, ColorBlueAlt3, ColorLightSlateGray)
}

func TestGreenColors(t *testing.T) {
	testColorShades(t, ColorGreen, ColorOlive, ColorLime, ColorSageGreen,
		ColorGreenAlt1, ColorGreenAlt2, ColorGreenAlt3, ColorGreenAlt4, ColorGreenAlt5, ColorGreenAlt6, ColorGreenAlt7)
}

func TestRedColors(t *testing.T) {
	testColorShades(t, ColorLightCoral, ColorCoral, ColorRed, ColorPink, ColorSalmon, ColorMaroon, ColorBrown, ColorChocolate,
		ColorRedAlt1, ColorRedAlt2, ColorRedAlt3, ColorRedAlt4)
}

func TestOrangeColors(t *testing.T) {
	testColorShades(t, ColorOrange, ColorOrangeAlt1, ColorOrangeAlt2, ColorOrangeAlt3, ColorOrangeAlt4)
}

func TestAquaColors(t *testing.T) {
	testColorShades(t, ColorAqua, ColorTeal, ColorTurquoise, ColorAquaAlt1, ColorAquaAlt2)
}

func TestYellowColors(t *testing.T) {
	testColorShades(t, ColorYellow, ColorGold, ColorYellowAlt1, ColorMustardYellow)
}

func TestTanColors(t *testing.T) {
	testColorShades(t, ColorAzure, ColorIvory, ColorBeige, ColorKhaki, ColorTan, ColorCoral, ColorSalmon, ColorLightCoral)
}

func TestPurpleColors(t *testing.T) {
	testColorShades(t, ColorLavender, ColorThistle, ColorPurple, ColorViolet, ColorIndigo, ColorPlum, ColorFuchsia,
		ColorPurpleAlt1, ColorPurpleAlt2)
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

	c = ParseColor("unknown")
	assert.Equal(t, ColorBlack, c)

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

func TestColorConvertGo(t *testing.T) {
	t.Parallel()

	goC := colornames.Lavender
	ourC := ColorConvertGo(goC)

	goR, goG, goB, goA := goC.RGBA()
	ourR, ourG, ourB, ourA := ourC.RGBA()

	assert.Equal(t, goR, ourR)
	assert.Equal(t, goG, ourG)
	assert.Equal(t, goB, ourB)
	assert.Equal(t, goA, ourA)
}
