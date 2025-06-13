package charts

import (
	"errors"

	"github.com/golang/freetype/truetype"
)

// ChartOption represents a generic method of representing a chart. This can be useful when you want to render
// different chart types with the same data and configuration.
type ChartOption struct {
	// OutputFormat specifies the output type of chart, "svg", "png" or "jpg", default value is "png".
	OutputFormat string
	// Width is the width of chart.
	Width int
	// Height is the height of chart.
	Height int
	// Theme specifies the colors used for the chart. Built in themes can be loaded using GetTheme with
	// "light", "dark", "vivid-light", "vivid-dark", "ant" or "grafana".
	Theme ColorPalette
	// Padding specifies the padding for chart, default padding is [20, 20, 20, 20].
	Padding Box
	// XAxis are options for the x-axis.
	XAxis XAxisOption
	// YAxis are options for the y-axis (at most two).
	YAxis []YAxisOption
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// Deprecated: Font is deprecated, instead the font needs to be set on the SeriesLabel, or other specific elements.
	Font *truetype.Font
	// Box specifies the canvas box for the chart.
	Box Box
	// SeriesList provides the population data for the charts, constructed through NewSeriesListGeneric.
	SeriesList GenericSeriesList
	// StackSeries if set to *true the lines will be layered or stacked. This option significantly changes the chart
	// visualization, please see the specific chart docs for full details.
	StackSeries *bool
	// RadarIndicators are radar indicator list for radar charts.
	RadarIndicators []RadarIndicator
	// Symbol specifies the symbols to draw at the data points. Empty (default) will vary based on the dataset.
	// Specify 'none' to enforce no symbol, or specify a desired symbol: 'circle', 'dot', 'square', 'diamond'.
	Symbol Symbol
	// LineStrokeWidth is the stroke width for line charts.
	LineStrokeWidth float64
	// FillArea set to *true to fill the area under the line in line charts
	FillArea *bool
	// FillOpacity is the opacity (alpha) of the area fill in line charts.
	FillOpacity uint8
	// Deprecated: BarWidth is deprecated, instead use BarSize.
	BarWidth int
	// Deprecated: BarHeight is deprecated, instead use BarSize.
	BarHeight int
	// BarSize represents the width of the bar, or if a horizontal bar chart the height.
	BarSize int
	// BarMargin specifies the margin between bars grouped together. BarWidth or BarHeight takes priority over the margin.
	BarMargin *float64
	// Radius default radius for pie and radar charts e.g.: 40%, default is "40%"
	Radius string
	// Children are Child charts to render together.
	Children []ChartOption
	parent   *Painter
	// ValueFormatter to format numeric values into labels.
	ValueFormatter ValueFormatter
}

// OptionFunc option function.
type OptionFunc func(opt *ChartOption)

// SVGOutputOptionFunc sets svg as the image type for the chart's output.
func SVGOutputOptionFunc() OptionFunc {
	return outputFormatOptionFunc(ChartOutputSVG)
}

// PNGOutputOptionFunc sets png as the image type for the chart's output.
func PNGOutputOptionFunc() OptionFunc {
	return outputFormatOptionFunc(ChartOutputPNG)
}

// JPGOutputOptionFunc sets jpg as the image type for the chart's output.
func JPGOutputOptionFunc() OptionFunc {
	return outputFormatOptionFunc(ChartOutputJPG)
}

// outputFormatOptionFunc set type of chart's output.
func outputFormatOptionFunc(t string) OptionFunc {
	return func(opt *ChartOption) {
		opt.OutputFormat = t
	}
}

// Deprecated: FontOptionFunc is deprecated, fonts should be set on the specific elements (SeriesLabel, Title, etc).
func FontOptionFunc(font *truetype.Font) OptionFunc {
	return func(opt *ChartOption) {
		opt.Font = font
	}
}

// ThemeNameOptionFunc sets the theme of the chart by name.
func ThemeNameOptionFunc(theme string) OptionFunc {
	return func(opt *ChartOption) {
		opt.Theme = GetTheme(theme)
	}
}

// ThemeOptionFunc sets the theme of the chart.
func ThemeOptionFunc(theme ColorPalette) OptionFunc {
	return func(opt *ChartOption) {
		opt.Theme = theme
	}
}

// TitleOptionFunc sets the title of the chart.
func TitleOptionFunc(title TitleOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.Title = title
	}
}

// TitleTextOptionFunc sets the title text of chart.
func TitleTextOptionFunc(text string, subtext ...string) OptionFunc {
	return func(opt *ChartOption) {
		opt.Title.Text = text
		if len(subtext) != 0 {
			opt.Title.Subtext = subtext[0]
		}
	}
}

// LegendOptionFunc sets the legend of the chart.
func LegendOptionFunc(legend LegendOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.Legend = legend
	}
}

// LegendLabelsOptionFunc sets the legend series name labels of the chart.
func LegendLabelsOptionFunc(labels []string) OptionFunc {
	return func(opt *ChartOption) {
		opt.Legend = LegendOption{
			SeriesNames: labels,
		}
	}
}

// XAxisOptionFunc sets the x-axis of the chart.
func XAxisOptionFunc(xAxisOption XAxisOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.XAxis = xAxisOption
	}
}

