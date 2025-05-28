package drawing

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-analyze/charts/chartdraw/matrix"
)

// Basic Colors from:
// https://www.w3.org/wiki/CSS/Properties/color/keywords
var (
	// ColorTransparent is a fully transparent color.
	ColorTransparent = Color{R: 255, G: 255, B: 255, A: 0}
	// ColorWhite is R: 255, G: 255, B: 255.
	ColorWhite = Color{R: 255, G: 255, B: 255, A: 255}
	// ColorBlack is R: 0, G: 0, B: 0.
	ColorBlack = Color{R: 0, G: 0, B: 0, A: 255}
	// ColorGray is R: 128, G: 128, B: 128.
	ColorGray = Color{R: 128, G: 128, B: 128, A: 255}
	// ColorRed is R: 255, G: 0, B: 0.
	ColorRed = Color{R: 255, G: 0, B: 0, A: 255}
	// ColorGreen is R: 0, G: 128, B: 0.
	ColorGreen = Color{R: 0, G: 128, B: 0, A: 255}
	// ColorBlue is R: 0, G: 0, B: 255.
	ColorBlue = Color{R: 0, G: 0, B: 255, A: 255}
	// ColorSilver is R: 192, G: 192, B: 192.
	ColorSilver = Color{R: 192, G: 192, B: 192, A: 255}
	// ColorMaroon is R: 128, G: 0, B: 0.
	ColorMaroon = Color{R: 128, G: 0, B: 0, A: 255}
	// ColorPurple is R: 128, G: 0, B: 128.
	ColorPurple = Color{R: 128, G: 0, B: 128, A: 255}
	// ColorFuchsia is R: 255, G: 0, B: 255.
	ColorFuchsia = Color{R: 255, G: 0, B: 255, A: 255}
	// ColorLime is R: 0, G: 255, B: 0.
	ColorLime = Color{R: 0, G: 255, B: 0, A: 255}
	// ColorOlive is R: 128, G: 128, B: 0.
	ColorOlive = Color{R: 128, G: 128, B: 0, A: 255}
	// ColorYellow is R: 255, G: 255, B: 0.
	ColorYellow = Color{R: 255, G: 255, B: 0, A: 255}
	// ColorNavy is R: 0, G: 0, B: 128.
	ColorNavy = Color{R: 0, G: 0, B: 128, A: 255}
	// ColorTeal is R: 0, G: 128, B: 128.
	ColorTeal = Color{R: 0, G: 128, B: 128, A: 255}
	// ColorAqua is R: 0, G: 255, B: 255.
	ColorAqua = Color{R: 0, G: 255, B: 255, A: 255}

	// select extended colors

	// ColorLightGray is R: 211, G: 211, B: 211.
	ColorLightGray = Color{R: 211, G: 211, B: 211, A: 255}
	// ColorSlateGray is R: 112, G: 128, B: 144.
	ColorSlateGray = Color{R: 112, G: 128, B: 144, A: 255}
	// ColorLightSlateGray is R: 119, G: 136, B: 211.
	ColorLightSlateGray = Color{R: 119, G: 136, B: 211, A: 255}
	// ColorAzure is R: 240, G: 255, B: 255.
	ColorAzure = Color{R: 240, G: 255, B: 255, A: 255}
	// ColorBeige is R: 245, G: 245, B: 220.
	ColorBeige = Color{R: 245, G: 245, B: 220, A: 255}
	// ColorBrown is R: 165, G: 42, B: 42.
	ColorBrown = Color{R: 165, G: 42, B: 42, A: 255}
	// ColorChocolate is R: 210, G: 105, B: 30.
	ColorChocolate = Color{R: 210, G: 105, B: 30, A: 255}
	// ColorCoral is R: 255, G: 127, B: 80.
	ColorCoral = Color{R: 255, G: 127, B: 80, A: 255}
	// ColorLightCoral is R: 240, G: 128, B: 128.
	ColorLightCoral = Color{R: 240, G: 128, B: 128, A: 255}
	// ColorGold is R: 255, G: 215, B: 0.
	ColorGold = Color{R: 255, G: 215, B: 0, A: 255}
	// ColorIndigo is R: 75, G: 0, B: 130.
	ColorIndigo = Color{R: 75, G: 0, B: 130, A: 255}
	// ColorIvory is R: 255, G: 255, B: 250.
	ColorIvory = Color{R: 255, G: 255, B: 250, A: 255}
	// ColorOrange is R: 255, G: 165, B: 0.
	ColorOrange = Color{R: 255, G: 165, B: 0, A: 255}
	// ColorPink is R: 255, G: 192, B: 203.
	ColorPink = Color{R: 255, G: 192, B: 203, A: 255}
	// ColorPlum is R: 221, G: 160, B: 221.
	ColorPlum = Color{R: 221, G: 160, B: 221, A: 255}
	// ColorSalmon is R: 250, G: 128, B: 114.
	ColorSalmon = Color{R: 250, G: 128, B: 114, A: 255}
	// ColorTan is R: 210, G: 180, B: 140.
	ColorTan = Color{R: 210, G: 180, B: 140, A: 255}
	// ColorKhaki is R: 240, G: 230, B: 140.
	ColorKhaki = Color{R: 240, G: 230, B: 140, A: 255}
	// ColorTurquoise is R: 64, G: 224, B: 208.
	ColorTurquoise = Color{R: 64, G: 224, B: 208, A: 255}
	// ColorViolet is R: 238, G: 130, B: 238.
	ColorViolet = Color{R: 238, G: 130, B: 238, A: 255}
	// ColorSkyBlue is R: 135, G: 206, B: 235.
	ColorSkyBlue = Color{R: 135, G: 206, B: 235, A: 255}
	// ColorLavender is R: 230, G: 230, B: 250.
	ColorLavender = Color{R: 230, G: 230, B: 250, A: 255}
	// ColorThistle is R: 216, G: 191, B: 216.
	ColorThistle = Color{R: 216, G: 191, B: 216, A: 255}
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
	if len(arr) > 3 { // if an alpha channel is specified
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
	case "grey", "gray":
		return ColorGray
	case "lightgrey", "lightgray":
		return ColorLightGray
	case "lightslategrey", "lightslategray":
		return ColorLightSlateGray
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
	case "cyan", "aqua":
		return ColorAqua
	case "azure":
		return ColorAzure
	case "beige":
		return ColorBeige
	case "brown":
		return ColorBrown
	case "chocolate":
		return ColorChocolate
	case "coral":
		return ColorCoral
	case "lightcoral":
		return ColorLightCoral
	case "gold":
		return ColorGold
	case "indigo":
		return ColorIndigo
	case "ivory":
		return ColorIvory
	case "orange":
		return ColorOrange
	case "pink":
		return ColorPink
	case "plum":
		return ColorPlum
	case "salmon":
		return ColorSalmon
	case "tan":
		return ColorTan
	case "turquoise":
		return ColorTurquoise
	case "violet":
		return ColorViolet
	case "skyblue":
		return ColorSkyBlue
	case "lavender":
		return ColorLavender
	case "thistle":
		return ColorThistle
	case "":
		return Color{}
	default:
		return ColorBlack
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
	return uint8(clamp(v, 0, 1) * 255)
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

// WithAdjustHSL adjusts the hue, saturation, and lightness of a color by given deltas.
// hueDelta is in degrees; satDelta and lightDelta are in the range [-1, 1].
func (c Color) WithAdjustHSL(hueDelta, satDelta, lightDelta float64) Color {
	h, s, l := c.HSL()
	// Adjust hue (wrap around 360°)
	h = math.Mod(h+hueDelta, 360)
	if h < 0 {
		h += 360
	}
	s = clamp(s+satDelta, 0, 1)
	l = clamp(l+lightDelta, 0, 1)
	r, g, b := hslToRGB(h, s, l)
	return Color{R: r, G: g, B: b, A: c.A}
}

// HSL converts outputs the color HSL (hue, saturation, lightness). H in degrees, S and L in [0,1].
func (c Color) HSL() (h, s, l float64) {
	return rgbToHSL(c.R, c.G, c.B)
}

// rgbToHSL converts an RGB color (with components 0–255) to HSL (H in degrees, S and L in [0,1]).
func rgbToHSL(ri, gi, bi uint8) (h, s, l float64) {
	// Normalize RGB values to [0,1]
	r := float64(ri) / 255.0
	g := float64(gi) / 255.0
	b := float64(bi) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l = (max + min) / 2

	if max != min {
		d := max - min
		s = d / (1 - math.Abs(2*l-1))
		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h *= 60
	} // else, achromatic (gray), s and h should be 0
	return
}

// hslToRGB converts an HSL color (H in degrees, S and L in [0,1]) back to RGB (0–255).
func hslToRGB(h, s, l float64) (r, g, b uint8) {
	var rf, gf, bf float64
	if s < matrix.DefaultEpsilon { // Achromatic (gray)
		rf, gf, bf = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		hk := h / 360
		rf = hue2rgb(p, q, hk+1.0/3.0)
		gf = hue2rgb(p, q, hk)
		bf = hue2rgb(p, q, hk-1.0/3.0)
	}
	return ColorChannelFromFloat(rf), ColorChannelFromFloat(gf), ColorChannelFromFloat(bf)
}

func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	} else if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	} else if t < 1.0/2.0 {
		return q
	} else if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

