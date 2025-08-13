package charts

import (
	"fmt"
	"strings"
)

// TitleOption configures rendering of a chart title.
type TitleOption struct {
	// Show specifies if the title should be rendered, set this to *false (through Ptr(false)) to hide the title.
	Show *bool
	// Theme specifies the colors used for the title.
	Theme ColorPalette
	// Text specifies the title text, supporting '\n' for new lines.
	Text string
	// Subtext to the title, supporting '\n' for new lines.
	Subtext string
	// Offset allows you to specify the position of the title component relative to the left and top side.
	Offset OffsetStr
	// FontStyle specifies the font, size, and style for rendering the title.
	FontStyle FontStyle
	// SubtextFontStyle specifies the font, size, and style for rendering the subtext.
	SubtextFontStyle FontStyle
	// BorderWidth can be set to a non-zero value to render a box around the title.
	BorderWidth float64
}

type titleMeasureOption struct {
	width  int
	height int
	text   string
	style  FontStyle
}

func splitTitleText(text string) []string {
	arr := strings.Split(text, "\n")
	result := make([]string, 0, len(arr))
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

type titlePainter struct {
	p   *Painter
	opt *TitleOption
}

// newTitlePainter returns a title renderer.
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
	} else if opt.Text == "" && opt.Subtext == "" {
		return BoxZero, nil
	}

	theme := opt.Theme
	if theme == nil {
		theme = getPreferredTheme(p.theme)
	}

	fontStyle := fillFontStyleDefaults(opt.FontStyle, defaultFontSize, theme.GetTitleTextColor())
	subtextFontStyle := fillFontStyleDefaults(opt.SubtextFontStyle,
		fontStyle.FontSize, fontStyle.FontColor, fontStyle.Font)

	mainSplit := splitTitleText(opt.Text)
	subSplit := splitTitleText(opt.Subtext)
	measureOptions := make([]titleMeasureOption, 0, len(mainSplit)+len(subSplit))
	for _, v := range mainSplit {
		measureOptions = append(measureOptions, titleMeasureOption{
			text:  v,
			style: fontStyle,
		})
	}
	for _, v := range subSplit {
		measureOptions = append(measureOptions, titleMeasureOption{
			text:  v,
			style: subtextFontStyle,
		})
	}
	var textMaxWidth, textMaxHeight, textTotalHeight int
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
	var titleX int
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
	var titleY int
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

	result := Box{
		Top:    startY,
		Bottom: titleY,
		Left:   titleX,
		Right:  titleX + width,
		IsSet:  true,
	}

	if opt.BorderWidth > 0 {
		boxPad := 10 // built in adjustment for possible measure vs render variations
		boxPoints := []Point{
			{X: result.Left - boxPad, Y: result.Bottom + boxPad},
			{X: result.Left - boxPad, Y: result.Top - boxPad},
			{X: result.Left + result.Width() + boxPad, Y: result.Top - boxPad},
			{X: result.Left + result.Width() + boxPad, Y: result.Bottom + boxPad},
			{X: result.Left - boxPad, Y: result.Bottom + boxPad},
		}
		p.LineStroke(boxPoints, theme.GetTitleBorderColor(), opt.BorderWidth)
	}
	return result, nil
}
