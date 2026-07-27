package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func TestJet(t *testing.T) {
	t.Parallel()

	t.Run("zero_range", func(t *testing.T) {
		assert.Equal(t, drawing.Color{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, Jet(5, 5, 5))
	})
	t.Run("mid_range", func(t *testing.T) {
		assert.Equal(t, drawing.Color{R: 0x00, G: 0xff, B: 0x00, A: 0xff}, Jet(0.5, 0, 1))
	})
}
