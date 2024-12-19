package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistogramSeries(t *testing.T) {
	t.Parallel()

	cs := ContinuousSeries{
		Name:    "Test Series",
		XValues: LinearRange(1.0, 20.0),
		YValues: LinearRange(10.0, -10.0),
	}

	hs := HistogramSeries{
		InnerSeries: cs,
	}

	for x := 0; x < hs.Len(); x++ {
		csx, csy := cs.GetValues(0)
		hsx, hsy1, hsy2 := hs.GetBoundedValues(0)
		assert.Equal(t, csx, hsx)
		assert.True(t, hsy1 > 0)
		assert.True(t, hsy2 <= 0)
		assert.True(t, csy < 0 || (csy > 0 && csy == hsy1))
		assert.True(t, csy > 0 || (csy < 0 && csy == hsy2))
	}
}
