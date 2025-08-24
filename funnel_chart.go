package charts

import (
	"errors"

	"github.com/dustin/go-humanize"
	"github.com/golang/freetype/truetype"
)

type funnelChart struct {
	p   *Painter
	opt *FunnelChartOption
}

// newFunnelChart returns a funnel chart renderer.
func newFunnelChart(p *Painter, opt FunnelChartOption) *funnelChart {
	return &funnelChart{
		p:   p,
		opt: &opt,
	}
}

// NewFunnelChartOptionWithData returns an initialized FunnelChartOption with the SeriesList set with the provided data slice.
func NewFunnelChartOptionWithData(data []float64) FunnelChartOption {
	return FunnelChartOption{
		SeriesList: NewSeriesListFunnel(data),
		Padding:    defaultPadding,
		Theme:      GetDefaultTheme(),
		Font:       GetDefaultFont(),
	}
}

// FunnelChartOption defines the options for rendering a funnel chart. Render the chart using Painter.FunnelChart.
type FunnelChartOption struct {
	// Theme specifies the colors used for the chart.
	Theme ColorPalette
	// Padding specifies the padding around the chart.
	Padding Box
	// Deprecated: Font is deprecated, instead the font needs to be set on the SeriesLabel, or other specific elements.
	Font *truetype.Font
	// SeriesList provides the data population for the chart. Typically constructed using NewSeriesListFunnel.
	SeriesList FunnelSeriesList
	// Title contains options for rendering the chart title.
	Title TitleOption
	// Legend contains options for the data legend.
	Legend LegendOption
	// Deprecated: ValueFormatter is deprecated, instead set the ValueFormatter at `SeriesList[*].Label.ValueFormatter`.
	ValueFormatter ValueFormatter
}

func (f *funnelChart) renderChart(result *defaultRenderResult) (Box, error) {
	opt := f.opt
	seriesCount := len(opt.SeriesList)
	if seriesCount == 0 {
		return BoxZero, errors.New("empty series list")
	}
	seriesPainter := result.seriesPainter
	max := opt.SeriesList[0].Value
	var min float64
	theme := opt.Theme
	gap := 2
	height := seriesPainter.Height()
	width := seriesPainter.Width()

	h := (height - gap*(seriesCount-1)) / seriesCount

	var y int
	widthList := make([]int, seriesCount)
	textList := make([]string, seriesCount)
	labelStyleList := make([]*LabelStyle, seriesCount)
	seriesNames := opt.SeriesList.names()
	offset := max - min
	for index, item := range opt.SeriesList {
		// if the maximum and minimum are consistent it's 100%
		widthPercent := 100.0
		if offset != 0 {
			widthPercent = (item.Value - min) / offset
		}
		w := int(widthPercent * float64(width))
		widthList[index] = w
		// if the maximum value is 0, the proportion is 100%
		percent := 1.0
		if max != 0 {
			percent = item.Value / max
		}
		if !flagIs(false, item.Label.Show) {
			if item.Label.LabelFormatter != nil {
				textList[index], labelStyleList[index] = item.Label.LabelFormatter(index, seriesNames[index], item.Value)
			} else if item.Label.FormatTemplate != "" {
				valueFormatter := item.Label.ValueFormatter
				if valueFormatter == nil {
					valueFormatter = opt.ValueFormatter
				}
				textList[index] = labelFormatFunnel(seriesNames[index], item.Label.FormatTemplate, valueFormatter,
					item.Value, percent)
			} else if item.Label.ValueFormatter != nil {
				textList[index] = item.Label.ValueFormatter(item.Value)
			} else { // default label
				textList[index] = seriesNames[index] + "(" + humanize.FtoaWithDigits(percent*100, 2) + "%)"
			}
		}
	}

	for index, w := range widthList {
		var nextWidth int
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

		seriesPainter.FillArea(points, theme.GetSeriesColor(index))

		text := textList[index]
		fontStyle := fillFontStyleDefaults(opt.SeriesList[index].Label.FontStyle,
			defaultLabelFontSize, theme.GetLabelTextColor(), opt.Font, seriesPainter.font)

		// Apply label style overrides if present
		var backgroundColor Color
		var cornerRadius int
		var borderColor Color
		var borderWidth float64
		if labelStyleList[index] != nil {
			fontStyle = mergeFontStyles(labelStyleList[index].FontStyle, fontStyle)
			backgroundColor = labelStyleList[index].BackgroundColor
			cornerRadius = labelStyleList[index].CornerRadius
			borderColor = labelStyleList[index].BorderColor
			borderWidth = labelStyleList[index].BorderWidth
		}

		textBox := seriesPainter.MeasureText(text, 0, fontStyle)
		textX := width>>1 - textBox.Width()>>1
		textY := y + h>>1
		drawLabelWithBackground(seriesPainter, text, textX, textY, 0, fontStyle, backgroundColor, cornerRadius, borderColor, borderWidth)
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
	if opt.Legend.Symbol == "" {
		opt.Legend.Symbol = SymbolSquare
	}

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:      opt.Theme,
		padding:    opt.Padding,
		seriesList: opt.SeriesList,
		xAxis: &XAxisOption{
			Show: Ptr(false),
		},
		yAxis: []YAxisOption{
			{
				Show: Ptr(false),
			},
		},
		title:  opt.Title,
		legend: &f.opt.Legend,
	})
	if err != nil {
		return BoxZero, err
	}
	return f.renderChart(renderResult)
}
