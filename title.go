package charts

import (
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"
)

type TitleOption struct {
	// The theme of chart
	Theme ColorPalette
	// Title text, support \n for new line
	Text string
	// Subtitle text, support \n for new line
	Subtext string
	// Distance between title component and the left side of the container.
	// It can be pixel value: 20, percentage value: 20%,
	// or position value: right, center.
	Left string
	// Distance between title component and the top side of the container.
	// It can be pixel value: 20.
	Top string
	// The font of label
	Font *truetype.Font
	// The font size of label
	FontSize float64
	// The color of label
	FontColor Color
	// The subtext font size of label
	SubtextFontSize float64
	// The subtext font color of label
	SubtextFontColor Color
}

type titleMeasureOption struct {
	width  int
	height int
	text   string
	style  Style
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

// NewTitlePainter returns a title renderer
func NewTitlePainter(p *Painter, opt TitleOption) *titlePainter {
	return &titlePainter{
		p:   p,
		opt: &opt,
	}
}

func (t *titlePainter) Render() (Box, error) {
	opt := t.opt
	p := t.p
	theme := opt.Theme

	if theme == nil {
		theme = p.theme
	}
	if opt.Text == "" && opt.Subtext == "" {
		return BoxZero, nil
	}

	measureOptions := make([]titleMeasureOption, 0)

	if opt.Font == nil {
		opt.Font = theme.GetFont()
	}
	if opt.FontColor.IsZero() {
		opt.FontColor = theme.GetTextColor()
	}
	if opt.FontSize == 0 {
		opt.FontSize = theme.GetFontSize()
	}
	if opt.SubtextFontColor.IsZero() {
		opt.SubtextFontColor = opt.FontColor
	}
	if opt.SubtextFontSize == 0 {
		opt.SubtextFontSize = opt.FontSize
	}

	titleTextStyle := Style{
		Font:      opt.Font,
		FontSize:  opt.FontSize,
		FontColor: opt.FontColor,
	}
	// 主标题
	for _, v := range splitTitleText(opt.Text) {
		measureOptions = append(measureOptions, titleMeasureOption{
			text:  v,
			style: titleTextStyle,
		})
	}
	subtextStyle := Style{
		Font:      opt.Font,
		FontSize:  opt.SubtextFontSize,
		FontColor: opt.SubtextFontColor,
	}
	// 副标题
	for _, v := range splitTitleText(opt.Subtext) {
		measureOptions = append(measureOptions, titleMeasureOption{
			text:  v,
			style: subtextStyle,
		})
	}
	textMaxWidth := 0
	textMaxHeight := 0
	for index, item := range measureOptions {
		p.OverrideTextStyle(item.style)
		textBox := p.MeasureText(item.text)

		w := textBox.Width()
		h := textBox.Height()
		if w > textMaxWidth {
			textMaxWidth = w
		}
		if h > textMaxHeight {
			textMaxHeight = h
		}
		measureOptions[index].height = h
		measureOptions[index].width = w
	}
	width := textMaxWidth

	titleX := 0
	switch opt.Left {
	case PositionRight:
		titleX = p.Width() - textMaxWidth
	case PositionCenter:
		titleX = p.Width()>>1 - (textMaxWidth >> 1)
	default:
		if strings.HasSuffix(opt.Left, "%") {
			value, _ := strconv.Atoi(strings.ReplaceAll(opt.Left, "%", ""))
			titleX = p.Width() * value / 100
		} else {
			value, _ := strconv.Atoi(opt.Left)
			titleX = value
		}
	}
	titleY := 0
	// TODO TOP 暂只支持数值
	if opt.Top != "" {
		value, _ := strconv.Atoi(opt.Top)
		titleY += value
	}
	for _, item := range measureOptions {
		p.OverrideTextStyle(item.style)
		x := titleX + (textMaxWidth-item.width)>>1
		y := titleY + item.height
		p.Text(item.text, x, y)
		titleY += item.height
	}

	return Box{
		Bottom: titleY,
		Right:  titleX + width,
	}, nil
}
