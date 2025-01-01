package charts

import (
	"errors"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

type horizontalBarChart struct {
	p   *Painter
	opt *HorizontalBarChartOption
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
	// XAxis are options for the x-axis.
	XAxis XAxisOption
	// YAxis are options for the y-axis (at most two).
	YAxis []YAxisOption
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// BarHeight specifies the height of each horizontal bar.
	BarHeight int
}

// NewHorizontalBarChart returns a horizontal bar chart renderer
func NewHorizontalBarChart(p *Painter, opt HorizontalBarChartOption) *horizontalBarChart {
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
	// margin between each block
	margin := 10
	// margin between each bar
	barMargin := 5
	if height < 20 {
		margin = 2
		barMargin = 2
	} else if height < 50 {
		margin = 5
		barMargin = 3
	}
	seriesCount := len(seriesList)
	if seriesCount == 0 {
		return BoxZero, errors.New("empty series list")
	}
	barHeight := (height - 2*margin - barMargin*(seriesCount-1)) / seriesCount
	if opt.BarHeight > 0 && opt.BarHeight < barHeight {
		barHeight = opt.BarHeight
		margin = (height - seriesCount*barHeight - barMargin*(seriesCount-1)) / 2
	}

	theme := opt.Theme

	min, max := seriesList.GetMinMax(0)
	xRange := newRange(p, seriesPainter.Width(), len(seriesList[0].Data), min, max, 1.0, 1.0)
	seriesNames := seriesList.Names()

	var rendererList []Renderer
	for index := range seriesList {
		series := seriesList[index]
		seriesColor := theme.GetSeriesColor(series.index)
		divideValues := yRange.AutoDivide()

		var labelPainter *seriesLabelPainter
		if series.Label.Show {
			labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, series.Label, opt.Theme, opt.Font)
			rendererList = append(rendererList, labelPainter)
		}
		for j, item := range series.Data {
			if j >= yRange.divideCount {
				continue
			}
			// display position switch
			j = yRange.divideCount - j - 1
			y := divideValues[j]
			y += margin
			if index != 0 {
				y += index * (barHeight + barMargin)
			}

			w := xRange.getHeight(item.Value)
			fillColor := seriesColor
			right := w
			seriesPainter.OverrideDrawingStyle(chartdraw.Style{
				FillColor: fillColor,
			}).Rect(chartdraw.Box{
				Top:    y,
				Left:   0,
				Right:  right,
				Bottom: y + barHeight,
				IsSet:  true,
			})
			// if the label does not need to be displayed, return
			if labelPainter == nil {
				continue
			}
			fontStyle := series.Label.FontStyle
			if fontStyle.FontColor.IsZero() {
				if isLightColor(fillColor) {
					fontStyle.FontColor = defaultLightFontColor
				} else {
					fontStyle.FontColor = defaultDarkFontColor
				}
			}
			labelValue := labelValue{
				vertical:  false, // label beside bar
				index:     index,
				value:     item.Value,
				x:         right,
				y:         y + (barHeight >> 1),
				offset:    series.Label.Offset,
				fontStyle: fontStyle,
			}
			if series.Label.Position == PositionLeft {
				labelValue.x = 0
			}
			labelPainter.Add(labelValue)
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
		Theme:        opt.Theme,
		Padding:      opt.Padding,
		SeriesList:   opt.SeriesList,
		XAxis:        opt.XAxis,
		YAxis:        opt.YAxis,
		Title:        opt.Title,
		Legend:       opt.Legend,
		axisReversed: true,
	})
	if err != nil {
		return BoxZero, err
	}
	seriesList := opt.SeriesList.Filter(ChartTypeHorizontalBar)
	return h.render(renderResult, seriesList)
}
