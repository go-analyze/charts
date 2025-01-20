package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateGridLines(t *testing.T) {
	t.Parallel()

	ticks := []Tick{
		{Value: 1.0, Label: "1.0"},
		{Value: 2.0, Label: "2.0"},
		{Value: 3.0, Label: "3.0"},
		{Value: 4.0, Label: "4.0"},
	}

	gl := GenerateGridLines(ticks, Style{}, Style{})
	require.Len(t, gl, 2)

	assert.InDelta(t, 2.0, gl[0].Value, 0)
	assert.InDelta(t, 3.0, gl[1].Value, 0)
}
