package chartdraw

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinMax(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		values      []float64
		expectedMin float64
		expectedMax float64
	}{
		{
			name:        "EmptySlice",
			values:      []float64{},
			expectedMin: 0.0,
			expectedMax: 0.0,
		},
		{
			name:        "SingleElement",
			values:      []float64{5.0},
			expectedMin: 5.0,
			expectedMax: 5.0,
		},
		{
			name:        "MultipleElements",
			values:      []float64{1.0, 2.0, 3.0, 4.0},
			expectedMin: 1.0,
			expectedMax: 4.0,
		},
		{
			name:        "IncludingNegatives",
			values:      []float64{-1.0, -2.0, 0.0, 3.0},
			expectedMin: -2.0,
			expectedMax: 3.0,
		},
		{
			name:        "AllNegative",
			values:      []float64{-1.0, -2.0, -3.0, -4.0},
			expectedMin: -4.0,
			expectedMax: -1.0,
		},
		{
			name:        "AllIdentical",
			values:      []float64{2.0, 2.0, 2.0, 2.0},
			expectedMin: 2.0,
			expectedMax: 2.0,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			min, max := MinMax(tc.values...)
			assert.Equal(t, tc.expectedMin, min)
			assert.Equal(t, tc.expectedMax, max)
		})
	}
}

func TestMinInt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		values      []int
		expectedMin int
	}{
		{
			name:        "EmptySlice",
			values:      []int{},
			expectedMin: 0,
		},
		{
			name:        "SingleElement",
			values:      []int{5},
			expectedMin: 5,
		},
		{
			name:        "MultipleElements",
			values:      []int{1, 2, 3, 4},
			expectedMin: 1,
		},
		{
			name:        "IncludingNegatives",
			values:      []int{-1, -2, 0, 3},
			expectedMin: -2,
		},
		{
			name:        "AllNegative",
			values:      []int{-1, -20, -3, -40},
			expectedMin: -40,
		},
		{
			name:        "AllIdentical",
			values:      []int{3, 3, 3, 3},
			expectedMin: 3,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			min := MinInt(tc.values...)
			assert.Equal(t, tc.expectedMin, min)
		})
	}
}

func TestMaxInt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		values      []int
		expectedMax int
	}{
		{
			name:        "EmptySlice",
			values:      []int{},
			expectedMax: 0,
		},
		{
			name:        "SingleElement",
			values:      []int{5},
			expectedMax: 5,
		},
		{
			name:        "MultipleElements",
			values:      []int{1, 2, 3, 4},
			expectedMax: 4,
		},
		{
			name:        "IncludingNegatives",
			values:      []int{-1, -2, 0, 3},
			expectedMax: 3,
		},
		{
			name:        "AllNegative",
			values:      []int{-1, -20, -3, -40},
			expectedMax: -1,
		},
		{
			name:        "AllIdentical",
			values:      []int{7, 7, 7, 7},
			expectedMax: 7,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			max := MaxInt(tc.values...)
			assert.Equal(t, tc.expectedMax, max)
		})
	}
}

func TestAbsInt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		value       int
		expectedAbs int
	}{
		{
			name:        "Positive",
			value:       5,
			expectedAbs: 5,
		},
		{
			name:        "Negative",
			value:       -5,
			expectedAbs: 5,
		},
		{
			name:        "Zero",
			value:       0,
			expectedAbs: 0,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			abs := AbsInt(tc.value)
			assert.Equal(t, tc.expectedAbs, abs)
		})
	}
}

func TestDegreesToRadians(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		degrees      float64
		expectedRads float64
	}{
		{
			name:         "ZeroDegrees",
			degrees:      0.0,
			expectedRads: 0.0,
		},
		{
			name:         "RightAngle",
			degrees:      90.0,
			expectedRads: _pi2,
		},
		{
			name:         "FullCircle",
			degrees:      360.0,
			expectedRads: _2pi,
		},
		{
			name:         "NegativeAngle",
			degrees:      -180.0,
			expectedRads: -math.Pi,
		},
		{
			name:         "Over360Degrees",
			degrees:      450.0,
			expectedRads: _pi2 + _2pi,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			rads := DegreesToRadians(tc.degrees)
			assert.Equal(t, tc.expectedRads, rads)
		})
	}
}

