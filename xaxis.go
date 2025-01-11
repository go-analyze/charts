package charts

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
	// centered between two axis ticks. Default is set based on the dataset density / size to produce an easy-to-read
	// graph. Specify a *bool (through charts.False() or charts.True()) to enforce a spacing.
	BoundaryGap *bool
	// FontStyle specifies the font configuration for each label.
	FontStyle FontStyle
	// TextRotation are the radians for rotating the label.
	TextRotation float64
	// LabelOffset is the offset of each label.
	LabelOffset OffsetInt
	// ValueFormatter defines how float values should be rendered to strings, notably for numeric axis labels.
	ValueFormatter ValueFormatter
	// Unit is a suggestion for how large the axis step is, this is a recommendation only. Larger numbers result in fewer labels.
	Unit float64
	// LabelCount is the number of labels to show on the axis.  Specify a smaller number to reduce writing collisions.
	LabelCount int
	// LabelCountAdjustment specifies a relative influence on how many labels should be rendered.
	// Typically, this is negative to result in cleaner graphs, positive values may result in text collisions.
	LabelCountAdjustment int
	isValueAxis          bool
}

const defaultXAxisHeight = 30
const boundaryGapDefaultThreshold = 40

func (opt *XAxisOption) toAxisOption() axisOption {
	position := PositionBottom
	if opt.Position == PositionTop {
		position = PositionTop
	}
	axisOpt := axisOption{
		Theme:                opt.Theme,
		Data:                 opt.Data,
		DataStartIndex:       opt.DataStartIndex,
		BoundaryGap:          opt.BoundaryGap,
		Position:             position,
		FontStyle:            opt.FontStyle,
		Show:                 opt.Show,
		Unit:                 opt.Unit,
		LabelCount:           opt.LabelCount,
		LabelCountAdjustment: opt.LabelCountAdjustment,
		TextRotation:         opt.TextRotation,
		LabelOffset:          opt.LabelOffset,
	}
	if opt.isValueAxis {
		axisOpt.SplitLineShow = true
		axisOpt.StrokeWidth = -1
		axisOpt.BoundaryGap = False()
	}
	return axisOpt
}

// newBottomXAxis returns a bottom x-axis renderer.
func newBottomXAxis(p *Painter, opt XAxisOption) *axisPainter {
	return newAxisPainter(p, opt.toAxisOption())
}
