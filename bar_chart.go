package charts

import (
	"errors"
	"math"
)

type barChart struct {
	p   *Painter
	opt *BarChartOption
}

// newBarChart returns a bar chart renderer.
func newBarChart(p *Painter, opt BarChartOption) *barChart {
	return &barChart{
		p:   p,
		opt: &opt,
	}
}

// NewBarChartOptionWithData returns an initialized BarChartOption with the SeriesList set with the provided data slice.
func NewBarChartOptionWithData(data [][]float64) BarChartOption {
	return NewBarChartOptionWithSeries(NewSeriesListBar(data))
}

// NewBarChartOptionWithSeries returns an initialized BarChartOption with the provided SeriesList.
func NewBarChartOptionWithSeries(sl BarSeriesList) BarChartOption {
	return BarChartOption{
		SeriesList:     sl,
		Padding:        defaultPadding,
		Theme:          GetDefaultTheme(),
		ValueAxis:      make([]ValueAxisOption, getSeriesYAxisCount(sl)),
		ValueFormatter: defaultValueFormatter,
	}
}

// BarChartOption defines the options for rendering a bar chart. Render the chart using Painter.BarChart.
type BarChartOption struct {
	// Theme specifies the colors used for the bar chart.
	Theme ColorPalette
	// Padding specifies the padding around the chart.
	Padding Box
	// Horizontal selects horizontal bar orientation when true, swapping the category and value axis.
	Horizontal bool
	// ValueAxis configures the value (numeric) axis. This axis represents the series values.
	ValueAxis []ValueAxisOption
	// CategoryAxis configures the category axis. This axis represents the series group.
	CategoryAxis CategoryAxisOption
	// SeriesList provides the data population for the chart. Typically constructed using NewSeriesListBar.
	SeriesList BarSeriesList
	// StackSeries when *true renders series stacked within one bar.
	// This ignores some options including BarMargin and SeriesLabelPosition.
	// MarkLine only renders for the first series and stacking only applies to the first y-axis.
	StackSeries *bool
	// SeriesLabelPosition specifies the label position for the series.
	// Vertical bars: "top" or "bottom". Horizontal bars: "left" or "right".
	SeriesLabelPosition string
	// Title contains options for rendering the chart title.
	Title TitleOption
	// Legend contains options for the data legend.
	Legend LegendOption
	// BarSize sets each bar's thickness as a ratio of the slot space allotted to it
	// (0.0–1.0, auto by default). With multiple series the slot is shared evenly, so the
	// ratio applies per bar. Vertical bars scale width, horizontal bars scale height.
	BarSize float64
	// BarMargin sets the spacing between grouped bars as a ratio of the category slot
	// (0.0–1.0, auto by default). BarSize takes priority over a set margin.
	BarMargin *float64
	// RoundedBarCaps when *true draws bars with rounded corners on the value-end of the bar.
	RoundedBarCaps *bool
	// ValueFormatter defines how float values are rendered to strings, notably for numeric axis labels.
	ValueFormatter ValueFormatter
}

// normalizeBarAxisPositions silently normalizes unsupported position values to their defaults.
func normalizeBarAxisPositions(horizontal bool, catAxis *CategoryAxisOption, valAxes []ValueAxisOption) {
	if horizontal {
		// category axis on Y: left and right supported, others normalize to left
		if p := catAxis.Position; p != "" && p != PositionLeft && p != PositionRight {
			catAxis.Position = PositionLeft
		}
		// value axis on X: only bottom supported
		// TODO - top-positioned value axis rendering is not yet supported
		for i := range valAxes {
			if p := valAxes[i].Position; p != "" && p != PositionBottom {
				valAxes[i].Position = PositionBottom
			}
		}
	} else {
		// category axis on X: only bottom supported
		// TODO - top-positioned category axis rendering for vertical bars is not yet supported
		if p := catAxis.Position; p != "" && p != PositionBottom {
			catAxis.Position = PositionBottom
		}
		// value axis on Y: left/right supported
		for i := range valAxes {
			p := valAxes[i].Position
			if p == "" {
				continue
			}
			if p != PositionLeft && p != PositionRight {
				if i == 0 {
					valAxes[i].Position = PositionLeft
				} else {
					valAxes[i].Position = PositionRight
				}
			}
		}
	}
}

