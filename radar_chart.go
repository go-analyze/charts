package charts

import (
	"errors"

	"github.com/dustin/go-humanize"
	"github.com/golang/freetype/truetype"
)

var radarDefaultValueFormatter = func(v float64) string {
	return humanize.FtoaWithDigits(v, 2)
}

type radarChart struct {
	p   *Painter
	opt *RadarChartOption
}

// RadarIndicator defines the dimensions of a radar chart axis.
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

// RadarChartOption defines the options for rendering a radar chart. Render the chart using Painter.RadarChart.
type RadarChartOption struct {
	// Theme specifies the colors used for the radar chart.
	Theme ColorPalette
	// Padding specifies the padding around the chart.
	Padding Box
	// Font is the font used to render the chart.
	Font *truetype.Font
	// SeriesList provides the data population for the chart, typically constructed using NewSeriesListRadar.
	SeriesList RadarSeriesList
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// RadarIndicators provides the radar indicator list.
	RadarIndicators []RadarIndicator
	// Radius sets the chart radius, for example "40%".
	// Default is "40%".
	Radius string
	// ValueFormatter defines how float values should be rendered to strings, notably for series labels.
	ValueFormatter ValueFormatter
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

func (r *radarChart) renderChart(result *defaultRenderResult) (Box, error) {
	opt := r.opt
	indicators := opt.RadarIndicators
	sides := len(indicators)
	if sides < 3 {
		return BoxZero, errors.New("indicator count should be at least 3")
	} else if len(opt.SeriesList) == 0 {
		return BoxZero, errors.New("empty series list")
	}
	maxValues := make([]float64, len(indicators))
	for _, series := range opt.SeriesList {
		for index, item := range series.Values {
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

	seriesPainter := result.seriesPainter
	theme := opt.Theme

	cx, cy, diameter := circleChartPosition(seriesPainter)
	radius := getFlexibleRadius(diameter, defaultPieRadiusFactor, opt.Radius)

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
		FontColor: theme.GetLabelTextColor(),
		FontSize:  defaultLabelFontSize,
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
	for index, series := range opt.SeriesList {
		valueFormatter := getPreferredValueFormatter(series.Label.ValueFormatter, opt.ValueFormatter,
			radarDefaultValueFormatter)
		linePoints := make([]Point, 0, maxCount)
		for j, item := range series.Values {
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
			if flagIs(true, series.Label.Show) && index < len(series.Values) {
				valueStr := valueFormatter(series.Values[index])
				b := seriesPainter.MeasureText(valueStr, 0, fontStyle)
				seriesPainter.Text(valueStr, point.X-b.Width()/2, point.Y, 0, fontStyle)
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
	if opt.Legend.Symbol == "" {
		// default to square symbol for this chart type
		opt.Legend.Symbol = SymbolSquare
	}

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:      opt.Theme,
		padding:    opt.Padding,
		seriesList: opt.SeriesList,
		xAxis: &XAxisOption{
			Show: Ptr(false),
		},
		yAxis: []YAxisOption{
			{
				Show: Ptr(false),
			},
		},
		title:  opt.Title,
		legend: &r.opt.Legend,
	})
	if err != nil {
		return BoxZero, err
	}
	return r.renderChart(renderResult)
}
