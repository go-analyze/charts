package charts

import (
	"github.com/golang/freetype/truetype"
)

type labelRenderValue struct {
	Text      string
	FontStyle FontStyle
	X         int
	Y         int
	Radians   float64
}

type labelValue struct {
	index     int
	value     float64
	x         int
	y         int
	radians   float64
	fontStyle FontStyle
	vertical  bool
	offset    OffsetInt
}

type seriesLabelPainter struct {
	p           *Painter
	seriesNames []string
	label       *SeriesLabel
	theme       ColorPalette
	font        *truetype.Font
	values      []labelRenderValue
}

func newSeriesLabelPainter(p *Painter, seriesNames []string, label SeriesLabel,
	theme ColorPalette, font *truetype.Font) *seriesLabelPainter {
	return &seriesLabelPainter{
		p:           p,
		seriesNames: seriesNames,
		label:       &label,
		theme:       theme,
		font:        font,
	}
}

func (o *seriesLabelPainter) Add(value labelValue) {
	label := o.label
	if flagIs(false, label.Show) {
		return
	}
	distance := label.Distance
	if distance == 0 {
		distance = 5
	}
	var text string
	if label.ValueFormatter != nil {
		text = label.ValueFormatter(value.value)
	} else {
		if label.FormatTemplate == "" {
			label.FormatTemplate = label.Formatter
		}
		text = labelFormatValue(o.seriesNames, label.FormatTemplate, value.index, value.value, -1)
	}
	labelStyle := FontStyle{
		FontColor: o.theme.GetTextColor(),
		FontSize:  labelFontSize,
		Font:      getPreferredFont(label.FontStyle.Font, value.fontStyle.Font, o.font),
	}
	if label.FontStyle.FontSize != 0 {
		labelStyle.FontSize = label.FontStyle.FontSize
	} else if value.fontStyle.FontSize != 0 {
		labelStyle.FontSize = value.fontStyle.FontSize
	}
	if !label.FontStyle.FontColor.IsZero() {
		labelStyle.FontColor = label.FontStyle.FontColor
	} else if !value.fontStyle.FontColor.IsZero() {
		labelStyle.FontColor = value.fontStyle.FontColor
	}
	p := o.p
	textBox := p.MeasureText(text, value.radians, labelStyle)
	renderValue := labelRenderValue{
		Text:      text,
		FontStyle: labelStyle,
		X:         value.x,
		Y:         value.y,
		Radians:   value.radians,
	}
	if value.vertical {
		renderValue.X -= textBox.Width() >> 1
		renderValue.Y -= distance
	} else {
		renderValue.X += distance
		renderValue.Y += textBox.Height() >> 1
		renderValue.Y -= 2
	}
	if value.radians != 0 {
		renderValue.X = value.x + (textBox.Width() >> 1) - 1
	}
	renderValue.X += value.offset.Left
	renderValue.Y += value.offset.Top
	o.values = append(o.values, renderValue)
}

func (o *seriesLabelPainter) Render() (Box, error) {
	for _, item := range o.values {
		if item.Text != "" {
			o.p.Text(item.Text, item.X, item.Y, item.Radians, item.FontStyle)
		}
	}
	return BoxZero, nil
}
