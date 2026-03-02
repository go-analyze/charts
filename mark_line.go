package charts

import (
	"github.com/golang/freetype/truetype"
)

// SeriesMarkLine configures mark lines for a series.
type SeriesMarkLine struct {
	// ValueFormatter is used to produce the label for the Mark Line.
	ValueFormatter ValueFormatter
	// Lines are the mark lines for the series.
	Lines SeriesMarkList
}

// AddLines adds mark lines for the series.
func (m *SeriesMarkLine) AddLines(markTypes ...string) {
	m.Lines = appendMarks(m.Lines, false, markTypes)
}

// AddGlobalLines adds "global" mark lines, which reference the sum of all series. These marks
// are only rendered when the Series is "Stacked" and the mark line is on the LAST Series of the SeriesList.
func (m *SeriesMarkLine) AddGlobalLines(markTypes ...string) {
	m.Lines = appendMarks(m.Lines, true, markTypes)
}

// NewMarkLine returns a mark line for the provided types. Set on a specific Series instance.
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
	seriesSummary  *PopulationSummary
	marklines      []SeriesMark
	axisRange      axisRange // For vertical bar charts: y-axis range; for horizontal bar charts: x-axis range
	valueFormatter ValueFormatter
	verticalLine   bool
	// horizontalUsesHeight switches horizontal mark-line Y mapping from getRestHeight (default)
	// to getHeight. Violin bucket space uses getHeight.
	horizontalUsesHeight bool
}

func (m *markLinePainter) Render() (Box, error) {
	painter := m.p
	for _, opt := range m.options {
		if len(opt.marklines) == 0 {
			continue
		}
		summary := resolveMarkLineSummary(opt)
		fontStyle := FontStyle{
			Font:      getPreferredFont(opt.font),
			FontColor: opt.fontColor,
			FontSize:  defaultLabelFontSize,
		}
		for _, markLine := range opt.marklines {
			value := resolveSeriesMarkLineValue(markLine.Type, summary)
			text := opt.valueFormatter(value)
			textBox := painter.MeasureText(text, 0, fontStyle)
			m.renderOne(opt, text, textBox, value, painter, fontStyle)
		}
	}
	return BoxZero, nil
}

func resolveMarkLineSummary(opt markLineRenderOption) PopulationSummary {
	if opt.seriesSummary != nil {
		return *opt.seriesSummary
	}
	return summarizePopulationData(opt.seriesValues)
}

func (m *markLinePainter) renderOne(opt markLineRenderOption, text string, textBox Box, value float64,
	painter *Painter, fontStyle FontStyle) {
	if opt.verticalLine {
		x := opt.axisRange.getHeight(value) // x coordinate for the mark line
		height := painter.Height() - 2
		if height <= 0 {
			return
		}
		painter.VerticalMarkLine(x, 2, height, opt.fillColor, opt.strokeColor, 1, []float64{4, 2})
		painter.Text(text, x-(textBox.Width()>>1)-1, 0, 0, fontStyle)
		return
	}
	y := opt.axisRange.getRestHeight(value)
	if opt.horizontalUsesHeight {
		y = opt.axisRange.getHeight(value)
	}
	width := painter.Width() - 2
	if width <= 0 {
		return
	}
	painter.HorizontalMarkLine(0, y, width, opt.fillColor, opt.strokeColor, 1, []float64{4, 2})
	painter.Text(text, painter.Width(), y+(textBox.Height()>>1)-2, 0, fontStyle)
}

func resolveSeriesMarkLineValue(markType string, summary PopulationSummary) float64 {
	switch markType {
	case SeriesMarkTypeMax:
		return summary.Max
	case SeriesMarkTypeMin:
		return summary.Min
	case SeriesMarkTypeMedian:
		return summary.Median
	default:
		return summary.Average
	}
}
