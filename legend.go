package charts

import (
	"fmt"
)

const (
	legendBuiltInSpacing    = 20
	legendTextOffset        = 2
	legendIconStandardWidth = 30
	legendIconHeight        = 20
)

// iconWidthForSymbol returns the width used by a legend icon for the given symbol.
func iconWidthForSymbol(symbol Symbol) int {
	switch symbol {
	case SymbolNone:
		return 0
	case SymbolDiamond:
		return 20
	case SymbolSquare, SymbolCircle, SymbolDot, symbolCandlestick:
		return legendIconStandardWidth
	default:
		return legendIconStandardWidth
	}
}

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
	// TODO - v0.6 - consider combining symbol with size into a SymbolStyle struct
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

// computeLayoutParams calculates common layout parameters for legend rendering.
func (l *legendPainter) computeLayoutParams() (
	theme ColorPalette,
	fontStyle FontStyle,
	vertical bool,
	padding Box,
	p *Painter,
	measureList []Box,
	iconWidths []int,
	maxTextWidth, itemMaxHeight, width, left, top int,
	err error,
) {
	opt := l.opt
	theme = getPreferredTheme(opt.Theme, l.p.theme)
	fontStyle = fillFontStyleDefaults(opt.FontStyle, defaultFontSize, theme.GetLegendTextColor(), l.p.font)
	vertical = flagIs(true, opt.Vertical)

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

	padding = opt.Padding
	if padding.IsZero() {
		padding.Top = 5
	}
	p = l.p.Child(PainterPaddingOption(padding))

	// measure items and resolve per-series icon widths
	measureList = make([]Box, len(opt.SeriesNames))
	iconWidths = make([]int, len(opt.SeriesNames))
	var totalIconWidth, maxIconWidth int
	for index, text := range opt.SeriesNames {
		b := p.MeasureText(text, 0, fontStyle)
		if b.Width() > maxTextWidth {
			maxTextWidth = b.Width()
		}
		if b.Height() > itemMaxHeight {
			itemMaxHeight = b.Height()
		}
		if !vertical {
			width += b.Width()
		}
		measureList[index] = b
		// resolve symbol for this series
		symbol := opt.Symbol
		if index < len(opt.seriesSymbols) && opt.seriesSymbols[index] != "" {
			symbol = opt.seriesSymbols[index]
		}
		iconWidths[index] = iconWidthForSymbol(symbol)
		totalIconWidth += iconWidths[index]
		if iconWidths[index] > maxIconWidth {
			maxIconWidth = iconWidths[index]
		}
	}

	if vertical {
		width = maxTextWidth + legendTextOffset + maxIconWidth
	} else {
		offsetValue := (len(opt.SeriesNames) - 1) * (legendBuiltInSpacing + legendTextOffset)
		width += offsetValue + totalIconWidth
	}

	// calculate left position
	switch offset.Left {
	case PositionLeft:
		// leave default of zero
	case PositionRight:
		left = p.Width() - width
	case PositionCenter:
		left = (p.Width() - width) >> 1
	default:
		if v, parseErr := parseFlexibleValue(offset.Left, float64(p.Width())); parseErr != nil {
			err = fmt.Errorf("error parsing legend position: %w", parseErr)
			return
		} else {
			left = int(v)
		}
	}
	if left < 0 {
		left = 0
	}

	// calculate top position
	var height int
	if vertical {
		height = legendBuiltInSpacing * len(opt.SeriesNames)
	} else {
		height = legendIconHeight
	}
	switch offset.Top {
	case "", PositionTop:
		// leave default of zero
	case PositionBottom:
		top = p.Height() - height
	default:
		if v, parseErr := parseFlexibleValue(offset.Top, float64(p.Height())); parseErr != nil {
			err = fmt.Errorf("error parsing legend position: %w", parseErr)
			return
		} else {
			top = int(v)
		}
	}

	return
}

// iterateLegendLayout walks through legend item positions, calling onItem for each if provided.
// Returns the final y0 position for bounding box calculation.
func (l *legendPainter) iterateLegendLayout(
	p *Painter,
	measureList []Box,
	iconWidths []int,
	maxTextWidth, itemMaxHeight, left, top int,
	vertical bool,
	onItem func(index, x0, y0, iconWidth int),
) int {
	opt := l.opt

	y := top + 10
	x0 := left
	y0 := y

	lastIndex := len(opt.SeriesNames) - 1
	for index := range opt.SeriesNames {
		iconWidth := iconWidths[index]
		if vertical {
			if opt.Align == AlignRight {
				// adjust x0 so that the text will start with a right alignment to the longest line
				x0 += maxTextWidth - measureList[index].Width()
			}
		} else {
			// check if item will overrun the right side boundary
			itemWidth := x0 + measureList[index].Width() + legendTextOffset + legendBuiltInSpacing + iconWidth
			if lastIndex == index {
				itemWidth = x0 + measureList[index].Width() + iconWidth
			}
			if itemWidth > p.Width() {
				newLineStart := left
				if opt.Align == AlignCenter {
					// calculate remaining width using pre-measured values
					remainingCount := len(measureList) - index
					var remainingWidth int
					var remainingIconWidth int
					for i2 := index; i2 < len(measureList); i2++ {
						remainingWidth += measureList[i2].Width()
						remainingIconWidth += iconWidths[i2]
					}
					remainingWidth += remainingIconWidth + (remainingCount-1)*(legendBuiltInSpacing+legendTextOffset)
					newLineStart = left + (p.Width()-remainingWidth)>>1
					if newLineStart < 0 {
						newLineStart = 0
					}
				}
				x0 = newLineStart
				y += itemMaxHeight
				y0 = y
			}
		}

		if onItem != nil {
			onItem(index, x0, y0, iconWidth)
		}

		x0 += iconWidth + legendTextOffset + measureList[index].Width()

		// advance to next row/position
		if vertical {
			y0 += legendBuiltInSpacing
			x0 = left
		} else {
			x0 += legendBuiltInSpacing
			y0 = y
		}
	}

	return y0
}

