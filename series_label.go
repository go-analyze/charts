package charts

import (
	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

type labelRenderValue struct {
	Text      string
	FontStyle FontStyle
	X         int
	Y         int
	Radians   float64
}

type LabelValue struct {
	Index     int
	Value     float64
	X         int
	Y         int
	Radians   float64
	FontStyle FontStyle
	Vertical  bool
	Offset    OffsetInt
}

type SeriesLabelPainter struct {
	p           *Painter
	seriesNames []string
	label       *SeriesLabel
	theme       ColorPalette
	font        *truetype.Font
	values      []labelRenderValue
}

type SeriesLabelPainterParams struct {
	P           *Painter
	SeriesNames []string
	Label       SeriesLabel
	Theme       ColorPalette
	Font        *truetype.Font
}

func NewSeriesLabelPainter(params SeriesLabelPainterParams) *SeriesLabelPainter {
	return &SeriesLabelPainter{
		p:           params.P,
		seriesNames: params.SeriesNames,
		label:       &params.Label,
		theme:       params.Theme,
		font:        params.Font,
		values:      make([]labelRenderValue, 0),
	}
}

func (o *SeriesLabelPainter) Add(value LabelValue) {
	label := o.label
	distance := label.Distance
	if distance == 0 {
		distance = 5
	}
	text := NewValueLabelFormatter(o.seriesNames, label.Formatter)(value.Index, value.Value, -1)
	labelStyle := FontStyle{
		FontColor: o.theme.GetTextColor(),
		FontSize:  labelFontSize,
		Font:      getPreferredFont(label.FontStyle.Font, value.FontStyle.Font, o.font),
	}
	if label.FontStyle.FontSize != 0 {
		labelStyle.FontSize = label.FontStyle.FontSize
	} else if value.FontStyle.FontSize != 0 {
		labelStyle.FontSize = value.FontStyle.FontSize
	}
	if !label.FontStyle.FontColor.IsZero() {
		labelStyle.FontColor = label.FontStyle.FontColor
	} else if !value.FontStyle.FontColor.IsZero() {
		labelStyle.FontColor = value.FontStyle.FontColor
	}
	p := o.p
	p.OverrideDrawingStyle(chartdraw.Style{FontStyle: labelStyle})
	rotated := value.Radians != 0
	if rotated {
		p.SetTextRotation(value.Radians)
	}
	textBox := p.MeasureText(text)
	renderValue := labelRenderValue{
		Text:      text,
		FontStyle: labelStyle,
		X:         value.X,
		Y:         value.Y,
		Radians:   value.Radians,
	}
	if value.Vertical {
		renderValue.X -= textBox.Width() >> 1
		renderValue.Y -= distance
	} else {
		renderValue.X += distance
		renderValue.Y += textBox.Height() >> 1
		renderValue.Y -= 2
	}
	if rotated {
		renderValue.X = value.X + textBox.Width()>>1 - 1
		p.ClearTextRotation()
	} else if textBox.Width()%2 != 0 {
		renderValue.X++
	}
	renderValue.X += value.Offset.Left
	renderValue.Y += value.Offset.Top
	o.values = append(o.values, renderValue)
}

func (o *SeriesLabelPainter) Render() (Box, error) {
	for _, item := range o.values {
		o.p.OverrideFontStyle(item.FontStyle)
		if item.Radians != 0 {
			o.p.TextRotation(item.Text, item.X, item.Y, item.Radians)
		} else {
			o.p.Text(item.Text, item.X, item.Y)
		}
	}
	return chartdraw.BoxZero, nil
}
