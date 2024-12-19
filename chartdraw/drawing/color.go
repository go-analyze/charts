package drawing

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Basic Colors from:
// https://www.w3.org/wiki/CSS/Properties/color/keywords
var (
	// ColorTransparent is a fully transparent color.
	ColorTransparent = Color{R: 255, G: 255, B: 255, A: 0}
	// ColorWhite is white.
	ColorWhite = Color{R: 255, G: 255, B: 255, A: 255}
	// ColorBlack is black.
	ColorBlack = Color{R: 0, G: 0, B: 0, A: 255}
	// ColorRed is red.
	ColorRed = Color{R: 255, G: 0, B: 0, A: 255}
	// ColorGreen is green.
	ColorGreen = Color{R: 0, G: 128, B: 0, A: 255}
	// ColorBlue is blue.
	ColorBlue = Color{R: 0, G: 0, B: 255, A: 255}
	// ColorSilver is a known color.
	ColorSilver = Color{R: 192, G: 192, B: 192, A: 255}
	// ColorMaroon is a known color.
	ColorMaroon = Color{R: 128, G: 0, B: 0, A: 255}
	// ColorPurple is a known color.
	ColorPurple = Color{R: 128, G: 0, B: 128, A: 255}
	// ColorFuchsia is a known color.
	ColorFuchsia = Color{R: 255, G: 0, B: 255, A: 255}
	// ColorLime is a known color.
	ColorLime = Color{R: 0, G: 255, B: 0, A: 255}
	// ColorOlive is a known color.
	ColorOlive = Color{R: 128, G: 128, B: 0, A: 255}
	// ColorYellow is a known color.
	ColorYellow = Color{R: 255, G: 255, B: 0, A: 255}
	// ColorNavy is a known color.
	ColorNavy = Color{R: 0, G: 0, B: 128, A: 255}
	// ColorTeal is a known color.
	ColorTeal = Color{R: 0, G: 128, B: 128, A: 255}
	// ColorAqua is a known color.
	ColorAqua = Color{R: 0, G: 255, B: 255, A: 255}
)

// ParseColor parses a color from a string.
func ParseColor(rawColor string) Color {
	if strings.HasPrefix(rawColor, "#") {
		return ColorFromHex(rawColor)
	} else if strings.HasPrefix(rawColor, "rgb") {
		return ColorFromRGBA(rawColor)
	}
	return ColorFromKnown(rawColor)
}

var rgbReg = regexp.MustCompile(`\(([^)]+)\)`)

// ColorFromRGBA returns a color from a `rgb(i,i,i)` or `rgba(i,i,i,f)` css function.
func ColorFromRGBA(color string) Color {
	var c Color

	// Attempt to parse rgb(...) or rgba(...)
	result := rgbReg.FindAllStringSubmatch(color, 1)
	if len(result) == 0 || len(result[0]) != 2 {
		return c
	}
	arr := strings.Split(result[0][1], ",")
	if len(arr) < 3 { // at a minimum we expect r,g,b to be specified
		return c
	}

	rVal, _ := strconv.ParseInt(strings.TrimSpace(arr[0]), 10, 16)
	c.R = uint8(rVal)
	gVal, _ := strconv.ParseInt(strings.TrimSpace(arr[1]), 10, 16)
	c.G = uint8(gVal)
	bVal, _ := strconv.ParseInt(strings.TrimSpace(arr[2]), 10, 16)
	c.B = uint8(bVal)
	if len(arr) > 3 { // if alpha channel is specified
		aVal, _ := strconv.ParseFloat(strings.TrimSpace(arr[3]), 64)
		if aVal < 0 {
			aVal = 0
		} else if aVal <= 1 {
			// correctly specified decimal, convert it to an integer scale
			aVal *= 255
		} // else, incorrectly specified value over 1, accept the value directly
		c.A = uint8(aVal)
	} else {
		c.A = 255 // default alpha channel to 255
	}

	return c
}

