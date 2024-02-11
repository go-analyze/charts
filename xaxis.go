package charts

import (
	"github.com/golang/freetype/truetype"
)

type XAxisOption struct {
	// Show specifies if the x-axis should be rendered, set this to *false (through False()) to hide the axis.
	Show *bool
	// Theme specifies the colors used for the x-axis.
	Theme ColorPalette
	// Data provides labels for the x-axis.
	Data []string
	// DataStartIndex specifies what index the Data values should start from.
	DataStartIndex int
	// Position describes the position of x-axis, it can be 'top' or 'bottom'.
	Position string
	// BoundaryGap specifies that the chart should have additional space on the left and right, with data points being
	// centered between two axis ticks.  Enabled by default, specify *false (through False()) to change the spacing.
	BoundaryGap *bool
	// FontSize specifies the font size of each label.
	FontSize float64
	// Font is the font used to render each label.
	Font *truetype.Font
	// FontColor is the color used for text rendered.
	FontColor Color
	// TextRotation are the radians for rotating the label.
	TextRotation float64
	// LabelOffset is the offset of each label.
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
		DataStartIndex: opt.DataStartIndex,
		BoundaryGap:    opt.BoundaryGap,
		Position:       position,
		FontSize:       opt.FontSize,
		Font:           opt.Font,
		FontColor:      opt.FontColor,
		Show:           opt.Show,
		Unit:           opt.Unit,
		LabelCount:     opt.LabelCount,
		TextRotation:   opt.TextRotation,
		LabelOffset:    opt.LabelOffset,
	}
	if opt.isValueAxis {
		axisOpt.SplitLineShow = true
		axisOpt.StrokeWidth = -1
		axisOpt.BoundaryGap = False()
	}
	return axisOpt
}

// NewBottomXAxis returns a bottom x axis renderer
func NewBottomXAxis(p *Painter, opt XAxisOption) *axisPainter {
	return NewAxisPainter(p, opt.ToAxisOption())
}
