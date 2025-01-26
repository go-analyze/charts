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
		YAxis:          make([]YAxisOption, sl.getYAxisCount()),
		ValueFormatter: defaultValueFormatter,
	}
}

type HorizontalBarChartOption struct {
	// Theme specifies the colors used for the chart.
	Theme ColorPalette
	// Padding specifies the padding of bar chart.
	Padding Box
	// Font is the font used to render the chart.
	Font *truetype.Font
	// SeriesList provides the data series.
	SeriesList SeriesList
	// StackSeries if set to *true a single bar with the colored series stacked together will be rendered.
	// This feature will result in some options being ignored, including BarMargin and SeriesLabelPosition.
	// MarkLine is also interpreted differently, only the first Series will have the MarkLine rendered (as it's the
	// base bar, other bars are influenced by prior values).
	StackSeries *bool
	// SeriesLabelPosition specifies the position of the label for the series. Currently supported values are
	// "left" or "right".
	SeriesLabelPosition string
	// XAxis are options for the x-axis.
	XAxis XAxisOption
	// YAxis are options for the y-axis (at most two).
	YAxis []YAxisOption
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

func (h *horizontalBarChart) render(result *defaultRenderResult, seriesList SeriesList) (Box, error) {
	p := h.p
	opt := h.opt
	seriesPainter := result.seriesPainter
	yRange := result.axisRanges[0]
	y0, y1 := yRange.GetRange(0)
	height := int(y1 - y0)
	seriesCount := len(seriesList)
	if seriesCount == 0 {
		return BoxZero, errors.New("empty series list")
	}
	stackedSeries := flagIs(true, opt.StackSeries)
	min, max, sumMax := seriesList.getMinMaxSumMax(0, stackedSeries)
	// If stacking, keep track of accumulated widths for each data index (after the “reverse” logic).
	var accumulatedWidths []int
	var margin, barMargin, barHeight int
	if stackedSeries {
		// If stacking, max should be the highest sum
		max = sumMax
		accumulatedWidths = make([]int, yRange.divideCount)
		margin, _, barHeight = calculateBarMarginsAndSize(1, height, opt.BarHeight, nil)
	} else {
		margin, barMargin, barHeight = calculateBarMarginsAndSize(seriesCount, height, opt.BarHeight, opt.BarMargin)
	}

	seriesNames := seriesList.Names()
	// xRange is used to convert data values into horizontal bar widths
	xRange := newRange(p, getPreferredValueFormatter(opt.XAxis.ValueFormatter, opt.ValueFormatter),
		seriesPainter.Width(), len(seriesList[0].Data), min, max, 1.0, 1.0)
	divideValues := yRange.AutoDivide()

	var rendererList []renderer
	for index := range seriesList {
		series := seriesList[index]
		seriesColor := opt.Theme.GetSeriesColor(series.index)

		var labelPainter *seriesLabelPainter
		if flagIs(true, series.Label.Show) {
			labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, series.Label, opt.Theme, opt.Font)
			rendererList = append(rendererList, labelPainter)
		}

		for j, item := range series.Data {
			if j >= yRange.divideCount {
				continue
			}
			// Reverse the category index for drawing from top to bottom
			reversedJ := yRange.divideCount - j - 1

			// Compute the top of this bar “row”
			y := divideValues[reversedJ] + margin

			// Determine the width (horizontal length) of the bar based on the data value
			w := xRange.getHeight(item)

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
				labelLeft := (opt.SeriesLabelPosition == PositionLeft || series.Label.Position == PositionLeft) && !stackedSeries
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

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:          opt.Theme,
		padding:        opt.Padding,
		seriesList:     opt.SeriesList,
		stackSeries:    flagIs(true, opt.StackSeries),
		xAxis:          &h.opt.XAxis,
		yAxis:          opt.YAxis,
		title:          opt.Title,
		legend:         &h.opt.Legend,
		valueFormatter: opt.ValueFormatter,
		axisReversed:   true,
	})
	if err != nil {
		return BoxZero, err
	}
	seriesList := opt.SeriesList.Filter(ChartTypeHorizontalBar)
	return h.render(renderResult, seriesList)
}
