package charts

import (
	"errors"

	"github.com/golang/freetype/truetype"
)

type funnelChart struct {
	p   *Painter
	opt *FunnelChartOption
}

// newFunnelChart returns a funnel chart renderer
func newFunnelChart(p *Painter, opt FunnelChartOption) *funnelChart {
	return &funnelChart{
		p:   p,
		opt: &opt,
	}
}

// NewFunnelChartOptionWithData returns an initialized FunnelChartOption with the SeriesList set for the provided data slice.
func NewFunnelChartOptionWithData(data []float64) FunnelChartOption {
	return FunnelChartOption{
		SeriesList: NewSeriesListFunnel(data),
		Padding:    defaultPadding,
		Theme:      GetDefaultTheme(),
		Font:       GetDefaultFont(),
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
	max := seriesList[0].Data[0]
	min := float64(0)
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
		value := item.Data[0]
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
		if !flagIs(false, item.Label.Show) {
			if item.Label.ValueFormatter != nil {
				textList[index] = item.Label.ValueFormatter(value)
			} else {
				if item.Label.FormatTemplate == "" {
					item.Label.FormatTemplate = item.Label.Formatter
				}
				textList[index] = labelFormatFunnel(seriesNames, item.Label.FormatTemplate, index, value, percent)
			}
		}
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

		seriesPainter.FillArea(points, theme.GetSeriesColor(series.index))

		text := textList[index]
		fontStyle := FontStyle{
			FontColor: theme.GetTextColor(),
			FontSize:  labelFontSize,
			Font:      opt.Font,
		}
		textBox := seriesPainter.MeasureText(text, 0, fontStyle)
		textX := width>>1 - textBox.Width()>>1
		textY := y + h>>1
		seriesPainter.Text(text, textX, textY, 0, fontStyle)
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
		theme:      opt.Theme,
		padding:    opt.Padding,
		seriesList: opt.SeriesList,
		xAxis: &XAxisOption{
			Show: False(),
		},
		yAxis: []YAxisOption{
			{
				Show: False(),
			},
		},
		title:  opt.Title,
		legend: &f.opt.Legend,
	})
	if err != nil {
		return BoxZero, err
	}
	seriesList := opt.SeriesList.Filter(ChartTypeFunnel)
	return f.render(renderResult, seriesList)
}
