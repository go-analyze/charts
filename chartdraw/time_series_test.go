package chartdraw

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeSeriesGetValue(t *testing.T) {
	t.Parallel()

	now := time.Now()
	ts := TimeSeries{
		Name: "Test",
		XValues: []time.Time{
			now.AddDate(0, 0, -5),
			now.AddDate(0, 0, -4),
			now.AddDate(0, 0, -3),
			now.AddDate(0, 0, -2),
			now.AddDate(0, 0, -1),
		},
		YValues: []float64{
			1.0, 2.0, 3.0, 4.0, 5.0,
		},
	}

	x0, y0 := ts.GetValues(0)
	assert.NotZero(t, x0)
	assert.InDelta(t, 1.0, y0, 0)
}

func TestTimeSeriesValidate(t *testing.T) {
	t.Parallel()

	now := time.Now()
	cs := TimeSeries{
		Name: "Test Series",
		XValues: []time.Time{
			now.AddDate(0, 0, -5),
			now.AddDate(0, 0, -4),
			now.AddDate(0, 0, -3),
			now.AddDate(0, 0, -2),
			now.AddDate(0, 0, -1),
		},
		YValues: []float64{
			1.0, 2.0, 3.0, 4.0, 5.0,
		},
	}
	require.NoError(t, cs.Validate())

	cs = TimeSeries{
		Name: "Test Series",
		XValues: []time.Time{
			now.AddDate(0, 0, -5),
			now.AddDate(0, 0, -4),
			now.AddDate(0, 0, -3),
			now.AddDate(0, 0, -2),
			now.AddDate(0, 0, -1),
		},
	}
	require.Error(t, cs.Validate())

	cs = TimeSeries{
		Name: "Test Series",
		YValues: []float64{
			1.0, 2.0, 3.0, 4.0, 5.0,
		},
	}
	require.Error(t, cs.Validate())
}
