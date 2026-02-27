package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViridis(t *testing.T) {
	t.Parallel()

	t.Run("normal_range", func(t *testing.T) {
		c := Viridis(0.5, 0, 1)
		assert.Equal(t, viridisColors[127], c)
	})
	t.Run("equal_min_max", func(t *testing.T) {
		c := Viridis(5, 5, 5)
		assert.Equal(t, viridisColors[0], c)
	})
	t.Run("below_min", func(t *testing.T) {
		c := Viridis(-1, 0, 10)
		assert.Equal(t, viridisColors[0], c)
	})
	t.Run("above_max", func(t *testing.T) {
		c := Viridis(20, 0, 10)
		assert.Equal(t, viridisColors[255], c)
	})
}
