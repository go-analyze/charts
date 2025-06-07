package drawing

import (
	"fmt"
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
)

type recordBuilder struct{ ops []string }

func (r *recordBuilder) LastPoint() (float64, float64) { return 0, 0 }
func (r *recordBuilder) MoveTo(x, y float64)           { r.ops = append(r.ops, fmt.Sprintf("M%.1f,%.1f", x, y)) }
func (r *recordBuilder) LineTo(x, y float64)           { r.ops = append(r.ops, fmt.Sprintf("L%.1f,%.1f", x, y)) }
func (r *recordBuilder) QuadCurveTo(cx, cy, x, y float64) {
	r.ops = append(r.ops, fmt.Sprintf("Q%.1f,%.1f,%.1f,%.1f", cx, cy, x, y))
}
func (r *recordBuilder) CubicCurveTo(cx1, cy1, cx2, cy2, x, y float64)   {}
func (r *recordBuilder) ArcTo(cx, cy, rx, ry, startAngle, angle float64) {}
func (r *recordBuilder) Close()                                          {}

func TestDrawContour(t *testing.T) {
	t.Parallel()

	contour := []truetype.Point{
		{X: 0, Y: 0, Flags: 0x01},
		{X: 64, Y: 0, Flags: 0x01},
		{X: 64, Y: 64, Flags: 0x00},
		{X: 0, Y: 64, Flags: 0x01},
	}

	rec := &recordBuilder{}
	DrawContour(rec, contour, 0, 0)

	expect := []string{
		"M0.0,0.0",
		"L1.0,0.0",
		"Q1.0,-1.0,0.0,-1.0",
		"L0.0,0.0",
	}

	assert.Equal(t, expect, rec.ops)
}
