package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConcatSeries(t *testing.T) {
	t.Parallel()

	s1 := ContinuousSeries{
		XValues: LinearRange(1.0, 10.0),
		YValues: LinearRange(1.0, 10.0),
	}

	s2 := ContinuousSeries{
		XValues: LinearRange(11, 20.0),
		YValues: LinearRange(10.0, 1.0),
	}

	s3 := ContinuousSeries{
		XValues: LinearRange(21, 30.0),
		YValues: LinearRange(1.0, 10.0),
	}

	cs := ConcatSeries([]Series{s1, s2, s3})
	assert.Equal(t, 30, cs.Len())

	x0, y0 := cs.GetValue(0)
	assert.InDelta(t, 1.0, x0, 0)
	assert.InDelta(t, 1.0, y0, 0)

	xm, ym := cs.GetValue(19)
	assert.InDelta(t, 20.0, xm, 0)
	assert.InDelta(t, 1.0, ym, 0)

	xn, yn := cs.GetValue(29)
	assert.InDelta(t, 30.0, xn, 0)
	assert.InDelta(t, 10.0, yn, 0)

	invalid := ConcatSeries([]Series{
		ContinuousSeries{
			XValues: LinearRange(1.0, 10.0),
		},
	})
	require.Error(t, invalid.Validate())
}
