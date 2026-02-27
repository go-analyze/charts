package chartdraw

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoxClone(t *testing.T) {
	t.Parallel()

	a := Box{Top: 5, Left: 5, Right: 15, Bottom: 15}
	b := a.Clone()
	assert.True(t, a.Equals(b))
	assert.True(t, b.Equals(a))
}

func TestBoxWith(t *testing.T) {
	t.Parallel()

	base := Box{
		Top:    100,
		Left:   100,
		Right:  100,
		Bottom: 100,
	}
	t.Run("top", func(t *testing.T) {
		updated := base.WithTop(200)
		assert.True(t, updated.IsSet)
		assert.Equal(t, 200, updated.Top)
		assert.Equal(t, 100, updated.Left)
		assert.Equal(t, 100, updated.Right)
		assert.Equal(t, 100, updated.Bottom)
	})
	t.Run("left", func(t *testing.T) {
		updated := base.WithLeft(200)
		assert.True(t, updated.IsSet)
		assert.Equal(t, 100, updated.Top)
		assert.Equal(t, 200, updated.Left)
		assert.Equal(t, 100, updated.Right)
		assert.Equal(t, 100, updated.Bottom)
	})
	t.Run("right", func(t *testing.T) {
		updated := base.WithRight(200)
		assert.True(t, updated.IsSet)
		assert.Equal(t, 100, updated.Top)
		assert.Equal(t, 100, updated.Left)
		assert.Equal(t, 200, updated.Right)
		assert.Equal(t, 100, updated.Bottom)
	})
	t.Run("bottom", func(t *testing.T) {
		updated := base.WithBottom(200)
		assert.True(t, updated.IsSet)
		assert.Equal(t, 100, updated.Top)
		assert.Equal(t, 100, updated.Left)
		assert.Equal(t, 100, updated.Right)
		assert.Equal(t, 200, updated.Bottom)
	})
}

func TestBoxEquals(t *testing.T) {
	t.Parallel()

	a := Box{Top: 5, Left: 5, Right: 15, Bottom: 15}
	b := Box{Top: 10, Left: 10, Right: 30, Bottom: 30}
	c := Box{Top: 5, Left: 5, Right: 15, Bottom: 15}
	assert.True(t, a.Equals(a))
	assert.True(t, a.Equals(c))
	assert.True(t, c.Equals(a))
	assert.False(t, a.Equals(b))
	assert.False(t, c.Equals(b))
	assert.False(t, b.Equals(a))
	assert.False(t, b.Equals(c))
}

func TestBoxIsBiggerThan(t *testing.T) {
	t.Parallel()

	a := Box{Top: 5, Left: 5, Right: 25, Bottom: 25}
	b := Box{Top: 10, Left: 10, Right: 20, Bottom: 20} // only half bigger
	c := Box{Top: 1, Left: 1, Right: 30, Bottom: 30}   //bigger
	assert.True(t, a.IsBiggerThan(b))
	assert.False(t, a.IsBiggerThan(c))
	assert.True(t, c.IsBiggerThan(a))
}

func TestBoxIsSmallerThan(t *testing.T) {
	t.Parallel()

	a := Box{Top: 5, Left: 5, Right: 25, Bottom: 25}
	b := Box{Top: 10, Left: 10, Right: 20, Bottom: 20} // only half bigger
	c := Box{Top: 1, Left: 1, Right: 30, Bottom: 30}   //bigger
	assert.False(t, a.IsSmallerThan(b))
	assert.True(t, a.IsSmallerThan(c))
	assert.False(t, c.IsSmallerThan(a))
}

