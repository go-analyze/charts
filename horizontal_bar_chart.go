package charts

import (
	"errors"

	"github.com/golang/freetype/truetype"
)

type horizontalBarChart struct {
	p   *Painter
	opt *HorizontalBarChartOption
}

// NewHorizontalBarChartOptionWithData returns an initialized HorizontalBarChartOption with the SeriesList set for the provided data slice.
func NewHorizontalBarChartOptionWithData(data [][]float64) HorizontalBarChartOption {
	sl := NewSeriesListHorizontalBar(data)
	return HorizontalBarChartOption{
		SeriesList:     sl,
		Padding:        defaultPadding,
		Theme:          GetDefaultTheme(),
		Font:           GetDefaultFont(),
		ValueFormatter: defaultValueFormatter,
	}
}

// HorizontalBarChartOption defines the options for rendering a horizontal bar chart.
// Render the chart using Painter.HorizontalBarChart.
type HorizontalBarChartOption struct {
	// Theme specifies the colors used for the chart.
	Theme ColorPalette
	// Padding specifies the padding around the chart.
	Padding Box
	// Deprecated: Font is deprecated, instead the font needs to be set on the SeriesLabel, or other specific elements.
	Font *truetype.Font
	// SeriesList provides the data population for the chart, typically constructed using NewSeriesListHorizontalBar.
	SeriesList HorizontalBarSeriesList
	// StackSeries, if true, renders the series stacked within one bar. This
	// causes some options, including BarMargin and SeriesLabelPosition, to be
	// ignored. MarkLine only renders for the first series.
	StackSeries *bool
	// SeriesLabelPosition specifies the position of the label for the series. Currently supported values are
	// "left" or "right".
	SeriesLabelPosition string
	// XAxis are options for the x-axis.
	XAxis XAxisOption
	// YAxis are options for the y-axis.
	YAxis YAxisOption
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// BarHeight specifies the height of each horizontal bar. Height may be reduced to ensure all series fit on the chart.
	BarHeight int
	// BarMargin specifies the margin between bars grouped together. BarHeight takes priority over the margin.
	BarMargin *float64
	// ValueFormatter defines how float values should be rendered to strings, notably for numeric axis labels.
	ValueFormatter ValueFormatter
}

// newHorizontalBarChart returns a horizontal bar chart renderer
func newHorizontalBarChart(p *Painter, opt HorizontalBarChartOption) *horizontalBarChart {
	return &horizontalBarChart{
		p:   p,
		opt: &opt,
	}
}

