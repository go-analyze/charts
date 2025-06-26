package chartdraw

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

// PNG returns a new png raster renderer.
func PNG(width, height int) Renderer {
	i := image.NewRGBA(image.Rect(0, 0, width, height))
	return &rasterRenderer{
		i:          i,
		gc:         drawing.NewRasterGraphicContext(i),
		encodeFunc: png.Encode,
	}
}

// JPG returns a new jpg raster renderer.
func JPG(width, height int) Renderer {
	i := image.NewRGBA(image.Rect(0, 0, width, height))
	return &rasterRenderer{
		i:  i,
		gc: drawing.NewRasterGraphicContext(i),
		encodeFunc: func(w io.Writer, i image.Image) error {
			return jpeg.Encode(w, i, &jpeg.Options{Quality: 90})
		},
	}
}

// rasterRenderer renders chart commands to a bitmap.
type rasterRenderer struct {
	i          *image.RGBA
	gc         *drawing.RasterGraphicContext
	encodeFunc func(w io.Writer, i image.Image) error

	rotateRadians *float64

	s Style
}

func (rr *rasterRenderer) ResetStyle() {
	rr.s = Style{
		FontStyle: FontStyle{
			Font: rr.s.Font,
		},
	}
	rr.ClearTextRotation()
}

// GetDPI returns the dpi.
func (rr *rasterRenderer) GetDPI() float64 {
	return rr.gc.GetDPI()
}

// SetDPI sets the rendering DPI (for Renderer interface).
func (rr *rasterRenderer) SetDPI(dpi float64) {
	rr.gc.SetDPI(dpi)
}

// SetClassName is ignored because raster images have no class names (for Renderer interface).
func (rr *rasterRenderer) SetClassName(_ string) {}

// SetStrokeColor sets the stroke color for future paths (for Renderer interface).
func (rr *rasterRenderer) SetStrokeColor(c drawing.Color) {
	rr.s.StrokeColor = c
}

// SetStrokeWidth sets the width of drawn lines (for Renderer interface).
func (rr *rasterRenderer) SetStrokeWidth(width float64) {
	rr.s.StrokeWidth = width
}

// SetStrokeDashArray sets the stroke dash array.
func (rr *rasterRenderer) SetStrokeDashArray(dashArray []float64) {
	rr.s.StrokeDashArray = dashArray
}

// SetFillColor sets the fill color for future paths (for Renderer interface).
func (rr *rasterRenderer) SetFillColor(c drawing.Color) {
	rr.s.FillColor = c
}

// MoveTo moves the drawing cursor to the given position (for PathBuilder interface).
func (rr *rasterRenderer) MoveTo(x, y int) {
	rr.gc.MoveTo(float64(x), float64(y))
}

// LineTo adds a line to the current path (for PathBuilder interface).
func (rr *rasterRenderer) LineTo(x, y int) {
	rr.gc.LineTo(float64(x), float64(y))
}

// QuadCurveTo adds a quadratic curve to the current path (for PathBuilder interface).
func (rr *rasterRenderer) QuadCurveTo(cx, cy, x, y int) {
	rr.gc.QuadCurveTo(float64(cx), float64(cy), float64(x), float64(y))
}

// ArcTo appends an elliptical arc to the current path (for PathBuilder interface).
func (rr *rasterRenderer) ArcTo(cx, cy int, rx, ry, startAngle, delta float64) {
	rr.gc.ArcTo(float64(cx), float64(cy), rx, ry, startAngle, delta)
}

// Close closes the current path (for PathBuilder interface).
func (rr *rasterRenderer) Close() {
	rr.gc.Close()
}

// Stroke renders the path outline without filling it (for PathBuilder interface).
func (rr *rasterRenderer) Stroke() {
	rr.gc.SetStrokeColor(rr.s.StrokeColor)
	rr.gc.SetLineWidth(rr.s.StrokeWidth)
	rr.gc.SetLineDash(rr.s.StrokeDashArray, 0)
	rr.gc.Stroke()
}

// Fill renders the path fill without stroking it (for PathBuilder interface).
func (rr *rasterRenderer) Fill() {
	rr.gc.SetFillColor(rr.s.FillColor)
	rr.gc.Fill()
}

