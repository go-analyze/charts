package charts

import (
	"fmt"

	"github.com/go-analyze/charts/chartdraw"
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
	// Data provides text for the legend.
	Data []string
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
}

// IsEmpty checks legend is empty
func (opt *LegendOption) IsEmpty() bool {
	for _, v := range opt.Data {
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
	p.SetFontStyle(fontStyle)

	// calculate the width and height of the display
	measureList := make([]Box, len(opt.Data))
	width := 0
	height := 0
	const builtInSpacing = 20
	const textOffset = 2
	const legendWidth = 30
	const legendHeight = 20
	maxTextWidth := 0
	itemMaxHeight := 0
	for index, text := range opt.Data {
		b := p.MeasureText(text)
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
		height = builtInSpacing * len(opt.Data)
	} else {
		height = legendHeight
		offsetValue := (len(opt.Data) - 1) * (builtInSpacing + textOffset)
		allLegendWidth := len(opt.Data) * legendWidth
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
			return BoxZero, fmt.Errorf("error parsing legend position: %v", err)
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
			return BoxZero, fmt.Errorf("error parsing legend position: %v", err)
		} else {
			top = int(v)
		}
	}

	startX := left
	y := top + 10
	startY := y
	x0 := startX
	y0 := y

	var drawIcon func(top, left int) int
	if opt.Icon == IconRect {
		drawIcon = func(top, left int) int {
			p.Rect(Box{
				Top:    top - legendHeight + 8,
				Left:   left,
				Right:  left + legendWidth,
				Bottom: top + 1,
				IsSet:  true,
			})
			return left + legendWidth
		}
	} else {
		drawIcon = func(top, left int) int {
			p.LegendLineDot(Box{
				Top:    top + 1,
				Left:   left,
				Right:  left + legendWidth,
				Bottom: top + legendHeight + 1,
				IsSet:  true,
			})
			return left + legendWidth
		}
	}

	lastIndex := len(opt.Data) - 1
	for index, text := range opt.Data {
		color := theme.GetSeriesColor(index)
		p.SetDrawingStyle(chartdraw.Style{
			FillColor:   color,
			StrokeColor: color,
		})
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
				newLineStart := startX
				if opt.Align == AlignCenter {
					// recalculate width and center based off remaining width
					var remainingWidth int
					for i2 := index; i2 < len(opt.Data); i2++ {
						b := p.MeasureText(opt.Data[i2])
						remainingWidth += b.Width()
					}
					remainingCount := len(opt.Data) - index
					remainingWidth += remainingCount * legendWidth
					remainingWidth += (remainingCount - 1) * (builtInSpacing + textOffset)

					newLineStart = startX + ((p.Width() >> 1) - (remainingWidth >> 1))
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
		p.Text(text, x0, y0)
		x0 += measureList[index].Width()
		if opt.Align == AlignRight {
			x0 += textOffset
			x0 = drawIcon(y0, x0)
		}

		if vertical {
			y0 += builtInSpacing
			x0 = startX
		} else {
			x0 += builtInSpacing
			y0 = y
		}
	}

	return Box{
		Top:    startY,
		Bottom: y0 + itemMaxHeight + padding.Bottom + padding.Top,
		Left:   startX,
		Right:  width,
		IsSet:  true,
	}, nil
}
