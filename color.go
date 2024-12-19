package charts

import (
	"math"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func isLightColor(c Color) bool {
	r := float64(c.R) * float64(c.R) * 0.299
	g := float64(c.G) * float64(c.G) * 0.587
	b := float64(c.B) * float64(c.B) * 0.114
	return math.Sqrt(r+g+b) > 127.5
}

// ParseColor parses a color from a string. The color can be specified in hex with a `#` prefix (for example '#313233'),
// in rgb(i,i,i) or rgba(i,i,i,f) format, or as a common name (for example 'red').
func ParseColor(color string) Color {
	return drawing.ParseColor(color)
}
