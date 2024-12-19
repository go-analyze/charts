package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogRangeTranslate(t *testing.T) {
	t.Parallel()

	values := []float64{1, 10, 100, 1000, 10000, 100000, 1000000}
	r := LogarithmicRange{Domain: 1000}
	r.Min, r.Max = MinMax(values...)

	assert.Equal(t, 0, r.Translate(0))          // goes to bottom
	assert.Equal(t, 0, r.Translate(1))          // goes to bottom
	assert.Equal(t, 160, r.Translate(10))       // roughly 1/6th of max
	assert.Equal(t, 500, r.Translate(1000))     // roughly 1/2 of max (1.0e6 / 1.0e3)
	assert.Equal(t, 1000, r.Translate(1000000)) // max value
}

func TestGetTicks(t *testing.T) {
	t.Parallel()

	values := []float64{35, 512, 1525122}
	r := LogarithmicRange{Domain: 1000}
	r.Min, r.Max = MinMax(values...)

	ticks := r.GetTicks(nil, Style{}, FloatValueFormatter)
	require.Len(t, ticks, 7)
	assert.Equal(t, float64(10), ticks[0].Value)
	assert.Equal(t, float64(100), ticks[1].Value)
	assert.Equal(t, float64(10000000), ticks[6].Value)
}

func TestGetTicksFromHigh(t *testing.T) {
	t.Parallel()

	values := []float64{1412, 352144, 1525122} // min tick should be 1000
	r := LogarithmicRange{}
	r.Min, r.Max = MinMax(values...)

	ticks := r.GetTicks(nil, Style{}, FloatValueFormatter)
	require.Len(t, ticks, 5)
	assert.Equal(t, float64(1000), ticks[0].Value)
	assert.Equal(t, float64(10000), ticks[1].Value)
	assert.Equal(t, float64(10000000), ticks[4].Value)
}
