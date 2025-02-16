package charts

import (
	"math"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"

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
	return sliceConversion(slice, conversion)
}

// IntSliceToFloat64 converts an int slice to a float64 slice so that it can be used for chart values.
func IntSliceToFloat64(slice []int) []float64 {
	return sliceConversion(slice, func(i int) float64 { return float64(i) })
}

func sliceConversion[I any, R any](input []I, conversion func(I) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = conversion(v)
	}
	return result
}

// sliceSplit will split a slice in half based on a conditional function passed in. The right result are true values,
// second result being values that tested false.
func sliceSplit[T any](slice []T, test func(v T) bool) ([]T, []T) {
	if len(slice) == 0 {
		return nil, nil
	}

	var splitIndex int
	first := test(slice[0])
	for splitIndex = 1; splitIndex < len(slice); splitIndex++ {
		if first != test(slice[splitIndex]) {
			break
		}
	}

	// If all are the same, return early
	if splitIndex == len(slice) {
		if first {
			return slice, nil
		} else {
			return nil, slice
		}
	}

	// Allocate slices and copy first segment
	remainingBuff := len(slice) - splitIndex
	if remainingBuff > 2048 {
		remainingBuff /= 2
	}
	var trueList, falseList []T
	if first {
		trueList = append(make([]T, 0, splitIndex+remainingBuff-1), slice[:splitIndex]...)
		falseList = append(make([]T, 0, remainingBuff), slice[splitIndex])
	} else {
		falseList = append(make([]T, 0, splitIndex+remainingBuff-1), slice[:splitIndex]...)
		trueList = append(make([]T, 0, remainingBuff), slice[splitIndex])
	}
	// Finish iterating appending remaining elements
	for i := splitIndex + 1; i < len(slice); i++ {
		if test(slice[i]) {
			trueList = append(trueList, slice[i])
		} else {
			falseList = append(falseList, slice[i])
		}
	}

	return trueList, falseList
}

// sliceFilter iterates over the slice, testing each element with the provided function. The returned slice are items
// which had a true result.
func sliceFilter[T any](slice []T, test func(v T) bool) []T {
	for falseIndex, v := range slice {
		if !test(v) {
			if falseIndex == 0 {
				// iterate until a true result is found, then start appending at that point
				var result []T
				for i := falseIndex + 1; i < len(slice); i++ {
					if test(slice[i]) {
						if result == nil {
							remainingBuff := len(slice) - i
							if remainingBuff > 2048 {
								remainingBuff /= 2
							}
							result = make([]T, 0, remainingBuff)
						}
						result = append(result, slice[i])
					}
				}
				return result
			} else {
				// copy all records that already passed, and then finish iteration to produce result
				remainingBuff := len(slice) - falseIndex - 1
				if remainingBuff > 2048 {
					remainingBuff /= 2
				}
				result := append(make([]T, 0, falseIndex+remainingBuff), slice[:falseIndex]...)
				for i := falseIndex + 1; i < len(slice); i++ {
					if test(slice[i]) {
						result = append(result, slice[i])
					}
				}
				return result
			}
		}
	}
	return slice // all records tested to true
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

// FormatValueHumanizeShort takes in a value and a specified precision, rounding to the specified precision and
// returning a human friendly number string including commas. If the value is over 1,000 it will be reduced to a
// shorter version with the appropriate k, M, G, T suffix.
func FormatValueHumanizeShort(value float64, decimals int, ensureTrailingZeros bool) string {
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

const defaultRadiusPercent = 0.4

func getRadius(diameter float64, radiusValue string) float64 {
	var radius float64
	if radiusValue != "" {
		radius, _ = parseFlexibleValue(radiusValue, diameter)
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
