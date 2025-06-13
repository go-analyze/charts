package charts

import (
	"errors"
	"math"

	"github.com/golang/freetype/truetype"
)

type barChart struct {
	p   *Painter
	opt *BarChartOption
}

// newBarChart returns a bar chart renderer
func newBarChart(p *Painter, opt BarChartOption) *barChart {
	return &barChart{
		p:   p,
		opt: &opt,
	}
}

// NewBarChartOptionWithData returns an initialized BarChartOption with the SeriesList set for the provided data slice.
func NewBarChartOptionWithData(data [][]float64) BarChartOption {
	return NewBarChartOptionWithSeries(NewSeriesListBar(data))
}

// NewBarChartOptionWithSeries returns an initialized BarChartOption with the provided SeriesList.
func NewBarChartOptionWithSeries(sl BarSeriesList) BarChartOption {
	return BarChartOption{
		SeriesList:     sl,
		Padding:        defaultPadding,
		Theme:          GetDefaultTheme(),
		Font:           GetDefaultFont(),
		YAxis:          make([]YAxisOption, getSeriesYAxisCount(sl)),
		ValueFormatter: defaultValueFormatter,
	}
}

// BarChartOption defines the options for rendering a bar chart. Render the chart using Painter.BarChart.
type BarChartOption struct {
	// Theme specifies the colors used for the bar chart.
	Theme ColorPalette
	// Padding specifies the padding around the chart.
	Padding Box
	// Deprecated: Font is deprecated, instead the font needs to be set on the SeriesLabel, or other specific elements.
	Font *truetype.Font
	// SeriesList provides the data population for the chart, typically constructed using NewSeriesListBar.
	SeriesList BarSeriesList
	// StackSeries, if true, renders the series stacked within one bar. This causes
	// some options, including BarMargin and SeriesLabelPosition, to be ignored.
	// MarkLine only renders for the first series and stacking only applies to the
	// first YAxis (index 0).
	StackSeries *bool
	// SeriesLabelPosition specifies the position of the label for the series. Currently supported values are
	// "top" or "bottom".
	SeriesLabelPosition string
	// XAxis are options for the x-axis.
	XAxis XAxisOption
	// YAxis are options for the y-axis (at most two).
	YAxis []YAxisOption
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// BarWidth specifies the width of each bar. Width may be reduced to ensure all series fit on the chart.
	BarWidth int
	// BarMargin specifies the margin between bars grouped together. BarWidth takes priority over the margin.
	BarMargin *float64
	// RoundedBarCaps, if true, draws bars with rounded top corners.
	RoundedBarCaps *bool
	// ValueFormatter defines how float values should be rendered to strings, notably for numeric axis labels.
	ValueFormatter ValueFormatter
}

func calculateBarMarginsAndSize(seriesCount, space int, configuredBarSize int, configuredBarMargin *float64) (int, int, int) {
	// default margins, adjusted below with config and series count
	margin := 10   // margin between each series block
	barMargin := 5 // margin between each bar
	if space < 20 {
		margin = 2
		barMargin = 2
	} else if space < 50 {
		margin = 5
		barMargin = 3
	}
	// check margin configuration if bar size allows margin
	if configuredBarSize+barMargin < space/seriesCount {
		// BarWidth is in range that we should also consider an optional margin configuration
		if configuredBarMargin != nil {
			barMargin = int(math.Round(*configuredBarMargin))
			if barMargin+configuredBarSize > space/seriesCount {
				barMargin = (space / seriesCount) - configuredBarSize
			}
		}
	} // else, bar width is out of range.  Ignore margin config

	barSize := (space - 2*margin - barMargin*(seriesCount-1)) / seriesCount
	// check bar size configuration, limited by the series count and space available
	if configuredBarSize > 0 && configuredBarSize < barSize {
		barSize = configuredBarSize
		// recalculate margin
		margin = (space - seriesCount*barSize - barMargin*(seriesCount-1)) / 2
	}

	return margin, barMargin, barSize
}

