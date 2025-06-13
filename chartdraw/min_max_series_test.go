package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinSeriesEnsureMinValue(t *testing.T) {
	t.Parallel()

	vp := mockValuesProvider{
		X: LinearRange(1.0, 5.0),
		Y: []float64{4, 2, 5, 1, 3},
	}

	ms := &MinSeries{InnerSeries: vp}
	ms.ensureMinValue()

	if assert.NotNil(t, ms.minValue) {
		assert.InDelta(t, 1.0, *ms.minValue, 0)
	}

	_, y := ms.GetValues(0)
	assert.InDelta(t, 1.0, y, 0)
}

func TestMaxSeriesEnsureMaxValue(t *testing.T) {
	t.Parallel()

	vp := mockValuesProvider{
		X: LinearRange(1.0, 5.0),
		Y: []float64{4, 2, 5, 1, 3},
	}

	ms := &MaxSeries{InnerSeries: vp}
	ms.ensureMaxValue()

	if assert.NotNil(t, ms.maxValue) {
		assert.InDelta(t, 5.0, *ms.maxValue, 0)
	}

	_, y := ms.GetValues(0)
	assert.InDelta(t, 5.0, y, 0)
}