// calculateGroupMarginsAndSize returns the group margin, inter-element margin, and element
// size in pixels for seriesCount elements sharing a slot of the given pixel space, honoring
// the optional configured pixel size and margin.
func calculateGroupMarginsAndSize(seriesCount, space int, configuredSize int, configuredMargin *float64) (int, int, int) {
	// default margins, adjusted below with config and series count
	margin := 10       // margin between each series group
	elementMargin := 5 // margin between each element
	if space < 20 {
		margin = 2
		elementMargin = 2
	} else if space < 50 {
		margin = 5
		elementMargin = 3
	}
	// check margin configuration if element size allows margin
	if configuredSize+elementMargin < space/seriesCount {
		// element size is in range that we should also consider an optional margin configuration
		if configuredMargin != nil {
			elementMargin = int(math.Round(*configuredMargin))
			if elementMargin+configuredSize > space/seriesCount {
				elementMargin = (space / seriesCount) - configuredSize
			}
		}
	} // else, element size is out of range.  Ignore margin config

	size := (space - 2*margin - elementMargin*(seriesCount-1)) / seriesCount
	// check size configuration, limited by the series count and space available
	if configuredSize > 0 && configuredSize < size {
		size = configuredSize
		// recalculate margin
		margin = (space - seriesCount*size - elementMargin*(seriesCount-1)) / 2
	}

	return margin, elementMargin, size
}

// resolveBarSizePixels converts a slot-ratio bar size to a per-bar pixel size, returning 0 (auto) when unset.
func resolveBarSizePixels(barSize float64, space, count int) int {
	if barSize <= 0 {
		return 0
	}
	return int(float64(space) * barSize / float64(count))
}

// resolveBarMarginPixels converts a slot-ratio margin to pixels, returning nil (auto) when unset.
func resolveBarMarginPixels(barMargin *float64, space int) *float64 {
	if barMargin == nil {
		return nil
	}
	px := float64(space) * *barMargin
	return &px
}

func (b *barChart) renderChart(result *defaultRenderResult) (Box, error) {
	if len(b.opt.SeriesList) == 0 {
		result.renderNoData(b.opt.Theme)
		return b.p.box, nil
	}
	if b.opt.Horizontal {
		return b.renderHorizontalBars(result)
	}
	return b.renderVerticalBars(result)
}

