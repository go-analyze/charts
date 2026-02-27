package charts

import (
	"math"
)

// XAxisOption configures the horizontal axis.
type XAxisOption struct {
	// Show specifies if the x-axis should be rendered. Set to *false (via Ptr(false)) to hide the axis.
	Show *bool
	// Theme specifies the colors used for the x-axis.
	Theme ColorPalette
	// Title specifies a name for the axis. If set, the title is rendered below the x-axis.
	Title string
	// TitleFontStyle specifies the font, size, and color for the axis title.
	TitleFontStyle FontStyle
	// Labels provides labels for each value on the x-axis. Indices must match series data indices.
	Labels []string
	// DataStartIndex specifies the starting index for data values.
	DataStartIndex int
	// Deprecated: Position is deprecated. Currently, when set to `bottom` and the labels would render on the top
	// side of the axis line. However, the line would remain at the bottom of the chart. This seems confusing, and
	// attempts to actually move the axis line to the top of the chart are currently very messy looking. For that
	// reason this is currently deprecated. If a top x-Axis is valuable to you, please open a feature request.
	Position string
	// BoundaryGap specifies that the chart should have additional space on the left and right, with data points being
	// centered between two axis ticks. Default is set based on the dataset density / size to produce an easy-to-read
	// graph. Specify a *bool (through charts.Ptr(false) or charts.Ptr(true)) to enforce a spacing.
	BoundaryGap *bool
	// Deprecated: FontStyle is deprecated, use LabelFontStyle.
	FontStyle FontStyle
	// LabelFontStyle specifies the font configuration for each label.
	LabelFontStyle FontStyle
	// LabelRotation is the rotation angle in radians for labels. Use DegreesToRadians(float64) to convert from degrees.
	LabelRotation float64
	// LabelOffset is the position offset for each label.
	LabelOffset OffsetInt
	// ValueFormatter defines how float values are rendered to strings, notably for numeric axis labels.
	ValueFormatter ValueFormatter
	// Unit suggests the axis step size (recommendation only). Larger values result in fewer labels.
	Unit float64
	// LabelCount is the number of labels to show on the axis. Use a smaller count to reduce text collisions.
	LabelCount int
	// LabelCountAdjustment specifies a relative influence on how many labels should be rendered.
	// Typically, this is negative to result in cleaner graphs, positive values may result in text collisions.
	LabelCountAdjustment int
	// TODO - add PreferNiceIntervals, see https://github.com/go-analyze/charts/issues/63
}

const boundaryGapDefaultThreshold = 40

func (opt *XAxisOption) prep(fallbackTheme ColorPalette) *XAxisOption {
	opt.Theme = getPreferredTheme(opt.Theme, fallbackTheme)
	if opt.LabelFontStyle.IsZero() {
		opt.LabelFontStyle = opt.FontStyle
	}
	opt.LabelFontStyle = fillFontStyleDefaults(opt.LabelFontStyle, defaultFontSize,
		opt.Theme.GetXAxisTextColor())
	opt.TitleFontStyle = fillFontStyleDefaults(opt.TitleFontStyle, math.Max(opt.LabelFontStyle.FontSize, defaultFontSize),
		opt.LabelFontStyle.FontColor, opt.LabelFontStyle.Font)
	return opt
}

// toAxisOption converts the XAxisOption to axisOption after prep has been invoked.
func (opt *XAxisOption) toAxisOption(xAxisRange axisRange) axisOption {
	position := PositionBottom
	if opt.Position == PositionTop {
		position = PositionTop
	}
	axisOpt := axisOption{
		show:               opt.Show,
		aRange:             xAxisRange,
		title:              opt.Title,
		titleFontStyle:     opt.TitleFontStyle,
		boundaryGap:        opt.BoundaryGap,
		position:           position,
		minimumAxisHeight:  minimumHorizontalAxisHeight,
		axisSplitLineColor: opt.Theme.GetAxisSplitLineColor(),
		axisColor:          opt.Theme.GetXAxisStrokeColor(),
		labelOffset:        opt.LabelOffset,
	}
	if !xAxisRange.isCategory {
		axisOpt.splitLineShow = true
		axisOpt.strokeWidth = -1
		axisOpt.boundaryGap = Ptr(false)
		axisOpt.labelMargin = 2
	}
	return axisOpt
}
