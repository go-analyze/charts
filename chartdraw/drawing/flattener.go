package drawing

// Liner receive segment definition.
type Liner interface {
	// LineTo Draw a line from the current position to the point (x, y).
	LineTo(x, y float64)
}

// Flattener receive segment definition.
type Flattener interface {
	// MoveTo Start a New line from the point (x, y).
	MoveTo(x, y float64)
	// LineTo Draw a line from the current position to the point (x, y).
	LineTo(x, y float64)
	// End mark the current line as finished.
	End()
}

// Flatten convert curves into straight segments keeping join segments info.
func Flatten(path *Path, flattener Flattener, scale float64) {
	var startX, startY float64 // moveTo point starting a path
	var x, y float64           // current point
	var i int
	for _, cmp := range path.Components {
		switch cmp {
		case MoveToComponent:
			x, y = path.Points[i], path.Points[i+1]
			startX, startY = x, y
			if i != 0 {
				flattener.End()
			}
			flattener.MoveTo(x, y)
			i += 2
		case LineToComponent:
			x, y = path.Points[i], path.Points[i+1]
			flattener.LineTo(x, y)
			i += 2
		case QuadCurveToComponent:
			// we include the previous point for the start of the curve
			TraceQuad(flattener, path.Points[i-2:], 0.5)
			x, y = path.Points[i+2], path.Points[i+3]
			flattener.LineTo(x, y)
			i += 4
		case CubicCurveToComponent:
			TraceCubic(flattener, path.Points[i-2:], 0.5)
			x, y = path.Points[i+4], path.Points[i+5]
			flattener.LineTo(x, y)
			i += 6
		case ArcToComponent:
			x, y = TraceArc(flattener, path.Points[i], path.Points[i+1], path.Points[i+2], path.Points[i+3], path.Points[i+4], path.Points[i+5], scale)
			flattener.LineTo(x, y)
			i += 6
		case CloseComponent:
			if x != startX || y != startY {
				flattener.LineTo(startX, startY)
			}
		}
	}
	flattener.End()
}

// SegmentedPath is a path of disparate point sections.
type SegmentedPath struct {
	Points []float64
}

// MoveTo records the first point of a new segment (for PathBuilder interface).
func (p *SegmentedPath) MoveTo(x, y float64) {
	p.Points = append(p.Points, x, y)
}

// LineTo appends a point to the current path segment (for PathBuilder interface).
func (p *SegmentedPath) LineTo(x, y float64) {
	p.Points = append(p.Points, x, y)
}

// End completes the current path segment (for PathBuilder interface).
func (p *SegmentedPath) End() {
	// Nothing to do
}
