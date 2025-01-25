package charts

import (
	"math"
	"strings"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

var (
	// ColorTransparent is a fully transparent color.
	ColorTransparent = drawing.ColorTransparent
	// ColorWhite is R: 255, G: 255, B: 255.
	ColorWhite = drawing.ColorWhite
	// ColorBlack is R: 0, G: 0, B: 0.
	ColorBlack = drawing.ColorBlack
	// ColorGray is R: 128, G: 128, B: 128.
	ColorGray = drawing.ColorGray
	// ColorRed is R: 255, G: 0, B: 0.
	ColorRed = drawing.ColorRed
	// ColorGreen is R: 0, G: 128, B: 0.
	ColorGreen = drawing.ColorGreen
	// ColorBlue is R: 0, G: 0, B: 255.
	ColorBlue = drawing.ColorBlue
	// ColorSilver is R: 192, G: 192, B: 192.
	ColorSilver = drawing.ColorSilver
	// ColorMaroon is R: 128, G: 0, B: 0.
	ColorMaroon = drawing.ColorMaroon
	// ColorPurple is R: 128, G: 0, B: 128.
	ColorPurple = drawing.ColorPurple
	// ColorFuchsia is R: 255, G: 0, B: 255.
	ColorFuchsia = drawing.ColorFuchsia
	// ColorLime is R: 0, G: 255, B: 0.
	ColorLime = drawing.ColorLime
	// ColorOlive is R: 128, G: 128, B: 0.
	ColorOlive = drawing.ColorOlive
	// ColorYellow is R: 255, G: 255, B: 0.
	ColorYellow = drawing.ColorYellow
	// ColorNavy is R: 0, G: 0, B: 128.
	ColorNavy = drawing.ColorNavy
	// ColorTeal is R: 0, G: 128, B: 128.
	ColorTeal = drawing.ColorTeal
	// ColorAqua (or Cyan) is R: 0, G: 255, B: 255.
	ColorAqua = drawing.ColorAqua
	// ColorDarkGray is R: 40, G: 40, B: 40.
	ColorDarkGray = Color{R: 40, G: 40, B: 40, A: 255}
	// ColorLightGray is R: 211, G: 211, B: 211.
	ColorLightGray = drawing.ColorLightGray
	// ColorSlateGray is R: 112, G: 128, B: 144.
	ColorSlateGray = drawing.ColorSlateGray
	// ColorLightSlateGray is R: 119, G: 136, B: 211.
	ColorLightSlateGray = drawing.ColorLightSlateGray
	// ColorAzure is R: 240, G: 255, B: 255.
	ColorAzure = drawing.ColorAzure
	// ColorBeige is R: 245, G: 245, B: 220.
	ColorBeige = drawing.ColorBeige
	// ColorBrown is R: 165, G: 42, B: 42.
	ColorBrown = drawing.ColorBrown
	// ColorChocolate is R: 210, G: 105, B: 30.
	ColorChocolate = drawing.ColorChocolate
	// ColorCoral is R: 255, G: 127, B: 80.
	ColorCoral = drawing.ColorCoral
	// ColorGold is R: 255, G: 215, B: 0.
	ColorGold = drawing.ColorGold
	// ColorIndigo is R: 75, G: 0, B: 130.
	ColorIndigo = drawing.ColorIndigo
	// ColorIvory is R: 255, G: 255, B: 250.
	ColorIvory = drawing.ColorIvory
	// ColorOrange is R: 255, G: 165, B: 0.
	ColorOrange = drawing.ColorOrange
	// ColorPink is R: 255, G: 192, B: 203.
	ColorPink = drawing.ColorPink
	// ColorPlum is R: 221, G: 160, B: 221.
	ColorPlum = drawing.ColorPlum
	// ColorSalmon is R: 250, G: 128, B: 114.
	ColorSalmon = drawing.ColorSalmon
	// ColorTan is R: 210, G: 180, B: 140.
	ColorTan = drawing.ColorTan
	// ColorKhaki is R: 240, G: 230, B: 140.
	ColorKhaki = drawing.ColorKhaki
	// ColorTurquoise is R: 64, G: 224, B: 208.
	ColorTurquoise = drawing.ColorTurquoise
	// ColorViolet is R: 238, G: 130, B: 238.
	ColorViolet = drawing.ColorViolet

	// alternate non-standard shades //

	// ColorBlackAlt1 is slightly lighter shade of black: R: 51, G: 51, B: 51.
	ColorBlackAlt1 = chartdraw.ColorBlack
	// ColorBlueAlt1 is lighter shade of blue: R:0, G: 116, B: 217.
	ColorBlueAlt1 = chartdraw.ColorBlue
	// ColorBlueAlt2 is a sea blue: R: 106, G: 195, B: 203.
	ColorBlueAlt2 = chartdraw.ColorAlternateBlue
	// ColorAquaAlt1 is a lighter aqua: R: 0, G: 217, B: 210.
	ColorAquaAlt1 = chartdraw.ColorCyan
	// ColorGreenAlt1 is lighter green: R: 0, G: 217, B: 101.
	ColorGreenAlt1 = chartdraw.ColorGreen
	// ColorGreenAlt2 is R: 42, G: 190, B: 137.
	ColorGreenAlt2 = chartdraw.ColorAlternateGreen
	// ColorGreenAlt3 is darker green: R: 59, G: 162, B: 114.
	ColorGreenAlt3 = Color{R: 59, G: 162, B: 114, A: 255}
	// ColorGreenAlt4 is darker green: R: 80, G: 134, B: 66.
	ColorGreenAlt4 = Color{R: 80, G: 143, B: 66, A: 255}
	// ColorRedAlt1 is slightly purple red: R: 217, G: 0, B: 116.
	ColorRedAlt1 = chartdraw.ColorRed
	// ColorRedAlt2 is darker purple red: R: 226, G: 77, B: 66.
	ColorRedAlt2 = Color{R: 226, G: 77, B: 66, A: 255}
	// ColorOrangeAlt1 is more typical orange: R: 217, G: 101, B: 0.
	ColorOrangeAlt1 = chartdraw.ColorOrange
	// ColorOrangeAlt2 is a lighter orange: R: 250, G: 200, B: 88.
	ColorOrangeAlt2 = Color{R: 250, G: 200, B: 88, A: 255}
	// ColorOrangeAlt3 is a lighter orange: R: 255, G: 152, B: 69.
	ColorOrangeAlt3 = Color{R: 255, G: 152, B: 69, A: 255}
	// ColorYellowAlt1 is a slightly darker yellow: R: 217, G: 210, B: 0.
	ColorYellowAlt1 = chartdraw.ColorYellow
)

func isLightColor(c Color) bool {
	r := float64(c.R) * float64(c.R) * 0.299
	g := float64(c.G) * float64(c.G) * 0.587
	b := float64(c.B) * float64(c.B) * 0.114
	return math.Sqrt(r+g+b) > 127.5
}

// ParseColor parses a color from a string. The color can be specified in hex with a `#` prefix (for example '#313233'),
// in rgb(i,i,i) or rgba(i,i,i,f) format, or as a common name (for example 'red').
func ParseColor(rawColor string) Color {
	if strings.HasPrefix(rawColor, "#") {
		return ColorFromHex(rawColor)
	} else if strings.HasPrefix(rawColor, "rgb") {
		return ColorFromRGBA(rawColor)
	}
	return ColorFromKnown(rawColor)
}

// ColorFromKnown returns an internal color from a known (basic) color name.
func ColorFromKnown(known string) Color {
	return drawing.ColorFromKnown(known)
}

// ColorFromHex returns a color from a css hex code.
//
// NOTE: it will trim a leading '#' character if present.
func ColorFromHex(hex string) Color {
	return drawing.ColorFromHex(hex)
}

// ColorFromRGBA returns a color from a `rgb(i,i,i)` or `rgba(i,i,i,f)` css function.
func ColorFromRGBA(color string) Color {
	return drawing.ColorFromRGBA(color)
}

// ColorFromRGBAValues returns the system alpha mixed rgba values.
func ColorFromRGBAValues(r, g, b, a uint32) Color {
	return drawing.ColorFromAlphaMixedRGBA(r, g, b, a)
}
