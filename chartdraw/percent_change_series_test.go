package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPercentageDifferenceSeries(t *testing.T) {
	t.Parallel()

	cs := ContinuousSeries{
		XValues: LinearRange(1.0, 10.0),
		YValues: LinearRange(1.0, 10.0),
	}

	pcs := PercentChangeSeries{
		Name:        "Test Series",
		InnerSeries: cs,
	}

	assert.Equal(t, "Test Series", pcs.GetName())
	assert.Equal(t, 10, pcs.Len())
	x0, y0 := pcs.GetValues(0)
	assert.InDelta(t, 1.0, x0, 0)
	assert.InDelta(t, 0.0, y0, 0)

	xn, yn := pcs.GetValues(9)
	assert.InDelta(t, 10.0, xn, 0)
	assert.InDelta(t, 9.0, yn, 0)

	xn, yn = pcs.GetLastValues()
	assert.InDelta(t, 10.0, xn, 0)
	assert.InDelta(t, 9.0, yn, 0)
}
