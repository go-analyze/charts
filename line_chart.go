package charts

import (
	"math"

	"github.com/golang/freetype/truetype"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type lineChart struct {
	p   *Painter
	opt *LineChartOption
}

// NewLineChart returns a line chart render
func NewLineChart(p *Painter, opt LineChartOption) *lineChart {
	if opt.Theme == nil {
		opt.Theme = defaultTheme
	}
	return &lineChart{
		p:   p,
		opt: &opt,
	}
}

type LineChartOption struct {
	// The theme
	Theme ColorPalette
	// The font size
	Font *truetype.Font
	// The data series list
	SeriesList SeriesList
	// The x axis option
	XAxis XAxisOption
	// The padding of line chart
	Padding Box
	// The y axis option
	YAxisOptions []YAxisOption
	// The option of title
	Title TitleOption
	// The legend option
	Legend LegendOption
	// The flag for show symbol of line, set this to *false will hide symbol
	SymbolShow *bool
	// The stroke width of line
	StrokeWidth float64
	// Fill the area of line
	FillArea bool
	// background is filled
	backgroundIsFilled bool
	// background fill (alpha) opacity
	Opacity uint8
}

func (l *lineChart) render(result *defaultRenderResult, seriesList SeriesList) (Box, error) {
	p := l.p
	opt := l.opt
	boundaryGap := true
	if isFalse(opt.XAxis.BoundaryGap) {
		boundaryGap = false
	}

	seriesPainter := result.seriesPainter

	xDivideCount := len(opt.XAxis.Data)
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
		drawingStyle := Style{
			StrokeColor: seriesColor,
			StrokeWidth: strokeWidth,
		}
		if len(series.Style.StrokeDashArray) > 0 {
			drawingStyle.StrokeDashArray = series.Style.StrokeDashArray
		}

		yRange := result.axisRanges[series.AxisIndex]
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
				h = int(math.MaxInt32)
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
				Index: index,
				Value: item.Value,
				X:     p.X,
				Y:     p.Y,
				FontSize: series.Label.FontSize,
			})
		}
		if opt.FillArea {
			areaPoints := make([]Point, len(points))
			copy(areaPoints, points)
			bottomY := yRange.getRestHeight(yRange.min)
			var opacity uint8 = 200
			if opt.Opacity != 0 {
				opacity = opt.Opacity
			}
			areaPoints = append(areaPoints, Point{
				X: areaPoints[len(areaPoints)-1].X,
				Y: bottomY,
			}, Point{
				X: areaPoints[0].X,
				Y: bottomY,
			}, areaPoints[0])
			seriesPainter.SetDrawingStyle(Style{
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
		if !isFalse(opt.SymbolShow) {
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
	err := doRender(rendererList...)
	if err != nil {
		return BoxZero, err
	}

	return p.box, nil
}

func (l *lineChart) Render() (Box, error) {
	p := l.p
	opt := l.opt

	renderResult, err := defaultRender(p, defaultRenderOption{
		Theme:              opt.Theme,
		Padding:            opt.Padding,
		SeriesList:         opt.SeriesList,
		XAxis:              opt.XAxis,
		YAxisOptions:       opt.YAxisOptions,
		TitleOption:        opt.Title,
		LegendOption:       opt.Legend,
		backgroundIsFilled: opt.backgroundIsFilled,
	})
	if err != nil {
		return BoxZero, err
	}
	seriesList := opt.SeriesList.Filter(ChartTypeLine)

	return l.render(renderResult, seriesList)
}