// calculateBox returns the bounding box without rendering.
func (l *legendPainter) calculateBox() (Box, error) {
	opt := l.opt
	if opt.IsEmpty() || flagIs(false, opt.Show) {
		return BoxZero, nil
	}

	_, _, vertical, padding, p, measureList, iconWidths, maxTextWidth, itemMaxHeight, width, left, top, err := l.computeLayoutParams()
	if err != nil {
		return BoxZero, err
	}

	y0 := l.iterateLegendLayout(p, measureList, iconWidths, maxTextWidth, itemMaxHeight, left, top, vertical, nil)

	bottom := y0 + padding.Bottom - 10
	if !vertical {
		bottom += itemMaxHeight
	}

	return Box{
		Top:    top - padding.Top,
		Bottom: bottom,
		Left:   left - padding.Left,
		Right:  left + width + padding.Right,
		IsSet:  true,
	}, nil
}

func (l *legendPainter) Render() (Box, error) {
	opt := l.opt
	if opt.IsEmpty() || flagIs(false, opt.Show) {
		return BoxZero, nil
	}

	theme, fontStyle, vertical, padding, p, measureList, iconWidths, maxTextWidth, itemMaxHeight, width, left, top, err := l.computeLayoutParams()
	if err != nil {
		return BoxZero, err
	}

	// draw each legend item and capture final y0 for bounding box
	y0 := l.iterateLegendLayout(p, measureList, iconWidths, maxTextWidth, itemMaxHeight, left, top, vertical,
		func(index, x0, y0, iconWidth int) {
			text := opt.SeriesNames[index]
			seriesSymbol := opt.Symbol
			if index < len(opt.seriesSymbols) && opt.seriesSymbols[index] != "" {
				seriesSymbol = opt.seriesSymbols[index]
			}

			drawIcon := l.makeIconDrawer(p, theme, index, seriesSymbol)

			if opt.Align != AlignRight {
				drawIcon(y0, x0)
				x0 += iconWidth + legendTextOffset
			}
			p.Text(text, x0, y0, 0, fontStyle)
			if opt.Align == AlignRight {
				x0 += measureList[index].Width() + legendTextOffset
				drawIcon(y0, x0)
			}
		})
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

// makeIconDrawer returns a function that draws the legend icon for a series.
func (l *legendPainter) makeIconDrawer(p *Painter, theme ColorPalette, index int, symbol Symbol) func(top, left int) {
	switch symbol {
	case SymbolSquare:
		return func(top, left int) {
			color := theme.GetSeriesColor(index)
			p.FilledRect(left, top-legendIconHeight+8, left+legendIconStandardWidth, top+1, color, color, 0)
		}
	case SymbolDiamond:
		return func(top, left int) {
			color := theme.GetSeriesColor(index)
			p.FilledDiamond(left+5, top-5, 15, 20, color, color, 0)
		}
	case SymbolNone:
		return func(top, left int) {}
	case symbolCandlestick:
		return func(top, left int) {
			upColor, downColor := theme.GetSeriesUpDownColors(index)

			rectTop := top - legendIconHeight + 8
			rectBottom := top + 1
			midX := left + (legendIconStandardWidth / 2)

			// draw up arrow (left half)
			upArrowTip := left + (legendIconStandardWidth / 4)
			p.moveTo(left, rectBottom)
			p.lineTo(midX, rectBottom)
			p.lineTo(upArrowTip, rectTop)
			p.lineTo(left, rectBottom)
			p.fillStroke(upColor, upColor, 0)

			// draw down arrow (right half)
			downArrowTip := left + (3 * legendIconStandardWidth / 4)
			p.moveTo(midX, rectTop)
			p.lineTo(left+legendIconStandardWidth, rectTop)
			p.lineTo(downArrowTip, rectBottom)
			p.lineTo(midX, rectTop)
			p.fillStroke(downColor, downColor, 0)
		}
	default:
		centerColor := ColorTransparent
		if symbol == SymbolCircle {
			centerColor = theme.GetBackgroundColor()
		}
		return func(top, left int) {
			color := theme.GetSeriesColor(index)
			p.legendLineDot(Box{
				Top:    top + 1,
				Left:   left,
				Right:  left + legendIconStandardWidth,
				Bottom: top + legendIconHeight + 1,
				IsSet:  true,
			}, color, 3, color, centerColor)
		}
	}
}