// Deprecated: ColorFromRGB is deprecated, use ColorFromRGBA to get colors from RGB or RGBA format strings.
func ColorFromRGB(rgb string) (output Color) {
	return ColorFromRGBA(rgb)
}

func parseHex(hex string) uint8 {
	v, _ := strconv.ParseInt(hex, 16, 16)
	return uint8(v)
}

// ColorFromHex returns a color from a css hex code.
//
// NOTE: it will trim a leading '#' character if present.
func ColorFromHex(hex string) Color {
	hex = strings.TrimPrefix(hex, "#")
	c := Color{A: 255}
	if len(hex) == 3 {
		c.R = parseHex(string(hex[0])) * 0x11
		c.G = parseHex(string(hex[1])) * 0x11
		c.B = parseHex(string(hex[2])) * 0x11
	} else {
		c.R = parseHex(hex[0:2])
		c.G = parseHex(hex[2:4])
		c.B = parseHex(hex[4:6])
	}
	return c
}

// ColorFromKnown returns an internal color from a known (basic) color name.
func ColorFromKnown(known string) Color {
	switch strings.ToLower(known) {
	case "transparent":
		return ColorTransparent
	case "white":
		return ColorWhite
	case "black":
		return ColorBlack
	case "red":
		return ColorRed
	case "blue":
		return ColorBlue
	case "green":
		return ColorGreen
	case "silver":
		return ColorSilver
	case "maroon":
		return ColorMaroon
	case "purple":
		return ColorPurple
	case "fuchsia":
		return ColorFuchsia
	case "lime":
		return ColorLime
	case "olive":
		return ColorOlive
	case "yellow":
		return ColorYellow
	case "navy":
		return ColorNavy
	case "teal":
		return ColorTeal
	case "aqua":
		return ColorAqua
	default:
		return Color{}
	}
}

// ColorFromAlphaMixedRGBA returns the system alpha mixed rgba values.
func ColorFromAlphaMixedRGBA(r, g, b, a uint32) Color {
	fa := float64(a) / 255.0
	return Color{
		R: uint8(float64(r) / fa),
		G: uint8(float64(g) / fa),
		B: uint8(float64(b) / fa),
		A: uint8(a | (a >> 8)),
	}
}

// ColorChannelFromFloat returns a normalized byte from a given float value.
func ColorChannelFromFloat(v float64) uint8 {
	return uint8(v * 255)
}

// Color is our internal color type because color.Color is bullshit.
type Color struct {
	R, G, B, A uint8
}

// RGBA returns the color as a pre-alpha mixed color set.
func (c Color) RGBA() (r, g, b, a uint32) {
	fa := float64(c.A) / 255.0
	r = uint32(float64(c.R) * fa)
	r |= r << 8
	g = uint32(float64(c.G) * fa)
	g |= g << 8
	b = uint32(float64(c.B) * fa)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

// IsZero returns if the color has been set or not.
func (c Color) IsZero() bool {
	return c.R == 0 && c.G == 0 && c.B == 0 && c.A == 0
}

// IsTransparent returns if the colors alpha channel is zero.
func (c Color) IsTransparent() bool {
	return c.A == 0
}

// WithAlpha returns a copy of the color with a given alpha.
func (c Color) WithAlpha(a uint8) Color {
	return Color{
		R: c.R,
		G: c.G,
		B: c.B,
		A: a,
	}
}

// Equals returns true if the color equals another.
func (c Color) Equals(other Color) bool {
	return c.R == other.R &&
		c.G == other.G &&
		c.B == other.B &&
		c.A == other.A
}

// AverageWith averages two colors.
func (c Color) AverageWith(other Color) Color {
	return Color{
		R: (c.R + other.R) >> 1,
		G: (c.G + other.G) >> 1,
		B: (c.B + other.B) >> 1,
		A: c.A,
	}
}

// String returns a css string representation of the color.
func (c Color) String() string {
	fa := float64(c.A) / float64(255)
	return fmt.Sprintf("rgba(%v,%v,%v,%.1f)", c.R, c.G, c.B, fa)
}
