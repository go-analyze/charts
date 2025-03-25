package charts

import (
	"strings"

	"github.com/dustin/go-humanize"
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
	values      []labelRenderValue
}

func newSeriesLabelPainter(p *Painter, seriesNames []string, label SeriesLabel,
	theme ColorPalette) *seriesLabelPainter {
	return &seriesLabelPainter{
		p:           p,
		seriesNames: seriesNames,
		label:       &label,
		theme:       theme,
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
	if label.ValueFormatter != nil && label.FormatTemplate == "" {
		text = label.ValueFormatter(value.value)
	} else {
		text = labelFormatValue(o.seriesNames, label.FormatTemplate, label.ValueFormatter,
			value.index, value.value, -1)
	}
	labelStyle := FontStyle{
		FontColor: o.theme.GetLabelTextColor(),
		FontSize:  defaultLabelFontSize,
		Font:      getPreferredFont(label.FontStyle.Font, value.fontStyle.Font),
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

// labelFormatPie formats the value for a pie chart label.
func labelFormatPie(seriesNames []string, layout string, valueFormatter ValueFormatter,
	index int, value float64, percent float64) string {
	if len(layout) == 0 {
		layout = "{b}: {d}"
	}
	return newLabelFormatter(seriesNames, layout, valueFormatter)(index, value, percent)
}

// labelFormatFunnel formats the value for a funnel chart label.
func labelFormatFunnel(seriesNames []string, layout string, valueFormatter ValueFormatter,
	index int, value float64, percent float64) string {
	if len(layout) == 0 {
		layout = "{b}({d})"
	}
	return newLabelFormatter(seriesNames, layout, valueFormatter)(index, value, percent)
}

// labelFormatValue returns a formatted value.
func labelFormatValue(seriesNames []string, layout string, valueFormatter ValueFormatter,
	index int, value float64, percent float64) string {
	if len(layout) == 0 {
		layout = "{c}"
	}
	return newLabelFormatter(seriesNames, layout, valueFormatter)(index, value, percent)
}

// newLabelFormatter returns a label formatter.
func newLabelFormatter(seriesNames []string, layout string, valueFormatter ValueFormatter) func(index int, value float64, percent float64) string {
	if valueFormatter == nil {
		valueFormatter = func(f float64) string {
			return humanize.FtoaWithDigits(f, 2)
		}
	}
	return func(index int, value, percent float64) string {
		var percentText string
		if percent >= 0 {
			percentText = humanize.FtoaWithDigits(percent*100, 2) + "%"
		}
		var name string
		if len(seriesNames) > index {
			name = seriesNames[index]
		}
		text := strings.ReplaceAll(layout, "{c}", valueFormatter(value))
		text = strings.ReplaceAll(text, "{d}", percentText)
		text = strings.ReplaceAll(text, "{b}", name)
		return text
	}
}