// XAxisLabelsOptionFunc sets the x-axis labels of the chart.
func XAxisLabelsOptionFunc(labels []string) OptionFunc {
	return func(opt *ChartOption) {
		opt.XAxis = XAxisOption{
			Labels: labels,
		}
	}
}

// YAxisOptionFunc sets the y-axis of chart, supports up to two y-axis.
func YAxisOptionFunc(yAxisOption ...YAxisOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.YAxis = yAxisOption
	}
}

// YAxisLabelsOptionFunc sets the y-axis labels of the chart.
func YAxisLabelsOptionFunc(labels []string) OptionFunc {
	return func(opt *ChartOption) {
		opt.YAxis = []YAxisOption{
			{
				Labels: labels,
			},
		}
	}
}

// DimensionsOptionFunc sets the width and height dimensions of the chart.
func DimensionsOptionFunc(width, height int) OptionFunc {
	return func(opt *ChartOption) {
		opt.Width = width
		opt.Height = height
	}
}

// PaddingOptionFunc sets the padding of the chart.
func PaddingOptionFunc(padding Box) OptionFunc {
	return func(opt *ChartOption) {
		opt.Padding = padding
	}
}

// SeriesShowLabel sets the series show label state for all series.
func SeriesShowLabel(show bool) OptionFunc {
	return func(opt *ChartOption) {
		for index := range opt.SeriesList {
			opt.SeriesList[index].Label.Show = Ptr(show)
		}
	}
}

// ChildOptionFunc adds a Child chart on top of the current one. Use Padding and Box for positioning.
func ChildOptionFunc(child ...ChartOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.Children = append(opt.Children, child...)
	}
}

// RadarIndicatorOptionFunc sets the radar indicator of chart
func RadarIndicatorOptionFunc(names []string, values []float64) OptionFunc {
	return func(opt *ChartOption) {
		opt.RadarIndicators = NewRadarIndicators(names, values)
	}
}

// MarkLineOptionFunc sets the mark line for series of the chart.
func MarkLineOptionFunc(seriesIndex int, markLineTypes ...string) OptionFunc {
	return func(opt *ChartOption) {
		if len(opt.SeriesList) <= seriesIndex {
			return
		}
		opt.SeriesList[seriesIndex].MarkLine = NewMarkLine(markLineTypes...)
	}
}

// MarkPointOptionFunc sets the mark point for series of the chart.
func MarkPointOptionFunc(seriesIndex int, markPointTypes ...string) OptionFunc {
	return func(opt *ChartOption) {
		if len(opt.SeriesList) <= seriesIndex {
			return
		}
		opt.SeriesList[seriesIndex].MarkPoint = NewMarkPoint(markPointTypes...)
	}
}

func (o *ChartOption) fillDefault() error {
	o.Width = getDefaultInt(o.Width, defaultChartWidth)
	o.Height = getDefaultInt(o.Height, defaultChartHeight)

	yaxisCount := getSeriesYAxisCount(o.SeriesList)
	if yaxisCount < 0 {
		return errors.New("series specified invalid y-axis index")
	}
	if len(o.YAxis) < yaxisCount {
		yAxisOptions := make([]YAxisOption, yaxisCount)
		copy(yAxisOptions, o.YAxis)
		o.YAxis = yAxisOptions
	}

	if o.Font == nil {
		o.Font = GetDefaultFont()
	}
	if o.Theme == nil {
		o.Theme = GetDefaultTheme()
	}
	fillThemeDefaults(o.Theme, &o.Title, &o.Legend, &o.XAxis, o.YAxis)

	if o.Padding.IsZero() {
		o.Padding = defaultPadding
	}
	return nil
}

func fillThemeDefaults(defaultTheme ColorPalette, title *TitleOption, legend *LegendOption,
	xaxis *XAxisOption, yaxisOptions []YAxisOption) {
	if title.Theme == nil {
		title.Theme = defaultTheme
	}
	if legend.Theme == nil {
		legend.Theme = defaultTheme
	}
	if xaxis.Theme == nil {
		xaxis.Theme = defaultTheme
	}
	for i := range yaxisOptions {
		if yaxisOptions[i].Theme == nil {
			yaxisOptions[i].Theme = defaultTheme
		}
	}
}

// LineRender renders a line chart.
func LineRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListGeneric(values, ChartTypeLine),
	}, opts...)
}

// ScatterRender renders a scatter chart.
func ScatterRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListGeneric(values, ChartTypeScatter),
	}, opts...)
}

// BarRender renders a bar chart.
func BarRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListGeneric(values, ChartTypeBar),
	}, opts...)
}

// HorizontalBarRender renders a horizontal bar chart.
func HorizontalBarRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListGeneric(values, ChartTypeHorizontalBar),
	}, opts...)
}

// PieRender renders a pie chart.
func PieRender(values []float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListPie(values).ToGenericSeriesList(),
	}, opts...)
}

// DoughnutRender renders a doughnut or ring chart.
func DoughnutRender(values []float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListDoughnut(values).ToGenericSeriesList(),
	}, opts...)
}

// RadarRender renders a radar chart.
func RadarRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListGeneric(values, ChartTypeRadar),
	}, opts...)
}

// FunnelRender renders a funnel chart.
func FunnelRender(values []float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListFunnel(values).ToGenericSeriesList(),
	}, opts...)
}
