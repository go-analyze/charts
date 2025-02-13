package charts

import (
	"fmt"
)

type legendPainter struct {
	p   *Painter
	opt *LegendOption
}

const IconRect = "rect"
const IconDot = "dot"

type LegendOption struct {
	// Show specifies if the legend should be rendered, set this to *false (through False()) to hide the legend.
	Show *bool
	// Theme specifies the colors used for the legend.
	Theme ColorPalette
	// Deprecated: Data is deprecated, use SeriesNames instead.
	Data []string
	// SeriesNames provides text labels for the legend.
	SeriesNames []string
	// FontStyle specifies the font, size, and style for rendering the legend.
	FontStyle FontStyle
	// Padding specifies space padding around the legend.
	Padding Box
	// Offset allows you to specify the position of the legend component relative to the left and top side.
	Offset OffsetStr
	// Align is the legend marker and text alignment, it can be 'left', 'right' or 'center', default is 'left'.
	Align string
	// Vertical can be set to *true to set the legend orientation to be vertical.
	Vertical *bool
	// Icon to show next to the labels.	Can be 'rect' or 'dot'.
	Icon string
	// OverlayChart can be set to *true to render the legend over the chart. Ignored if Vertical is set to true (always overlapped).
	OverlayChart *bool
	// BorderWidth can be set to a non-zero value to render a box around the legend.
	BorderWidth float64
}

// IsEmpty checks legend is empty
func (opt *LegendOption) IsEmpty() bool {
	if len(opt.SeriesNames) == 0 {
		opt.SeriesNames = opt.Data
	}
	for _, v := range opt.SeriesNames {
		if v != "" {
			return false
		}
	}
	return true
}

// newLegendPainter returns a legend renderer
func newLegendPainter(p *Painter, opt LegendOption) *legendPainter {
	return &legendPainter{
		p:   p,
		opt: &opt,
	}
}

