package charts

import (
	"bytes"
	"math"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/matrix"
)

// ValueFormatter defines a function that formats numeric values into string representations for display on charts.
type ValueFormatter func(float64) string

var defaultValueFormatter = func(val float64) string {
	return FormatValueHumanizeShort(val, 2, false)
}

func getPreferredValueFormatter(t ...ValueFormatter) ValueFormatter {
	for _, vf := range t {
		if vf != nil {
			return vf
		}
	}
	return defaultValueFormatter
}

// Painter is the primary struct for drawing charts/graphs.
type Painter struct {
	render       chartdraw.Renderer
	outputFormat string
	box          Box
	theme        ColorPalette
	font         *truetype.Font
}

// PainterOptions contains parameters for creating a new Painter.
type PainterOptions struct {
	// OutputFormat specifies the output type: "svg", "png", "jpg". Default is "png".
	OutputFormat string
	// Width is the width of the painter canvas.
	Width int
	// Height is the height of the painter canvas.
	Height int
	// Font is the default font for rendering text.
	Font *truetype.Font
	// Theme is the default theme used when charts don't specify one.
	Theme ColorPalette
}

// PainterOptionFunc defines a function that can modify a Painter after creation.
type PainterOptionFunc func(*Painter)

type ticksOption struct {
	firstIndex  int
	length      int
	vertical    bool
	tickCount   int
	tickSpaces  int
	strokeWidth float64
	strokeColor Color
}

type multiTextOption struct {
	textList       []string
	fontStyle      FontStyle
	vertical       bool
	centerLabels   bool
	align          string
	textRotation   float64
	offset         OffsetInt
	firstIndex     int
	labelCount     int
	labelSkipCount int
}

// PainterPaddingOption sets the padding within the painter canvas.
func PainterPaddingOption(padding Box) PainterOptionFunc {
	return func(p *Painter) {
		p.box.Left += padding.Left
		p.box.Top += padding.Top
		p.box.Right -= padding.Right
		p.box.Bottom -= padding.Bottom
	}
}

// PainterBoxOption sets a specific drawing area for the Painter.
func PainterBoxOption(box Box) PainterOptionFunc {
	return func(p *Painter) {
		if box.IsZero() {
			return
		}
		p.box = box
	}
}

// PainterThemeOption sets the default theme for the Painter.
// Used when specific chart options don't have a theme set.
func PainterThemeOption(theme ColorPalette) PainterOptionFunc {
	return func(p *Painter) {
		p.theme = getPreferredTheme(theme)
	}
}

// PainterFontOption sets the default font face for the Painter.
// Used when FontStyle in chart configs doesn't specify a font.
func PainterFontOption(font *truetype.Font) PainterOptionFunc {
	return func(p *Painter) {
		p.font = getPreferredFont(font)
	}
}

// NewPainter creates a painter for rendering charts.
func NewPainter(opts PainterOptions, opt ...PainterOptionFunc) *Painter {
	if opts.Width <= 0 {
		opts.Width = defaultChartWidth
	}
	if opts.Height <= 0 {
		opts.Height = defaultChartHeight
	}
	fn := chartdraw.PNG
	if opts.OutputFormat == ChartOutputJPG {
		fn = chartdraw.JPG
	} else if opts.OutputFormat == ChartOutputSVG {
		fn = chartdraw.SVG
	}

	p := &Painter{
		outputFormat: opts.OutputFormat,
		render:       fn(opts.Width, opts.Height),
		box: Box{
			Right:  opts.Width,
			Bottom: opts.Height,
			IsSet:  true,
		},
		font:  opts.Font,
		theme: opts.Theme,
	}
	p.setOptions(opt...)
	return p
}

func (p *Painter) setOptions(opts ...PainterOptionFunc) {
	for _, fn := range opts {
		fn(p)
	}
}

// Child returns a painter with the provided options applied. Useful for rendering relative to a portion of the canvas via PainterBoxOption.
func (p *Painter) Child(opt ...PainterOptionFunc) *Painter {
	child := &Painter{
		outputFormat: p.outputFormat,
		render:       p.render,
		box:          p.box.Clone(),
		theme:        p.theme,
		font:         p.font,
	}
	child.setOptions(opt...)
	return child
}

