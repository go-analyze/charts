package charts

import (
	"github.com/golang/freetype/truetype"
)

type XAxisOption struct {
	// The font of x-axis
	Font *truetype.Font
	// The boundary gap on both sides of a coordinate axis.
	// Nil or *true means the center part of two axis ticks
	BoundaryGap *bool
	// The data value of x-axis
	Data []string
	// The theme of chart
	Theme ColorPalette
	// The font size of x-axis label
	FontSize float64
	// The flag for show axis, set this to *false will hide axis
	Show *bool
	// The position of axis, it can be 'top' or 'bottom'
	Position string
	// The line color of axis
	StrokeColor Color
	// The color of label
	FontColor Color
	// The text rotation of label
	TextRotation float64
	// The first axis
	FirstAxis int
	// The offset of label
	LabelOffset Box
	// Unit is a suggestion for how large the axis step is, this is a recommendation only. Larger numbers result in fewer labels.
	Unit float64
	// LabelCount is the number of labels to show on the axis.  Specify a smaller number to reduce writing collisions.
	LabelCount  int
	isValueAxis bool
}

const defaultXAxisHeight = 30

// NewXAxisOption returns a x axis option
func NewXAxisOption(data []string, boundaryGap ...*bool) XAxisOption {
	opt := XAxisOption{
		Data: data,
	}
	if len(boundaryGap) != 0 {
		opt.BoundaryGap = boundaryGap[0]
	}
	return opt
}

func (opt *XAxisOption) ToAxisOption() AxisOption {
	position := PositionBottom
	if opt.Position == PositionTop {
		position = PositionTop
	}
	axisOpt := AxisOption{
		Theme:          opt.Theme,
		Data:           opt.Data,
		BoundaryGap:    opt.BoundaryGap,
		Position:       position,
		StrokeColor:    opt.StrokeColor,
		FontSize:       opt.FontSize,
		Font:           opt.Font,
		FontColor:      opt.FontColor,
		Show:           opt.Show,
		Unit:           opt.Unit,
		LabelCount:     opt.LabelCount,
		SplitLineColor: opt.Theme.GetAxisSplitLineColor(),
		TextRotation:   opt.TextRotation,
		LabelOffset:    opt.LabelOffset,
		FirstAxis:      opt.FirstAxis,
	}
	if opt.isValueAxis {
		axisOpt.SplitLineShow = true
		axisOpt.StrokeWidth = -1
		axisOpt.BoundaryGap = FalseFlag()
	}
	return axisOpt
}

// NewBottomXAxis returns a bottom x axis renderer
func NewBottomXAxis(p *Painter, opt XAxisOption) *axisPainter {
	return NewAxisPainter(p, opt.ToAxisOption())
}
