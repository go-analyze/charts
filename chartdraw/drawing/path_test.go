package drawing

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type recordFlattener struct {
	moves []string
}

func (r *recordFlattener) MoveTo(x, y float64) {
	r.moves = append(r.moves, fmt.Sprintf("M%.1f,%.1f", x, y))
}

func (r *recordFlattener) LineTo(x, y float64) {
	r.moves = append(r.moves, fmt.Sprintf("L%.1f,%.1f", x, y))
}

func (r *recordFlattener) End() {}

func TestPathBasicOps(t *testing.T) {
	t.Parallel()

	p := &Path{}
	p.LineTo(1, 2)
	p.LineTo(3, 4)
	assert.Equal(t, []PathComponent{MoveToComponent, LineToComponent, LineToComponent}, p.Components)
	assert.InDeltaSlice(t, []float64{0, 0, 1, 2, 3, 4}, p.Points, 0.0001)
}

func TestPathArcTo(t *testing.T) {
	t.Parallel()

	p := &Path{}
	p.ArcTo(0, 0, 1, 1, 0, math.Pi/2)

	expectX := 0.0
	expectY := 1.0
	assert.InDelta(t, expectX, p.x, 0.0001)
	assert.InDelta(t, expectY, p.y, 0.0001)
	assert.InDeltaSlice(t, []float64{1, 0, 0, 0, 1, 1, 0, math.Pi / 2}, p.Points, 0.0001)
	assert.Equal(t, MoveToComponent, p.Components[0])
	assert.Equal(t, ArcToComponent, p.Components[1])
}

func TestTransformer(t *testing.T) {
	t.Parallel()

	rec := &recordFlattener{}
	tr := Transformer{Tr: NewTranslationMatrix(2, 3), Flattener: rec}
	tr.MoveTo(1, 1)
	tr.LineTo(2, 2)
	tr.End()

	assert.Equal(t, []string{"M3.0,4.0", "L4.0,5.0"}, rec.moves)
}

func TestPathCurveAndClose(t *testing.T) {
	t.Parallel()

	p := &Path{}
	p.QuadCurveTo(1, 1, 2, 2)
	p.CubicCurveTo(3, 3, 4, 4, 5, 5)
	p.Close()

	expComp := []PathComponent{MoveToComponent, QuadCurveToComponent, CubicCurveToComponent, CloseComponent}
	assert.Equal(t, expComp, p.Components)
	assert.InDeltaSlice(t, []float64{0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5}, p.Points, 0.0001)
	assert.InDelta(t, 5.0, p.x, 0.0001)
	assert.InDelta(t, 5.0, p.y, 0.0001)
}

func TestPathCopyClearIsEmpty(t *testing.T) {
	t.Parallel()

	p := &Path{}
	p.LineTo(1, 1)
	copyP := p.Copy()

	p.Clear()
	assert.True(t, p.IsEmpty())
	assert.False(t, copyP.IsEmpty())

	p2 := &Path{}
	assert.True(t, p2.IsEmpty())
}

func TestPathString(t *testing.T) {
	t.Parallel()

	p := &Path{}
	p.MoveTo(0, 0)
	p.LineTo(1, 1)
	p.QuadCurveTo(2, 2, 3, 3)
	p.CubicCurveTo(4, 4, 5, 5, 6, 6)
	p.Close()

	got := p.String()
	expect := "" +
		"MoveTo: 0.000000, 0.000000\n" +
		"LineTo: 1.000000, 1.000000\n" +
		"QuadCurveTo: 2.000000, 2.000000, 3.000000, 3.000000\n" +
		"CubicCurveTo: 4.000000, 4.000000, 5.000000, 5.000000, 6.000000, 6.000000\n" +
		"Close\n"
	assert.Equal(t, expect, got)
}

func TestPathLastPointAndArc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		ops   func(*Path)
		x, y  float64
		delta float64
	}{
		{
			name: "line and quad",
			ops: func(p *Path) {
				p.MoveTo(1, 1)
				p.LineTo(2, 2)
				p.QuadCurveTo(3, 3, 4, 4)
			},
			x: 4, y: 4, delta: 0,
		},
		{
			name: "arc positive",
			ops: func(p *Path) {
				p.ArcTo(0, 0, 1, 1, 0, math.Pi/2)
			},
			x: 0, y: 1, delta: 0.0001,
		},
		{
			name: "arc negative",
			ops: func(p *Path) {
				p.ArcTo(0, 0, 1, 1, math.Pi/2, -math.Pi/2)
			},
			x: 1, y: 0, delta: 0.0001,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := &Path{}
			tc.ops(p)
			x, y := p.LastPoint()
			assert.InDelta(t, tc.x, x, tc.delta)
			assert.InDelta(t, tc.y, y, tc.delta)
		})
	}
}

func TestRectangleHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		build  func() *Path
		rect   bool
		bounds [4]int
	}{
		{
			name: "rectangle",
			build: func() *Path {
				p := &Path{}
				p.MoveTo(0, 0)
				p.LineTo(2, 0)
				p.LineTo(2, 2)
				p.LineTo(0, 2)
				p.LineTo(0, 0)
				return p
			},
			rect:   true,
			bounds: [4]int{0, 0, 2, 2},
		},
		{
			name: "non-rectangle",
			build: func() *Path {
				p := &Path{}
				p.MoveTo(0, 0)
				p.LineTo(1, 1)
				return p
			},
			rect: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.build()
			if tc.rect {
				assert.True(t, isRectanglePath(p))
				x1, y1, x2, y2 := getRectangleBounds(p)
				assert.Equal(t, tc.bounds[0], x1)
				assert.Equal(t, tc.bounds[1], y1)
				assert.Equal(t, tc.bounds[2], x2)
				assert.Equal(t, tc.bounds[3], y2)
			} else {
				assert.False(t, isRectanglePath(p))
			}
		})
	}
}
