package charts

import (
	"errors"
	"math"
)

type candlestickChart struct {
	p   *Painter
	opt *CandlestickChartOption
}

// newCandlestickChart returns a candlestick chart renderer.
func newCandlestickChart(p *Painter, opt CandlestickChartOption) *candlestickChart {
	return &candlestickChart{
		p:   p,
		opt: &opt,
	}
}

// CandlestickChartOption defines options for rendering candlestick charts. Render the chart using Painter.CandlestickChart.
type CandlestickChartOption struct {
	// Theme specifies the colors used for the candlestick chart.
	Theme ColorPalette
	// Padding specifies the padding around the chart.
	Padding Box
	// SeriesList provides the OHLC data population for the chart. Typically constructed using NewSeriesListCandlestick.
	SeriesList CandlestickSeriesList
	// XAxis contains options for the x-axis.
	XAxis XAxisOption
	// YAxis contains options for the y-axis. At most two y-axes are supported.
	YAxis []YAxisOption
	// Title contains options for rendering the chart title.
	Title TitleOption
	// Legend contains options for the data legend.
	Legend LegendOption
	// CandleWidth sets body width ratio (0.0–1.0, default 0.8).
	CandleWidth float64
	// ShowWicks controls whether high-low wicks are displayed by default. When nil, wicks are shown.
	// Individual series can override this setting.
	ShowWicks *bool
	// WickWidth sets wick stroke width in pixels (default 1.0).
	WickWidth float64
	// CandleMargin sets inter-series spacing ratio (0.0–1.0, auto by default).
	// Only applies with multiple candlestick series.
	CandleMargin *float64
	// ValueFormatter formats numeric values.
	ValueFormatter ValueFormatter
}

// NewCandlestickOptionWithData creates a CandlestickChartOption from OHLC data slices.
func NewCandlestickOptionWithData(data ...[]OHLCData) CandlestickChartOption {
	seriesList := make(CandlestickSeriesList, len(data))
	for i, ohlcData := range data {
		seriesList[i] = CandlestickSeries{Data: ohlcData}
	}
	return NewCandlestickOptionWithSeries(seriesList...)
}

// NewCandlestickOptionWithSeries returns an initialized CandlestickChartOption with the provided Series.
func NewCandlestickOptionWithSeries(series ...CandlestickSeries) CandlestickChartOption {
	seriesList := make(CandlestickSeriesList, len(series))
	copy(seriesList, series)
	return CandlestickChartOption{
		SeriesList:     seriesList,
		Padding:        defaultPadding,
		Theme:          GetDefaultTheme(),
		YAxis:          make([]YAxisOption, getSeriesYAxisCount(seriesList)),
		ValueFormatter: defaultValueFormatter,
		CandleWidth:    0.8, // Default 80% of available space
		WickWidth:      1.0,
	}
}

// createPatternAwareLabelFormatter creates a label formatter that can handle pattern detection
// while respecting user-provided label formatters based on Replace/Complement mode
func createPatternAwareLabelFormatter(originalSeries *CandlestickSeries, seriesIndex int, theme ColorPalette,
	patternMap map[int][]PatternDetectionResult) SeriesLabelFormatter {
	return func(index int, name string, val float64) (string, *LabelStyle) {
		// Check for patterns at this index using pre-computed map
		var patterns []PatternDetectionResult
		if patternMap != nil {
			patterns = patternMap[index]
		}
		if !originalSeries.PatternConfig.PreferPatternLabels && originalSeries.Label.LabelFormatter != nil &&
			flagIs(true, originalSeries.Label.Show) { // prefer user label if provided
			userText, userStyle := originalSeries.Label.LabelFormatter(index, name, val)
			if userText != "" {
				return userText, userStyle // User label takes precedence
			}
		}
		var patternText string
		var patternStyle *LabelStyle
		if originalSeries.PatternConfig.PatternFormatter != nil {
			patternText, patternStyle = originalSeries.PatternConfig.PatternFormatter(patterns, originalSeries.Name, val)
		} else {
			patternText, patternStyle = formatPatternsDefault(patterns, seriesIndex, theme)
		}
		if patternText != "" {
			// either no user input, or configured to replace and we have a matching pattern
			return patternText, patternStyle
		}

		if flagIs(true, originalSeries.Label.Show) { // user explicitly requested a label, show something
			if originalSeries.Label.LabelFormatter != nil {
				return originalSeries.Label.LabelFormatter(index, name, val)
			} else if originalSeries.Label.ValueFormatter != nil {
				return originalSeries.Label.ValueFormatter(val), nil
			}
			return defaultValueFormatter(val), nil
		}
		return "", nil
	}
}

