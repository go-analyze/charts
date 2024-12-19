package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeqEach(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	values.Each(func(i int, v float64) {
		assert.Equal(t, float64(i), v-1)
	})
}

func TestSeqMap(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	mapped := values.Map(func(i int, v float64) float64 {
		assert.Equal(t, float64(i), v-1)
		return v * 2
	})
	assert.Equal(t, 4, mapped.Len())
}

func TestSeqFoldLeft(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	ten := values.FoldLeft(func(_ int, vp, v float64) float64 {
		return vp + v
	})
	assert.Equal(t, float64(10), ten)

	orderTest := Seq{NewArray(10, 3, 2, 1)}
	four := orderTest.FoldLeft(func(_ int, vp, v float64) float64 {
		return vp - v
	})
	assert.Equal(t, float64(4), four)
}

func TestSeqFoldRight(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	ten := values.FoldRight(func(_ int, vp, v float64) float64 {
		return vp + v
	})
	assert.Equal(t, float64(10), ten)

	orderTest := Seq{NewArray(10, 3, 2, 1)}
	notFour := orderTest.FoldRight(func(_ int, vp, v float64) float64 {
		return vp - v
	})
	assert.Equal(t, float64(-14), notFour)
}

func TestSeqSum(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	assert.Equal(t, float64(10), values.Sum())
}

func TestSeqAverage(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4)}
	assert.Equal(t, 2.5, values.Average())

	valuesOdd := Seq{NewArray(1, 2, 3, 4, 5)}
	assert.Equal(t, float64(3), valuesOdd.Average())
}

func TestSequenceVariance(t *testing.T) {
	t.Parallel()

	values := Seq{NewArray(1, 2, 3, 4, 5)}
	assert.Equal(t, float64(2), values.Variance())
}

func TestSequenceNormalize(t *testing.T) {
	t.Parallel()

	normalized := ValueSequence(1, 2, 3, 4, 5).Normalize().Values()

	assert.NotEmpty(t, normalized)
	require.Len(t, normalized, 5)
	assert.Equal(t, 0.0, normalized[0])
	assert.Equal(t, 0.25, normalized[1])
	assert.Equal(t, 1.0, normalized[4])
}

func TestLinearRange(t *testing.T) {
	t.Parallel()

	values := LinearRange(1, 100)
	require.Len(t, values, 100)
	assert.Equal(t, float64(1), values[0])
	assert.Equal(t, float64(100), values[99])
}

func TestLinearRangeWithStep(t *testing.T) {
	t.Parallel()

	values := LinearRangeWithStep(0, 100, 5)
	assert.Equal(t, float64(100), values[20])
	assert.Len(t, values, 21)
}

func TestLinearRangeReversed(t *testing.T) {
	t.Parallel()

	values := LinearRange(10.0, 1.0)
	require.Len(t, values, 10)
	assert.Equal(t, 10.0, values[0])
	assert.Equal(t, 1.0, values[9])
}

func TestLinearSequenceRegression(t *testing.T) {
	t.Parallel()

	// note; this assumes a 1.0 step is implicitly set in the constructor.
	linearProvider := NewLinearSequence().WithStart(1.0).WithEnd(100.0)
	assert.Equal(t, float64(1), linearProvider.Start())
	assert.Equal(t, float64(100), linearProvider.End())
	assert.Equal(t, 100, linearProvider.Len())

	values := Seq{linearProvider}.Values()
	require.Len(t, values, 100)
	assert.Equal(t, 1.0, values[0])
	assert.Equal(t, 100.0, values[99])
}
