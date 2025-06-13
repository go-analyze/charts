package charts

import (
	"errors"
	"math"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

type scatterChart struct {
	p   *Painter
	opt *ScatterChartOption
}

// newScatterChart returns a scatter chart renderer.
func newScatterChart(p *Painter, opt ScatterChartOption) *scatterChart {
	return &scatterChart{
		p:   p,
		opt: &opt,
	}
}

// NewScatterChartOptionWithData returns an initialized ScatterChartOption with the SeriesList set for the provided data slice.
func NewScatterChartOptionWithData(data [][]float64) ScatterChartOption {
	return NewScatterChartOptionWithSeries(NewSeriesListScatter(data))
}

// NewScatterChartOptionWithSeries returns an initialized ScatterChartOption with the provided SeriesList.
func NewScatterChartOptionWithSeries(sl ScatterSeriesList) ScatterChartOption {
	return ScatterChartOption{
		SeriesList: sl,
		Padding:    defaultPadding,
		Theme:      GetDefaultTheme(),
		Font:       GetDefaultFont(),
		XAxis: XAxisOption{
			Labels: make([]string, getSeriesMaxDataCount(sl)),
		},
		YAxis:          make([]YAxisOption, getSeriesYAxisCount(sl)),
		ValueFormatter: defaultValueFormatter,
	}
}

// ScatterChartOption defines the options for rendering a scatter chart. Render the chart using Painter.ScatterChart.
type ScatterChartOption struct {
	// Theme specifies the colors used for the scatter chart.
	Theme ColorPalette
	// Padding specifies the padding around the chart.
	Padding Box
	// Deprecated: Font is deprecated, instead the font needs to be set on the SeriesLabel, or other specific elements.
	Font *truetype.Font
	// SeriesList provides the data population for the chart, typically constructed using NewSeriesListScatter or
	// NewSeriesListScatterMultiValue.
	SeriesList ScatterSeriesList
	// XAxis are options for the x-axis.
	XAxis XAxisOption
	// YAxis are options for the y-axis (at most two).
	YAxis []YAxisOption
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// Symbol specifies the default shape to draw each point with. The default
	// is 'dot'. Valid options are 'circle', 'dot', 'square' and 'diamond'. This
	// can also be set per-series.
	Symbol Symbol
	// SymbolSize specifies the size for each data point, default is 2.0.
	SymbolSize float64
	// ValueFormatter defines how float values should be rendered to strings, notably for numeric axis labels.
	ValueFormatter ValueFormatter
}

const defaultSymbolSize = 2.0

func (s *scatterChart) renderChart(result *defaultRenderResult) (Box, error) {
	p := s.p
	opt := s.opt
	if len(opt.SeriesList) == 0 {
		return BoxZero, errors.New("empty series list")
	}
	seriesPainter := result.seriesPainter

	symbolSize := defaultSymbolSize
	if opt.SymbolSize > 0 {
		symbolSize = opt.SymbolSize
	}
	xValues := boundaryGapAxisPositions(seriesPainter.Width(), flagIs(true, opt.XAxis.BoundaryGap),
		chartdraw.MaxInt(getSeriesMaxDataCount(opt.SeriesList), len(opt.XAxis.Labels)))

	markLinePainter := newMarkLinePainter(seriesPainter)
	trendLinePainter := newTrendLinePainter(seriesPainter)
	rendererList := []renderer{markLinePainter, trendLinePainter}

	seriesNames := opt.SeriesList.names()
	var points []Point
	for index, series := range opt.SeriesList {
		seriesSymbol := series.Symbol
		if seriesSymbol == "" {
			seriesSymbol = opt.Symbol
		}
		seriesColor := opt.Theme.GetSeriesColor(index)
		yRange := result.yaxisRanges[series.YAxisIndex]
		var labelPainter *seriesLabelPainter
		if flagIs(true, series.Label.Show) {
			labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, series.Label, opt.Theme)
			rendererList = append(rendererList, labelPainter)
		}

		if points == nil {
			points = make([]Point, 0, len(series.Values))
		} else {
			points = points[:0]
		}
		for i, sampleValues := range series.Values {
			allNull := true
			for _, item := range sampleValues {
				if item == GetNullValue() {
					continue
				}
				allNull = false
				p := Point{
					X: xValues[i],
					Y: yRange.getRestHeight(item),
				}
				points = append(points, p)

				if series.Label.FontStyle.Font == nil {
					series.Label.FontStyle.Font = opt.Font
				}
				if labelPainter != nil {
					labelPainter.Add(labelValue{
						index:     index,
						value:     item,
						x:         p.X,
						y:         p.Y,
						fontStyle: series.Label.FontStyle,
					})
				}
			}
			if allNull {
				points = append(points, Point{X: xValues[i], Y: math.MaxInt32})
			}
		}

		// Draw points
		switch seriesSymbol {
		case SymbolCircle:
			seriesPainter.Dots(points, opt.Theme.GetBackgroundColor(), seriesColor, 1.0, symbolSize)
		case SymbolSquare:
			seriesPainter.squares(points, seriesColor, seriesColor, 1.0, ceilFloatToInt(symbolSize*2.0))
		case SymbolDiamond:
			seriesPainter.diamonds(points, seriesColor, seriesColor, 1.0, ceilFloatToInt(symbolSize*2.8))
		default:
			seriesPainter.Dots(points, seriesColor, seriesColor, 1.0, symbolSize)
		}

		if len(series.MarkLine.Lines) > 0 {
			markLinePainter.add(markLineRenderOption{
				fillColor:    seriesColor,
				fontColor:    opt.Theme.GetMarkTextColor(),
				strokeColor:  seriesColor,
				font:         opt.Font,
				marklines:    series.MarkLine.Lines.filterGlobal(false),
				seriesValues: series.getValues(),
				axisRange:    yRange,
				valueFormatter: getPreferredValueFormatter(series.MarkLine.ValueFormatter,
					series.Label.ValueFormatter, opt.ValueFormatter),
			})
		}
		if len(series.TrendLine) > 0 {
			trendLinePainter.add(trendLineRenderOption{
				defaultStrokeColor: opt.Theme.GetSeriesTrendColor(index),
				xValues:            xValues,
				seriesValues:       series.avgValues(),
				axisRange:          yRange,
				trends:             series.TrendLine,
			})
		}
	}

	if err := doRender(rendererList...); err != nil {
		return BoxZero, err
	}
	return p.box, nil
}

func (s *scatterChart) Render() (Box, error) {
	p := s.p
	opt := s.opt
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.theme)
	}
	// boundary gap default must be set here as it's used by the x-axis as well
	if opt.XAxis.BoundaryGap == nil {
		opt.XAxis.BoundaryGap = Ptr(false)
	}
	if opt.Legend.Symbol == "" {
		if opt.Symbol == "" {
			opt.Legend.Symbol = SymbolDot
		} else {
			opt.Legend.Symbol = opt.Symbol
		}
	}

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:          opt.Theme,
		padding:        opt.Padding,
		seriesList:     opt.SeriesList,
		xAxis:          &s.opt.XAxis,
		yAxis:          opt.YAxis,
		title:          opt.Title,
		legend:         &s.opt.Legend,
		valueFormatter: opt.ValueFormatter,
	})
	if err != nil {
		return BoxZero, err
	}
	return s.renderChart(renderResult)
}
