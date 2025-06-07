package drawing

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackGraphicContextSaveRestore(t *testing.T) {
	t.Parallel()

	gc := NewStackGraphicContext()
	gc.SetLineWidth(2)
	gc.MoveTo(1, 1)
	gc.Save()
	gc.SetLineWidth(4)
	gc.LineTo(2, 2)
	gc.Restore()
	assert.InDelta(t, 2.0, gc.current.LineWidth, 0.0001)
	x, y := gc.LastPoint()
	assert.InDelta(t, 1.0, x, 0.0001)
	assert.InDelta(t, 1.0, y, 0.0001)
}

func TestStackGraphicContextTransforms(t *testing.T) {
	t.Parallel()

	gc := NewStackGraphicContext()
	gc.Translate(2, 3)
	tr := gc.GetMatrixTransform()
	x, y := tr.TransformPoint(0, 0)
	assert.InDelta(t, 2.0, x, 0.0001)
	assert.InDelta(t, 3.0, y, 0.0001)
	gc.Rotate(math.Pi / 2)
	tr = gc.GetMatrixTransform()
	x, y = tr.TransformPoint(1, 0)
	assert.InDelta(t, 2.0, x, 0.0001)
	assert.InDelta(t, 4.0, y, 0.0001)
}

func TestStackGraphicContextColors(t *testing.T) {
	t.Parallel()

	gc := NewStackGraphicContext()
	gc.SetStrokeColor(ColorRed)
	gc.SetFillColor(ColorBlue)
	assert.Equal(t, ColorRed, gc.current.StrokeColor)
	assert.Equal(t, ColorBlue, gc.current.FillColor)
}

func TestStackMatrixTransform(t *testing.T) {
	t.Parallel()

	gc := NewStackGraphicContext()
	tr := NewTranslationMatrix(5, 7)
	gc.SetMatrixTransform(tr)
	got := gc.GetMatrixTransform()
	x, y := got.TransformPoint(0, 0)
	assert.InDelta(t, 5.0, x, 0.0001)
	assert.InDelta(t, 7.0, y, 0.0001)
}
