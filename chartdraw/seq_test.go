package chartdraw

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeqEach(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	values.Each(func(i int, v float64) {
		assert.InDelta(t, float64(i), v-1, 0)
	})
}

func TestSeqMap(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	mapped := values.Map(func(i int, v float64) float64 {
		assert.InDelta(t, float64(i), v-1, 0)
		return v * 2
	})
	assert.Equal(t, []float64{2, 4, 6, 8}, mapped.Values())
	assert.Equal(t, 4, mapped.Len())

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		empty := Seq{}.Map(func(i int, v float64) float64 { return v * 2 })
		assert.Empty(t, empty.Values())
	})

	t.Run("single_value", func(t *testing.T) {
		t.Parallel()
		mapped := ValueSequence(7).Map(func(i int, v float64) float64 { return v + 3 })
		assert.Equal(t, []float64{10}, mapped.Values())
	})
}

func TestSeqFoldLeft(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	ten := values.FoldLeft(func(_ int, vp, v float64) float64 {
		return vp + v
	})
	assert.InDelta(t, float64(10), ten, 0)

	orderTest := Seq{NewArray(10, 3, 2, 1)}
	four := orderTest.FoldLeft(func(_ int, vp, v float64) float64 {
		return vp - v
	})
	assert.InDelta(t, float64(4), four, 0)
}

func TestSeqFoldRight(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	ten := values.FoldRight(func(_ int, vp, v float64) float64 {
		return vp + v
	})
	assert.InDelta(t, float64(10), ten, 0)

	orderTest := Seq{NewArray(10, 3, 2, 1)}
	notFour := orderTest.FoldRight(func(_ int, vp, v float64) float64 {
		return vp - v
	})
	assert.InDelta(t, float64(-14), notFour, 0)
}

func TestSeqSum(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	assert.InDelta(t, float64(10), values.Sum(), 0)
}

func TestSeqAverage(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	assert.InDelta(t, 2.5, values.Average(), 0)

	valuesOdd := Seq{NewArray(1, 2, 3, 4, 5)}
	assert.InDelta(t, float64(3), valuesOdd.Average(), 0)
}

func TestSequenceVariance(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4, 5)}
	assert.InDelta(t, float64(2), values.Variance(), 0)
}

func TestSequenceNormalize(t *testing.T) {
	t.Parallel()

	normalized := ValueSequence(1, 2, 3, 4, 5).Normalize().Values()

	assert.NotEmpty(t, normalized)
	require.Len(t, normalized, 5)
	assert.InDelta(t, 0.0, normalized[0], 0)
	assert.InDelta(t, 0.25, normalized[1], 0)
	assert.InDelta(t, 1.0, normalized[4], 0)
}

func TestLinearRange(t *testing.T) {
	t.Parallel()

	values := LinearRange(1, 100)
	require.Len(t, values, 100)
	assert.InDelta(t, float64(1), values[0], 0)
	assert.InDelta(t, float64(100), values[99], 0)
}

func TestLinearRangeWithStep(t *testing.T) {
	t.Parallel()

	values := LinearRangeWithStep(0, 100, 5)
	assert.InDelta(t, float64(100), values[20], 0)
	assert.Len(t, values, 21)
}

func TestLinearRangeReversed(t *testing.T) {
	t.Parallel()

	values := LinearRange(10.0, 1.0)
	require.Len(t, values, 10)
	assert.InDelta(t, 10.0, values[0], 0)
	assert.InDelta(t, 1.0, values[9], 0)
}

func TestLinearSequenceRegression(t *testing.T) {
	t.Parallel()

	// note; this assumes a 1.0 step is implicitly set in the constructor.
	linearProvider := NewLinearSequence().WithStart(1.0).WithEnd(100.0)
	assert.InDelta(t, float64(1), linearProvider.Start(), 0)
	assert.InDelta(t, float64(100), linearProvider.End(), 0)
	assert.Equal(t, 100, linearProvider.Len())

	values := Seq{linearProvider}.Values()
	require.Len(t, values, 100)
	assert.InDelta(t, 1.0, values[0], 0)
	assert.InDelta(t, 100.0, values[99], 0)
}

func TestSeqPercentileOutOfRange(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	assert.InDelta(t, 1.0, values.Percentile(-0.1), 0)
	assert.InDelta(t, 1.0, values.Percentile(math.Inf(-1)), 0)
	assert.InDelta(t, 4.0, values.Percentile(1.0), 0)
	assert.InDelta(t, 4.0, values.Percentile(1.5), 0)
	assert.InDelta(t, 4.0, values.Percentile(math.Inf(1)), 0)
	assert.InDelta(t, 0.0, values.Percentile(math.NaN()), 0)
}

func TestSeqMedian(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"empty", nil, 0},
		{"single_value", []float64{5}, 5},
		{"two_values", []float64{1, 2}, 1.5},
		{"odd_length", []float64{3, 1, 2}, 2},
		{"even_length", []float64{1, 2, 3, 4}, 2.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, ValueSequence(tt.values...).Median(), 0)
		})
	}

	t.Run("nil_interface", func(t *testing.T) {
		t.Parallel()
		// Seq{} has a truly nil embedded interface; must not panic
		var s Seq
		assert.InDelta(t, 0.0, s.Median(), 0)
	})
}
