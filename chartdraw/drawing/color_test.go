package drawing

import (
	"image/color"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorFromHex(t *testing.T) {
	t.Parallel()

	white := ColorFromHex("FFFFFF")
	assert.Equal(t, ColorWhite, white)

	shortWhite := ColorFromHex("FFF")
	assert.Equal(t, ColorWhite, shortWhite)

	black := ColorFromHex("000000")
	assert.Equal(t, ColorBlack, black)

	shortBlack := ColorFromHex("000")
	assert.Equal(t, ColorBlack, shortBlack)

	red := ColorFromHex("FF0000")
	assert.Equal(t, ColorRed, red)

	shortRed := ColorFromHex("F00")
	assert.Equal(t, ColorRed, shortRed)

	green := ColorFromHex("008000")
	assert.Equal(t, ColorGreen, green)

	// shortGreen := ColorFromHex("0F0")
	// assert.Equal(t, ColorGreen, shortGreen)

	blue := ColorFromHex("0000FF")
	assert.Equal(t, ColorBlue, blue)

	shortBlue := ColorFromHex("00F")
	assert.Equal(t, ColorBlue, shortBlue)
}

func TestColorFromHex_handlesHash(t *testing.T) {
	t.Parallel()

	withHash := ColorFromHex("#FF0000")
	assert.Equal(t, ColorRed, withHash)

	withoutHash := ColorFromHex("#FF0000")
	assert.Equal(t, ColorRed, withoutHash)
}

func TestColorFromAlphaMixedRGBA(t *testing.T) {
	t.Parallel()

	black := ColorFromAlphaMixedRGBA(color.Black.RGBA())
	assert.True(t, black.Equals(ColorBlack), black.String())

	white := ColorFromAlphaMixedRGBA(color.White.RGBA())
	assert.True(t, white.Equals(ColorWhite), white.String())
}

func Test_ColorFromRGBA(t *testing.T) {
	t.Parallel()

	value := "rgba(192, 192, 192, 1.0)"
	parsed := ColorFromRGBA(value)
	assert.Equal(t, ColorSilver, parsed)

	value = "rgba(192,192,192,1.0)"
	parsed = ColorFromRGBA(value)
	assert.Equal(t, ColorSilver, parsed)

	value = "rgba(192,192,192,255)"
	parsed = ColorFromRGBA(value)
	assert.Equal(t, ColorSilver, parsed)
}

func TestParseColor(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		Input    string
		Expected Color
	}{
		{"", Color{}},
		{"unknown", ColorBlack},
		{"white", ColorWhite},
		{"WHITE", ColorWhite}, // caps!
		{"black", ColorBlack},
		{"red", ColorRed},
		{"gray", ColorGray},
		{"grey", ColorGray},
		{"green", ColorGreen},
		{"blue", ColorBlue},
		{"silver", ColorSilver},
		{"maroon", ColorMaroon},
		{"purple", ColorPurple},
		{"fuchsia", ColorFuchsia},
		{"lime", ColorLime},
		{"olive", ColorOlive},
		{"yellow", ColorYellow},
		{"navy", ColorNavy},
		{"teal", ColorTeal},
		{"aqua", ColorAqua},

		{"rgba(192, 192, 192, 1.0)", ColorSilver},
		{"rgba(192,192,192,1.0)", ColorSilver},
		{"rgb(192, 192, 192)", ColorSilver},
		{"rgb(192,192,192)", ColorSilver},

		{"#FF0000", ColorRed},
		{"#008000", ColorGreen},
		{"#0000FF", ColorBlue},
		{"#F00", ColorRed},
		{"#080", Color{0, 136, 0, 255}},
		{"#00F", ColorBlue},
	}

	for index, tc := range testCases {
		t.Run(strconv.Itoa(index), func(t *testing.T) {
			actual := ParseColor(tc.Input)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}
func TestColorHelperMethods(t *testing.T) {
	t.Parallel()

	chTests := []struct {
		f      float64
		expect uint8
	}{
		{-0.1, 0},
		{0.5, 127},
		{1.5, 255},
	}
	for i, tc := range chTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.expect, ColorChannelFromFloat(tc.f))
		})
	}

	c := Color{R: 10, G: 20, B: 30, A: 255}
	r, g, b, a := c.RGBA()
	assert.Equal(t, uint32(2570), r)
	assert.Equal(t, uint32(5140), g)
	assert.Equal(t, uint32(7710), b)
	assert.Equal(t, uint32(65535), a)

	zero := Color{}
	assert.True(t, zero.IsZero())
	assert.True(t, zero.IsTransparent())
	assert.False(t, c.IsZero())

	withAlpha := c.WithAlpha(128)
	assert.Equal(t, uint8(128), withAlpha.A)

	avg := ColorRed.AverageWith(ColorBlue)
	assert.Equal(t, Color{R: 127, G: 0, B: 127, A: 255}, avg)

	assert.Equal(t, "rgb(10,20,30)", c.StringRGB())
	assert.Equal(t, "rgba(10,20,30,0.5)", c.WithAlpha(128).StringRGBA())
}

func TestColorHSLConversions(t *testing.T) {
	t.Parallel()

	h, s, l := ColorRed.HSL()
	assert.InDelta(t, 0.0, h, 0.001)
	assert.InDelta(t, 1.0, s, 0.001)
	assert.InDelta(t, 0.5, l, 0.001)

	r, g, b := hslToRGB(h, s, l)
	assert.Equal(t, ColorRed.R, r)
	assert.Equal(t, ColorRed.G, g)
	assert.Equal(t, ColorRed.B, b)

	adjusted := ColorRed.WithAdjustHSL(120, 0, 0)
	assert.Equal(t, ColorLime.R, adjusted.R)
	assert.Equal(t, ColorLime.G, adjusted.G)
	assert.Equal(t, ColorLime.B, adjusted.B)

	assert.InDelta(t, 0.25, clamp(0.25, 0, 1), 0.0)
	assert.InDelta(t, 0.0, clamp(-0.5, 0, 1), 0.0)
	assert.InDelta(t, 1.0, clamp(2, 0, 1), 0.0)

	assert.InDelta(t, 0.0, hue2rgb(0, 1, 0), 0.0001)
	assert.InDelta(t, 1.0, hue2rgb(0, 1, 1.0/6.0), 0.0001)
}
