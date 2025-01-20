package chartdraw

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"golang.org/x/image/font"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

// SVG returns a new png/raster renderer.
func SVG(width, height int) Renderer {
	buffer := bytes.NewBuffer([]byte{})
	canvas := newCanvas(buffer)
	canvas.Start(width, height)
	return &vectorRenderer{
		b:   buffer,
		c:   canvas,
		s:   &Style{},
		p:   []string{},
		dpi: DefaultDPI,
	}
}

// SVGWithCSS returns a new png/raster renderer with attached custom CSS
// The optional nonce argument sets a CSP nonce.
func SVGWithCSS(css string, nonce string) func(width, height int) Renderer {
	return func(width, height int) Renderer {
		buffer := bytes.NewBuffer([]byte{})
		canvas := newCanvas(buffer)
		canvas.css = css
		canvas.nonce = nonce
		canvas.Start(width, height)
		return &vectorRenderer{
			b:   buffer,
			c:   canvas,
			s:   &Style{},
			p:   []string{},
			dpi: DefaultDPI,
		}
	}
}

// vectorRenderer renders chart commands to a bitmap.
type vectorRenderer struct {
	dpi float64
	b   *bytes.Buffer
	c   *canvas
	s   *Style
	p   []string
	fc  *font.Drawer
}

func (vr *vectorRenderer) ResetStyle() {
	vr.s = &Style{
		FontStyle: FontStyle{
			Font: vr.s.Font,
		},
	}
	vr.fc = nil
}

// GetDPI returns the dpi.
func (vr *vectorRenderer) GetDPI() float64 {
	return vr.dpi
}

// SetDPI implements the interface method.
func (vr *vectorRenderer) SetDPI(dpi float64) {
	vr.dpi = dpi
	vr.c.dpi = dpi
}

// SetClassName implements the interface method.
func (vr *vectorRenderer) SetClassName(classname string) {
	vr.s.ClassName = classname
}

// SetStrokeColor implements the interface method.
func (vr *vectorRenderer) SetStrokeColor(c drawing.Color) {
	vr.s.StrokeColor = c
}

// SetFillColor implements the interface method.
func (vr *vectorRenderer) SetFillColor(c drawing.Color) {
	vr.s.FillColor = c
}

// SetStrokeWidth implements the interface method.
func (vr *vectorRenderer) SetStrokeWidth(width float64) {
	vr.s.StrokeWidth = width
}

// SetStrokeDashArray sets the stroke dash array.
func (vr *vectorRenderer) SetStrokeDashArray(dashArray []float64) {
	vr.s.StrokeDashArray = dashArray
}

// MoveTo implements the interface method.
func (vr *vectorRenderer) MoveTo(x, y int) {
	vr.p = append(vr.p, "M "+strconv.Itoa(x)+" "+strconv.Itoa(y))
}

// LineTo implements the interface method.
func (vr *vectorRenderer) LineTo(x, y int) {
	vr.p = append(vr.p, "L "+strconv.Itoa(x)+" "+strconv.Itoa(y))
}

// QuadCurveTo draws a quad curve.
func (vr *vectorRenderer) QuadCurveTo(cx, cy, x, y int) {
	vr.p = append(vr.p, "Q"+strconv.Itoa(cx)+","+strconv.Itoa(cy)+" "+strconv.Itoa(x)+","+strconv.Itoa(y))
}

func (vr *vectorRenderer) ArcTo(cx, cy int, rx, ry, startAngle, delta float64) {
	startAngle = RadianAdd(startAngle, _pi2)
	endAngle := RadianAdd(startAngle, delta)

	startx := cx + int(rx*math.Sin(startAngle))
	starty := cy - int(ry*math.Cos(startAngle))

	if len(vr.p) > 0 {
		vr.LineTo(startx, starty)
	} else {
		vr.MoveTo(startx, starty)
	}

	endx := cx + int(rx*math.Sin(endAngle))
	endy := cy - int(ry*math.Cos(endAngle))

	dd := RadiansToDegrees(delta)

	largeArcFlag := 0
	if delta > _pi {
		largeArcFlag = 1
	}

	vr.p = append(vr.p, fmt.Sprintf("A %d %d %0.2f %d 1 %d %d", int(rx), int(ry), dd, largeArcFlag, endx, endy))
}

// Close closes a shape.
func (vr *vectorRenderer) Close() {
	vr.p = append(vr.p, "Z")
}

