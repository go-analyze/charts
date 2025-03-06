package charts

import (
	"strings"
)

const axisMargin = 4
const minimumAxisLabels = 2            // 2 labels so range is fully shown
const minimumHorizontalAxisHeight = 24 // too small looks too crowded to the chart data, notable for horizontal bar charts

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
	title          string
	titleFontStyle FontStyle
	// labels provides string values for each value on the axis.
	labels []string
	// dataStartIndex specifies what index the label values should start from.
	dataStartIndex int
	// formatter for replacing axis text values.
	formatter string
	// position describes the position of axis, it can be 'left', 'top', 'right' or 'bottom'.
	position string
	// boundaryGap specifies that the chart should have additional space on the left and right, with data points being
	// centered between two axis ticks. Default is set based on the dataset density / size to produce an easy-to-read
	// graph. Specify a *bool to enforce a spacing.
	boundaryGap        *bool
	strokeWidth        float64
	minimumAxisHeight  int
	tickLength         int
	labelMargin        int
	labelFontStyle     FontStyle
	axisSplitLineColor Color
	axisColor          Color
	splitLineShow      bool
	// labelRotation are the radians for rotating the label.
	labelRotation        float64
	labelOffset          OffsetInt
	unit                 float64
	labelCount           int
	labelCountAdjustment int
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

	defaultTheme := getPreferredTheme(top.theme)
	strokeWidth := opt.strokeWidth
	if strokeWidth == 0 {
		strokeWidth = 1
	} else if strokeWidth < 0 {
		strokeWidth = 0 // set to zero as negative values will confuse calculations below
	}
	if opt.axisColor.IsZero() {
		if isVertical {
			opt.axisColor = defaultTheme.GetYAxisStrokeColor()
		} else {
			opt.axisColor = defaultTheme.GetXAxisStrokeColor()
		}
	}
	if opt.axisSplitLineColor.IsZero() {
		opt.axisSplitLineColor = defaultTheme.GetAxisSplitLineColor()
	}
	opt.labelFontStyle.Font = getPreferredFont(opt.labelFontStyle.Font, top.font)
	if opt.labelFontStyle.FontColor.IsZero() {
		if isVertical {
			opt.labelFontStyle.FontColor = defaultTheme.GetYAxisTextColor()
		} else {
			opt.labelFontStyle.FontColor = defaultTheme.GetXAxisTextColor()
		}
	}
	if opt.labelFontStyle.FontSize == 0 {
		opt.labelFontStyle.FontSize = defaultFontSize
	}

	if opt.formatter != "" {
		for i, text := range opt.labels {
			opt.labels[i] = strings.ReplaceAll(opt.formatter, "{value}", text)
		}
	}

	// Measure labels to determine space needed
	textMaxWidth, textMaxHeight := top.measureTextMaxWidthHeight(opt.labels, opt.labelRotation, opt.labelFontStyle)

	// Basic measurements for ticks + labels
	tickLength := getDefaultInt(opt.tickLength, 5)
	labelMargin := getDefaultInt(opt.labelMargin, 5)
	var axisNeededWidth, axisNeededHeight int
	if isVertical {
		axisNeededWidth = labelMargin + textMaxWidth + axisMargin
		axisNeededHeight = top.Height()
	} else {
		axisNeededWidth = top.Width()
		labelMargin += textMaxHeight // add height to move label past line
		axisNeededHeight = labelMargin + axisMargin
	}

	// Measure axis title and add its needed space
	var titleBox Box
	var titleShift int
	if opt.title != "" {
		// Fill title font defaults with label's color if not set
		opt.titleFontStyle =
			fillFontStyleDefaults(opt.titleFontStyle, defaultFontSize, opt.labelFontStyle.FontColor, top.font)
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
	if axisNeededHeight < opt.minimumAxisHeight {
		axisNeededHeight = opt.minimumAxisHeight
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

	// If we have a title, draw it onto the child painter.
	// We'll place it in the correct spot relative to the axis line.
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

	// Draw the main axis line at whichever edge is "touching" the chart region
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
		}, opt.axisColor, strokeWidth)
	}

	// Decide how many of the labels to draw (labelCount + logic for skipping, etc.)
	dataCount := len(opt.labels)
	labelCount := opt.labelCount
	if labelCount <= 0 {
		var maxLabelCount int
		// Add 10px and remove one for some minimal extra padding so that letters don't collide
		if isVertical {
			maxLabelCount = (child.Height() / (textMaxHeight + 10)) - 1
		} else {
			maxLabelCount = (child.Width() / (textMaxWidth + 10)) - 1
		}
		if maxLabelCount < minimumAxisLabels {
			maxLabelCount = minimumAxisLabels // required to prevent infinite loop if less than zero
		}
		if opt.unit > 0 {
			// If the user gave a 'unit', figure out how many 'units' fit
			multiplier := 1.0
			for {
				labelCount = ceilFloatToInt(float64(dataCount) / (opt.unit * multiplier))
				if labelCount > maxLabelCount {
					multiplier++
				} else {
					break
				}
			}
		} else {
			// TODO - check if a small adjustment allows for a better unit
			labelCount = maxLabelCount
		}
	}
	if labelCount > dataCount {
		labelCount = dataCount
	}
	labelCount += opt.labelCountAdjustment
	if labelCount < minimumAxisLabels {
		labelCount = minimumAxisLabels
	}

	// Decide whether to center the labels between ticks or align them
	centerLabels := true
	if opt.boundaryGap != nil {
		centerLabels = *opt.boundaryGap
	} else if dataCount > 1 && top.Width()/dataCount <= boundaryGapDefaultThreshold {
		// for dense datasets it's visually better to have the label aligned to the tick mark
		// this default is also handled in the chart rendering to ensure data aligns with the labels
		centerLabels = false
	}

	tickSpaces := dataCount
	tickCount := labelCount
	if centerLabels {
		// In order to center the labels we need an extra tick mark to center the labels between
		tickCount++
	} else {
		// there is always one more tick than data sample, and if we are centering labels we use that extra tick to
		// center the label against, if not centering then we need one less tick spacing
		// passing the tickSpaces reduces the need to copy the logic from painter.go:multiText
		tickSpaces--
	}

	// Child painter for drawing the *tick marks* themselves
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
			labelCount:  labelCount,
			tickCount:   tickCount,
			tickSpaces:  tickSpaces,
			length:      tickLength,
			vertical:    isVertical,
			firstIndex:  opt.dataStartIndex,
			strokeWidth: strokeWidth,
			strokeColor: opt.axisColor,
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
		textList:       opt.labels,
		vertical:       isVertical,
		centerLabels:   centerLabels,
		align:          alignSide,
		textRotation:   opt.labelRotation,
		offset:         opt.labelOffset,
		firstIndex:     opt.dataStartIndex,
		labelCount:     labelCount,
		tickCount:      tickCount,
		labelSkipCount: opt.labelSkipCount,
		fontStyle:      opt.labelFontStyle,
	})

	if opt.splitLineShow { // show auxiliary lines
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
			// Typically skip the last one to avoid re-drawing the axis line
			if len(yValues) > 0 {
				yValues = yValues[:len(yValues)-1]
			}
			for _, yy := range yValues {
				top.LineStroke([]Point{
					{X: x0Split, Y: yy},
					{X: x1Split, Y: yy},
				}, opt.axisSplitLineColor, 1)
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
				}, opt.axisSplitLineColor, 1)
			}
		}
	}

	// Return the “used” dimension for this axis
	// This consumed space will be removed from the chart space in defaultRender
	return Box{
		Right:  axisNeededWidth + ceilFloatToInt(strokeWidth),
		Bottom: axisNeededHeight + ceilFloatToInt(strokeWidth),
		IsSet:  true,
	}, nil
}
