package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-analyze/charts/chartdraw/matrix"
)

func TestPolynomialRegression(t *testing.T) {
	t.Parallel()

	var xv []float64
	var yv []float64

	for i := 0; i < 100; i++ {
		xv = append(xv, float64(i))
		yv = append(yv, float64(i*i))
	}

	values := ContinuousSeries{
		XValues: xv,
		YValues: yv,
	}

	poly := &PolynomialRegressionSeries{
		InnerSeries: values,
		Degree:      2,
	}

	for i := 0; i < 100; i++ {
		_, y := poly.GetValues(i)
		assert.InDelta(t, float64(i*i), y, matrix.DefaultEpsilon)
	}
}
