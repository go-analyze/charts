package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/matrix"
)

func TestPolynomialRegressionSeries_GetValues(t *testing.T) {
	t.Parallel()

	xs := []float64{0, 1, 2, 3}
	ys := []float64{5, 10, 19, 32}

	series := ContinuousSeries{
		XValues: xs,
		YValues: ys,
	}

	prs := &PolynomialRegressionSeries{
		InnerSeries: series,
		Degree:      2,
	}

	for i, x := range xs {
		_, y := prs.GetValues(i)
		expect := 5 + 3*x + 2*x*x
		assert.InDelta(t, expect, y, matrix.DefaultEpsilon)
	}

	require.Error(t, (&PolynomialRegressionSeries{}).Validate())
	require.Error(t, (&PolynomialRegressionSeries{
		InnerSeries: series,
		Offset:      len(xs),
	}).Validate())
	require.Error(t, (&PolynomialRegressionSeries{
		InnerSeries: series,
		Degree:      len(xs),
	}).Validate())
}

func TestPolynomialRegressionSeries_InvalidInputSafe(t *testing.T) {
	t.Parallel()

	series := ContinuousSeries{
		XValues: []float64{0, 1, 2, 3},
		YValues: []float64{5, 10, 19, 32},
	}

	invalidDegree := &PolynomialRegressionSeries{
		InnerSeries: series,
		Degree:      len(series.XValues),
	}
	x, y := invalidDegree.GetValues(0)
	assert.InDelta(t, 0.0, x, 0.0)
	assert.InDelta(t, 0.0, y, 0.0)

	invalidOffset := &PolynomialRegressionSeries{
		InnerSeries: series,
		Offset:      len(series.XValues),
		Degree:      1,
	}
	x, y = invalidOffset.GetFirstValues()
	assert.InDelta(t, 0.0, x, 0.0)
	assert.InDelta(t, 0.0, y, 0.0)
	x, y = invalidOffset.GetLastValues()
	assert.InDelta(t, 0.0, x, 0.0)
	assert.InDelta(t, 0.0, y, 0.0)
}

func TestPolynomialRegressionSeries_IndexClamped(t *testing.T) {
	t.Parallel()

	series := ContinuousSeries{
		XValues: []float64{0, 1, 2, 3},
		YValues: []float64{5, 10, 19, 32},
	}
	prs := &PolynomialRegressionSeries{
		InnerSeries: series,
		Degree:      2,
	}

	// Out-of-range lookup should safely clamp to the last sample.
	x, _ := prs.GetValues(999)
	assert.InDelta(t, 3.0, x, matrix.DefaultEpsilon)
}
