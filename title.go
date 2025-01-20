package charts

import (
	"fmt"
	"strings"
)

type TitleOption struct {
	// Show specifies if the title should be rendered, set this to *false (through False()) to hide the title.
	Show *bool
	// Theme specifies the colors used for the title.
	Theme ColorPalette
	// Text specifies the title text, supporting \n for new lines.
	Text string
	// Subtext to the title, supporting \n for new lines.
	Subtext string
	// Offset allows you to specify the position of the title component relative to the left and top side.
	Offset OffsetStr
	// FontStyle specifies the font, size, and style for rendering the title.
	FontStyle FontStyle
	// SubtextFontStyle specifies the font, size, and style for rendering the subtext.
	SubtextFontStyle FontStyle
}

type titleMeasureOption struct {
	width  int
	height int
	text   string
	style  FontStyle
}

func splitTitleText(text string) []string {
	arr := strings.Split(text, "\n")
	result := make([]string, 0)
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		result = append(result, v)
	}
	return result
}

type titlePainter struct {
	p   *Painter
	opt *TitleOption
}

// newTitlePainter returns a title renderer
func newTitlePainter(p *Painter, opt TitleOption) *titlePainter {
	return &titlePainter{
		p:   p,
		opt: &opt,
	}
}

func (t *titlePainter) Render() (Box, error) {
	opt := t.opt
	p := t.p
	if flagIs(false, opt.Show) {
		return BoxZero, nil
	}

	theme := opt.Theme
	if theme == nil {
		theme = getPreferredTheme(p.theme)
	}
	if opt.Text == "" && opt.Subtext == "" {
		return BoxZero, nil
	}

	measureOptions := make([]titleMeasureOption, 0)

	fontStyle := opt.FontStyle
	if fontStyle.Font == nil {
		fontStyle.Font = GetDefaultFont()
	}
	if fontStyle.FontColor.IsZero() {
		fontStyle.FontColor = theme.GetTextColor()
	}
	if fontStyle.FontSize == 0 {
		fontStyle.FontSize = defaultFontSize
	}
	subtextFontStyle := opt.SubtextFontStyle
	if subtextFontStyle.Font == nil {
		subtextFontStyle.Font = fontStyle.Font
	}
	if subtextFontStyle.FontColor.IsZero() {
		subtextFontStyle.FontColor = fontStyle.FontColor
	}
	if subtextFontStyle.FontSize == 0 {
		subtextFontStyle.FontSize = fontStyle.FontSize
	}

	// main title
	for _, v := range splitTitleText(opt.Text) {
		measureOptions = append(measureOptions, titleMeasureOption{
			text:  v,
			style: fontStyle,
		})
	}
	// subtitle
	for _, v := range splitTitleText(opt.Subtext) {
		measureOptions = append(measureOptions, titleMeasureOption{
			text:  v,
			style: subtextFontStyle,
		})
	}
	textMaxWidth := 0
	textMaxHeight := 0
	textTotalHeight := 0
	for index, item := range measureOptions {
		textBox := p.MeasureText(item.text, 0, item.style)

		w := textBox.Width()
		h := textBox.Height()
		if w > textMaxWidth {
			textMaxWidth = w
		}
		if h > textMaxHeight {
			textMaxHeight = h
		}
		textTotalHeight += h
		measureOptions[index].height = h
		measureOptions[index].width = w
	}
	width := textMaxWidth

	offset := opt.Offset
	titleX := 0
	switch offset.Left {
	case "", PositionLeft:
		// no-op
	case PositionRight:
		titleX = p.Width() - textMaxWidth
	case PositionCenter:
		titleX = p.Width()>>1 - (textMaxWidth >> 1)
	default:
		if v, err := parseFlexibleValue(offset.Left, float64(p.Width())); err != nil {
			return BoxZero, fmt.Errorf("error parsing title position: %w", err)
		} else {
			titleX = int(v)
		}
	}
	titleY := 0
	switch offset.Top {
	case "", PositionTop:
		// leave default of zero
	case PositionBottom:
		titleY = p.Height() - textTotalHeight
	default:
		if v, err := parseFlexibleValue(offset.Top, float64(p.Height())); err != nil {
			return BoxZero, fmt.Errorf("error parsing title position: %w", err)
		} else {
			titleY = int(v)
		}
	}
	startY := titleY
	for _, item := range measureOptions {
		x := titleX + (textMaxWidth-item.width)>>1
		y := titleY + item.height
		p.Text(item.text, x, y, 0, item.style)
		titleY = y
	}

	return Box{
		Top:    startY,
		Bottom: titleY,
		Left:   titleX,
		Right:  titleX + width,
		IsSet:  true,
	}, nil
}
