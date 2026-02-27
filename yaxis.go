package charts

import (
	"math"
)

// YAxisOption configures the vertical axis.
type YAxisOption struct {
	// Show specifies if the y-axis should be rendered. Set to *false (via Ptr(false)) to hide the axis.
	Show *bool
	// Theme specifies the colors used for the y-axis.
	Theme ColorPalette
	// Title specifies a name for the axis. If set, the axis name is rendered on the outside of the y-Axis.
	Title string
	// TitleFontStyle specifies the font, size, and color for the axis title.
	TitleFontStyle FontStyle
	// Min forces the minimum value of the y-axis when set (Use Ptr(float64)).
	Min *float64
	// Max forces the maximum value of the y-axis when set (Use Ptr(float64)).
	Max *float64
	// RangeValuePaddingScale suggests a padding scale to apply to the max and min values.
	RangeValuePaddingScale *float64
	// Labels provides labels for each value on the y-axis.
	Labels []string
	// Position describes the y-axis position: 'left' or 'right'.
	Position string
	// Deprecated: FontStyle is deprecated, use LabelFontStyle.
	FontStyle FontStyle
	// LabelFontStyle specifies the font configuration for each label.
	LabelFontStyle FontStyle
	// LabelRotation is the rotation angle in radians for labels. Use DegreesToRadians(float64) to convert from degrees.
	LabelRotation float64
	// Deprecated: Formatter is deprecated, use ValueFormatter instead.
	Formatter string
	// Unit suggests the axis step size (recommendation only). Larger values result in fewer labels.
	Unit float64
	// LabelCount is the number of labels to show on the axis. Use a smaller count to reduce text collisions.
	LabelCount int
	// LabelCountAdjustment specifies relative influence on label count.
	// Negative values result in cleaner graphs; positive values may cause text collisions.
	LabelCountAdjustment int
	// PreferNiceIntervals allows the label count to flex slightly to produce rounder axis intervals.
	// TODO - reconsider default for v0.6.0 (possibly with padding default changes)
	PreferNiceIntervals *bool
	// LabelSkipCount specifies a qty of lines between labels that show only horizontal lines without labels.
	LabelSkipCount int
	// SplitLineShow when set to *true shows horizontal axis split lines.
	SplitLineShow *bool
	// SpineLineShow controls whether the vertical spine line is shown.
	// Default is hidden unless it's a category axis.
	SpineLineShow *bool
	// ValueFormatter defines how float values are rendered to strings, notably for numeric axis labels.
	ValueFormatter ValueFormatter
	isCategoryAxis bool
}

func (opt *YAxisOption) prep(fallbackTheme ColorPalette) *YAxisOption {
	opt.Theme = getPreferredTheme(opt.Theme, fallbackTheme)
	if opt.LabelFontStyle.IsZero() {
		opt.LabelFontStyle = opt.FontStyle
	}
	opt.LabelFontStyle = fillFontStyleDefaults(opt.LabelFontStyle, defaultFontSize,
		opt.Theme.GetYAxisTextColor())
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
