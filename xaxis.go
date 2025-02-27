package charts

type XAxisOption struct {
	// Show specifies if the x-axis should be rendered, set this to *false (through Ptr(false)) to hide the axis.
	Show *bool
	// Theme specifies the colors used for the x-axis.
	Theme ColorPalette
	// Title specifies a name for the axis, if specified the axis name is rendered below the X-Axis.
	Title string
	// TitleFontStyle provides the font, size, and color for the axis title.
	TitleFontStyle FontStyle
	// Labels provides labels for each value on the x-axis (index matching to the series index).
	Labels []string
	// DataStartIndex specifies what index the Data values should start from.
	DataStartIndex int
	// Position describes the position of x-axis, it can be 'top' or 'bottom'.
	Position string
	// BoundaryGap specifies that the chart should have additional space on the left and right, with data points being
	// centered between two axis ticks. Default is set based on the dataset density / size to produce an easy-to-read
	// graph. Specify a *bool (through charts.Ptr(false) or charts.Ptr(true)) to enforce a spacing.
	BoundaryGap *bool
	// Deprecated: FontStyle is deprecated, use LabelFontStyle.
	FontStyle FontStyle
	// LabelFontStyle specifies the font configuration for each label.
	LabelFontStyle FontStyle
	// LabelRotation are the radians for rotating the label. Convert from degrees using DegreesToRadians(float64).
	LabelRotation float64
	// LabelOffset is the offset of each label.
	LabelOffset OffsetInt
	// ValueFormatter defines how float values should be rendered to strings, notably for numeric axis labels.
	ValueFormatter ValueFormatter
	// Unit is a suggestion for how large the axis step is, this is a recommendation only. Larger numbers result in fewer labels.
	Unit float64
	// LabelCount is the number of labels to show on the axis. Specify a smaller number to reduce writing collisions.
	LabelCount int
	// LabelCountAdjustment specifies a relative influence on how many labels should be rendered.
	// Typically, this is negative to result in cleaner graphs, positive values may result in text collisions.
	LabelCountAdjustment int
	isValueAxis          bool
}

const defaultXAxisHeight = 30
const boundaryGapDefaultThreshold = 40

func (opt *XAxisOption) toAxisOption(fallbackTheme ColorPalette) axisOption {
	position := PositionBottom
	if opt.Position == PositionTop {
		position = PositionTop
	}
	theme := getPreferredTheme(opt.Theme, fallbackTheme)
	if opt.LabelFontStyle.IsZero() {
		opt.LabelFontStyle = opt.FontStyle
	}
	if opt.LabelFontStyle.FontColor.IsZero() {
		opt.LabelFontStyle.FontColor = theme.GetXAxisTextColor()
	}
	if opt.TitleFontStyle.FontColor.IsZero() {
		opt.TitleFontStyle.FontColor = opt.LabelFontStyle.FontColor
	}
	axisOpt := axisOption{
		show:                 opt.Show,
		title:                opt.Title,
		titleFontStyle:       opt.TitleFontStyle,
		labels:               opt.Labels,
		dataStartIndex:       opt.DataStartIndex,
		boundaryGap:          opt.BoundaryGap,
		position:             position,
		labelFontStyle:       opt.LabelFontStyle,
		axisSplitLineColor:   theme.GetAxisSplitLineColor(),
		axisColor:            theme.GetXAxisStrokeColor(),
		unit:                 opt.Unit,
		labelCount:           opt.LabelCount,
		labelCountAdjustment: opt.LabelCountAdjustment,
		labelRotation:        opt.LabelRotation,
		labelOffset:          opt.LabelOffset,
	}
	if opt.isValueAxis {
		axisOpt.splitLineShow = true
		axisOpt.strokeWidth = -1
		axisOpt.boundaryGap = Ptr(false)
	}
	return axisOpt
}

// newBottomXAxis returns a bottom x-axis renderer.
func newBottomXAxis(p *Painter, opt XAxisOption) *axisPainter {
	return newAxisPainter(p, opt.toAxisOption(p.theme))
}