// FillStroke fills and then strokes the current path (for PathBuilder interface).
func (rr *rasterRenderer) FillStroke() {
	rr.gc.SetFillColor(rr.s.FillColor)
	rr.gc.SetStrokeColor(rr.s.StrokeColor)
	rr.gc.SetLineWidth(rr.s.StrokeWidth)
	rr.gc.SetLineDash(rr.s.StrokeDashArray, 0)
	rr.gc.FillStroke()
}

// Circle fully draws a circle at a given point but does not apply the fill or stroke (for PathBuilder interface).
func (rr *rasterRenderer) Circle(radius float64, x, y int) {
	xf, yf := float64(x), float64(y)
	rr.gc.MoveTo(xf-radius, yf) // explicit MoveTo to avoid LineTo if components already on raster, see issue #78
	rr.gc.ArcTo(xf, yf, radius, radius, 0, _2pi)
}

// SetFont sets the font used for text drawing (for Renderer interface).
func (rr *rasterRenderer) SetFont(f *truetype.Font) {
	rr.s.Font = f
}

// SetFontSize sets the font size in points (for Renderer interface).
func (rr *rasterRenderer) SetFontSize(size float64) {
	rr.s.FontSize = size
}

// SetFontColor sets the color used for text drawing (for Renderer interface).
func (rr *rasterRenderer) SetFontColor(c drawing.Color) {
	rr.s.FontColor = c
}

// Text draws the provided string at the given coordinates (for Renderer interface).
func (rr *rasterRenderer) Text(body string, x, y int) {
	if body == "" {
		return
	}
	xf, yf := rr.getCoords(x, y)
	rr.gc.SetFont(rr.s.Font)
	rr.gc.SetFontSize(rr.s.FontSize)
	rr.gc.SetFillColor(rr.s.FontColor)
	// TODO - handle error?
	_, _ = rr.gc.CreateStringPath(body, float64(xf), float64(yf))
	rr.gc.Fill()
}

// MeasureText returns the height and width in pixels of a string.
func (rr *rasterRenderer) MeasureText(body string) Box {
	rr.gc.SetFont(rr.s.Font)
	rr.gc.SetFontSize(rr.s.FontSize)
	rr.gc.SetFillColor(rr.s.FontColor)
	l, t, r, b, err := rr.gc.GetStringBounds(body)
	if err != nil {
		return Box{}
	}
	if l < 0 {
		r -= l // equivalent to r+(-1*l)
		l = 0
	} else if l > 0 {
		r += l
		l = 0
	}
	if t < 0 {
		b -= t
		t = 0
	} else if t > 0 {
		b += t
		t = 0
	}

	textBox := Box{
		Top:    int(math.Ceil(t)),
		Left:   int(math.Ceil(l)),
		Right:  int(math.Ceil(r)),
		Bottom: int(math.Ceil(b)),
		IsSet:  true,
	}
	if rr.rotateRadians == nil {
		return textBox
	}

	return textBox.Corners().Rotate(RadiansToDegrees(*rr.rotateRadians)).Box()
}

// SetTextRotation sets a text rotation.
func (rr *rasterRenderer) SetTextRotation(radians float64) {
	rr.rotateRadians = &radians
}

func (rr *rasterRenderer) getCoords(x, y int) (xf, yf int) {
	if rr.rotateRadians == nil {
		xf = x
		yf = y
		return
	}

	rr.gc.Translate(float64(x), float64(y))
	rr.gc.Rotate(*rr.rotateRadians)
	return
}

// ClearTextRotation clears text rotation.
func (rr *rasterRenderer) ClearTextRotation() {
	rr.gc.SetMatrixTransform(drawing.NewIdentityMatrix())
	rr.rotateRadians = nil
}

// Save writes the rendered image to the provided writer (for Renderer interface).
func (rr *rasterRenderer) Save(w io.Writer) error {
	if typed, isTyped := w.(RGBACollector); isTyped {
		typed.SetRGBA(rr.i)
		return nil
	} else if rr.encodeFunc != nil {
		return rr.encodeFunc(w, rr.i)
	}
	return png.Encode(w, rr.i)
}
