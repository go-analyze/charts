package charts

import (
	"errors"
	"math"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

type lineChart struct {
	p   *Painter
	opt *LineChartOption
}

// newLineChart returns a line chart render
func newLineChart(p *Painter, opt LineChartOption) *lineChart {
	return &lineChart{
		p:   p,
		opt: &opt,
	}
}

// NewLineChartOptionWithData returns an initialized LineChartOption with the SeriesList set for the provided data slice.
func NewLineChartOptionWithData(data [][]float64) LineChartOption {
	sl := NewSeriesListLine(data)
	return LineChartOption{
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

type LineChartOption struct {
	// Theme specifies the colors used for the line chart.
	Theme ColorPalette
	// Padding specifies the padding of line chart.
	Padding Box
	// Font is the font used to render the chart.
	Font *truetype.Font
	// SeriesList provides the data population for the chart, typically constructed using NewSeriesListLine.
	SeriesList LineSeriesList
	// StackSeries if set to *true the lines will be layered over each other, with the last series value representing
	// the sum of all the values. Enabling this will also enable FillArea (which until v0.5 can't be disabled).
	// Some options will be ignored when StackedSeries is enabled, this includes StrokeSmoothingTension.
	// MarkLine is also interpreted differently, only the first Series will have the MarkLine rendered (as it's the
	// base bar, other bars are influenced by prior values). Additionally only the 0 index y-axis is stacked,
	// allowing a non-stacked line to also be included on y-axis 1.
	StackSeries *bool
	// XAxis are options for the x-axis.
	XAxis XAxisOption
	// YAxis are options for the y-axis (at most two).
	YAxis []YAxisOption
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// Symbol specifies the symbols to draw at the data points. Empty (default) will vary based on the dataset.
	// Specify 'none' to enforce no symbol, or specify a desired symbol: 'circle', 'dot', 'square', 'diamond'. This can
	// also be set on each series specifically.
	Symbol Symbol
	// LineStrokeWidth is the width of the rendered line.
	LineStrokeWidth float64
	// StrokeSmoothingTension should be between 0 and 1. At 0 perfectly straight lines will be used with 1 providing
	// smoother lines. Because the tension smooths out the line, the line will no longer hit the data points exactly.
	// The more variable the points, and the higher the tension, the more the line will be moved from the points.
	StrokeSmoothingTension float64
	// FillArea set this to true to fill the area below the line.
	FillArea *bool
	// FillOpacity is the opacity (alpha) of the area fill.
	FillOpacity uint8
	// ValueFormatter defines how float values should be rendered to strings, notably for numeric axis labels.
	ValueFormatter ValueFormatter
}

const showSymbolDefaultThreshold = 100

func boundaryGapAxisPositions(painterWidth int, boundaryGap bool, xDivideCount int) []int {
	if !boundaryGap {
		xDivideCount--
	}
	xDivideValues := autoDivide(painterWidth, xDivideCount)
	var xValues []int
	if boundaryGap {
		xValues = make([]int, len(xDivideValues)-1)
		for i := 0; i < len(xDivideValues)-1; i++ {
			xValues[i] = (xDivideValues[i] + xDivideValues[i+1]) >> 1
		}
	} else {
		xValues = xDivideValues
	}
	return xValues
}

func (l *lineChart) render(result *defaultRenderResult, seriesList LineSeriesList) (Box, error) {
	p := l.p
	opt := l.opt
	seriesCount := len(seriesList)
	if seriesCount == 0 {
		return BoxZero, errors.New("empty series list")
	}
	seriesPainter := result.seriesPainter

	stackedSeries := flagIs(true, opt.StackSeries)
	fillAreaY0 := stackedSeries // fill area defaults to on if the series is stacked
	fillAreaY1 := false
	if opt.FillArea != nil {
		fillAreaY0 = *opt.FillArea
		fillAreaY1 = *opt.FillArea
	}
	boundaryGap := !fillAreaY0 // boundary gap default enabled unless fill area is set
	if opt.XAxis.BoundaryGap != nil {
		boundaryGap = *opt.XAxis.BoundaryGap
	}
	xDivideCount := chartdraw.MaxInt(getSeriesMaxDataCount(opt.SeriesList), len(opt.XAxis.Labels))
	if boundaryGap && xDivideCount > 1 && seriesPainter.Width()/xDivideCount <= boundaryGapDefaultThreshold {
		// boundary gap would be so small it's visually better to disable the line spacing adjustment.
		// Although label changes can be forced to center, this behavior is unconditional for the line
		boundaryGap = false
	}
	xValues := boundaryGapAxisPositions(seriesPainter.Width(), boundaryGap, xDivideCount)
	dataCount := getSeriesMaxDataCount(seriesList)
	// accumulatedValues is used for stacking: it holds the summed data values at each X index
	var accumulatedValues []float64
	if stackedSeries {
		accumulatedValues = make([]float64, dataCount)
	}
	strokeWidth := opt.LineStrokeWidth
	if strokeWidth == 0 {
		strokeWidth = defaultStrokeWidth
	}
	symbol := opt.Symbol
	if symbol == "" {
		showSymbol := dataCount < showSymbolDefaultThreshold // default enable when data count is reasonable
		if opt.StrokeSmoothingTension > 0 {
			showSymbol = false // default disable symbols on curved lines since the dots won't hit the line exactly
		}
		if showSymbol {
			symbol = SymbolCircle
		} else {
			symbol = SymbolNone
		}
	}

	// render list must start with the markPointPainter, as it can influence label painters (if enabled)
	markPointPainter := newMarkPointPainter(seriesPainter)
	markLinePainter := newMarkLinePainter(seriesPainter)
	rendererList := []renderer{markPointPainter, markLinePainter}

	seriesNames := seriesList.names()
	var priorSeriesPoints []Point
	for index, series := range seriesList {
		stackSeries := stackedSeries && series.YAxisIndex == 0
		seriesColor := opt.Theme.GetSeriesColor(index)
		yRange := result.axisRanges[series.YAxisIndex]
		points := make([]Point, len(series.Values))
		var labelPainter *seriesLabelPainter
		if flagIs(true, series.Label.Show) {
			labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, series.Label, opt.Theme, opt.Font)
			rendererList = append(rendererList, labelPainter)
		}

		for i, item := range series.Values {
			if item == GetNullValue() {
				points[i] = Point{X: xValues[i], Y: math.MaxInt32}
			} else if stackSeries {
				accumulatedValues[i] += item
				points[i] = Point{
					X: xValues[i],
					Y: yRange.getRestHeight(accumulatedValues[i]),
				}
			} else {
				points[i] = Point{
					X: xValues[i],
					Y: yRange.getRestHeight(item),
				}
			}

			if labelPainter != nil {
				labelPainter.Add(labelValue{
					index:     index,
					value:     item,
					x:         points[i].X,
					y:         points[i].Y,
					fontStyle: series.Label.FontStyle,
				})
			}
		}

		if (series.YAxisIndex == 0 && fillAreaY0) || fillAreaY1 {
			areaPoints := make([]Point, len(points))
			copy(areaPoints, points)
			for i, p := range areaPoints {
				if p.Y != math.MaxInt32 {
					if i > 0 {
						areaPoints = areaPoints[i:] // truncate null points from array
					}
					break
				}
			}
			bottomY := yRange.getRestHeight(yRange.min)
			if stackSeries && len(priorSeriesPoints) > 0 {
				// Fill between current line (areaPoints) and priorSeriesPoints
				for i := len(priorSeriesPoints) - 1; i >= 0; i-- {
					areaPoints = append(areaPoints, priorSeriesPoints[i])
				}
				// Close the shape by re-appending the first of point
				areaPoints = append(areaPoints, areaPoints[0])
			} else {
				// Not stacked or first stacked series: fill down to bottom and then back to the first point
				areaPoints = append(areaPoints,
					Point{
						X: areaPoints[len(areaPoints)-1].X,
						Y: bottomY,
					}, Point{
						X: areaPoints[0].X,
						Y: bottomY,
					},
					areaPoints[0],
				)
			}

			var opacity uint8 = 200
			if opt.FillOpacity > 0 {
				opacity = opt.FillOpacity
			}
			fillColor := seriesColor.WithAlpha(opacity)

			// If smoothing is enabled, do a smooth fill (not currently supported for stacked series)
			if !stackSeries && opt.StrokeSmoothingTension > 0 {
				seriesPainter.smoothFillChartArea(areaPoints, opt.StrokeSmoothingTension, fillColor)
			} else {
				seriesPainter.FillArea(areaPoints, fillColor)
			}
		}

		// Draw the line
		if opt.StrokeSmoothingTension > 0 {
			seriesPainter.SmoothLineStroke(points, opt.StrokeSmoothingTension, seriesColor, strokeWidth)
		} else {
			seriesPainter.LineStroke(points, seriesColor, strokeWidth)
		}

		// Draw symbols if enabled
		seriesSymbol := symbol
		if series.Symbol != "" {
			seriesSymbol = series.Symbol
		}
		switch seriesSymbol {
		case SymbolCircle:
			seriesPainter.Dots(points, opt.Theme.GetBackgroundColor(), seriesColor, 1, strokeWidth)
		case SymbolDot:
			radius := 1.0
			if strokeWidth > 0 {
				radius = strokeWidth * 1.5
			}
			seriesPainter.Dots(points, seriesColor, seriesColor, 1, radius)
		case SymbolSquare:
			size := 2
			if strokeWidth > 1 {
				size = ceilFloatToInt(strokeWidth * 2.0)
			}
			seriesPainter.squares(points, seriesColor, seriesColor, 1, size)
		case SymbolDiamond:
			size := 4
			if strokeWidth > 1 {
				size = ceilFloatToInt(strokeWidth * 4.0)
			}
			seriesPainter.diamonds(points, seriesColor, seriesColor, 1, size)
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
				// In stacked mode we only support the line painter for the first series
				markLinePainter.add(markLineRenderOption{
					fillColor:      seriesColor,
					fontColor:      opt.Theme.GetMarkTextColor(),
					strokeColor:    seriesColor,
					font:           opt.Font,
					marklines:      seriesMarks,
					seriesValues:   series.Values,
					axisRange:      yRange,
					valueFormatter: markLineValueFormatter,
				})
			}
			if len(globalMarks) > 0 {
				if globalSeriesData == nil {
					globalSeriesData = sumSeriesData(seriesList, series.YAxisIndex)
				}
				markLinePainter.add(markLineRenderOption{
					fillColor:      defaultGlobalMarkFillColor,
					fontColor:      opt.Theme.GetMarkTextColor(),
					strokeColor:    defaultGlobalMarkFillColor,
					font:           opt.Font,
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
					font:               opt.Font,
					symbolSize:         series.MarkPoint.SymbolSize,
					points:             points,
					markpoints:         seriesMarks,
					seriesValues:       series.Values,
					valueFormatter:     markPointValueFormatter,
					seriesLabelPainter: labelPainter,
				})
			}
			if len(globalMarks) > 0 {
				if globalSeriesData == nil {
					globalSeriesData = sumSeriesData(seriesList, series.YAxisIndex)
				}
				markPointPainter.add(markPointRenderOption{
					fillColor:          defaultGlobalMarkFillColor,
					font:               opt.Font,
					symbolSize:         series.MarkPoint.SymbolSize,
					points:             points,
					markpoints:         globalMarks,
					seriesValues:       globalSeriesData,
					valueFormatter:     markPointValueFormatter,
					seriesLabelPainter: labelPainter,
				})
			}
		}

		// Save these points as "priorSeriesPoints" for the next series to stack onto (if needed)
		priorSeriesPoints = points
	}

	if err := doRender(rendererList...); err != nil {
		return BoxZero, err
	}
	return p.box, nil
}

func (l *lineChart) Render() (Box, error) {
	p := l.p
	opt := l.opt
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.theme)
	}
	// boundary gap default must be set here as it's used by the x-axis as well
	if opt.XAxis.BoundaryGap == nil {
		fillArea := flagIs(true, opt.StackSeries) // fill area default based on StackedSeries state
		if opt.FillArea != nil {                  // default override
			fillArea = *opt.FillArea
		}
		boundaryGap := !fillArea // boundary gap default enabled unless fill area is set
		l.opt.XAxis.BoundaryGap = &boundaryGap
	}

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:          opt.Theme,
		padding:        opt.Padding,
		seriesList:     opt.SeriesList,
		stackSeries:    flagIs(true, opt.StackSeries),
		xAxis:          &l.opt.XAxis,
		yAxis:          opt.YAxis,
		title:          opt.Title,
		legend:         &l.opt.Legend,
		valueFormatter: opt.ValueFormatter,
	})
	if err != nil {
		return BoxZero, err
	}
	return l.render(renderResult, opt.SeriesList)
}