// Stroke draws the path with no fill.
func (vr *vectorRenderer) Stroke() {
	vr.drawPath()
}

// Fill draws the path with no stroke.
func (vr *vectorRenderer) Fill() {
	vr.drawPath()
}

// FillStroke draws the path with both fill and stroke.
func (vr *vectorRenderer) FillStroke() {
	vr.drawPath()
}

// drawPath draws the path set into the p slice.
func (vr *vectorRenderer) drawPath() {
	vr.c.Path(vr.p, vr.s.GetFillAndStrokeOptions())
	vr.p = vr.p[:0] // clear the path
}

// Circle implements the interface method.
func (vr *vectorRenderer) Circle(radius float64, x, y int) {
	vr.c.Circle(x, y, int(radius), vr.s.GetFillAndStrokeOptions())
}

// SetFont implements the interface method.
func (vr *vectorRenderer) SetFont(f *truetype.Font) {
	vr.s.Font = f
}

// SetFontColor implements the interface method.
func (vr *vectorRenderer) SetFontColor(c drawing.Color) {
	vr.s.FontColor = c
}

// SetFontSize implements the interface method.
func (vr *vectorRenderer) SetFontSize(size float64) {
	vr.s.FontSize = size
}

// Text draws a text blob.
func (vr *vectorRenderer) Text(body string, x, y int) {
	vr.c.Text(x, y, body, vr.s.GetTextOptions())
}

// MeasureText uses the truetype font drawer to measure the width of text.
func (vr *vectorRenderer) MeasureText(body string) (box Box) {
	if vr.s.GetFont() != nil {
		vr.fc = &font.Drawer{
			Face: truetype.NewFace(vr.s.GetFont(), &truetype.Options{
				DPI:  vr.dpi,
				Size: vr.s.FontSize,
			}),
		}
		w := vr.fc.MeasureString(body).Ceil()

		box.Right = w
		box.Bottom = int(drawing.PointsToPixels(vr.dpi, vr.s.FontSize))
		box.IsSet = true
		if vr.c.textTheta == nil {
			return
		}
		box = box.Corners().Rotate(RadiansToDegrees(*vr.c.textTheta)).Box()
	}
	return
}

// SetTextRotation sets the text rotation.
func (vr *vectorRenderer) SetTextRotation(radians float64) {
	vr.c.textTheta = &radians
}

// ClearTextRotation clears the text rotation.
func (vr *vectorRenderer) ClearTextRotation() {
	vr.c.textTheta = nil
}

// Save saves the renderer's contents to a writer.
func (vr *vectorRenderer) Save(w io.Writer) error {
	vr.c.End()
	_, err := w.Write(vr.b.Bytes())
	return err
}

func newCanvas(w io.Writer) *canvas {
	return &canvas{
		w:   w,
		dpi: DefaultDPI,
	}
}

type canvas struct {
	w         io.Writer
	dpi       float64
	textTheta *float64
	width     int
	height    int
	css       string
	nonce     string
}

func (c *canvas) Start(width, height int) {
	c.width = width
	c.height = height
	// TODO - error handling on write
	_, _ = c.w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 ` + strconv.Itoa(c.width) + ` ` + strconv.Itoa(c.height) + `">`))
	if c.css != "" {
		_, _ = c.w.Write([]byte(`<style type="text/css"`))
		if c.nonce != "" {
			// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy
			_, _ = c.w.Write([]byte(` nonce="` + c.nonce + `"`))
		}
		// To avoid compatibility issues between XML and CSS (f.e. with child selectors) we should encapsulate the CSS with CDATA.
		_, _ = c.w.Write([]byte(`><![CDATA[` + c.css + `]]></style>`))
	}
}

func (c *canvas) Path(parts []string, style Style) {
	if len(parts) == 0 {
		return
	}
	bb := bytes.NewBuffer(make([]byte, 0, 80))
	bb.WriteString(`<path `)
	c.writeStrokeDashArray(bb, style)
	bb.WriteString(` d="`)
	for i, p := range parts {
		if i > 0 {
			bb.WriteRune('\n')
		}
		bb.WriteString(p)
	}
	bb.WriteString(`" `)
	c.styleAsSVG(bb, style, false)
	bb.WriteString(`/>`)

	_, _ = c.w.Write(bb.Bytes())
}