// TODO - v0.6 - attempt to de-duplicate with calculateBarMarginsAndSize
// calculateCandleMarginsAndSize calculates margins and candle sizes similar to bar charts.
func calculateCandleMarginsAndSize(seriesCount, space int, configuredCandleSize int, configuredCandleMargin *float64) (int, int, int) {
	// default margins, adjusted below with config and series count
	margin := 10      // margin between each series group
	candleMargin := 5 // margin between each candle
	if space < 20 {
		margin = 2
		candleMargin = 2
	} else if space < 50 {
		margin = 5
		candleMargin = 3
	}
	// check margin configuration if candle size allows margin
	if configuredCandleSize+candleMargin < space/seriesCount {
		// CandleWidth is in range that we should also consider an optional margin configuration
		if configuredCandleMargin != nil {
			candleMargin = int(math.Round(*configuredCandleMargin))
			if candleMargin+configuredCandleSize > space/seriesCount {
				candleMargin = (space / seriesCount) - configuredCandleSize
			}
		}
	} // else, candle width is out of range.  Ignore margin config

	candleSize := (space - 2*margin - candleMargin*(seriesCount-1)) / seriesCount
	// check candle size configuration, limited by the series count and space available
	if configuredCandleSize > 0 && configuredCandleSize < candleSize {
		candleSize = configuredCandleSize
		// recalculate margin
		margin = (space - seriesCount*candleSize - candleMargin*(seriesCount-1)) / 2
	}

	return margin, candleMargin, candleSize
}

