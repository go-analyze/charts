package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func TestIsLightColor(t *testing.T) {
	t.Parallel()

	assert.True(t, isLightColor(drawing.Color{R: 255, G: 255, B: 255}))
	assert.True(t, isLightColor(drawing.Color{R: 145, G: 204, B: 117}))

	assert.False(t, isLightColor(drawing.Color{R: 88, G: 112, B: 198}))
	assert.False(t, isLightColor(drawing.Color{R: 0, G: 0, B: 0}))
	assert.False(t, isLightColor(drawing.Color{R: 16, G: 12, B: 42}))
}

func TestParseColor(t *testing.T) {
	t.Parallel()

	c := parseColor("")
	assert.True(t, c.IsZero())

	c = parseColor("#333")
	assert.Equal(t, drawing.Color{R: 51, G: 51, B: 51, A: 255}, c)

	c = parseColor("#313233")
	assert.Equal(t, drawing.Color{R: 49, G: 50, B: 51, A: 255}, c)

	c = parseColor("rgb(31,32,33)")
	assert.Equal(t, drawing.Color{R: 31, G: 32, B: 33, A: 255}, c)

	c = parseColor("rgba(50,51,52,250)")
	assert.Equal(t, drawing.Color{R: 50, G: 51, B: 52, A: 250}, c)
}
