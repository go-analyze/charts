package drawing

// DemuxFlattener is a slice of Flattener.
type DemuxFlattener struct {
	Flatteners []Flattener
}

// MoveTo forwards a move command to all child flatteners (for PathBuilder interface).
func (dc DemuxFlattener) MoveTo(x, y float64) {
	for _, flattener := range dc.Flatteners {
		flattener.MoveTo(x, y)
	}
}

// LineTo forwards a line command to all child flatteners (for PathBuilder interface).
func (dc DemuxFlattener) LineTo(x, y float64) {
	for _, flattener := range dc.Flatteners {
		flattener.LineTo(x, y)
	}
}

// End signals completion to all child flatteners (for PathBuilder interface).
func (dc DemuxFlattener) End() {
	for _, flattener := range dc.Flatteners {
		flattener.End()
	}
}
