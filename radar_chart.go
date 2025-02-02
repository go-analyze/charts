package charts

import (
	"errors"

	"github.com/dustin/go-humanize"
	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

type radarChart struct {
	p   *Painter
	opt *RadarChartOption
}

type RadarIndicator struct {
	// Name specifies a name for the iIndicator.
	Name string
	// Max is the maximum value of indicator.
	Max float64
	// Min is the minimum value of indicator.
	Min float64
}

// NewRadarChartOptionWithData returns an initialized RadarChartOption with the SeriesList set for the provided data slice.
func NewRadarChartOptionWithData(data [][]float64, names []string, values []float64) RadarChartOption {
	return RadarChartOption{
		SeriesList:      NewSeriesListRadar(data),
		RadarIndicators: NewRadarIndicators(names, values),
		Padding:         defaultPadding,
		Theme:           GetDefaultTheme(),
		Font:            GetDefaultFont(),
	}
}

type RadarChartOption struct {
	// Theme specifies the colors used for the pie chart.
	Theme ColorPalette
	// Padding specifies the padding of pie chart.
	Padding Box
	// Font is the font used to render the chart.
	Font *truetype.Font
	// SeriesList provides the data series.
	SeriesList SeriesList
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// RadarIndicators provides the radar indicator list.
	RadarIndicators []RadarIndicator
	// backgroundIsFilled is set to true if the background is filled.
	backgroundIsFilled bool
}

// NewRadarIndicators returns a radar indicator list
func NewRadarIndicators(names []string, values []float64) []RadarIndicator {
	if len(names) != len(values) {
		return nil
	}
	indicators := make([]RadarIndicator, len(names))
	for index, name := range names {
		indicators[index] = RadarIndicator{
			Name: name,
			Max:  values[index],
		}
	}
	return indicators
}

// newRadarChart returns a radar chart renderer.
func newRadarChart(p *Painter, opt RadarChartOption) *radarChart {
	return &radarChart{
		p:   p,
		opt: &opt,
	}
}

func (r *radarChart) render(result *defaultRenderResult, seriesList SeriesList) (Box, error) {
	opt := r.opt
	indicators := opt.RadarIndicators
	sides := len(indicators)
	if sides < 3 {
		return BoxZero, errors.New("the count of indicator should be >= 3")
	}
	maxValues := make([]float64, len(indicators))
	for _, series := range seriesList {
		for index, item := range series.Data {
			if index < len(maxValues) && item > maxValues[index] {
				maxValues[index] = item
			}
		}
	}
	for index, indicator := range indicators {
		if indicator.Max <= 0 {
			indicators[index].Max = maxValues[index]
		}
	}

	radiusValue := ""
	for _, series := range seriesList {
		if len(series.Radius) != 0 {
			radiusValue = series.Radius
		}
	}

	seriesPainter := result.seriesPainter
	theme := opt.Theme

	cx := seriesPainter.Width() >> 1
	cy := seriesPainter.Height() >> 1
	diameter := chartdraw.MinInt(seriesPainter.Width(), seriesPainter.Height())
	radius := getRadius(float64(diameter), radiusValue)

	divideCount := 5
	divideRadius := float64(int(radius / float64(divideCount)))
	radius = divideRadius * float64(divideCount)

	center := Point{X: cx, Y: cy}
	for i := 0; i < divideCount; i++ {
		seriesPainter.Polygon(center, divideRadius*float64(i+1), sides, theme.GetAxisSplitLineColor(), 1)
	}
	points := getPolygonPoints(center, radius, sides)
	for _, p := range points {
		seriesPainter.moveTo(center.X, center.Y)
		seriesPainter.lineTo(p.X, p.Y)
		seriesPainter.stroke(theme.GetAxisSplitLineColor(), 1)
	}
	fontStyle := FontStyle{
		FontColor: theme.GetTextColor(),
		FontSize:  labelFontSize,
		Font:      opt.Font,
	}
	offset := 5
	// text generation
	for index, p := range points {
		name := indicators[index].Name
		b := seriesPainter.MeasureText(name, 0, fontStyle)
		isXCenter := p.X == center.X
		isYCenter := p.Y == center.Y
		isRight := p.X > center.X
		isLeft := p.X < center.X
		isTop := p.Y < center.Y
		isBottom := p.Y > center.Y
		x := p.X
		y := p.Y
		if isXCenter {
			x -= b.Width() >> 1
			if isTop {
				y -= b.Height()
			} else {
				y += b.Height()
			}
		}
		if isYCenter {
			y += b.Height() >> 1
		}
		if isTop {
			y += offset
		}
		if isBottom {
			y += offset
		}
		if isRight {
			x += offset
		}
		if isLeft {
			x -= b.Width() + offset
		}
		seriesPainter.Text(name, x, y, 0, fontStyle)
	}

	// radar chart
	angles := getPolygonPointAngles(sides)
	maxCount := len(indicators)
	for index, series := range seriesList {
		linePoints := make([]Point, 0, maxCount)
		for j, item := range series.Data {
			if j >= maxCount {
				continue
			}
			indicator := indicators[j]
			var percent float64
			offset := indicator.Max - indicator.Min
			if offset > 0 {
				percent = (item - indicator.Min) / offset
			}
			r := percent * radius
			p := getPolygonPoint(center, r, angles[j])
			linePoints = append(linePoints, p)
		}
		color := theme.GetSeriesColor(index)
		dotFillColor := ColorWhite
		if theme.IsDark() {
			dotFillColor = color
		}
		linePoints = append(linePoints, linePoints[0])
		seriesPainter.LineStroke(linePoints, color, defaultStrokeWidth)
		seriesPainter.FillArea(linePoints, color.WithAlpha(20))
		dotWith := defaultDotWidth
		for index, point := range linePoints {
			seriesPainter.Circle(dotWith, point.X, point.Y, dotFillColor, color, defaultStrokeWidth)
			if flagIs(true, series.Label.Show) && index < len(series.Data) {
				value := humanize.FtoaWithDigits(series.Data[index], 2)
				b := seriesPainter.MeasureText(value, 0, fontStyle)
				seriesPainter.Text(value, point.X-b.Width()/2, point.Y, 0, fontStyle)
			}
		}
	}

	return r.p.box, nil
}

func (r *radarChart) Render() (Box, error) {
	p := r.p
	opt := r.opt
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.theme)
	}

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:      opt.Theme,
		padding:    opt.Padding,
		seriesList: opt.SeriesList,
		xAxis: &XAxisOption{
			Show: False(),
		},
		yAxis: []YAxisOption{
			{
				Show: False(),
			},
		},
		title:              opt.Title,
		legend:             &r.opt.Legend,
		backgroundIsFilled: opt.backgroundIsFilled,
	})
	if err != nil {
		return BoxZero, err
	}
	seriesList := opt.SeriesList.Filter(ChartTypeRadar)
	return r.render(renderResult, seriesList)
}
