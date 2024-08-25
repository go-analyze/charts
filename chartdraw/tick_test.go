package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateContinuousTicks(t *testing.T) {
	f, err := GetDefaultFont()
	require.NoError(t, err)

	r, err := PNG(1024, 1024)
	require.NoError(t, err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:    0.0,
		Max:    10.0,
		Domain: 256,
	}

	vf := FloatValueFormatter

	ticks := GenerateContinuousTicks(r, ra, false, Style{}, vf)
	assert.NotEmpty(t, ticks)
	assert.Len(t, ticks, 11)
	assert.Equal(t, 0.0, ticks[0].Value)
	assert.Equal(t, 10.0, ticks[len(ticks)-1].Value)
}

func TestGenerateContinuousTicksDescending(t *testing.T) {
	f, err := GetDefaultFont()
	require.NoError(t, err)

	r, err := PNG(1024, 1024)
	require.NoError(t, err)
	r.SetFont(f)

	ra := &ContinuousRange{
		Min:        0.0,
		Max:        10.0,
		Domain:     256,
		Descending: true,
	}

	vf := FloatValueFormatter

	ticks := GenerateContinuousTicks(r, ra, false, Style{}, vf)
	assert.NotEmpty(t, ticks)
	assert.Len(t, ticks, 11)
	assert.Equal(t, 10.0, ticks[0].Value)
	assert.Equal(t, 9.0, ticks[1].Value)
	assert.Equal(t, 1.0, ticks[len(ticks)-2].Value)
	assert.Equal(t, 0.0, ticks[len(ticks)-1].Value)
}
