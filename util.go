package charts

import (
	"math"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/go-analyze/bulk"

	"github.com/go-analyze/charts/chartdraw"
)

// Ptr is a helper function to help build config options which reference pointers.
func Ptr[T any](val T) *T {
	return &val
}

// flagIs returns true if the flag is not-nil and matches the comparison argument.
func flagIs(is bool, flag *bool) bool {
	if flag == nil {
		return false
	}
	return *flag == is
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
		var end int
		for index, v := range spans {
			end += v
			newValues[index+1] = values[end]
		}
		values = newValues
	}
	return values
}

func reverseSlice[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// SliceToFloat64 converts a slice of arbitrary types to float64 to be used as chart values.
func SliceToFloat64[T any](slice []T, conversion func(T) float64) []float64 {
	return bulk.SliceTransform(conversion, slice)
}

// IntSliceToFloat64 converts an int slice to a float64 slice so that it can be used for chart values.
func IntSliceToFloat64(slice []int) []float64 {
	return bulk.SliceTransform(func(i int) float64 { return float64(i) }, slice)
}

func sliceMaxLen[T any](values ...[]T) int {
	result := 0
	for _, slice := range values {
		if len(slice) > result {
			result = len(slice)
		}
	}
	return result
}

func parseFlexibleValue(value string, percentTotal float64) (float64, error) {
	if strings.HasSuffix(value, "%") {
		percent, err := strconv.ParseFloat(strings.TrimSuffix(value, "%"), 64)
		if err != nil {
			return 0, err
		}
		return (percent / 100.0) * percentTotal, nil
	} else {
		return strconv.ParseFloat(value, 64)
	}
}

const kValue = float64(1000)
const mValue = kValue * kValue
const gValue = mValue * kValue
const tValue = gValue * kValue

// DegreesToRadians returns degrees as radians.
func DegreesToRadians(degrees float64) float64 {
	return chartdraw.DegreesToRadians(degrees)
}

// RadiansToDegrees translates a radian value to a degree value.
func RadiansToDegrees(value float64) float64 {
	return chartdraw.RadiansToDegrees(value)
}

// normalizeAngle brings the angle into the range [0, 2π).
func normalizeAngle(radians float64) float64 {
	radians = math.Mod(radians, _2pi)
	if radians < 0 {
		radians += _2pi
	}
	return radians
}

// FormatValueHumanizeShort takes in a value and a specified precision, rounding to the specified precision and
// returning a human friendly number string including commas. If the value is over 1,000 it will be reduced to a
// shorter version with the appropriate k, M, G, T suffix.
func FormatValueHumanizeShort(value float64, decimals int, ensureTrailingZeros bool) string {
	// Handle negative values by delegating to the positive case and then re-applying the sign
	if value < 0 {
		return "-" + FormatValueHumanizeShort(-value, decimals, ensureTrailingZeros)
	}

	if value >= tValue {
		return FormatValueHumanize(value/tValue, decimals, ensureTrailingZeros) + "T"
	} else if value >= gValue {
		return FormatValueHumanize(value/gValue, decimals, ensureTrailingZeros) + "G"
	} else if value >= mValue {
		return FormatValueHumanize(value/mValue, decimals, ensureTrailingZeros) + "M"
	} else if value >= kValue {
		return FormatValueHumanize(value/kValue, decimals, ensureTrailingZeros) + "k"
	} else {
		return FormatValueHumanize(value, decimals, ensureTrailingZeros)
	}
}

// FormatValueHumanize takes in a value and a specified precision, rounding to the specified precision and returning a
// human friendly number string including commas.
func FormatValueHumanize(value float64, decimals int, ensureTrailingZeros bool) string {
	if decimals < 0 {
		decimals = 0
	}
	multiplier := math.Pow(10, float64(decimals))
	roundedValue := math.Round(value*multiplier) / multiplier

	result := humanize.CommafWithDigits(roundedValue, decimals)

	if ensureTrailingZeros && decimals > 0 {
		if decimalIndex := strings.IndexAny(result, "."); decimalIndex == -1 {
			return result + "." + strings.Repeat("0", decimals)
		} else if existingDecimals := len(result) - decimalIndex - 1; existingDecimals < decimals {
			return result + strings.Repeat("0", decimals-existingDecimals)
		}
	}

	return result
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
