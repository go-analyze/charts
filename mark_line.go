package charts

import (
	"github.com/golang/freetype/truetype"
)

// NewMarkLine returns a series mark line
func NewMarkLine(markLineTypes ...string) SeriesMarkLine {
	data := make([]SeriesMarkData, len(markLineTypes))
	for index, t := range markLineTypes {
		data[index] = SeriesMarkData{
			Type: t,
		}
	}
	return SeriesMarkLine{
		Data: data,
	}
}

type markLinePainter struct {
	p       *Painter
	options []markLineRenderOption
}

func (m *markLinePainter) Add(opt markLineRenderOption) {
	if len(opt.series.MarkLine.Data) > 0 {
		m.options = append(m.options, opt)
	}
}

// newMarkLinePainter returns a mark line renderer
func newMarkLinePainter(p *Painter) *markLinePainter {
	return &markLinePainter{
		p:       p,
		options: make([]markLineRenderOption, 0),
	}
}

type markLineRenderOption struct {
	fillColor      Color
	fontColor      Color
	strokeColor    Color
	font           *truetype.Font
	series         Series
	axisRange      axisRange
	valueFormatter ValueFormatter
}

func (m *markLinePainter) Render() (Box, error) {
	painter := m.p
	for _, opt := range m.options {
		s := opt.series
		if len(s.MarkLine.Data) == 0 {
			continue
		}
		summary := s.Summary()
		fontStyle := FontStyle{
			Font:      getPreferredFont(opt.font),
			FontColor: opt.fontColor,
			FontSize:  labelFontSize,
		}
		valueFormatter := getPreferredValueFormatter(opt.series.MarkLine.ValueFormatter,
			opt.series.Label.ValueFormatter, opt.valueFormatter)
		for _, markLine := range s.MarkLine.Data {
			var value float64
			switch markLine.Type {
			case SeriesMarkDataTypeMax:
				value = summary.Max
			case SeriesMarkDataTypeMin:
				value = summary.Min
			default:
				value = summary.Average
			}
			y := opt.axisRange.getRestHeight(value)
			text := valueFormatter(value)
			textBox := painter.MeasureText(text, 0, fontStyle)
			width := painter.Width()
			painter.MarkLine(0, y, width-2, opt.fillColor, opt.strokeColor, 1, []float64{4, 2})
			painter.Text(text, width, y+textBox.Height()>>1-2, 0, fontStyle)
		}
	}
	return BoxZero, nil
}
