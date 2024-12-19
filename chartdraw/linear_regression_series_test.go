package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinearRegressionSeries(t *testing.T) {
	t.Parallel()

	mainSeries := ContinuousSeries{
		Name:    "A test series",
		XValues: LinearRange(1.0, 100.0),
		YValues: LinearRange(1.0, 100.0),
	}

	linRegSeries := &LinearRegressionSeries{
		InnerSeries: mainSeries,
	}

	lrx0, lry0 := linRegSeries.GetValues(0)
	assert.InDelta(t, 1.0, lrx0, 0.0000001)
	assert.InDelta(t, 1.0, lry0, 0.0000001)

	lrxn, lryn := linRegSeries.GetLastValues()
	assert.InDelta(t, 100.0, lrxn, 0.0000001)
	assert.InDelta(t, 100.0, lryn, 0.0000001)
}

func TestLinearRegressionSeriesDesc(t *testing.T) {
	t.Parallel()

	mainSeries := ContinuousSeries{
		Name:    "A test series",
		XValues: LinearRange(100.0, 1.0),
		YValues: LinearRange(100.0, 1.0),
	}

	linRegSeries := &LinearRegressionSeries{
		InnerSeries: mainSeries,
	}

	lrx0, lry0 := linRegSeries.GetValues(0)
	assert.InDelta(t, 100.0, lrx0, 0.0000001)
	assert.InDelta(t, 100.0, lry0, 0.0000001)

	lrxn, lryn := linRegSeries.GetLastValues()
	assert.InDelta(t, 1.0, lrxn, 0.0000001)
	assert.InDelta(t, 1.0, lryn, 0.0000001)
}

func TestLinearRegressionSeriesWindowAndOffset(t *testing.T) {
	t.Parallel()

	mainSeries := ContinuousSeries{
		Name:    "A test series",
		XValues: LinearRange(100.0, 1.0),
		YValues: LinearRange(100.0, 1.0),
	}

	linRegSeries := &LinearRegressionSeries{
		InnerSeries: mainSeries,
		Offset:      10,
		Limit:       10,
	}

	assert.Equal(t, 10, linRegSeries.Len())

	lrx0, lry0 := linRegSeries.GetValues(0)
	assert.InDelta(t, 90.0, lrx0, 0.0000001)
	assert.InDelta(t, 90.0, lry0, 0.0000001)

	lrxn, lryn := linRegSeries.GetLastValues()
	assert.InDelta(t, 80.0, lrxn, 0.0000001)
	assert.InDelta(t, 80.0, lryn, 0.0000001)
}