func (b *barChart) renderChart(result *defaultRenderResult) (Box, error) {
	p := b.p
	opt := b.opt
	seriesCount := len(opt.SeriesList)
	if seriesCount == 0 {
		return BoxZero, errors.New("empty series list")
	}
	seriesPainter := result.seriesPainter

	x0, x1 := result.xaxisRange.getRange(0)
	width := int(x1 - x0)
	barMaxHeight := seriesPainter.Height() // total vertical space for bars
	seriesNames := opt.SeriesList.names()
	divideValues := result.xaxisRange.autoDivide()
	stackedSeries := flagIs(true, opt.StackSeries)
	var margin, barMargin, barWidth int
	var accumulatedHeights []int // prior heights for stacking to avoid recalculating the heights
	if stackedSeries {
		barCount := getSeriesYAxisCount(opt.SeriesList) // only two bars if two y-axis
		configuredMargin := opt.BarMargin
		if barCount == 1 {
			configuredMargin = nil // no margin needed with a single bar
		}
		margin, _, barWidth = calculateBarMarginsAndSize(barCount, width, opt.BarWidth, configuredMargin)
		accumulatedHeights = make([]int, result.xaxisRange.divideCount)
	} else {
		margin, barMargin, barWidth = calculateBarMarginsAndSize(seriesCount, width, opt.BarWidth, opt.BarMargin)
	}

	markPointPainter := newMarkPointPainter(seriesPainter)
	markLinePainter := newMarkLinePainter(seriesPainter)
	// render list must start with the markPointPainter, as it can influence label painters (if enabled)
	rendererList := []renderer{markPointPainter, markLinePainter}

	for index, series := range opt.SeriesList {
		stackSeries := stackedSeries && series.YAxisIndex == 0
		yRange := result.yaxisRanges[series.YAxisIndex]
		seriesColor := opt.Theme.GetSeriesColor(index)

		var labelPainter *seriesLabelPainter
		if flagIs(true, series.Label.Show) {
			labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, series.Label, opt.Theme)
			rendererList = append(rendererList, labelPainter)
		}

		points := make([]Point, len(series.Values)) // used for mark points
		for j, item := range series.Values {
			if j >= result.xaxisRange.divideCount {
				break
			}

			// Compute bar placement differently for stacked vs non-stacked.
			var x, top, bottom int
			h := yRange.getHeight(item)

			if stackSeries {
				// Use accumulatedHeights to stack
				x = divideValues[j] + margin
				top = barMaxHeight - (accumulatedHeights[j] + h)
				bottom = barMaxHeight - accumulatedHeights[j]
				accumulatedHeights[j] += h
			} else {
				// Non-stacked: offset each series in its own lane
				x = divideValues[j] + margin + index*(barWidth+barMargin)
				top = barMaxHeight - h
				bottom = barMaxHeight - 1 // or -0, depending on your style
			}

			// In stacked mode, only round caps on the last series
			if flagIs(true, opt.RoundedBarCaps) && (!stackSeries || index == seriesCount-1) {
				seriesPainter.roundedRect(
					Box{Top: top, Left: x, Right: x + barWidth, Bottom: bottom, IsSet: true},
					barWidth, true, false, seriesColor, seriesColor, 0.0)
			} else {
				seriesPainter.FilledRect(x, top, x+barWidth, bottom, seriesColor, seriesColor, 0.0)
			}

			// Prepare point for mark points
			points[j] = Point{
				X: x + (barWidth >> 1), // center of the bar horizontally
				Y: top,                 // top of bar
			}

			if labelPainter != nil {
				labelY := top
				var radians float64
				fontStyle := series.Label.FontStyle
				labelBottom := opt.SeriesLabelPosition == PositionBottom && !stackSeries
				if labelBottom {
					labelY = barMaxHeight
					radians = -math.Pi / 2 // Rotated label at the bottom
				}
				if fontStyle.FontColor.IsZero() {
					var testColor Color
					if labelBottom {
						testColor = seriesColor
					} else if stackSeries && index+1 < seriesCount {
						testColor = opt.Theme.GetSeriesColor(index + 1)
					}
					if !testColor.IsZero() {
						if isLightColor(testColor) {
							fontStyle.FontColor = defaultLightFontColor
						} else {
							fontStyle.FontColor = defaultDarkFontColor
						}
					}
				}
				if fontStyle.Font == nil {
					fontStyle.Font = opt.Font
				}

				labelPainter.Add(labelValue{
					vertical:  true, // label is vertically oriented
					index:     index,
					value:     item,
					fontStyle: fontStyle,
					x:         x + (barWidth >> 1),
					y:         labelY,
					radians:   radians,
					offset:    series.Label.Offset,
				})
			}
		}

		var globalSeriesData []float64 // lazily initialized
		if len(series.MarkLine.Lines) > 0 {
			markLineValueFormatter := getPreferredValueFormatter(series.MarkLine.ValueFormatter,
				series.Label.ValueFormatter, opt.ValueFormatter)
			var seriesMarks, globalMarks SeriesMarkList
			if stackSeries && index == seriesCount-1 { // global is only allowed when stacked and on the last series
				seriesMarks, globalMarks = series.MarkLine.Lines.splitGlobal()
			} else {
				seriesMarks = series.MarkLine.Lines.filterGlobal(false)
			}
			if len(seriesMarks) > 0 && (!stackSeries || index == 0) {
				// in stacked mode we only support the line painter for the first series
				markLinePainter.add(markLineRenderOption{
					fillColor:      seriesColor,
					fontColor:      opt.Theme.GetMarkTextColor(),
					strokeColor:    seriesColor,
					font:           getPreferredFont(series.Label.FontStyle.Font, opt.Font),
					marklines:      seriesMarks,
					seriesValues:   series.Values,
					axisRange:      yRange,
					valueFormatter: markLineValueFormatter,
				})
			}
			if len(globalMarks) > 0 {
				if globalSeriesData == nil {
					globalSeriesData = sumSeriesData(opt.SeriesList, series.YAxisIndex)
				}
				markLinePainter.add(markLineRenderOption{
					fillColor:      defaultGlobalMarkFillColor,
					fontColor:      opt.Theme.GetMarkTextColor(),
					strokeColor:    defaultGlobalMarkFillColor,
					font:           getPreferredFont(series.Label.FontStyle.Font, opt.Font),
					marklines:      globalMarks,
					seriesValues:   globalSeriesData,
					axisRange:      yRange,
					valueFormatter: markLineValueFormatter,
				})
			}
		}
		if len(series.MarkPoint.Points) > 0 {
			markPointValueFormatter := getPreferredValueFormatter(series.MarkPoint.ValueFormatter,
				series.Label.ValueFormatter, opt.ValueFormatter)
			var seriesMarks, globalMarks SeriesMarkList
			if stackSeries && index == seriesCount-1 { // global is only allowed when stacked and on the last series
				seriesMarks, globalMarks = series.MarkPoint.Points.splitGlobal()
			} else {
				seriesMarks = series.MarkPoint.Points.filterGlobal(false)
			}
			if len(seriesMarks) > 0 {
				markPointPainter.add(markPointRenderOption{
					fillColor:          seriesColor,
					font:               getPreferredFont(series.Label.FontStyle.Font, opt.Font),
					symbolSize:         series.MarkPoint.SymbolSize,
					markpoints:         seriesMarks,
					seriesValues:       series.Values,
					points:             points,
					valueFormatter:     markPointValueFormatter,
					seriesLabelPainter: labelPainter,
				})
			}
			if len(globalMarks) > 0 {
				if globalSeriesData == nil {
					globalSeriesData = sumSeriesData(opt.SeriesList, series.YAxisIndex)
				}
				markPointPainter.add(markPointRenderOption{
					fillColor:          defaultGlobalMarkFillColor,
					font:               getPreferredFont(series.Label.FontStyle.Font, opt.Font),
					symbolSize:         series.MarkPoint.SymbolSize,
					markpoints:         globalMarks,
					seriesValues:       globalSeriesData,
					points:             points,
					valueFormatter:     markPointValueFormatter,
					seriesLabelPainter: labelPainter,
				})
			}
		}
	}

	if err := doRender(rendererList...); err != nil {
		return BoxZero, err
	}
	return p.box, nil
}

func (b *barChart) Render() (Box, error) {
	p := b.p
	opt := b.opt
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.theme)
	}
	if opt.Legend.Symbol == "" {
		// default to rectangle symbol for this chart type
		opt.Legend.Symbol = SymbolSquare
	}

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:          opt.Theme,
		padding:        opt.Padding,
		seriesList:     opt.SeriesList,
		stackSeries:    flagIs(true, opt.StackSeries),
		xAxis:          &b.opt.XAxis,
		yAxis:          opt.YAxis,
		title:          opt.Title,
		legend:         &b.opt.Legend,
		valueFormatter: opt.ValueFormatter,
	})
	if err != nil {
		return BoxZero, err
	}
	return b.renderChart(renderResult)
}
