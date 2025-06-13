package drawing

import (
	"strconv"
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/math/fixed"
)

func TestPointsToPixels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		dpi, points float64
		want        float64
	}{
		{72, 72, 72},
		{96, 72, 96},
		{96, 36, 48},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := PointsToPixels(tt.dpi, tt.points)
			assert.InDelta(t, tt.want, got, 0.0001)
		})
	}
}

func TestDistanceFuncs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x1, y1, x2, y2 float64
		want           float64
	}{
		{0, 0, 3, 4, 5},
		{1, 2, 1, 2, 0},
		{-1, -1, -4, -5, 5},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.InDelta(t, tt.want, distance(tt.x1, tt.y1, tt.x2, tt.y2), 0.0001)
		})
	}
}

func TestVectorDistance(t *testing.T) {
	t.Parallel()
	tests := []struct{ dx, dy, want float64 }{
		{3, 4, 5},
		{0, 0, 0},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.InDelta(t, tt.want, vectorDistance(tt.dx, tt.dy), 0.0001)
		})
	}
}

func TestFUnitsConversion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   fixed.Int26_6
		want float64
	}{
		{64, 1},
		{96, 1.5},
		{-64, -1},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.InDelta(t, tt.want, fUnitsToFloat64(tt.in), 0.0001)
		})
	}
}

func TestPointToF64Point(t *testing.T) {
	t.Parallel()

	p := truetype.Point{X: 128, Y: -64}
	x, y := pointToF64Point(p)
	assert.InDelta(t, 2.0, x, 0.0001)
	// Y is negated inside function
	assert.InDelta(t, 1.0, y, 0.0001)
}

func TestAbsInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in, want int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tt.want, absInt(tt.in))
		})
	}
}
