package charts

import (
	"math"
)

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
	// LabelRotation are the radians for rotating the label. Convert from degrees using DegreesToRadians(float64).
	LabelRotation float64
	// Deprecated: Formatter is deprecated, use ValueFormatter instead.
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

func (opt *YAxisOption) prep(fallbackTheme ColorPalette) *YAxisOption {
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

// toAxisOption converts the YAxisOption to axisOption after prep has been invoked.
func (opt *YAxisOption) toAxisOption(yAxisRange axisRange) axisOption {
	axisOpt := axisOption{
		show:               opt.Show,
		aRange:             yAxisRange,
		title:              opt.Title,
		titleFontStyle:     opt.TitleFontStyle,
		position:           opt.Position,
		axisSplitLineColor: opt.Theme.GetAxisSplitLineColor(),
		axisColor:          opt.Theme.GetYAxisStrokeColor(),
		strokeWidth:        -1,
		boundaryGap:        Ptr(false),
		labelSkipCount:     opt.LabelSkipCount,
		splitLineShow:      true,
	}
	if opt.isCategoryAxis || yAxisRange.isCategory {
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