func TestRadiansToDegrees(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		radians      float64
		expectedDegs float64
	}{
		{
			name:         "Zero",
			radians:      0.0,
			expectedDegs: 0.0,
		},
		{
			name:         "Pi",
			radians:      math.Pi,
			expectedDegs: 180.0,
		},
		{
			name:         "TwoPi",
			radians:      _2pi,
			expectedDegs: 0.0, // Tests modulus wrap-around
		},
		{
			name:         "NegativePi",
			radians:      -math.Pi,
			expectedDegs: 180.0, // Tests negative wrap-around
		},
		{
			name:         "OverTwoPi",
			radians:      3 * math.Pi,
			expectedDegs: 180.0, // Tests multiple modulus wrap-around
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			degs := RadiansToDegrees(tc.radians)
			assert.Equal(t, tc.expectedDegs, degs)
		})
	}
}

func TestPercentToRadians(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		percent      float64
		expectedRads float64
	}{
		{
			name:         "Zero",
			percent:      0.0,
			expectedRads: 0.0,
		},
		{
			name:         "Fifty",
			percent:      0.5,
			expectedRads: math.Pi,
		},
		{
			name:         "Full",
			percent:      1.0,
			expectedRads: _2pi,
		},
		{
			name:         "Negative",
			percent:      -0.5,
			expectedRads: -math.Pi,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			rads := PercentToRadians(tc.percent)
			assert.Equal(t, tc.expectedRads, rads)
		})
	}
}

func TestRadianAdd(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		base, delta   float64
		expectedValue float64
	}{
		{
			name:          "AddZero",
			base:          math.Pi,
			delta:         0.0,
			expectedValue: math.Pi,
		},
		{
			name:          "AddPi",
			base:          _pi,
			delta:         _pi,
			expectedValue: 0.0, // Tests wrap-around
		},
		{
			name:          "NegativeDelta",
			base:          0.0,
			delta:         -_pi2,
			expectedValue: _3pi2, // Tests negative wrap-around
		},
		{
			name:          "OverTwoPi",
			base:          _2pi,
			delta:         _pi,
			expectedValue: math.Pi, // Tests modulus wrap-around
		},
		{
			name:          "UnderNegativeTwoPi",
			base:          0.0,
			delta:         -_3pi2,
			expectedValue: _pi2, // Tests modulus wrap-around
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			result := RadianAdd(tc.base, tc.delta)
			assert.Equal(t, tc.expectedValue, result)
		})
	}
}

func TestDegreesAdd(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                      string
		baseDegrees, deltaDegrees float64
		expectedValue             float64
	}{
		{
			name:          "AddZero",
			baseDegrees:   180.0,
			deltaDegrees:  0.0,
			expectedValue: 180.0,
		},
		{
			name:          "Add180",
			baseDegrees:   180.0,
			deltaDegrees:  180.0,
			expectedValue: 0.0, // Tests 360 wrap-around
		},
		{
			name:          "NegativeDelta",
			baseDegrees:   0.0,
			deltaDegrees:  -90.0,
			expectedValue: 270.0, // Tests negative wrap-around
		},
		{
			name:          "Over360",
			baseDegrees:   360.0,
			deltaDegrees:  90.0,
			expectedValue: 90.0, // Tests modulus wrap-around
		},
		{
			name:          "UnderNegative360",
			baseDegrees:   0.0,
			deltaDegrees:  -450.0,
			expectedValue: 270.0, // Tests modulus negative wrap-around
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			result := DegreesAdd(tc.baseDegrees, tc.deltaDegrees)
			assert.Equal(t, tc.expectedValue, result)
		})
	}
}

func TestCirclePoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		cx, cy       int
		radius       float64
		thetaRadians float64
		expectedX    int
		expectedY    int
	}{
		{
			name:         "ZeroRadius",
			cx:           10,
			cy:           10,
			radius:       0,
			thetaRadians: 0.0,
			expectedX:    10,
			expectedY:    10,
		},
		{
			name:         "QuarterCircle",
			cx:           10,
			cy:           10,
			radius:       10,
			thetaRadians: _pi2,
			expectedX:    20,
			expectedY:    10,
		},
		{
			name:         "HalfCircle",
			cx:           10,
			cy:           10,
			radius:       10,
			thetaRadians: math.Pi,
			expectedX:    10,
			expectedY:    20,
		},
		{
			name:         "FullCircle",
			cx:           10,
			cy:           10,
			radius:       10,
			thetaRadians: _2pi,
			expectedX:    10,
			expectedY:    0,
		},
		{
			name:         "NegativeQuarterCircle",
			cx:           10,
			cy:           10,
			radius:       10,
			thetaRadians: -_pi2,
			expectedX:    0,
			expectedY:    10,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			x, y := CirclePoint(tc.cx, tc.cy, tc.radius, tc.thetaRadians)
			assert.Equal(t, tc.expectedX, x, "x failure")
			assert.Equal(t, tc.expectedY, y, "y failure")
		})
	}
}

