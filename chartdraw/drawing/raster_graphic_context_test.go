package drawing

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
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
