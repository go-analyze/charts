package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLastValueAnnotationSeries(t *testing.T) {
	series := ContinuousSeries{
		XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
		YValues: []float64{5.0, 3.0, 3.0, 2.0, 1.0},
	}

	lva := LastValueAnnotationSeries(series)
	assert.NotEmpty(t, lva.Annotations)
	lvaa := lva.Annotations[0]
	assert.Equal(t, float64(5), lvaa.XValue)
	assert.Equal(t, float64(1), lvaa.YValue)
}
