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
	// Show specifies if the axis should be rendered, set this to *false (through False()) to hide the axis.
	Show *bool
	// Theme specifies the colors used for the axis.
	Theme ColorPalette
	// Labels provides labels for each value on the axis.
	Labels []string
	// DataStartIndex specifies what index the Data values should start from.
	DataStartIndex int
	// Formatter for replacing axis text values.
	Formatter string
	// Position describes the position of axis, it can be 'left', 'top', 'right' or 'bottom'.
	Position string
	// BoundaryGap specifies that the chart should have additional space on the left and right, with data points being
	// centered between two axis ticks. Default is set based on the dataset density / size to produce an easy-to-read
	// graph. Specify a *bool (through charts.False() or charts.True()) to enforce a spacing.
	BoundaryGap *bool
	// StrokeWidth is the axis line width.
	StrokeWidth float64
	// TickLength is the length of each axis tick.
	TickLength int
	// LabelMargin specifies the margin value of each label.
	LabelMargin int
	// FontStyle specifies the font configuration for each label.
	FontStyle FontStyle
	// SplitLineShow, set this to true will show axis split line.
	SplitLineShow bool
	// LabelRotation are the radians for rotating the label.
	LabelRotation float64
	// LabelOffset is the offset of each label.
	LabelOffset OffsetInt
	// Unit is a suggestion for how large the axis step is, this is a recommendation only. Larger numbers result in fewer labels.
	Unit float64
	// LabelCount is the number of labels to show on the axis. Specify a smaller number to reduce writing collisions. This value takes priority over Unit.
	LabelCount int
	// LabelCountAdjustment specifies a relative influence on how many labels should be rendered.
	// Typically, this is negative to result in cleaner graphs, positive values may result in text collisions.
	LabelCountAdjustment int
	// LabelSkipCount specifies a number of lines between labels where there will be no label,
	// but a horizontal line will still be drawn.
	LabelSkipCount int
}

func (a *axisPainter) Render() (Box, error) {
	opt := a.opt
	if flagIs(false, opt.Show) {
		return BoxZero, nil
	}
	top := a.p
	theme := getPreferredTheme(opt.Theme, top.theme)

	strokeWidth := opt.StrokeWidth
	if strokeWidth == 0 {
		strokeWidth = 1
	}

	fontStyle := opt.FontStyle
	fontStyle.Font = getPreferredFont(fontStyle.Font, a.p.font)
	if fontStyle.FontColor.IsZero() {
		fontStyle.FontColor = theme.GetTextColor()
	}
	if fontStyle.FontSize == 0 {
		fontStyle.FontSize = defaultFontSize
	}

	formatter := opt.Formatter
	if formatter != "" {
		for index, text := range opt.Labels {
			opt.Labels[index] = strings.ReplaceAll(formatter, "{value}", text)
		}
	}

	isVertical := opt.Position == PositionLeft || opt.Position == PositionRight

	// if less than zero, it means not processing
	tickLength := getDefaultInt(opt.TickLength, 5)
	labelMargin := getDefaultInt(opt.LabelMargin, 5)

	textMaxWidth, textMaxHeight := top.measureTextMaxWidthHeight(opt.Labels, opt.LabelRotation, fontStyle)

	var width, height int
	if isVertical {
		width = textMaxWidth + tickLength<<1
		height = top.Height()
	} else {
		width = top.Width()
		height = tickLength<<1 + textMaxHeight
	}
	padding := Box{IsSet: true}
	switch opt.Position {
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

	switch opt.Position {
	case PositionTop:
		if opt.LabelRotation != 0 {
			flatWidth, flatHeight := top.measureTextMaxWidthHeight(opt.Labels, 0, fontStyle)
			labelPaddingTop = flatHeight - textRotationHeightAdjustment(flatWidth, flatHeight, opt.LabelRotation)
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
		if opt.LabelRotation != 0 {
			flatWidth, flatHeight := top.measureTextMaxWidthHeight(opt.Labels, 0, fontStyle)
			labelPaddingTop = tickLength<<1 + (textMaxHeight - textRotationHeightAdjustment(flatWidth, flatHeight, opt.LabelRotation))
		} else {
			labelPaddingTop = height
		}
		x1 = p.Width()
	}

	dataCount := len(opt.Labels)
	labelCount := opt.LabelCount
	if labelCount <= 0 {
		var maxLabelCount int
		// Add 10px and remove one for some minimal extra padding so that letters don't collide
		if isVertical {
			maxLabelCount = (top.Height() / (textMaxHeight + 10)) - 1
		} else {
			maxLabelCount = (top.Width() / (textMaxWidth + 10)) - 1
		}
		if opt.Unit > 0 {
			multiplier := 1.0
			for {
				labelCount = int(math.Ceil(float64(dataCount) / (opt.Unit * multiplier)))
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
	labelCount += opt.LabelCountAdjustment
	if labelCount < 2 {
		labelCount = 2
	}

	centerLabels := true
	if opt.BoundaryGap != nil {
		centerLabels = *opt.BoundaryGap
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
		strokeColor := theme.GetAxisStrokeColor()
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
			firstIndex:  opt.DataStartIndex,
			strokeWidth: strokeWidth,
			strokeColor: strokeColor,
		})
		p.LineStroke([]Point{
			{X: x0, Y: y0},
			{X: x1, Y: y1},
		}, strokeColor, strokeWidth)
	}

	p.Child(PainterPaddingOption(Box{
		Left:  labelPaddingLeft,
		Top:   labelPaddingTop,
		Right: labelPaddingRight,
		IsSet: true,
	})).multiText(multiTextOption{
		firstIndex:     opt.DataStartIndex,
		align:          textAlign,
		textList:       opt.Labels,
		fontStyle:      fontStyle,
		vertical:       isVertical,
		labelCount:     labelCount,
		tickCount:      tickCount,
		labelSkipCount: opt.LabelSkipCount,
		centerLabels:   centerLabels,
		textRotation:   opt.LabelRotation,
		offset:         opt.LabelOffset,
	})

	if opt.SplitLineShow { // show auxiliary lines
		if isVertical {
			x0 := p.Width()
			x1 := top.Width()
			if opt.Position == PositionRight {
				x0 = 0
				x1 = top.Width() - p.Width()
			}
			yValues := autoDivide(height, tickSpaces)
			yValues = yValues[0 : len(yValues)-1]
			for _, y := range yValues {
				top.LineStroke([]Point{
					{X: x0, Y: y},
					{X: x1, Y: y},
				}, theme.GetAxisSplitLineColor(), 1)
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
				}, theme.GetAxisSplitLineColor(), 1)
			}
		}
	}

	return Box{
		Bottom: height,
		Right:  width,
		IsSet:  true,
	}, nil
}
