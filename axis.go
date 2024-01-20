package charts

import (
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/wcharczuk/go-chart/v2"
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
	// The theme of chart
	Theme ColorPalette
	// Formatter for y axis text value
	Formatter string
	// The label of axis
	Data []string
	// The boundary gap on both sides of a coordinate axis.
	// Nil or *true means the center part of two axis ticks
	BoundaryGap *bool
	// The flag for show axis, set this to *false will hide axis
	Show *bool
	// The position of axis, it can be 'left', 'top', 'right' or 'bottom'
	Position string
	// Number of segments that the axis is split into. Note that this number serves only as a recommendation.
	SplitNumber int
	// The line color of axis
	StrokeColor Color
	// The line width
	StrokeWidth float64
	// The length of the axis tick
	TickLength int
	// The first axis
	FirstAxis int
	// The margin value of label
	LabelMargin int
	// The font size of label
	FontSize float64
	// The font of label
	Font *truetype.Font
	// The color of label
	FontColor Color
	// The flag for show axis split line, set this to true will show axis split line
	SplitLineShow bool
	// The color of split line
	SplitLineColor Color
	// The text rotation of label
	TextRotation float64
	// The offset of label
	LabelOffset Box
	Unit        int
}

func (a *axisPainter) Render() (Box, error) {
	opt := a.opt
	top := a.p
	theme := opt.Theme
	if theme == nil {
		theme = top.theme
	}
	if isFalse(opt.Show) {
		return BoxZero, nil
	}

	strokeWidth := opt.StrokeWidth
	if strokeWidth == 0 {
		strokeWidth = 1
	}

	font := opt.Font
	if font == nil {
		font = a.p.font
	}
	if font == nil {
		font = theme.GetFont()
	}
	fontColor := opt.FontColor
	if fontColor.IsZero() {
		fontColor = theme.GetTextColor()
	}
	fontSize := opt.FontSize
	if fontSize == 0 {
		fontSize = theme.GetFontSize()
	}
	strokeColor := opt.StrokeColor
	if strokeColor.IsZero() {
		strokeColor = theme.GetAxisStrokeColor()
	}

	data := opt.Data
	formatter := opt.Formatter
	if len(formatter) != 0 {
		for index, text := range data {
			data[index] = strings.ReplaceAll(formatter, "{value}", text)
		}
	}
	dataCount := len(data)
	tickCount := dataCount

	boundaryGap := true
	if isFalse(opt.BoundaryGap) {
		boundaryGap = false
	}
	isVertical := opt.Position == PositionLeft ||
		opt.Position == PositionRight

	labelPosition := ""
	if !boundaryGap {
		tickCount--
		labelPosition = PositionLeft
	}
	if isVertical && boundaryGap {
		labelPosition = PositionCenter
	}

	// if less than zero, it means not processing
	tickLength := getDefaultInt(opt.TickLength, 5)
	labelMargin := getDefaultInt(opt.LabelMargin, 5)

	style := Style{
		StrokeColor: strokeColor,
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
	textMaxWidth, textMaxHeight := top.MeasureTextMaxWidthHeight(data)
	if isTextRotation {
		top.ClearTextRotation()
	}

	// Add 30px to calculate text display area
	textFillWidth := float64(textMaxWidth + 20)
	// Calculate more suitable display item based on text width
	fitTextCount := ceilFloatToInt(float64(top.Width()) / textFillWidth)

	unit := opt.Unit
	if unit <= 0 {
		unit = ceilFloatToInt(float64(dataCount) / float64(fitTextCount))
		unit = chart.MaxInt(unit, opt.SplitNumber)
		// even number
		if unit%2 == 0 && dataCount%(unit+1) == 0 {
			unit++
		}
	}

	width := 0
	height := 0
	// vertical
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

	if strokeWidth > 0 {
		p.Child(PainterPaddingOption(Box{
			Top:  ticksPaddingTop,
			Left: ticksPaddingLeft,
		})).Ticks(TicksOption{
			Count:  tickCount,
			Length: tickLength,
			Unit:   unit,
			Orient: orient,
			First:  opt.FirstAxis,
		})
		p.LineStroke([]Point{
			{
				X: x0,
				Y: y0,
			},
			{
				X: x1,
				Y: y1,
			},
		})
	}

	p.Child(PainterPaddingOption(Box{
		Left:  labelPaddingLeft,
		Top:   labelPaddingTop,
		Right: labelPaddingRight,
	})).MultiText(MultiTextOption{
		First:        opt.FirstAxis,
		Align:        textAlign,
		TextList:     data,
		Orient:       orient,
		Unit:         unit,
		Position:     labelPosition,
		TextRotation: opt.TextRotation,
		Offset:       opt.LabelOffset,
	})
	// show auxilary lines
	if opt.SplitLineShow {
		style.StrokeColor = opt.SplitLineColor
		style.StrokeWidth = 1
		top.OverrideDrawingStyle(style)
		if isVertical {
			x0 := p.Width()
			x1 := top.Width()
			if opt.Position == PositionRight {
				x0 = 0
				x1 = top.Width() - p.Width()
			}
			yValues := autoDivide(height, tickCount)
			yValues = yValues[0 : len(yValues)-1]
			for _, y := range yValues {
				top.LineStroke([]Point{
					{
						X: x0,
						Y: y,
					},
					{
						X: x1,
						Y: y,
					},
				})
			}
		} else {
			y0 := p.Height() - defaultXAxisHeight
			y1 := top.Height() - defaultXAxisHeight
			for index, x := range autoDivide(width, tickCount) {
				if index == 0 {
					continue
				}
				top.LineStroke([]Point{
					{
						X: x,
						Y: y0,
					},
					{
						X: x,
						Y: y1,
					},
				})
			}
		}
	}

	return Box{
		Bottom: height,
		Right:  width,
	}, nil
}