// clamp confines a value to the given range.
func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
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

// String returns a CSS string representation of the color.
func (c Color) String() string {
	switch c {
	case ColorWhite:
		return "white"
	case ColorBlack:
		return "black"
	case ColorRed:
		return "red"
	case ColorBlue:
		return "blue"
	case ColorGreen:
		return "green"
	case ColorSilver:
		return "silver"
	case ColorMaroon:
		return "maroon"
	case ColorPurple:
		return "purple"
	case ColorFuchsia:
		return "fuchsia"
	case ColorLime:
		return "lime"
	case ColorOlive:
		return "olive"
	case ColorYellow:
		return "yellow"
	case ColorNavy:
		return "navy"
	case ColorTeal:
		return "teal"
	case ColorAqua:
		return "aqua"
	default:
		if c.A == 255.0 {
			return c.StringRGB()
		} else {
			return c.StringRGBA()
		}
	}
}

// StringRGB returns a CSS RGB string representation of the color.
func (c Color) StringRGB() string {
	return "rgb(" + strconv.Itoa(int(c.R)) + "," +
		strconv.Itoa(int(c.G)) + "," +
		strconv.Itoa(int(c.B)) + ")"
}

// StringRGBA returns a CSS RGBA string representation of the color.
func (c Color) StringRGBA() string {
	return fmt.Sprintf("rgba(%v,%v,%v,%.1f)", c.R, c.G, c.B, float64(c.A)/255.0)
}
