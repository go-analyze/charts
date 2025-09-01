package charts

import (
	"github.com/golang/freetype/truetype"
)

// SeriesMarkPoint configures mark points for a series.
type SeriesMarkPoint struct {
	// SymbolSize is the width of symbol, default value is 28.
	SymbolSize int
	// ValueFormatter is used to produce the label for the Mark Point.
	ValueFormatter ValueFormatter
	// Points are the mark points for the series.
	Points SeriesMarkList
}

// AddPoints adds mark points for the series.
func (m *SeriesMarkPoint) AddPoints(markTypes ...string) {
	m.Points = appendMarks(m.Points, false, markTypes)
}

// AddGlobalPoints adds "global" mark points, which reference the sum of all series. These marks
// are only rendered when the Series is "Stacked" and the mark point is on the LAST Series of the SeriesList.
func (m *SeriesMarkPoint) AddGlobalPoints(markTypes ...string) {
	m.Points = appendMarks(m.Points, true, markTypes)
}

// NewMarkPoint returns a mark point for the provided types. Set on a specific Series instance.
func NewMarkPoint(markPointTypes ...string) SeriesMarkPoint {
	return SeriesMarkPoint{
		Points: NewSeriesMarkList(markPointTypes...),
	}
}

type markPointPainter struct {
	p       *Painter
	options []markPointRenderOption
}

func (m *markPointPainter) add(opt markPointRenderOption) {
	if opt.valueFormatter == nil {
		opt.valueFormatter = defaultValueFormatter
	}
	if opt.symbolSize == 0 {
		opt.symbolSize = 28
	}
	m.options = append(m.options, opt)
}

type markPointRenderOption struct {
	fillColor          Color
	font               *truetype.Font
	symbolSize         int
	seriesValues       []float64
	markpoints         []SeriesMark
	seriesLabelPainter *seriesLabelPainter
	points             []Point
	valueFormatter     ValueFormatter
}

// newMarkPointPainter returns a mark point renderer.
func newMarkPointPainter(p *Painter) *markPointPainter {
	return &markPointPainter{
		p: p,
	}
}

func (m *markPointPainter) Render() (Box, error) {
	painter := m.p
	for _, opt := range m.options {
		if len(opt.markpoints) == 0 {
			continue
		}
		summary := summarizePopulationData(opt.seriesValues)
		textStyle := FontStyle{
			FontSize: defaultLabelFontSize,
			Font:     opt.font,
		}
		if isLightColor(opt.fillColor) {
			textStyle.FontColor = defaultLightFontColor
		} else {
			textStyle.FontColor = defaultDarkFontColor
		}
		for _, markPointData := range opt.markpoints {
			textStyle.FontSize = defaultLabelFontSize
			index := summary.MinIndex
			value := summary.Min
			switch markPointData.Type {
			case SeriesMarkTypeMax:
				index = summary.MaxIndex
				value = summary.Max
			}
			p := opt.points[index]
			if opt.seriesLabelPainter != nil {
				// the series label has been replaced by our MarkPoint
				// This is why MarkPoints must be rendered BEFORE series labels
				opt.seriesLabelPainter.values[index].text = ""
			}

			painter.Pin(p.X, p.Y-opt.symbolSize>>1, opt.symbolSize, opt.fillColor, opt.fillColor, 0.0)
			text := opt.valueFormatter(value)
			textBox := painter.MeasureText(text, 0, textStyle)
			if textStyle.FontSize > smallLabelFontSize && textBox.Width() > opt.symbolSize {
				textStyle.FontSize = smallLabelFontSize
				textBox = painter.MeasureText(text, 0, textStyle)
			}
			painter.Text(text, p.X-textBox.Width()>>1, p.Y-opt.symbolSize>>1-2, 0, textStyle)
		}
	}
	return BoxZero, nil
}
