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
		{"white", ColorWhite},
		{"WHITE", ColorWhite}, // caps!
		{"black", ColorBlack},
		{"red", ColorRed},
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
