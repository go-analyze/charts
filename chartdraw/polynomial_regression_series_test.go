package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
}