func TestOverlaps(t *testing.T) {
	t.Parallel()

	t.Run("identical_box", func(t *testing.T) {
		box1 := Box{Top: 0, Left: 0, Right: 10, Bottom: 10}
		box2 := Box{Top: 0, Left: 0, Right: 10, Bottom: 10}
		assert.True(t, box1.Overlaps(box2))
		assert.True(t, box2.Overlaps(box1))
	})
	t.Run("partial_overlap", func(t *testing.T) {
		box1 := Box{Top: 0, Left: 0, Right: 10, Bottom: 10}
		box2 := Box{Top: 5, Left: 5, Right: 15, Bottom: 15}
		assert.True(t, box1.Overlaps(box2))
		assert.True(t, box2.Overlaps(box1))
	})
	t.Run("corner_touch_not_overlap", func(t *testing.T) {
		box1 := Box{Top: 0, Left: 0, Right: 10, Bottom: 10}
		// This box starts exactly where box1 ends (corner at 10,10)
		box2 := Box{Top: 10, Left: 10, Right: 20, Bottom: 20}
		assert.False(t, box1.Overlaps(box2))
		assert.False(t, box2.Overlaps(box1))
	})
	t.Run("completely_inside", func(t *testing.T) {
		outer := Box{Top: 0, Left: 0, Right: 10, Bottom: 10}
		inner := Box{Top: 2, Left: 2, Right: 8, Bottom: 8}
		assert.True(t, outer.Overlaps(inner))
		assert.True(t, inner.Overlaps(outer))
	})
	t.Run("no_overlap_zero", func(t *testing.T) {
		box := Box{Top: 0, Left: 0, Right: 10, Bottom: 10}
		assert.False(t, box.Overlaps(BoxZero))
		assert.False(t, BoxZero.Overlaps(box))
	})
	t.Run("no_overlap_zero_center", func(t *testing.T) {
		box := Box{Top: 10, Left: 10, Right: 20, Bottom: 20}
		assert.False(t, box.Overlaps(BoxZero))
		assert.False(t, BoxZero.Overlaps(box))
	})
	t.Run("no_overlap_horizontally", func(t *testing.T) {
		box1 := Box{Top: 0, Left: 0, Right: 10, Bottom: 10}
		box2 := Box{Top: 0, Left: 11, Right: 20, Bottom: 10}
		assert.False(t, box1.Overlaps(box2))
		assert.False(t, box2.Overlaps(box1))
	})
	t.Run("no_overlap_vertically", func(t *testing.T) {
		box1 := Box{Top: 0, Left: 0, Right: 10, Bottom: 10}
		box2 := Box{Top: 11, Left: 0, Right: 10, Bottom: 20}
		assert.False(t, box1.Overlaps(box2))
		assert.False(t, box2.Overlaps(box1))
	})
}

func TestBoxGrow(t *testing.T) {
	t.Parallel()

	a := Box{Top: 1, Left: 2, Right: 15, Bottom: 15}
	b := Box{Top: 4, Left: 5, Right: 30, Bottom: 35}
	c := a.Grow(b)
	assert.False(t, c.Equals(b))
	assert.False(t, c.Equals(a))
	assert.Equal(t, 1, c.Top)
	assert.Equal(t, 2, c.Left)
	assert.Equal(t, 30, c.Right)
	assert.Equal(t, 35, c.Bottom)
}

func TestBoxFit(t *testing.T) {
	t.Parallel()

	a := Box{Top: 64, Left: 64, Right: 192, Bottom: 192}
	b := Box{Top: 16, Left: 16, Right: 256, Bottom: 170}
	c := Box{Top: 16, Left: 16, Right: 170, Bottom: 256}

	fab := a.Fit(b)
	assert.Equal(t, a.Left, fab.Left)
	assert.Equal(t, a.Right, fab.Right)
	assert.Less(t, fab.Top, fab.Bottom)
	assert.Less(t, fab.Left, fab.Right)
	assert.Less(t, math.Abs(b.Aspect()-fab.Aspect()), 0.02)

	fac := a.Fit(c)
	assert.Equal(t, a.Top, fac.Top)
	assert.Equal(t, a.Bottom, fac.Bottom)
	assert.Less(t, math.Abs(c.Aspect()-fac.Aspect()), 0.02)
}

