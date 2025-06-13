package drawing

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/go-analyze/charts/chartdraw/roboto"
	"github.com/golang/freetype/raster"
	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/math/fixed"
)

func TestRasterGraphicContextBasic(t *testing.T) {
	t.Parallel()

	img := image.NewRGBA(image.Rect(0, 0, 3, 3))
	rgc := NewRasterGraphicContext(img)
	assert.InDelta(t, defaultDPI, rgc.GetDPI(), 0.0)
	rgc.SetDPI(72)
	assert.InDelta(t, 72.0, rgc.GetDPI(), 0.0)
}

func TestRasterFillRectAndClear(t *testing.T) {
	t.Parallel()

	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	rgc := NewRasterGraphicContext(img)
	rgc.SetFillColor(color.RGBA{255, 0, 0, 255})
	rgc.FillRect(0, 0, 2, 2)
	_, _, _, a := img.At(1, 1).RGBA()
	assert.Equal(t, uint32(0xffff), a)

	rgc.Clear()
	_, _, _, a = img.At(1, 1).RGBA()
	assert.Equal(t, uint32(0xffff), a)
}

func TestRasterFillRectanglePath(t *testing.T) {
	t.Parallel()

	img := image.NewRGBA(image.Rect(0, 0, 3, 3))
	rgc := NewRasterGraphicContext(img)
	rgc.SetFillColor(color.RGBA{0, 255, 0, 255})
	p := &Path{}
	p.MoveTo(0, 0)
	p.LineTo(2, 0)
	p.LineTo(2, 2)
	p.LineTo(0, 2)
	p.LineTo(0, 0)
	rgc.Fill(p)
	_, _, _, a := img.At(1, 1).RGBA()
	assert.Equal(t, uint32(0xffff), a)
}

func TestRasterDrawImage(t *testing.T) {
	t.Parallel()

	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	src.Set(0, 0, color.White)
	dst := image.NewRGBA(image.Rect(0, 0, 3, 3))
	rgc := NewRasterGraphicContext(dst)
	rgc.DrawImage(src)
	_, _, _, a := dst.At(0, 0).RGBA()
	assert.Equal(t, uint32(0xffff), a)
}

func TestRasterStrokeAndFillStroke(t *testing.T) {
	t.Parallel()

	img := image.NewRGBA(image.Rect(0, 0, 3, 3))
	rgc := NewRasterGraphicContext(img)
	rgc.SetLineWidth(1)
	rgc.SetStrokeColor(color.Black)
	rgc.SetFillColor(color.RGBA{0, 0, 255, 255})

	p := &Path{}
	p.MoveTo(0, 0)
	p.LineTo(2, 0)
	p.LineTo(2, 2)
	p.LineTo(0, 2)
	p.Close()

	rgc.FillStroke(p)
	_, _, _, a := img.At(1, 1).RGBA()
	assert.Equal(t, uint32(0xffff), a) // fill
	_, _, _, a = img.At(0, 0).RGBA()
	assert.Equal(t, uint32(0xffff), a) // stroke

	img2 := image.NewRGBA(image.Rect(0, 0, 3, 3))
	rgc = NewRasterGraphicContext(img2)
	rgc.SetLineWidth(1)
	rgc.SetStrokeColor(color.Black)
	p = &Path{}
	p.MoveTo(0, 0)
	p.LineTo(2, 0)
	rgc.Stroke(p)
	_, _, _, a = img2.At(1, 0).RGBA()
	assert.Equal(t, uint32(0xffff), a)
}

func TestRasterFontFunctions(t *testing.T) {
	t.Parallel()

	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	rgc := NewRasterGraphicContext(img)

	f, err := truetype.Parse(roboto.Roboto)
	require.NoError(t, err)
	rgc.SetFont(f)
	assert.Equal(t, f, rgc.GetFont())

	rgc.SetFontSize(12)
	assert.InDelta(t, 12.0, rgc.GetFontSize(), 0.0)

	rgc.SetFontSize(8)
	wSmall, err := rgc.CreateStringPath("A", 0, 0)
	require.NoError(t, err)

	rgc.SetFontSize(16)
	wLarge, err := rgc.CreateStringPath("A", 0, 0)
	require.NoError(t, err)
	assert.Greater(t, wLarge, wSmall)
}