func TestRotateCoordinate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		cx, cy, x, y         int
		thetaRadians         float64
		expectedX, expectedY int
	}{
		{
			name:         "NoRotation",
			cx:           0,
			cy:           0,
			x:            1,
			y:            0,
			thetaRadians: 0.0,
			expectedX:    1,
			expectedY:    0,
		},
		{
			name:         "QuarterRotation",
			cx:           0,
			cy:           0,
			x:            1,
			y:            0,
			thetaRadians: _pi2,
			expectedX:    0,
			expectedY:    1,
		},
		{
			name:         "ThreeQuarterRotation",
			cx:           0,
			cy:           0,
			x:            1,
			y:            0,
			thetaRadians: 3 * _pi2,
			expectedX:    0,
			expectedY:    -1,
		},
		{
			name:         "FullRotation",
			cx:           0,
			cy:           0,
			x:            1,
			y:            0,
			thetaRadians: _2pi,
			expectedX:    1,
			expectedY:    0,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			rx, ry := RotateCoordinate(tc.cx, tc.cy, tc.x, tc.y, tc.thetaRadians)
			assert.Equal(t, tc.expectedX, rx)
			assert.Equal(t, tc.expectedY, ry)
		})
	}
}

func TestRoundUp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		value, roundTo float64
		expectedValue  float64
	}{
		{
			name:          "Exact",
			value:         2.5,
			roundTo:       0.5,
			expectedValue: 2.5,
		},
		{
			name:          "Fraction",
			value:         2.3,
			roundTo:       0.5,
			expectedValue: 2.5,
		},
		{
			name:          "Large",
			value:         2.9999,
			roundTo:       0.1,
			expectedValue: 3.0,
		},
		{
			name:          "Negative",
			value:         -2.1,
			roundTo:       0.5,
			expectedValue: -2.0,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			rounded := RoundUp(tc.value, tc.roundTo)
			assert.Equal(t, tc.expectedValue, rounded)
		})
	}
}

func TestRoundDown(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		value, roundTo float64
		expectedValue  float64
	}{
		{
			name:          "Exact",
			value:         2.5,
			roundTo:       0.5,
			expectedValue: 2.5,
		},
		{
			name:          "Fraction",
			value:         2.7,
			roundTo:       0.5,
			expectedValue: 2.5,
		},
		{
			name:          "Large",
			value:         2.9999,
			roundTo:       1.0,
			expectedValue: 2.0,
		},
		{
			name:          "Negative",
			value:         -2.7,
			roundTo:       0.5,
			expectedValue: -3.0,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			rounded := RoundDown(tc.value, tc.roundTo)
			assert.Equal(t, tc.expectedValue, rounded)
		})
	}
}

func TestNormalize(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		values        []float64
		expectedNorms []float64
	}{
		{
			name:          "UniformValues",
			values:        []float64{10, 10, 10, 10},
			expectedNorms: []float64{0.25, 0.25, 0.25, 0.25},
		},
		{
			name:          "VaryingValues",
			values:        []float64{5, 15, 30, 50},
			expectedNorms: []float64{0.05, 0.15, 0.3, 0.5},
		},
		{
			name:          "SingleValue",
			values:        []float64{20},
			expectedNorms: []float64{1.0},
		},
		{
			name:          "IncludingZero",
			values:        []float64{0, 20, 30, 50},
			expectedNorms: []float64{0.0, 0.2, 0.3, 0.5},
		},
		{
			name:          "AllZeros",
			values:        []float64{0, 0, 0, 0},
			expectedNorms: []float64{0.0, 0.0, 0.0, 0.0},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			norms := Normalize(tc.values...)
			for i, norm := range norms {
				assert.InDelta(t, tc.expectedNorms[i], norm, 0.0001)
			}
		})
	}
}

