package charts

type YAxisOption struct {
	// Show specifies if the y-axis should be rendered, set this to *false (through Ptr(false)) to hide the axis.
	Show *bool
	// Theme specifies the colors used for the x-axis.
	Theme ColorPalette
	// Title specifies a name for the axis, if specified the axis name is rendered on the outside of the Y-Axis.
	Title string
	// TitleFontStyle provides the font, size, and color for the axis title.
	TitleFontStyle FontStyle
	// Min, if set this will force the minimum value of y-axis.
	Min *float64
	// Max, if set this will force the maximum value of y-axis.
	Max *float64
	// RangeValuePaddingScale suggest a scale of padding added to the max and min values.
	RangeValuePaddingScale *float64
	// Labels provides labels for each value on the y-axis.
	Labels []string
	// Position describes the position of y-axis, it can be 'left' or 'right'.
	Position string
	// Deprecated: FontStyle is deprecated, use LabelFontStyle.
	FontStyle FontStyle
	// LabelFontStyle specifies the font configuration for each label.
	LabelFontStyle FontStyle
	// Formatter for replacing y-axis text values.
	Formatter string
	// Unit is a suggestion for how large the axis step is, this is a recommendation only. Larger numbers result in fewer labels.
	Unit float64
	// LabelCount is the number of labels to show on the axis. Specify a smaller number to reduce writing collisions.
	LabelCount int
	// LabelCountAdjustment specifies a relative influence on how many labels should be rendered.
	// Typically, this is negative to result in cleaner graphs, positive values may result in text collisions.
	LabelCountAdjustment int
	// LabelSkipCount specifies a number of lines between labels where there will be no label and instead just a horizontal line.
	LabelSkipCount int
	// SplitLineShow for showing axis split line, set this to true to show the horizontal axis split lines.
	SplitLineShow *bool
	// SpineLineShow can be set to enforce if the vertical spine on the axis should be shown or not.
	// By default, not shown unless a category axis.
	SpineLineShow *bool
	// ValueFormatter defines how float values should be rendered to strings, notably for numeric axis labels.
	ValueFormatter ValueFormatter
	isCategoryAxis bool
}

func (opt *YAxisOption) toAxisOption(fallbackTheme ColorPalette) axisOption {
	position := PositionLeft
	if opt.Position == PositionRight {
		position = PositionRight
	}
	theme := getPreferredTheme(opt.Theme, fallbackTheme)
	if opt.LabelFontStyle.IsZero() {
		opt.LabelFontStyle = opt.FontStyle
	}
	if opt.LabelFontStyle.FontColor.IsZero() {
		opt.LabelFontStyle.FontColor = theme.GetYAxisTextColor()
	}
	if opt.TitleFontStyle.FontColor.IsZero() {
		opt.TitleFontStyle.FontColor = opt.LabelFontStyle.FontColor
	}
	axisOpt := axisOption{
		show:                 opt.Show,
		title:                opt.Title,
		titleFontStyle:       opt.TitleFontStyle,
		labels:               opt.Labels,
		formatter:            opt.Formatter,
		position:             position,
		labelFontStyle:       opt.LabelFontStyle,
		axisSplitLineColor:   theme.GetAxisSplitLineColor(),
		axisColor:            theme.GetYAxisStrokeColor(),
		strokeWidth:          -1,
		boundaryGap:          Ptr(false),
		unit:                 opt.Unit,
		labelCount:           opt.LabelCount,
		labelCountAdjustment: opt.LabelCountAdjustment,
		labelSkipCount:       opt.LabelSkipCount,
		splitLineShow:        true,
	}
	if opt.isCategoryAxis {
		axisOpt.boundaryGap = Ptr(true)
		axisOpt.strokeWidth = 1
		axisOpt.splitLineShow = false
	}
	if opt.SplitLineShow != nil {
		axisOpt.splitLineShow = *opt.SplitLineShow
	}
	if opt.SpineLineShow != nil {
		if *opt.SpineLineShow {
			axisOpt.strokeWidth = 1
		} else {
			axisOpt.strokeWidth = -1
		}
	}
	return axisOpt
}