func (c *canvas) Text(x, y int, body string, style Style) {
	if body == "" {
		return
	}
	bb := bytes.NewBuffer(make([]byte, 0, 128))
	bb.WriteString(`<text x="`)
	bb.WriteString(strconv.Itoa(x))
	bb.WriteString(`" y="`)
	bb.WriteString(strconv.Itoa(y))
	bb.WriteString(`" `)
	c.styleAsSVG(bb, style, true)
	if c.textTheta != nil {
		bb.WriteString(fmt.Sprintf(` transform="rotate(%0.2f,%d,%d)"`, RadiansToDegrees(*c.textTheta), x, y))
	}
	bb.WriteRune('>')
	bb.WriteString(body)
	bb.WriteString("</text>")

	_, _ = c.w.Write(bb.Bytes())
}

func (c *canvas) Circle(x, y, r int, style Style) {
	bb := bytes.NewBuffer(make([]byte, 0, 80))
	bb.WriteString(`<circle cx="`)
	bb.WriteString(strconv.Itoa(x))
	bb.WriteString(`" cy="`)
	bb.WriteString(strconv.Itoa(y))
	bb.WriteString(`" r="`)
	bb.WriteString(strconv.Itoa(r))
	bb.WriteString(`" `)
	c.styleAsSVG(bb, style, true)
	bb.WriteString(`/>`)

	_, _ = c.w.Write(bb.Bytes())
}

func (c *canvas) End() {
	_, _ = c.w.Write([]byte("</svg>"))
}

// writeStrokeDashArray writes the stroke-dasharray property of a style.
func (c *canvas) writeStrokeDashArray(bb *bytes.Buffer, s Style) {
	if len(s.StrokeDashArray) > 0 {
		bb.WriteString("stroke-dasharray=\"")
		for i, v := range s.StrokeDashArray {
			if i > 0 {
				bb.WriteString(", ")
			}
			bb.WriteString(fmt.Sprintf("%0.1f", v))
		}
		bb.WriteString("\"")
	}
}

// GetFontFace returns the font face for the style.
func (c *canvas) getFontFace(s Style) string {
	family := "sans-serif"
	if s.GetFont() != nil {
		name := s.GetFont().Name(truetype.NameIDFontFamily)
		if name != "" {
			family = `'` + name + `',` + family
		}
	}
	return "font-family:" + family
}

// styleAsSVG returns the style as a svg style or class string.
func (c *canvas) styleAsSVG(bb *bytes.Buffer, s Style, applyText bool) {
	sw := s.StrokeWidth
	sc := s.StrokeColor
	fc := s.FillColor
	fs := s.FontSize
	fnc := s.FontColor

	if s.ClassName != "" {
		bb.WriteString("class=\"")
		bb.WriteString(s.ClassName)
		if !sc.IsZero() {
			bb.WriteRune(' ')
			bb.WriteString("stroke")
		}
		if !fc.IsZero() {
			bb.WriteRune(' ')
			bb.WriteString("fill")
		}
		if applyText && (fs != 0 || s.Font != nil) {
			bb.WriteRune(' ')
			bb.WriteString("text")
		}
		bb.WriteString("\"")
		return
	}

	bb.WriteString("style=\"")

	if sw != 0 && !sc.IsTransparent() {
		bb.WriteString("stroke-width:")
		bb.WriteString(formatFloatMinimized(sw))
		bb.WriteString(";stroke:")
		bb.WriteString(sc.String())
	} else {
		bb.WriteString("stroke:none")
	}

	if applyText && !fnc.IsTransparent() {
		bb.WriteString(";fill:")
		bb.WriteString(fnc.String())
	} else if !fc.IsTransparent() {
		bb.WriteString(";fill:")
		bb.WriteString(fc.String())
	} else {
		bb.WriteString(";fill:none")
	}

	if applyText {
		if fs != 0 {
			bb.WriteString(";font-size:")
			bb.WriteString(formatFloatMinimized(drawing.PointsToPixels(c.dpi, fs)))
			bb.WriteString("px")
		}
		if s.Font != nil {
			bb.WriteRune(';')
			bb.WriteString(c.getFontFace(s))
		}
	}

	bb.WriteRune('"')
}

// formatFloatNoTrailingZero formats a float without trailing zeros, so it is as small as possible.
func formatFloatMinimized(val float64) string {
	if val == float64(int(val)) {
		return strconv.Itoa(int(val))
	}
	str := fmt.Sprintf("%.1f", val)   // e.g. "1.20"
	str = strings.TrimRight(str, "0") // e.g. "1.2"
	str = strings.TrimRight(str, ".") // a rounding condition where an int is acceptable
	return str
}
