package drawing

import (
	"math"
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/roboto"
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

func TestStackComposeMatrixTransform(t *testing.T) {
	t.Parallel()

	gc := NewStackGraphicContext()
	gc.SetMatrixTransform(NewTranslationMatrix(5, 7))
	gc.ComposeMatrixTransform(NewTranslationMatrix(3, 4))
	got := gc.GetMatrixTransform()
	x, y := got.TransformPoint(0, 0)
	assert.InDelta(t, 8.0, x, 0.0001)
	assert.InDelta(t, 11.0, y, 0.0001)
}

func TestStackLineDash(t *testing.T) {
	t.Parallel()

	gc := NewStackGraphicContext()
	dash := []float64{1, 2, 3}
	gc.SetLineDash(dash, 0.5)
	assert.Equal(t, dash, gc.current.Dash)
	assert.InDelta(t, 0.5, gc.current.DashOffset, 0.0001)
}

func TestStackFontRoundTrip(t *testing.T) {
	t.Parallel()

	f, err := truetype.Parse(roboto.Roboto)
	require.NoError(t, err)

	gc := NewStackGraphicContext()
	gc.SetFont(f)
	gc.SetFontSize(13.0)
	assert.Equal(t, f, gc.GetFont())
	assert.InDelta(t, 13.0, gc.GetFontSize(), 0.0001)
}

func TestStackScale(t *testing.T) {
	t.Parallel()

	gc := NewStackGraphicContext()
	before := gc.GetMatrixTransform()
	gc.Scale(2, 3)
	after := gc.GetMatrixTransform()

	assert.False(t, before.Equals(after))
	sx, sy := after.GetScaling()
	assert.InDelta(t, 2.0, sx, 0.0001)
	assert.InDelta(t, 3.0, sy, 0.0001)
}

func TestStackSetFillRule(t *testing.T) {
	t.Parallel()

	gc := NewStackGraphicContext()
	gc.SetFillRule(FillRuleWinding)
	assert.Equal(t, FillRuleWinding, gc.current.FillRule)
	gc.SetFillRule(FillRuleEvenOdd)
	assert.Equal(t, FillRuleEvenOdd, gc.current.FillRule)
}

func TestStackBeginPathIsEmpty(t *testing.T) {
	t.Parallel()

	gc := NewStackGraphicContext()
	gc.MoveTo(1, 1)
	gc.LineTo(2, 2)
	require.False(t, gc.IsEmpty())
	gc.BeginPath()
	assert.True(t, gc.IsEmpty())
	assert.Empty(t, gc.current.Path.Points)
}

func TestStackPathCurvesAndClose(t *testing.T) {
	t.Parallel()

	gc := NewStackGraphicContext()
	gc.QuadCurveTo(1, 1, 2, 2)
	gc.CubicCurveTo(3, 3, 4, 4, 5, 5)
	gc.ArcTo(0, 0, 1, 1, 0, math.Pi/2)
	gc.Close()

	expComp := []PathComponent{
		MoveToComponent,
		QuadCurveToComponent,
		CubicCurveToComponent,
		LineToComponent,
		ArcToComponent,
		CloseComponent,
	}
	assert.Equal(t, expComp, gc.current.Path.Components)

	expPts := []float64{
		0, 0,
		1, 1, 2, 2,
		3, 3, 4, 4, 5, 5,
		1, 0,
		0, 0, 1, 1, 0, math.Pi / 2,
	}
	assert.InDeltaSlice(t, expPts, gc.current.Path.Points, 0.0001)
}
