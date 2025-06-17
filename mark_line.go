package charts

import (
	"github.com/golang/freetype/truetype"
)

// NewMarkLine returns a mark line for the provided types, this is set on a specific instance within a Series.
func NewMarkLine(markLineTypes ...string) SeriesMarkLine {
	return SeriesMarkLine{
		Lines: NewSeriesMarkList(markLineTypes...),
	}
}

type markLinePainter struct {
	p       *Painter
	options []markLineRenderOption
}

func (m *markLinePainter) add(opt markLineRenderOption) {
	if opt.valueFormatter == nil {
		opt.valueFormatter = defaultValueFormatter
	}
	m.options = append(m.options, opt)
}

// newMarkLinePainter returns a mark line renderer.
func newMarkLinePainter(p *Painter) *markLinePainter {
	return &markLinePainter{
		p: p,
	}
}

type markLineRenderOption struct {
	fillColor      Color
	fontColor      Color
	strokeColor    Color
	font           *truetype.Font
	seriesValues   []float64
	marklines      []SeriesMark
	axisRange      axisRange // For vertical bar charts this is y-axis range; for horizontal bar charts this is the x-axis range
	valueFormatter ValueFormatter
	verticalLine   bool
}

func (m *markLinePainter) Render() (Box, error) {
	painter := m.p
	for _, opt := range m.options {
		if len(opt.marklines) == 0 {
			continue
		}
		summary := summarizePopulationData(opt.seriesValues)
		fontStyle := FontStyle{
			Font:      getPreferredFont(opt.font),
			FontColor: opt.fontColor,
			FontSize:  defaultLabelFontSize,
		}
		for _, markLine := range opt.marklines {
			var value float64
			switch markLine.Type {
			case SeriesMarkTypeMax:
				value = summary.Max
			case SeriesMarkTypeMin:
				value = summary.Min
			default:
				value = summary.Average
			}
			text := opt.valueFormatter(value)
			textBox := painter.MeasureText(text, 0, fontStyle)
			if opt.verticalLine {
				x := opt.axisRange.getHeight(value) // x coordinate for the mark line
				height := painter.Height()
				painter.VerticalMarkLine(x, 2, height-2, opt.fillColor, opt.strokeColor, 1, []float64{4, 2})
				painter.Text(text, x-(textBox.Width()>>1)-1, 0, 0, fontStyle)
			} else { // horizontal mark line
				y := opt.axisRange.getRestHeight(value) // y coordinate for the mark line.
				width := painter.Width()
				painter.HorizontalMarkLine(0, y, width-2, opt.fillColor, opt.strokeColor, 1, []float64{4, 2})
				painter.Text(text, width, y+(textBox.Height()>>1)-2, 0, fontStyle)
			}
		}
	}
	return BoxZero, nil
}
