package drawing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type recordFlattenerEnd struct {
	moves []string
}

func (r *recordFlattenerEnd) MoveTo(x, y float64) {
	r.moves = append(r.moves, fmt.Sprintf("M%.1f,%.1f", x, y))
}

func (r *recordFlattenerEnd) LineTo(x, y float64) {
	r.moves = append(r.moves, fmt.Sprintf("L%.1f,%.1f", x, y))
}

func (r *recordFlattenerEnd) End() {
	r.moves = append(r.moves, "E")
}

func TestDashVertexConverterLineTo(t *testing.T) {
	t.Parallel()

	rec := &recordFlattenerEnd{}
	d := NewDashVertexConverter([]float64{2, 2}, 0, rec)
	d.MoveTo(0, 0)
	d.LineTo(5, 0)
	d.End()

	expect := []string{"M0.0,0.0", "L2.0,0.0", "E", "M4.0,0.0", "L5.0,0.0", "E"}
	assert.Equal(t, expect, rec.moves)
}
