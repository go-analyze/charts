package charts

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func isLightColor(c Color) bool {
	r := float64(c.R) * float64(c.R) * 0.299
	g := float64(c.G) * float64(c.G) * 0.587
	b := float64(c.B) * float64(c.B) * 0.114
	return math.Sqrt(r+g+b) > 127.5
}

func parseColor(color string) Color {
	c := Color{}
	if color == "" {
		return c
	}
	if strings.HasPrefix(color, "#") {
		return drawing.ColorFromHex(color[1:])
	}
	reg := regexp.MustCompile(`\((\S+)\)`)
	result := reg.FindAllStringSubmatch(color, 1)
	if len(result) == 0 || len(result[0]) != 2 {
		return c
	}
	arr := strings.Split(result[0][1], ",")
	if len(arr) < 3 {
		return c
	}
	// set the default value to 255
	c.A = 255
	for index, v := range arr {
		value, _ := strconv.Atoi(strings.TrimSpace(v))
		ui8 := uint8(value)
		switch index {
		case 0:
			c.R = ui8
		case 1:
			c.G = ui8
		case 2:
			c.B = ui8
		default:
			c.A = ui8
		}
	}
	return c
}
