package charts

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDefaultInt(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 1, getDefaultInt(0, 1))
	assert.Equal(t, 10, getDefaultInt(10, 1))
}

func TestCeilFloatToInt(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 1, ceilFloatToInt(0.8))
	assert.Equal(t, 1, ceilFloatToInt(1.0))
	assert.Equal(t, 2, ceilFloatToInt(1.2))
	assert.Equal(t, math.MaxInt, ceilFloatToInt(math.MaxFloat64))
	assert.Equal(t, math.MaxInt, ceilFloatToInt(float64(math.MaxInt)+1))
	assert.Equal(t, math.MinInt, ceilFloatToInt(float64(math.MinInt)-1))
}

func TestFormatValueHumanizeShort(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "1", FormatValueHumanizeShort(1.2, 0, false))
	assert.Equal(t, "1.2", FormatValueHumanizeShort(1.2, 2, false))
	assert.Equal(t, "1.20", FormatValueHumanizeShort(1.2, 2, true))
	assert.Equal(t, "1.21", FormatValueHumanizeShort(1.21, 2, false))
	assert.Equal(t, "1.21", FormatValueHumanizeShort(1.21231, 2, false))
	assert.Equal(t, "1.22", FormatValueHumanizeShort(1.216, 2, false))

	assert.Equal(t, "1.2k", FormatValueHumanizeShort(1200.121, 2, false))
	assert.Equal(t, "1.20k", FormatValueHumanizeShort(1200.121, 2, true))
	assert.Equal(t, "1.21k", FormatValueHumanizeShort(1211.121, 2, false))
	assert.Equal(t, "1.22k", FormatValueHumanizeShort(1216.121, 2, false))
	assert.Equal(t, "1.2M", FormatValueHumanizeShort(1200000.121, 2, false))
	assert.Equal(t, "1.20M", FormatValueHumanizeShort(1200000.121, 2, true))
}

func TestFormatValueHumanize(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "1,234,567,890", FormatValueHumanize(1234567890, 2, false))
	assert.Equal(t, "1", FormatValueHumanize(1.2, 0, false))
	assert.Equal(t, "1.2", FormatValueHumanize(1.2, 2, false))
	assert.Equal(t, "1.21", FormatValueHumanize(1.21, 2, false))
	assert.Equal(t, "1.21", FormatValueHumanize(1.21231, 2, false))
	assert.Equal(t, "1.22", FormatValueHumanize(1.216, 2, false))
	assert.Equal(t, "1", FormatValueHumanize(1.216, -1, false)) // invalid decimal count reset to zero
	assert.Equal(t, "1,200.12", FormatValueHumanize(1200.121, 2, false))
	assert.Equal(t, "1,200.13", FormatValueHumanize(1200.126, 2, false))

	assert.Equal(t, "1.00", FormatValueHumanize(1, 2, true))
	assert.Equal(t, "1.20", FormatValueHumanize(1.2, 2, true))
}

func BenchmarkFormatValueHumanizeAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FormatValueHumanize(float64(i), 8, true)
		_ = FormatValueHumanize(float64(i*11), 8, true)
		_ = FormatValueHumanize(float64(i*41024320), 8, true)
		_ = FormatValueHumanize(1234567890, 8, true)
		_ = FormatValueHumanize(1234567890.2, 8, true)
		_ = FormatValueHumanize(123456789.012, 8, true)
		_ = FormatValueHumanize(122333444455555666666777777788888888999999999, 8, true)
		_ = FormatValueHumanize(122333444455555666666777777788888888999999999.012, 8, true)
		_ = FormatValueHumanize(122333444455555666666777777788888888999999999.987654321, 10, true)
	}
}

func BenchmarkFormatValueHumanizeTruncate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FormatValueHumanize(1234567890.21, 2, true)
		_ = FormatValueHumanize(1.21231, 2, true)
		_ = FormatValueHumanize(1200.126, 2, true)
		_ = FormatValueHumanize(1200.0123456789, 8, true)
		_ = FormatValueHumanize(122333444455555666666777777788888888999999999.0123456789, 8, true)
	}
}

func TestAutoDivide(t *testing.T) {
	assert.Equal(t, []int{
		0, 85, 171, 257, 342, 428, 514, 600,
	}, autoDivide(600, 7))
}

func TestGetRadius(t *testing.T) {
	t.Parallel()

	assert.InDelta(t, 50.0, getRadius(100, "50%"), 0)
	assert.InDelta(t, 30.0, getRadius(100, "30"), 0)
	assert.InDelta(t, 40.0, getRadius(100, ""), 0)
}

func TestReverseSlice(t *testing.T) {
	t.Parallel()

	arr := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	reverseSlice(arr)
	assert.Equal(t, []string{
		"Sun", "Sat", "Fri", "Thu", "Wed", "Tue", "Mon",
	}, arr)

	numbers := []int{1, 3, 5, 7, 9}
	reverseSlice(numbers)
	assert.Equal(t, []int{9, 7, 5, 3, 1}, numbers)
}

func TestSliceToFloat64(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		expected := []float64{1.0, 2.0, 3.0, 4.0}
		result := SliceToFloat64(input, func(i int) float64 { return float64(i) })
		assert.Equal(t, expected, result)
	})
	t.Run("string", func(t *testing.T) {
		input := []string{"1.5", "2.5", "3.5"}
		expected := []float64{1.5, 2.5, 3.5}
		result := SliceToFloat64(input, func(s string) float64 {
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				return f
			}
			return 0
		})
		assert.Equal(t, expected, result)
	})
	t.Run("empty", func(t *testing.T) {
		input := []string{}
		expected := []float64{}
		result := SliceToFloat64(input, func(s string) float64 { return 0 })
		assert.Equal(t, expected, result)
	})
	t.Run("nil", func(t *testing.T) {
		var input []int
		expected := []float64{}
		result := SliceToFloat64(input, func(i int) float64 { return float64(i) })
		assert.Equal(t, expected, result)
	})
}

// TestIntSliceToFloat64 tests the IntSliceToFloat64 convenience function.
func TestIntSliceToFloat64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []int
		want  []float64
	}{
		{
			name:  "positive",
			input: []int{1, 2, 3},
			want:  []float64{1.0, 2.0, 3.0},
		},
		{
			name:  "negative",
			input: []int{0, -1, -2},
			want:  []float64{0.0, -1.0, -2.0},
		},
		{
			name:  "empty",
			input: []int{},
			want:  []float64{},
		},
		{
			name:  "nil",
			input: nil,
			want:  []float64{},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IntSliceToFloat64(tt.input))
		})
	}
}

func TestParseFlexibleValue(t *testing.T) {
	t.Parallel()

	t.Run("percent", func(t *testing.T) {
		result, err := parseFlexibleValue("10%", 200)
		require.NoError(t, err)
		assert.InDelta(t, 20.0, result, 0)
	})
	t.Run("value", func(t *testing.T) {
		result, err := parseFlexibleValue("10", 200)
		require.NoError(t, err)
		assert.InDelta(t, 10.0, result, 0)
	})
	t.Run("error_percent", func(t *testing.T) {
		result, err := parseFlexibleValue("a%", 100)
		require.Error(t, err)
		assert.InDelta(t, 0.0, result, 0)
	})
	t.Run("error_val", func(t *testing.T) {
		result, err := parseFlexibleValue("a", 100)
		require.Error(t, err)
		assert.InDelta(t, 0.0, result, 0)
	})
}
