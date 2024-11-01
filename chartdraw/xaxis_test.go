package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXAxisGetTicks(t *testing.T) {
	r, err := PNG(1024, 1024)
	require.NoError(t, err)

	f, err := GetDefaultFont()
	require.NoError(t, err)

	xa := XAxis{}
	xr := &ContinuousRange{Min: 10, Max: 100, Domain: 1024}
	styleDefaults := Style{
		FontStyle: FontStyle{
			Font:     f,
			FontSize: 10.0,
		},
	}
	vf := FloatValueFormatter
	ticks := xa.GetTicks(r, xr, styleDefaults, vf)
	assert.Len(t, ticks, 16)
}

func TestXAxisGetTicksWithUserDefaults(t *testing.T) {
	r, err := PNG(1024, 1024)
	require.NoError(t, err)

	f, err := GetDefaultFont()
	require.NoError(t, err)

	xa := XAxis{
		Ticks: []Tick{{Value: 1.0, Label: "1.0"}},
	}
	xr := &ContinuousRange{Min: 10, Max: 100, Domain: 1024}
	styleDefaults := Style{
		FontStyle: FontStyle{
			Font:     f,
			FontSize: 10.0,
		},
	}
	vf := FloatValueFormatter
	ticks := xa.GetTicks(r, xr, styleDefaults, vf)
	assert.Len(t, ticks, 1)
}

func TestXAxisMeasure(t *testing.T) {
	f, err := GetDefaultFont()
	require.NoError(t, err)
	style := Style{
		FontStyle: FontStyle{
			Font:     f,
			FontSize: 10.0,
		},
	}
	r, err := PNG(100, 100)
	require.NoError(t, err)
	ticks := []Tick{{Value: 1.0, Label: "1.0"}, {Value: 2.0, Label: "2.0"}, {Value: 3.0, Label: "3.0"}}
	xa := XAxis{}
	xab := xa.Measure(r, NewBox(0, 0, 100, 100), &ContinuousRange{Min: 1.0, Max: 3.0, Domain: 100}, style, ticks)
	assert.Equal(t, 122, xab.Width())
	assert.Equal(t, 21, xab.Height())
}