func TestRasterCreateStringPathAndBounds(t *testing.T) {
	t.Parallel()

	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	rgc := NewRasterGraphicContext(img)
	f, err := truetype.Parse(roboto.Roboto)
	require.NoError(t, err)
	rgc.SetFont(f)
	rgc.SetFontSize(10)

	idx := f.Index('A')
	expected := fUnitsToFloat64(f.HMetric(fixed.Int26_6(rgc.current.Scale), idx).AdvanceWidth)
	cursor, err := rgc.CreateStringPath("A", 0, 0)
	require.NoError(t, err)
	assert.InDelta(t, expected, cursor, 0.001)
	assert.False(t, rgc.current.Path.IsEmpty())

	left, top, right, bottom, err := rgc.GetStringBounds("A")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, cursor, right-left)

	pbLeft, pbTop, pbRight, pbBottom := pathBounds(rgc.current.Path)
	assert.InDelta(t, left, pbLeft, 0.001)
	assert.InDelta(t, top, pbTop, 0.001)
	assert.InDelta(t, right, pbRight, 0.001)
	assert.InDelta(t, bottom, pbBottom, 0.001)
}

func TestRasterFillAndStrokeStringAt(t *testing.T) {
	t.Parallel()

	f, err := truetype.Parse(roboto.Roboto)
	require.NoError(t, err)

	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	rgc := NewRasterGraphicContext(img)
	rgc.SetFont(f)
	rgc.SetFontSize(10)
	rgc.SetFillColor(color.White)

	left, top, right, bottom, err := rgc.GetStringBounds("A")
	require.NoError(t, err)
	x, y := 10.0, 30.0
	_, err = rgc.FillStringAt("A", x, y)
	require.NoError(t, err)

	x1 := int(math.Floor(left + x))
	y1 := int(math.Floor(top + y))
	found := false
	for yy := y1; yy < int(math.Ceil(bottom+y)) && !found; yy++ {
		for xx := x1; xx < int(math.Ceil(right+x)) && !found; xx++ {
			_, _, _, a := img.At(xx, yy).RGBA()
			if a != 0 {
				found = true
			}
		}
	}
	assert.True(t, found, "filled text not drawn")

	img2 := image.NewRGBA(image.Rect(0, 0, 50, 50))
	rgc2 := NewRasterGraphicContext(img2)
	rgc2.SetFont(f)
	rgc2.SetFontSize(10)
	rgc2.SetStrokeColor(color.White)

	_, err = rgc2.StrokeStringAt("A", x, y)
	require.NoError(t, err)
	found = false
	for yy := y1; yy < int(math.Ceil(bottom+y)) && !found; yy++ {
		for xx := x1; xx < int(math.Ceil(right+x)) && !found; xx++ {
			_, _, _, a := img2.At(xx, yy).RGBA()
			if a != 0 {
				found = true
			}
		}
	}
	assert.True(t, found, "stroked text not drawn")
}

func pathBounds(p *Path) (left, top, right, bottom float64) {
	if len(p.Points) == 0 {
		return
	}
	left, top = p.Points[0], p.Points[1]
	right, bottom = left, top
	for i := 2; i < len(p.Points); i += 2 {
		x, y := p.Points[i], p.Points[i+1]
		if x < left {
			left = x
		}
		if y < top {
			top = y
		}
		if x > right {
			right = x
		}
		if y > bottom {
			bottom = y
		}
	}
	return
}

func TestNewRasterGraphicContextWithPainter(t *testing.T) {
	t.Parallel()

	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	p := raster.NewRGBAPainter(img)
	rgc := NewRasterGraphicContextWithPainter(img, p)
	if rgc.painter != p {
		t.Fatalf("painter not set")
	}
}
