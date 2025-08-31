package charts

import (
	"fmt"
)

type legendPainter struct {
	p   *Painter
	opt *LegendOption
}

// LegendOption defines the configuration for rendering the chart legend.
type LegendOption struct {
	// Show specifies if the legend should be rendered. Set to *false (via Ptr(false)) to hide the legend.
	Show *bool
	// Theme specifies the colors used for the legend.
	Theme ColorPalette
	// SeriesNames provides the text labels for the data series.
	SeriesNames []string
	// FontStyle specifies the font, size, and style for legend text.
	FontStyle FontStyle
	// Padding specifies the space around the legend.
	Padding Box
	// Offset specifies the position of the legend relative to the left and top sides.
	Offset OffsetStr
	// Align is the legend marker and text alignment: 'left', 'right', or 'center'. Default is 'left'.
	Align string
	// Vertical when set to *true makes the legend orientation vertical.
	Vertical *bool
	// Symbol defines the icon shape next to each label. Options: 'square', 'dot', 'diamond', 'circle'.
	Symbol Symbol // TODO - should Symbol configuration be changed now that we support per-series symbols
	// OverlayChart when set to *true renders the legend over the chart. Ignored if Vertical is true (Vertical always forces overlay).
	OverlayChart *bool
	// BorderWidth can be set to a non-zero value to render a box around the legend.
	BorderWidth float64
	// seriesSymbols provides custom symbols for each series.
	seriesSymbols []Symbol
}

// IsEmpty checks if the legend is empty.
func (opt *LegendOption) IsEmpty() bool {
	for _, v := range opt.SeriesNames {
		if v != "" {
			return false
		}
	}
	return true
}

// newLegendPainter returns a legend renderer.
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

	theme := opt.Theme
	if theme == nil {
		theme = getPreferredTheme(l.p.theme)
	}
	fontStyle := fillFontStyleDefaults(opt.FontStyle, defaultFontSize, theme.GetLegendTextColor(), l.p.font)
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
		left = (p.Width() - width) >> 1
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

	lastIndex := len(opt.SeriesNames) - 1
	for index, text := range opt.SeriesNames {
		seriesSymbol := opt.Symbol
		if index < len(opt.seriesSymbols) && opt.seriesSymbols[index] != "" {
			seriesSymbol = opt.seriesSymbols[index]
		}
		var drawIcon func(top, left int) int
		switch seriesSymbol {
		case SymbolSquare:
			drawIcon = func(top, left int) int {
				color := theme.GetSeriesColor(index)

				p.FilledRect(left, top-legendHeight+8, left+legendWidth, top+1, color, color, 0)
				return left + legendWidth
			}
		case SymbolDiamond:
			drawIcon = func(top, left int) int {
				color := theme.GetSeriesColor(index)

				p.FilledDiamond(left+5, top-5, 15, 20, color, color, 0)
				return left + legendHeight
			}
		case SymbolNone:
			drawIcon = func(top, left int) int { return left }
		case symbolCandlestick:
			drawIcon = func(top, left int) int {
				upColor, downColor := theme.GetSeriesUpDownColors(index)

				rectTop := top - legendHeight + 8
				rectBottom := top + 1
				midX := left + (legendWidth / 2)

				// Draw up arrow (left half) - triangle pointing up
				upArrowTip := left + (legendWidth / 4)
				p.moveTo(left, rectBottom)    // bottom left
				p.lineTo(midX, rectBottom)    // bottom right of up triangle
				p.lineTo(upArrowTip, rectTop) // tip of up arrow
				p.lineTo(left, rectBottom)    // back to bottom left
				p.fillStroke(upColor, upColor, 0)

				// Draw down arrow (right half) - triangle pointing down
				downArrowTip := left + (3 * legendWidth / 4)
				p.moveTo(midX, rectTop)             // top left of down triangle
				p.lineTo(left+legendWidth, rectTop) // top right
				p.lineTo(downArrowTip, rectBottom)  // tip of down arrow
				p.lineTo(midX, rectTop)             // back to top left
				p.fillStroke(downColor, downColor, 0)

				return left + legendWidth
			}
		default:
			centerColor := ColorTransparent
			if seriesSymbol == SymbolCircle {
				centerColor = opt.Theme.GetBackgroundColor()
			}
			drawIcon = func(top, left int) int {
				color := theme.GetSeriesColor(index)

				p.legendLineDot(Box{
					Top:    top + 1,
					Left:   left,
					Right:  left + legendWidth,
					Bottom: top + legendHeight + 1,
					IsSet:  true,
				}, color, 3, color, centerColor)
				return left + legendWidth
			}
		}

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
			x0 = drawIcon(y0, x0)
			x0 += textOffset
		}
		p.Text(text, x0, y0, 0, fontStyle)
		x0 += measureList[index].Width()
		if opt.Align == AlignRight {
			x0 += textOffset
			x0 = drawIcon(y0, x0)
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
		p.LineStroke(boxPoints, theme.GetLegendBorderColor(), opt.BorderWidth)
	}

	return result, nil
}
