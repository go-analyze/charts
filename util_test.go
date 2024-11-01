package charts

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw"
)

func TestGetDefaultInt(t *testing.T) {
	assert.Equal(t, 1, getDefaultInt(0, 1))
	assert.Equal(t, 10, getDefaultInt(10, 1))
}

func TestCeilFloatToInt(t *testing.T) {
	assert.Equal(t, 1, ceilFloatToInt(0.8))
	assert.Equal(t, 1, ceilFloatToInt(1.0))
	assert.Equal(t, 2, ceilFloatToInt(1.2))
	assert.Equal(t, math.MaxInt, ceilFloatToInt(math.MaxFloat64))
	assert.Equal(t, math.MaxInt, ceilFloatToInt(float64(math.MaxInt)+1))
	assert.Equal(t, math.MinInt, ceilFloatToInt(float64(math.MinInt)-1))
}

func TestCommafWithDigits(t *testing.T) {
	assert.Equal(t, "1.2", commafWithDigits(1.2))
	assert.Equal(t, "1.21", commafWithDigits(1.21231))

	assert.Equal(t, "1.20k", commafWithDigits(1200.121))
	assert.Equal(t, "1.20M", commafWithDigits(1200000.121))
}

func TestAutoDivide(t *testing.T) {
	assert.Equal(t, []int{
		0,
		85,
		171,
		257,
		342,
		428,
		514,
		600,
	}, autoDivide(600, 7))
}

func TestSumInt(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		assert.Equal(t, 3, sumInt([]int{1, 2}))
	})
	t.Run("overflow-add", func(t *testing.T) {
		assert.Equal(t, math.MaxInt, sumInt([]int{1, math.MaxInt}))
		assert.Equal(t, math.MaxInt, sumInt([]int{1, math.MaxInt - 1}))
		assert.Equal(t, math.MaxInt, sumInt([]int{math.MaxInt, math.MaxInt}))
	})
	t.Run("overflow-sub", func(t *testing.T) {
		assert.Equal(t, math.MinInt, sumInt([]int{-1, math.MinInt}))
		assert.Equal(t, math.MinInt, sumInt([]int{-1, math.MinInt + 1}))
		assert.Equal(t, math.MinInt, sumInt([]int{math.MinInt, math.MinInt}))
	})
}

func TestGetRadius(t *testing.T) {
	assert.Equal(t, 50.0, getRadius(100, "50%"))
	assert.Equal(t, 30.0, getRadius(100, "30"))
	assert.Equal(t, 40.0, getRadius(100, ""))
}

func TestMeasureTextMaxWidthHeight(t *testing.T) {
	p, err := NewPainter(PainterOptions{
		Width:  400,
		Height: 300,
	})
	require.NoError(t, err)
	style := FontStyle{
		FontSize: 10,
	}
	p.SetStyle(chartdraw.Style{FontStyle: style})

	maxWidth, maxHeight := measureTextMaxWidthHeight([]string{
		"Mon",
		"Tue",
		"Wed",
		"Thu",
		"Fri",
		"Sat",
		"Sun",
	}, p)
	assert.Equal(t, 31, maxWidth)
	assert.Equal(t, 12, maxHeight)
}

func TestReverseSlice(t *testing.T) {
	arr := []string{
		"Mon",
		"Tue",
		"Wed",
		"Thu",
		"Fri",
		"Sat",
		"Sun",
	}
	reverseStringSlice(arr)
	assert.Equal(t, []string{
		"Sun",
		"Sat",
		"Fri",
		"Thu",
		"Wed",
		"Tue",
		"Mon",
	}, arr)

	numbers := []int{1, 3, 5, 7, 9}
	reverseIntSlice(numbers)
	assert.Equal(t, []int{9, 7, 5, 3, 1}, numbers)
}

func TestParseFlexibleValue(t *testing.T) {
	t.Run("percent", func(t *testing.T) {
		result, err := parseFlexibleValue("10%", 200)
		assert.NoError(t, err)
		assert.Equal(t, 20.0, result)
	})
	t.Run("value", func(t *testing.T) {
		result, err := parseFlexibleValue("10", 200)
		assert.NoError(t, err)
		assert.Equal(t, 10.0, result)
	})
}

func TestConvertPercent(t *testing.T) {
	verifyConvertPercent(t, -1.0, "1")
	verifyConvertPercent(t, -1.0, "a%")
	verifyConvertPercent(t, 0.1, "10%")
}

func verifyConvertPercent(t *testing.T, expected float64, input string) {
	t.Helper()

	v, err := convertPercent(input)
	if expected == -1 {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
	assert.Equal(t, expected, v)
}
