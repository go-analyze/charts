package charts

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/wcharczuk/go-chart/v2"
)

// True returns a pointer to a true bool, useful for configuration.
func True() *bool {
	return BoolPointer(true)
}

// False returns a pointer to a false bool, useful for configuration.
func False() *bool {
	return BoolPointer(false)
}

// BoolPointer returns a pointer to the given bool value, useful for configuration.
func BoolPointer(b bool) *bool {
	return &b
}

// FloatPointer returns a pointer to the given float64 value, useful for configuration.
func FloatPointer(f float64) *float64 {
	return &f
}

// flagIs returns true if the flag is not-nil and matches the comparison argument.
func flagIs(is bool, flag *bool) bool {
	if flag == nil {
		return false
	}
	return *flag == is
}

func containsInt(values []int, value int) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func ceilFloatToInt(value float64) int {
	if value >= float64(math.MaxInt) {
		return math.MaxInt
	} else if value <= float64(math.MinInt) {
		return math.MinInt
	}

	i := int(value)
	if value == float64(i) {
		return i
	}
	return i + 1
}

func getDefaultInt(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

func autoDivide(max, size int) []int {
	unit := float64(max) / float64(size)

	values := make([]int, size+1)
	for i := 0; i < size+1; i++ {
		if i == size {
			values[i] = max
		} else {
			values[i] = int(float64(i) * unit)
		}
	}
	return values
}

func autoDivideSpans(max, size int, spans []int) []int {
	values := autoDivide(max, size)
	// re-merge
	if len(spans) != 0 {
		newValues := make([]int, len(spans)+1)
		newValues[0] = 0
		end := 0
		for index, v := range spans {
			end += v
			newValues[index+1] = values[end]
		}
		values = newValues
	}
	return values
}

func sumInt(values []int) int {
	// chart.SumInt() also exists, but it does not handle overflow like we are
	sum := 0
	for _, v := range values {
		if v > 0 && (math.MaxInt-sum) < v {
			return math.MaxInt
		} else if v < 0 && math.MinInt-sum > v {
			return math.MinInt
		}
		sum += v
	}
	return sum
}

// measureTextMaxWidthHeight returns maxWidth and maxHeight of text list
func measureTextMaxWidthHeight(textList []string, p *Painter) (int, int) {
	maxWidth := 0
	maxHeight := 0
	for _, text := range textList {
		box := p.MeasureText(text)
		maxWidth = chart.MaxInt(maxWidth, box.Width())
		maxHeight = chart.MaxInt(maxHeight, box.Height())
	}
	return maxWidth, maxHeight
}

func reverseStringSlice(stringList []string) {
	for i, j := 0, len(stringList)-1; i < j; i, j = i+1, j-1 {
		stringList[i], stringList[j] = stringList[j], stringList[i]
	}
}

func reverseIntSlice(intList []int) {
	for i, j := 0, len(intList)-1; i < j; i, j = i+1, j-1 {
		intList[i], intList[j] = intList[j], intList[i]
	}
}

func parseFlexibleValue(value string, percentTotal float64) (float64, error) {
	if strings.HasSuffix(value, "%") {
		percent, err := convertPercent(value)
		if err != nil {
			return 0, err
		}
		return percent * percentTotal, nil
	} else {
		return strconv.ParseFloat(value, 64)
	}
}

func convertPercent(value string) (float64, error) {
	if !strings.HasSuffix(value, "%") {
		return -1, fmt.Errorf("not a percent input: %s", value)
	}
	v, err := strconv.ParseFloat(strings.ReplaceAll(value, "%", ""), 64)
	if err != nil {
		return -1, err
	}
	return v / 100.0, nil
}

const K_VALUE = float64(1000)
const M_VALUE = K_VALUE * K_VALUE
const G_VALUE = M_VALUE * K_VALUE
const T_VALUE = G_VALUE * K_VALUE

func commafWithDigits(value float64) string {
	decimals := 2
	if value >= T_VALUE {
		return humanize.CommafWithDigits(value/T_VALUE, decimals) + "T"
	}
	if value >= G_VALUE {
		return humanize.CommafWithDigits(value/G_VALUE, decimals) + "G"
	}
	if value >= M_VALUE {
		return humanize.CommafWithDigits(value/M_VALUE, decimals) + "M"
	}
	if value >= K_VALUE {
		return humanize.CommafWithDigits(value/K_VALUE, decimals) + "k"
	}
	return humanize.CommafWithDigits(value, decimals)
}

const defaultRadiusPercent = 0.4

func getRadius(diameter float64, radiusValue string) float64 {
	var radius float64
	if len(radiusValue) != 0 {
		v, _ := convertPercent(radiusValue)
		if v != -1 {
			radius = diameter * v
		} else {
			radius, _ = strconv.ParseFloat(radiusValue, 64)
		}
	}
	if radius <= 0 {
		radius = diameter * defaultRadiusPercent
	}
	return radius
}

func getPolygonPointAngles(sides int) []float64 {
	angles := make([]float64, sides)
	for i := 0; i < sides; i++ {
		angle := 2*math.Pi/float64(sides)*float64(i) - (math.Pi / 2)
		angles[i] = angle
	}
	return angles
}

func getPolygonPoint(center Point, radius, angle float64) Point {
	x := center.X + int(radius*math.Cos(angle))
	y := center.Y + int(radius*math.Sin(angle))
	return Point{X: x, Y: y}
}

func getPolygonPoints(center Point, radius float64, sides int) []Point {
	points := make([]Point, sides)
	for i, angle := range getPolygonPointAngles(sides) {
		points[i] = getPolygonPoint(center, radius, angle)
	}
	return points
}
