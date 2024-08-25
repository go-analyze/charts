package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstValueAnnotation(t *testing.T) {
	series := ContinuousSeries{
		XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
		YValues: []float64{5.0, 3.0, 3.0, 2.0, 1.0},
	}

	fva := FirstValueAnnotation(series)
	assert.NotEmpty(t, fva.Annotations)
	fvaa := fva.Annotations[0]
	assert.Equal(t, float64(1), fvaa.XValue)
	assert.Equal(t, float64(5), fvaa.YValue)
}
