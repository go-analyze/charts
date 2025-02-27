package charts

import (
	"math"
	"strings"
)

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
}

func (a *axisPainter) Render() (Box, error) {
	opt := a.opt
	if flagIs(false, opt.show) {
		return BoxZero, nil
	}
	top := a.p
	defaultTheme := getPreferredTheme(top.theme)
	isVertical := opt.position == PositionLeft || opt.position == PositionRight

	strokeWidth := opt.strokeWidth
	if strokeWidth == 0 {
		strokeWidth = 1
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

	formatter := opt.formatter
	if formatter != "" {
		for index, text := range opt.labels {
			opt.labels[index] = strings.ReplaceAll(formatter, "{value}", text)
		}
	}

	tickLength := getDefaultInt(opt.tickLength, 5)
	labelMargin := getDefaultInt(opt.labelMargin, 5)

	textMaxWidth, textMaxHeight := top.measureTextMaxWidthHeight(opt.labels, opt.labelRotation, opt.labelFontStyle)

	var width, height int
	if isVertical {
		width = textMaxWidth + (tickLength << 1)
		height = top.Height()
	} else {
		width = top.Width()
		height = textMaxHeight + (tickLength << 1)
	}
	padding := Box{IsSet: true}
	switch opt.position {
	case PositionTop:
		padding.Top = top.Height() - height
	case PositionLeft:
		padding.Right = top.Width() - width
	case PositionRight:
		padding.Left = top.Width() - width
	default: // PositionBottom
		padding.Top = top.Height() - defaultXAxisHeight
	}

	// Measure and render the axis title to shift the axis
	var titleShift int
	if opt.title != "" {
		opt.titleFontStyle = fillFontStyleDefaults(opt.titleFontStyle, defaultFontSize,
			opt.labelFontStyle.FontColor, opt.titleFontStyle.Font, top.font)

		var titleX, titleY int
		var titleRotation float64
		switch opt.position {
		case PositionLeft:
			titleRotation = DegreesToRadians(270)
			titleBox := top.MeasureText(opt.title, titleRotation, opt.titleFontStyle)
			titleX = 16
			titleY = (top.Height() / 2) + (titleBox.Height() / 2)
			titleShift = titleBox.Width() + 8
			padding.Right -= titleShift
		case PositionRight:
			titleRotation = DegreesToRadians(90)
			titleBox := top.MeasureText(opt.title, titleRotation, opt.titleFontStyle)
			titleX = top.Width() - 16
			titleY = (top.Height() / 2) - (titleBox.Height() / 2)
			titleShift = titleBox.Width() + 8
			padding.Left -= titleShift
		default: // horizontal top / bottom
			titleBox := top.MeasureText(opt.title, titleRotation, opt.titleFontStyle)
			titleX = (top.Width() / 2) - (titleBox.Width() / 2)
			titleY = top.Height() - 2
			titleShift = titleBox.Height()
			padding.Top -= titleShift
		}

		top.Text(opt.title, titleX, titleY, titleRotation, opt.titleFontStyle)
	}

	// Create a child painter using the adjusted padding for the axis
	p := top.Child(PainterPaddingOption(padding))

	// Set up variables used for drawing ticks and labels
	var x0, y0, x1, y1 int
	var ticksPaddingTop, ticksPaddingLeft int
	var labelPaddingTop, labelPaddingLeft, labelPaddingRight int
	var textAlign string
	switch opt.position {
	case PositionTop:
		if opt.labelRotation != 0 {
			flatWidth, flatHeight := top.measureTextMaxWidthHeight(opt.labels, 0, opt.labelFontStyle)
			labelPaddingTop = flatHeight - textRotationHeightAdjustment(flatWidth, flatHeight, opt.labelRotation)
		} else {
			labelPaddingTop = 0
		}
		x1 = p.Width()
		y0 = labelMargin + int(opt.labelFontStyle.FontSize)
		ticksPaddingTop = int(opt.labelFontStyle.FontSize)
		y1 = y0
	case PositionLeft:
		x0 = p.Width()
		y0 = 0
		x1 = p.Width()
		y1 = p.Height()
		textAlign = AlignRight
		ticksPaddingLeft = textMaxWidth + tickLength
		labelPaddingRight = width - textMaxWidth
	case PositionRight:
		y1 = p.Height()
		labelPaddingLeft = width - textMaxWidth
	default: // PositionBottom
		if opt.labelRotation != 0 {
			flatWidth, flatHeight := top.measureTextMaxWidthHeight(opt.labels, 0, opt.labelFontStyle)
			labelPaddingTop = (tickLength << 1) + (textMaxHeight - textRotationHeightAdjustment(flatWidth, flatHeight, opt.labelRotation))
		} else {
			labelPaddingTop = height
		}
		x1 = p.Width()
	}

	dataCount := len(opt.labels)
	labelCount := opt.labelCount
	if labelCount <= 0 {
		var maxLabelCount int
		// Add 10px and remove one for some minimal extra padding so that letters don't collide
		if isVertical {
			maxLabelCount = (top.Height() / (textMaxHeight + 10)) - 1
		} else {
			maxLabelCount = (top.Width() / (textMaxWidth + 10)) - 1
		}
		if opt.unit > 0 {
			multiplier := 1.0
			for {
				labelCount = int(math.Ceil(float64(dataCount) / (opt.unit * multiplier)))
				if labelCount > maxLabelCount {
					multiplier++
				} else {
					break
				}
			}
		} else {
			labelCount = maxLabelCount
		}
	}
	if labelCount > dataCount {
		labelCount = dataCount
	}
	labelCount += opt.labelCountAdjustment
	if labelCount < 2 {
		labelCount = 2
	}

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

	// Render ticks and axis line
	if strokeWidth > 0 {
		p.Child(PainterPaddingOption(Box{
			Top:   ticksPaddingTop,
			Left:  ticksPaddingLeft,
			IsSet: true,
		})).ticks(ticksOption{
			labelCount:  labelCount,
			tickCount:   tickCount,
			tickSpaces:  tickSpaces,
			length:      tickLength,
			vertical:    isVertical,
			firstIndex:  opt.dataStartIndex,
			strokeWidth: strokeWidth,
			strokeColor: opt.axisColor,
		})
		p.LineStroke([]Point{
			{X: x0, Y: y0},
			{X: x1, Y: y1},
		}, opt.axisColor, strokeWidth)
	}

	// Render tick labels
	p.Child(PainterPaddingOption(Box{
		Left:  labelPaddingLeft,
		Top:   labelPaddingTop,
		Right: labelPaddingRight,
		IsSet: true,
	})).multiText(multiTextOption{
		firstIndex:     opt.dataStartIndex,
		align:          textAlign,
		textList:       opt.labels,
		fontStyle:      opt.labelFontStyle,
		vertical:       isVertical,
		labelCount:     labelCount,
		tickCount:      tickCount,
		labelSkipCount: opt.labelSkipCount,
		centerLabels:   centerLabels,
		textRotation:   opt.labelRotation,
		offset:         opt.labelOffset,
	})

	if opt.splitLineShow { // show auxiliary lines
		if isVertical {
			x0 := p.Width()
			x1 := top.Width()
			if opt.position == PositionRight {
				x0 = 0
				x1 = top.Width() - p.Width()
			}
			yValues := autoDivide(height, tickSpaces)
			yValues = yValues[0 : len(yValues)-1]
			for _, y := range yValues {
				top.LineStroke([]Point{
					{X: x0, Y: y},
					{X: x1, Y: y},
				}, opt.axisSplitLineColor, 1)
			}
		} else {
			y0 := p.Height() - defaultXAxisHeight
			y1 := top.Height() - defaultXAxisHeight
			xValues := autoDivide(width, tickSpaces)
			for index, x := range xValues {
				if index == 0 {
					continue
				}
				top.LineStroke([]Point{
					{X: x, Y: y0},
					{X: x, Y: y1},
				}, opt.axisSplitLineColor, 1)
			}
		}
	}

	// we need to adjust the dimensions for the title here, adjusting earlier would change the axis draw box
	if isVertical {
		width += titleShift
	} else {
		height += titleShift
	}
	return Box{
		Bottom: height,
		Right:  width,
		IsSet:  true,
	}, nil
}