func (l *legendPainter) Render() (Box, error) {
	opt := l.opt
	if opt.IsEmpty() || flagIs(false, opt.Show) {
		return BoxZero, nil
	}

	if len(opt.SeriesNames) == 0 {
		opt.SeriesNames = opt.Data
	}
	theme := opt.Theme
	if theme == nil {
		theme = getPreferredTheme(l.p.theme)
	}
	fontStyle := opt.FontStyle
	if fontStyle.FontSize == 0 {
		fontStyle.FontSize = defaultFontSize
	}
	if fontStyle.FontColor.IsZero() {
		fontStyle.FontColor = theme.GetTextColor()
	}
	vertical := flagIs(true, opt.Vertical)
	offset := opt.Offset
	if offset.Left == "" {
		if vertical {
			// in the vertical orientation it's more visually appealing to default to the right side or left side
			if opt.Align != "" {
				offset.Left = opt.Align
			} else {
				offset.Left = PositionLeft
			}
		} else {
			offset.Left = PositionCenter
		}
	}
	padding := opt.Padding
	if padding.IsZero() {
		padding.Top = 5
	}
	p := l.p.Child(PainterPaddingOption(padding))

	// calculate the width and height of the display
	measureList := make([]Box, len(opt.SeriesNames))
	var width, height int
	const builtInSpacing = 20
	const textOffset = 2
	const legendWidth = 30
	const legendHeight = 20
	var maxTextWidth, itemMaxHeight int
	for index, text := range opt.SeriesNames {
		b := p.MeasureText(text, 0, fontStyle)
		if b.Width() > maxTextWidth {
			maxTextWidth = b.Width()
		}
		if b.Height() > itemMaxHeight {
			itemMaxHeight = b.Height()
		}
		if flagIs(true, opt.Vertical) {
			height += b.Height()
		} else {
			width += b.Width()
		}
		measureList[index] = b
	}

	// add padding
	if vertical {
		width = maxTextWidth + textOffset + legendWidth
		height = builtInSpacing * len(opt.SeriesNames)
	} else {
		height = legendHeight
		offsetValue := (len(opt.SeriesNames) - 1) * (builtInSpacing + textOffset)
		allLegendWidth := len(opt.SeriesNames) * legendWidth
		width += offsetValue + allLegendWidth
	}

	// calculate starting position
	var left int
	switch offset.Left {
	case PositionLeft:
		// leave default of zero
	case PositionRight:
		left = p.Width() - width
	case PositionCenter:
		left = p.Width()>>1 - (width >> 1)
	default:
		if v, err := parseFlexibleValue(offset.Left, float64(p.Width())); err != nil {
			return BoxZero, fmt.Errorf("error parsing legend position: %w", err)
		} else {
			left = int(v)
		}
	}
	if left < 0 {
		left = 0
	}

	var top int
	switch offset.Top {
	case "", PositionTop:
		// leave default of zero
	case PositionBottom:
		top = p.Height() - height
	default:
		if v, err := parseFlexibleValue(offset.Top, float64(p.Height())); err != nil {
			return BoxZero, fmt.Errorf("error parsing legend position: %w", err)
		} else {
			top = int(v)
		}
	}

	y := top + 10
	x0 := left
	y0 := y

	var drawIcon func(top, left int, color Color) int
	if opt.Icon == IconRect {
		drawIcon = func(top, left int, color Color) int {
			p.FilledRect(left, top-legendHeight+8, left+legendWidth, top+1, color, color, 0)
			return left + legendWidth
		}
	} else {
		drawIcon = func(top, left int, color Color) int {
			p.legendLineDot(Box{
				Top:    top + 1,
				Left:   left,
				Right:  left + legendWidth,
				Bottom: top + legendHeight + 1,
				IsSet:  true,
			}, color, 3, color)
			return left + legendWidth
		}
	}

	lastIndex := len(opt.SeriesNames) - 1
	for index, text := range opt.SeriesNames {
		color := theme.GetSeriesColor(index)
		if vertical {
			if opt.Align == AlignRight {
				// adjust x0 so that the text will start with a right alignment to the longest line
				x0 += maxTextWidth - measureList[index].Width()
			}
		} else {
			// check if item will overrun the right side boundary
			itemWidth := x0 + measureList[index].Width() + textOffset + builtInSpacing + legendWidth
			if lastIndex == index {
				itemWidth = x0 + measureList[index].Width() + legendWidth
			}
			if itemWidth > p.Width() {
				newLineStart := left
				if opt.Align == AlignCenter {
					// recalculate width and center based off remaining width
					var remainingWidth int
					for i2 := index; i2 < len(opt.SeriesNames); i2++ {
						b := p.MeasureText(opt.SeriesNames[i2], 0, fontStyle)
						remainingWidth += b.Width()
					}
					remainingCount := len(opt.SeriesNames) - index
					remainingWidth += remainingCount * legendWidth
					remainingWidth += (remainingCount - 1) * (builtInSpacing + textOffset)

					newLineStart = left + ((p.Width() >> 1) - (remainingWidth >> 1))
					if newLineStart < 0 {
						newLineStart = 0
					}
				}
				x0 = newLineStart
				y += itemMaxHeight
				y0 = y
			}
		}

		if opt.Align != AlignRight {
			x0 = drawIcon(y0, x0, color)
			x0 += textOffset
		}
		p.Text(text, x0, y0, 0, fontStyle)
		x0 += measureList[index].Width()
		if opt.Align == AlignRight {
			x0 += textOffset
			x0 = drawIcon(y0, x0, color)
		}

		if vertical {
			y0 += builtInSpacing
			x0 = left
		} else {
			x0 += builtInSpacing
			y0 = y
		}
	}

	bottom := y0 + padding.Bottom - 10
	if !vertical {
		bottom += itemMaxHeight
	}

	result := Box{
		Top:    top - padding.Top,
		Bottom: bottom,
		Left:   left - padding.Left,
		Right:  left + width + padding.Right,
		IsSet:  true,
	}

	if opt.BorderWidth > 0 {
		// TODO - if drawn over the chart this can look awkward, we should try to draw this first
		boxPad := 10 // built in adjustment for possible measure vs render variations
		boxPoints := []Point{
			{X: result.Left - boxPad, Y: result.Bottom + boxPad},
			{X: result.Left - boxPad, Y: result.Top - boxPad},
			{X: result.Left + result.Width() + boxPad, Y: result.Top - boxPad},
			{X: result.Left + result.Width() + boxPad, Y: result.Bottom + boxPad},
			{X: result.Left - boxPad, Y: result.Bottom + boxPad},
		}
		// TODO - allow color to be configured via theme or configuration
		p.LineStroke(boxPoints, ColorBlack, opt.BorderWidth)
	}

	return result, nil
}
