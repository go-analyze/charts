package charts

import (
	"math"
)

type CategoryAxisOption struct {
	// Show specifies if the axis should be rendered. Set to *false (via Ptr(false)) to hide the axis.
	Show *bool
	// Theme specifies the colors used for the axis.
	Theme ColorPalette
	// Title specifies a name for the axis. If set, the title is rendered adjacent to the axis.
	Title string
	// TitleFontStyle specifies the font, size, and color for the axis title.
	TitleFontStyle FontStyle
	// Labels provides labels for each value on the axis. Indices must match series data indices.
	Labels []string
	// Position controls the physical axis placement. All four position constants are accepted.
	// TODO - top-positioned category axis rendering for vertical bars is not yet supported.
	Position string
	// BoundaryGap specifies that the chart should have additional space on the left and right, with data points being
	// centered between two axis ticks. Default is set based on the dataset density / size to produce an easy-to-read
	// graph. Specify a *bool (through charts.Ptr(false) or charts.Ptr(true)) to enforce a spacing.
	BoundaryGap *bool
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
}

// XAxisOption is an alias for CategoryAxisOption. Use whatever the chart type accepts.
type XAxisOption = CategoryAxisOption

// prepAxisStyles resolves theme, label, and title font styles for either axis option type.
func prepAxisStyles(theme *ColorPalette, fallbackTheme ColorPalette, isVertical bool,
	labelFontStyle *FontStyle, titleFontStyle *FontStyle) {
	*theme = getPreferredTheme(*theme, fallbackTheme)
	textColor := (*theme).GetXAxisTextColor()
	if isVertical {
		textColor = (*theme).GetYAxisTextColor()
	}
	*labelFontStyle = fillFontStyleDefaults(*labelFontStyle, defaultFontSize, textColor)
	*titleFontStyle = fillFontStyleDefaults(*titleFontStyle, math.Max(labelFontStyle.FontSize, defaultFontSize),
		labelFontStyle.FontColor, labelFontStyle.Font)
}

func (opt *CategoryAxisOption) prep(fallbackTheme ColorPalette, isVertical bool) *CategoryAxisOption {
	prepAxisStyles(&opt.Theme, fallbackTheme, isVertical, &opt.LabelFontStyle, &opt.TitleFontStyle)
	return opt
}

// toAxisOption converts the CategoryAxisOption to axisOption after prep has been invoked.
func (opt *CategoryAxisOption) toAxisOption(xAxisRange axisRange) axisOption {
	return axisOption{
		show:           opt.Show,
		theme:          opt.Theme,
		aRange:         xAxisRange,
		title:          opt.Title,
		titleFontStyle: opt.TitleFontStyle,
		boundaryGap:    opt.BoundaryGap,
		position:       opt.Position,
		labelOffset:    opt.LabelOffset,
	}
}

// ValueAxisOption configures the value (numeric / range) axis.
type ValueAxisOption struct {
	// Show specifies if the axis should be rendered. Set to *false (via Ptr(false)) to hide the axis.
	Show *bool
	// Theme specifies the colors used for the axis.
	Theme ColorPalette
	// Title specifies a name for the axis. If set, the axis name is rendered on the outside of the axis.
	Title string
	// TitleFontStyle specifies the font, size, and color for the axis title.
	TitleFontStyle FontStyle
	// Min forces the minimum value of the axis when set (Use Ptr(float64)).
	Min *float64
	// Max forces the maximum value of the axis when set (Use Ptr(float64)).
	Max *float64
	// RangeValuePaddingScale suggests a padding scale to apply to the max and min values.
	RangeValuePaddingScale *float64
	// Labels provides labels for each value on the axis.
	Labels []string
	// Position controls the physical axis placement. All four position constants are accepted.
	// TODO - top-positioned value axis rendering is not yet supported.
	Position string
	// LabelFontStyle specifies the font configuration for each label.
	LabelFontStyle FontStyle
	// LabelRotation is the rotation angle in radians for labels. Use DegreesToRadians(float64) to convert from degrees.
	LabelRotation float64
	// Unit suggests the axis step size (recommendation only). Larger values result in fewer labels.
	Unit float64
	// LabelCount is the number of labels to show on the axis. Use a smaller count to reduce text collisions.
	LabelCount int
	// LabelCountAdjustment specifies relative influence on label count.
	// Negative values result in cleaner graphs; positive values may cause text collisions.
	LabelCountAdjustment int
	// PreferNiceIntervals allows the label count to flex slightly to produce rounder axis intervals.
	// Enabled by default when no explicit LabelCount is set; set to *false to disable.
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
	// TODO - isCategoryAxis is a hack used only by heat map so its Y-position axis
	// renders with category styling. Remove when defaultRender supports dual category axes.
	isCategoryAxis bool
}

