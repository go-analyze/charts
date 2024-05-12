package charts

import (
	"math"
	"strings"

	"github.com/golang/freetype/truetype"
)

type axisPainter struct {
	p   *Painter
	opt *AxisOption
}

func NewAxisPainter(p *Painter, opt AxisOption) *axisPainter {
	return &axisPainter{
		p:   p,
		opt: &opt,
	}
}

type AxisOption struct {
	// Show specifies if the axis should be rendered, set this to *false (through False()) to hide the axis.
	Show *bool
	// Theme specifies the colors used for the axis.
	Theme ColorPalette
	// Data provides labels for the axis.
	Data []string
	// DataStartIndex specifies what index the Data values should start from.
	DataStartIndex int
	// Formatter for replacing axis text values.
	Formatter string
	// Position describes the position of axis, it can be 'left', 'top', 'right' or 'bottom'.
	Position string
	// BoundaryGap specifies that the chart should have additional space on the left and right, with data points being
	// centered between two axis ticks.  Enabled by default, specify *false (through False()) to change the spacing.
	BoundaryGap *bool
	// StrokeWidth is the axis line width.
	StrokeWidth float64
	// TickLength is the length of each axis tick.
	TickLength int
	// LabelMargin specifies the margin value of each label.
	LabelMargin int
	// FontSize specifies the font size of each label.
	FontSize float64
	// Font is the font used to render each label.
	Font *truetype.Font
	// FontColor is the color used for text rendered.
	FontColor Color
	// SplitLineShow, set this to true will show axis split line.
	SplitLineShow bool
	// TextRotation are the radians for rotating the label.
	TextRotation float64
	// LabelOffset is the offset of each label.
	LabelOffset Offset
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
	theme := opt.Theme
	if theme == nil {
		theme = top.theme
	}

	strokeWidth := opt.StrokeWidth
	if strokeWidth == 0 {
		strokeWidth = 1
	}

	font := getPreferredFont(opt.Font, a.p.font)
	fontColor := opt.FontColor
	if fontColor.IsZero() {
		fontColor = theme.GetTextColor()
	}
	fontSize := opt.FontSize
	if fontSize == 0 {
		fontSize = defaultFontSize
	}

	formatter := opt.Formatter
	if len(formatter) != 0 {
		for index, text := range opt.Data {
			opt.Data[index] = strings.ReplaceAll(formatter, "{value}", text)
		}
	}
	dataCount := len(opt.Data)

	centerLabels := !flagIs(false, opt.BoundaryGap)
	isVertical := opt.Position == PositionLeft || opt.Position == PositionRight

	// if less than zero, it means not processing
	tickLength := getDefaultInt(opt.TickLength, 5)
	labelMargin := getDefaultInt(opt.LabelMargin, 5)

	style := Style{
		StrokeColor: theme.GetAxisStrokeColor(),
		StrokeWidth: strokeWidth,
		Font:        font,
		FontColor:   fontColor,
		FontSize:    fontSize,
	}
	top.SetDrawingStyle(style).OverrideTextStyle(style)

	isTextRotation := opt.TextRotation != 0

	if isTextRotation {
		top.SetTextRotation(opt.TextRotation)
	}
	textMaxWidth, textMaxHeight := top.MeasureTextMaxWidthHeight(opt.Data)
	if isTextRotation {
		top.ClearTextRotation()
	}

	width := 0
	height := 0
	if isVertical {
		width = textMaxWidth + tickLength<<1
		height = top.Height()
	} else {
		width = top.Width()
		height = tickLength<<1 + textMaxHeight
	}
	padding := Box{}
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

	x0 := 0
	y0 := 0
	x1 := 0
	y1 := 0
	ticksPaddingTop := 0
	ticksPaddingLeft := 0
	labelPaddingTop := 0
	labelPaddingLeft := 0
	labelPaddingRight := 0
	orient := ""
	textAlign := ""

	switch opt.Position {
	case PositionTop:
		labelPaddingTop = 0
		x1 = p.Width()
		y0 = labelMargin + int(opt.FontSize)
		ticksPaddingTop = int(opt.FontSize)
		y1 = y0
		orient = OrientHorizontal
	case PositionLeft:
		x0 = p.Width()
		y0 = 0
		x1 = p.Width()
		y1 = p.Height()
		orient = OrientVertical
		textAlign = AlignRight
		ticksPaddingLeft = textMaxWidth + tickLength
		labelPaddingRight = width - textMaxWidth
	case PositionRight:
		orient = OrientVertical
		y1 = p.Height()
		labelPaddingLeft = width - textMaxWidth
	default:
		labelPaddingTop = height
		x1 = p.Width()
		orient = OrientHorizontal
	}

	labelCount := opt.LabelCount
	if labelCount <= 0 {
		var maxLabelCount int
		// Add 10px and remove one for some minimal extra padding so that letters don't collide
		if orient == OrientVertical {
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
	tickSpaces := dataCount
	if !centerLabels {
		// there is always one more tick than data sample, and if we are centering labels we use that extra tick to
		// center the label against, if not centering then we need one less tick spacing
		// passing the tickSpaces reduces the need to copy the logic from painter.go:MultiText
		tickSpaces--
	}

	if strokeWidth > 0 {
		p.Child(PainterPaddingOption(Box{
			Top:  ticksPaddingTop,
			Left: ticksPaddingLeft,
		})).Ticks(TicksOption{
			LabelCount: labelCount,
			TickSpaces: tickSpaces,
			Length:     tickLength,
			Orient:     orient,
			First:      opt.DataStartIndex,
		})
		p.LineStroke([]Point{
			{X: x0, Y: y0},
			{X: x1, Y: y1},
		})
	}

	p.Child(PainterPaddingOption(Box{
		Left:  labelPaddingLeft,
		Top:   labelPaddingTop,
		Right: labelPaddingRight,
	})).MultiText(MultiTextOption{
		First:          opt.DataStartIndex,
		Align:          textAlign,
		TextList:       opt.Data,
		Orient:         orient,
		LabelCount:     labelCount,
		LabelSkipCount: opt.LabelSkipCount,
		CenterLabels:   centerLabels,
		TextRotation:   opt.TextRotation,
		Offset:         opt.LabelOffset,
	})

	if opt.SplitLineShow { // show auxiliary lines
		style.StrokeColor = theme.GetAxisSplitLineColor()
		style.StrokeWidth = 1
		top.OverrideDrawingStyle(style)
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
				})
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
				})
			}
		}
	}

	return Box{Bottom: height, Right: width}, nil
}
