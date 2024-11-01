package charts

import (
	"math"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

type lineChart struct {
	p   *Painter
	opt *LineChartOption
}

// NewLineChart returns a line chart render
func NewLineChart(p *Painter, opt LineChartOption) *lineChart {
	return &lineChart{
		p:   p,
		opt: &opt,
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
	// XAxis are options for the x-axis.
	XAxis XAxisOption
	// YAxis are options for the y-axis (at most two).
	YAxis []YAxisOption
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// SymbolShow set this to *false (through False()) to hide symbols.
	SymbolShow *bool
	// StrokeWidth is the width of the rendered line.
	StrokeWidth float64
	// FillArea set this to true to fill the area below the line.
	FillArea bool
	// FillOpacity is the opacity (alpha) of the area fill.
	FillOpacity uint8
	// backgroundIsFilled is set to true if the background is filled.
	backgroundIsFilled bool
}

func (l *lineChart) render(result *defaultRenderResult, seriesList SeriesList) (Box, error) {
	p := l.p
	opt := l.opt
	seriesPainter := result.seriesPainter

	boundaryGap := !flagIs(false, opt.XAxis.BoundaryGap)
	xDivideCount := len(opt.XAxis.Data)
	if boundaryGap && xDivideCount > 1 && seriesPainter.Width()/xDivideCount <= 10 {
		// boundary gap would be so small it's visually better to disable the line spacing adjustment and just keep
		// the label changes only
		boundaryGap = false
	}
	if !boundaryGap {
		xDivideCount--
	}
	xDivideValues := autoDivide(seriesPainter.Width(), xDivideCount)
	xValues := make([]int, len(xDivideValues)-1)
	if boundaryGap {
		for i := 0; i < len(xDivideValues)-1; i++ {
			xValues[i] = (xDivideValues[i] + xDivideValues[i+1]) >> 1
		}
	} else {
		xValues = xDivideValues
	}
	markPointPainter := NewMarkPointPainter(seriesPainter)
	markLinePainter := NewMarkLinePainter(seriesPainter)
	rendererList := []Renderer{
		markPointPainter,
		markLinePainter,
	}
	strokeWidth := opt.StrokeWidth
	if strokeWidth == 0 {
		strokeWidth = defaultStrokeWidth
	}
	seriesNames := seriesList.Names()
	for index := range seriesList {
		series := seriesList[index]
		seriesColor := opt.Theme.GetSeriesColor(series.index)
		drawingStyle := chartdraw.Style{
			StrokeColor: seriesColor,
			StrokeWidth: strokeWidth,
		}

		yRange := result.axisRanges[series.YAxisIndex]
		points := make([]Point, 0)
		var labelPainter *SeriesLabelPainter
		if series.Label.Show {
			labelPainter = NewSeriesLabelPainter(SeriesLabelPainterParams{
				P:           seriesPainter,
				SeriesNames: seriesNames,
				Label:       series.Label,
				Theme:       opt.Theme,
				Font:        opt.Font,
			})
			rendererList = append(rendererList, labelPainter)
		}
		for i, item := range series.Data {
			h := yRange.getRestHeight(item.Value)
			if item.Value == nullValue {
				h = math.MaxInt32
			}
			p := Point{
				X: xValues[i],
				Y: h,
			}
			points = append(points, p)

			// if the label does not need to be displayed, return
			if labelPainter == nil {
				continue
			}
			labelPainter.Add(LabelValue{
				Index:     index,
				Value:     item.Value,
				X:         p.X,
				Y:         p.Y,
				FontStyle: series.Label.FontStyle,
			})
		}
		if opt.FillArea {
			areaPoints := make([]Point, len(points))
			copy(areaPoints, points)
			bottomY := yRange.getRestHeight(yRange.min)
			var opacity uint8 = 200
			if opt.FillOpacity > 0 {
				opacity = opt.FillOpacity
			}
			areaPoints = append(areaPoints, Point{
				X: areaPoints[len(areaPoints)-1].X,
				Y: bottomY,
			}, Point{
				X: areaPoints[0].X,
				Y: bottomY,
			}, areaPoints[0])
			seriesPainter.SetDrawingStyle(chartdraw.Style{
				FillColor: seriesColor.WithAlpha(opacity),
			})
			seriesPainter.FillArea(areaPoints)
		}
		seriesPainter.SetDrawingStyle(drawingStyle)

		// draw line
		seriesPainter.LineStroke(points)

		// draw dots
		if opt.Theme.IsDark() {
			drawingStyle.FillColor = drawingStyle.StrokeColor
		} else {
			drawingStyle.FillColor = drawing.ColorWhite
		}
		drawingStyle.StrokeWidth = 1
		seriesPainter.SetDrawingStyle(drawingStyle)
		if !flagIs(false, opt.SymbolShow) {
			seriesPainter.Dots(points)
		}
		markPointPainter.Add(markPointRenderOption{
			FillColor: seriesColor,
			Font:      opt.Font,
			Points:    points,
			Series:    series,
		})
		markLinePainter.Add(markLineRenderOption{
			FillColor:   seriesColor,
			FontColor:   opt.Theme.GetTextColor(),
			StrokeColor: seriesColor,
			Font:        opt.Font,
			Series:      series,
			Range:       yRange,
		})
	}
	// the largest and smallest mark point
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

	renderResult, err := defaultRender(p, defaultRenderOption{
		Theme:              opt.Theme,
		Padding:            opt.Padding,
		SeriesList:         opt.SeriesList,
		XAxis:              opt.XAxis,
		YAxis:              opt.YAxis,
		Title:              opt.Title,
		Legend:             opt.Legend,
		backgroundIsFilled: opt.backgroundIsFilled,
	})
	if err != nil {
		return BoxZero, err
	}
	seriesList := opt.SeriesList.Filter(ChartTypeLine)

	return l.render(renderResult, seriesList)
}
