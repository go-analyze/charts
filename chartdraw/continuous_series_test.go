package chartdraw

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContinuousSeries(t *testing.T) {
	t.Parallel()

	cs := ContinuousSeries{
		Name:    "Test Series",
		XValues: LinearRange(1.0, 10.0),
		YValues: LinearRange(1.0, 10.0),
	}

	assert.Equal(t, "Test Series", cs.GetName())
	assert.Equal(t, 10, cs.Len())
	x0, y0 := cs.GetValues(0)
	assert.Equal(t, 1.0, x0)
	assert.Equal(t, 1.0, y0)

	xn, yn := cs.GetValues(9)
	assert.Equal(t, 10.0, xn)
	assert.Equal(t, 10.0, yn)

	xn, yn = cs.GetLastValues()
	assert.Equal(t, 10.0, xn)
	assert.Equal(t, 10.0, yn)
}

func TestContinuousSeriesValueFormatter(t *testing.T) {
	t.Parallel()

	cs := ContinuousSeries{
		XValueFormatter: func(v interface{}) string {
			return fmt.Sprintf("%f foo", v)
		},
		YValueFormatter: func(v interface{}) string {
			return fmt.Sprintf("%f bar", v)
		},
	}

	xf, yf := cs.GetValueFormatters()
	assert.Equal(t, "0.100000 foo", xf(0.1))
	assert.Equal(t, "0.100000 bar", yf(0.1))
}

func TestContinuousSeriesValidate(t *testing.T) {
	t.Parallel()

	cs := ContinuousSeries{
		Name:    "Test Series",
		XValues: LinearRange(1.0, 10.0),
		YValues: LinearRange(1.0, 10.0),
	}
	require.NoError(t, cs.Validate())

	cs = ContinuousSeries{
		Name:    "Test Series",
		XValues: LinearRange(1.0, 10.0),
	}
	require.Error(t, cs.Validate())

	cs = ContinuousSeries{
		Name:    "Test Series",
		YValues: LinearRange(1.0, 10.0),
	}
	require.Error(t, cs.Validate())
}