func (k *candlestickChart) renderChart(result *defaultRenderResult) (Box, error) {
	p := k.p
	opt := k.opt
	seriesList := opt.SeriesList
	if seriesList.len() == 0 {
		return BoxZero, errors.New("empty series list")
	}
	seriesPainter := result.seriesPainter

	// Find maximum data count across all series
	maxDataCount := getSeriesMaxDataCount(seriesList)
	if maxDataCount == 0 {
		return BoxZero, errors.New("no data in any series")
	}
	width := seriesPainter.Width()
	if width <= 0 {
		return BoxZero, errors.New("invalid painter width")
	}
	seriesCount := seriesList.len()

	// Calculate candle width using CandleWidth ratio (default 80%)
	candleWidthRatio := opt.CandleWidth
	if candleWidthRatio <= 0 {
		candleWidthRatio = 0.8 // Default 80% of available space
	} else if candleWidthRatio > 1 {
		candleWidthRatio = 1
	}
	candleWidth := int(float64(width) * candleWidthRatio / float64(maxDataCount))
	if candleWidth < 1 {
		candleWidth = 1
	}

	// Calculate candleWidthPerSeries for body rendering
	candleWidthPerSeries := candleWidth / seriesCount
	if candleWidthPerSeries < 1 {
		candleWidthPerSeries = 1
	}

	// Use autoDivide for positioning
	divideValues := result.xaxisRange.autoDivide()

	// Center positions for each series index
	seriesCenterValues := make([][]int, seriesList.len())

	// render list must start with the markPointPainter, as it can influence label painters (if enabled)
	markPointPainter := newMarkPointPainter(seriesPainter)
	markLinePainter := newMarkLinePainter(seriesPainter)
	trendLinePainter := newTrendLinePainter(seriesPainter)
	rendererList := []renderer{markPointPainter, markLinePainter, trendLinePainter}

	seriesNames := seriesList.names()

	// Store points and label painters for each series
	seriesClosePoints := make([][]Point, seriesList.len())
	seriesOpenPoints := make([][]Point, seriesList.len())
	seriesHighPoints := make([][]Point, seriesList.len())
	seriesLowPoints := make([][]Point, seriesList.len())
	allLabelPainters := make([]*seriesLabelPainter, seriesList.len())

	// Render each series
	for seriesIndex := 0; seriesIndex < seriesList.len(); seriesIndex++ {
		series := seriesList.getSeries(seriesIndex).(*CandlestickSeries)

		// Bounds check for Y axis index to prevent panic
		if series.YAxisIndex >= len(result.yaxisRanges) {
			return BoxZero, errors.New("candlestick series YAxisIndex out of bounds")
		}
		yRange := result.yaxisRanges[series.YAxisIndex]

		// pre-compute patterns for this series
		var patternMap map[int][]PatternDetectionResult
		if series.PatternConfig != nil {
			patternMap = scanForCandlestickPatterns(series.Data, *series.PatternConfig)
		}

		// Create labelPainter only when labels are enabled or patterns were detected
		var labelPainter *seriesLabelPainter
		if flagIs(true, series.Label.Show) || len(patternMap) > 0 {
			if len(patternMap) > 0 {
				labelCopy := series.Label
				labelCopy.LabelFormatter = createPatternAwareLabelFormatter(series, seriesIndex, opt.Theme, patternMap)
				labelCopy.Show = Ptr(true) // Enable labels when patterns are detected, even if user labels are disabled
				labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, labelCopy, opt.Theme, opt.Padding.Right)
			} else {
				labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, series.Label, opt.Theme, opt.Padding.Right)
			}
			rendererList = append(rendererList, labelPainter)
		}
		allLabelPainters[seriesIndex] = labelPainter

		seriesThemeIndex := seriesIndex
		if series.absThemeIndex != nil {
			seriesThemeIndex = *series.absThemeIndex
		}
		upColor, downColor := opt.Theme.GetSeriesUpDownColors(seriesThemeIndex)

		// Initialize point arrays for all OHLC values for this series
		seriesClosePoints[seriesIndex] = make([]Point, len(series.Data))
		seriesOpenPoints[seriesIndex] = make([]Point, len(series.Data))
		seriesHighPoints[seriesIndex] = make([]Point, len(series.Data))
		seriesLowPoints[seriesIndex] = make([]Point, len(series.Data))
		seriesCenterValues[seriesIndex] = make([]int, len(series.Data))
		// Render each candlestick in this series
		for j, ohlc := range series.Data {
			if j >= maxDataCount || j >= len(divideValues) {
				continue
			}

			// Position calculation: center candlesticks in each time period
			// Calculate the center of each time period section
			var sectionWidth int
			if j < len(divideValues)-1 {
				sectionWidth = divideValues[j+1] - divideValues[j]
			} else {
				// Last section uses same width as previous
				if j > 0 {
					sectionWidth = divideValues[j] - divideValues[j-1]
				} else {
					sectionWidth = seriesPainter.Width() / maxDataCount
				}
			}

			// Calculate margins and positioning exactly like bar charts
			var groupMargin, candleMargin, candleWidth int
			if seriesList.len() == 1 {
				// Single series: use simple centering
				groupMargin = 0
				candleMargin = 0
				candleWidth = candleWidthPerSeries
			} else {
				// Multiple series: use bar chart margin calculation logic
				// Convert CandleMargin percentage to pixels
				var candleMarginFloat *float64
				if opt.CandleMargin != nil {
					// Convert percentage to absolute pixels
					marginPixels := float64(sectionWidth) * (*opt.CandleMargin)
					candleMarginFloat = &marginPixels
				}

				// Use bar chart logic: calculateBarMarginsAndSize equivalent
				groupMargin, candleMargin, candleWidth =
					calculateCandleMarginsAndSize(seriesList.len(),
						sectionWidth, candleWidthPerSeries, candleMarginFloat)
			}

			var centerX int
			if seriesList.len() == 1 {
				// Single series: center in the time period section
				centerX = divideValues[j] + sectionWidth/2
			} else {
				// Multiple series: use exact bar chart positioning formula
				// x = divideValues[j] + margin + index*(barWidth+barMargin)
				x := divideValues[j] + groupMargin + seriesIndex*(candleWidth+candleMargin)
				centerX = x + candleWidth/2
			}
			seriesCenterValues[seriesIndex][j] = centerX

			if !validateOHLCData(ohlc) { // if invalid mark as null and skip
				// Mark all OHLC points as invalid
				invalidPoint := Point{X: centerX, Y: math.MaxInt32}
				seriesClosePoints[seriesIndex][j] = invalidPoint
				seriesOpenPoints[seriesIndex][j] = invalidPoint
				seriesHighPoints[seriesIndex][j] = invalidPoint
				seriesLowPoints[seriesIndex][j] = invalidPoint
				continue
			}

			leftX := centerX - candleWidth/2
			rightX := centerX + candleWidth/2

			highY := yRange.getRestHeight(ohlc.High)
			lowY := yRange.getRestHeight(ohlc.Low)
			openY := yRange.getRestHeight(ohlc.Open)
			closeY := yRange.getRestHeight(ohlc.Close)

			bodyTop := int(math.Min(float64(openY), float64(closeY)))
			bodyBottom := int(math.Max(float64(openY), float64(closeY)))

			// Determine colors and style
			isBullish := ohlc.Close >= ohlc.Open
			candleStyle := series.CandleStyle
			if candleStyle == "" {
				candleStyle = CandleStyleFilled
			}

			var bodyColor, wickColor Color
			if isBullish {
				bodyColor = upColor
			} else {
				bodyColor = downColor
			}

			wickColor = opt.Theme.GetCandleWickColor()
			if wickColor.IsZero() {
				wickColor = bodyColor
			}

			// Draw high-low wick (if enabled)
			showWicks := !flagIs(false, opt.ShowWicks)
			if series.ShowWicks != nil {
				showWicks = *series.ShowWicks
			}
			wickWidth := opt.WickWidth
			if wickWidth <= 0 {
				wickWidth = 1.0
			}
			if showWicks {
				if highY < bodyTop {
					seriesPainter.LineStroke([]Point{
						{X: centerX, Y: highY},
						{X: centerX, Y: bodyTop},
					}, wickColor, wickWidth)
				}
				if lowY > bodyBottom {
					seriesPainter.LineStroke([]Point{
						{X: centerX, Y: bodyBottom},
						{X: centerX, Y: lowY},
					}, wickColor, wickWidth)
				}

				// Calculate cap width (based on series candle width)
				capWidth := candleWidthPerSeries / 4
				if capWidth < 1 {
					capWidth = 1
				}

				// Draw horizontal cap at high point
				seriesPainter.LineStroke([]Point{
					{X: centerX - capWidth, Y: highY},
					{X: centerX + capWidth, Y: highY},
				}, wickColor, wickWidth)

				// Draw horizontal cap at low point
				seriesPainter.LineStroke([]Point{
					{X: centerX - capWidth, Y: lowY},
					{X: centerX + capWidth, Y: lowY},
				}, wickColor, wickWidth)
			}

			// Draw open-close body based on style
			if bodyTop == bodyBottom { // Doji (open == close)
				// Draw thin line instead of rectangle
				seriesPainter.LineStroke([]Point{
					{X: leftX, Y: bodyTop},
					{X: rightX, Y: bodyTop},
				}, bodyColor, 1.0)
			} else {
				switch candleStyle {
				case CandleStyleFilled:
					seriesPainter.FilledRect(leftX, bodyTop, rightX, bodyBottom,
						bodyColor, bodyColor, 0.0)

				case CandleStyleTraditional:
					if isBullish { // Hollow body for bullish
						seriesPainter.FilledRect(leftX, bodyTop, rightX, bodyBottom,
							ColorTransparent, bodyColor, wickWidth)
					} else { // Filled body for bearish
						seriesPainter.FilledRect(leftX, bodyTop, rightX, bodyBottom,
							bodyColor, bodyColor, 0.0)
					}

				case CandleStyleOutline:
					seriesPainter.FilledRect(leftX, bodyTop, rightX, bodyBottom,
						ColorTransparent, bodyColor, wickWidth)
				}
			}

			// Store points for all OHLC values for mark points
			seriesClosePoints[seriesIndex][j] = Point{X: centerX, Y: closeY}
			seriesOpenPoints[seriesIndex][j] = Point{X: centerX, Y: openY}
			seriesHighPoints[seriesIndex][j] = Point{X: centerX, Y: highY}
			seriesLowPoints[seriesIndex][j] = Point{X: centerX, Y: lowY}

			// Add label if enabled (pattern logic is now handled in the label formatter)
			if labelPainter != nil {
				labelPainter.Add(labelValue{
					index:     j,          // Data point index (candlestick position), not series index
					value:     ohlc.Close, // Use close price for label
					x:         centerX,
					y:         closeY,
					fontStyle: series.Label.FontStyle,
					offset:    series.Label.Offset,
				})
			}
		}
	}

	// Handle mark lines, mark points, and trend lines for each series and OHLC component
	for seriesIndex := 0; seriesIndex < seriesList.len(); seriesIndex++ {
		series := seriesList.getSeries(seriesIndex).(*CandlestickSeries)

		// Bounds check for Y axis index to prevent panic
		if series.YAxisIndex >= len(result.yaxisRanges) {
			continue // Skip this series if YAxisIndex is out of bounds
		}
		yRange := result.yaxisRanges[series.YAxisIndex]
		seriesThemeIndex := seriesIndex
		if series.absThemeIndex != nil {
			seriesThemeIndex = *series.absThemeIndex
		}
		seriesColor := opt.Theme.GetSeriesColor(seriesThemeIndex)

		ohlcComponents := []struct {
			markLine    SeriesMarkLine
			markPoint   SeriesMarkPoint
			trendLines  []SeriesTrendLine
			extractFunc func(*CandlestickSeries) []float64
			points      []Point
		}{
			{
				markLine:    series.OpenMarkLine,
				markPoint:   series.OpenMarkPoint,
				trendLines:  series.OpenTrendLine,
				extractFunc: (*CandlestickSeries).ExtractOpenPrices,
				points:      seriesOpenPoints[seriesIndex],
			},
			{
				markLine:    series.HighMarkLine,
				markPoint:   series.HighMarkPoint,
				trendLines:  series.HighTrendLine,
				extractFunc: (*CandlestickSeries).ExtractHighPrices,
				points:      seriesHighPoints[seriesIndex],
			},
			{
				markLine:    series.LowMarkLine,
				markPoint:   series.LowMarkPoint,
				trendLines:  series.LowTrendLine,
				extractFunc: (*CandlestickSeries).ExtractLowPrices,
				points:      seriesLowPoints[seriesIndex],
			},
			{
				markLine:    series.CloseMarkLine,
				markPoint:   series.CloseMarkPoint,
				trendLines:  series.CloseTrendLine,
				extractFunc: (*CandlestickSeries).ExtractClosePrices,
				points:      seriesClosePoints[seriesIndex],
			},
		}

		for _, component := range ohlcComponents {
			var values []float64

			// Handle mark lines
			if seriesMarks := component.markLine.Lines.filterGlobal(false); len(seriesMarks) > 0 {
				if values == nil {
					values = component.extractFunc(series)
				}
				markLineValueFormatter := getPreferredValueFormatter(component.markLine.ValueFormatter,
					series.Label.ValueFormatter, opt.ValueFormatter)
				markLinePainter.add(markLineRenderOption{
					fillColor:      seriesColor,
					fontColor:      opt.Theme.GetMarkTextColor(),
					strokeColor:    seriesColor,
					font:           getPreferredFont(series.Label.FontStyle.Font),
					marklines:      seriesMarks,
					seriesValues:   values,
					axisRange:      yRange,
					valueFormatter: markLineValueFormatter,
				})
			}

			// Handle mark points
			if seriesMarks := component.markPoint.Points.filterGlobal(false); len(seriesMarks) > 0 {
				if values == nil {
					values = component.extractFunc(series)
				}
				markPointValueFormatter := getPreferredValueFormatter(component.markPoint.ValueFormatter,
					series.Label.ValueFormatter, opt.ValueFormatter)
				markPointPainter.add(markPointRenderOption{
					fillColor:          seriesColor,
					font:               getPreferredFont(series.Label.FontStyle.Font),
					symbolSize:         component.markPoint.SymbolSize,
					points:             component.points,
					markpoints:         seriesMarks,
					seriesValues:       values,
					valueFormatter:     markPointValueFormatter,
					seriesLabelPainter: allLabelPainters[seriesIndex],
				})
			}

			// Handle trend lines
			if len(component.trendLines) > 0 {
				if values == nil {
					values = component.extractFunc(series)
				}
				trendLinePainter.add(trendLineRenderOption{
					defaultStrokeColor: opt.Theme.GetSeriesTrendColor(seriesThemeIndex),
					xValues:            seriesCenterValues[seriesIndex],
					seriesValues:       values,
					axisRange:          yRange,
					trends:             component.trendLines,
					dashed:             false, // Default for candlestick charts
				})
			}
		}
	}

	if err := doRender(rendererList...); err != nil {
		return BoxZero, err
	}
	return p.box, nil
}

func (k *candlestickChart) Render() (Box, error) {
	p := k.p
	opt := k.opt
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.theme)
	}
	opt.Legend.Symbol = symbolCandlestick

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:          opt.Theme,
		padding:        opt.Padding,
		seriesList:     &opt.SeriesList,
		xAxis:          &opt.XAxis,
		yAxis:          opt.YAxis,
		title:          opt.Title,
		legend:         &opt.Legend,
		valueFormatter: opt.ValueFormatter,
	})
	if err != nil {
		return BoxZero, err
	}
	return k.renderChart(renderResult)
}
