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

type ticksOption struct {
	firstIndex int
	length     int
	vertical   bool
	labelCount int
	tickCount  int
	tickSpaces int
}

type multiTextOption struct {
	textList       []string
	vertical       bool
	centerLabels   bool
	align          string
	textRotation   float64
	offset         OffsetInt
	firstIndex     int
	labelCount     int
	tickCount      int
	labelSkipCount int
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

// NewPainter creates a painter which can be used to render charts to (using for example newLineChart).
func NewPainter(opts PainterOptions, opt ...PainterOptionFunc) *Painter {
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
	r := fn(opts.Width, opts.Height)
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
	return p
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

func (p *Painter) setStyle(style chartdraw.Style) {
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

func (p *Painter) OverrideDrawingStyle(style chartdraw.Style) {
	// TODO - we should alias parts of Style we want to support drawing on
	s := overrideStyle(p.style, style)
	p.SetDrawingStyle(s)
}

func (p *Painter) SetDrawingStyle(style chartdraw.Style) {
	style.WriteDrawingOptionsToRenderer(p.render)
}

func (p *Painter) SetFontStyle(style chartdraw.FontStyle) {
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
}

func (p *Painter) OverrideFontStyle(style chartdraw.FontStyle) {
	s := overrideStyle(p.style, chartdraw.Style{FontStyle: style})
	p.SetFontStyle(s.FontStyle)
}

func (p *Painter) resetStyle() {
	p.style.WriteToRenderer(p.render)
}

// Bytes returns the data of draw canvas.
func (p *Painter) Bytes() ([]byte, error) {
	buffer := bytes.Buffer{}
	if err := p.render.Save(&buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// moveTo moves the cursor to a given point.
func (p *Painter) moveTo(x, y int) {
	p.render.MoveTo(x+p.box.Left, y+p.box.Top)
}

// arcTo renders an arc from the current cursor and the given parameters.
func (p *Painter) arcTo(cx, cy int, rx, ry, startAngle, delta float64) {
	p.render.ArcTo(cx+p.box.Left, cy+p.box.Top, rx, ry, startAngle, delta)
}

// Line renders a line from the first xy coordinates to the second point.
func (p *Painter) Line(x1, y1, x2, y2 int) {
	p.moveTo(x1, y1)
	p.lineTo(x2, y2)
}

func (p *Painter) quadCurveTo(cx, cy, x, y int) {
	p.render.QuadCurveTo(cx+p.box.Left, cy+p.box.Top, x+p.box.Left, y+p.box.Top)
}

// lineTo renders a line from the current cursor to the given point.
func (p *Painter) lineTo(x, y int) {
	p.render.LineTo(x+p.box.Left, y+p.box.Top)
}

func (p *Painter) Pin(x, y, width int) {
	r := float64(width) / 2
	y -= width / 4
	angle := chartdraw.DegreesToRadians(15)
	box := p.box

	startAngle := math.Pi/2 + angle
	delta := 2*math.Pi - 2*angle
	p.arcTo(x, y, r, r, startAngle, delta)
	p.lineTo(x, y)
	p.close()
	p.fillStroke()

	startX := x - int(r)
	startY := y
	endX := x + int(r)
	endY := y
	p.moveTo(startX, startY)

	left := box.Left
	top := box.Top
	cx := x
	cy := y + int(r*2.5)
	p.render.QuadCurveTo(cx+left, cy+top, endX+left, endY+top)
	p.close()
	p.render.Fill()
}

func (p *Painter) arrow(x, y, width, height int, direction string) {
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
	p.fillStroke()
}

// ArrowLeft draws an arrow at the given point and dimensions pointing left.
func (p *Painter) ArrowLeft(x, y, width, height int) {
	p.arrow(x, y, width, height, PositionLeft)
}

// ArrowRight draws an arrow at the given point and dimensions pointing right.
func (p *Painter) ArrowRight(x, y, width, height int) {
	p.arrow(x, y, width, height, PositionRight)
}

// ArrowUp draws an arrow at the given point and dimensions pointing up.
func (p *Painter) ArrowUp(x, y, width, height int) {
	p.arrow(x, y, width, height, PositionTop)
}

// ArrowDown draws an arrow at the given point and dimensions pointing down.
func (p *Painter) ArrowDown(x, y, width, height int) {
	p.arrow(x, y, width, height, PositionBottom)
}

// Circle draws a circle at the given coords with a given radius.
func (p *Painter) Circle(radius float64, x, y int) {
	p.render.Circle(radius, x+p.box.Left, y+p.box.Top)
}

func (p *Painter) stroke() {
	p.render.Stroke()
}

func (p *Painter) close() {
	p.render.Close()
}

func (p *Painter) fillStroke() {
	p.render.FillStroke()
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

func (p *Painter) measureTextMaxWidthHeight(textList []string) (int, int) {
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
func (p *Painter) LineStroke(points []Point) {
	var valid []Point
	for _, pt := range points {
		if pt.Y == math.MaxInt32 {
			// If we encounter a break, draw the accumulated segment
			p.drawStraightPath(valid, true)
			p.stroke()
			valid = valid[:0] // reset
			continue
		}
		valid = append(valid, pt)
	}

	// Draw the last segment if there is one
	p.drawStraightPath(valid, true)
	p.stroke()
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
	p.moveTo(points[0].X, points[0].Y)
	for i := 1; i < pointCount; i++ {
		p.lineTo(points[i].X, points[i].Y)
	}
}

// SmoothLineStroke draws a smooth curve through the given points using Quadratic Bézier segments and a
// `tension` parameter in [0..1] with 0 providing straight lines between midpoints and 1 providing a smoother line.
// Because the tension smooths out the line, the line will no longer hit the provided points exactly. The more variable
// the points, and the higher the tension, the more the line will be
func (p *Painter) SmoothLineStroke(points []Point, tension float64) {
	if tension <= 0 {
		p.LineStroke(points)
		return
	} else if tension > 1 {
		tension = 1
	}

	var valid []Point // Slice to hold valid points between breaks
	for _, pt := range points {
		if pt.Y == math.MaxInt32 {
			// When a line break is found, draw the curve for the accumulated valid points if any
			p.drawSmoothCurve(valid, tension, true)
			p.stroke()
			valid = valid[:0] // reset
			continue
		}

		valid = append(valid, pt)
	}
	// draw any remaining points collected
	p.drawSmoothCurve(valid, tension, true)
	p.stroke()
}

// drawSmoothCurve handles the actual path drawing (MoveTo/LineTo/QuadCurveTo)
// but does NOT call Stroke() or Fill(). This allows us to reuse this path
// logic for either SmoothLineStroke or smoothFillArea.
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

func (p *Painter) SetBackground(width, height int, color Color, inside ...bool) {
	r := p.render
	s := chartdraw.Style{
		FillColor: color,
	}
	// background color
	p.SetDrawingStyle(s)
	defer p.resetStyle()
	if len(inside) != 0 && inside[0] {
		p.moveTo(0, 0)
		p.lineTo(width, 0)
		p.lineTo(width, height)
		p.lineTo(0, height)
		p.lineTo(0, 0)
	} else {
		// setting the background color does not use boxes
		r.MoveTo(0, 0)
		r.LineTo(width, 0)
		r.LineTo(width, height)
		r.LineTo(0, height)
		r.LineTo(0, 0)
	}
	p.fillStroke()
}

// MarkLine draws a horizontal line with a small circle and arrow at the right.
func (p *Painter) MarkLine(x, y, width int) {
	arrowWidth := 16
	arrowHeight := 10
	endX := x + width
	radius := 3
	p.Circle(3, x+radius, y)
	p.render.Fill()
	p.Line(x+radius*3, y, endX-arrowWidth, y)
	p.stroke()
	p.ArrowRight(endX, y, arrowWidth, arrowHeight)
}

// Polygon draws a polygon with the specified center, radius, and number of sides.
func (p *Painter) Polygon(center Point, radius float64, sides int) {
	points := getPolygonPoints(center, radius, sides)
	for i, item := range points {
		if i == 0 {
			p.moveTo(item.X, item.Y)
		} else {
			p.lineTo(item.X, item.Y)
		}
	}
	p.lineTo(points[0].X, points[0].Y)
	p.stroke()
}

// FillArea draws a filled polygon through the given points, skipping "null" (MaxInt32) break values (filling the area
// flat between them).
func (p *Painter) FillArea(points []Point) {
	if len(points) == 0 {
		return
	}

	var valid []Point
	for _, pt := range points {
		if pt.Y == math.MaxInt32 {
			// If we encounter a break, fill the accumulated segment
			p.drawStraightPath(valid, false)
			p.render.Fill()
			valid = valid[:0] // reset
			continue
		}
		valid = append(valid, pt)
	}

	// Fill the last segment if there is one
	p.drawStraightPath(valid, false)
	p.render.Fill()
}

// smoothFillArea draws a smooth curve for the "top" portion of points but uses straight lines for the bottom corners,
// producing a fill with sharp corners.
func (p *Painter) smoothFillChartArea(points []Point, tension float64) {
	if tension <= 0 {
		p.FillArea(points)
		return
	} else if tension > 1 {
		tension = 1
	}

	// Typically, areaPoints has the shape:
	//   [ top data points... ] + [ bottom-right corner, bottom-left corner, first top point ]
	// We'll separate them:
	if len(points) < 3 {
		// Not enough to separate top from bottom
		p.FillArea(points)
		return
	}

	// The final 3 points are the corners + repeated first point
	top := points[:len(points)-3]
	bottom := points[len(points)-3:] // [ corner1, corner2, firstTopAgain ]

	// If top portion is empty or 1 point, just fill straight
	if len(top) < 2 {
		p.FillArea(points)
		return
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
		p.FillArea(points) // No actual top segment was drawn, fallback to straight fill
		return
	}

	// Add sharp lines to close the shape at the bottom
	// The path is currently at the last top point we drew. Now we need to draw to corner1 -> corner2 -> firstTopAgain
	for i := 0; i < len(bottom); i++ {
		p.lineTo(bottom[i].X, bottom[i].Y)
	}

	p.render.Fill()
}

func (p *Painter) Text(body string, x, y int) {
	p.render.Text(body, x+p.box.Left, y+p.box.Top)
}

func (p *Painter) TextRotation(body string, x, y int, radians float64) {
	p.render.SetTextRotation(radians)
	p.render.Text(body, x+p.box.Left, y+p.box.Top)
	p.render.ClearTextRotation()
}

func (p *Painter) setTextRotation(radians float64) {
	p.render.SetTextRotation(radians)
}
func (p *Painter) clearTextRotation() {
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

func (p *Painter) ticks(opt ticksOption) {
	if opt.labelCount <= 0 || opt.length <= 0 {
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
			})
		} else {
			p.LineStroke([]Point{
				{X: value, Y: opt.length},
				{X: value, Y: 0},
			})
		}
	}
}

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
	isTextRotation := opt.textRotation != 0
	positionCount := len(positions)

	skippedLabels := opt.labelSkipCount // specify the skip count to ensure the top value is listed
	for index, start := range positions {
		if opt.centerLabels && index == positionCount-1 {
			break // positions have one item more than we can map to text, this extra value is used to center against
		} else if index < opt.firstIndex {
			continue
		} else if !opt.vertical &&
			index != count-1 && // one off case for last label due to values and label qty difference
			!isTick(positionCount-opt.firstIndex, opt.tickCount, index-opt.firstIndex) {
			continue
		} else if index != count-1 && // ensure the bottom value is always printed
			skippedLabels < opt.labelSkipCount {
			skippedLabels++
			continue
		} else {
			skippedLabels = 0
		}

		if isTextRotation {
			p.clearTextRotation()
			p.setTextRotation(opt.textRotation)
		}
		text := opt.textList[index]
		box := p.MeasureText(text)
		x := 0
		y := 0
		if opt.vertical {
			if opt.centerLabels {
				start = (positions[index] + positions[index+1]) >> 1
			} else {
				start = positions[index]
			}
			y = start + box.Height()>>1
			switch opt.align {
			case AlignRight:
				x = width - box.Width()
			case AlignCenter:
				x = width - box.Width()>>1
			default:
				x = 0
			}
		} else {
			if opt.centerLabels {
				// graphs with limited data samples generally look better with the samples directly below the label
				// for that reason we will exactly center these graphs, but graphs with higher sample counts will
				// attempt to space the labels better rather than line up directly to the graph points
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
		p.Text(text, x, y)
	}
	if isTextRotation {
		p.clearTextRotation()
	}
}

// Dots prints filled circles for the given points.
func (p *Painter) Dots(points []Point) {
	for _, item := range points {
		p.Circle(2, item.X, item.Y)
	}
	p.fillStroke()
}

// Rect will draw a box with the given coordinates.
func (p *Painter) Rect(box Box) {
	p.moveTo(box.Left, box.Top)
	p.lineTo(box.Right, box.Top)
	p.lineTo(box.Right, box.Bottom)
	p.lineTo(box.Left, box.Bottom)
	p.lineTo(box.Left, box.Top)
}

// filledRect will draw a filled box with the given coordinates.
func (p *Painter) filledRect(box Box) {
	p.Rect(box)
	p.fillStroke()
}

// roundedRect is similar to filledRect except the top and bottom will be rounded.
func (p *Painter) roundedRect(box Box, radius int, roundTop, roundBottom bool) {
	r := (box.Right - box.Left) / 2
	if radius > r {
		radius = r
	}
	rx := float64(radius)
	ry := float64(radius)

	if roundTop {
		// Start at the appropriate point depending on rounding at the top
		p.Line(box.Left+radius, box.Top, box.Right-radius, box.Top)

		// right top
		cx := box.Right - radius
		cy := box.Top + radius
		p.arcTo(cx, cy, rx, ry, -math.Pi/2, math.Pi/2)
	} else {
		p.Line(box.Left, box.Top, box.Right, box.Top)
	}

	if roundBottom {
		p.lineTo(box.Right, box.Bottom-radius)

		// right bottom
		cx := box.Right - radius
		cy := box.Bottom - radius
		p.arcTo(cx, cy, rx, ry, 0, math.Pi/2)

		p.lineTo(box.Left+radius, box.Bottom)

		// left bottom
		cx = box.Left + radius
		cy = box.Bottom - radius
		p.arcTo(cx, cy, rx, ry, math.Pi/2, math.Pi/2)
	} else {
		p.lineTo(box.Right, box.Bottom)
		p.lineTo(box.Left, box.Bottom)
	}

	if roundTop {
		// left top
		p.lineTo(box.Left, box.Top+radius)
		cx := box.Left + radius
		cy := box.Top + radius
		p.arcTo(cx, cy, rx, ry, math.Pi, math.Pi/2)
	} else {
		p.lineTo(box.Left, box.Top)
	}

	p.close()
	p.fillStroke()
	p.render.Fill()
}

func (p *Painter) legendLineDot(box Box) {
	width := box.Width()
	height := box.Height()
	strokeWidth := 3
	dotHeight := 5

	p.render.SetStrokeWidth(float64(strokeWidth))
	center := (height-strokeWidth)>>1 - 1
	p.Line(box.Left, box.Top-center, box.Right, box.Top-center)
	p.stroke()
	p.Circle(float64(dotHeight), box.Left+width>>1, box.Top-center)
	p.fillStroke()
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

// PieChart renders a pie chart with the provided configuration to the painter.
func (p *Painter) PieChart(opt PieChartOption) error {
	_, err := newPieChart(p, opt).Render()
	return err
}

// RadarChart renders a radar chart with the provided configuration to the painter.
func (p *Painter) RadarChart(opt RadarChartOption) error {
	_, err := newRadarChart(p, opt).Render()
	return err
}

// TableChart renders a table with the provided configuration to the painter.
func (p *Painter) TableChart(opt TableChartOption) error {
	_, err := newTableChart(p, opt).Render()
	return err
}