func (b *barChart) renderVerticalBars(result *defaultRenderResult) (Box, error) {
	p := b.p
	opt := b.opt
	seriesCount := len(opt.SeriesList)
	seriesPainter := result.seriesPainter

	x0, x1 := result.categoryAxisRange.getRange(0)
	width := int(x1 - x0)
	barMaxHeight := seriesPainter.Height() // total vertical space for bars
	seriesNames := opt.SeriesList.names()
	divideValues := result.categoryAxisRange.autoDivide()
	stackedSeries := flagIs(true, opt.StackSeries)
	barSize := opt.BarSize
	var margin, barMargin, barWidth int
	var accumulatedHeights []int // prior heights for stacking to avoid recalculating the heights
	if stackedSeries {
		barCount := getSeriesYAxisCount(opt.SeriesList) // only two bars if two y-axis
		configuredMargin := opt.BarMargin
		if barCount == 1 {
			configuredMargin = nil // no margin needed with a single bar
		}
		margin, _, barWidth = calculateGroupMarginsAndSize(barCount, width,
			resolveBarSizePixels(barSize, width, barCount), resolveBarMarginPixels(configuredMargin, width))
		accumulatedHeights = make([]int, result.categoryAxisRange.divideCount)
	} else {
		margin, barMargin, barWidth = calculateGroupMarginsAndSize(seriesCount, width,
			resolveBarSizePixels(barSize, width, seriesCount), resolveBarMarginPixels(opt.BarMargin, width))
	}

	markPointPainter := newMarkPointPainter(seriesPainter)
	markLinePainter := newMarkLinePainter(seriesPainter)
	// render list must start with the markPointPainter, as it can influence label painters (if enabled)
	rendererList := []renderer{markPointPainter, markLinePainter}

	for index, series := range opt.SeriesList {
		stackSeries := stackedSeries && series.YAxisIndex == 0
		yRange := result.valueAxisRanges[series.YAxisIndex]
		seriesThemeIndex := index
		if series.absThemeIndex != nil {
			seriesThemeIndex = *series.absThemeIndex
		}
		seriesColor := opt.Theme.GetSeriesColor(seriesThemeIndex)

		var labelPainter *seriesLabelPainter
		if flagIs(true, series.Label.Show) {
			labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, series.Label, opt.Theme, opt.Padding.Right)
			rendererList = append(rendererList, labelPainter)
		}

		points := make([]Point, len(series.Values)) // used for mark points
		for j, item := range series.Values {
			if j >= result.categoryAxisRange.divideCount {
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
					barWidth, roundTopLeft|roundTopRight, seriesColor, seriesColor, 0.0)
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
				labelPainter.Add(labelValue{
					vertical:  true, // label is vertically oriented
					index:     index,
					dataIndex: j,
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
					font:           series.Label.FontStyle.Font,
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
					font:           series.Label.FontStyle.Font,
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
					font:               series.Label.FontStyle.Font,
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
					font:               series.Label.FontStyle.Font,
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

	// TODO - support dual horizontal value axis when top x-axis rendering is added
	if opt.Horizontal && len(opt.ValueAxis) > 1 {
		return BoxZero, errors.New("dual value axes with horizontal bars is not supported")
	}
	valueAxis := opt.ValueAxis
	if len(valueAxis) == 0 {
		valueAxis = []ValueAxisOption{{}}
	}
	categoryAxis := opt.CategoryAxis
	normalizeBarAxisPositions(opt.Horizontal, &categoryAxis, valueAxis)

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:          opt.Theme,
		padding:        opt.Padding,
		seriesList:     opt.SeriesList,
		stackSeries:    flagIs(true, opt.StackSeries),
		categoryAxis:   &categoryAxis,
		valueAxis:      valueAxis,
		title:          opt.Title,
		legend:         &b.opt.Legend,
		valueFormatter: opt.ValueFormatter,
		categoryY:      opt.Horizontal,
	})
	if err != nil {
		return BoxZero, err
	}
	return b.renderChart(renderResult)
}

func (b *barChart) renderHorizontalBars(result *defaultRenderResult) (Box, error) {
	p := b.p
	opt := b.opt
	seriesCount := len(opt.SeriesList)
	seriesPainter := result.seriesPainter
	yRange := result.categoryAxisRange
	y0, y1 := yRange.getRange(0)
	height := int(y1 - y0)
	stackedSeries := flagIs(true, opt.StackSeries)
	// TODO - propagate per-axis reversed flag once horizontal bars support multiple value axes
	reversed := result.valueAxisRanges[0].reversed // bars grow from the right when the category axis is on the right
	plotWidth := seriesPainter.Width()
	// baselineX is the value=0 edge; dir is the direction bars grow.
	baselineX, dir := 0, 1
	if reversed {
		baselineX, dir = plotWidth, -1
	}
	barSize := opt.BarSize

	// if stacking, keep track of accumulated widths for each data index (after the "reverse" logic)
	var accumulatedWidths []int
	var margin, barMargin, barHeight int
	if stackedSeries {
		accumulatedWidths = make([]int, yRange.divideCount)
		margin, _, barHeight = calculateGroupMarginsAndSize(1, height,
			resolveBarSizePixels(barSize, height, 1), nil)
	} else {
		margin, barMargin, barHeight = calculateGroupMarginsAndSize(seriesCount, height,
			resolveBarSizePixels(barSize, height, seriesCount), resolveBarMarginPixels(opt.BarMargin, height))
	}

	seriesNames := opt.SeriesList.names()
	divideValues := yRange.autoDivide()

	markPointPainter := newMarkPointPainter(seriesPainter)
	markLinePainter := newMarkLinePainter(seriesPainter)
	// render list must start with the markPointPainter, as it can influence label painters (if enabled)
	rendererList := []renderer{markPointPainter, markLinePainter}

	for index, series := range opt.SeriesList {
		seriesThemeIndex := index
		if series.absThemeIndex != nil {
			seriesThemeIndex = *series.absThemeIndex
		}
		seriesColor := opt.Theme.GetSeriesColor(seriesThemeIndex)

		var labelPainter *seriesLabelPainter
		if flagIs(true, series.Label.Show) {
			labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, series.Label, opt.Theme, opt.Padding.Right)
			rendererList = append(rendererList, labelPainter)
		}

		points := make([]Point, len(series.Values))
		for j, item := range series.Values {
			if j >= yRange.divideCount {
				break
			}
			// Reverse the category index for drawing from top to bottom
			reversedJ := yRange.divideCount - j - 1

			// Compute the top of this bar "row"
			y := divideValues[reversedJ] + margin

			// Determine the width (horizontal length) of the bar based on the data value
			w := result.valueAxisRanges[0].getHeight(item)

			// stackBase is the bar's category-axis-side edge; tipX is the value-end edge.
			var stackBase, tipX int
			if stackedSeries {
				stackBase = baselineX + dir*accumulatedWidths[reversedJ]
				accumulatedWidths[reversedJ] += w
			} else {
				// Offset each series in its own lane
				if index != 0 {
					y += index * (barHeight + barMargin)
				}
				stackBase = baselineX
			}
			tipX = stackBase + dir*w
			left, right := stackBase, tipX
			if left > right {
				left, right = right, left
			}

			// In stacked mode, only round caps on the last series
			roundedCorners := roundTopRight | roundBottomRight
			if reversed {
				roundedCorners = roundTopLeft | roundBottomLeft
			}
			if flagIs(true, opt.RoundedBarCaps) && (!stackedSeries || index == seriesCount-1) {
				seriesPainter.roundedRect(
					Box{Top: y, Left: left, Right: right, Bottom: y + barHeight, IsSet: true},
					barHeight, roundedCorners, seriesColor, seriesColor, 0.0)
			} else {
				seriesPainter.FilledRect(left, y, right, y+barHeight, seriesColor, seriesColor, 0.0)
			}

			// Prepare point for mark points (anchor at the bar's value-end)
			points[j] = Point{
				X: tipX,
				Y: y + (barHeight >> 1), // vertical center of bar
			}

			if labelPainter != nil {
				fontStyle := series.Label.FontStyle
				labelX := tipX
				labelY := y + (barHeight >> 1)
				labelLeft := opt.SeriesLabelPosition == PositionLeft && !stackedSeries
				if labelLeft {
					labelX = baselineX // anchor to the category-axis side
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
					dataIndex: j,
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
					font:           series.Label.FontStyle.Font,
					marklines:      seriesMarks,
					seriesValues:   series.Values,
					axisRange:      result.valueAxisRanges[0],
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
					font:           series.Label.FontStyle.Font,
					marklines:      globalMarks,
					seriesValues:   globalSeriesData,
					axisRange:      result.valueAxisRanges[0],
					valueFormatter: markLineValueFormatter,
				})
			}
		}
		if len(series.MarkPoint.Points) > 0 {
			markPointValueFormatter := getPreferredValueFormatter(series.MarkPoint.ValueFormatter,
				series.Label.ValueFormatter, opt.ValueFormatter)
			// pin tail points toward the data point.
			// With y-down screen coordinates, +pi/2 rotates the default down pin to left pointing.
			markPointRotation := math.Pi / 2
			if reversed {
				markPointRotation = -math.Pi / 2
			}
			var seriesMarks, globalMarks SeriesMarkList
			if stackedSeries && index == seriesCount-1 { // global is only allowed when stacked and on the last series
				seriesMarks, globalMarks = series.MarkPoint.Points.splitGlobal()
			} else {
				seriesMarks = series.MarkPoint.Points.filterGlobal(false)
			}
			if len(seriesMarks) > 0 {
				markPointPainter.add(markPointRenderOption{
					fillColor:          seriesColor,
					font:               series.Label.FontStyle.Font,
					symbolSize:         series.MarkPoint.SymbolSize,
					rotationRadians:    markPointRotation,
					markpoints:         seriesMarks,
					seriesValues:       series.Values,
					points:             points,
					valueFormatter:     markPointValueFormatter,
					seriesLabelPainter: labelPainter,
				})
			}
			if len(globalMarks) > 0 {
				if globalSeriesData == nil {
					globalSeriesData = sumSeriesData(opt.SeriesList, 0)
				}
				markPointPainter.add(markPointRenderOption{
					fillColor:          defaultGlobalMarkFillColor,
					font:               series.Label.FontStyle.Font,
					symbolSize:         series.MarkPoint.SymbolSize,
					rotationRadians:    markPointRotation,
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
