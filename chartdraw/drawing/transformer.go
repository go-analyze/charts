// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 13/12/2010 by Laurent Le Goff

package drawing

// Transformer apply the Matrix transformation tr.
type Transformer struct {
	Tr        Matrix
	Flattener Flattener
}

// MoveTo transforms and passes a move command downstream (for PathBuilder interface).
func (t Transformer) MoveTo(x, y float64) {
	u := x*t.Tr[0] + y*t.Tr[2] + t.Tr[4]
	v := x*t.Tr[1] + y*t.Tr[3] + t.Tr[5]
	t.Flattener.MoveTo(u, v)
}

// LineTo transforms and forwards a line command (for PathBuilder interface).
func (t Transformer) LineTo(x, y float64) {
	u := x*t.Tr[0] + y*t.Tr[2] + t.Tr[4]
	v := x*t.Tr[1] + y*t.Tr[3] + t.Tr[5]
	t.Flattener.LineTo(u, v)
}

// End finishes the transformed path (for PathBuilder interface).
func (t Transformer) End() {
	t.Flattener.End()
}