func (h *horizontalBarChart) renderChart(result *defaultRenderResult) (Box, error) {
	p := h.p
	opt := h.opt
	seriesCount := len(opt.SeriesList)
	if seriesCount == 0 {
		return BoxZero, errors.New("empty series list")
	}
	seriesPainter := result.seriesPainter
	yRange := result.yaxisRanges[0]
	y0, y1 := yRange.getRange(0)
	height := int(y1 - y0)
	stackedSeries := flagIs(true, opt.StackSeries)
	// If stacking, keep track of accumulated widths for each data index (after the “reverse” logic).
	var accumulatedWidths []int
	var margin, barMargin, barHeight int
	if stackedSeries {
		accumulatedWidths = make([]int, yRange.divideCount)
		margin, _, barHeight = calculateBarMarginsAndSize(1, height, opt.BarHeight, nil)
	} else {
		margin, barMargin, barHeight = calculateBarMarginsAndSize(seriesCount, height, opt.BarHeight, opt.BarMargin)
	}

	seriesNames := opt.SeriesList.names()
	divideValues := yRange.autoDivide()

	markLinePainter := newMarkLinePainter(seriesPainter)
	rendererList := []renderer{markLinePainter}

	for index, series := range opt.SeriesList {
		seriesColor := opt.Theme.GetSeriesColor(index)

		var labelPainter *seriesLabelPainter
		if flagIs(true, series.Label.Show) {
			labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, series.Label, opt.Theme)
			rendererList = append(rendererList, labelPainter)
		}

		for j, item := range series.Values {
			if j >= yRange.divideCount {
				break
			}
			// Reverse the category index for drawing from top to bottom
			reversedJ := yRange.divideCount - j - 1

			// Compute the top of this bar “row”
			y := divideValues[reversedJ] + margin

			// Determine the width (horizontal length) of the bar based on the data value
			w := result.xaxisRange.getHeight(item)

			var left, right int
			if stackedSeries {
				// Start where the previous series ended
				left = accumulatedWidths[reversedJ]
				right = left + w
				accumulatedWidths[reversedJ] = right
			} else {
				// Offset each series in its own lane
				if index != 0 {
					y += index * (barHeight + barMargin)
				}
				left = 0
				right = w
			}

			seriesPainter.FilledRect(left, y, right, y+barHeight, seriesColor, seriesColor, 0.0)

			if labelPainter != nil {
				fontStyle := series.Label.FontStyle
				labelX := right
				labelY := y + (barHeight >> 1)
				labelLeft := opt.SeriesLabelPosition == PositionLeft && !stackedSeries
				if labelLeft {
					labelX = 0
				}
				if fontStyle.FontColor.IsZero() {
					var testColor Color
					if labelLeft {
						testColor = seriesColor
					} else if stackedSeries && index+1 < seriesCount {
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
					vertical:  false, // horizontal label
					index:     index,
					value:     item,
					x:         labelX,
					y:         labelY,
					offset:    series.Label.Offset,
					fontStyle: fontStyle,
				})
			}
		}

		var globalSeriesData []float64 // lazily initialized
		if len(series.MarkLine.Lines) > 0 {
			markLineValueFormatter := getPreferredValueFormatter(series.MarkLine.ValueFormatter,
				series.Label.ValueFormatter, opt.ValueFormatter)
			var seriesMarks, globalMarks SeriesMarkList
			if stackedSeries && index == seriesCount-1 { // global is only allowed when stacked and on the last series
				seriesMarks, globalMarks = series.MarkLine.Lines.splitGlobal()
			} else {
				seriesMarks = series.MarkLine.Lines.filterGlobal(false)
			}
			if len(seriesMarks) > 0 && (!stackedSeries || index == 0) {
				// in stacked mode we only support the line painter for the first series
				markLinePainter.add(markLineRenderOption{
					verticalLine:   true,
					fillColor:      seriesColor,
					fontColor:      opt.Theme.GetMarkTextColor(),
					strokeColor:    seriesColor,
					font:           opt.Font,
					marklines:      seriesMarks,
					seriesValues:   series.Values,
					axisRange:      result.xaxisRange,
					valueFormatter: markLineValueFormatter,
				})
			}
			if len(globalMarks) > 0 {
				if globalSeriesData == nil {
					globalSeriesData = sumSeriesData(opt.SeriesList, 0)
				}
				markLinePainter.add(markLineRenderOption{
					verticalLine:   true,
					fillColor:      defaultGlobalMarkFillColor,
					fontColor:      opt.Theme.GetMarkTextColor(),
					strokeColor:    defaultGlobalMarkFillColor,
					font:           opt.Font,
					marklines:      globalMarks,
					seriesValues:   globalSeriesData,
					axisRange:      result.xaxisRange,
					valueFormatter: markLineValueFormatter,
				})
			}
		}
	}

	if err := doRender(rendererList...); err != nil {
		return BoxZero, err
	}
	return p.box, nil
}

func (h *horizontalBarChart) Render() (Box, error) {
	p := h.p
	opt := h.opt
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
		xAxis:          &h.opt.XAxis,
		yAxis:          []YAxisOption{opt.YAxis},
		title:          opt.Title,
		legend:         &h.opt.Legend,
		valueFormatter: opt.ValueFormatter,
		axisReversed:   true,
	})
	if err != nil {
		return BoxZero, err
	}
	return h.renderChart(renderResult)
}
