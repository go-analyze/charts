package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, 1.0, x0)
	assert.Equal(t, 1.0, y0)

	xm, ym := cs.GetValue(19)
	assert.Equal(t, 20.0, xm)
	assert.Equal(t, 1.0, ym)

	xn, yn := cs.GetValue(29)
	assert.Equal(t, 30.0, xn)
	assert.Equal(t, 10.0, yn)
}
