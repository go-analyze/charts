package chartdraw

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBollingerBandSeries(t *testing.T) {
	s1 := mockValuesProvider{
		X: LinearRange(1.0, 100.0),
		Y: RandomValuesWithMax(100, 1024),
	}

	bbs := &BollingerBandsSeries{
		InnerSeries: s1,
	}

	xvalues := make([]float64, 100)
	y1values := make([]float64, 100)
	y2values := make([]float64, 100)

	for x := 0; x < 100; x++ {
		xvalues[x], y1values[x], y2values[x] = bbs.GetBoundedValues(x)
	}

	for x := bbs.GetPeriod(); x < 100; x++ {
		assert.True(t, y1values[x] > y2values[x], fmt.Sprintf("%v vs. %v", y1values[x], y2values[x]))
	}
}

func TestBollingerBandLastValue(t *testing.T) {
	s1 := mockValuesProvider{
		X: LinearRange(1.0, 100.0),
		Y: LinearRange(1.0, 100.0),
	}

	bbs := &BollingerBandsSeries{
		InnerSeries: s1,
	}

	x, y1, y2 := bbs.GetBoundedLastValues()
	assert.Equal(t, 100.0, x)
	assert.Equal(t, float64(101), math.Floor(y1))
	assert.Equal(t, float64(83), math.Floor(y2))
}
