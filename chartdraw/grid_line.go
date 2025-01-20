package chartdraw

// GridLineProvider is a type that provides grid lines.
type GridLineProvider interface {
	GetGridLines(ticks []Tick, isVertical bool, majorStyle, minorStyle Style) []GridLine
}

// GridLine is a line on a graph canvas.
type GridLine struct {
	IsMinor bool
	Style   Style
	Value   float64
}

// Major returns if the gridline is a `major` line.
func (gl GridLine) Major() bool {
	return !gl.IsMinor
}

// Minor returns if the gridline is a `minor` line.
func (gl GridLine) Minor() bool {
	return gl.IsMinor
}

// Render renders the gridline
func (gl GridLine) Render(r Renderer, canvasBox Box, ra Range, isVertical bool, defaults Style) {
	r.SetStrokeColor(gl.Style.GetStrokeColor(defaults.GetStrokeColor()))
	r.SetStrokeWidth(gl.Style.GetStrokeWidth(defaults.GetStrokeWidth()))
	r.SetStrokeDashArray(gl.Style.GetStrokeDashArray(defaults.GetStrokeDashArray()))

	if isVertical {
		lineLeft := canvasBox.Left + ra.Translate(gl.Value)
		lineBottom := canvasBox.Bottom
		lineTop := canvasBox.Top

		r.MoveTo(lineLeft, lineBottom)
		r.LineTo(lineLeft, lineTop)
		r.Stroke()
	} else {
		lineLeft := canvasBox.Left
		lineRight := canvasBox.Right
		lineHeight := canvasBox.Bottom - ra.Translate(gl.Value)

		r.MoveTo(lineLeft, lineHeight)
		r.LineTo(lineRight, lineHeight)
		r.Stroke()
	}
}

// GenerateGridLines generates grid lines.
func GenerateGridLines(ticks []Tick, majorStyle, minorStyle Style) []GridLine {
	if len(ticks) < 3 {
		return []GridLine{}
	}

	isMinor := false
	gl := make([]GridLine, 0, len(ticks)-2)
	for _, t := range ticks[1 : len(ticks)-1] {
		s := majorStyle
		if isMinor {
			s = minorStyle
		}
		gl = append(gl, GridLine{
			Style:   s,
			IsMinor: isMinor,
			Value:   t.Value,
		})
		isMinor = !isMinor
	}
	return gl
}
