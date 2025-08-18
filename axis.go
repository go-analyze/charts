package charts

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
	aRange         axisRange
	title          string
	titleFontStyle FontStyle
	// position describes the axis position: 'left', 'top', 'right', or 'bottom'.
	position string
	// boundaryGap specifies that the chart should have additional space on the left and right, with data points being
	// centered between two axis ticks. Default is set based on the dataset density / size to produce an easy-to-read
	// graph. Specify a *bool to enforce a spacing.
	boundaryGap          *bool
	strokeWidth          float64
	minimumAxisHeight    int
	tickLength           int
	labelMargin          int
	axisSplitLineColor   Color
	axisColor            Color
	splitLineShow        bool
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
	strokeWidth := opt.strokeWidth
	if strokeWidth == 0 {
		strokeWidth = 1
	} else if strokeWidth < 0 {
		strokeWidth = 0 // set to zero as negative values will confuse calculations below
	}

	// calculate how much space the axis line + tick marks + labels need
	tickLength := getDefaultInt(opt.tickLength, 5)
	labelMargin := getDefaultInt(opt.labelMargin, 5)
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
		}, opt.axisColor, strokeWidth)
	}

	rangeLabels := opt.aRange.labels
	if isVertical { // reverse to match multitext expectations (draws from top down)
		// make copy first to avoid changing range slice
		rangeLabels = make([]string, len(opt.aRange.labels))
		copy(rangeLabels, opt.aRange.labels)
		reverseSlice(rangeLabels)
	}

	// Decide whether to center the labels between ticks or align them
	centerLabels := true
	if opt.boundaryGap != nil {
		centerLabels = *opt.boundaryGap
	} else if opt.aRange.divideCount > 1 && top.Width()/opt.aRange.divideCount <= boundaryGapDefaultThreshold {
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
			firstIndex:  opt.aRange.dataStartIndex,
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
		firstIndex:     opt.aRange.dataStartIndex,
		labelCount:     opt.aRange.labelCount,
		labelSkipCount: opt.labelSkipCount,
		fontStyle:      opt.aRange.labelFontStyle,
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
			// Skip the last one to avoid re-drawing the axis line
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
