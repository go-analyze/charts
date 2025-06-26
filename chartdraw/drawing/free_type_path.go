package drawing

import (
	"github.com/golang/freetype/raster"
	"golang.org/x/image/math/fixed"
)

// FtLineBuilder is a builder for freetype raster glyphs.
type FtLineBuilder struct {
	Adder raster.Adder
}

// MoveTo starts a new segment at the given point (for PathBuilder interface).
func (liner FtLineBuilder) MoveTo(x, y float64) {
	liner.Adder.Start(fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)})
}

// LineTo adds a line segment from the current point to the specified point (for PathBuilder interface).
func (liner FtLineBuilder) LineTo(x, y float64) {
	liner.Adder.Add1(fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)})
}

// End finalizes the current path (for PathBuilder interface).
func (liner FtLineBuilder) End() {}
