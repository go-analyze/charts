package charts

import (
	"bytes"
	"math"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

// ValueFormatter defines a function that can be used to format numeric values.
type ValueFormatter func(float64) string

var defaultValueFormatter = func(val float64) string {
	return FormatValueHumanizeShort(val, 2, false)
}

// Painter is the primary struct for drawing charts/graphs.
type Painter struct {
	render         chartdraw.Renderer
	box            Box
	style          chartdraw.Style
	theme          ColorPalette
	font           *truetype.Font
	outputFormat   string
	valueFormatter ValueFormatter
}

// PainterOptions contains parameters for creating a new Painter.
type PainterOptions struct {
	// OutputFormat specifies the output type, "svg" or "png", default is "png".
	OutputFormat string
	// Width is the width of the draw painter.
	Width int
	// Height is the height of the draw painter.
	Height int
	// Font is the font used for rendering text.
	Font *truetype.Font
}

type PainterOptionFunc func(*Painter)

type TicksOption struct {
	// the first tick index
	First      int
	Length     int
	Vertical   bool
	LabelCount int
	TickCount  int
	TickSpaces int
}

type MultiTextOption struct {
	TextList     []string
	Vertical     bool
	CenterLabels bool
	Align        string
	TextRotation float64
	Offset       OffsetInt
	// The first text index
	First          int
	LabelCount     int
	TickCount      int
	LabelSkipCount int
}

type GridOption struct {
	// Columns is the count of columns in the grid.
	Columns int
	// Rows are the count of rows in the grid.
	Rows int
	// ColumnSpans specifies the span for each column.
	ColumnSpans []int
	// IgnoreColumnLines specifies index for columns to not display.
	IgnoreColumnLines []int
	// IgnoreRowLines specifies index for rows to not display.
	IgnoreRowLines []int
}

// PainterPaddingOption sets the padding of draw painter.
func PainterPaddingOption(padding Box) PainterOptionFunc {
	return func(p *Painter) {
		p.box.Left += padding.Left
		p.box.Top += padding.Top
		p.box.Right -= padding.Right
		p.box.Bottom -= padding.Bottom
	}
}

// PainterBoxOption sets a specific box for the Painter to draw within.
func PainterBoxOption(box Box) PainterOptionFunc {
	return func(p *Painter) {
		if box.IsZero() {
			return
		}
		p.box = box
	}
}

// PainterThemeOption sets a color palette theme default for the Painter.
// This theme is used if the specific chart options don't have a theme set.
func PainterThemeOption(theme ColorPalette) PainterOptionFunc {
	return func(p *Painter) {
		if theme == nil {
			theme = GetDefaultTheme()
		}
		p.theme = theme
	}
}

// TODO - try to remove the error return
// NewPainter creates a painter which can be used to render charts to (using for example NewLineChart).
func NewPainter(opts PainterOptions, opt ...PainterOptionFunc) (*Painter, error) {
	if opts.Width <= 0 {
		opts.Width = defaultChartWidth
	}
	if opts.Height <= 0 {
		opts.Height = defaultChartHeight
	}
	if opts.Font == nil {
		opts.Font = GetDefaultFont()
	}
	fn := chartdraw.PNG
	if opts.OutputFormat == ChartOutputSVG {
		fn = chartdraw.SVG
	}
	width := opts.Width
	height := opts.Height
	r, err := fn(width, height)
	if err != nil {
		return nil, err
	}
	r.SetFont(opts.Font)

	p := &Painter{
		render: r,
		box: Box{
			Right:  opts.Width,
			Bottom: opts.Height,
			IsSet:  true,
		},
		font:         opts.Font,
		outputFormat: opts.OutputFormat,
	}
	p.setOptions(opt...)
	if p.theme == nil {
		p.theme = GetDefaultTheme()
	}
	return p, nil
}

func (p *Painter) setOptions(opts ...PainterOptionFunc) {
	for _, fn := range opts {
		fn(p)
	}
}

// Child returns a painter with the passed in options applied to it. This can be most useful when you want to render
// relative to only a portion of the canvas using PainterBoxOption.
func (p *Painter) Child(opt ...PainterOptionFunc) *Painter {
	child := &Painter{
		render:         p.render,
		box:            p.box.Clone(),
		style:          p.style,
		theme:          p.theme,
		font:           p.font,
		outputFormat:   p.outputFormat,
		valueFormatter: p.valueFormatter,
	}
	child.setOptions(opt...)
	return child
}

func (p *Painter) SetStyle(style chartdraw.Style) {
	if style.Font == nil {
		style.Font = p.font
	}
	p.style = style
	style.WriteToRenderer(p.render)
}

func overrideStyle(defaultStyle chartdraw.Style, style chartdraw.Style) chartdraw.Style {
	if style.StrokeWidth == 0 {
		style.StrokeWidth = defaultStyle.StrokeWidth
	}
	if style.StrokeColor.IsZero() {
		style.StrokeColor = defaultStyle.StrokeColor
	}
	if style.StrokeDashArray == nil {
		style.StrokeDashArray = defaultStyle.StrokeDashArray
	}
	if style.DotColor.IsZero() {
		style.DotColor = defaultStyle.DotColor
	}
	if style.DotWidth == 0 {
		style.DotWidth = defaultStyle.DotWidth
	}
	if style.FillColor.IsZero() {
		style.FillColor = defaultStyle.FillColor
	}
	if style.FontSize == 0 {
		style.FontSize = defaultStyle.FontSize
	}
	if style.FontColor.IsZero() {
		style.FontColor = defaultStyle.FontColor
	}
	if style.Font == nil {
		style.Font = defaultStyle.Font
	}
	return style
}

func (p *Painter) OverrideDrawingStyle(style chartdraw.Style) *Painter {
	// TODO - we should alias parts of Style we want to support drawing on
	s := overrideStyle(p.style, style)
	p.SetDrawingStyle(s)
	return p
}

func (p *Painter) SetDrawingStyle(style chartdraw.Style) *Painter {
	style.WriteDrawingOptionsToRenderer(p.render)
	return p
}

func (p *Painter) SetFontStyle(style chartdraw.FontStyle) *Painter {
	if style.Font == nil {
		style.Font = p.font
	}
	if style.FontColor.IsZero() {
		style.FontColor = p.style.FontColor
	}
	if style.FontSize == 0 {
		style.FontSize = p.style.FontSize
	}
	style.WriteTextOptionsToRenderer(p.render)
	return p
}
func (p *Painter) OverrideFontStyle(style chartdraw.FontStyle) *Painter {
	s := overrideStyle(p.style, chartdraw.Style{FontStyle: style})
	p.SetFontStyle(s.FontStyle)
	return p
}

func (p *Painter) ResetStyle() *Painter {
	p.style.WriteToRenderer(p.render)
	return p
}

// Bytes returns the data of draw canvas
func (p *Painter) Bytes() ([]byte, error) {
	buffer := bytes.Buffer{}
	if err := p.render.Save(&buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// MoveTo moves the cursor to a given point
func (p *Painter) MoveTo(x, y int) *Painter {
	p.render.MoveTo(x+p.box.Left, y+p.box.Top)
	return p
}

func (p *Painter) ArcTo(cx, cy int, rx, ry, startAngle, delta float64) *Painter {
	p.render.ArcTo(cx+p.box.Left, cy+p.box.Top, rx, ry, startAngle, delta)
	return p
}

func (p *Painter) LineTo(x, y int) *Painter {
	p.render.LineTo(x+p.box.Left, y+p.box.Top)
	return p
}

func (p *Painter) QuadCurveTo(cx, cy, x, y int) *Painter {
	p.render.QuadCurveTo(cx+p.box.Left, cy+p.box.Top, x+p.box.Left, y+p.box.Top)
	return p
}

func (p *Painter) Pin(x, y, width int) *Painter {
	r := float64(width) / 2
	y -= width / 4
	angle := chartdraw.DegreesToRadians(15)
	box := p.box

	startAngle := math.Pi/2 + angle
	delta := 2*math.Pi - 2*angle
	p.ArcTo(x, y, r, r, startAngle, delta)
	p.LineTo(x, y)
	p.Close()
	p.FillStroke()

	startX := x - int(r)
	startY := y
	endX := x + int(r)
	endY := y
	p.MoveTo(startX, startY)

	left := box.Left
	top := box.Top
	cx := x
	cy := y + int(r*2.5)
	p.render.QuadCurveTo(cx+left, cy+top, endX+left, endY+top)
	p.Close()
	p.Fill()
	return p
}

func (p *Painter) arrow(x, y, width, height int, direction string) *Painter {
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
		p.MoveTo(x0, y0)
		p.LineTo(x0+halfWidth, y1)
		p.LineTo(x1, y0)
		p.LineTo(x0+halfWidth, y+dy)
		p.LineTo(x0, y0)
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
		p.MoveTo(x0, y0)
		p.LineTo(x1, y0+halfHeight)
		p.LineTo(x0, y0+height)
		p.LineTo(x0+dx, y0+halfHeight)
		p.LineTo(x0, y0)
	}
	p.FillStroke()
	return p
}

func (p *Painter) ArrowLeft(x, y, width, height int) *Painter {
	p.arrow(x, y, width, height, PositionLeft)
	return p
}

func (p *Painter) ArrowRight(x, y, width, height int) *Painter {
	p.arrow(x, y, width, height, PositionRight)
	return p
}

func (p *Painter) ArrowTop(x, y, width, height int) *Painter {
	p.arrow(x, y, width, height, PositionTop)
	return p
}
func (p *Painter) ArrowBottom(x, y, width, height int) *Painter {
	p.arrow(x, y, width, height, PositionBottom)
	return p
}

func (p *Painter) Circle(radius float64, x, y int) *Painter {
	p.render.Circle(radius, x+p.box.Left, y+p.box.Top)
	return p
}

func (p *Painter) Stroke() *Painter {
	p.render.Stroke()
	return p
}

func (p *Painter) Close() *Painter {
	p.render.Close()
	return p
}

func (p *Painter) FillStroke() *Painter {
	p.render.FillStroke()
	return p
}

func (p *Painter) Fill() *Painter {
	p.render.Fill()
	return p
}

// Width returns the drawable width of the painter's box.
func (p *Painter) Width() int {
	return p.box.Width()
}

// Height returns the drawable height of the painter's box.
func (p *Painter) Height() int {
	return p.box.Height()
}

// MeasureText will provide the rendered size of the text for the provided font style.
func (p *Painter) MeasureText(text string) Box {
	return p.render.MeasureText(text)
}

func (p *Painter) MeasureTextMaxWidthHeight(textList []string) (int, int) {
	maxWidth := 0
	maxHeight := 0
	for _, text := range textList {
		box := p.MeasureText(text)
		if maxWidth < box.Width() {
			maxWidth = box.Width()
		}
		if maxHeight < box.Height() {
			maxHeight = box.Height()
		}
	}
	return maxWidth, maxHeight
}

// LineStroke draws a line in the graph from point to point with the specified stroke color/width.
// Points with values of math.MaxInt32 will be skipped, resulting in a gap.
// Single or isolated points will result in just a dot being drawn at the point.
func (p *Painter) LineStroke(points []Point) *Painter {
	var valid []Point
	for _, pt := range points {
		if pt.Y == math.MaxInt32 {
			// If we encounter a break, draw the accumulated segment
			p.drawStraightPath(valid, true)
			p.Stroke()
			valid = valid[:0] // reset
			continue
		}
		valid = append(valid, pt)
	}

	// Draw the last segment if there is one
	p.drawStraightPath(valid, true)
	return p.Stroke()
}

// drawStraightPath draws a simple (non-curved) path for the given points.
// If dotForSinglePoint is true, single points are drawn as 2px radius dots.
func (p *Painter) drawStraightPath(points []Point, dotForSinglePoint bool) {
	pointCount := len(points)
	if pointCount == 0 {
		return
	} else if pointCount == 1 {
		if dotForSinglePoint {
			p.Dots(points)
		}
	}
	p.MoveTo(points[0].X, points[0].Y)
	for i := 1; i < pointCount; i++ {
		p.LineTo(points[i].X, points[i].Y)
	}
}

// smoothLineStroke draws a smooth curve through the given points using Quadratic Bézier segments and a
// `tension` parameter in [0..1] with 0 providing straight lines between midpoints and 1 providing a smoother line.
// Because the tension smooths out the line, the line will no longer hit the provided points exactly. The more variable
// the points, and the higher the tension, the more the line will be
func (p *Painter) smoothLineStroke(points []Point, tension float64) *Painter {
	if tension <= 0 {
		return p.LineStroke(points)
	} else if tension > 1 {
		tension = 1
	}

	var valid []Point // Slice to hold valid points between breaks
	for _, pt := range points {
		if pt.Y == math.MaxInt32 {
			// When a line break is found, draw the curve for the accumulated valid points if any
			p.drawSmoothCurve(valid, tension, true)
			p.Stroke()
			valid = valid[:0] // reset
			continue
		}

		valid = append(valid, pt)
	}
	// draw any remaining points collected
	p.drawSmoothCurve(valid, tension, true)
	return p.Stroke()
}

// drawSmoothCurve handles the actual path drawing (MoveTo/LineTo/QuadCurveTo)
// but does NOT call Stroke() or Fill(). This allows us to reuse this path
// logic for either smoothLineStroke or smoothFillArea.
func (p *Painter) drawSmoothCurve(points []Point, tension float64, dotForSinglePoint bool) {
	if len(points) < 3 { // Not enough points to form a curve, draw a line
		p.drawStraightPath(points, dotForSinglePoint)
		return
	}

	p.MoveTo(points[0].X, points[0].Y) // Start from the first valid point

	// Handle each segment between points with quadratic Bézier curves
	for i := 1; i < len(points)-1; i++ {
		x1, y1 := points[i].X, points[i].Y
		x2, y2 := points[i+1].X, points[i+1].Y

		mx := float64(x1+x2) / 2.0
		my := float64(y1+y2) / 2.0

		cx := float64(x1) + tension*(mx-float64(x1))
		cy := float64(y1) + tension*(my-float64(y1))

		p.QuadCurveTo(x1, y1, int(cx), int(cy))
	}

	// Connect the second-to-last point to the last point
	n := len(points)
	p.QuadCurveTo(points[n-2].X, points[n-2].Y, points[n-1].X, points[n-1].Y)
}

// SmoothLineStroke is Deprecated. This implementation produced sharp joints at the point, and will be removed in v0.4.0.
func (p *Painter) SmoothLineStroke(points []Point) *Painter {
	prevX := 0
	prevY := 0
	for index, point := range points {
		x := point.X
		y := point.Y
		if index == 0 {
			p.MoveTo(x, y)
		} else {
			cx := prevX + (x-prevX)/5
			cy := y + (y-prevY)/2
			p.QuadCurveTo(cx, cy, x, y)
		}
		prevX = x
		prevY = y
	}
	p.Stroke()
	return p
}

func (p *Painter) SetBackground(width, height int, color Color, inside ...bool) *Painter {
	r := p.render
	s := chartdraw.Style{
		FillColor: color,
	}
	// background color
	p.SetDrawingStyle(s)
	defer p.ResetStyle()
	if len(inside) != 0 && inside[0] {
		p.MoveTo(0, 0)
		p.LineTo(width, 0)
		p.LineTo(width, height)
		p.LineTo(0, height)
		p.LineTo(0, 0)
	} else {
		// setting the background color does not use boxes
		r.MoveTo(0, 0)
		r.LineTo(width, 0)
		r.LineTo(width, height)
		r.LineTo(0, height)
		r.LineTo(0, 0)
	}
	p.FillStroke()
	return p
}

// MarkLine draws a horizontal line with a small circle and arrow at the right.
func (p *Painter) MarkLine(x, y, width int) *Painter {
	arrowWidth := 16
	arrowHeight := 10
	endX := x + width
	radius := 3
	p.Circle(3, x+radius, y)
	p.render.Fill()
	p.MoveTo(x+radius*3, y)
	p.LineTo(endX-arrowWidth, y)
	p.Stroke()
	p.ArrowRight(endX, y, arrowWidth, arrowHeight)
	return p
}

// Polygon draws a polygon with the specified center, radius, and number of sides.
func (p *Painter) Polygon(center Point, radius float64, sides int) *Painter {
	points := getPolygonPoints(center, radius, sides)
	for i, item := range points {
		if i == 0 {
			p.MoveTo(item.X, item.Y)
		} else {
			p.LineTo(item.X, item.Y)
		}
	}
	p.LineTo(points[0].X, points[0].Y)
	p.Stroke()
	return p
}

// FillArea draws a filled polygon through the given points, skipping "null" (MaxInt32) break values (filling the area
// flat between them).
func (p *Painter) FillArea(points []Point) *Painter {
	if len(points) == 0 {
		return p
	}

	var valid []Point
	for _, pt := range points {
		if pt.Y == math.MaxInt32 {
			// If we encounter a break, fill the accumulated segment
			p.drawStraightPath(valid, false)
			p.Fill()
			valid = valid[:0] // reset
			continue
		}
		valid = append(valid, pt)
	}

	// Fill the last segment if there is one
	p.drawStraightPath(valid, false)
	p.Fill()

	return p
}

// smoothFillArea draws a smooth curve for the "top" portion of points but uses straight lines for the bottom corners,
// producing a fill with sharp corners.
func (p *Painter) smoothFillChartArea(points []Point, tension float64) *Painter {
	if tension <= 0 {
		return p.FillArea(points)
	} else if tension > 1 {
		tension = 1
	}

	// Typically, areaPoints has the shape:
	//   [ top data points... ] + [ bottom-right corner, bottom-left corner, first top point ]
	// We'll separate them:
	if len(points) < 3 {
		// Not enough to separate top from bottom
		return p.FillArea(points)
	}

	// The final 3 points are the corners + repeated first point
	top := points[:len(points)-3]
	bottom := points[len(points)-3:] // [ corner1, corner2, firstTopAgain ]

	// If top portion is empty or 1 point, just fill straight
	if len(top) < 2 {
		return p.FillArea(points)
	}

	// Build the smooth path for the top portion
	var currentSegment []Point
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
		return p.FillArea(points) // No actual top segment was drawn, fallback to straight fill
	}

	// Add sharp lines to close the shape at the bottom
	// The path is currently at the last top point we drew. Now we need to draw to corner1 -> corner2 -> firstTopAgain
	for i := 0; i < len(bottom); i++ {
		p.LineTo(bottom[i].X, bottom[i].Y)
	}

	p.Fill()
	return p
}

func (p *Painter) Text(body string, x, y int) *Painter {
	p.render.Text(body, x+p.box.Left, y+p.box.Top)
	return p
}

func (p *Painter) TextRotation(body string, x, y int, radians float64) {
	p.render.SetTextRotation(radians)
	p.render.Text(body, x+p.box.Left, y+p.box.Top)
	p.render.ClearTextRotation()
}

func (p *Painter) SetTextRotation(radians float64) {
	p.render.SetTextRotation(radians)
}
func (p *Painter) ClearTextRotation() {
	p.render.ClearTextRotation()
}

func (p *Painter) TextFit(body string, x, y, width int, textAligns ...string) chartdraw.Box {
	style := p.style
	textWarp := style.TextWrap
	style.TextWrap = chartdraw.TextWrapWord
	r := p.render
	lines := chartdraw.Text.WrapFit(r, body, width, style)
	p.SetFontStyle(style.FontStyle)
	var output chartdraw.Box

	textAlign := ""
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
		p.Text(line, x0, y0)
		output.Right = chartdraw.MaxInt(lineBox.Right, output.Right)
		output.Bottom += lineBox.Height()
		if index < len(lines)-1 {
			output.Bottom += +style.GetTextLineSpacing()
		}
	}
	p.style.TextWrap = textWarp
	return output
}

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
	// the actual position of this tick on the discrete scale of data indices, and rounds it
	// to ensure it aligns with an exact index in the array.
	actualTickIndex := int(float64(predictedTickIndex)*step + 0.5)
	return actualTickIndex == index
}

