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

func (m *markPointPainter) Add(opt markPointRenderOption) {
	m.options = append(m.options, opt)
}

type markPointRenderOption struct {
	FillColor Color
	Font      *truetype.Font
	Series    Series
	Points    []Point
}

// NewMarkPointPainter returns a mark point renderer
func NewMarkPointPainter(p *Painter) *markPointPainter {
	return &markPointPainter{
		p:       p,
		options: make([]markPointRenderOption, 0),
	}
}

func (m *markPointPainter) Render() (Box, error) {
	painter := m.p
	for _, opt := range m.options {
		if len(opt.Series.MarkPoint.Data) == 0 {
			continue
		}
		points := opt.Points
		summary := opt.Series.Summary()
		symbolSize := opt.Series.MarkPoint.SymbolSize
		if symbolSize == 0 {
			symbolSize = 28
		}
		textStyle := chartdraw.Style{
			FontStyle: FontStyle{
				FontSize: labelFontSize,
				Font:     opt.Font,
			},
			StrokeWidth: 1,
		}
		if isLightColor(opt.FillColor) {
			textStyle.FontColor = defaultLightFontColor
		} else {
			textStyle.FontColor = defaultDarkFontColor
		}
		painter.OverrideDrawingStyle(chartdraw.Style{
			FillColor: opt.FillColor,
		}).OverrideFontStyle(textStyle.FontStyle)
		for _, markPointData := range opt.Series.MarkPoint.Data {
			textStyle.FontSize = labelFontSize
			painter.OverrideFontStyle(textStyle.FontStyle)
			p := points[summary.MinIndex]
			value := summary.MinValue
			switch markPointData.Type {
			case SeriesMarkDataTypeMax:
				p = points[summary.MaxIndex]
				value = summary.MaxValue
			}

			painter.Pin(p.X, p.Y-symbolSize>>1, symbolSize)
			text := commafWithDigits(value)
			textBox := painter.MeasureText(text)
			if textBox.Width() > symbolSize {
				textStyle.FontSize = smallLabelFontSize
				painter.OverrideFontStyle(textStyle.FontStyle)
				textBox = painter.MeasureText(text)
			}
			painter.Text(text, p.X-textBox.Width()>>1, p.Y-symbolSize>>1-2)
		}
	}
	return BoxZero, nil
}