// Bytes returns the final rendered data as a byte slice.
func (p *Painter) Bytes() ([]byte, error) {
	buffer := bytes.Buffer{}
	if err := p.render.Save(&buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// moveTo sets the current path cursor to a given point.
func (p *Painter) moveTo(x, y int) {
	p.render.MoveTo(x+p.box.Left, y+p.box.Top)
}

// arcTo renders an arc from the current cursor.
func (p *Painter) arcTo(cx, cy int, rx, ry, startAngle, delta float64) {
	p.render.ArcTo(cx+p.box.Left, cy+p.box.Top, rx, ry, startAngle, delta)
}

// quadCurveTo draws a quadratic curve from the current cursor using a control point (cx, cy) and ending at (x, y).
func (p *Painter) quadCurveTo(cx, cy, x, y int) {
	p.render.QuadCurveTo(cx+p.box.Left, cy+p.box.Top, x+p.box.Left, y+p.box.Top)
}

// lineTo draws a line from the current path cursor to the given point.
func (p *Painter) lineTo(x, y int) {
	p.render.LineTo(x+p.box.Left, y+p.box.Top)
}

// close finalizes a shape as drawn by the current path.
func (p *Painter) close() {
	p.render.Close()
}

// stroke performs a stroke using the provided color and width, then resets style.
func (p *Painter) stroke(strokeColor Color, strokeWidth float64) {
	defer p.render.ResetStyle()
	p.render.SetStrokeColor(strokeColor)
	p.render.SetStrokeWidth(strokeWidth)
	p.render.Stroke()
}

// fill performs a fill with the given color, then resets style.
func (p *Painter) fill(fillColor Color) {
	defer p.render.ResetStyle()
	p.render.SetFillColor(fillColor)
	p.render.Fill()
}

// fillStroke performs a fill+stroke with the given colors and stroke width, then resets the style.
func (p *Painter) fillStroke(fillColor, strokeColor Color, strokeWidth float64) {
	defer p.render.ResetStyle()
	p.render.SetFillColor(fillColor)
	if strokeWidth > 0 && !strokeColor.IsTransparent() {
		p.render.SetStrokeColor(strokeColor)
		p.render.SetStrokeWidth(strokeWidth)
		p.render.FillStroke()
	} else {
		p.render.Fill()
	}
}

// Width returns the drawable width of the painter's box.
func (p *Painter) Width() int {
	return p.box.Width()
}

// Height returns the drawable height of the painter's box.
func (p *Painter) Height() int {
	return p.box.Height()
}

// MeasureText returns the rendered size of the text for the provided font style.
func (p *Painter) MeasureText(text string, textRotation float64, fontStyle FontStyle) Box {
	if text == "" || fontStyle.FontSize == 0 || fontStyle.FontColor.IsTransparent() {
		return BoxZero
	}
	if fontStyle.Font == nil {
		fontStyle.Font = getPreferredFont(p.font)
	}
	if textRotation != 0 {
		defer p.render.ClearTextRotation()
		p.render.SetTextRotation(textRotation)
	}
	defer p.render.ResetStyle()
	p.render.SetFont(fontStyle.Font)
	p.render.SetFontSize(fontStyle.FontSize)
	p.render.SetFontColor(fontStyle.FontColor)
	box := p.render.MeasureText(text)
	return box
}

func (p *Painter) measureTextMaxWidthHeight(textList []string, textRotation float64, fontStyle FontStyle) (int, int) {
	if fontStyle.Font == nil {
		fontStyle.Font = getPreferredFont(p.font)
	}
	var maxWidth, maxHeight int
	for _, text := range textList {
		box := p.MeasureText(text, textRotation, fontStyle)
		if maxWidth < box.Width() {
			maxWidth = box.Width()
		}
		if maxHeight < box.Height() {
			maxHeight = box.Height()
		}
	}
	return maxWidth, maxHeight
}

// Circle draws a circle at the given coords with a given radius.
func (p *Painter) Circle(radius float64, x, y int, fillColor, strokeColor Color, strokeWidth float64) {
	// This function has a slight behavior difference between png and svg.
	// We need to set the style attributes before the `Circle` call for SVG.
	defer p.render.ResetStyle()
	p.render.SetFillColor(fillColor)
	p.render.SetStrokeColor(strokeColor)
	p.render.SetStrokeWidth(strokeWidth)
	p.render.Circle(radius, x+p.box.Left, y+p.box.Top)
	p.render.FillStroke()
}

// LineStroke draws a line in the graph from point to point with the specified stroke color/width.
// Points with values of math.MaxInt32 are skipped, resulting in a gap.
// Single or isolated points result in just a dot being drawn at the point.
func (p *Painter) LineStroke(points []Point, strokeColor Color, strokeWidth float64) {
	valid := make([]Point, 0, len(points))
	for _, pt := range points {
		if pt.Y == math.MaxInt32 {
			// If we encounter a break, draw the accumulated segment
			if len(valid) > 0 {
				p.drawStraightPath(valid, true)
				p.stroke(strokeColor, strokeWidth)
				valid = valid[:0] // reset
			}
			continue
		}
		valid = append(valid, pt)
	}
	p.drawStraightPath(valid, true)
	p.stroke(strokeColor, strokeWidth)
}

// drawStraightPath draws a simple (non-curved) path for the given points.
// If dotForSinglePoint is true, single points are drawn as 2 px radius dots.
func (p *Painter) drawStraightPath(points []Point, dotForSinglePoint bool) {
	pointCount := len(points)
	if pointCount == 0 {
		return
	} else if pointCount == 1 {
		if dotForSinglePoint {
			p.render.Circle(2.0, points[0].X+p.box.Left, points[0].Y+p.box.Top)
		}
		return
	}
	p.moveTo(points[0].X, points[0].Y)
	for i := 1; i < pointCount; i++ {
		p.lineTo(points[i].X, points[i].Y)
	}
}

// SmoothLineStroke draws a smooth curve through the given points using Quadratic Bézier segments and a
// tension parameter in [0..1] with 0 providing straight lines between midpoints and 1 providing a smoother line.
// Because tension smooths out the line, the line no longer hits the provided points exactly. The more variable
// the points and the higher the tension, the more the line deviates.
func (p *Painter) SmoothLineStroke(points []Point, tension float64, strokeColor Color, strokeWidth float64) {
	if tension <= 0 {
		p.LineStroke(points, strokeColor, strokeWidth)
		return
	} else if tension > 1 {
		tension = 1
	}

	valid := make([]Point, 0, len(points)) // Slice to hold valid points between breaks
	for _, pt := range points {
		if pt.Y == math.MaxInt32 {
			// When a line break is found, draw the curve for the accumulated valid points if any
			if len(valid) > 0 {
				p.drawSmoothCurve(valid, tension, true)
				p.stroke(strokeColor, strokeWidth)
				valid = valid[:0] // reset
			}
			continue
		}
		valid = append(valid, pt)
	}
	// draw any remaining points collected
	p.drawSmoothCurve(valid, tension, true)
	p.stroke(strokeColor, strokeWidth)
}

// drawSmoothCurve handles the actual path drawing (MoveTo/LineTo/QuadCurveTo)
// but does NOT call Stroke() or Fill(), letting caller do it.
func (p *Painter) drawSmoothCurve(points []Point, tension float64, dotForSinglePoint bool) {
	if len(points) < 3 { // Not enough points to form a curve, draw a line
		p.drawStraightPath(points, dotForSinglePoint)
		return
	}

	p.moveTo(points[0].X, points[0].Y) // Start from the first valid point

	// Handle each segment between points with quadratic Bézier curves
	for i := 1; i < len(points)-1; i++ {
		x1, y1 := points[i].X, points[i].Y
		x2, y2 := points[i+1].X, points[i+1].Y

		mx := float64(x1+x2) / 2.0
		my := float64(y1+y2) / 2.0

		cx := float64(x1) + tension*(mx-float64(x1))
		cy := float64(y1) + tension*(my-float64(y1))

		p.quadCurveTo(x1, y1, int(cx), int(cy))
	}

	// Connect the second-to-last point to the last point
	n := len(points)
	p.quadCurveTo(points[n-2].X, points[n-2].Y, points[n-1].X, points[n-1].Y)
}

// DashedLineStroke draws line segments through the provided points with dashed strokes.
// Points with values of math.MaxInt32 are skipped, resulting in a gap.
func (p *Painter) DashedLineStroke(points []Point, strokeColor Color, strokeWidth float64, strokeDashArray []float64) {
	valid := make([]Point, 0, len(points))
	for _, pt := range points {
		if pt.Y == math.MaxInt32 {
			// If we encounter a break, draw the accumulated segment
			if len(valid) > 0 {
				p.drawStraightPath(valid, true)
				p.dashedStroke(strokeColor, strokeWidth, strokeDashArray)
				valid = valid[:0] // reset
			}
			continue
		}
		valid = append(valid, pt)
	}
	p.drawStraightPath(valid, true)
	p.dashedStroke(strokeColor, strokeWidth, strokeDashArray)
}

// SmoothDashedLineStroke draws a smooth dashed curve through the given points.
func (p *Painter) SmoothDashedLineStroke(points []Point, tension float64, strokeColor Color, strokeWidth float64, strokeDashArray []float64) {
	if tension <= 0 {
		p.DashedLineStroke(points, strokeColor, strokeWidth, strokeDashArray)
		return
	} else if tension > 1 {
		tension = 1
	}

	valid := make([]Point, 0, len(points))
	for _, pt := range points {
		if pt.Y == math.MaxInt32 {
			// If we encounter a break, draw the accumulated segment
			if len(valid) > 0 {
				p.drawSmoothCurve(valid, tension, true)
				p.dashedStroke(strokeColor, strokeWidth, strokeDashArray)
				valid = valid[:0] // reset
			}
			continue
		}
		valid = append(valid, pt)
	}
	p.drawSmoothCurve(valid, tension, true)
	p.dashedStroke(strokeColor, strokeWidth, strokeDashArray)
}

// dashedStroke applies the stroke with the given dash array.
func (p *Painter) dashedStroke(strokeColor Color, strokeWidth float64, strokeDashArray []float64) {
	defer p.render.ResetStyle()
	p.render.SetStrokeColor(strokeColor)
	p.render.SetStrokeWidth(strokeWidth)
	p.render.SetStrokeDashArray(strokeDashArray)
	p.render.Stroke()
}

// drawBackground fills the entire painter area with the given color.
func (p *Painter) drawBackground(color Color) {
	p.FilledRect(0, 0, p.Width(), p.Height(), color, color, 0.0)
}

// FilledRect draws a filled box with the given coordinates.
func (p *Painter) FilledRect(x1, y1, x2, y2 int, fillColor, strokeColor Color, strokeWidth float64) {
	p.rectMoveLine(x1, y1, x2, y2)
	p.fillStroke(fillColor, strokeColor, strokeWidth)
}

func (p *Painter) rectMoveLine(x1, y1, x2, y2 int) {
	p.moveTo(x1, y1)
	p.lineTo(x2, y1)
	p.lineTo(x2, y2)
	p.lineTo(x1, y2)
	p.lineTo(x1, y1)
}

// FilledDiamond draws a filled diamond centered at (cx, cy) with the given width and height.
func (p *Painter) FilledDiamond(cx, cy, width, height int, fillColor, strokeColor Color, strokeWidth float64) {
	p.diamondMoveLine(cx, cy, width, height)
	p.fillStroke(fillColor, strokeColor, strokeWidth)
}

func (p *Painter) diamondMoveLine(cx, cy, width, height int) {
	// Calculate the four corners of the diamond
	hw, hh := width/2, height/2 // Half-width and height
	p1x, p1y := cx, cy-hh       // Top
	p2x, p2y := cx+hw, cy       // Right
	p3x, p3y := cx, cy+hh       // Bottom
	p4x, p4y := cx-hw, cy       // Left

	p.moveTo(p1x, p1y)
	p.lineTo(p2x, p2y)
	p.lineTo(p3x, p3y)
	p.lineTo(p4x, p4y)
	p.lineTo(p1x, p1y)
}

// HorizontalMarkLine draws a horizontal line with a small circle and arrow at the right.
func (p *Painter) HorizontalMarkLine(x, y, width int, fillColor, strokeColor Color, strokeWidth float64, strokeDashArray []float64) {
	const arrowWidth = 16
	const arrowHeight = 10
	endX := x + width
	radius := 3

	// Set up stroke style before drawing
	defer p.render.ResetStyle()
	p.render.SetStrokeColor(strokeColor)
	p.render.SetStrokeWidth(strokeWidth)
	p.render.SetStrokeDashArray(strokeDashArray)
	p.render.SetFillColor(fillColor)

	// Draw the circle at the starting point
	p.render.Circle(float64(radius), x+radius+p.box.Left, y+p.box.Top)
	p.render.Fill() // only fill the circle, do not stroke

	// Draw the line from the end of the circle to near the arrow start
	p.moveTo(x+radius*3, y)
	p.lineTo(endX-arrowWidth, y)
	p.render.Stroke() // apply stroke with the dash array

	p.ArrowRight(endX, y, arrowWidth, arrowHeight, fillColor, strokeColor, strokeWidth)
}

// VerticalMarkLine draws a vertical line with a small dot at the bottom and an arrow at the top.
func (p *Painter) VerticalMarkLine(x, y, height int, fillColor, strokeColor Color, strokeWidth float64, strokeDashArray []float64) {
	const arrowHeight = 16
	const arrowWidth = 10
	endY := y + height
	radius := 3

	// Set up stroke style before drawing
	defer p.render.ResetStyle()
	p.render.SetStrokeColor(strokeColor)
	p.render.SetStrokeWidth(strokeWidth)
	p.render.SetStrokeDashArray(strokeDashArray)
	p.render.SetFillColor(fillColor)

	// Draw the dot at the bottom of the line
	p.render.Circle(float64(radius), x+p.box.Left, endY-radius+p.box.Top)
	p.render.Fill() // fill the circle

	// Draw the vertical line
	p.moveTo(x, y)
	p.lineTo(x, endY)
	p.render.Stroke() // apply stroke with the dash array

	p.ArrowUp(x, y+arrowHeight, arrowWidth, arrowHeight, fillColor, strokeColor, strokeWidth)
}

// Polygon draws a polygon with the specified center, radius, and number of sides.
func (p *Painter) Polygon(center Point, radius float64, sides int, strokeColor Color, strokeWidth float64) {
	points := getPolygonPoints(center, radius, sides)
	p.drawStraightPath(points, false)
	p.lineTo(points[0].X, points[0].Y)
	p.stroke(strokeColor, strokeWidth)
}

const (
	_pi2  = math.Pi / 2.0
	_2pi  = 2 * math.Pi
	_3pi2 = (3 * math.Pi) / 2.0
)

// Pin draws a pin shape (circle + curved tail).
func (p *Painter) Pin(x, y, width int, fillColor, strokeColor Color, strokeWidth float64) {
	r := float64(width) / 2
	y -= width / 4
	angle := DegreesToRadians(15)

	// Draw the pin head with fill and stroke
	startAngle := _pi2 + angle
	delta := _2pi - 2*angle
	p.arcTo(x, y, r, r, startAngle, delta)
	p.lineTo(x, y)
	p.close()
	p.fillStroke(fillColor, strokeColor, strokeWidth)

	// The curved tail
	startX := x - int(r)
	startY := y
	endX := x + int(r)
	endY := y
	p.moveTo(startX, startY)
	cx := x
	cy := y + int(r*2.5)
	p.quadCurveTo(cx, cy, endX, endY)
	p.close()

	// Apply both fill and stroke to the tail
	p.fillStroke(fillColor, strokeColor, strokeWidth)
}

// arrow draws an arrow shape in the given direction, then fill+stroke with the given style.
func (p *Painter) arrow(x, y, width, height int, direction string,
	fillColor, strokeColor Color, strokeWidth float64) {
	halfWidth := width >> 1
	halfHeight := height >> 1
	if direction == PositionTop || direction == PositionBottom {
		x0 := x - halfWidth
		x1 := x0 + width
		dy := -height / 3
		y0 := y
		y1 := y0 - height
		if direction == PositionBottom {
			y0 = y - height
			y1 = y
			dy = 2 * dy
		}
		p.moveTo(x0, y0)
		p.lineTo(x0+halfWidth, y1)
		p.lineTo(x1, y0)
		p.lineTo(x0+halfWidth, y+dy)
		p.lineTo(x0, y0)
	} else {
		x0 := x + width
		x1 := x0 - width
		y0 := y - halfHeight
		dx := -width / 3
		if direction == PositionRight {
			x0 = x - width
			dx = -dx
			x1 = x0 + width
		}
		p.moveTo(x0, y0)
		p.lineTo(x1, y0+halfHeight)
		p.lineTo(x0, y0+height)
		p.lineTo(x0+dx, y0+halfHeight)
		p.lineTo(x0, y0)
	}
	p.fillStroke(fillColor, strokeColor, strokeWidth)
}

// ArrowLeft draws an arrow at the given point and dimensions pointing left.
func (p *Painter) ArrowLeft(x, y, width, height int,
	fillColor, strokeColor Color, strokeWidth float64) {
	p.arrow(x, y, width, height, PositionLeft, fillColor, strokeColor, strokeWidth)
}

// ArrowRight draws an arrow at the given point and dimensions pointing right.
func (p *Painter) ArrowRight(x, y, width, height int,
	fillColor, strokeColor Color, strokeWidth float64) {
	p.arrow(x, y, width, height, PositionRight, fillColor, strokeColor, strokeWidth)
}

// ArrowUp draws an arrow at the given point and dimensions pointing up.
func (p *Painter) ArrowUp(x, y, width, height int,
	fillColor, strokeColor Color, strokeWidth float64) {
	p.arrow(x, y, width, height, PositionTop, fillColor, strokeColor, strokeWidth)
}

// ArrowDown draws an arrow at the given point and dimensions pointing down.
func (p *Painter) ArrowDown(x, y, width, height int,
	fillColor, strokeColor Color, strokeWidth float64) {
	p.arrow(x, y, width, height, PositionBottom, fillColor, strokeColor, strokeWidth)
}

// FillArea draws a filled polygon through the given points, skipping "null" (MaxInt32) break values
// (filling the area flat between them).
func (p *Painter) FillArea(points []Point, fillColor Color) {
	if len(points) == 0 {
		return
	}

	valid := make([]Point, 0, len(points))
	for _, pt := range points {
		if pt.Y == math.MaxInt32 {
			// If we encounter a break, fill the accumulated segment
			if len(valid) > 0 {
				p.drawStraightPath(valid, false)
				p.fill(fillColor)
				valid = valid[:0] // reset
			}
			continue
		}
		valid = append(valid, pt)
	}

	// Fill the last segment if there is one
	p.drawStraightPath(valid, false)
	p.fill(fillColor)
}

// smoothFillChartArea draws a smooth curve for the "top" portion of points but uses straight lines for
// the bottom corners, producing a fill with sharp corners.
func (p *Painter) smoothFillChartArea(points []Point, tension float64, fillColor Color) {
	pointCount := len(points)
	if tension <= 0 || pointCount < 4 /* need at least 4 points to curve the line */ {
		p.FillArea(points, fillColor)
		return
	} else if tension > 1 {
		tension = 1
	}

	// Typically, areaPoints has the shape:
	//   [ top data points... ] + [ bottom-right corner, bottom-left corner, first top point ]
	// The final 3 points are the corners + repeated first point
	top := points[:pointCount-3]
	bottom := points[pointCount-3:] // [corner1, corner2, firstTopAgain]

	// If top portion is empty or 1 point, just fill straight
	if len(top) < 2 {
		p.FillArea(points, fillColor)
		return
	}

	// Build the smooth path for the top portion
	currentSegment := make([]Point, 0, len(top))
	firstPointSet := false
	for _, pt := range top {
		if pt.Y == math.MaxInt32 {
			// If we encounter a break, fill the accumulated segment
			if len(currentSegment) > 0 {
				p.drawSmoothCurve(currentSegment, tension, false)
				firstPointSet = true
				currentSegment = currentSegment[:0]
			}
			continue
		}
		currentSegment = append(currentSegment, pt)
	}

	// Draw the remaining top segment
	if len(currentSegment) > 0 {
		p.drawSmoothCurve(currentSegment, tension, false)
		firstPointSet = true
	}

	if !firstPointSet {
		p.FillArea(points, fillColor) // No actual top segment was drawn, fallback to straight fill
		return
	}

	// Add sharp lines to close the shape at the bottom
	// The path is currently at the last top point we drew. Now we need to draw to corner1 -> corner2 -> firstTopAgain
	for i := 0; i < len(bottom); i++ {
		p.lineTo(bottom[i].X, bottom[i].Y)
	}
	p.fill(fillColor)
}

// Text draws the given string at the specified position using the given font style.
// Specifying radians rotates the text.
func (p *Painter) Text(body string, x, y int, radians float64, fontStyle FontStyle) {
	if fontStyle.Font == nil {
		fontStyle.Font = getPreferredFont(p.font)
	}
	defer p.render.ResetStyle()
	p.render.SetFont(fontStyle.Font)
	p.render.SetFontSize(fontStyle.FontSize)
	p.render.SetFontColor(fontStyle.FontColor)

	if radians != 0 {
		defer p.render.ClearTextRotation()
		p.render.SetTextRotation(radians)
	}
	p.render.Text(body, x+p.box.Left, y+p.box.Top)
}

// TextFit draws multi-line text constrained to a given width.
func (p *Painter) TextFit(body string, x, y, width int, fontStyle FontStyle, textAligns ...string) Box {
	if fontStyle.Font == nil {
		fontStyle.Font = getPreferredFont(p.font)
	}
	style := chartdraw.Style{
		FontStyle: fontStyle,
		TextWrap:  chartdraw.TextWrapWord,
	}
	r := p.render
	defer r.ResetStyle()
	r.SetFont(fontStyle.Font)
	r.SetFontSize(fontStyle.FontSize)
	r.SetFontColor(fontStyle.FontColor)

	lines := chartdraw.Text.WrapFit(r, body, width, style)

	var output Box
	var textAlign string
	if len(textAligns) != 0 {
		textAlign = textAligns[0]
	}

	for index, line := range lines {
		if line == "" {
			continue
		}
		x0 := x
		y0 := y + output.Height()
		lineBox := r.MeasureText(line)
		switch textAlign {
		case AlignRight:
			x0 += width - lineBox.Width()
		case AlignCenter:
			x0 += (width - lineBox.Width()) >> 1
		}

		p.render.Text(line, x0+p.box.Left, y0+p.box.Top)
		output.Right = chartdraw.MaxInt(lineBox.Right, output.Right)
		output.Bottom += lineBox.Height()
		if index < len(lines)-1 {
			output.Bottom += style.GetTextLineSpacing()
		}
	}
	output.IsSet = true
	return output
}

// isTick determines whether the given index is a "tick" mark out of numTicks.
func isTick(totalRange int, numTicks int, index int) bool {
	if numTicks >= totalRange {
		return true
	} else if index == 0 || index == totalRange-1 {
		return true // shortcut to always define tick at start and end of range
	}
	step := float64(totalRange-1) / float64(numTicks-1)
	// predictedTickIndex calculates the nearest theoretical tick position based on a continuous scale.
	// It divides the current index by the step size to determine how many ticks fit into the index,
	// then rounds to the nearest whole number to find the closest tick index.
	predictedTickIndex := int(float64(index)/step + 0.5)
	// actualTickIndex translates the predictedTickIndex back to the actual data index.
	// It does this by multiplying the predictedTickIndex by the step size, effectively finding
	// the actual position of this tick on the discrete scale of data indices. It rounds it
	// to ensure it aligns with an exact index in the array.
	actualTickIndex := int(float64(predictedTickIndex)*step + 0.5)
	return actualTickIndex == index
}

// ticks draws small lines to indicate tick marks, using a fixed stroke color/width.
func (p *Painter) ticks(opt ticksOption) {
	if opt.tickCount <= 0 || opt.length <= 0 {
		return
	}
	var values []int
	if opt.vertical {
		values = autoDivide(p.Height(), opt.tickSpaces)
	} else {
		values = autoDivide(p.Width(), opt.tickSpaces)
	}
	for index, value := range values {
		if index < opt.firstIndex {
			continue
		} else if !isTick(len(values)-opt.firstIndex, opt.tickCount, index-opt.firstIndex) {
			continue
		}
		if opt.vertical {
			p.LineStroke([]Point{
				{X: 0, Y: value},
				{X: opt.length, Y: value},
			}, opt.strokeColor, opt.strokeWidth)
		} else {
			p.LineStroke([]Point{
				{X: value, Y: opt.length},
				{X: value, Y: 0},
			}, opt.strokeColor, opt.strokeWidth)
		}
	}
}

// multiText prints multiple lines of text for axis labels.
func (p *Painter) multiText(opt multiTextOption) {
	if len(opt.textList) == 0 {
		return
	}
	count := len(opt.textList)
	width := p.Width()
	height := p.Height()
	var positions []int
	if opt.vertical {
		if opt.centerLabels {
			positions = autoDivide(height, count)
		} else {
			positions = autoDivide(height, count-1)
		}
	} else {
		if opt.centerLabels {
			positions = autoDivide(width, count)
		} else {
			positions = autoDivide(width, count-1)
		}
	}
	if opt.textRotation != 0 {
		defer p.render.ClearTextRotation()
		p.render.SetTextRotation(opt.textRotation)
	}
	positionCount := len(positions)
	tickCount := opt.labelCount
	if opt.centerLabels {
		tickCount++
	}

	skippedLabels := opt.labelSkipCount // specify the skip count to ensure the top value is listed
	for index, start := range positions {
		if opt.centerLabels && index == positionCount-1 {
			break // positions have one item more than we can map to text, this extra value is used to center against
		} else if index < opt.firstIndex {
			continue
		} else if index != count-1 && // one off case for last label due to values and label qty difference
			!isTick(positionCount-opt.firstIndex, tickCount, index-opt.firstIndex) {
			continue
		} else if index != count-1 && // ensure the bottom value is always printed
			skippedLabels < opt.labelSkipCount {
			skippedLabels++
			continue
		} else {
			skippedLabels = 0
		}

		text := opt.textList[index]
		box := p.MeasureText(text, opt.textRotation, opt.fontStyle)
		var x, y int
		if opt.vertical {
			if opt.centerLabels {
				start = (positions[index] + positions[index+1]) >> 1
			} else {
				start = positions[index]
			}
			y = start + (box.Height() >> 1)
			switch opt.align {
			case AlignRight:
				x = width - box.Width()
			case AlignCenter:
				x = width - (box.Width() >> 1)
			default:
				x = 0
			}
		} else {
			if opt.centerLabels {
				// Graphs with limited data samples generally look better with the samples directly below the label.
				// For that reason, we will exactly center these graphs, but graphs with higher sample counts will
				// attempt to space the labels better rather than line up directly to the graph points.
				exactLabels := count == opt.labelCount
				if !exactLabels && index == 0 {
					x = start - 1 // align to the actual start (left side of tick space)
				} else if !exactLabels && index == count-1 {
					x = width - box.Width() // align to the right side of tick space
				} else {
					start = (positions[index] + positions[index+1]) >> 1
					x = start - box.Width()>>1 // align to center of tick space
				}
			} else {
				if index == count-1 {
					x = width - box.Width() // align to the right side of tick space
				} else {
					x = start - 1 // align to the left side of the tick space
				}
			}
		}
		x += opt.offset.Left
		y += opt.offset.Top
		p.Text(text, x, y, opt.textRotation, opt.fontStyle)
	}
}

// textRotationHeightAdjustment calculates how much vertical adjustment is needed
// after rotating the text around the bottom-right corner.
//
// The caller typically subtracts this returned value from the existing y-position so that the text
// stays aligned with the bottom position. For this calculation, the provided text dimensions should be
// WITHOUT rotation applied.
func textRotationHeightAdjustment(textWidth, textHeight int, radians float64) int {
	r := normalizeAngle(radians)

	switch {
	// Very close to 0 radians: no vertical adjustment needed
	case r < matrix.DefaultEpsilon:
		return 0
	// 0 to π (0 to 180 degrees)
	case r < math.Pi:
		// Compute vertical displacement needed to maintain alignment at the bottom
		// sin(r) gives the vertical component of the text width as it rotates
		return int(math.Round(float64(textWidth) * math.Sin(r)))
	// π to 3π/2 (180 to 270 degrees)
	case r >= math.Pi && r < _3pi2:
		// Adjust the text downward as it rotates past 180 degrees
		// cos(angle) gives the horizontal overlap, subtract from height to get adjustment
		return textHeight - int(math.Round(float64(textHeight)*math.Cos(r-_3pi2)))
	// 3π/2 to 2π (270 to 360 degrees)
	default:
		// No adjustment needed as the text aligns back towards zero position
		return 0
	}
}

// Dots prints filled circles for the given points.
func (p *Painter) Dots(points []Point, fillColor, strokeColor Color, strokeWidth float64, dotRadius float64) {
	defer p.render.ResetStyle()
	p.render.SetFillColor(fillColor)
	p.render.SetStrokeColor(strokeColor)
	p.render.SetStrokeWidth(strokeWidth)
	for _, item := range points {
		p.render.Circle(dotRadius, item.X+p.box.Left, item.Y+p.box.Top)
	}
	p.render.FillStroke()
}

// squares prints filled squares for the given points.
func (p *Painter) squares(points []Point, fillColor, strokeColor Color, strokeWidth float64, size int) {
	defer p.render.ResetStyle()
	p.render.SetFillColor(fillColor)
	p.render.SetStrokeColor(strokeColor)
	p.render.SetStrokeWidth(strokeWidth)
	halfSize := int(math.Round(float64(size) / 2.0))
	for _, item := range points {
		x1 := item.X - halfSize
		y1 := item.Y - halfSize
		p.rectMoveLine(x1, y1, x1+size, y1+size)
	}
	p.render.FillStroke()
}

// diamonds prints filled diamonds for the given points.
func (p *Painter) diamonds(points []Point, fillColor, strokeColor Color, strokeWidth float64, size int) {
	defer p.render.ResetStyle()
	p.render.SetFillColor(fillColor)
	p.render.SetStrokeColor(strokeColor)
	p.render.SetStrokeWidth(strokeWidth)
	for _, item := range points {
		p.diamondMoveLine(item.X, item.Y, size, size)
	}
	p.render.FillStroke()
}

// roundedRect is similar to filledRect except the top and bottom are rounded.
func (p *Painter) roundedRect(box Box, radius int, roundTop, roundBottom bool,
	fillColor, strokeColor Color, strokeWidth float64) {
	r := (box.Right - box.Left) / 2
	if radius > r {
		radius = r
	}
	rx := float64(radius)
	ry := float64(radius)

	if roundTop {
		// Start at the appropriate point depending on rounding at the top
		p.moveTo(box.Left+radius, box.Top)
		p.lineTo(box.Right-radius, box.Top)

		// right top
		cx := box.Right - radius
		cy := box.Top + radius
		p.arcTo(cx, cy, rx, ry, -_pi2, _pi2)
	} else {
		p.moveTo(box.Left, box.Top)
		p.lineTo(box.Right, box.Top)
	}

	if roundBottom {
		p.lineTo(box.Right, box.Bottom-radius)

		// right bottom
		cx := box.Right - radius
		cy := box.Bottom - radius
		p.arcTo(cx, cy, rx, ry, 0, _pi2)

		p.lineTo(box.Left+radius, box.Bottom)

		// left bottom
		cx = box.Left + radius
		cy = box.Bottom - radius
		p.arcTo(cx, cy, rx, ry, _pi2, _pi2)
	} else {
		p.lineTo(box.Right, box.Bottom)
		p.lineTo(box.Left, box.Bottom)
	}

	if roundTop {
		// left top
		p.lineTo(box.Left, box.Top+radius)
		cx := box.Left + radius
		cy := box.Top + radius
		p.arcTo(cx, cy, rx, ry, math.Pi, _pi2)
	} else {
		p.lineTo(box.Left, box.Top)
	}

	p.close()
	p.fillStroke(fillColor, strokeColor, strokeWidth)
}

// legendLineDot draws a small horizontal line with a dot in the middle, often used in legends.
func (p *Painter) legendLineDot(box Box, strokeColor Color, strokeWidth float64, dotColor, centerColor Color) {
	center := (box.Height()-int(strokeWidth))>>1 - 1

	defer p.render.ResetStyle()
	p.render.SetStrokeColor(strokeColor)
	p.render.SetStrokeWidth(strokeWidth)
	p.moveTo(box.Left, box.Top-center)
	p.lineTo(box.Right, box.Top-center)
	p.render.Stroke()

	// draw dot in the middle
	midX := box.Left + (box.Width() >> 1)
	p.Circle(5, midX, box.Top-center, dotColor, dotColor, 3)
	if dotColor != centerColor && !centerColor.IsTransparent() {
		p.Circle(2, midX, box.Top-center, centerColor, centerColor, 3)
	}
}

// BarChart renders a bar chart with the provided configuration to the painter.
func (p *Painter) BarChart(opt BarChartOption) error {
	_, err := newBarChart(p, opt).Render()
	return err
}

// HorizontalBarChart renders a horizontal bar chart with the provided configuration to the painter.
func (p *Painter) HorizontalBarChart(opt HorizontalBarChartOption) error {
	_, err := newHorizontalBarChart(p, opt).Render()
	return err
}

// FunnelChart renders a funnel chart with the provided configuration to the painter.
func (p *Painter) FunnelChart(opt FunnelChartOption) error {
	_, err := newFunnelChart(p, opt).Render()
	return err
}

// LineChart renders a line chart with the provided configuration to the painter.
func (p *Painter) LineChart(opt LineChartOption) error {
	_, err := newLineChart(p, opt).Render()
	return err
}

// ScatterChart renders a scatter chart with the provided configuration to the painter.
func (p *Painter) ScatterChart(opt ScatterChartOption) error {
	_, err := newScatterChart(p, opt).Render()
	return err
}

// PieChart renders a pie chart with the provided configuration to the painter.
func (p *Painter) PieChart(opt PieChartOption) error {
	_, err := newPieChart(p, opt).Render()
	return err
}

// DoughnutChart renders a doughnut or ring chart with the provided configuration to the painter.
func (p *Painter) DoughnutChart(opt DoughnutChartOption) error {
	_, err := newDoughnutChart(p, opt).Render()
	return err
}

// RadarChart renders a radar chart with the provided configuration to the painter.
func (p *Painter) RadarChart(opt RadarChartOption) error {
	_, err := newRadarChart(p, opt).Render()
	return err
}

// HeatMapChart renders a heat map with the provided configuration to the painter.
func (p *Painter) HeatMapChart(opt HeatMapOption) error {
	_, err := newHeatMapChart(p, opt).Render()
	return err
}

// TableChart renders a table with the provided configuration to the painter.
func (p *Painter) TableChart(opt TableChartOption) error {
	_, err := newTableChart(p, opt).Render()
	return err
}

// CandlestickChart renders a candlestick chart with the provided configuration to the painter.
func (p *Painter) CandlestickChart(opt CandlestickChartOption) error {
	_, err := newCandlestickChart(p, opt).Render()
	return err
}
