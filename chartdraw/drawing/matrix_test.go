package drawing

import (
	"math"
	"testing"

	"github.com/go-analyze/charts/chartdraw/matrix"
	"github.com/stretchr/testify/assert"
)

func TestMatrixTransformInverse(t *testing.T) {
	t.Parallel()

	m := NewTranslationMatrix(5, 7)
	x, y := m.TransformPoint(1, 2)
	assert.InDelta(t, 6.0, x, 0.0)
	assert.InDelta(t, 9.0, y, 0.0)

	ix, iy := m.InverseTransformPoint(x, y)
	assert.InDelta(t, 1.0, ix, matrix.DefaultEpsilon)
	assert.InDelta(t, 2.0, iy, matrix.DefaultEpsilon)
}

func TestMatrixComposeInverse(t *testing.T) {
	t.Parallel()

	m := NewScaleMatrix(2, 3)
	m.Translate(4, 5)
	m.Rotate(math.Pi / 3)

	inv := m.Copy()
	inv.Inverse()
	inv.Compose(m)

	assert.True(t, inv.IsIdentity())
}

func TestMatrixTransformRectangle(t *testing.T) {
	t.Parallel()

	m := NewScaleMatrix(2, 3)
	x0, y0, x1, y1 := m.TransformRectangle(1, 2, 3, 4)
	assert.InDelta(t, 2.0, x0, 0.0)
	assert.InDelta(t, 6.0, y0, 0.0)
	assert.InDelta(t, 6.0, x1, 0.0)
	assert.InDelta(t, 12.0, y1, 0.0)
}

func TestMatrixIdentityHelpers(t *testing.T) {
	t.Parallel()

	id := NewIdentityMatrix()
	assert.True(t, id.IsTranslation())
	assert.True(t, id.IsIdentity())

	tr := NewTranslationMatrix(2, 3)
	tx, ty := tr.GetTranslation()
	assert.InDelta(t, 2.0, tx, 0.0)
	assert.InDelta(t, 3.0, ty, 0.0)
}

func TestMatrixTransformSlice(t *testing.T) {
	t.Parallel()

	m := NewScaleMatrix(2, 3)
	m.Translate(4, -5)
	pts := []float64{1, 2, 3, 4}
	expect := append([]float64(nil), pts...)
	m.Transform(pts)
	m.InverseTransform(pts)
	assert.InDeltaSlice(t, expect, pts, matrix.DefaultEpsilon)
}

func TestMatrixVectorTransform(t *testing.T) {
	t.Parallel()

	m := NewRotationMatrix(math.Pi / 2)
	vec := []float64{1, 0}
	m.VectorTransform(vec)
	assert.InDeltaSlice(t, []float64{0, 1}, vec, matrix.DefaultEpsilon)
}

func TestMatrixScaleTranslateRotate(t *testing.T) {
	t.Parallel()

	m := NewIdentityMatrix()
	m.Scale(2, 3)
	sx, sy := m.GetScaling()
	assert.InDelta(t, 2.0, sx, 0.0)
	assert.InDelta(t, 3.0, sy, 0.0)

	m.Translate(4, 5)
	assert.InDeltaSlice(t, []float64{2, 0, 0, 3, 8, 15}, m[:], matrix.DefaultEpsilon)

	m.Rotate(math.Pi / 2)
	expected := Matrix{0, 3, -2, 0, 8, 15}
	assert.True(t, m.Equals(expected))
}

func TestMatrixGetScale(t *testing.T) {
	t.Parallel()

	m := NewScaleMatrix(2, 2)
	assert.InDelta(t, 2.0, m.GetScale(), matrix.DefaultEpsilon)

	m.Rotate(math.Pi / 4)
	assert.InDelta(t, 2.0, m.GetScale(), matrix.DefaultEpsilon)
}
