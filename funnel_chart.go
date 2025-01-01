package charts

import (
	"errors"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

type funnelChart struct {
	p   *Painter
	opt *FunnelChartOption
}

// NewFunnelSeriesList returns a series list for funnel
func NewFunnelSeriesList(values []float64) SeriesList {
	seriesList := make(SeriesList, len(values))
	for index, value := range values {
		seriesList[index] = NewSeriesFromValues([]float64{
			value,
		}, ChartTypeFunnel)
	}
	return seriesList
}

// NewFunnelChart returns a funnel chart renderer
func NewFunnelChart(p *Painter, opt FunnelChartOption) *funnelChart {
	return &funnelChart{
		p:   p,
		opt: &opt,
	}
}

type FunnelChartOption struct {
	// Theme specifies the colors used for the chart.
	Theme ColorPalette
	// Padding specifies the padding of funnel chart.
	Padding Box
	// Font is the font used to render the chart.
	Font *truetype.Font
	// SeriesList provides the data series.
	SeriesList SeriesList
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
}

func (f *funnelChart) render(result *defaultRenderResult, seriesList SeriesList) (Box, error) {
	opt := f.opt
	seriesPainter := result.seriesPainter
	max := seriesList[0].Data[0].Value
	min := float64(0)
	for _, item := range seriesList {
		if item.Max != nil {
			max = *item.Max
		}
		if item.Min != nil {
			min = *item.Min
		}
	}
	theme := opt.Theme
	gap := 2
	height := seriesPainter.Height()
	width := seriesPainter.Width()
	count := len(seriesList)
	if count == 0 {
		return BoxZero, errors.New("empty series list")
	}

	h := (height - gap*(count-1)) / count

	y := 0
	widthList := make([]int, len(seriesList))
	textList := make([]string, len(seriesList))
	seriesNames := seriesList.Names()
	offset := max - min
	for index, item := range seriesList {
		value := item.Data[0].Value
		// if the maximum and minimum are consistent it's 100%
		widthPercent := 100.0
		if offset != 0 {
			widthPercent = (value - min) / offset
		}
		w := int(widthPercent * float64(width))
		widthList[index] = w
		// if the maximum value is 0, the proportion is 100%
		percent := 1.0
		if max != 0 {
			percent = value / max
		}
		textList[index] = labelFormatFunnel(seriesNames, item.Label.Formatter, index, value, percent)
	}

	for index, w := range widthList {
		series := seriesList[index]
		nextWidth := 0
		if index+1 < len(widthList) {
			nextWidth = widthList[index+1]
		}
		topStartX := (width - w) >> 1
		topEndX := topStartX + w
		bottomStartX := (width - nextWidth) >> 1
		bottomEndX := bottomStartX + nextWidth
		points := []Point{
			{
				X: topStartX,
				Y: y,
			},
			{
				X: topEndX,
				Y: y,
			},
			{
				X: bottomEndX,
				Y: y + h,
			},
			{
				X: bottomStartX,
				Y: y + h,
			},
			{
				X: topStartX,
				Y: y,
			},
		}
		color := theme.GetSeriesColor(series.index)

		seriesPainter.OverrideDrawingStyle(chartdraw.Style{
			FillColor: color,
		}).FillArea(points)

		// text
		text := textList[index]
		seriesPainter.OverrideFontStyle(FontStyle{
			FontColor: theme.GetTextColor(),
			FontSize:  labelFontSize,
			Font:      opt.Font,
		})
		textBox := seriesPainter.MeasureText(text)
		textX := width>>1 - textBox.Width()>>1
		textY := y + h>>1
		seriesPainter.Text(text, textX, textY)
		y += h + gap
	}

	return f.p.box, nil
}

func (f *funnelChart) Render() (Box, error) {
	p := f.p
	opt := f.opt
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.theme)
	}

	renderResult, err := defaultRender(p, defaultRenderOption{
		Theme:      opt.Theme,
		Padding:    opt.Padding,
		SeriesList: opt.SeriesList,
		XAxis: XAxisOption{
			Show: False(),
		},
		YAxis: []YAxisOption{
			{
				Show: False(),
			},
		},
		Title:  opt.Title,
		Legend: opt.Legend,
	})
	if err != nil {
		return BoxZero, err
	}
	seriesList := opt.SeriesList.Filter(ChartTypeFunnel)
	return f.render(renderResult, seriesList)
}
