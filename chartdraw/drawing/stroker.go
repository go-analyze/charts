// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 13/12/2010 by Laurent Le Goff

package drawing

// NewLineStroker creates a new line stroker.
func NewLineStroker(flattener Flattener) *LineStroker {
	l := new(LineStroker)
	l.Flattener = flattener
	l.HalfLineWidth = 0.5
	return l
}

// LineStroker draws the stroke portion of a line.
type LineStroker struct {
	Flattener     Flattener
	HalfLineWidth float64
	vertices      []float64
	rewind        []float64
	x, y, nx, ny  float64
}

// MoveTo records the starting point for a stroked segment (for PathBuilder interface).
func (l *LineStroker) MoveTo(x, y float64) {
	l.x, l.y = x, y
}

// LineTo adds a stroked line to the path (for PathBuilder interface).
func (l *LineStroker) LineTo(x, y float64) {
	l.line(l.x, l.y, x, y)
}

func (l *LineStroker) line(x1, y1, x2, y2 float64) {
	dx := x2 - x1
	dy := y2 - y1
	d := vectorDistance(dx, dy)
	if d != 0 {
		nx := dy * l.HalfLineWidth / d
		ny := -(dx * l.HalfLineWidth / d)
		l.appendVertex(x1+nx, y1+ny, x2+nx, y2+ny, x1-nx, y1-ny, x2-nx, y2-ny)
		l.x, l.y, l.nx, l.ny = x2, y2, nx, ny
	}
}

// End emits the completed stroked path to the next flattener (for PathBuilder interface).
func (l *LineStroker) End() {
	if len(l.vertices) > 1 {
		l.Flattener.MoveTo(l.vertices[0], l.vertices[1])
		for i, j := 2, 3; j < len(l.vertices); i, j = i+2, j+2 {
			l.Flattener.LineTo(l.vertices[i], l.vertices[j])
		}
	}
	for i, j := len(l.rewind)-2, len(l.rewind)-1; j > 0; i, j = i-2, j-2 {
		l.Flattener.LineTo(l.rewind[i], l.rewind[j])
	}
	if len(l.vertices) > 1 {
		l.Flattener.LineTo(l.vertices[0], l.vertices[1])
	}
	l.Flattener.End()
	// reinit vertices
	l.vertices = l.vertices[0:0]
	l.rewind = l.rewind[0:0]
	l.x, l.y, l.nx, l.ny = 0, 0, 0, 0
}

func (l *LineStroker) appendVertex(vertices ...float64) {
	s := len(vertices) / 2
	l.vertices = append(l.vertices, vertices[:s]...)
	l.rewind = append(l.rewind, vertices[s:]...)
}
