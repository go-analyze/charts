package charts

import (
	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

// NewMarkPoint returns a series mark point
func NewMarkPoint(markPointTypes ...string) SeriesMarkPoint {
	data := make([]SeriesMarkData, len(markPointTypes))
	for index, t := range markPointTypes {
		data[index] = SeriesMarkData{
			Type: t,
		}
	}
	return SeriesMarkPoint{
		Data: data,
	}
}

type markPointPainter struct {
	p       *Painter
	options []markPointRenderOption
}

func (m *markPointPainter) add(opt markPointRenderOption) {
	if len(opt.series.MarkPoint.Data) > 0 {
		m.options = append(m.options, opt)
	}
}

type markPointRenderOption struct {
	fillColor          Color
	font               *truetype.Font
	series             Series
	seriesLabelPainter *seriesLabelPainter
	points             []Point
	valueFormatter     ValueFormatter
}

// newMarkPointPainter returns a mark point renderer
func newMarkPointPainter(p *Painter) *markPointPainter {
	return &markPointPainter{
		p:       p,
		options: make([]markPointRenderOption, 0),
	}
}

func (m *markPointPainter) Render() (Box, error) {
	painter := m.p
	for _, opt := range m.options {
		if len(opt.series.MarkPoint.Data) == 0 {
			continue
		}
		points := opt.points
		summary := opt.series.Summary()
		symbolSize := opt.series.MarkPoint.SymbolSize
		if symbolSize == 0 {
			symbolSize = 28
		}
		textStyle := chartdraw.Style{
			FontStyle: FontStyle{
				FontSize: labelFontSize,
				Font:     opt.font,
			},
			StrokeWidth: 1,
		}
		if isLightColor(opt.fillColor) {
			textStyle.FontColor = defaultLightFontColor
		} else {
			textStyle.FontColor = defaultDarkFontColor
		}
		valueFormatter := getPreferredValueFormatter(opt.series.MarkPoint.ValueFormatter,
			opt.series.Label.ValueFormatter, opt.valueFormatter)
		for _, markPointData := range opt.series.MarkPoint.Data {
			textStyle.FontSize = labelFontSize
			index := summary.MinIndex
			value := summary.Min
			switch markPointData.Type {
			case SeriesMarkDataTypeMax:
				index = summary.MaxIndex
				value = summary.Max
			}
			p := points[index]
			if opt.seriesLabelPainter != nil {
				// the series label has been replaced by our MarkPoint
				// This is why MarkPoints must be rendered BEFORE series labels
				opt.seriesLabelPainter.values[index].Text = ""
			}

			painter.Pin(p.X, p.Y-symbolSize>>1, symbolSize, opt.fillColor, opt.fillColor, 0.0)
			text := valueFormatter(value)
			textBox := painter.MeasureText(text, 0, textStyle.FontStyle)
			if textStyle.FontSize > smallLabelFontSize && textBox.Width() > symbolSize {
				textStyle.FontSize = smallLabelFontSize
				textBox = painter.MeasureText(text, 0, textStyle.FontStyle)
			}
			painter.Text(text, p.X-textBox.Width()>>1, p.Y-symbolSize>>1-2, 0, textStyle.FontStyle)
		}
	}
	return BoxZero, nil
}
