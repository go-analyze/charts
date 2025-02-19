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
	show *bool
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
	fontStyle          FontStyle
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
		opt.axisColor = defaultTheme.GetAxisStrokeColor()
	}
	if opt.axisSplitLineColor.IsZero() {
		opt.axisSplitLineColor = defaultTheme.GetAxisSplitLineColor()
	}
	fontStyle := opt.fontStyle
	fontStyle.Font = getPreferredFont(fontStyle.Font, a.p.font)
	if fontStyle.FontColor.IsZero() {
		fontStyle.FontColor = defaultTheme.GetTextColor()
	}
	if fontStyle.FontSize == 0 {
		fontStyle.FontSize = defaultFontSize
	}

	formatter := opt.formatter
	if formatter != "" {
		for index, text := range opt.labels {
			opt.labels[index] = strings.ReplaceAll(formatter, "{value}", text)
		}
	}

	// if less than zero, it means not processing
	tickLength := getDefaultInt(opt.tickLength, 5)
	labelMargin := getDefaultInt(opt.labelMargin, 5)

	textMaxWidth, textMaxHeight := top.measureTextMaxWidthHeight(opt.labels, opt.labelRotation, fontStyle)

	var width, height int
	if isVertical {
		width = textMaxWidth + tickLength<<1
		height = top.Height()
	} else {
		width = top.Width()
		height = tickLength<<1 + textMaxHeight
	}
	padding := Box{IsSet: true}
	switch opt.position {
	case PositionTop:
		padding.Top = top.Height() - height
	case PositionLeft:
		padding.Right = top.Width() - width
	case PositionRight:
		padding.Left = top.Width() - width
	default:
		padding.Top = top.Height() - defaultXAxisHeight
	}

	p := top.Child(PainterPaddingOption(padding))

	var x0, y0, x1, y1 int
	var ticksPaddingTop, ticksPaddingLeft int
	var labelPaddingTop, labelPaddingLeft, labelPaddingRight int
	var textAlign string

	switch opt.position {
	case PositionTop:
		if opt.labelRotation != 0 {
			flatWidth, flatHeight := top.measureTextMaxWidthHeight(opt.labels, 0, fontStyle)
			labelPaddingTop = flatHeight - textRotationHeightAdjustment(flatWidth, flatHeight, opt.labelRotation)
		} else {
			labelPaddingTop = 0
		}
		x1 = p.Width()
		y0 = labelMargin + int(fontStyle.FontSize)
		ticksPaddingTop = int(fontStyle.FontSize)
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
	default:
		if opt.labelRotation != 0 {
			flatWidth, flatHeight := top.measureTextMaxWidthHeight(opt.labels, 0, fontStyle)
			labelPaddingTop = tickLength<<1 + (textMaxHeight - textRotationHeightAdjustment(flatWidth, flatHeight, opt.labelRotation))
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
	} else if dataCount > 1 && a.p.Width()/dataCount <= boundaryGapDefaultThreshold {
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

	p.Child(PainterPaddingOption(Box{
		Left:  labelPaddingLeft,
		Top:   labelPaddingTop,
		Right: labelPaddingRight,
		IsSet: true,
	})).multiText(multiTextOption{
		firstIndex:     opt.dataStartIndex,
		align:          textAlign,
		textList:       opt.labels,
		fontStyle:      fontStyle,
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

	return Box{
		Bottom: height,
		Right:  width,
		IsSet:  true,
	}, nil
}
