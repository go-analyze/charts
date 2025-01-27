package charts

import (
	"math"

	"github.com/golang/freetype/truetype"
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
			Data: make([]string, len(data[0])),
		},
		YAxis:          make([]YAxisOption, sl.getYAxisCount()),
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
	// SeriesList provides the data series.
	SeriesList SeriesList
	// StackSeries if set to *true the lines will be layered over each other, with the last series value representing
	// the sum of all the values. Enabling this will also enable FillArea (which until v0.5 can't be disabled).
	// Some options will be ignored when StackedSeries is enabled, this includes StrokeSmoothingTension.
	// MarkLine is also interpreted differently, only the first Series will have the MarkLine rendered (as it's the
	// base bar, other bars are influenced by prior values).
	StackSeries *bool
	// XAxis are options for the x-axis.
	XAxis XAxisOption
	// YAxis are options for the y-axis (at most two).
	YAxis []YAxisOption
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// SymbolShow set this to *false or *true (using False() or True()) to force if the symbols should be shown or hidden.
	SymbolShow *bool
	// LineStrokeWidth is the width of the rendered line.
	LineStrokeWidth float64
	// StrokeSmoothingTension should be between 0 and 1. At 0 perfectly straight lines will be used with 1 providing
	// smoother lines. Because the tension smooths out the line, the line will no longer hit the data points exactly.
	// The more variable the points, and the higher the tension, the more the line will be moved from the points.
	StrokeSmoothingTension float64
	// TODO - make FillArea a pointer so that it can be disabled for stacking, update StackSeries docs when done
	// FillArea set this to true to fill the area below the line.
	FillArea bool
	// FillOpacity is the opacity (alpha) of the area fill.
	FillOpacity uint8
	// ValueFormatter defines how float values should be rendered to strings, notably for numeric axis labels.
	ValueFormatter ValueFormatter
	// backgroundIsFilled is set to true if the background is filled.
	backgroundIsFilled bool
}

const showSymbolDefaultThreshold = 100

func (l *lineChart) render(result *defaultRenderResult, seriesList SeriesList) (Box, error) {
	p := l.p
	opt := l.opt
	seriesPainter := result.seriesPainter

	boundaryGap := !opt.FillArea // boundary gap default enabled unless fill area is set
	if opt.XAxis.BoundaryGap != nil {
		boundaryGap = *opt.XAxis.BoundaryGap
	}
	xDivideCount := len(opt.XAxis.Data)
	if boundaryGap && xDivideCount > 1 && seriesPainter.Width()/xDivideCount <= boundaryGapDefaultThreshold {
		// boundary gap would be so small it's visually better to disable the line spacing adjustment.
		// Although label changes can be forced to center, this behavior is unconditional for the line
		boundaryGap = false
	}
	if !boundaryGap {
		xDivideCount--
	}
	xDivideValues := autoDivide(seriesPainter.Width(), xDivideCount)
	xValues := make([]int, len(xDivideValues)-1)
	// accumulatedValues is used for stacking: it holds the summed data values at each X index
	var accumulatedValues []float64
	if boundaryGap {
		for i := 0; i < len(xDivideValues)-1; i++ {
			xValues[i] = (xDivideValues[i] + xDivideValues[i+1]) >> 1
		}
	} else {
		xValues = xDivideValues
	}

	// render list must start with the markPointPainter, as it can influence label painters (if enabled)
	markPointPainter := newMarkPointPainter(seriesPainter)
	markLinePainter := newMarkLinePainter(seriesPainter)
	rendererList := []renderer{
		markPointPainter,
		markLinePainter,
	}

	strokeWidth := opt.LineStrokeWidth
	if strokeWidth == 0 {
		strokeWidth = defaultStrokeWidth
	}

	var dataCount int
	for _, s := range seriesList {
		l := len(s.Data)
		if l > dataCount {
			dataCount = l
		}
	}

	stackedSeries := flagIs(true, opt.StackSeries)
	if stackedSeries {
		accumulatedValues = make([]float64, dataCount)
	}

	showSymbol := dataCount < showSymbolDefaultThreshold // default enable when data count is reasonable
	if opt.StrokeSmoothingTension > 0 {
		showSymbol = false // default disable symbols on curved lines since the dots won't hit the line exactly
	}
	if opt.SymbolShow != nil {
		showSymbol = *opt.SymbolShow
	}

	seriesNames := seriesList.Names()
	var priorSeriesPoints []Point
	for index := range seriesList {
		series := seriesList[index]
		seriesColor := opt.Theme.GetSeriesColor(series.index)
		yRange := result.axisRanges[series.YAxisIndex]
		points := make([]Point, len(series.Data))
		var labelPainter *seriesLabelPainter
		if flagIs(true, series.Label.Show) {
			labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, series.Label, opt.Theme, opt.Font)
			rendererList = append(rendererList, labelPainter)
		}

		for i, item := range series.Data {
			if item == GetNullValue() {
				points[i] = Point{X: xValues[i], Y: math.MaxInt32}
			} else if stackedSeries {
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

		if stackedSeries || opt.FillArea {
			areaPoints := make([]Point, len(points))
			copy(areaPoints, points)
			bottomY := yRange.getRestHeight(yRange.min)
			if stackedSeries && len(priorSeriesPoints) > 0 {
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
			if !stackedSeries && opt.StrokeSmoothingTension > 0 {
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

		// Draw dots if enabled
		if showSymbol {
			dotFillColor := ColorWhite
			if opt.Theme.IsDark() {
				dotFillColor = seriesColor
			}
			seriesPainter.Dots(points, dotFillColor, seriesColor, 1, 2)
		}

		// In stacked mode we only support the line painter for the first series
		if !stackedSeries || index == 0 {
			markLinePainter.Add(markLineRenderOption{
				fillColor:      seriesColor,
				fontColor:      opt.Theme.GetTextColor(),
				strokeColor:    seriesColor,
				font:           opt.Font,
				series:         series,
				axisRange:      yRange,
				valueFormatter: opt.ValueFormatter,
			})
		}
		markPointPainter.Add(markPointRenderOption{
			fillColor:          seriesColor,
			font:               opt.Font,
			points:             points,
			series:             series,
			valueFormatter:     opt.ValueFormatter,
			seriesLabelPainter: labelPainter,
		})

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
		boundaryGap := !opt.FillArea // boundary gap default enabled unless fill area is set
		l.opt.XAxis.BoundaryGap = &boundaryGap
	}

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:              opt.Theme,
		padding:            opt.Padding,
		seriesList:         opt.SeriesList,
		stackSeries:        flagIs(true, opt.StackSeries),
		xAxis:              &l.opt.XAxis,
		yAxis:              opt.YAxis,
		title:              opt.Title,
		legend:             &l.opt.Legend,
		valueFormatter:     opt.ValueFormatter,
		backgroundIsFilled: opt.backgroundIsFilled,
	})
	if err != nil {
		return BoxZero, err
	}
	seriesList := opt.SeriesList.Filter(ChartTypeLine)

	return l.render(renderResult, seriesList)
}