func TestMeanFloat64(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		values       []float64
		expectedMean float64
	}{
		{
			name:         "PositiveValues",
			values:       []float64{10, 20, 30},
			expectedMean: 20.0,
		},
		{
			name:         "IncludingNegative",
			values:       []float64{-10, 0, 10},
			expectedMean: 0.0,
		},
		{
			name:         "AllZero",
			values:       []float64{0, 0, 0},
			expectedMean: 0.0,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			mean := MeanFloat64(tc.values...)
			assert.Equal(t, tc.expectedMean, mean)
		})
	}
}

func TestMeanInt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		values       []int
		expectedMean int
	}{
		{
			name:         "PositiveValues",
			values:       []int{10, 20, 30},
			expectedMean: 20,
		},
		{
			name:         "IncludingNegative",
			values:       []int{-10, 0, 10},
			expectedMean: 0,
		},
		{
			name:         "AllZero",
			values:       []int{0, 0, 0},
			expectedMean: 0,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			mean := MeanInt(tc.values...)
			assert.Equal(t, tc.expectedMean, mean)
		})
	}
}

func TestSumFloat64(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		values      []float64
		expectedSum float64
	}{
		{
			name:        "PositiveValues",
			values:      []float64{1.0, 2.0, 3.0},
			expectedSum: 6.0,
		},
		{
			name:        "IncludingNegative",
			values:      []float64{-1.0, 2.0, -3.0},
			expectedSum: -2.0,
		},
		{
			name:        "AllZeros",
			values:      []float64{0.0, 0.0, 0.0},
			expectedSum: 0.0,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			sum := SumFloat64(tc.values...)
			assert.Equal(t, tc.expectedSum, sum)
		})
	}
}

func TestSumInt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		values      []int
		expectedSum int
	}{
		{
			name:        "PositiveValues",
			values:      []int{1, 2, 3},
			expectedSum: 6,
		},
		{
			name:        "IncludingNegative",
			values:      []int{-1, 2, -3},
			expectedSum: -2,
		},
		{
			name:        "AllZeros",
			values:      []int{0, 0, 0},
			expectedSum: 0,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			sum := SumInt(tc.values...)
			assert.Equal(t, tc.expectedSum, sum)
		})
	}
}

func TestPercentDifference(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		v1, v2             float64
		expectedDifference float64
	}{
		{
			name:               "ZeroInitial",
			v1:                 0,
			v2:                 10,
			expectedDifference: 0,
		},
		{
			name:               "PositiveIncrease",
			v1:                 10,
			v2:                 20,
			expectedDifference: 1.0,
		},
		{
			name:               "NegativeDecrease",
			v1:                 20,
			v2:                 10,
			expectedDifference: -0.5,
		},
		{
			name:               "NoChange",
			v1:                 10,
			v2:                 10,
			expectedDifference: 0.0,
		},
		{
			name:               "NegativeToPositive",
			v1:                 -10,
			v2:                 10,
			expectedDifference: -2.0,
		},
		{
			name:               "LargeIncrease",
			v1:                 1,
			v2:                 1000,
			expectedDifference: 999.0,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			diff := PercentDifference(tc.v1, tc.v2)
			assert.Equal(t, tc.expectedDifference, diff)
		})
	}
}

func TestRoundPlaces(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		input          float64
		places         int
		expectedOutput float64
	}{
		{
			name:           "RoundToTwoPlaces",
			input:          3.14159,
			places:         2,
			expectedOutput: 3.14,
		},
		{
			name:           "RoundToZeroPlaces",
			input:          3.14159,
			places:         0,
			expectedOutput: 3.0,
		},
		{
			name:           "NegativeValue",
			input:          -3.14159,
			places:         2,
			expectedOutput: -3.14,
		},
		{
			name:           "SmallDecimal",
			input:          0.00098765,
			places:         5,
			expectedOutput: 0.00099,
		},
		{
			name:           "NegativeDecimal",
			input:          -0.00098765,
			places:         5,
			expectedOutput: -0.00099,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			rounded := RoundPlaces(tc.input, tc.places)
			assert.Equal(t, tc.expectedOutput, rounded)
		})
	}
}
