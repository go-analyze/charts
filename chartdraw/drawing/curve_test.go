package drawing

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type point struct {
	X, Y float64
}

type mockLine struct {
	inner []point
}

func (ml *mockLine) LineTo(x, y float64) {
	ml.inner = append(ml.inner, point{x, y})
}

func (ml *mockLine) Len() int {
	return len(ml.inner)
}

func TestTraceQuad(t *testing.T) {
	t.Parallel()

	// Quad
	// x1, y1, cpx1, cpy2, x2, y2 float64
	// do the 9->12 circle segment
	quad := []float64{10, 20, 20, 20, 20, 10}
	liner := &mockLine{}
	TraceQuad(liner, quad, 0.5)
	assert.NotZero(t, liner.Len())
}

func TestSubdivideCubic(t *testing.T) {
	t.Parallel()

	cubic := []float64{0, 0, 0, 1, 1, 1, 1, 0}
	c1 := make([]float64, 8)
	c2 := make([]float64, 8)
	SubdivideCubic(cubic, c1, c2)

	expectC1 := []float64{0, 0, 0, 0.5, 0.25, 0.75, 0.5, 0.75}
	expectC2 := []float64{0.5, 0.75, 0.75, 0.75, 1, 0.5, 1, 0}
	assert.InDeltaSlice(t, expectC1, c1, 0.0001)
	assert.InDeltaSlice(t, expectC2, c2, 0.0001)
}

func TestTraceCubicAndArc(t *testing.T) {
	t.Parallel()

	cubic := []float64{0, 0, 0, 1, 1, 1, 1, 0}
	liner := &mockLine{}
	TraceCubic(liner, cubic, 0.1)
	last := liner.inner[len(liner.inner)-1]
	assert.InDelta(t, 1.0, last.X, 0.0001)
	assert.InDelta(t, 0.0, last.Y, 0.0001)

	liner = &mockLine{}
	lx, ly := TraceArc(liner, 0, 0, 1, 1, 0, math.Pi/2, 1)
	assert.InDelta(t, 0.0, lx, 0.0001)
	assert.InDelta(t, 1.0, ly, 0.0001)
	assert.NotZero(t, liner.Len())
}