// YAxisOption is an alias for ValueAxisOption. Use whatever the chart type accepts.
type YAxisOption = ValueAxisOption

func (opt *ValueAxisOption) prep(fallbackTheme ColorPalette, isVertical bool) *ValueAxisOption {
	prepAxisStyles(&opt.Theme, fallbackTheme, isVertical, &opt.LabelFontStyle, &opt.TitleFontStyle)
	return opt
}

// toAxisOption converts the ValueAxisOption to axisOption after prep has been invoked.
func (opt *ValueAxisOption) toAxisOption(yAxisRange axisRange) axisOption {
	return axisOption{
		show:           opt.Show,
		theme:          opt.Theme,
		aRange:         yAxisRange,
		title:          opt.Title,
		titleFontStyle: opt.TitleFontStyle,
		position:       opt.Position,
		splitLineShow:  opt.SplitLineShow,
		spineLineShow:  opt.SpineLineShow,
		isCategoryAxis: opt.isCategoryAxis,
		labelSkipCount: opt.LabelSkipCount,
	}
}

const axisMargin = 4
const minimumAxisLabels = 2            // 2 labels so range is fully shown
const minimumHorizontalAxisHeight = 24 // too small looks too crowded to the chart data, notable for horizontal bar charts
const boundaryGapDefaultThreshold = 40

type axisPainter struct {
	p   *Painter
	opt *axisOption
}

func newAxisPainter(p *Painter, opt axisOption) *axisPainter {
	return &axisPainter{
		p:   p,
		opt: &opt,
	}
}

type axisOption struct {
	show           *bool
	theme          ColorPalette
	aRange         axisRange
	title          string
	titleFontStyle FontStyle
	// position describes the axis position: 'left', 'top', 'right', or 'bottom'.
	position string
	// boundaryGap specifies that the chart should have additional space on the left and right, with data points being
	// centered between two axis ticks. Default is set based on the dataset density / size to produce an easy-to-read
	// graph. Specify a *bool to enforce a spacing.
	boundaryGap   *bool
	splitLineShow *bool // nil = painter decides based on isCategory
	spineLineShow *bool // nil = painter decides based on isCategory
	// TODO - isCategoryAxis is a hack used only by heat map so its Y-position axis
	// renders with category styling. Remove when defaultRender supports dual category axes.
	isCategoryAxis       bool
	tickLength           int
	labelMargin          int
	labelOffset          OffsetInt
	labelSkipCount       int
	painterPrePositioned bool
}

