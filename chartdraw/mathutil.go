package chartdraw

import "math"

const (
	_pi  = math.Pi
	_2pi = 2 * math.Pi
	//_3pi4 = (3 * math.Pi) / 4.0
	//_4pi3 = (4 * math.Pi) / 3.0
	_3pi2 = (3 * math.Pi) / 2.0
	//_5pi4 = (5 * math.Pi) / 4.0
	//_7pi4 = (7 * math.Pi) / 4.0
	_pi2 = math.Pi / 2.0
	//_pi4  = math.Pi / 4.0
	_d2r = math.Pi / 180.0
	_r2d = 180.0 / math.Pi
)

// MinMax returns the minimum and maximum of a given set of values.
func MinMax(values ...float64) (min, max float64) {
	if len(values) == 0 {
		return
	}

	max = values[0]
	min = values[0]
	var value float64
	for index := 1; index < len(values); index++ {
		value = values[index]
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return
}

// MinInt returns the minimum int.
func MinInt(values ...int) int {
	min := math.MaxInt32
	for x := 0; x < len(values); x++ {
		if values[x] < min {
			min = values[x]
		}
	}
	return min
}

// MaxInt returns the maximum int.
func MaxInt(values ...int) int {
	max := math.MinInt32
	for x := 0; x < len(values); x++ {
		if values[x] > max {
			max = values[x]
		}
	}
	return max
}

// AbsInt returns the absolute value of an int.
func AbsInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

// DegreesToRadians returns degrees as radians.
func DegreesToRadians(degrees float64) float64 {
	return degrees * _d2r
}

// RadiansToDegrees translates a radian value to a degree value.
func RadiansToDegrees(value float64) float64 {
	degrees := value * _r2d
	return DegreesAdd(degrees, 0)
}

// PercentToRadians converts a normalized value (0,1) to radians.
func PercentToRadians(pct float64) float64 {
	return DegreesToRadians(360.0 * pct)
}

// RadianAdd adds a delta to a base in radians.
func RadianAdd(base, delta float64) float64 {
	value := base + delta
	return math.Mod(math.Mod(value, _2pi)+_2pi, _2pi)
}

// DegreesAdd adds a delta to a base in degrees.
func DegreesAdd(baseDegrees, deltaDegrees float64) float64 {
	value := baseDegrees + deltaDegrees
	return math.Mod(math.Mod(value, 360.0)+360.0, 360.0)
}

// CirclePoint returns the absolute position of a circle diameter point given
// by the radius and the theta.
func CirclePoint(cx, cy int, radius, thetaRadians float64) (x, y int) {
	x = cx + int(radius*math.Sin(thetaRadians))
	y = cy - int(radius*math.Cos(thetaRadians))
	return
}

// RotateCoordinate rotates a coordinate around a given center by a theta in radians.
func RotateCoordinate(cx, cy, x, y int, thetaRadians float64) (rx, ry int) {
	tempX, tempY := float64(x-cx), float64(y-cy)
	rotatedX := tempX*math.Cos(thetaRadians) - tempY*math.Sin(thetaRadians)
	rotatedY := tempX*math.Sin(thetaRadians) + tempY*math.Cos(thetaRadians)
	rx = int(rotatedX) + cx
	ry = int(rotatedY) + cy
	return
}

// RoundUp rounds up to a given roundTo value.
func RoundUp(value, roundTo float64) float64 {
	if roundTo < 0.000000000000001 {
		return value
	}
	d1 := math.Ceil(value / roundTo)
	return d1 * roundTo
}

// RoundDown rounds down to a given roundTo value.
func RoundDown(value, roundTo float64) float64 {
	if roundTo < 0.000000000000001 {
		return value
	}
	d1 := math.Floor(value / roundTo)
	return d1 * roundTo
}

// Normalize returns a set of numbers on the interval [0,1] for a given set of inputs.
// An example: 4,3,2,1 => 0.4, 0.3, 0.2, 0.1
// Caveat; the total may be < 1.0; there are going to be issues with irrational numbers etc.
func Normalize(values ...float64) []float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	output := make([]float64, len(values))
	if total != 0 {
		for x, v := range values {
			output[x] = RoundDown(v/total, 0.0001)
		}
	}
	return output
}

// MeanFloat64 returns the mean of a set of values
func MeanFloat64(values ...float64) float64 {
	return SumFloat64(values...) / float64(len(values))
}

// MeanInt returns the mean of a set of integer values.
func MeanInt(values ...int) int {
	return SumInt(values...) / len(values)
}

// SumFloat64 sums a set of values.
func SumFloat64(values ...float64) float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	return total
}

// SumInt sums a set of values, if an overflow occurs math.MaxInt will be returned exactly. Similar on an underflow
// math.MinInt is returned.
func SumInt(values ...int) int {
	var sum int
	for _, v := range values {
		if v > 0 && sum > math.MaxInt-v {
			return math.MaxInt
		} else if v < 0 && sum < math.MinInt-v {
			return math.MinInt
		}
		sum += v
	}
	return sum
}

// PercentDifference computes the percentage difference between two values.
// The formula is (v2-v1)/v1.
func PercentDifference(v1, v2 float64) float64 {
	if v1 == 0 {
		return 0
	}
	return (v2 - v1) / v1
}

// GetRoundToForDelta returns a `roundTo` value for a given delta.
func GetRoundToForDelta(delta float64) float64 {
	startingDeltaBound := math.Pow(10.0, 10.0)
	for cursor := startingDeltaBound; cursor > 0; cursor /= 10.0 {
		if delta > cursor {
			return cursor / 10.0
		}
	}

	return 0.0
}

// RoundPlaces rounds an input to a given places.
func RoundPlaces(input float64, places int) (rounded float64) {
	if math.IsNaN(input) {
		return 0.0
	}

	sign := 1.0
	if input < 0 {
		sign = -1
		input *= -1
	}

	precision := math.Pow(10, float64(places))
	digit := input * precision
	_, decimal := math.Modf(digit)

	if decimal >= 0.5 {
		rounded = math.Ceil(digit)
	} else {
		rounded = math.Floor(digit)
	}

	return rounded / precision * sign
}

func f64i(value float64) int {
	r := RoundPlaces(value, 0)
	return int(r)
}
