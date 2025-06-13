package drawing

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type recordFloat struct {
	ops  []string
	xs   []float64
	ys   []float64
	ends int
}

func (r *recordFloat) MoveTo(x, y float64) {
	r.ops = append(r.ops, "M")
	r.xs = append(r.xs, x)
	r.ys = append(r.ys, y)
}

func (r *recordFloat) LineTo(x, y float64) {
	r.ops = append(r.ops, "L")
	r.xs = append(r.xs, x)
	r.ys = append(r.ys, y)
}

func (r *recordFloat) End() {
	r.ends++
}

func TestFlattenMixed(t *testing.T) {
	t.Parallel()

	p := &Path{}
	p.MoveTo(0, 0)
	p.LineTo(1, 0)
	p.QuadCurveTo(1.5, 0, 2, 0)
	p.CubicCurveTo(2.5, 0, 2.75, 0, 3, 0)
	p.ArcTo(4, 0, 1, 1, 0, math.Pi/2)
	p.Close()

	rec := &recordFloat{}
	Flatten(p, rec, 1.0)

	expectOps := []string{"M", "L", "L", "L", "L", "L", "L", "L", "L"}
	expectX := []float64{0, 1, 2, 3, 3, 5, 4.580247, 4, 0}
	expectY := []float64{0, 0, 0, 0, 0, 0, 0.814441, 1, 0}

	assert.Equal(t, 1, rec.ends)
	assert.Equal(t, expectOps, rec.ops)
	assert.InDeltaSlice(t, expectX, rec.xs, 0.0001)
	assert.InDeltaSlice(t, expectY, rec.ys, 0.0001)
}

func TestFlattenMultiMove(t *testing.T) {
	t.Parallel()

	p := &Path{}
	p.MoveTo(0, 0)
	p.LineTo(1, 0)
	p.MoveTo(2, 0)
	p.LineTo(3, 0)

	rec := &recordFloat{}
	Flatten(p, rec, 1.0)

	expectOps := []string{"M", "L", "M", "L"}
	expectX := []float64{0, 1, 2, 3}
	expectY := []float64{0, 0, 0, 0}

	assert.Equal(t, 2, rec.ends)
	assert.Equal(t, expectOps, rec.ops)
	assert.InDeltaSlice(t, expectX, rec.xs, 0.0001)
	assert.InDeltaSlice(t, expectY, rec.ys, 0.0001)
}

func TestSegmentedPathPoints(t *testing.T) {
	t.Parallel()

	segments := []struct {
		startX, startY float64
		endX, endY     float64
	}{
		{0, 0, 1, 1},
		{2, 2, 3, 3},
	}

	sp := &SegmentedPath{}
	for _, seg := range segments {
		sp.MoveTo(seg.startX, seg.startY)
		sp.LineTo(seg.endX, seg.endY)
		sp.End()
	}
	assert.InDeltaSlice(t, []float64{0, 0, 1, 1, 2, 2, 3, 3}, sp.Points, 0.0001)
}

func TestSegmentedPathEnd(t *testing.T) {
	t.Parallel()

	sp := &SegmentedPath{}
	sp.MoveTo(1, 1)
	sp.LineTo(2, 2)

	expect := append([]float64(nil), sp.Points...)
	sp.End()

	assert.InDeltaSlice(t, expect, sp.Points, 0.0001)
}
