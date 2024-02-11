package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func TestGetDefaultInt(t *testing.T) {
	require.Equal(t, 1, getDefaultInt(0, 1))
	require.Equal(t, 10, getDefaultInt(10, 1))
}

func TestCeilFloatToInt(t *testing.T) {
	require.Equal(t, 1, ceilFloatToInt(0.8))
	require.Equal(t, 1, ceilFloatToInt(1.0))
	require.Equal(t, 2, ceilFloatToInt(1.2))
}

func TestCommafWithDigits(t *testing.T) {
	require.Equal(t, "1.2", commafWithDigits(1.2))
	require.Equal(t, "1.21", commafWithDigits(1.21231))

	require.Equal(t, "1.20k", commafWithDigits(1200.121))
	require.Equal(t, "1.20M", commafWithDigits(1200000.121))
}

func TestAutoDivide(t *testing.T) {
	require.Equal(t, []int{
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
	style := chart.Style{
		FontSize: 10,
	}
	p.SetStyle(style)

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

func TestConvertPercent(t *testing.T) {
	assert.Equal(t, -1.0, convertPercent("1"))
	assert.Equal(t, -1.0, convertPercent("a%"))
	assert.Equal(t, 0.1, convertPercent("10%"))
}

func TestParseColor(t *testing.T) {
	c := parseColor("")
	assert.True(t, c.IsZero())

	c = parseColor("#333")
	assert.Equal(t, drawing.Color{R: 51, G: 51, B: 51, A: 255}, c)

	c = parseColor("#313233")
	assert.Equal(t, drawing.Color{R: 49, G: 50, B: 51, A: 255}, c)

	c = parseColor("rgb(31,32,33)")
	assert.Equal(t, drawing.Color{R: 31, G: 32, B: 33, A: 255}, c)

	c = parseColor("rgba(50,51,52,250)")
	assert.Equal(t, drawing.Color{R: 50, G: 51, B: 52, A: 250}, c)
}

func TestIsLightColor(t *testing.T) {
	assert.True(t, isLightColor(drawing.Color{R: 255, G: 255, B: 255}))
	assert.True(t, isLightColor(drawing.Color{R: 145, G: 204, B: 117}))

	assert.False(t, isLightColor(drawing.Color{R: 88, G: 112, B: 198}))
	assert.False(t, isLightColor(drawing.Color{R: 0, G: 0, B: 0}))
	assert.False(t, isLightColor(drawing.Color{R: 16, G: 12, B: 42}))
}