func TestBoxConstrain(t *testing.T) {
	t.Parallel()

	a := Box{Top: 64, Left: 64, Right: 192, Bottom: 192}
	b := Box{Top: 16, Left: 16, Right: 256, Bottom: 170}
	c := Box{Top: 16, Left: 16, Right: 170, Bottom: 256}

	cab := a.Constrain(b)
	assert.Equal(t, 64, cab.Top)
	assert.Equal(t, 64, cab.Left)
	assert.Equal(t, 192, cab.Right)
	assert.Equal(t, 170, cab.Bottom)

	cac := a.Constrain(c)
	assert.Equal(t, 64, cac.Top)
	assert.Equal(t, 64, cac.Left)
	assert.Equal(t, 170, cac.Right)
	assert.Equal(t, 192, cac.Bottom)
}

func TestBoxOuterConstrain(t *testing.T) {
	t.Parallel()

	box := NewBox(0, 0, 100, 100)
	canvas := NewBox(5, 5, 95, 95)
	taller := NewBox(-10, 5, 50, 50)

	c := canvas.OuterConstrain(box, taller)
	assert.Equal(t, 15, c.Top, c.String())
	assert.Equal(t, 5, c.Left, c.String())
	assert.Equal(t, 95, c.Right, c.String())
	assert.Equal(t, 95, c.Bottom, c.String())

	wider := NewBox(5, 5, 110, 50)
	d := canvas.OuterConstrain(box, wider)
	assert.Equal(t, 5, d.Top, d.String())
	assert.Equal(t, 5, d.Left, d.String())
	assert.Equal(t, 85, d.Right, d.String())
	assert.Equal(t, 95, d.Bottom, d.String())
}

func TestBoxValidate(t *testing.T) {
	t.Parallel()

	require.Error(t, (Box{Left: -1}).Validate())
	require.Error(t, (Box{Right: -1}).Validate())
	require.Error(t, (Box{Top: -1}).Validate())
	require.Error(t, (Box{Bottom: -1}).Validate())
	require.NoError(t, (Box{Top: 1, Left: 1, Right: 1, Bottom: 1}).Validate())
}

func TestBoxShift(t *testing.T) {
	t.Parallel()

	b := Box{
		Top:    5,
		Left:   5,
		Right:  10,
		Bottom: 10,
	}

	shifted := b.Shift(1, 2)
	assert.Equal(t, 7, shifted.Top)
	assert.Equal(t, 6, shifted.Left)
	assert.Equal(t, 11, shifted.Right)
	assert.Equal(t, 12, shifted.Bottom)
}

func TestBoxCenter(t *testing.T) {
	t.Parallel()

	b := Box{
		Top:    10,
		Left:   10,
		Right:  20,
		Bottom: 30,
	}
	cx, cy := b.Center()
	assert.Equal(t, 15, cx)
	assert.Equal(t, 20, cy)
}

func TestBoxCornersCenter(t *testing.T) {
	t.Parallel()

	bc := BoxCorners{
		TopLeft:     Point{5, 5},
		TopRight:    Point{15, 5},
		BottomRight: Point{15, 15},
		BottomLeft:  Point{5, 15},
	}

	cx, cy := bc.Center()
	assert.Equal(t, 10, cx)
	assert.Equal(t, 10, cy)
}

func TestBoxCornersRotate(t *testing.T) {
	t.Parallel()

	bc := BoxCorners{
		TopLeft:     Point{5, 5},
		TopRight:    Point{15, 5},
		BottomRight: Point{15, 15},
		BottomLeft:  Point{5, 15},
	}

	rotated := bc.Rotate(45)
	assert.True(t, rotated.TopLeft.Equals(Point{10, 3}), rotated.String())
}

func TestDistanceTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		p1       Point
		p2       Point
		expected float64
	}{
		{"same_point", Point{0, 0}, Point{0, 0}, 0},
		{"positive_int", Point{0, 0}, Point{3, 4}, 5},
		{"negative_and_positive", Point{-1, -1}, Point{1, 1}, 2.8284271247461903},
		{"negative_int", Point{-3, -4}, Point{0, 0}, 5},
		{"one_axis_zero", Point{0, 5}, Point{0, -5}, 10},
		{"both_axes_same", Point{2, 2}, Point{-2, -2}, 5.656854249492381},
		{"large_numbers", Point{100000, 100000}, Point{-100000, -100000}, 282842.71247461904},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.InDelta(t, tc.expected, tc.p1.DistanceTo(tc.p2), 0)
		})
	}
}
