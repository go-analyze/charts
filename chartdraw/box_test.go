package chartdraw

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoxClone(t *testing.T) {
	a := Box{Top: 5, Left: 5, Right: 15, Bottom: 15}
	b := a.Clone()
	assert.True(t, a.Equals(b))
	assert.True(t, b.Equals(a))
}

func TestBoxEquals(t *testing.T) {
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
	a := Box{Top: 5, Left: 5, Right: 25, Bottom: 25}
	b := Box{Top: 10, Left: 10, Right: 20, Bottom: 20} // only half bigger
	c := Box{Top: 1, Left: 1, Right: 30, Bottom: 30}   //bigger
	assert.True(t, a.IsBiggerThan(b))
	assert.False(t, a.IsBiggerThan(c))
	assert.True(t, c.IsBiggerThan(a))
}

func TestBoxIsSmallerThan(t *testing.T) {
	a := Box{Top: 5, Left: 5, Right: 25, Bottom: 25}
	b := Box{Top: 10, Left: 10, Right: 20, Bottom: 20} // only half bigger
	c := Box{Top: 1, Left: 1, Right: 30, Bottom: 30}   //bigger
	assert.False(t, a.IsSmallerThan(b))
	assert.True(t, a.IsSmallerThan(c))
	assert.False(t, c.IsSmallerThan(a))
}

func TestBoxGrow(t *testing.T) {
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
	a := Box{Top: 64, Left: 64, Right: 192, Bottom: 192}
	b := Box{Top: 16, Left: 16, Right: 256, Bottom: 170}
	c := Box{Top: 16, Left: 16, Right: 170, Bottom: 256}

	fab := a.Fit(b)
	assert.Equal(t, a.Left, fab.Left)
	assert.Equal(t, a.Right, fab.Right)
	assert.True(t, fab.Top < fab.Bottom)
	assert.True(t, fab.Left < fab.Right)
	assert.True(t, math.Abs(b.Aspect()-fab.Aspect()) < 0.02)

	fac := a.Fit(c)
	assert.Equal(t, a.Top, fac.Top)
	assert.Equal(t, a.Bottom, fac.Bottom)
	assert.True(t, math.Abs(c.Aspect()-fac.Aspect()) < 0.02)
}

func TestBoxConstrain(t *testing.T) {
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

func TestBoxShift(t *testing.T) {
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
	bc := BoxCorners{
		TopLeft:     Point{5, 5},
		TopRight:    Point{15, 5},
		BottomRight: Point{15, 15},
		BottomLeft:  Point{5, 15},
	}

	rotated := bc.Rotate(45)
	assert.True(t, rotated.TopLeft.Equals(Point{10, 3}), rotated.String())
}
