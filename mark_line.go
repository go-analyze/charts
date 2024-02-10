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
	m.options = append(m.options, opt)
}

// NewMarkLinePainter returns a mark line renderer
func NewMarkLinePainter(p *Painter) *markLinePainter {
	return &markLinePainter{
		p:       p,
		options: make([]markLineRenderOption, 0),
	}
}

type markLineRenderOption struct {
	FillColor   Color
	FontColor   Color
	StrokeColor Color
	Font        *truetype.Font
	Series      Series
	Range       axisRange
}

func (m *markLinePainter) Render() (Box, error) {
	painter := m.p
	for _, opt := range m.options {
		s := opt.Series
		if len(s.MarkLine.Data) == 0 {
			continue
		}
		font := opt.Font
		if font == nil {
			font = GetDefaultFont()
		}
		summary := s.Summary()
		for _, markLine := range s.MarkLine.Data {
			// since the mark line will modify the style, it must be reset every time
			painter.OverrideDrawingStyle(Style{
				FillColor:   opt.FillColor,
				StrokeColor: opt.StrokeColor,
				StrokeWidth: 1,
				StrokeDashArray: []float64{
					4,
					2,
				},
			}).OverrideTextStyle(Style{
				Font:      font,
				FontColor: opt.FontColor,
				FontSize:  labelFontSize,
			})
			value := float64(0)
			switch markLine.Type {
			case SeriesMarkDataTypeMax:
				value = summary.MaxValue
			case SeriesMarkDataTypeMin:
				value = summary.MinValue
			default:
				value = summary.AverageValue
			}
			y := opt.Range.getRestHeight(value)
			width := painter.Width()
			text := commafWithDigits(value)
			textBox := painter.MeasureText(text)
			painter.MarkLine(0, y, width-2)
			painter.Text(text, width, y+textBox.Height()>>1-2)
		}
	}
	return BoxZero, nil
}
