package charts

import (
	"math"

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
// are only rendered when the Series is "Stacked" and the mark point is on the LAST stacked Series.
func (m *SeriesMarkPoint) AddGlobalPoints(markTypes ...string) {
	m.Points = appendMarks(m.Points, true, markTypes)
}

// NewMarkPoint returns a mark point for the provided types. Set on a specific Series instance.
func NewMarkPoint(markPointTypes ...string) SeriesMarkPoint {
	return SeriesMarkPoint{
		Points: NewSeriesMarkList(markPointTypes...),
	}
}

// effectiveSymbolSize returns the rendered pin width (configured SymbolSize, or the default when set
// without an explicit size), or 0 when no points are configured.
func (m *SeriesMarkPoint) effectiveSymbolSize() int {
	if len(m.Points) == 0 {
		return 0
	} else if m.SymbolSize > 0 {
		return m.SymbolSize
	}
	return 28
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
	rotationRadians    float64
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
			Font:     getPreferredFont(opt.font),
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
			if index < 0 {
				continue // no valid data (all null)
			}
			p := opt.points[index]
			if opt.seriesLabelPainter != nil {
				// blank the label the MarkPoint replaces (rendered before series labels)
				for i := range opt.seriesLabelPainter.values {
					if opt.seriesLabelPainter.values[i].dataIndex == index {
						opt.seriesLabelPainter.values[i].text = ""
						break
					}
				}
			}

			// In the default (down-pointing) pin, the anchor passed to MarkPin is the base of
			// the tail where it joins the head; the head center sits at (0, -3w/4) relative to
			// the bar anchor p, and the outer edge of the head (opposite the tail tip) sits at
			// (0, -5w/4). Rotating these offsets gives the same positions for a rotated pin.
			anchorOffsetX, anchorOffsetY := 0, -(opt.symbolSize >> 1)
			headOuterOffsetX := 0
			headOuterOffsetY := -(5 * opt.symbolSize) / 4
			if opt.rotationRadians != 0 {
				cos := math.Cos(opt.rotationRadians)
				sin := math.Sin(opt.rotationRadians)
				rotate := func(dx, dy int) (int, int) {
					fx, fy := float64(dx), float64(dy)
					return int(math.Round(cos*fx - sin*fy)),
						int(math.Round(sin*fx + cos*fy))
				}
				anchorOffsetX, anchorOffsetY = rotate(anchorOffsetX, anchorOffsetY)
				headOuterOffsetX, _ = rotate(headOuterOffsetX, headOuterOffsetY)
			}
			drawnAnchorX, drawnAnchorY := p.X+anchorOffsetX, p.Y+anchorOffsetY
			painter.MarkPin(drawnAnchorX, drawnAnchorY, opt.symbolSize,
				opt.rotationRadians, opt.fillColor, opt.fillColor, 0.0)
			text := opt.valueFormatter(value)
			textBox := painter.MeasureText(text, 0, textStyle)
			if textStyle.FontSize > smallLabelFontSize && textBox.Width() > opt.symbolSize {
				textStyle.FontSize = smallLabelFontSize
				textBox = painter.MeasureText(text, 0, textStyle)
			}
			const textInset = 2
			var textX, textY int
			if opt.rotationRadians == 0 {
				// vertical pin: keep the horizontal centering on the head center
				textX = drawnAnchorX - (textBox.Width() >> 1)
				textY = drawnAnchorY - textInset
			} else {
				// horizontal pin: anchor the text just inside the outer head curve
				if headOuterOffsetX >= 0 {
					textX = p.X + headOuterOffsetX - textInset - textBox.Width()
				} else {
					textX = p.X + headOuterOffsetX + textInset
				}
				textY = drawnAnchorY - textInset + (textBox.Height() >> 1)
			}
			painter.Text(text, textX, textY, 0, textStyle)
		}
	}
	return BoxZero, nil
}
