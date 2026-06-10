package matrix

import (
	"math"
	"strconv"
)

func f64s(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func roundToEpsilon(value, epsilon float64) float64 {
	// TODO - epsilon is not used here, this does not appear to be as the function describes
	return math.Nextafter(value, value)
}