func (a *axisPainter) Render() (Box, error) {
	opt := a.opt
	if flagIs(false, opt.show) {
		return BoxZero, nil
	}

	top := a.p
	isVertical := opt.position == PositionLeft || opt.position == PositionRight
	isCategory := opt.isCategoryAxis || opt.aRange.isCategory

	// rendering defaults derived from physical position, theme, and axis data type
	axisTheme := getPreferredTheme(opt.theme, top.theme)
	axisSplitLineColor := axisTheme.GetAxisSplitLineColor()
	var axisColor Color
	var minimumAxisHeight int
	if isVertical {
		axisColor = axisTheme.GetYAxisStrokeColor()
	} else {
		axisColor = axisTheme.GetXAxisStrokeColor()
		if top.Height() >= 100 { // don't reserve space if chart is too small
			minimumAxisHeight = minimumHorizontalAxisHeight
		}
	}

	// spine line: category axes show spine, value axes hide it, user override wins
	var strokeWidth float64 = 1
	if !isCategory {
		strokeWidth = -1
	}
	if opt.spineLineShow != nil {
		if *opt.spineLineShow {
			strokeWidth = 1
		} else {
			strokeWidth = -1
		}
	}
	if strokeWidth < 0 {
		strokeWidth = 0
	}

	// split lines: value axes show them, category axes don't, user override wins
	splitLineShow := !isCategory
	if opt.splitLineShow != nil {
		splitLineShow = *opt.splitLineShow
	}

	// label margin: tighter for horizontal axes
	tickLength := getDefaultInt(opt.tickLength, 5)
	labelMargin := getDefaultInt(opt.labelMargin, 5)
	if !isVertical {
		labelMargin = 2
	}
	var axisNeededWidth, axisNeededHeight int
	if isVertical {
		axisNeededWidth = labelMargin + opt.aRange.textMaxWidth + axisMargin
		axisNeededHeight = top.Height()
	} else {
		axisNeededWidth = top.Width()
		labelMargin += opt.aRange.textMaxHeight // add height to move label past line
		axisNeededHeight = labelMargin + axisMargin
	}

	// Measure axis title and add its needed space
	var titleBox Box
	var titleShift int
	if opt.title != "" {
		// measured without rotation, we choose measurement side as appropriate
		titleBox = top.MeasureText(opt.title, 0, opt.titleFontStyle)
		// measuring without rotation also allows us to simply refer to the height as the shift for any orientation
		titleShift = titleBox.Height() + axisMargin
		if isVertical {
			axisNeededWidth += titleShift
		} else {
			axisNeededHeight += titleShift
		}
	}
	if axisNeededHeight < minimumAxisHeight {
		axisNeededHeight = minimumAxisHeight
	}

	// Build a Box to reduce the parent's painter area, so that
	// this child painter is exactly the region for the axis
	padding := Box{IsSet: true}
	if !opt.painterPrePositioned {
		switch opt.position {
		case PositionLeft:
			padding.Left = axisMargin
			padding.Right = top.Width() - axisNeededWidth
		case PositionRight:
			padding.Left = top.Width() - axisNeededWidth // margin not needed here
		case PositionTop:
			padding.Top = top.Height() - axisNeededHeight
		default: // PositionBottom
			padding.Top = top.Height() - axisNeededHeight - axisMargin
		}
	}
	child := top.Child(PainterPaddingOption(padding))

	// draw axis title
	if opt.title != "" {
		switch opt.position {
		case PositionLeft:
			cx := child.Height() >> 1
			xTitle := titleShift >> 1 // No margin adjustment on this one for some reason
			yTitle := cx + (titleBox.Width() >> 1)
			child.Text(opt.title, xTitle, yTitle, DegreesToRadians(270), opt.titleFontStyle)
		case PositionRight:
			cx := child.Height() >> 1
			xTitle := child.Width() - (titleShift >> 1) - axisMargin
			yTitle := cx - (titleBox.Width() >> 1)
			child.Text(opt.title, xTitle, yTitle, DegreesToRadians(90), opt.titleFontStyle)
		case PositionTop:
			xTitle := (child.Width() - titleBox.Width()) >> 1
			yTitle := titleShift >> 1
			child.Text(opt.title, xTitle, yTitle, 0, opt.titleFontStyle)
		default: // PositionBottom
			xTitle := (child.Width() - titleBox.Width()) >> 1
			yTitle := child.Height() - axisMargin
			child.Text(opt.title, xTitle, yTitle, 0, opt.titleFontStyle)
		}
	}

	// draw axis line
	if strokeWidth > 0 {
		var x0, y0, x1, y1 int
		switch opt.position {
		case PositionLeft:
			x0, x1 = child.Width(), child.Width()
			y0, y1 = 0, child.Height()
		case PositionRight:
			x0, x1 = 0, 0
			y0, y1 = 0, child.Height()
		case PositionTop:
			x0, x1 = 0, child.Width()
			y0, y1 = child.Height(), child.Height()
		default: // PositionBottom
			x0, x1 = 0, child.Width()
			y0, y1 = 0, 0
		}
		child.LineStroke([]Point{
			{X: x0, Y: y0},
			{X: x1, Y: y1},
		}, axisColor, strokeWidth)
	}

	rangeLabels := opt.aRange.labels
	// vertical: reverse to match multitext expectations (draws from top down)
	// horizontal + reversed: right-anchored value axis runs max-to-min left-to-right
	if isVertical || opt.aRange.reversed {
		// TODO - replace copy with slices.Clone with go update
		rangeLabels = make([]string, len(opt.aRange.labels))
		copy(rangeLabels, opt.aRange.labels)
		reverseSlice(rangeLabels)
	}

	// Decide whether to center the labels between ticks or align them
	centerLabels := isCategory // category axes default to boundary gap, value axes don't
	if opt.boundaryGap != nil {
		centerLabels = *opt.boundaryGap
	} else if opt.aRange.isCategory && !isVertical && opt.aRange.divideCount > 1 &&
		top.Width()/opt.aRange.divideCount <= boundaryGapDefaultThreshold {
		// for dense datasets it's visually better to have the label aligned to the tick mark
		// this default is also handled in the chart rendering to ensure data aligns with the labels
		centerLabels = false
	}

	tickSpaces := opt.aRange.tickCount
	tickCount := opt.aRange.tickCount
	if centerLabels {
		// In order to center the labels we need an extra tick mark to center the labels between
		tickCount++
	} else {
		// there is always one more tick than data sample, and if we are centering labels we use that extra tick to
		// center the label against, if not centering then we need one less tick spacing
		// passing the tickSpaces reduces the need to copy the logic from painter.go:multiText
		tickSpaces--
	}

	// draw tick marks
	if strokeWidth > 0 {
		var tickPaddingBox Box
		tickPaddingBox.IsSet = true
		switch opt.position {
		case PositionLeft:
			tickPaddingBox.Left = child.Width() - tickLength
		case PositionRight:
			tickPaddingBox.Right = tickLength
		case PositionTop:
			tickPaddingBox.Top = child.Height() - tickLength
		default: // PositionBottom
			tickPaddingBox.Bottom = tickLength
		}
		tickPainter := child.Child(PainterPaddingOption(tickPaddingBox))
		tickPainter.ticks(ticksOption{
			tickCount:   tickCount,
			tickSpaces:  tickSpaces,
			length:      tickLength,
			vertical:    isVertical,
			strokeWidth: strokeWidth,
			strokeColor: axisColor,
		})
	}

	// Render tick labels
	labelPadding := Box{IsSet: true}
	switch opt.position {
	case PositionLeft:
		// Place labels to the left, so we leave margin on the right
		labelPadding.Right = tickLength + labelMargin
		labelPadding.Top = -2 // TODO - it's unclear why this adjustment is needed for vertical to position labels right
		labelPadding.Bottom = 4
	case PositionRight:
		labelPadding.Left = tickLength + labelMargin
		labelPadding.Top = -2
		labelPadding.Bottom = 4
	case PositionTop:
		labelPadding.Bottom = tickLength + labelMargin
	default: // PositionBottom
		labelPadding.Top = tickLength + labelMargin
		if opt.aRange.labelRotation != 0 {
			flatWidth, flatHeight :=
				top.measureTextMaxWidthHeight(opt.aRange.labels, 0, opt.aRange.labelFontStyle)
			labelPadding.Top -= textRotationHeightAdjustment(flatWidth, flatHeight, opt.aRange.labelRotation)
		}
	}
	labelPainter := child.Child(PainterPaddingOption(labelPadding))
	alignSide := AlignCenter
	if isVertical {
		if opt.position == PositionLeft {
			alignSide = AlignRight
		} else {
			alignSide = AlignLeft
		}
	}
	labelPainter.multiText(multiTextOption{
		textList:       rangeLabels,
		vertical:       isVertical,
		centerLabels:   centerLabels,
		align:          alignSide,
		textRotation:   opt.aRange.labelRotation,
		offset:         opt.labelOffset,
		labelCount:     opt.aRange.labelCount,
		labelSkipCount: opt.labelSkipCount,
		fontStyle:      opt.aRange.labelFontStyle,
	})

	if splitLineShow { // show auxiliary lines
		if isVertical {
			var x0Split, x1Split int
			if opt.position == PositionLeft {
				x0Split = child.Width()
				x1Split = top.Width()
			} else { // right axis
				x0Split = 0
				x1Split = top.Width() - child.Width()
			}
			yValues := autoDivide(child.Height(), tickSpaces)
			// Skip the last one to avoid re-drawing the axis line
			if len(yValues) > 0 {
				yValues = yValues[:len(yValues)-1]
			}
			for _, yy := range yValues {
				top.LineStroke([]Point{
					{X: x0Split, Y: yy},
					{X: x1Split, Y: yy},
				}, axisSplitLineColor, 1)
			}
		} else {
			var y0Split, y1Split int
			if opt.position == PositionTop {
				y0Split = child.Height()
				y1Split = top.Height()
			} else { // PositionBottom
				y0Split = 0
				y1Split = top.Height() - child.Height()
			}
			xValues := autoDivide(child.Width(), tickSpaces)
			for i, xx := range xValues {
				if i == 0 {
					continue // skip the first, so we don't overlap the axis line
				}
				top.LineStroke([]Point{
					{X: xx, Y: y0Split},
					{X: xx, Y: y1Split},
				}, axisSplitLineColor, 1)
			}
		}
	}

	// Return the "used" dimension for this axis
	// This consumed space will be removed from the chart space in defaultRender
	return Box{
		Right:  axisNeededWidth + ceilFloatToInt(strokeWidth),
		Bottom: axisNeededHeight + ceilFloatToInt(strokeWidth),
		IsSet:  true,
	}, nil
}
