package drawing

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidDash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		dash     []float64
		expected bool
	}{
		{name: "normal", dash: []float64{5, 5}, expected: true},
		{name: "zero_gap", dash: []float64{5, 0}, expected: true},
		{name: "zero_dash", dash: []float64{0, 5}, expected: true},
		{name: "single_entry", dash: []float64{5}, expected: true},
		{name: "all_zero", dash: []float64{0, 0}, expected: false},
		{name: "leading_negative", dash: []float64{-5, 5}, expected: false},
		{name: "trailing_negative", dash: []float64{5, -5}, expected: false},
		{name: "nan_entry", dash: []float64{5, math.NaN()}, expected: false},
		{name: "inf_entry", dash: []float64{5, math.Inf(1)}, expected: false},
		{name: "empty", dash: []float64{}, expected: false},
		{name: "nil", dash: nil, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ValidDash(tt.dash))
		})
	}
}

func TestDashable(t *testing.T) {
	t.Parallel()

	assert.True(t, dashable([]float64{5, 5}, 2))
	assert.False(t, dashable([]float64{0, 0}, 2))
	assert.False(t, dashable([]float64{5, 5}, math.Inf(1)))
	assert.False(t, dashable([]float64{5, 5}, math.NaN()))
}

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

	t.Run("single_segment", func(t *testing.T) {
		rec := &recordFlattenerEnd{}
		d := NewDashVertexConverter([]float64{2, 2}, 0, rec)
		d.MoveTo(0, 0)
		d.LineTo(5, 0)
		d.End()

		expect := []string{"M0.0,0.0", "L2.0,0.0", "E", "M4.0,0.0", "L5.0,0.0", "E"}
		assert.Equal(t, expect, rec.moves)
	})

	t.Run("short_segments_accumulate", func(t *testing.T) {
		rec := &recordFlattenerEnd{}
		d := NewDashVertexConverter([]float64{10, 10}, 0, rec)
		d.MoveTo(0, 0)
		for x := 3.0; x <= 60; x += 3 {
			d.LineTo(x, 0)
		}
		d.End()

		// phase advances across short segments: dash ends at x=10, gaps, resumes at x=20
		assert.Contains(t, rec.moves, "L10.0,0.0")
		assert.Contains(t, rec.moves, "M20.0,0.0")
		assert.Contains(t, rec.moves, "E")
	})
}