func (p *Painter) Ticks(opt TicksOption) *Painter {
	if opt.LabelCount <= 0 || opt.Length <= 0 {
		return p
	}
	var values []int
	if opt.Vertical {
		values = autoDivide(p.Height(), opt.TickSpaces)
	} else {
		values = autoDivide(p.Width(), opt.TickSpaces)
	}
	for index, value := range values {
		if index < opt.First {
			continue
		} else if !isTick(len(values)-opt.First, opt.TickCount, index-opt.First) {
			continue
		}
		if opt.Vertical {
			p.LineStroke([]Point{
				{X: 0, Y: value},
				{X: opt.Length, Y: value},
			})
		} else {
			p.LineStroke([]Point{
				{X: value, Y: opt.Length},
				{X: value, Y: 0},
			})
		}
	}
	return p
}

func (p *Painter) MultiText(opt MultiTextOption) *Painter {
	if len(opt.TextList) == 0 {
		return p
	}
	count := len(opt.TextList)
	width := p.Width()
	height := p.Height()
	var positions []int
	if opt.Vertical {
		if opt.CenterLabels {
			positions = autoDivide(height, count)
		} else {
			positions = autoDivide(height, count-1)
		}
	} else {
		if opt.CenterLabels {
			positions = autoDivide(width, count)
		} else {
			positions = autoDivide(width, count-1)
		}
	}
	isTextRotation := opt.TextRotation != 0
	positionCount := len(positions)

	skippedLabels := opt.LabelSkipCount // specify the skip count to ensure the top value is listed
	for index, start := range positions {
		if opt.CenterLabels && index == positionCount-1 {
			break // positions have one item more than we can map to text, this extra value is used to center against
		} else if index < opt.First {
			continue
		} else if !opt.Vertical &&
			index != count-1 && // one off case for last label due to values and label qty difference
			!isTick(positionCount-opt.First, opt.TickCount, index-opt.First) {
			continue
		} else if index != count-1 && // ensure the bottom value is always printed
			skippedLabels < opt.LabelSkipCount {
			skippedLabels++
			continue
		} else {
			skippedLabels = 0
		}

		if isTextRotation {
			p.ClearTextRotation()
			p.SetTextRotation(opt.TextRotation)
		}
		text := opt.TextList[index]
		box := p.MeasureText(text)
		x := 0
		y := 0
		if opt.Vertical {
			if opt.CenterLabels {
				start = (positions[index] + positions[index+1]) >> 1
			} else {
				start = positions[index]
			}
			y = start + box.Height()>>1
			switch opt.Align {
			case AlignRight:
				x = width - box.Width()
			case AlignCenter:
				x = width - box.Width()>>1
			default:
				x = 0
			}
		} else {
			if opt.CenterLabels {
				// graphs with limited data samples generally look better with the samples directly below the label
				// for that reason we will exactly center these graphs, but graphs with higher sample counts will
				// attempt to space the labels better rather than line up directly to the graph points
				exactLabels := count == opt.LabelCount
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
		x += opt.Offset.Left
		y += opt.Offset.Top
		p.Text(text, x, y)
	}
	if isTextRotation {
		p.ClearTextRotation()
	}
	return p
}

func (p *Painter) Grid(opt GridOption) *Painter {
	width := p.Width()
	height := p.Height()
	drawLines := func(values []int, ignoreIndexList []int, isVertical bool) {
		for index, v := range values {
			if containsInt(ignoreIndexList, index) {
				continue
			}
			x0 := 0
			y0 := 0
			x1 := 0
			y1 := 0
			if isVertical {
				x0 = v
				x1 = v
				y1 = height
			} else {
				x1 = width
				y0 = v
				y1 = v
			}
			p.LineStroke([]Point{
				{X: x0, Y: y0},
				{X: x1, Y: y1},
			})
		}
	}
	columnCount := sumInt(opt.ColumnSpans)
	if columnCount == 0 {
		columnCount = opt.Columns
	}
	if columnCount > 0 {
		values := autoDivideSpans(width, columnCount, opt.ColumnSpans)
		drawLines(values, opt.IgnoreColumnLines, true)
	}
	if opt.Rows > 0 {
		values := autoDivide(height, opt.Rows)
		drawLines(values, opt.IgnoreRowLines, false)
	}
	return p
}

func (p *Painter) Dots(points []Point) *Painter {
	for _, item := range points {
		p.Circle(2, item.X, item.Y)
	}
	p.FillStroke()
	return p
}

func (p *Painter) Rect(box Box) *Painter {
	p.MoveTo(box.Left, box.Top)
	p.LineTo(box.Right, box.Top)
	p.LineTo(box.Right, box.Bottom)
	p.LineTo(box.Left, box.Bottom)
	p.LineTo(box.Left, box.Top)
	p.FillStroke()
	return p
}

func (p *Painter) RoundedRect(box Box, radius int, roundTop, roundBottom bool) *Painter {
	r := (box.Right - box.Left) / 2
	if radius > r {
		radius = r
	}
	rx := float64(radius)
	ry := float64(radius)

	if roundTop {
		// Start at the appropriate point depending on rounding at the top
		p.MoveTo(box.Left+radius, box.Top)
		p.LineTo(box.Right-radius, box.Top)

		// right top
		cx := box.Right - radius
		cy := box.Top + radius
		p.ArcTo(cx, cy, rx, ry, -math.Pi/2, math.Pi/2)
	} else {
		p.MoveTo(box.Left, box.Top)
		p.LineTo(box.Right, box.Top)
	}

	if roundBottom {
		p.LineTo(box.Right, box.Bottom-radius)

		// right bottom
		cx := box.Right - radius
		cy := box.Bottom - radius
		p.ArcTo(cx, cy, rx, ry, 0, math.Pi/2)

		p.LineTo(box.Left+radius, box.Bottom)

		// left bottom
		cx = box.Left + radius
		cy = box.Bottom - radius
		p.ArcTo(cx, cy, rx, ry, math.Pi/2, math.Pi/2)
	} else {
		p.LineTo(box.Right, box.Bottom)
		p.LineTo(box.Left, box.Bottom)
	}

	if roundTop {
		// left top
		p.LineTo(box.Left, box.Top+radius)
		cx := box.Left + radius
		cy := box.Top + radius
		p.ArcTo(cx, cy, rx, ry, math.Pi, math.Pi/2)
	} else {
		p.LineTo(box.Left, box.Top)
	}

	p.Close()
	p.FillStroke()
	p.Fill()
	return p
}

func (p *Painter) LegendLineDot(box Box) *Painter {
	width := box.Width()
	height := box.Height()
	strokeWidth := 3
	dotHeight := 5

	p.render.SetStrokeWidth(float64(strokeWidth))
	center := (height-strokeWidth)>>1 - 1
	p.MoveTo(box.Left, box.Top-center)
	p.LineTo(box.Right, box.Top-center)
	p.Stroke()
	p.Circle(float64(dotHeight), box.Left+width>>1, box.Top-center)
	p.FillStroke()
	return p
}
