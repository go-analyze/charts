package charts

import (
	"errors"
	"math"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

type barChart struct {
	p   *Painter
	opt *BarChartOption
}

// NewBarChart returns a bar chart renderer
func NewBarChart(p *Painter, opt BarChartOption) *barChart {
	return &barChart{
		p:   p,
		opt: &opt,
	}
}

type BarChartOption struct {
	// Theme specifies the colors used for the bar chart.
	Theme ColorPalette
	// Padding specifies the padding of bar chart.
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
	// BarWidth specifies the width of each bar.
	BarWidth int
	// RoundedBarCaps set to `true` to produce a bar graph where the bars have rounded tops.
	RoundedBarCaps *bool
}

func (b *barChart) render(result *defaultRenderResult, seriesList SeriesList) (Box, error) {
	p := b.p
	opt := b.opt
	seriesPainter := result.seriesPainter

	xRange := newRange(b.p, seriesPainter.Width(), len(opt.XAxis.Data), 0.0, 0.0, 0.0, 0.0)
	x0, x1 := xRange.GetRange(0)
	width := int(x1 - x0)
	// margin between each block
	margin := 10
	// margin between each bar
	barMargin := 5
	if width < 20 {
		margin = 2
		barMargin = 2
	} else if width < 50 {
		margin = 5
		barMargin = 3
	}
	seriesCount := len(seriesList)
	if seriesCount == 0 {
		return BoxZero, errors.New("empty series list")
	}
	barWidth := (width - 2*margin - barMargin*(seriesCount-1)) / seriesCount
	if opt.BarWidth > 0 && opt.BarWidth < barWidth {
		barWidth = opt.BarWidth
		// recalculate margin
		margin = (width - seriesCount*barWidth - barMargin*(seriesCount-1)) / 2
	}
	barMaxHeight := seriesPainter.Height()
	theme := opt.Theme
	seriesNames := seriesList.Names()

	markPointPainter := newMarkPointPainter(seriesPainter)
	markLinePainter := newMarkLinePainter(seriesPainter)
	rendererList := []Renderer{
		markPointPainter,
		markLinePainter,
	}
	for index := range seriesList {
		series := seriesList[index]
		yRange := result.axisRanges[series.YAxisIndex]
		seriesColor := theme.GetSeriesColor(series.index)

		divideValues := xRange.AutoDivide()
		points := make([]Point, len(series.Data))
		var labelPainter *seriesLabelPainter
		if flagIs(true, series.Label.Show) {
			labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, series.Label, opt.Theme, opt.Font)
			rendererList = append(rendererList, labelPainter)
		}

		for j, item := range series.Data {
			if j >= xRange.divideCount {
				continue
			}
			x := divideValues[j]
			x += margin
			if index != 0 {
				x += index * (barWidth + barMargin)
			}

			h := yRange.getHeight(item)
			fillColor := seriesColor
			top := barMaxHeight - h

			seriesPainter.OverrideDrawingStyle(chartdraw.Style{
				FillColor: fillColor,
			})
			if flagIs(true, opt.RoundedBarCaps) {
				seriesPainter.roundedRect(Box{
					Top:    top,
					Left:   x,
					Right:  x + barWidth,
					Bottom: barMaxHeight - 1,
					IsSet:  true,
				}, barWidth, true, false)
			} else {
				seriesPainter.filledRect(Box{
					Top:    top,
					Left:   x,
					Right:  x + barWidth,
					Bottom: barMaxHeight - 1,
					IsSet:  true,
				})
			}
			// generate marker point by hand
			points[j] = Point{
				X: x + (barWidth >> 1), // centered position
				Y: top,
			}
			// return if the label does not need to be displayed
			if labelPainter == nil {
				continue
			}
			y := barMaxHeight - h
			radians := float64(0)
			fontStyle := series.Label.FontStyle
			if series.Label.Position == PositionBottom {
				y = barMaxHeight
				radians = -math.Pi / 2
				if fontStyle.FontColor.IsZero() {
					if isLightColor(fillColor) {
						fontStyle.FontColor = defaultLightFontColor
					} else {
						fontStyle.FontColor = defaultDarkFontColor
					}
				}
			}
			labelPainter.Add(labelValue{
				vertical:  true, // label is above bar
				index:     index,
				value:     item,
				fontStyle: fontStyle,
				x:         x + (barWidth >> 1),
				y:         y,
				radians:   radians,
				offset:    series.Label.Offset,
			})
		}

		markPointPainter.Add(markPointRenderOption{
			FillColor: seriesColor,
			Font:      opt.Font,
			Series:    series,
			Points:    points,
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

func (b *barChart) Render() (Box, error) {
	p := b.p
	opt := b.opt
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.theme)
	}
	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:      opt.Theme,
		padding:    opt.Padding,
		seriesList: opt.SeriesList,
		xAxis:      opt.XAxis,
		yAxis:      opt.YAxis,
		title:      opt.Title,
		legend:     opt.Legend,
	})
	if err != nil {
		return BoxZero, err
	}
	seriesList := opt.SeriesList.Filter(ChartTypeBar)
	return b.render(renderResult, seriesList)
}
