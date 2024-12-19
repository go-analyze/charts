package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPoly(t *testing.T) {
	t.Parallel()

	var xGiven = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var yGiven = []float64{1, 6, 17, 34, 57, 86, 121, 162, 209, 262, 321}
	var degree = 2

	c, err := Poly(xGiven, yGiven, degree)
	require.NoError(t, err)
	assert.Len(t, c, 3)

	assert.InDelta(t, c[0], 0.999999999, DefaultEpsilon)
	assert.InDelta(t, c[1], 2, DefaultEpsilon)
	assert.InDelta(t, c[2], 3, DefaultEpsilon)
}
