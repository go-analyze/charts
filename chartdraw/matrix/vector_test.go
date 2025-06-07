package matrix

import "testing"

import "github.com/stretchr/testify/assert"

func TestVectorDotProduct(t *testing.T) {
	t.Parallel()

	v1 := Vector{1, 2, 3}
	v2 := Vector{4, 5, 6}

	result, err := v1.DotProduct(v2)
	assert.NoError(t, err)
	assert.InDelta(t, float64(32), result, 0)
}

func TestVectorDotProductDimensionMismatch(t *testing.T) {
	t.Parallel()

	_, err := Vector{1, 2}.DotProduct(Vector{1})
	assert.ErrorIs(t, err, ErrDimensionMismatch)
}
