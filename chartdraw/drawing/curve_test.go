package drawing

import (
	"math"
	"strconv"
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

func TestSubdivideQuadAndHelpers(t *testing.T) {
	t.Parallel()

	t.Run("SubdivideQuad", func(t *testing.T) {
		tests := []struct {
			name    string
			quad    []float64
			expect1 []float64
			expect2 []float64
		}{
			{
				name:    "simple",
				quad:    []float64{0, 0, 1, 1, 2, 0},
				expect1: []float64{0, 0, 0.5, 0.5, 1, 0.5},
				expect2: []float64{1, 0.5, 1.5, 0.5, 2, 0},
			},
			{
				name:    "offset",
				quad:    []float64{0, 1, 1, 0, 2, 1},
				expect1: []float64{0, 1, 0.5, 0.5, 1, 0.5},
				expect2: []float64{1, 0.5, 1.5, 0.5, 2, 1},
			},
		}

		for i, tt := range tests {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				c1 := make([]float64, 6)
				c2 := make([]float64, 6)
				SubdivideQuad(tt.quad, c1, c2)
				assert.InDeltaSlice(t, tt.expect1, c1, 0.0001)
				assert.InDeltaSlice(t, tt.expect2, c2, 0.0001)
			})
		}
	})

	t.Run("traceWindowIndices", func(t *testing.T) {
		tests := []struct {
			idx        int
			start, end int
		}{
			{0, 0, 6},
			{1, 6, 12},
			{2, 12, 18},
		}

		for i, tt := range tests {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				s, e := traceWindowIndices(tt.idx)
				assert.Equal(t, tt.start, s)
				assert.Equal(t, tt.end, e)
			})
		}
	})

	t.Run("traceCalcDeltas", func(t *testing.T) {
		tests := []struct {
			name      string
			c         []float64
			dx, dy, d float64
		}{
			{"flat", []float64{0, 0, 1, 1, 2, 0}, 2, 0, 2},
			{"diag", []float64{0, 0, 1, 2, 2, 2}, 2, 2, 2},
		}

		for i, tt := range tests {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				dx, dy, d := traceCalcDeltas(tt.c)
				assert.InDelta(t, tt.dx, dx, 0.0001)
				assert.InDelta(t, tt.dy, dy, 0.0001)
				assert.InDelta(t, tt.d, d, 0.0001)
			})
		}
	})

	t.Run("traceIsFlat", func(t *testing.T) {
		tests := []struct {
			name                 string
			dx, dy, d, threshold float64
			expect               bool
		}{
			{"true", 2, 0, 1, 1, true},
			{"false", 2, 0, 3, 0.5, false},
		}

		for i, tt := range tests {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				got := traceIsFlat(tt.dx, tt.dy, tt.d, tt.threshold)
				assert.Equal(t, tt.expect, got)
			})
		}
	})

	t.Run("traceGetWindow", func(t *testing.T) {
		curves := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
		tests := []struct {
			idx    int
			expect []float64
		}{
			{0, []float64{0, 1, 2, 3, 4, 5}},
			{1, []float64{6, 7, 8, 9, 10, 11}},
		}

		for i, tt := range tests {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				window := traceGetWindow(curves, tt.idx)
				assert.Equal(t, tt.expect, window)
			})
		}
	})
}
